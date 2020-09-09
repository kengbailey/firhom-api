package main

import (
	"database/sql"

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

	database, _ := sql.Open("sqlite3", "./waverunner.db")
	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS waves (id INTEGER PRIMARY KEY, youtube_url NOT NULL TEXT, youtube_id NOT NULL TEXT, wave_name NOT NULL TEXT, insert_dt TEXT DEFAULT CURRENT_TIMESTAMP)")
	statement.Exec()

	statement, _ = database.Prepare("INSERT INTO waves (youtube_url, youtube_id, wave_name) VALUES (?, ?, ?)")
	statement.Exec("Nic", "Raboy")

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
youtube_url NOT NULL TEXT,
youtube_id NOT NULL TEXT,
wave_name NOT NULL TEXT,
insert_dt TEXT DEFAULT CURRENT_TIMESTAMP

*/
