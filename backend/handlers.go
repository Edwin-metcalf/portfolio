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
func dotaWinLoseHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		http.Error(w, "method not aloud", http.StatusMethodNotAllowed)
		return
	}
	winLossPercentage, err := getWinLose("287883142")

	if err != nil {
		http.Error(w, "Error with open dota api", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(winLossPercentage)
}
