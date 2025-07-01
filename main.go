/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package main

import (
"github.com/Mr-Aaryan/blockit/cmd"
"github.com/Mr-Aaryan/blockit/database"

)

func main() {
	database.InitDB()
	defer database.DB.Close()
	cmd.Execute()
}
