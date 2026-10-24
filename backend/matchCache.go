package main

import (
	"log"
	"sync"
	"time"
)

type matchCacheEntry struct {
	match     *MatchDetail
	expiresAt time.Time
}

var (
	matchCache   = make(map[int]*matchCacheEntry)
	matchCacheMu sync.RWMutex
	// Tune this based off your usage. Match details are immutable so a longer TTL is fine
	matchCacheTTL = 30 * time.Minute
)

// StartMatchCacheJanitor starts a background goroutine that purges expired cache entries
// every `interval`. Call this once at startup (e.g. in main()).
func StartMatchCacheJanitor(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for range ticker.C {
			now := time.Now()
			evicted := 0
			matchCacheMu.Lock()
			for id, entry := range matchCache {
				if now.After(entry.expiresAt) {
					delete(matchCache, id)
					evicted++
				}
			}
			remaining := len(matchCache)
			matchCacheMu.Unlock()
			if evicted > 0 {
				log.Printf("[CACHE] janitor evicted=%d remaining=%d", evicted, remaining)
			}
		}
	}()
}

func init() {
	// Start the janitor so expired entries don't live forever
	StartMatchCacheJanitor(5 * time.Minute)
}

// GetMatchDetailsCached fetches match details with a TTL based in-memory cache
func GetMatchDetailsCached(id int) (*MatchDetail, error) {
	// Fast path: check cache with read lock
	matchCacheMu.RLock()
	if entry, ok := matchCache[id]; ok {
		if time.Now().Before(entry.expiresAt) {
			match := entry.match
			matchCacheMu.RUnlock()
			log.Printf("[CACHE] hit matchID=%d", id)
			return match, nil
		}
	}
	matchCacheMu.RUnlock()

	// Cache miss or expired: fetch from API
	log.Printf("[CACHE] miss matchID=%d", id)
	match, err := getMatchDetails(id)
	if err != nil {
		return nil, err
	}

	// Store in cache with new TTL
	matchCacheMu.Lock()
	matchCache[id] = &matchCacheEntry{
		match:     match,
		expiresAt: time.Now().Add(matchCacheTTL),
	}
	matchCacheMu.Unlock()

	return match, nil
}
