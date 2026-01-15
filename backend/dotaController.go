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

type winLose struct {
	Win  int `json:"win"`
	Lose int `json:"lose"`
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	APIKEY = os.Getenv("OPENDOTA_API_KEY")
}
func getWinLose(playerID string) {
	if APIKEY == "" {
		log.Fatal("OPENDOTA API KEY NOT SET")
	}
	url := fmt.Sprintf("https://api.opendota.com/api/players/%s/wl?api_key=%s", playerID, APIKEY)
	var winLosePull winLose

	response, err := http.Get(url)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer response.Body.Close()
	err = json.NewDecoder(response.Body).Decode(&winLosePull)

	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}
	fmt.Printf("Wins: %d, Losses: %d\n", winLosePull.Win, winLosePull.Lose)

}

func testing() {
	apikey := APIKEY

	if apikey == "" {
		log.Fatal("OPENDOTA_API_KEY not set")
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
