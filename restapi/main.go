package main

import (
	"os"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	db "github.com/zach/database"
	"github.com/zach/fileuploader"
	"github.com/zach/models"
)

func upload(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	r.ParseMultipartForm(32 << 20)
	m := r.MultipartForm
	files := m.File["uploadfile"]
	for i, _ := range files {
		file, _ := files[i].Open()
		defer file.Close()
		song := fileuploader.CreateSong(file)
		fileuploader.UploadFile(file, song)
	}
	w.WriteHeader(http.StatusOK)
	// song := fileuploader.CreateSong(w, r)
	// fileuploader.UploadFile(w, r, song)
}

func songs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	// songs := models.GetSongs()
	songs := models.GetSongsAll()
	js, err := json.Marshal(songs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(js)
}
func artists(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	artists := models.GetArtists()
	js, err := json.Marshal(artists)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(js)
}
func albums(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	albums := models.GetAlbums()
	js, err := json.Marshal(albums)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(js)
}

func albumById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	album, _ := models.GetAlbumById(vars["id"])
	album.Songs = models.GetSongsByAlbumId(vars["id"])
	js, err := json.Marshal(album)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(js)
}

func artistById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	artist := models.GetArtistById(vars["id"])
	// TODO: loop through albums to add songs
	artist.Albums= models.GetAlbumsByArtistId(vars["id"])
	artist.Songs = models.GetSongsByArtistId(vars["id"])
	js, err := json.Marshal(artist)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(js)
}

func main() {
	var err error
	cfg := db.Config{
		Host:     "localhost",
		Port:     "5432",
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Database: os.Getenv("DB_DATABASE")}
	db.DbConn, err = db.New(cfg)
	if err != nil {
		log.Print(err)
	}
	defer db.DbConn.Close()
	r := mux.NewRouter()
	r.HandleFunc("/songs", songs).Methods("GET")
	r.HandleFunc("/artists", artists).Methods("GET")
	r.HandleFunc("/artist/{id}", artistById).Methods("GET")
	r.HandleFunc("/albums", albums).Methods("GET")
	r.HandleFunc("/album/{id}", albumById).Methods("GET")
	r.HandleFunc("/upload", upload).Methods("POST")
	r.PathPrefix("/music/").Handler(http.StripPrefix("/music/", http.FileServer(http.Dir("music"))))
	log.Fatal(http.ListenAndServe(":8000", r))
}
