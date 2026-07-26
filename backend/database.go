package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
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
	//this is for the space invaders games
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
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL, 
		score INTEGER NOT NULL
		);`
	}
	_, err = DB.Exec(sqlStmt)
	if err != nil {
		return nil, err
	}
	//table for the clash royale stats
	var CRTrophyStmt string
	if driver == "postgres" {
		CRTrophyStmt = `
		CREATE TABLE IF NOT EXISTS ladderGames (
		id SERIAL PRIMARY KEY,
		battle_time TEXT NOT NULL,
		starting_trophies INTEGER NOT NULL,
		trophy_change INTEGER NOT NULL,
		result SMALLINT
		);`
	} else {
		CRTrophyStmt = `
		CREATE TABLE IF NOT EXISTS ladderGames (
		id INTEGER PRIMARY KEY,
		battle_time TEXT NOT NULL,
		starting_trophies INTEGER NOT NULL,
		trophy_change INTEGER NOT NULL,
		result TINYINT
		);`
	}
	_, err = DB.Exec(CRTrophyStmt)
	if err != nil {
		return nil, err
	}

	var CRFriendlyStmt string
	if driver == "postgres" {
		CRFriendlyStmt = `
		CREATE TABLE IF NOT EXISTS friendlyGames (
		id SERIAL PRIMARY KEY,
		battle_time TIMESTAMP NOT NULL,
		result SMALLINT,
		my_deck JSON NOT NULL,
		enemy_deck JSON NOT NULL
		);`
	} else {
		CRFriendlyStmt = `
		CREATE TABLE IF NOT EXISTS friendlyGames (
		id INTEGER PRIMARY KEY,
		battle_time TEXT NOT NULL,
		result TINYINT,
		my_deck TEXT NOT NULL,
		enemy_deck TEXT NOT NULL
		);`
	}

	_, err = DB.Exec(CRFriendlyStmt)
	if err != nil {
		return nil, err
	}
	return DB, err
}

// space invaders section ====
func getTopScores(db *sql.DB, limit int) ([]SpaceInvadersEntry, error) {
	if limit == 0 {
		limit = 20
	}
	var entries []SpaceInvadersEntry
	row, err := db.Query("SELECT name, score FROM leaderboard ORDER BY score DESC LIMIT $1;", limit)
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
	_, err := db.Exec("INSERT INTO leaderboard (name, score) VALUES($1,$2);", name, score)

	return err
}

func deleteLeaderboard(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM leaderboard;")
	return err
}

// Clash royale section ------------

// replace the string int pair with this guy could be here or in the other file
type CRLadderGame struct {
	BattleTime       string `json:"battleTime"`
	StartingTrophies int    `json:"startingTrophies"`
	TrophyChange     int    `json:"trophyChange"`
	Result           int    `json:"result"`
}

func getLadderHistory(db *sql.DB) ([]CRLadderGame, error) {
	var games []CRLadderGame
	rows, err := db.Query("SELECT battle_time, starting_trophies, trophy_change, result FROM ladderGames ORDER BY battle_time;")
	if err != nil {
		log.Printf("ERROR in getLadder History query: %v", err)
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		item := CRLadderGame{}
		err := rows.Scan(&item.BattleTime, &item.StartingTrophies, &item.TrophyChange, &item.Result)
		if err != nil {
			return nil, err
		}
		games = append(games, item)
	}

	return games, nil
}
func addLadderEntry(db *sql.DB, game CRLadderGame) error {
	if game.Result < -1 || game.Result > 1 {
		return fmt.Errorf("result value needs to be 1, 0 or -1, it is: %v", game.Result)
	}
	_, err := db.Exec(`INSERT INTO ladderGames (battle_time, starting_trophies, trophy_change, result) 
					VALUES($1,$2,$3,$4);`, game.BattleTime, game.StartingTrophies, game.TrophyChange, game.Result)

	return err
}

func getMostRecentCRGame(db *sql.DB, table string) (string, error) {
	query := fmt.Sprintf("SELECT battle_time FROM %v ORDER BY battle_time DESC LIMIT 1;", table)
	var date string
	err := db.QueryRow(query).Scan(&date)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		log.Printf("Error in getting recent clash royale game: %v", err)
		return "", err
	}
	return date, nil

}

func getFriendlyHistory(db *sql.DB) ([]CRFriendlyDataBaseEntry, error) {
	var friendyGameHistory []CRFriendlyDataBaseEntry
	rows, err := db.Query("SELECT battle_time, result, my_deck, enemy_deck FROM friendlyGames ORDER BY battle_time")
	if err != nil {
		log.Printf("ERROR in get friendy history query %v", err)
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		item := CRFriendlyDataBaseEntry{}
		err := rows.Scan(&item.BattleTime, &item.Result, &item.MyDeck, &item.EnemyDeck)
		if err != nil {
			return nil, err
		}
		friendyGameHistory = append(friendyGameHistory, item)

	}

	return friendyGameHistory, nil
}
func addFriendlyEntry(db *sql.DB, battle_time string, result int, my_deck json.RawMessage, enemy_deck json.RawMessage) error {
	if result < -1 || result > 1 {
		return fmt.Errorf("result value needs to be 1, 0 or -1, it is: %v", result)
	}
	_, err := db.Exec(`INSERT INTO friendlyGames (battle_time, result, my_deck, enemy_deck)
					VALUES($1, $2, $3, $4);`, battle_time, result, my_deck, enemy_deck)
	return err
}
