package fileuploader

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"mime/multipart"

	"github.com/dhowden/tag"
	"github.com/zach/models"
)

func CreateSong(file multipart.File) (song models.SongRow) {

	m, er := tag.ReadFrom(file)
	if er != nil {
		log.Fatal(er)
	}
	artist := models.ArtistRow{Name: m.Artist()}
	album := models.AlbumRow{
		Title: m.Album(),
		Year:  m.Year()}
	filename, _ := tag.Sum(file)
	song = models.SongRow{
		Title:    m.Title(),
		FileType: strings.ToLower(string(m.FileType())),
		FileName: filename}

	a, _ := models.InsertArtist(artist)
	al, _ := models.InsertAlbum(album, a)
	models.InsertSong(song, al, a)
	file.Seek(0,0)
	return
}

func UploadFile(file multipart.File, song models.SongRow) {

	f, err := os.OpenFile("../restapi/music/"+song.FileName+".mp3", os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer f.Close()
	_, er := io.Copy(f, file)
	if er != nil {
		fmt.Println(er)
		return
	}
}
