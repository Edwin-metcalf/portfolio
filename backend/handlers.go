package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type SpaceInvadersEntry struct {
	Name  string `json:"name"`
	Score int    `json:"score"`
}

// space invader handlers
func addSpaceInvadersEntry(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var newEntry SpaceInvadersEntry

	err := json.NewDecoder(r.Body).Decode(&newEntry)
	if err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	err = addScore(DB, newEntry.Name, newEntry.Score)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	fmt.Printf("Recieved entry: %+v\n", newEntry)
	//tthis is where I would add the new Entry to the database type beat
	json.NewEncoder(w).Encode(newEntry)
}

func getSpaceInvadersLeaderboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	//fake data for testing before the database get put in this john as fuck as fuck

	leaderboard, err := getTopScores(DB, 10)

	if err != nil {
		http.Error(w, "databse error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(leaderboard)
}
func clearLeaderboardHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	apiKey := r.Header.Get("Authorization")
	secretKey := os.Getenv("ADMIN_API_KEY")

	if secretKey == "" {
		http.Error(w, "Server Configuration Error", http.StatusInternalServerError)
		return
	}

	if apiKey != "Bearer "+secretKey {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err := deleteLeaderboard(DB)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Leaderboard cleared successfylly",
	})
}

// a handler to check if the backend is running
func healthHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if DB != nil {
		err := DB.Ping()
		if err != nil {
			log.Printf("ERROR: Database ping failed: %v", err)
			json.NewEncoder(w).Encode(map[string]string{
				"status":  "error",
				"message": "Database connection failed",
				"error":   err.Error(),
			})
			return
		}
	}
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "ok",
		"message": "Somehow the server is running",
	})
}

// dota stuff handlers
type DotaPlayerStats struct {
	Overall *winLossRate       `json:"overall"`
	Recent  *recentGameCleaned `json:"recent"`
}

func dotaStatsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	playerID := "287883142"
	stats, err := returnGamesSats(playerID)

	if err != nil {
		log.Printf("Error fetching dota stats: %v", err)
		http.Error(w, "Error fetching player stats", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(stats)
}

// struct to return my profile and my battle history
type CRLadderChartPoint struct {
	BattleTime string `json:"battleTime"`
	Trophies   int    `json:"trophies"`
}

type ClashRoyaleLoadReturn struct {
	Profile   *PlayerProfileReturn `json:"profile"`
	BattleLog []CRLadderChartPoint `json:"battleLog"`
}

func clashRoyaleLoadHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	myPlayerId := "#P9L0U88GQ"
	//friendPlayerId := "#QQCJYR0Y8"
	CRclient := newClashRoyaleClient()
	playerInfo, err := CRclient.getPlayerProfile(myPlayerId)
	if err != nil {
		log.Printf("Error fetching clash Royale stats: %v", err)
		http.Error(w, "Error fetching clash royale stats", http.StatusInternalServerError)
		return
	}
	//sync before getting hisotry
	if err := SyncLadderGames(DB, CRclient, myPlayerId); err != nil {
		log.Printf("Sync Warning: %v", err)
		//DB has history if fails dont kill everything
	}

	ladderGames, err := getLadderHistory(DB)
	if err != nil {
		log.Printf("Error fetching ladder history: %v", err)
		http.Error(w, "Error fetching ladder history", http.StatusInternalServerError)
		return
	}
	var chartPoints []CRLadderChartPoint
	for _, game := range ladderGames {
		chartPoints = append(chartPoints, CRLadderChartPoint{
			BattleTime: game.BattleTime,
			Trophies:   game.StartingTrophies + game.TrophyChange,
		})
	}

	// do stuff for friendlies currently just against friend but then add against whoever

	result := ClashRoyaleLoadReturn{
		Profile:   playerInfo,
		BattleLog: chartPoints,
	}
	json.NewEncoder(w).Encode(result)

	//gonna need my clash id
	//stats, err := ClashRoyaleStatsReturn(playerID)
}

/*
func clashRoyaleLadderHistoryHandler(w http.ResponseWriter, r *http.Request, playerId string) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	CRclient := newClashRoyaleClient()
	playerBattleHistory, err  := CRclient.getPlayerBattleLog(playerId)
	if err != nil {
		log.Printf("Error fetching Clash Royale stats: %v", err)
		http.Error(w, "Error fetching clash royale stats", http.StatusInternalServerError)
		return
	}
	// this may cause error I dont know why i need to have a 0 in here i want the whole list
	result := DateRankList{playerBattleHistory.RankList[0]}
	json.NewEncoder(w).Encode(result)
}

// should this function only take 1 player and then compare results latter?
func clashRoyaleFriendlyHisoryHandler(w http.ResponseWriter, r *http.Request, player1 string, player2 string) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	CRclient := newClashRoyaleClient()

}
*/

// stuff for Clash royale database entries
type CRLadderDataBaseEntry struct {
	BattleTime       string `json:"batteTime"`
	StartingTrophies int    `json:"startingTrophies"`
	TrophyChange     int    `json:"trophyChange"`
	Result           int    `json:"result"`
}

type CRFriendlyDataBaseEntry struct {
	BattleTime string          `json:"battleTime"`
	Result     int             `json:"result"`
	MyDeck     json.RawMessage `json:"myDeck"`
	EnemyDeck  json.RawMessage `json:"enemyDeck"`
}
