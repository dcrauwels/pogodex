package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

func GetLocations(u string) (pokeapiResponse, error) {
	var locations pokeapiResponse
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

	// unmarshal JSON to struct

	if err = json.Unmarshal(body, &locations); err != nil {
		return locations, fmt.Errorf("error unmarshalling data: %w", err)
	}

	return locations, nil
}
