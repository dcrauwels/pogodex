package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/dcrauwels/pogodex/internal/pokecache"
)

// marker interface
type APIResponse interface {
	locationAreaResponse | locationAreaEncounters // Type constraint.
}

// type pokeapiResponse
type locationAreaResponse struct {
	Count    int     `json:"count"`
	Next     string  `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	} `json:"results"`
}

// type locationAreaEncounters: used in commandExplore
type locationAreaEncounters struct {
	/*
		EncounterMethodRates struct{} `json:"encounter_method_rates"` // NOTE I am functionally ignoring this data during unmarshalling
		GameIndex            int      `json:"game_index"`
		ID                   int      `json:"id"`
		Location             struct{} `json:"location"` // NOTE I am functionally ignoring this data during unmarshalling.
		Name                 string   `json:"name"`
		Names                struct{} `json:"names"` // NOTE I am functionally ignoring this data during unmarshalling.
	*/
	Ignored           struct{} `json:"-"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		Ignored struct{} `json:"-"` // NOTE I am functionally ignoring this data during unmarshalling.

	} `json:"pokemon_encounters"`
}

// generic function for GET request to API, unmarshalling
func GetAPIResource[T APIResponse](u string, c *pokecache.Cache) (T, error) {
	var resource T
	// check if u already in cache and return if so
	if value, ok := c.Entry[u]; ok {
		if err := json.Unmarshal(value.Val, &resource); err != nil {
			return resource, fmt.Errorf("error unmarshalling data: %w", err)
		}
		return resource, nil
	}

	// GET request
	res, err := http.Get(u)
	if err != nil {
		return resource, fmt.Errorf("error getting request: %w", err)
	}
	defer res.Body.Close()

	// parse response
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return resource, fmt.Errorf("error parsing response: %w", err)
	}

	// write response to Cache
	c.Add(u, body)

	// unmarshal JSON to struct
	if err = json.Unmarshal(body, &resource); err != nil {
		return resource, fmt.Errorf("error unmarshalling data: %w", err)
	}

	return resource, nil

}

func GetEncounters(u string, c *pokecache.Cache) (locationAreaEncounters, error) {
	return GetAPIResource[locationAreaEncounters](u, c)
}

func GetLocations(u string, c *pokecache.Cache) (locationAreaResponse, error) {
	return GetAPIResource[locationAreaResponse](u, c)
}
