package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	//"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

//there is some reference stuff in the golearn you know this

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		allowedOrigins := []string{
			"https://edwinmetcalf.com",
			"https://www.edwinmetcalf.com",
		}

		origin := r.Header.Get("Origin")
		for _, allowed := range allowedOrigins {
			if origin == allowed {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				break
			}
		}
		if origin == "" || w.Header().Get("Access-Control-Allow-Origin") == "" {
			w.Header().Set("Access-Control-Allow-Origin", "*")

		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

type SpaceInvadersEntry struct {
	Name  string `json:"name"`
	Score int    `json:"score"`
}

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

func healthHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "ok",
		"message": "Somehow the server is running",
	})
}

func main() {
	var err error
	DB, err = initDB()
	if err != nil {
		log.Fatal("failed to intialize database: ", err)
	}
	defer DB.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("/api/health", healthHandler)
	mux.HandleFunc("/api/space-invaders/entry", addSpaceInvadersEntry)
	mux.HandleFunc("/api/space-invaders/leaderboard", getSpaceInvadersLeaderboard)
	handler := corsMiddleware(mux)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("server up on god @ port %s\n", port)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, handler))
}
