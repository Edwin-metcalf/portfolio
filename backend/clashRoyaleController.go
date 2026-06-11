package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
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
	Tag            string `json:"tag"`
	Name           string `json:"name"`
	Trophies       int    `json:"trophies"`
	BestTrophies   int    `json:"bestTrophies"`
	Wins           int    `json:"wins"`
	Losses         int    `json:"losses"`
	BattleCount    int    `json:"battleCount"`
	ThreeCrownWins int    `json:"threeCrownWins"`
	Clan           Clan   `json:"clan"`
	Arena          Arena  `json:"arena"`
}
type Clan struct {
	Tag  string `json:"tag"`
	Name string `json:"name"`
}
type Arena struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}
type PlayerProfileReturn struct {
	Tag            string  `json:"tag"`
	Name           string  `json:"name"`
	Trophies       int     `json:"trophies"`
	BestTrophies   int     `json:"bestTrophies"`
	Wins           int     `json:"wins"`
	Losses         int     `json:"losses"`
	WinRate        float64 `json:"winRate"`
	BattleCount    int     `json:"battleCount"`
	ThreeCrownWins int     `json:"threeCrownWins"`
	Clan           Clan    `json:"clan"`
	Arena          Arena   `json:"arena"`
}
type GameMode struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

// ther is a lot more that you can pull out of these dont know yet
type PlayerBattleDataList struct {
	Tag  string `json:"tag"`
	Name string `json:"name"`
}

// do not really know how the opponent list works
type Battle struct {
	GameMode GameMode               `json:"gameMode"`
	Opponent []PlayerBattleDataList `json:"playerBattleList"`
}
type PlayerBattleLog struct {
	Battles []Battle `json:"battleList"`
}
type PlayerBattleLogReturn struct {
}

// api key stuff is actually a JSON Web Token kinda cool something new

var CR_APIKEY string

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file, using environment variables")
	}
	//set this to the production one when needed now its the dev one this may break a bunch lmao
	//CR_APIKEY = os.Getenv("CLASH_ROYALE_API_KEY")
	CR_APIKEY = os.Getenv("CLASH_ROYALE_JWT_DEV")

	if CR_APIKEY == "" {
		log.Fatal("clash royale JWT is empty/not set ")
	}
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
		//THIS IS THE DEV login you will need to change this to the correct deployment one
		// Before:
		// baseURL: "https://api.clashroyale.com/v1",

		// After:
		baseURL: "https://proxy.royaleapi.dev/v1",
		jwt:     CR_APIKEY,
		client:  &http.Client{},
	}
}

func (c *ClashRoyaleClient) getPlayerProfile(playerTag string) (*PlayerProfileReturn, error) {
	encodedTag := url.PathEscape(playerTag)
	req, err := http.NewRequest("GET", c.baseURL+"/players/"+encodedTag, nil)
	if err != nil {
		log.Println("error creating request")
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+c.jwt)
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	defer resp.Body.Close()
	var p PlayerProfile
	err = json.NewDecoder(resp.Body).Decode(&p)
	if err != nil {
		return nil, err
	}
	// do something with this
	var retProfile PlayerProfileReturn

	//calculate win rate
	winRate := 0.0
	if p.Wins+p.Losses > 0 {
		winRate = float64(p.Wins) / float64(p.Wins+p.Losses)

	}

	retProfile.Tag = p.Tag
	retProfile.Name = p.Name
	retProfile.Trophies = p.Trophies
	retProfile.BestTrophies = p.BestTrophies
	retProfile.Wins = p.Wins
	retProfile.Losses = p.Losses
	retProfile.WinRate = winRate
	retProfile.BattleCount = p.BattleCount
	retProfile.ThreeCrownWins = p.ThreeCrownWins
	retProfile.Clan = p.Clan
	retProfile.Arena = p.Arena

	return &retProfile, err
}
func (c *ClashRoyaleClient) getPlayerBattlelog(playerTag string) (*PlayerBattleLogReturn, error) {
	encodedTag := url.PathEscape(playerTag)
	req, err := http.NewRequest("GET", c.baseURL+"/players/"+encodedTag+"/battlelog", nil)
	if err != nil {
		log.Println("Error creating request")
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+c.jwt)
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("API returned status %d: %s", resp.StatusCode, string(body))
		return nil, fmt.Errorf("api error: status %d", resp.StatusCode)
	}

	var battleLog PlayerBattleLog
	err = json.NewDecoder(resp.Body).Decode(&battleLog)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// this might not need the player Tag as this could be the on load without input
// the matchups one should need IDs
func returnClashRoyaleStats(playerTag string) (*ClashRoyaleStatsReturn, error) {
	CRclient := newClashRoyaleClient()
	//maybe pass it in or not hard code it
	_, err := CRclient.getPlayerProfile(playerTag)
	if err != nil {
		log.Printf("Error fetching clash Royale stats: %v", err)
		return nil, err
	}
	var CRStatsLoad ClashRoyaleStatsReturn

	return &CRStatsLoad, nil
}

/*
testign guy if needed
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
*/
