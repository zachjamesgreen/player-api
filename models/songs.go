package models

import (
	"database/sql"
	"log"
	"strconv"
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"
	db "github.com/zach/database"
)

// SongRow represents in a `struct` the information we
// can get from the table (some fields are insertable but
// not all - ID and CreatedAt are generated when we `insert`,
// thus, these can only be retrieved).
type SongRow struct {
	Id       int64
	Title    string
	FileType string
	FilePath string
	FileName string
	AlbumId  int
	ArtistId int
}

type SongAll struct {
	Id         int64
	Title      string
	FileType   string
	FilePath   string
	FileName   string
	AlbumTitle string
	ArtistName string
	AlbumId  int
	ArtistId int
}

func InsertSong(song SongRow, album AlbumRow, artist ArtistRow) (newRow SongRow, err error) {
	if song.Title == "" {
		err = errors.Errorf("Can't create song without Type (%s)",
			spew.Sdump(song))
		return
	}

	s := GetSongByFileName(song.FileName)
	if s.FileName == song.FileName {
		newRow = s
	} else {
		const qry = `
	    INSERT INTO songs (title, file_type, file_name, albumId, artistId)
	    VALUES ($1, $2, $3, $4, $5)
	    RETURNING id, title`

		err = db.DbConn.
			QueryRow(qry, song.Title, song.FileType, song.FileName, album.Id, artist.Id).
			Scan(&newRow.Id, &newRow.Title)
		if err != nil {
			err = errors.Wrapf(err,
				"Couldn't insert user song into DB (%s)",
				spew.Sdump(song))
			return
		}

		return
	}
	return
}

func GetSongsByTitle(songTitle string) (rows []SongRow, err error) {
	if songTitle == "" {
		err = errors.Errorf("Can't get song rows with empty type")
		return
	}

	const qry = `
    SELECT
    	id, title
    FROM
    	songs
    WHERE
    	title = $1%`

	iterator, err := db.DbConn.Query(qry, songTitle)
	if err != nil {
		err = errors.Wrapf(err,
			"Song listing failed (type=%s)",
			songTitle)
		return
	}

	defer iterator.Close()

	for iterator.Next() {
		var row = SongRow{}

		err = iterator.Scan(&row.Id, &row.Title)
		if err != nil {
			err = errors.Wrapf(err,
				"Song row scanning failed (type=%s)",
				songTitle)
			return
		}

		rows = append(rows, row)
	}

	if err = iterator.Err(); err != nil {
		err = errors.Wrapf(err,
			"Errored while looping through songs listing (type=%s)",
			songTitle)
		return
	}

	return
}


func GetSongsByArtistId(id string) (rows []SongAll){
	if id == "" {
		log.Print("Can't get song rows with empty type")
		return
	}
	i, _ := strconv.Atoi(id)

	const qry = `
    SELECT songs.id, songs.title, songs.file_type, songs.file_name, albums.title as album_title, artists.name, songs.artistId, songs.albumId
		FROM songs
		inner join albums on songs.albumId = albums.id
		inner join artists on songs.artistId = artists.id
		WHERE songs.artistId = $1`

	iterator, err := db.DbConn.Query(qry, i)
	if err != nil {
		log.Print("Song listing failed")
		log.Print(err)
		return
	}

	defer iterator.Close()

	for iterator.Next() {
		var row = SongAll{}

		err = iterator.Scan(&row.Id, &row.Title, &row.FileType, &row.FileName, &row.AlbumTitle, &row.ArtistName, &row.ArtistId, &row.AlbumId)
		if err != nil {
			log.Print("Song row scanning failed")
			log.Print(err)
			return
		}

		rows = append(rows, row)
	}

	if err = iterator.Err(); err != nil {
		log.Print("Errored while looping through songs listing")
		log.Print(err)
		return
	}
	return
}

