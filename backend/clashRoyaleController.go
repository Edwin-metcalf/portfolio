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
	"sort"
	"strconv"
	"strings"

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

type PathOfLegendSeasonResult struct {
	LeagueNumber int `json:"leagueNumber"`
	//tophies 0  int
	//rank null int
	//these seem to be set to 0 and null so dont need em
}

type PlayerProfile struct {
	Tag                             string                   `json:"tag"`
	Name                            string                   `json:"name"`
	Trophies                        int                      `json:"trophies"`
	BestTrophies                    int                      `json:"bestTrophies"`
	Wins                            int                      `json:"wins"`
	Losses                          int                      `json:"losses"`
	BattleCount                     int                      `json:"battleCount"`
	ThreeCrownWins                  int                      `json:"threeCrownWins"`
	Clan                            Clan                     `json:"clan"`
	Arena                           Arena                    `json:"arena"`
	CurrentPathOfLegendSeasonResult PathOfLegendSeasonResult `json:"currentPathOfLegendSeasonResult"`
	LastPathOfLegendSeasonResult    PathOfLegendSeasonResult `json:"lastPathOfLegendSeasonResult"`
	BestPathOfLegendSeasonResult    PathOfLegendSeasonResult `json:"bestPathOfLegendSeasonResult"`
}
type Clan struct {
	Tag  string `json:"tag"`
	Name string `json:"name"`
}
type Arena struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}
type IconUrls struct {
	Medium          string `json:"medium"`
	EvolutionMedium string `json:"evolutionMedium"`
	HeroMedium      string `json:"HeroMedium"`
}
type Card struct {
	Name           string `json:"name"`
	Id             int    `json:"id"`
	EvolutionLevel int    `json:"evolutionLevel"`
	//dont need these right now
	//level int
	//rarity string
	//elixerCost int
	IconUrls IconUrls `json:"iconUrls"`
}
type PlayerProfileReturn struct {
	Tag                             string                   `json:"tag"`
	Name                            string                   `json:"name"`
	Trophies                        int                      `json:"trophies"`
	BestTrophies                    int                      `json:"bestTrophies"`
	Wins                            int                      `json:"wins"`
	Losses                          int                      `json:"losses"`
	WinRate                         float64                  `json:"winRate"`
	BattleCount                     int                      `json:"battleCount"`
	ThreeCrownWins                  int                      `json:"threeCrownWins"`
	Clan                            Clan                     `json:"clan"`
	Arena                           Arena                    `json:"arena"`
	CurrentPathOfLegendSeasonResult PathOfLegendSeasonResult `json:"currentPathOfLegendSeasonResult"`
	LastPathOfLegendSeasonResult    PathOfLegendSeasonResult `json:"lastPathOfLegendSeasonResult"`
	BestPathOfLegendSeasonResult    PathOfLegendSeasonResult `json:"bestPathOfLegendSeasonResult"`
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
	Cards            []Card `json:"cards"`
}
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
	MyTag         string                    `json:"myTag"`
	FriendTag     string                    `json:"friendTag"`
	MyName        string                    `json:"myName"`
	FriendName    string                    `json:"friendName"`
	Wins          int                       `json:"wins"`
	Losses        int                       `json:"losses"`
	Ties          int                       `json:"ties"`
	WinRate       float64                   `json:"winRate"`
	Games         []CRFriendlyDataBaseEntry `json:"games"`
	CardWinRates  map[string]*CardRecord    `json:"cardWinRates"`
	MyMatchup     PlayerCardMatchup         `json:"myMatchup"`
	FriendMatchup PlayerCardMatchup         `json:"friendMatchup"`
}
type PlayerCardMatchup struct {
	WithCards    map[string]*CardRecord `json:"withCards"`
	AgainstCards map[string]*CardRecord `json:"againstCards"`
}
type DeckCardWinRates struct {
	Deck      []Card                 `json:"deck"`
	CardStats map[string]*CardRecord `json:"cardStats"`
}
type CRRankedLoadReturn struct {
	Games            []CRRankedDataBaseEntry `json:"games"`
	Wins             int                     `json:"wins"`
	Losses           int                     `json:"losses"`
	Ties             int                     `json:"ties"`
	WinRate          float64                 `json:"winRate"`
	DeckCardWinRates []DeckCardWinRates      `json:"deckCardWinRates"`
}
type CardRecord struct {
	Wins    int     `json:"wins"`
	Losses  int     `json:"losses"`
	WinRate float64 `json:"winRate"`
}
type GameResult struct {
	MyDeck    []Card `json:"myDeck"`
	EnemyDeck []Card `json:"enemyDeck"`
	Result    int    `json:"result"`
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
	fmt.Printf("current POL %v", p.CurrentPathOfLegendSeasonResult)
	fmt.Printf("last POL %v", p.LastPathOfLegendSeasonResult)
	fmt.Printf("best POL %v", p.BestPathOfLegendSeasonResult)

	retProfile.CurrentPathOfLegendSeasonResult = p.CurrentPathOfLegendSeasonResult
	retProfile.LastPathOfLegendSeasonResult = p.LastPathOfLegendSeasonResult
	retProfile.BestPathOfLegendSeasonResult = p.BestPathOfLegendSeasonResult

	return &retProfile, err
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

func syncRankedGames(db *sql.DB, client *ClashRoyaleClient, playerTag string) error {
	mostRecentGame, err := getMostRecentCRGame(db, "ladderGames")
	if err != nil {
		return fmt.Errorf("failed to get most recent game %w", err)
	}

	battleLog, err := client.fetchBattleLog(playerTag)
	if err != nil {
		return fmt.Errorf("failed to fetch battleog: %w", err)
	}

	for _, battle := range battleLog {
		if battle.Type != "pathOfLegend" {
			continue
		}
		if mostRecentGame != "" && battle.BattleTime <= mostRecentGame {
			continue
		}
		if len(battle.Team) == 0 || len(battle.Opponent) == 0 {
			continue
		}

		result := -1
		if battle.Team[0].Crowns > battle.Opponent[0].Crowns {
			result = 1
		} else if battle.Team[0].Crowns < battle.Opponent[0].Crowns {
			result = 0
		}
		myDeckJson, err := json.Marshal(battle.Team[0].Cards)
		if err != nil {
			log.Printf("failed to marshal my deck: %v", err)
		}
		enemyDeckJson, err := json.Marshal(battle.Opponent[0].Cards)
		if err != nil {
			log.Printf("failed to marshal enemy deck: %v", err)
		}

		if err := addRankedEntry(db, battle.BattleTime, result, myDeckJson, enemyDeckJson); err != nil {
			log.Printf("failed to insert battle %s: %v", battle.BattleTime, err)
		}
	}

	return nil
}

func loadRankedStats(db *sql.DB, client *ClashRoyaleClient, playerTag string) (*CRRankedLoadReturn, error) {
	if err := syncRankedGames(db, client, playerTag); err != nil {
		log.Printf("Ranked SYNC warning: %v", err)
	}

	games, err := getRankedHistory(db)
	if err != nil {
		return nil, fmt.Errorf("failed to get ranked history: %w", err)
	}
	var result CRRankedLoadReturn
	deckToGamesMap := make(map[string]*deckGroup)

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

		var myCards []Card
		if err := json.Unmarshal(g.MyDeck, &myCards); err != nil {
			log.Printf("failed to unmarshal my deck: %v", err)
			continue
		}

		var enemyCards []Card
		if err := json.Unmarshal(g.EnemyDeck, &enemyCards); err != nil {
			log.Printf("failed to unmarshal enemy deck: %v", err)
			continue
		}
		if len(myCards) == 0 {
			continue
		}

		key := deckKey(myCards)
		group, ok := deckToGamesMap[key]
		if !ok {
			group = &deckGroup{deck: myCards}
			deckToGamesMap[key] = group
		}

		group.games = append(group.games, GameResult{
			MyDeck:    myCards,
			EnemyDeck: enemyCards,
			Result:    g.Result,
		})
	}
	result.DeckCardWinRates = deckToCardWinRate(deckToGamesMap)

	if len(games) > 0 {
		result.WinRate = float64(result.Wins) / float64(len(games))
	}

	return &result, nil
}

