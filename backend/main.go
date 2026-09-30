package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	//"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

//there is some reference stuff in the golearn you know this

func corsMiddleware(next http.Handler) http.Handler {
	//this is cors middleware to give some security
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		allowedOrigins := map[string]bool{
			"https://edwinmetcalf.com":     true,
			"https://www.edwinmetcalf.com": true,
		}

		origin := r.Header.Get("Origin")
		if strings.HasPrefix(origin, "http://localhost") ||
			strings.HasPrefix(origin, "http://127.0.0.1") {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		} else if allowedOrigins[origin] {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}

		w.Header().Set("Vary", "Origin")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// handlers are now moved to handlers.go

func main() {
	//testing()
	//getWinLose("287883142")
	var err error
	DB, err = initDB()
	if err != nil {
		log.Fatal("failed to intialize database: ", err)
	}
	defer DB.Close()
	getRecentGames("287883142")

	mux := http.NewServeMux()
	mux.HandleFunc("/api/health", healthHandler)
	mux.HandleFunc("/api/games/space-invaders/entry", addSpaceInvadersEntry)
	mux.HandleFunc("/api/games/space-invaders/leaderboard", getSpaceInvadersLeaderboard)
	mux.HandleFunc("/api/games/space-invaders/clear", clearLeaderboardHandler)
	mux.HandleFunc("/api/games/dota", dotaStatsHandler)
	mux.HandleFunc("/api/games/clash-royale", clashRoyaleLoadHandler)
	mux.HandleFunc("/api/games/clash-royale/matchup-generator", matchupGeneratorHandler)
	handler := corsMiddleware(mux)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("server up on god @ port %s\n", port)

	//set up auto sync every 6 hours
	syncTimer := time.NewTicker(6 * time.Hour)
	//syncDone := make(chan bool) this never really stops running....

	go func() {
		synchClashRoyaleDatabase()
		for t := range syncTimer.C {
			log.Printf("[sync] ticker fires at %s", t.Format(time.RFC3339))
			synchClashRoyaleDatabase()

		}
	}()

	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, handler))
}
