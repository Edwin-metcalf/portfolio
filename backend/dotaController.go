package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	//"golang.org/x/tools/playground"
)

//check appflowy for your stuff
// what are next steps i need to define structs to take in the JSON
//then handle the data and turn it into meaningful stuff I believe

var APIKEY string

// for the api
type winLose struct {
	Win  int `json:"win"`
	Lose int `json:"lose"`
}

// for the return
type winLossRate struct {
	Wins    int     `json:"wins"`
	Losses  int     `json:"lose"`
	WinRate float32 `json:"winRate"`
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file, using environment variables")
	}
	APIKEY = os.Getenv("OPENDOTA_API_KEY")
}

//grab the win lose and return both along with the win rate

func getWinLose(playerID string) (*winLossRate, error) {
	if APIKEY == "" {
		log.Println("OPENDOTA API KEY NOT SET")
	}
	url := fmt.Sprintf("https://api.opendota.com/api/players/%s/wl?api_key=%s", playerID, APIKEY)
	var winLosePull winLose

	response, err := http.Get(url)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	defer response.Body.Close()
	err = json.NewDecoder(response.Body).Decode(&winLosePull)

	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return nil, err
	}

	winPercentage := float32(winLosePull.Win) / float32(winLosePull.Win+winLosePull.Lose)

	winLossPercentage := winLossRate{
		Wins:    winLosePull.Win,
		Losses:  winLosePull.Lose,
		WinRate: winPercentage,
	}
	fmt.Printf("this is what we are returning:  %+v", winLossPercentage)
	return &winLossPercentage, nil
}

// grab recent games and structs for it
type recentGame struct {
	MatchID      int    `json:"match_id"`
	PlayerSlot   int    `json:"player_slot"`
	RadiantWin   bool   `json:"radiant_win"`
	Duration     int    `json:"duration"`
	GameMode     int    `json:"game_mode"`
	LobbyType    int    `json:"lobby_type"`
	HeroID       int    `json:"hero_id"`
	StartTime    int    `json:"start_time"`
	Version      string `json:"version"`
	Kills        int    `json:"kills"`
	Deaths       int    `json:"deaths"`
	Assists      int    `json:"assists"`
	AverageRank  int    `json:"average_rank"`
	LeaverStatus int    `json:"leaver_status"`
	PartySize    int    `json:"party_size"`
	HeroVariant  int    `json:"hero_variant"`
}

func getRecentGames(playerID string) ([]recentGame, error) {
	if APIKEY == "" {
		log.Println("OPENDOTA API KEY NOT SET")
	}
	numberOfGames := 7
	url := fmt.Sprintf("https://api.opendota.com/api/players/%s/matches?api_key=%s&limit=%s", playerID, APIKEY, numberOfGames)

	var recentGameList []recentGame

	response, err := http.Get(url)

	if err != nil {
		fmt.Println(err.Error())
	}

	defer response.Body.Close()

	err = json.NewDecoder(response.Body).Decode(&recentGameList)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("Number of games fetched: %d\n", len(recentGameList))
	fmt.Printf("First game: %+v\n", recentGameList[6])
	return recentGameList, err
}

/*func processRawRecentGames(rawData []recentGame) { //needs to return struct we want to send to the front end
	processedData := make([])
	I need to make a decision on wether this should just return wins loses and winrate or should it be more indepth
	lowkey maybe all the above then let the front end deal with it
} */

func testing() {
	apikey := APIKEY

	if apikey == "" {
		log.Println("OPENDOTA_API_KEY not set")
	}
	//good tempelate for calling information from the api
	url := fmt.Sprintf("https://api.opendota.com/api/players/287883142?api_key=%s", apikey)
	response, err := http.Get(url)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer response.Body.Close()

	responseData, err := io.ReadAll(response.Body)

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(responseData))
}
