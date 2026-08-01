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
	Profile   *PlayerProfileReturn          `json:"profile"`
	BattleLog []CRLadderChartPoint          `json:"battleLog"`
	Friendly  ClashRoyaleFriendlyLoadReturn `json:"friendly"`
}

func clashRoyaleLoadHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	//me and Ry Guy
	myPlayerId := "#P9L0U88GQ"
	friendPlayerId := "#QQCJYR0Y8"
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
	friendlyStats, err := loadFriendlyStats(DB, CRclient, myPlayerId, friendPlayerId)
	if err != nil {
		log.Printf("Error loading friendly stats: %v", err)
		http.Error(w, "Error fetching friendly stats", http.StatusInternalServerError)
		return
	}
	result := ClashRoyaleLoadReturn{
		Profile:   playerInfo,
		BattleLog: chartPoints,
		Friendly:  *friendlyStats,
	}
	json.NewEncoder(w).Encode(result)
}

func clashRoyaleFriendlyHandler(w http.ResponseWriter, r *http.Request) {
	//will have to edit this to be able to take arbitrary tags in future
	w.Header().Set("Content-Type", "application/json")

	myTag := "#P9L0U88GQ"
	friendTag := "#QQCJYR0Y8"

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	CRclient := newClashRoyaleClient()
	stats, err := loadFriendlyStats(DB, CRclient, myTag, friendTag)
	if err != nil {
		http.Error(w, "Error fetching friendly stats", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(stats)
}

type matchupGeneratorTags struct {
	Tag1 string `json:"tag1"`
	Tag2 string `json:"tag2"`
}

func matchupGeneratorHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var CRTags matchupGeneratorTags

	err := json.NewDecoder(r.Body).Decode(&CRTags)

	if err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	CRclient := newClashRoyaleClient()
	matchupGeneratorReturn, err := matchupGeneratorhelper(CRclient, CRTags.Tag1, CRTags.Tag2)
	if err != nil {
		http.Error(w, "Error fetching matchup stats stats", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(matchupGeneratorReturn)

}

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
