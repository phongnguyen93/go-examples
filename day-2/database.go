package database

import (
	"database/sql"
	"fmt"

	"github.com/go-sql-driver/mysql"
)

type Album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

var db *sql.DB

func GetAlbumsByArtist(name string) ([]Album, error) {
	var albums []Album
	rows, err := db.Query("SELECT * FROM album WHERE artist = ?", name)
	if err != nil {
		return nil, fmt.Errorf("DB: GetAlbumsByArtist %q: %v", name, err)
	}
	defer rows.Close()
	for rows.Next() {
		var album Album
		if err := rows.Scan(&album.ID, &album.Title, &album.Artist, &album.Price); err != nil {
			return nil, fmt.Errorf("DB: GetAlbumsByArtist %q: %v", name, err)
		}
		albums = append(albums, album)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("DB: GetAlbumsByArtist %q: %v", name, err)
	}
	return albums, nil
}

func GetAllAlbums() ([]Album, error) {
	var albums []Album
	rows, err := db.Query("SELECT * FROM album")
	if err != nil {
		return nil, fmt.Errorf("DB: GetAllAlbums %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var album Album
		if err := rows.Scan(&album.ID, &album.Title, &album.Artist, &album.Price); err != nil {
			return nil, fmt.Errorf("DB: GetAllAlbums %v", err)
		}
		albums = append(albums, album)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("DB: GetAllAlbums %v", err)
	}
	return albums, nil
}

func GetAlbumByID(id int64) (Album, error) {
	var album Album
	row := db.QueryRow("SELECT * FROM album WHERE id = ?", id)
	if err := row.Scan(&album.ID, &album.Title, &album.Artist, &album.Price); err != nil {
		if err == sql.ErrNoRows {
			return album, fmt.Errorf("DB: GetAlbumByID %d no such album", id)
		}
		return album, fmt.Errorf("DB: GetAlbumByID %d: %v", id, err)
	}
	return album, nil
}

func InsertNewAlbum(album Album) (int64, error) {
	result, err := db.Exec("INSERT INTO album (title, artist, price) VALUES (?,?,?)", album.Title, album.Artist, album.Price)
	if err != nil {
		return 0, fmt.Errorf("DB: InsertNewAlbum %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("DB: InsertNewAlbum %v", err)
	}
	return id, nil
}

func Init() error {
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
