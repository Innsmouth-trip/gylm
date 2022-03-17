package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/joho/godotenv"
	"github.com/shkh/lastfm-go/lastfm"
)

var (
	trackName  string
	artistName string
	imageData  string
	userUrl    string
)

type TrackData struct {
	TrackName  string
	ArtistName string
	Image      string
	UserUrl    string
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("No .env file found")
	}
}

func main() {
	key := os.Getenv("KEY")
	secret := os.Getenv("SECRET")
	user := os.Getenv("USERNAME")
	limit := os.Getenv("LIMIT")
	userUrl = os.Getenv("YANDEX_URL")

	http.HandleFunc("/yandex", func(w http.ResponseWriter, r *http.Request) {

		api := lastfm.New(key, secret)

		result, err := api.User.GetRecentTracks(lastfm.P{"limit": limit,
			"user": user})
		if err != nil {
			return
		}

		for _, track := range result.Tracks {
			if track.NowPlaying == "true" {
				trackName = track.Name
				artistName = track.Artist.Name
				imageData = track.Images[3].Url

				fmt.Println(track.Artist.Name, track.Name, "now: ", track.NowPlaying, track.Images[3].Url)
				break
			} else {
				rand.Seed(time.Now().UnixNano())
				track := result.Tracks[rand.Intn(len(result.Tracks))]

				trackName = track.Name
				artistName = track.Artist.Name
				imageData = track.Images[3].Url

				fmt.Println(track.Artist.Name, track.Name)
				break
			}
		}

		data := TrackData{
			TrackName:  trackName,
			ArtistName: artistName,
			Image:      imageData,
			UserUrl:    userUrl,
		}

		tmpl, err := template.ParseFiles("templates/index.html")
		if err != nil {
			fmt.Printf("template execution: %s", err)
		}

		w.Header().Set("Content-Type", "image/svg+xml")
		w.Header().Set("Cache-Control", "max-age=0")
		if err := tmpl.Execute(w, data); err != nil {
			log.Printf("Error executing template: %v", err)
		}

	})
	//server log
	err := http.ListenAndServe(":1984", nil)
	if err != nil {
		fmt.Println("ListenAndServe: ", err)
	}

}
