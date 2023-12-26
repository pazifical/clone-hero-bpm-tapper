package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

//go:embed static
var static embed.FS

var tapTimes []float64

var songDirectory = "songs"
var chartDirectory = "charts"

var songs []string

type SongInfo struct {
	Name   string `json:"name"`
	Artist string `json:"artist"`
}

func prepareSongDirectory() error {
	os.Mkdir(songDirectory, 0775)
	os.Mkdir(chartDirectory, 0775)

	files, err := os.ReadDir(songDirectory)
	if err != nil {
		return err
	}

	for _, file := range files {
		songs = append(songs, file.Name())
	}
	return nil
}

func main() {
	prepareSongDirectory()

	fileSystem, err := fs.Sub(static, "static")
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/", http.FileServer(http.FS(fileSystem)))
	http.Handle("/songs/", http.StripPrefix("/songs/", http.FileServer(http.Dir("./songs"))))
	http.HandleFunc("/api/songs", HandleSong)
	http.HandleFunc("/api/charts", HandleChart)
	http.HandleFunc("/api/taps", HandleTaps)

	log.Println("Serving")
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func HandleSong(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(songs)
	} else {
		w.WriteHeader(405)
	}
}

func HandleChart(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var songInfo SongInfo
		err := json.NewDecoder(r.Body).Decode(&songInfo)
		if err != nil {
			log.Println(err)
			w.WriteHeader(422)
			return
		}
		err = writeChartFile(songInfo)
		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
			return
		}
		w.WriteHeader(200)
	} else {
		w.WriteHeader(405)
	}
}

func writeChartFile(songInfo SongInfo) error {
	fileName := fmt.Sprintf("%s-%s.chart", songInfo.Artist, songInfo.Name)
	filePath := filepath.Join(chartDirectory, fileName)
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	f.WriteString(fmt.Sprintf(`[Song]
{
  Name = "%s"
  Artist = "%s"
}
[SyncTrack]
{
	
}`, songInfo.Name, songInfo.Artist))
	return nil
}

func HandleTaps(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tapTimes)
	} else if r.Method == "POST" {
		var times []float64
		defer r.Body.Close()
		err := json.NewDecoder(r.Body).Decode(&times)
		if err != nil {
			log.Println(err)
			w.WriteHeader(422)
			return
		}
		tapTimes = times
		w.WriteHeader(200)
		log.Println(tapTimes)
	} else {
		w.WriteHeader(405)
	}
}
