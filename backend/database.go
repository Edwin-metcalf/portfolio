package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func initDB() (*sql.DB, error) {
	var err error
	DB, err = sql.Open("sqlite3", "./space-invaders.db")
	if err != nil {
		log.Fatal(err)
	}

	sqlStmt := `
	CREATE TABLE IF NOT EXISTS leaderboard (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL, 
	score INTEGER NOT NULL
	);`

	_, err = DB.Exec(sqlStmt)
	if err != nil {
		return nil, err
	}
	return DB, err
}

func getTopScores(db *sql.DB, limit int) ([]SpaceInvadersEntry, error) {
	if limit == 0 {
		limit = 20
	}
	var entries []SpaceInvadersEntry
	row, err := db.Query("SELECT name, score FROM leaderboard ORDER BY score DESC LIMIT ?", limit)
	if err != nil {
		return nil, err
	}
	defer row.Close()
	for row.Next() {
		item := SpaceInvadersEntry{}
		err := row.Scan(&item.Name, &item.Score)
		if err != nil {
			return nil, err
		}

		entries = append(entries, item)
	}
	if err = row.Err(); err != nil {
		return nil, err
	}
	return entries, nil
}

func addScore(db *sql.DB, name string, score int) error {
	_, err := db.Exec("INSERT INTO leaderboard (name, score) VALUES(?,?)", name, score)

	return err
}
