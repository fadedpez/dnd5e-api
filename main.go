package main

import (
	"log"
	"net/http"
	"time"

	"github.com/fadedpez/dnd5e-api/clients/dnd5e"
)

func main() {
	// Create HTTP client with reasonable timeout
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Create the base D&D 5e API client
	baseClient, err := dnd5e.NewDND5eAPI(&dnd5e.DND5eAPIConfig{
		Client: httpClient,
	})
	if err != nil {
		log.Fatalf("Failed to create D&D 5e API client: %v", err)
	}

	// Wrap with caching layer
	// Using 24-hour TTL since D&D 5e data is static
	cachedClient := dnd5e.NewCachedClient(baseClient, 24*time.Hour)

	// Example usage - the first call will hit the API
	log.Println("Fetching races (first call - API)...")
	races, err := cachedClient.ListRaces()
	if err != nil {
		log.Printf("Error fetching races: %v", err)
	} else {
		log.Printf("Found %d races", len(races))
	}

	// Second call will use cache
	log.Println("Fetching races (second call - cache)...")
	races2, err := cachedClient.ListRaces()
	if err != nil {
		log.Printf("Error fetching races: %v", err)
	} else {
		log.Printf("Found %d races (from cache)", len(races2))
	}

	// Example of fetching specific data
	if len(races) > 0 {
		log.Printf("Fetching details for %s...", races[0].Name)
		race, err := cachedClient.GetRace(races[0].Key)
		if err != nil {
			log.Printf("Error fetching race: %v", err)
		} else {
			log.Printf("Race: %s, Speed: %d", race.Name, race.Speed)
		}
	}

	log.Println("D&D 5e API client with caching initialized successfully")
}