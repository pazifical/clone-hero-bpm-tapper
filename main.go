package main

import (
	"bytes"
	"clone-hero-bpm-tapper/internal"
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

//go:embed static
var static embed.FS

var tapTimes []float64

var songDirectory = "songs"
var chartDirectory = "charts"

var songs []string

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
	fmt.Println(fileSystem)

	// http.Handle("/", http.FileServer(http.FS(fileSystem)))
	http.Handle("/", http.FileServer(http.Dir("static")))
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
	} else if r.Method == "POST" {
		fmt.Println("SONG POST")
		err := r.ParseMultipartForm(32 << 20)
		if err != nil {
			log.Println(err)
			w.WriteHeader(422)
			return
		}

		var buf bytes.Buffer

		file, header, err := r.FormFile("file")
		if err != nil {
			log.Println(err)
			w.WriteHeader(422)
			return
		}
		defer file.Close()
		name := strings.Split(header.Filename, ".")
		fmt.Printf("File name %s\n", name[0])

		_, err = io.Copy(&buf, file)
		if err != nil {
			log.Println(err)
			w.WriteHeader(422)
			return
		}

		f, err := os.Create(filepath.Join(songDirectory, header.Filename))
		if err != nil {
			log.Println(err)
			w.WriteHeader(422)
			return
		}
		defer f.Close()
		f.Write(buf.Bytes())
		buf.Reset()

		w.WriteHeader(200)
	} else {
		w.WriteHeader(405)
	}
}

func HandleChart(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		err := r.ParseMultipartForm(100000)
		if err != nil {
			log.Println(err)
			w.WriteHeader(422)
			return
		}
		beatsPerBar, err := strconv.Atoi(r.FormValue("beats_per_bar"))
		if err != nil {
			log.Println(err)
			w.WriteHeader(422)
			return
		}
		abc, err := strconv.Atoi(r.FormValue("average_beat_count"))
		if err != nil {
			log.Println(err)
			w.WriteHeader(422)
			return
		}

		chartInfo := internal.ChartInfo{
			Name:             r.FormValue("name"),
			Artist:           r.FormValue("artist"),
			BeatsPerBar:      beatsPerBar,
			AverageBeatCount: abc,
		}

		err = writeChartFile(chartInfo)
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

func writeChartFile(chartInfo internal.ChartInfo) error {
	fileName := fmt.Sprintf("%s-%s.chart", chartInfo.Artist, chartInfo.Name)
	filePath := filepath.Join(chartDirectory, fileName)
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	bpmParts := internal.CalculateBPMParts(tapTimes, chartInfo)
	var builder strings.Builder
	for _, bpmPart := range bpmParts {
		builder.WriteString(fmt.Sprintf("  %d = B %d\n", bpmPart.Position, int64(1000*bpmPart.BPM)))
	}
	bpmPartsLines := builder.String()
	fmt.Println(bpmPartsLines)

	f.WriteString(fmt.Sprintf(`[Song]
{
  Name = "%s"
  Artist = "%s"
}
[SyncTrack]
{
%s}`, chartInfo.Name, chartInfo.Artist, bpmPartsLines))
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
		log.Printf("INFO: received taps: %v", tapTimes)
		log.Println(tapTimes)
	} else {
		w.WriteHeader(405)
	}
}
