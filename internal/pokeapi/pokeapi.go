package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/dcrauwels/pogodex/internal/pokecache"
)

// type locationArea
type locationArea struct {
	Name string
	Url  string
}

// type pokeapiResponse
type pokeapiResponse struct {
	Count    int
	Next     string
	Previous *string
	Results  []locationArea
}

func GetLocations(u string, c *pokecache.Cache) (pokeapiResponse, error) {
	var locations pokeapiResponse

	// check if u already in cache and return if so
	if value, ok := c.Entry[u]; ok {
		if err := json.Unmarshal(value.Val, &locations); err != nil {
			return locations, fmt.Errorf("error unmarshalling data: %w", err)
		}
		return locations, nil
	}

	// GET request
	res, err := http.Get(u)
	if err != nil {
		return locations, fmt.Errorf("error getting request: %w", err)
	}
	defer res.Body.Close()

	// parse response
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return locations, fmt.Errorf("error parsing response: %w", err)
	}

	// write response to Cache
	c.Add(u, body)

	// unmarshal JSON to struct

	if err = json.Unmarshal(body, &locations); err != nil {
		return locations, fmt.Errorf("error unmarshalling data: %w", err)
	}

	return locations, nil
}
