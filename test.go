package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// url := "https://www.youtube.com/watch?v=lmku-exLxxw"
	// cmd := exec.Command(`youtube-dl`, `--get-filename`, `-o`, `%(id)s|%(title)s`, url)
	// cmdOut, err := cmd.Output()
	// if err != nil {
	// 	log.Printf("Failed to execute youtube-dl command: %s\n", err.Error())
	// }
	// vidInfo := strings.Split(string(cmdOut), "|")
	// fmt.Println(vidInfo[0], vidInfo[1])

	database, err := sql.Open("sqlite3", "./waverunner.db")
	if err != nil {
		log.Fatalf("Failed to open sqlite database connection: %v\n", err)
	}
	// statement, err := database.Prepare("CREATE TABLE IF NOT EXISTS waves (id INTEGER PRIMARY KEY, youtube_url TEXT NOT NULL, youtube_id TEXT NOT NULL, wave_name TEXT NOT NULL, insert_dt DATETIME DEFAULT CURRENT_TIMESTAMP);")
	// if err != nil {
	// 	log.Fatalf("Failed to execute sqlite statement: %v", err)
	// }
	// statement.Exec()

	// statement, err := database.Prepare("INSERT INTO waves (youtube_url, youtube_id, wave_name) VALUES (?, ?, ?)")
	// if err != nil {
	// 	log.Fatalf("Failed to insert new item into the database: %v", err)
	// }
	// statement.Exec("one", "two", "three")

	rows, err := database.Query(fmt.Sprintf("Select youtube_id from waves where id = %d;", 1))
	if err != nil {
		log.Fatalf("Failed to select item from database: %v", err)
	}
	var youtubeURL string
	for rows.Next() {
		rows.Scan(&youtubeURL)
		fmt.Println(youtubeURL)
	}
	// rows, _ := database.Query("SELECT id, firstname, lastname FROM people")
	// var id int
	// var firstname string
	// var lastname string
	// for rows.Next() {
	// 	rows.Scan(&id, &firstname, &lastname)
	// 	fmt.Println(strconv.Itoa(id) + ": " + firstname + " " + lastname)
	// }
}

/*
What do we need for this db?

Table Name --> Waves

id INTEGER PRIMARY KEY AUTOINCREMENT,
youtube_url TEXT NOT NULL,
youtube_id TEXT NOT NULL,
wave_name TEXT NOT NULL,
insert_dt DATETIME DEFAULT CURRENT_TIMESTAMP

*/
