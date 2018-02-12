package models

import (
	"database/sql"
	"log"
	"strconv"

	"github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"
	db "github.com/zach/database"
)

type AlbumRow struct {
	Id       int64
	Title    string
	Year     int
	ArtistId int
	Songs []SongAll
}

func InsertAlbum(album AlbumRow, artist ArtistRow) (newRow AlbumRow, err error) {
	if album.Title == "" {
		err = errors.Errorf("Can't create song without title (%s)",
			spew.Sdump(album))
		return
	}

	a, err := GetAlbumByTitle(album.Title)
	if a.Title == album.Title {
		newRow = a
		return
	} else {
	}

	const qry = `
      INSERT INTO albums (title, year, artistId)
      VALUES ($1, $2, $3)
      RETURNING id, title`

	err = db.DbConn.
		QueryRow(qry, album.Title, album.Year, artist.Id).
		Scan(&newRow.Id, &newRow.Title)
	if err != nil {
		err = errors.Wrapf(err,
			"Couldn't insert user album into DB (%s)",
			spew.Sdump(album))
		return
	}
	return
}

func GetAlbumById(id string) (row AlbumRow, err error) {
	if id == "" {
		err = errors.Errorf("Can't find album without id (%s)",
			spew.Sdump(id))
		return
	}

	i , _ := strconv.Atoi(id)
	const qry = `
		SELECT * FROM albums
		WHERE id = $1`

	err = db.DbConn.QueryRow(qry, i).Scan(&row.Id, &row.Title, &row.Year, &row.ArtistId)
	switch {
	case err == sql.ErrNoRows:
		log.Printf("None by that id.")
	case err != nil:
		log.Fatal(err)
	}
	return
}

func GetAlbumsByArtistId(id string)(rows []AlbumRow) {
	if id == "" {
		log.Print("Can't get song rows with empty type")
		return
	}
	i, _ := strconv.Atoi(id)

	const qry = `
    SELECT id, title, year, artistId
		FROM albums
		WHERE artistId = $1`

	iterator, err := db.DbConn.Query(qry, i)
	if err != nil {
		log.Print("Album listing failed")
		log.Print(err)
		return
	}

	defer iterator.Close()

	for iterator.Next() {
		var row = AlbumRow{}

		err = iterator.Scan(&row.Id, &row.Title, &row.Year, &row.ArtistId)
		if err != nil {
			log.Print("Album row scanning failed")
			log.Print(err)
			return
		}

		rows = append(rows, row)
	}

	if err = iterator.Err(); err != nil {
		log.Print("Errored while looping through Albums listing")
		log.Print(err)
		return
	}
	log.Printf("Got Albums")
	return
}

func GetAlbumByTitle(title string) (row AlbumRow, err error) {
	if title == "" {
		err = errors.Errorf("Can't find album without title (%s)",
			spew.Sdump(title))
		return
	}

	const qry = `
		SELECT id, title FROM albums
		WHERE title = $1`

	err = db.DbConn.QueryRow(qry, title).Scan(&row.Id, &row.Title)
	switch {
	case err == sql.ErrNoRows:
		log.Printf("None by that title.")
	case err != nil:
		log.Fatal(err)
	}
	return
}

func GetAlbumsByArtist() {}

func GetAlbums() (albums []AlbumRow) {
	const qry = `
		SELECT id, title, year, artistId
		FROM albums`

	iterator, err := db.DbConn.Query(qry)
	if err != nil {
		log.Print(err)
		return
	}
	defer iterator.Close()
	for iterator.Next() {
		var album = AlbumRow{}
		err = iterator.Scan(&album.Id, &album.Title, &album.Year, &album.ArtistId)
		if err != nil {
			log.Print(err)
			return
		}

		albums = append(albums, album)
	}
	return
}