func GetSongsByAlbumId(id string) (rows []SongAll){
	if id == "" {
		log.Print("Can't get song rows with empty type")
		return
	}
	i, _ := strconv.Atoi(id)

	const qry = `
    SELECT songs.id, songs.title, songs.file_type, songs.file_name, albums.title as album_title, artists.name, songs.artistId, songs.albumId
		FROM songs
		inner join albums on songs.albumId = albums.id
		inner join artists on songs.artistId = artists.id
		WHERE albumId = $1`

	iterator, err := db.DbConn.Query(qry, i)
	if err != nil {
		log.Print("Song listing failed")
		log.Print(err)
		return
	}

	defer iterator.Close()

	for iterator.Next() {
		var row = SongAll{}

		err = iterator.Scan(&row.Id, &row.Title, &row.FileType, &row.FileName, &row.AlbumTitle, &row.ArtistName, &row.ArtistId, &row.AlbumId)
		if err != nil {
			log.Print("Song row scanning failed")
			log.Print(err)
			return
		}

		rows = append(rows, row)
	}

	if err = iterator.Err(); err != nil {
		log.Print("Errored while looping through songs listing")
		log.Print(err)
		return
	}
	return
}

func GetSongByFileName(filename string) (song SongRow) {
	const qry = `
		SELECT id, title, file_name, artistId, albumId
		FROM songs
		WHERE file_name = $1`

	err := db.DbConn.QueryRow(qry, filename).Scan(&song.Id, &song.Title, &song.FileName, &song.AlbumId, &song.ArtistId)
	switch {
	case err == sql.ErrNoRows:
		log.Printf("None by that name.")
	case err != nil:
		log.Fatal(err)
	}
	return
}

func GetSongs() (songs []SongRow) {
	const qry = `
		SELECT id, title, file_type, file_name, artistId, albumId
		FROM songs`

	iterator, err := db.DbConn.Query(qry)
	if err != nil {
		log.Print(err)
		return
	}
	defer iterator.Close()
	for iterator.Next() {
		var song = SongRow{}
		err = iterator.Scan(&song.Id, &song.Title, &song.FileType, &song.FileName, &song.ArtistId, &song.AlbumId)
		if err != nil {
			log.Print(err)
			return
		}

		songs = append(songs, song)
	}
	return
}

func GetSongsAll() (songs []SongAll) {
	const qry = `
		SELECT songs.id, songs.title, songs.file_type, songs.file_name, songs.albumId, songs.artistId, albums.title as album_title, artists.name
		from songs
		inner join albums on songs.albumId = albums.id
		inner join artists on songs.artistId = artists.id;`

	iterator, err := db.DbConn.Query(qry)
	if err != nil {
		log.Print(err)
		return
	}
	defer iterator.Close()
	for iterator.Next() {
		var song = SongAll{}
		err = iterator.Scan(&song.Id, &song.Title, &song.FileType, &song.FileName, &song.AlbumId, &song.ArtistId, &song.AlbumTitle, &song.ArtistName)
		if err != nil {
			log.Print(err)
			return
		}

		songs = append(songs, song)
	}
	return
}

func SongsSearch(term string) (songs []SongAll) {
	qry := fmt.Sprint(`
		SELECT songs.id, songs.title, songs.file_type, songs.file_name, songs.albumId, songs.artistId, albums.title as album_title, artists.name
		from songs
		inner join albums on songs.albumId = albums.id
		inner join artists on songs.artistId = artists.id
		WHERE songs.title ILIKE '`, term, "%'")

	iterator, err := db.DbConn.Query(qry)
	if err != nil {
 		log.Print(err)
		return
	}
	defer iterator.Close()
	for iterator.Next() {
		var song = SongAll{}
		err = iterator.Scan(&song.Id, &song.Title, &song.FileType, &song.FileName, &song.AlbumId, &song.ArtistId, &song.AlbumTitle, &song.ArtistName)
		if err != nil {
			log.Print(err)
			return
		}

		songs = append(songs, song)
	}
	return
}
