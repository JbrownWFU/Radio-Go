package storage

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

type DBConn struct {
	Filepath string
	Conn     *sql.DB
}

type Station struct {
	ID   int
	Name string
	Url  string
}

type LikedSong struct {
	ID          int
	Artist      string
	Title       string
	StationName string
	LikedAt     string
}

// InitDB opens the database file and creates the schema if it doesn't exist.
func InitDB(dbPath string) (*DBConn, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	schema := `
	CREATE TABLE IF NOT EXISTS stations (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		url TEXT NOT NULL
	);
	CREATE TABLE IF NOT EXISTS liked_songs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		artist TEXT NOT NULL,
		title TEXT NOT NULL,
		station_name TEXT,
		liked_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`

	if _, err := db.Exec(schema); err != nil {
		return nil, err
	}

	return &DBConn{
		Filepath: dbPath,
		Conn:     db,
	}, nil
}

func (dbConn *DBConn) GetStations() ([]Station, error) {
	query := "SELECT id, name, url FROM stations"
	rows, err := dbConn.Conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stations []Station
	for rows.Next() {
		var s Station
		if err := rows.Scan(&s.ID, &s.Name, &s.Url); err != nil {
			return nil, err
		}
		stations = append(stations, s)
	}
	return stations, nil
}

func (dbConn *DBConn) InsertStation(name, url string) error {
	query := "INSERT INTO stations (name, url) VALUES (?, ?)"
	_, err := dbConn.Conn.Exec(query, name, url)
	return err
}

func (dbConn *DBConn) DeleteStation(id string) error {
	query := "DELETE FROM stations WHERE id = ?"
	_, err := dbConn.Conn.Exec(query, id)
	return err
}

func (dbConn *DBConn) LikeSong(artist, title, stationName string) error {
	query := "INSERT INTO liked_songs (artist, title, station_name) VALUES (?, ?, ?)"
	_, err := dbConn.Conn.Exec(query, artist, title, stationName)
	return err
}

func (dbConn *DBConn) GetLikedSongs() ([]LikedSong, error) {
	query := "SELECT id, artist, title, station_name, liked_at FROM liked_songs ORDER BY liked_at DESC"
	rows, err := dbConn.Conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var songs []LikedSong
	for rows.Next() {
		var s LikedSong
		if err := rows.Scan(&s.ID, &s.Artist, &s.Title, &s.StationName, &s.LikedAt); err != nil {
			return nil, err
		}
		songs = append(songs, s)
	}
	return songs, nil
}

func (dbConn *DBConn) DeleteLikedSong(id string) error {
	query := "DELETE FROM liked_songs WHERE id = ?"
	_, err := dbConn.Conn.Exec(query, id)
	return err
}
