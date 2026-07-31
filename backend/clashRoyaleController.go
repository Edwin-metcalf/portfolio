package main

import (
	"database/sql"
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
type PlayerBattleData struct {
	Tag              string `json:"tag"`
	Name             string `json:"name"`
	StartingTrophies int    `json:"startingTrophies"`
	TrophyChange     int    `json:"trophyChange"`
	Crowns           int    `json:"crowns"`
}

//type PlayerBattleDataList struct {
//	List PlayerBattleData `json:"playerBattleData"`
//}

type Battle struct {
	Type       string             `json:"type"`
	BattleTime string             `json:"battleTime"`
	GameMode   GameMode           `json:"gameMode"`
	Team       []PlayerBattleData `json:"team"`
	Opponent   []PlayerBattleData `json:"opponent"`
}
type BattleList []Battle

// used for date and rank
type StringIntPair struct {
	Text  string `json:"text"`
	Value int    `json:"value"`
}
type DateRankList []StringIntPair

type PlayerBattleLogReturn struct {
	Tag      string       `json:"tag"`
	Name     string       `json:"name"`
	RankList DateRankList `json:"rankList"`
}
type ClashRoyaleFriendlyLoadReturn struct {
	MyTag     string                    `json:"myTag"`
	FriendTag string                    `json:"friendTag"`
	Wins      int                       `json:"wins"`
	Losses    int                       `json:"losses"`
	Ties      int                       `json:"ties"`
	WinRate   float64                   `json:"winRate"`
	Games     []CRFriendlyDataBaseEntry `json:"games"`
}

// api key stuff is actually a JSON Web Token kinda cool something new

var CR_APIKEY string

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file, using environment variables")
	}
	//set this to the production one when needed now its the dev one this may break a bunch lmao
	CR_APIKEY = os.Getenv("CLASH_ROYALE_JWT_KEY")

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

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("API returned status %d: %s", resp.StatusCode, string(body))
		return nil, fmt.Errorf("api error: status %d", resp.StatusCode)
	}
	var p PlayerProfile
	err = json.NewDecoder(resp.Body).Decode(&p)
	if err != nil {
		return nil, fmt.Errorf("decode failed: %w", err)
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

func createDateRankList(battleList BattleList) *DateRankList {
	var DRList DateRankList
	for i := 0; i < len(battleList); i++ {
		battle := battleList[i]
		if battle.Type == "PvP" {
			var DR StringIntPair
			DR.Text = battle.BattleTime
			DR.Value = battle.Team[0].StartingTrophies
			DRList = append(DRList, DR)
		}
	}
	return &DRList
}
func (c *ClashRoyaleClient) getPlayerBattleLog(playerTag string) (*PlayerBattleLogReturn, error) {
	//this specifically pulls ladder games ie trophy road
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

	var battleLog BattleList
	err = json.NewDecoder(resp.Body).Decode(&battleLog)
	if err != nil {
		return nil, err
	}

	var battleReturn PlayerBattleLogReturn
	//get the name I dont really know if this is the best way but it does consider ones first game being a 2v2
	var name string
	for _, player := range battleLog[0].Team {
		if player.Tag == playerTag {
			name = player.Name
		}
	}

	battleReturn.Tag = playerTag
	battleReturn.Name = name
	battleReturn.RankList = *createDateRankList(battleLog)

	return &battleReturn, nil
}

// similar to get PlayerBattleLog but more simple want to just get the list
func (c *ClashRoyaleClient) fetchBattleLog(playerTag string) (BattleList, error) {
	encodedTag := url.PathEscape(playerTag)
	req, err := http.NewRequest("GET", c.baseURL+"/players/"+encodedTag+"/battlelog", nil)
	if err != nil {
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
		return nil, fmt.Errorf("api error: status %d, body %s", resp.StatusCode, string(body))
	}

	var battleLog BattleList
	err = json.NewDecoder(resp.Body).Decode(&battleLog)
	if err != nil {
		return nil, err
	}
	return battleLog, nil

}

func SyncLadderGames(db *sql.DB, client *ClashRoyaleClient, playerTag string) error {
	mostRecentGame, err := getMostRecentCRGame(db, "ladderGames")
	if err != nil {
		return fmt.Errorf("failed to get most recent game %w", err)
	}

	battleLog, err := client.fetchBattleLog(playerTag)
	if err != nil {
		return fmt.Errorf("failed to fetch battlelog: %w", err)
	}

	for _, battle := range battleLog {
		if battle.Type != "PvP" {
			continue
		}
		if mostRecentGame != "" && battle.BattleTime <= mostRecentGame {
			continue
		}

		for _, player := range battle.Team {
			if player.Tag == playerTag {
				//calculate win is 1 loss is 0 tie is -1(very rare)
				result := -1
				if player.TrophyChange > 0 {
					result = 1
				} else if player.TrophyChange < 0 {
					result = 0
				}
				game := CRLadderGame{
					BattleTime:       battle.BattleTime,
					StartingTrophies: player.StartingTrophies,
					TrophyChange:     player.TrophyChange,
					Result:           result,
				}
				if err := addLadderEntry(db, game); err != nil {
					log.Printf("failed to insert battle %s: %v", battle.BattleTime, err)
				}
			}
		}
	}

	return nil
}
func (c *ClashRoyaleClient) getHeadToHeadBattles(myTag string, friendTag string) (BattleList, error) {
	myEncodedTag := url.PathEscape(myTag)
	req, err := http.NewRequest("GET", c.baseURL+"/players/"+myEncodedTag+"/battlelog", nil)
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

	var battleLog BattleList
	if err := json.NewDecoder(resp.Body).Decode(&battleLog); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	var friendlyBattleList BattleList
	for _, battle := range battleLog {
		if len(battle.Team) == 0 || len(battle.Opponent) == 0 {
			continue
		}
		if battle.Team[0].Tag == myTag && battle.Opponent[0].Tag == friendTag {
			friendlyBattleList = append(friendlyBattleList, battle)
		}
	}

	return friendlyBattleList, nil
}

func (c *ClashRoyaleClient) filterHeadToHeadBattles(myTag string, friendTag string) (BattleList, error) {
	myGamesList, err := c.getHeadToHeadBattles(myTag, friendTag)
	if err != nil {
		log.Println("issue getting head to head battles")
		return nil, err
	}

	friendsGamesList, err := c.getHeadToHeadBattles(friendTag, myTag)
	if err != nil {
		log.Println("issue getting head to head battles")
		return nil, err
	}

	combinedGames := append(myGamesList, friendsGamesList...)

	seen := make(map[string]struct{}, len(combinedGames))
	unique := BattleList{}

	for _, battle := range combinedGames {
		key := battle.BattleTime

		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}
		unique = append(unique, battle)
	}
	return unique, nil

}

