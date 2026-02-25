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
	Losses  int     `json:"losses"`
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
	MatchID      int  `json:"match_id"`
	PlayerSlot   int  `json:"player_slot"`
	RadiantWin   bool `json:"radiant_win"`
	Duration     int  `json:"duration"`
	GameMode     int  `json:"game_mode"`
	LobbyType    int  `json:"lobby_type"`
	HeroID       int  `json:"hero_id"`
	StartTime    int  `json:"start_time"`
	Version      int  `json:"version"`
	Kills        int  `json:"kills"`
	Deaths       int  `json:"deaths"`
	Assists      int  `json:"assists"`
	AverageRank  int  `json:"average_rank"`
	LeaverStatus int  `json:"leaver_status"`
	PartySize    int  `json:"party_size"`
	HeroVariant  int  `json:"hero_variant"`
}

func getRecentGames(playerID string) (*[]recentGame, error) {
	if APIKEY == "" {
		log.Println("OPENDOTA API KEY NOT SET")
	}
	numberOfGames := 7
	url := fmt.Sprintf("https://api.opendota.com/api/players/%s/matches?api_key=%s&limit=%d", playerID, APIKEY, numberOfGames)

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
	return &recentGameList, err
}

type recentGameCleaned struct {
	Wins    int       `json:"wins"`
	Losses  int       `json:"losses"`
	WinRate float32   `json:"winRate"`
	AvgKDA  []float32 `json:"avgKDA"`
	// could do some cool stuff like getting kills per game and stuff

}

// should maybe call the recent games and the process togehter so the handler only calls one
func processRawRecentGames(rawData []recentGame) (*recentGameCleaned, error) { //needs to return struct we want to send to the front end
	var cleanedGames recentGameCleaned
	cleanedGames.AvgKDA = make([]float32, 3)
	positionMap := make(map[int]int)
	var totalKda [3]int
	numberOfGames := len(rawData)

	for _, val := range rawData {
		var team string
		//process the team and the position
		if val.PlayerSlot > 6 {
			team = "radiant"
			positionMap[val.PlayerSlot] += 1
		} else {
			team = "dire"
			positionMap[val.PlayerSlot-128] += 1 //or whatever the correct bit number is
		}

		//process if the game was a win
		if team == "radiant" && val.RadiantWin || team == "dire" && !val.RadiantWin {
			cleanedGames.Wins += 1
		} else {
			cleanedGames.Losses += 1
		}
		//get the total KDA
		totalKda[0] += val.Kills
		totalKda[1] += val.Deaths
		totalKda[2] += val.Assists

	}
	//do some averaging for better return struct
	cleanedGames.AvgKDA[0] = float32(totalKda[0] / numberOfGames)
	cleanedGames.AvgKDA[1] = float32(totalKda[1] / numberOfGames)
	cleanedGames.AvgKDA[2] = float32(totalKda[2] / numberOfGames)

	cleanedGames.WinRate = float32(cleanedGames.Wins) / float32(cleanedGames.Wins+cleanedGames.Losses)
	return &cleanedGames, nil
}

type dotaStatsReturn struct {
	Wins          int       `json:"wins"`
	Losses        int       `json:"losses"`
	WinRate       float32   `json:"winRate"`
	RecentWins    int       `json:"recentWins"`
	RecentLosses  int       `json:"recentLosses"`
	RecentWinRate float32   `json:"recentWinRate"`
	AvgKDA        []float32 `json:"avgKDA"`
}

func returnGamesSats(playerID string) (*dotaStatsReturn, error) {
	winLoss, err := getWinLose(playerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get win/loss stats: %w", err)
	}

	recentGames, err := getRecentGames(playerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent games: %w", err)
	}

	recentStats, err := processRawRecentGames(*recentGames)
	if err != nil {
		return nil, fmt.Errorf("failed to process recent games: %w", err)
	}

	stats := &dotaStatsReturn{
		Wins:          winLoss.Wins,
		Losses:        winLoss.Losses,
		WinRate:       winLoss.WinRate,
		RecentWins:    recentStats.Wins,
		RecentLosses:  recentStats.Losses,
		RecentWinRate: recentStats.WinRate,
		AvgKDA:        recentStats.AvgKDA,
	}

	return stats, nil
}

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
