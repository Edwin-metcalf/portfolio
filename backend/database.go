package main

import (
	"database/sql"
	"log"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func initDB() (*sql.DB, error) {
	var err error
	var driver, dsn string

	databaseURL := os.Getenv("DATABASE_URL")

	if databaseURL == "" {
		log.Println("using local fallback which is a sqlite")
		driver = "sqlite3"
		dsn = "./local.db"
		//databaseURL = "postgres://localhost/space_invaders?sslmode=disable"
	} else {
		log.Println("DATABASE_URL found, connecting to Postgres...")
		driver = "postgres"
		dsn = databaseURL

		//possible fix for heroku postgres issues
		if strings.HasPrefix(dsn, "postgres://") {
			dsn = strings.Replace(dsn, "postgres://", "postgresql://", 1)
		}

	}
	DB, err = sql.Open(driver, dsn)
	if err != nil {
		log.Fatal(err)
	}

	err = DB.Ping()
	if err != nil {
		return nil, err
	}
	var sqlStmt string
	if driver == "postgres" {
		sqlStmt = `
		CREATE TABLE IF NOT EXISTS leaderboard (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL, 
		score INTEGER NOT NULL
		);`
	} else {
		sqlStmt = `
		CREATE TABLE IF NOT EXISTS leaderboard (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL, 
		score INTEGER NOT NULL
		);`
	}

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
	row, err := db.Query("SELECT name, score FROM leaderboard ORDER BY score DESC LIMIT $1", limit)
	if err != nil {
		log.Printf("ERROR in getTopScores query: %v", err)
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
	_, err := db.Exec("INSERT INTO leaderboard (name, score) VALUES($1,$2)", name, score)

	return err
}

func deleteLeaderboard(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM leaderboard")
	return err
}
