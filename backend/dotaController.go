package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/joho/godotenv"
	//"golang.org/x/tools/playground"
)

//check appflowy for your stuff

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
	numberOfGames := 10
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

// gonna need a struct for the return value of recent heros played
// need to figure out how we want to send the data to the front end
// prolly want nested maps tbh
// ex nestedmap := make(map[string]map[string]int)
type heroStats struct {
	Wins    int       `json:"wins"`
	Losses  int       `json:"losses"`
	Games   int       `json:"games"`
	Kills   int       `json:"kills"`
	Deaths  int       `json:"deaths"`
	Assists int       `json:"assists"`
	AvgKDA  []float32 `json:"avgKDA"`
}

func recentHerosPlayed(rawData []recentGame) (*map[string]heroStats, error) {
	heroMap, err := getHeroes()

	if err != nil {
		return nil, fmt.Errorf("failed to get heroes: %w", err)
	}

	heroMapStats := make(map[string]heroStats)

	for _, game := range rawData {
		var won bool
		//Check if the game was won or not and update correctly
		if game.PlayerSlot < 128 {
			won = game.RadiantWin
		} else {
			won = !game.RadiantWin
		}
		heroName := heroMap[game.HeroID]
		stats := heroMapStats[heroName]
		if stats.Games == 0 {
			stats.AvgKDA = make([]float32, 3)
		}
		stats.Games += 1

		if won {
			stats.Wins += 1
		} else {
			stats.Losses += 1
		}

		//KDA
		stats.Kills += game.Kills
		stats.Deaths += game.Deaths
		stats.Assists += game.Assists

		heroMapStats[heroName] = stats
	}
	// avg it out I know thing a list mighta been better but oh well iterate through this guy
	for _, stats := range heroMapStats {
		stats.AvgKDA[0] = float32(stats.Kills / stats.Games)
		stats.AvgKDA[1] = float32(stats.Deaths / stats.Games)
		stats.AvgKDA[2] = float32(stats.Assists / stats.Games)
	}

	return &heroMapStats, nil
}

type Hero struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	LocalizedName string `json:"localized_name"`
}

// this is for getting the hero names just hit the website for them
func getHeroes() (map[int]string, error) {
	url := "https://api.opendota.com/api/heroes"

	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var heroes []Hero

	err = json.NewDecoder(response.Body).Decode(&heroes)
	if err != nil {
		return nil, err
	}

	heroMap := make(map[int]string)

	for _, hero := range heroes {
		heroMap[hero.ID] = hero.LocalizedName
	}
	return heroMap, nil
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
	cleanedGames.AvgKDA[0] = float32(totalKda[0]) / float32(numberOfGames)
	cleanedGames.AvgKDA[1] = float32(totalKda[1]) / float32(numberOfGames)
	cleanedGames.AvgKDA[2] = float32(totalKda[2]) / float32(numberOfGames)

	cleanedGames.WinRate = float32(cleanedGames.Wins) / float32(cleanedGames.Wins+cleanedGames.Losses)
	return &cleanedGames, nil
}

type MatchDetail struct {
	Players []struct {
		HeroID     int `json:"hero_id"`
		PlayerSlot int `json:"player_slot"`
	} `json:"players"`
}

func getMatchDetails(matchID int) (*MatchDetail, error) {
	url := fmt.Sprintf("https://api.opendota.com/api/matches/%d", matchID)

	response, err := http.Get(url)

	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var matchDetail MatchDetail
	err = json.NewDecoder(response.Body).Decode(&matchDetail)
	if err != nil {
		return nil, err
	}

	return &matchDetail, nil
}

func getEnemyHeroes(match *MatchDetail, yourTeam string, heroMap map[int]string) []string {
	var enemyHeroes []string

	for _, player := range match.Players {
		playerTeam := "radiant"
		if player.PlayerSlot > 128 {
			playerTeam = "dire"
		}

		if playerTeam != yourTeam {
			enemyHeroes = append(enemyHeroes, heroMap[player.HeroID])
		}
	}
	return enemyHeroes
}

