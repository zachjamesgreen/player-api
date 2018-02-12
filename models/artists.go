package models

import (
	"database/sql"
	"log"
	"strconv"

	"github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"
	db "github.com/zach/database"
)

type ArtistRow struct {
	Id   int64
	Name string
	Songs []SongAll
	Albums []AlbumRow
}

func InsertArtist(row ArtistRow) (newRow ArtistRow, err error) {
	if row.Name == "" {
		err = errors.Errorf("Can't create song without Type (%s)",
			spew.Sdump(row))
		return
	}
	a, err := GetArtistByName(row.Name)
	if a.Name == row.Name {
		newRow = a
		return
	} else {
		const qry = `
	      INSERT INTO artists (name)
	      VALUES ($1)
	      RETURNING id, name`

		err = db.DbConn.
			QueryRow(qry, row.Name).
			Scan(&newRow.Id, &newRow.Name)
		if err != nil {
			err = errors.Wrapf(err,
				"Couldn't insert user row into DB (%s)",
				spew.Sdump(row))
			return
		}
	}
	return
}

func GetArtistById(id string) (row ArtistRow) {
	if id == "" {
		log.Printf("Can't find artist without id (%s)", id)
		return
	}
	i, _ := strconv.Atoi(id)

	const qry = "SELECT * FROM artists WHERE id = $1 LIMIT 1"
	log.Print(qry)

	err := db.DbConn.QueryRow(qry, i).Scan(&row.Id, &row.Name)
	if err != nil {
		log.Printf("Couldn't find artist by id (%d)", i)
		return
	}
	return
}

func GetArtistByName(name string) (row ArtistRow, err error) {
	if name == "" {
		err = errors.Errorf("Can't find artist without name (%s)",
			spew.Sdump(name))
		return
	}

	const qry = `
    SELECT * FROM artists
    WHERE name = $1`

	err = db.DbConn.QueryRow(qry, name).Scan(&row.Id, &row.Name)
	switch {
	case err == sql.ErrNoRows:
		log.Printf("None by that name.")
	case err != nil:
		log.Fatal(err)
	}
	return
}

func GetArtists() (artists []ArtistRow) {
	const qry = `
		SELECT id, name
		FROM artists`

	iterator, err := db.DbConn.Query(qry)
	if err != nil {
		log.Print(err)
		return
	}
	defer iterator.Close()
	for iterator.Next() {
		var artist = ArtistRow{}
		err = iterator.Scan(&artist.Id, &artist.Name)
		if err != nil {
			log.Print(err)
			return
		}

		artists = append(artists, artist)
	}
	return
}
