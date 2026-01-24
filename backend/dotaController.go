package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
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
	fmt.Printf("Wins: %d, Losses: %d\n", winLosePull.Win, winLosePull.Lose)
	winPercentage := winLosePull.Win / (winLosePull.Win + winLosePull.Lose)
	winLossPercentage := winLossRate{
		Wins:    winLosePull.Win,
		Losses:  winLosePull.Lose,
		WinRate: float32(winPercentage),
	}
	/* there  must be a better way to initalize and declare a struct like this
	winLossPercentage {
		wins: winLosePull.Win,
		losses: winLosePull.Lose,
		winRate: winPercentage,
	}
	*/
	return &winLossPercentage, nil
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
