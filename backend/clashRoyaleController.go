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

//===============================================
// THE CLASH ROYALE API DOES NOT DO API KEY IT DOES JSON WEB TOKENS

// type section this is the final return struct
// lowkey very subject to change
type ClashRoyaleStatsReturn struct {
	Wins            int   `json:"wins"`
	Losses          int   `json:"losses"`
	WinRate         int   `json:"winRate"`
	CurrentTrophies int   `json:"currentTrophies"`
	TrophyProgress  []int `json:"trophyProgress"`
}

type PlayerProfile struct {
}
type PlayerProfileReturn struct {
}

// api key stuff is actually a JSON Web Token kinda cool something new

var CR_APIKEY string

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file, using environment variables")
	}
	CR_APIKEY = os.Getenv("CLASH_ROYALE_API_KEY")
}

// helper to wrap the api calls in
type ClashRoyaleClient struct {
	baseURL string
	jwt     string
	client  *http.Client
}

func newClashRoyaleClient() *ClashRoyaleClient {
	// make a new client and add the url and our jwt to it
	return &ClashRoyaleClient{
		baseURL: "https://api.clashroyale.com/v1",
		jwt:     CR_APIKEY,
		client:  &http.Client{},
	}
}

func (c *ClashRoyaleClient) getPlayerProfile(playerTag string) (*PlayerProfile, error) {
	req, err := http.NewRequest("GET", c.baseURL+"/players/"+playerTag, nil)
	if err != nil {
		log.Println("error creating request")
	}

	req.Header.Add("Authorization", "Bearer "+c.jwt)
	resp, err := c.client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
	}

	defer resp.Body.Close()
	var profile PlayerProfile
	err = json.NewDecoder(resp.Body).Decode(&profile)
	if err != nil {
		return nil, err
	}
	// do something with this
	return &profile, err
}

//helper to get the wins and losses and create the winrate

func returnClashRoyaleStats(playerID string) (*ClashRoyaleStatsReturn, error) {
	return nil, nil
}
func testing() {
	jwt := CR_APIKEY

	if jwt == "" {
		log.Println("CLASH_ROYALE_API_KEY not set")
	}
	client := &http.Client{}

	url := fmt.Sprintf("https://api.clashroyale.com/v1/players/")

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Authorization", jwt)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
	}

	defer resp.Body.Close()
	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(responseData))
}
