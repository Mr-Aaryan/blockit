package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/Mr-Aaryan/blockit/database"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

// helper func
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func handleFile(lineNumbers []int, selec int) error {
	fp := "/etc/hosts"

	content, err := os.ReadFile(fp)
	if err != nil {
		return fmt.Errorf("failed to read file '%s': %w", fp, err)
	}

	lines := strings.Split(string(content), "\n")

	for _, lineIndex := range lineNumbers {
		if lineIndex < 0 || lineIndex >= len(lines) {
			log.Printf("Warning: Line index %d is out of range (file has %d lines). Skipping this line.", lineIndex, len(lines))
			continue
		}

		targetLine := lines[lineIndex]

		if strings.TrimSpace(targetLine) != "" {
			switch selec {
			case 1: // selec = 1 means to BLOCK (comment out)
				if !strings.HasPrefix(targetLine, "# ") {
					lines[lineIndex] = "# " + targetLine
				}
			case 0: // selec = 0 means to UNBLOCK (uncomment)
				if strings.HasPrefix(targetLine, "# ") {
					lines[lineIndex] = strings.TrimPrefix(targetLine, "# ")
				}
			}
		}
	}

	newContent := strings.Join(lines, "\n")

	if err := os.WriteFile(fp, []byte(newContent), 0644); err != nil {
		return fmt.Errorf("failed to write to file '%s': %w", fp, err)
	}

	return nil
}

var rootCmd = &cobra.Command{
	Use:   "block_site",
	Short: "Block or unblock websites by modifying the /etc/hosts file.",
	Long: `This application allows you to easily block or unblock websites
by managing entries in your /etc/hosts file. It uses a database to store
website configurations and their corresponding line numbers for efficient updates.

Note: This tool modifies system files and requires root privileges to run.
Please execute it using 'sudo', for example: 'sudo block_site'.`,

	Run: func(cmd *cobra.Command, args []string) {
		// database.InsertIntoDBLoop()
		allBlocks := database.ReadDB()
		if len(allBlocks) == 0 {
			fmt.Println("No websites found in the database. Nothing to select.")
			return
		}

		// Pre-allocate slices with known capacity to reduce memory allocations
		blockMap := make(map[string][]int, len(allBlocks))
		allSites := make([]string, 0, len(allBlocks))

		for _, item := range allBlocks {
			allSites = append(allSites, item.Title)
			parts := strings.Split(item.BlockId, ",")
			lineNumbers := make([]int, 0, len(parts)) // Pre-allocate with capacity

			for _, part := range parts {
				if num, err := strconv.Atoi(strings.TrimSpace(part)); err == nil {
					lineNumbers = append(lineNumbers, num)
				} else {
					log.Printf("Warning: Could not convert '%s' to integer for site '%s'. Skipping.", part, item.Title)
				}
			}
			blockMap[item.Title] = lineNumbers
		}

		blockedSites, err := database.GetBlockedList()
		if err != nil {
			log.Fatal("Error loading currently blocked items from database: ", err)
		}

		var selectedSites []string
		// Set default selected sites to currently blocked sites
		selectedSites = append(selectedSites, blockedSites...)

		// Create options for Huh MultiSelect
		options := make([]huh.Option[string], len(allSites))
		for i, site := range allSites {
			options[i] = huh.NewOption(site, site)
		}

		form := huh.NewForm(
			huh.NewGroup(
				huh.NewMultiSelect[string]().
					Title("Select websites to block").
					Description("Use arrow keys to navigate, space to select/deselect, enter to confirm").
					Options(options...).
					Value(&selectedSites).
					Limit(len(allSites)).
					Height(15),
			),
		)

		err = form.Run()
		if err != nil {
			log.Fatal("Error during form prompt: ", err)
		}

		changesApplied := false

		for _, title := range allSites {
			if contains(selectedSites, title) {
				database.SelectBlocked(title)
				if lineNumbers, exists := blockMap[title]; exists {
					if err := handleFile(lineNumbers, 0); err != nil { // selec=1 to UNCOMMENT (BLOCK)
						log.Printf("Error blocking '%s' in /etc/hosts: %v", title, err)
					} else {
						changesApplied = true
					}
				} else {
					log.Printf("Warning: No line numbers found for site '%s' in blockMap. Cannot modify /etc/hosts.", title)
				}
			} else {
				database.UnselectBlocked(title)
				if lineNumbers, exists := blockMap[title]; exists {
					if err := handleFile(lineNumbers, 1); err != nil { // selec=0 to COMMENT (UNBLOCK)
						log.Printf("Error unblocking '%s' in /etc/hosts: %v", title, err)
					} else {
						changesApplied = true
					}
				} else {
					log.Printf("Warning: No line numbers found for site '%s' in blockMap. Cannot modify /etc/hosts.", title)
				}
			}
		}

		if changesApplied {
			fmt.Println("Restarting NetworkManager to apply changes...")
			restartCmd := exec.Command("systemctl", "restart", "NetworkManager")
			output, err := restartCmd.CombinedOutput()
			if err != nil {
				fmt.Printf("Failed to restart NetworkManager: %v\n", err)
				if len(output) > 0 {
					fmt.Printf("Output: %s\n", string(output))
				}
			} else {
				fmt.Println("NetworkManager restarted successfully.")
			}
		} else {
			fmt.Println("No changes were applied to /etc/hosts file.")
		}

		fmt.Printf("\nOperation completed. Selected sites to be blocked: %v\n", selectedSites)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {}