func syncFriendlyGames(db *sql.DB, client *ClashRoyaleClient, myTag string, friendTag string) error {
	mostRecentGame, err := getMostRecentCRGame(db, "friendlyGames")
	if err != nil {
		return fmt.Errorf("failed to get most recent game %w", err)
	}
	battleLog, err := client.filterHeadToHeadBattles(myTag, friendTag)
	if err != nil {
		return fmt.Errorf("fauled to fetch friendly battlelog: %w", err)
	}

	for _, battle := range battleLog {
		if mostRecentGame != "" && battle.BattleTime <= mostRecentGame {
			continue
		}
		if len(battle.Team) == 0 || len(battle.Opponent) == 0 {
			continue
		}

		result := -1

		if battle.Team[0].Tag == myTag && battle.Team[0].Crowns > battle.Opponent[0].Crowns {
			result = 1
		} else if battle.Opponent[0].Tag == myTag && battle.Opponent[0].Crowns > battle.Team[0].Crowns {
			result = 1
		} else if battle.Opponent[0].Crowns == battle.Team[0].Crowns {
			result = -1
		} else {
			result = 0
		}

		//game := CRFriendlyDataBaseEntry{
		//	BattleTime: battle.BattleTime,
		//	Result:     result,
		//}

		if err := addFriendlyEntry(db, battle.BattleTime, result, json.RawMessage("null"), json.RawMessage("null")); err != nil {
			log.Printf("failed to insert battle %s: %v", battle.BattleTime, err)
		}

	}
	return nil

}
func loadFriendlyStats(db *sql.DB, client *ClashRoyaleClient, myTag string, friendTag string) (*ClashRoyaleFriendlyLoadReturn, error) {
	// helper to load friendly games gonna be more important later on find matchups
	if err := syncFriendlyGames(db, client, myTag, friendTag); err != nil {
		log.Printf("friendy Sync warning: %v", err)
	}

	games, err := getFriendlyHistory(db)
	if err != nil {
		return nil, fmt.Errorf("failed to get friendly history: %w", err)
	}
	var result ClashRoyaleFriendlyLoadReturn
	for _, g := range games {
		result.Games = append(result.Games, g)
		switch g.Result {
		case 1:
			result.Wins++
		case 0:
			result.Losses++
		default:
			result.Ties++
		}
	}
	result.MyTag = myTag
	result.FriendTag = friendTag
	if len(games) > 0 {
		result.WinRate = float64(result.Wins) / float64(len(games))
	}
	return &result, nil
}

// this might not need the player Tag as this could be the on load without input
// the matchups one should need IDs
/*func returnClashRoyaleStats(playerTag string) (*ClashRoyaleStatsReturn, error) {
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
*/
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