type deckGroup struct {
	deck  []Card
	games []GameResult
}

func deckKey(cards []Card) string {
	//this makes the combination of card IDs into a string we can use as a key, need to sort to make sure
	//two decks in different order are the same ID key
	sorted := make([]Card, len(cards))

	copy(sorted, cards)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].Id < sorted[j].Id })

	ids := make([]string, len(sorted))

	for i, c := range sorted {
		ids[i] = strconv.Itoa(c.Id)
	}
	return strings.Join(ids, "-")
}
func deckToCardWinRate(deckToGamesMap map[string]*deckGroup) []DeckCardWinRates {

	var result []DeckCardWinRates
	for _, group := range deckToGamesMap {
		cardStats := winRateByCard(group.games, func(g GameResult) []Card { return g.EnemyDeck })
		result = append(result, DeckCardWinRates{
			Deck:      group.deck,
			CardStats: cardStats,
		})
	}

	return result
}
func winRateByCard(games []GameResult, side func(GameResult) []Card) map[string]*CardRecord {
	// if side is func(g) { return g.EnemyDeck } then it means thats my win rate against the card
	// if side is func(g) { return g.MyDeck } then it is for each of my cards that is my win rate with it
	records := make(map[string]*CardRecord)

	for _, g := range games {
		for _, c := range side(g) {
			rec, ok := records[c.Name]
			if !ok {
				rec = &CardRecord{}
				records[c.Name] = rec
			}

			if g.Result == 1 {
				rec.Wins += 1
			} else if g.Result == 0 {
				rec.Losses += 1
			}
		}
	}
	for _, rec := range records {
		total := rec.Wins + rec.Losses
		if total > 0 {
			rec.WinRate = float64(rec.Wins) / float64(total)
		}
	}
	return records
}
func buildPlayerMatchup(games []GameResult) PlayerCardMatchup {
	return PlayerCardMatchup{
		WithCards:    winRateByCard(games, func(g GameResult) []Card { return g.MyDeck }),
		AgainstCards: winRateByCard(games, func(g GameResult) []Card { return g.EnemyDeck }),
	}
}
func flipGames(games []GameResult) []GameResult {
	flipped := make([]GameResult, len(games))

	for i, g := range games {
		r := g.Result

		if r == 1 {
			r = 0
		} else if r == 0 {
			r = 1
		}
		flipped[i] = GameResult{MyDeck: g.EnemyDeck, EnemyDeck: g.MyDeck, Result: r}
	}
	return flipped
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
		return fmt.Errorf("failed to fetch friendly battlelog: %w", err)
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

		var myCards, enemyCards []Card

		if battle.Team[0].Tag == myTag {
			myCards = battle.Team[0].Cards
			enemyCards = battle.Opponent[0].Cards
		} else {
			myCards = battle.Opponent[0].Cards
			enemyCards = battle.Team[0].Cards
		}

		myDeckJson, err := json.Marshal(myCards)
		if err != nil {
			log.Printf("failed to marshal my deck: %v", err)
		}
		enemyDeckJson, err := json.Marshal(enemyCards)
		if err != nil {
			log.Printf("failed to marshal enemy deck: %v", err)
		}

		if err := addFriendlyEntry(db, battle.BattleTime, result, myDeckJson, enemyDeckJson); err != nil {
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

	var gameResults []GameResult
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
		var myCards []Card
		if err := json.Unmarshal(g.MyDeck, &myCards); err != nil {
			log.Printf("failed to unmarshal my deck: %v", err)
			continue
		}

		var enemyCards []Card
		if err := json.Unmarshal(g.EnemyDeck, &enemyCards); err != nil {
			log.Printf("failed to unmarshal enemy deck: %v", err)
			continue
		}
		if len(myCards) == 0 {
			continue
		}
		gameResults = append(gameResults, GameResult{MyDeck: myCards, EnemyDeck: enemyCards, Result: g.Result})
	}

	result.MyTag = myTag
	result.FriendTag = friendTag
	if len(games) > 0 {
		result.WinRate = float64(result.Wins) / float64(len(games))
	}

	result.CardWinRates = winRateByCard(gameResults, func(g GameResult) []Card { return g.MyDeck })
	result.MyMatchup = buildPlayerMatchup(gameResults)
	result.FriendMatchup = buildPlayerMatchup(flipGames(gameResults))

	return &result, nil
}

