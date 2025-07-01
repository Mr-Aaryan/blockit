package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Block struct {
	Id      int
	Title   string
	Blocked bool
	BlockId string
}

var DB *sql.DB
var init_db_statement string = `
	CREATE TABLE IF NOT EXISTS blocks (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		title TEXT,
		blocked INTEGER,
		block_id TEXT
	);`
var read_db_statement string = "SELECT * from blocks"
var read_blocked_statement string = "SELECT title from blocks WHERE blocked=1"
var select_block string = `UPDATE blocks SET blocked = 1 WHERE title = ?`
var unselect_block string = `UPDATE blocks SET blocked = 0 WHERE title = ?`
var add_into_blocklist_statement = `INSERT INTO blocks(title, block_id, blocked) VALUES(?,?, 0)`

func InitDB() {
	var err error
	dbPath := "/etc/blocklist/data.db"
	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	_, err = DB.Exec(init_db_statement)
	if err != nil {
		log.Fatalf("Error creating table: %q: %s\n", err, init_db_statement)
	}
}

func ReadDB() []Block {
	rows, err := DB.Query(read_db_statement)
	if err != nil {
		log.Fatal("Error getting data from database")
	}
	defer rows.Close()

	blocks := []Block{}
	for rows.Next() {
		var block Block
		if err := rows.Scan(&block.Id, &block.Title, &block.Blocked, &block.BlockId); err != nil {
			log.Fatal("Error reading next row from db", err)
		}
		blocks = append(blocks, block)
	}
	return blocks
}

func GetBlockedList() ([]string, error) {
	rows, err := DB.Query(read_blocked_statement)
	if err != nil {
		return nil, fmt.Errorf("error getting blocked list from database: %w", err)
	}
	defer rows.Close()

	var blockedItems []string
	for rows.Next() {
		var title string
		if err := rows.Scan(&title); err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		blockedItems = append(blockedItems, title)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return blockedItems, nil
}

func SelectBlocked(title string) {
	if _, err := DB.Exec(select_block, title); err != nil {
		log.Fatalf("Error adding into blocklist: %v", err)
	}
}

func UnselectBlocked(title string) {
	if _, err := DB.Exec(unselect_block, title); err != nil {
		log.Fatalf("Error removing from blocklist: %v", err)
	}
}

func InsertIntoDBLoop() {
	websites := []struct {
		name  string
		lines string
	}{
		{"reddit", "0,1"},
		{"linkedin", "2,3"},
		{"twitter", "4,5,6,7"},
		{"facebook", "8,9"},
		{"netflix", "10,11"},
		{"twitch", "12,13"},
		{"youtube", "14,15"},
		{"instagram", "16,17"},
		{"tiktok", "18,19"},
		{"pinterest", "20,21"},
		{"snapchat", "22,23"},
		{"discord", "24,25"},
		{"roblox", "26,27"},
		{"steamcommunity", "28,29"},
		{"amazon", "30,31"},
		{"ebay", "32,33"},
		{"aliexpress", "34,35"},
	}

	for _, website := range websites {
		if _, err := DB.Exec(add_into_blocklist_statement, website.name, website.lines); err != nil {
			log.Fatalf("Error inserting %s into blocklist: %v", website.name, err)
		}
	}

	log.Println("Successfully inserted all websites into blocklist database")
}