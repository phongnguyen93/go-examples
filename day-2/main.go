package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
)

type Album struct {
	ID     int64
	Title  string
	Artist string
	Price  float32
}

var db *sql.DB

func getAlbumsByArtist(name string) ([]Album, error) {
	var albums []Album
	rows, err := db.Query("SELECT * FROM album WHERE artist = ?", name)
	if err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
	}
	defer rows.Close()
	for rows.Next() {
		var album Album
		if err := rows.Scan(&album.ID, &album.Title, &album.Artist, &album.Price); err != nil {
			return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
		}
		albums = append(albums, album)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
	}
	return albums, nil
}

func getAlbumByID(id int64) (Album, error) {
	var album Album
	row := db.QueryRow("SELECT * FROM album WHERE id = ?", id)
	if err := row.Scan(&album.ID, &album.Title, &album.Artist, &album.Price); err != nil {
		if err == sql.ErrNoRows {
			return album, fmt.Errorf("albumByID %d no such album", id)
		}
		return album, fmt.Errorf("albumByID %d: %v", id, err)
	}
	return album, nil
}

func insertNewAlbum(album Album) (int64, error) {
	result, err := db.Exec("INSERT INTO album (title, artist, price) VALUES (?,?,?)", album.Title, album.Artist, album.Price)
	if err != nil {
		return 0, fmt.Errorf("insertNewAlbum %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("insertNewAlbum %v", err)
	}
	return id, nil
}

func connectDB() error {
	cfg := mysql.Config{
		User:   "root",
		Passwd: "9232780a",
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "recordings",
	}
	// Get a database handle.
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return fmt.Errorf("open db failed: %v", err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		return fmt.Errorf("ping db failed: %v", pingErr)
	}
	return nil
}

func main() {
	err := connectDB()
	if err == nil {
		albums, err := getAlbumsByArtist("Phong Nguyen")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Albums: %v\n", albums)

		album, err := getAlbumByID(3)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Album: %v\n", album)

		id, err := insertNewAlbum(Album{1, "My Go Lang", "Phong Nguyen", 100.01})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Album insert: %d\n", id)
	}
}