type matchupStats struct {
	EnemyHeroName string  `json:"enemyHeroName"`
	Wins          int     `json:"wins"`
	Losses        int     `json:"losses"`
	Games         int     `json:"games"`
	WinRate       float32 `json:"winRate"`
}

func getMatchupStats(rawData []recentGame, heroMap map[int]string) (*map[string]matchupStats, error) {
	matchupMap := make(map[string]matchupStats)

	for _, game := range rawData {
		var won bool
		var yourTeam string
		if game.PlayerSlot < 128 {
			won = game.RadiantWin
			yourTeam = "radiant"

		} else {
			won = !game.RadiantWin
			yourTeam = "dire"
		}

		//yourHero := heroMap[game.HeroID]

		matchDetails, err := getMatchDetails(game.MatchID)
		if err != nil {
			log.Printf("Failed to get match details for %d: %v", game.MatchID, err)
			continue
		}

		enemyHeroes := getEnemyHeroes(matchDetails, yourTeam, heroMap)

		for _, enemyHero := range enemyHeroes {
			stats := matchupMap[enemyHero]
			stats.EnemyHeroName = enemyHero
			stats.Games += 1

			if won {
				stats.Wins += 1
			} else {
				stats.Losses += 1
			}

			matchupMap[enemyHero] = stats
		}
	}

	for key, stats := range matchupMap {
		if stats.Games > 0 {
			stats.WinRate = float32(stats.Wins) / float32(stats.Games)
		}
		matchupMap[key] = stats
	}

	return &matchupMap, nil
}

type dotaStatsReturn struct {
	Wins            int                     `json:"wins"`
	Losses          int                     `json:"losses"`
	WinRate         float32                 `json:"winRate"`
	RecentWins      int                     `json:"recentWins"`
	RecentLosses    int                     `json:"recentLosses"`
	RecentWinRate   float32                 `json:"recentWinRate"`
	AvgKDA          []float32               `json:"avgKDA"`
	RecentHeroStats map[string]heroStats    `json:"recentHeroStats"`
	MatchupStats    map[string]matchupStats `json:"matchupStats"`
}

func returnGamesSats(playerID string) (*dotaStatsReturn, error) {
	var wg sync.WaitGroup

	var winLoss *winLossRate
	var recentGames *[]recentGame
	var heroMap map[int]string
	var winLossErr, recentGamesErr, heroMapErr error

	wg.Go(func() {
		var err error
		winLoss, err = getWinLose(playerID)
		winLossErr = err
	})

	wg.Go(func() {
		var err error
		recentGames, err = getRecentGames(playerID)
		recentGamesErr = err
	})

	wg.Go(func() {
		var err error
		heroMap, err = getHeroes()
		heroMapErr = err
	})

	wg.Wait()

	if winLossErr != nil {
		return nil, fmt.Errorf("failed to get win/loss stats: %w", winLossErr)
	}

	if recentGamesErr != nil {
		return nil, fmt.Errorf("failed to get recent games: %w", recentGamesErr)
	}

	//for the enemy stuff
	if heroMapErr != nil {
		return nil, fmt.Errorf("Failed to get heroes: %w", heroMapErr)
	}

	recentStats, err := processRawRecentGames(*recentGames)
	if err != nil {
		return nil, fmt.Errorf("failed to process recent games: %w", err)
	}

	recentHeroStats, err := recentHerosPlayed(*recentGames)
	if err != nil {
		return nil, fmt.Errorf("failed to process recent hero stats: %w", err)
	}

	matchupStats, err := getMatchupStats(*recentGames, heroMap)
	if err != nil {
		return nil, fmt.Errorf("failed to get matchup stats: %w", err)

	}
	stats := &dotaStatsReturn{
		Wins:            winLoss.Wins,
		Losses:          winLoss.Losses,
		WinRate:         winLoss.WinRate,
		RecentWins:      recentStats.Wins,
		RecentLosses:    recentStats.Losses,
		RecentWinRate:   recentStats.WinRate,
		AvgKDA:          recentStats.AvgKDA,
		RecentHeroStats: *recentHeroStats,
		MatchupStats:    *matchupStats,
	}

	return stats, nil
}

/*
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
*/
