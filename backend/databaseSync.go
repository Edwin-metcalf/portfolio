package main

import (
	"log"
	"time"
)

func synchClashRoyaleDatabase() {
	// need these things
	start := time.Now()
	log.Println("[sync] starting clash royale sync")

	myPlayerId := "#P9L0U88GQ"
	friendPlayerId := "#QQCJYR0Y8"
	CRclient := newClashRoyaleClient()

	//then sync
	if err := SyncLadderGames(DB, CRclient, myPlayerId); err != nil {
		log.Printf("[sync] ladder failed: %v", err)
	} else {
		log.Println("[sync] ladder ok")
	}

	if err := syncRankedGames(DB, CRclient, myPlayerId); err != nil {
		log.Printf("[sync] ranked failed: %v", err)
	} else {
		log.Println("[sync] ranked ok")
	}

	if err := syncFriendlyGames(DB, CRclient, myPlayerId, friendPlayerId); err != nil {
		log.Printf("[sync] firendly failed: %v", err)
	} else {
		log.Println("[sync] friendly ok")
	}

	log.Printf("[sync] finished in %s", time.Since(start).Round(time.Millisecond))

}
