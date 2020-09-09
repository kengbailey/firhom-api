package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

// processWavHandler ... This function kicks off the downloading of a youtube video to wav file.
// It does not wait for the download to finish.
// Instead it checks for the successful start of the process and returns a success to the application.
// Video details are logged to a sqlite3 database.
func processWavHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request!", http.StatusNotFound)
		return
	}

	// get url
	videoURL := r.FormValue("url")

	// extract youtube id + name
	cmd := exec.Command(`youtube-dl`, `--get-filename`, `-o`, `%(id)s|%(title)s`, videoURL)
	cmdOut, err := cmd.Output()
	if err != nil {
		log.Printf("Failed to execute youtube-dl command: %v\n", err)
		http.Error(w, "Failed to download WAV!", http.StatusBadRequest)
		return
	}
	vidInfo := strings.Split(string(cmdOut), "|")

	// create db record
	// TODO: Check for existing record before we insert new record
	statement, err := db.Prepare("INSERT INTO waves (youtube_url, youtube_id, wave_name) VALUES (?, ?, ?);")
	if err != nil {
		log.Printf("Failed to insert new record into sqlite3 database: %v\n", err)
		http.Error(w, "Failed to download WAV!", http.StatusBadRequest)
		return
	}
	statement.Exec(videoURL, vidInfo[0], vidInfo[1])

	// kick off ytdl to wav --> youtube-dl -f bestaudio --extract-audio --audio-format wav --audio-quality 0 --output "./wavs/%(id)s.%(ext)s" <Video-URL>
	// We don't wait for this to finish
	cmd = exec.Command(`youtube-dl`,
		`-f`, `bestaudio`,
		`--extract-audio`,
		`--audio-format`, `wav`,
		`--audio-quality`, `0`,
		`--output`, `./wavs/%(id)s.%(ext)s`,
		videoURL)
	err = cmd.Start()
	if err != nil {
		log.Printf("Failed to execute youtube-dl command: %v\n", err)
		http.Error(w, "Failed to download WAV!", http.StatusBadRequest)
		return
	}

	return
}

func downloadWavHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request!", http.StatusNotFound)
		return
	}

	// get wav id
	wavID := r.URL.Query()["id"][0]

	// get record from db
	var fileID string
	err := db.QueryRow("select youtube_id from waves where id=?", wavID).Scan(&fileID)
	if err != nil {
		log.Printf("Failed to retrieve wav from sqlite database: %v", err)
		http.Error(w, "Failed to find WAV!", http.StatusNotFound)
		return
	}

	// download file to user
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment; filename="+fileID+".wav")
	// todo: content length
	diskFile, err := os.Open(fmt.Sprintf("./wavs/%s.wav", fileID))
	if err != nil {
		log.Printf("Failed to open wav file on disk: %v", err)
		http.Error(w, "Failed to open file!", http.StatusNotFound)
		return
	}
	defer diskFile.Close()
	io.Copy(w, diskFile)

	return
}

// wavListItem ... returned to UI upon page load.
type wavListItem struct {
	ID     int
	Name   string
	URL    string
	DlDate string
}

func getWavListHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request!", http.StatusNotFound)
		return
	}

	// fetch from db
	var WavList []wavListItem
	rows, err := db.Query("Select id, wav_name, youtube_url, insert_dt from waves;")
	if err != nil {
		log.Println("Failed to retrieve list of waves: %v", err)
		http.Error(w, "Failed to retrieve wav list!", http.StatusNotFound)
		return
	}
	for rows.Next() {
		var tempItem wavListItem
		rows.Scan(&tempItem.ID, &tempItem.Name, &tempItem.URL, &tempItem.DlDate)
		WavList = append(WavList, tempItem)
	}

	// return json marshalled list

	return
}

var db *sql.DB

// initDb starts the database connection for the application.
// Initializes the global db variable.
func initDb() {
	db, err := sql.Open("sqlite3", "./waverunner.db")
	if err != nil {
		log.Fatal("Failed to open sqlite3 database: %v", err)
	}
	statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS waves (id INTEGER PRIMARY KEY, youtube_url NOT NULL TEXT, youtube_id NOT NULL TEXT, wave_name NOT NULL TEXT, insert_dt TEXT DEFAULT CURRENT_TIMESTAMP)")
	if err != nil {
		log.Fatal("Failed to initialize waves table: %v", err)
	}
	statement.Exec()
	statement.Close()
}

func main() {

	// init db
	initDb()

	// endpoints
	http.HandleFunc("/processwav", processWavHandler)
	http.HandleFunc("/downloadwav", downloadWavHandler)
	http.HandleFunc("/getwavlist", getWavListHandler)

	// serve
	log.Println("Listening on localhost:3000...")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal(err)
	}
}
