package main

import (
	"database/sql"
	"fmt"
	_ "modernc.org/sqlite"
)

func main() {
	db, err := sql.Open("sqlite", "")
	if err != nil {
		fmt.Println("Open err:", err)
		return
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT NOT NULL)")
	if err != nil {
		fmt.Println("Create err:", err)
		return
	}
	var name string
	err = db.QueryRow("SELECT name FROM users ORDER BY id LIMIT 1;").Scan(&name)
	fmt.Println("Query err:", err)
}
