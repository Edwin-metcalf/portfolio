package main

import (
	"database/sql"
	"encoding/json"
	"testing"
	"time"
)

func mustParseBattleTime(t *testing.T, value string) time.Time {
	t.Helper()
	parsed, err := parseBattleTime(value)
	if err != nil {
		t.Fatalf("parseBattleTime(%q): %v", value, err)
	}
	return parsed
}

func TestFirstMondayAt(t *testing.T) {
	cases := []struct {
		year  int
		month time.Month
		want  int
	}{
		{2026, time.September, 7},
		{2026, time.October, 5},
		{2026, time.November, 2},
	}

	for _, tc := range cases {
		got := firstMondayAt(tc.year, tc.month)
		if got.Day() != tc.want || got.Weekday() != time.Monday {
			t.Errorf("firstMondayAt(%d, %s) = %s, want day %d (Monday)", tc.year, tc.month, got, tc.want)
		}
		if got.Hour() != crSeasonStartHourUTC {
			t.Errorf("firstMondayAt(%d, %s) hour = %d, want %d", tc.year, tc.month, got.Hour(), crSeasonStartHourUTC)
		}
	}
}

func TestSeasonIDForBattleTime(t *testing.T) {
	cases := []struct {
		battleTime string
		want       string
	}{
		{"20260916T222642.000Z", "2026-09"},
		{"20261031T235959.000Z", "2026-10"},
		{"20261101T225717.000Z", "2026-10"},
		{"20261102T095959.000Z", "2026-10"},
		{"20261102T100000.000Z", "2026-11"},
		{"20261102T170850.000Z", "2026-11"},
	}

	for _, tc := range cases {
		got := seasonIDForTime(mustParseBattleTime(t, tc.battleTime))
		if got != tc.want {
			t.Errorf("seasonIDForTime(%s) = %s, want %s", tc.battleTime, got, tc.want)
		}
	}
}

func TestSeasonLabel(t *testing.T) {
	if got := seasonLabel("2026-11"); got != "November 2026" {
		t.Errorf("seasonLabel(2026-11) = %q, want %q", got, "November 2026")
	}
	if got := seasonLabel("bogus"); got != "bogus" {
		t.Errorf("seasonLabel(bogus) = %q, want %q", got, "bogus")
	}
}

func rankedTestGame(t *testing.T, battleTime string, result int) CRRankedDataBaseEntry {
	t.Helper()

	myDeck, err := json.Marshal([]Card{{Id: 1, Name: "Knight"}, {Id: 2, Name: "Archers"}})
	if err != nil {
		t.Fatal(err)
	}
	enemyDeck, err := json.Marshal([]Card{{Id: 3, Name: "Hog Rider"}, {Id: 4, Name: "Fireball"}})
	if err != nil {
		t.Fatal(err)
	}

	return CRRankedDataBaseEntry{
		BattleTime: battleTime,
		Result:     result,
		MyDeck:     myDeck,
		EnemyDeck:  enemyDeck,
	}
}

func TestAggregateRankedStatsBySeason(t *testing.T) {
	games := []CRRankedDataBaseEntry{
		rankedTestGame(t, "20260916T222642.000Z", 1),
		rankedTestGame(t, "20261020T120000.000Z", 0),
		rankedTestGame(t, "20261021T120000.000Z", 1),
		rankedTestGame(t, "20261102T170850.000Z", 1), // after the Nov 2 10:00 UTC reset
	}

	result := aggregateRankedStats(games)

	if result.Wins != 3 || result.Losses != 1 || result.Ties != 0 {
		t.Fatalf("totals = %dW %dL %dT, want 3W 1L 0T", result.Wins, result.Losses, result.Ties)
	}
	if result.WinRate != 0.75 {
		t.Errorf("win rate = %v, want 0.75", result.WinRate)
	}

	byID := make(map[string]CRRankedSeason)
	for _, season := range result.Seasons {
		byID[season.Id] = season
	}

	september, ok := byID["2026-09"]
	if !ok {
		t.Fatal("missing 2026-09 season")
	}
	if september.Wins != 1 || september.Losses != 0 || september.Label != "September 2026" {
		t.Errorf("2026-09 = %+v, want 1W 0L September 2026", september)
	}

	october, ok := byID["2026-10"]
	if !ok {
		t.Fatal("missing 2026-10 season")
	}
	if october.Wins != 1 || october.Losses != 1 {
		t.Errorf("2026-10 = %dW %dL, want 1W 1L", october.Wins, october.Losses)
	}
	if october.WinRate != 0.5 {
		t.Errorf("2026-10 win rate = %v, want 0.5", october.WinRate)
	}

	november, ok := byID["2026-11"]
	if !ok {
		t.Fatal("missing 2026-11 season")
	}
	if november.Wins != 1 || november.Losses != 0 {
		t.Errorf("2026-11 = %dW %dL, want 1W 0L", november.Wins, november.Losses)
	}

	if len(result.Seasons) < 3 {
		t.Fatalf("got %d seasons, want at least 3", len(result.Seasons))
	}
	for i := 1; i < len(result.Seasons); i++ {
		if result.Seasons[i-1].Id <= result.Seasons[i].Id {
			t.Errorf("seasons not sorted newest first: %s then %s", result.Seasons[i-1].Id, result.Seasons[i].Id)
		}
	}

	if result.CurrentSeasonID == "" {
		t.Error("CurrentSeasonID is empty")
	}
	if _, ok := byID[result.CurrentSeasonID]; !ok {
		t.Errorf("current season %s missing from seasons", result.CurrentSeasonID)
	}
}

func TestAddRankedEntryIgnoresDuplicateBattleTime(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE rankedGames (
		id INTEGER PRIMARY KEY,
		battle_time TEXT NOT NULL UNIQUE,
		result TINYINT,
		my_deck TEXT NOT NULL,
		enemy_deck TEXT NOT NULL
	);`)
	if err != nil {
		t.Fatal(err)
	}

	deck := json.RawMessage(`[]`)
	if err := addRankedEntry(db, "20261102T170850.000Z", 1, deck, deck); err != nil {
		t.Fatal(err)
	}
	if err := addRankedEntry(db, "20261102T170850.000Z", 1, deck, deck); err != nil {
		t.Fatalf("duplicate insert should be ignored, got: %v", err)
	}

	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM rankedGames").Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Errorf("rankedGames count = %d, want 1", count)
	}
}