func matchupGeneratorhelper(client *ClashRoyaleClient, tag1 string, tag2 string) (*ClashRoyaleFriendlyLoadReturn, error) {
	battleLog, err := client.filterHeadToHeadBattles(tag1, tag2)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch friendly battlelog: %w", err)
	}
	if len(battleLog) == 0 {
		return nil, fmt.Errorf("no battles found between %s and %s", tag1, tag2)
	}

	numGames := len(battleLog)
	var matchupReturn ClashRoyaleFriendlyLoadReturn
	gamesList := make([]CRFriendlyDataBaseEntry, numGames)
	var wins int
	var losses int
	var ties int
	myName := findPlayerName(battleLog[0], tag1)
	friendName := findPlayerName(battleLog[0], tag2)

	var gameResults []GameResult

	for i, battle := range battleLog {
		if len(battle.Team) == 0 || len(battle.Opponent) == 0 {
			continue
		}
		gamesList[i].BattleTime = battle.BattleTime
		var myCards, enemyCards []Card

		if battle.Team[0].Tag == tag1 {
			myCards = battle.Team[0].Cards
			enemyCards = battle.Opponent[0].Cards
		} else {
			myCards = battle.Opponent[0].Cards
			enemyCards = battle.Team[0].Cards
		}

		myDeckBytes, err := json.Marshal(battle.Team[0].Cards)
		if err != nil {
			log.Printf("failed to marshal my deck: %v", err)
		}
		enemyDeckBytes, err := json.Marshal(battle.Opponent[0].Cards)
		if err != nil {
			log.Printf("failed to marshal enemy deck: %v", err)
		}
		//			THIS NEEDS TO CHANGE TO UNMARSHAL THE BYTES THEN DO THE WIN RATE BY CARD

		gamesList[i].MyDeck = myDeckBytes
		gamesList[i].EnemyDeck = enemyDeckBytes

		result := -1

		if battle.Team[0].Tag == tag1 && battle.Team[0].Crowns > battle.Opponent[0].Crowns {
			result = 1
			wins += 1
		} else if battle.Opponent[0].Tag == tag1 && battle.Opponent[0].Crowns > battle.Team[0].Crowns {
			result = 1
			wins += 1
		} else if battle.Opponent[0].Crowns == battle.Team[0].Crowns {
			result = -1
			ties += 1
		} else {
			result = 0
			losses += 1
		}

		gamesList[i].Result = result

		gameResults = append(gameResults, GameResult{MyDeck: myCards, EnemyDeck: enemyCards, Result: result})

	}

	matchupReturn.CardWinRates = winRateByCard(gameResults, func(g GameResult) []Card { return g.MyDeck })

	matchupReturn.MyTag = tag1
	matchupReturn.FriendTag = tag2
	matchupReturn.MyName = myName
	matchupReturn.FriendName = friendName
	matchupReturn.Wins = wins
	matchupReturn.Losses = losses
	matchupReturn.Ties = ties
	if wins+losses > 0 {
		matchupReturn.WinRate = float64(wins) / float64(wins+losses)
	}
	matchupReturn.Games = gamesList

	return &matchupReturn, nil

}
func findPlayerName(battle Battle, tag string) string {
	for _, p := range battle.Team {
		if p.Tag == tag {
			return p.Name
		}
	}
	for _, p := range battle.Opponent {
		if p.Tag == tag {
			return p.Name
		}
	}
	return ""
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
