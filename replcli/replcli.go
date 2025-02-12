package replcli

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"

	"github.com/dcrauwels/pogodex/internal/pokeapi"
	"github.com/dcrauwels/pogodex/internal/pokecache"
	"github.com/dcrauwels/pogodex/internal/stringutils"
)

// define struct for commands
type cliCommand struct {
	name        string
	description string
	callback    func(argument ...string) error
	argument    string // currently just the single string, but maybe a slice is better?
}

type Pokemon struct {
	name   string
	id     int
	height int
	weight int
	stats  map[string]int

	types []string
}

// define struct for REPL and associated function
type REPL struct {
	commands    map[string]cliCommand
	nextURL     string // store what URL we are at when using commandMap
	previousURL string // see above
	cache       *pokecache.Cache
	pokemon     map[string]Pokemon
}

func NewREPL(interval int) *REPL {

	return &REPL{
		commands:    make(map[string]cliCommand),
		nextURL:     "",
		previousURL: "",
		cache:       pokecache.NewCache(interval),
		pokemon:     make(map[string]Pokemon),
	}
}

// add a command to the REPL, for use in the ReplCLI function to not have this huge definition at the start
func (r *REPL) RegisterCommand(name string, description string, callback func(argument ...string) error, argument string) {
	r.commands[name] = cliCommand{
		name:        name,
		description: description,
		callback:    callback,
		argument:    argument,
	}
}

// first command: 'exit' pokedex
func (r *REPL) commandExit(argument ...string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")

	// any cleanup functions (that might produce errors) go here

	os.Exit(0)
	return nil
}

// second command: print 'help'
func (r *REPL) commandHelp(argument ...string) error {
	// header fluff
	fmt.Println("Welcome to the Pokedex!\nUsage:")

	// sanity check
	if len(r.commands) < 1 {
		return fmt.Errorf("no commands are implemented (other than this one)")
	}

	// loop over all commands
	for _, command := range r.commands {
		fmt.Printf(" %s: %s.\n  Arguments: %s\n", command.name, command.description, command.argument)
	}
	return nil
}

// third command: print area names on the 'map'
func (r *REPL) commandMap(argument ...string) error {

	// construct URL
	var u string
	if r.nextURL != "" {
		u = r.nextURL // this means we have a URL stored from previous commandMap() calls
	} else {
		u = "https://pokeapi.co/api/v2/location-area/"
	}

	locations, err := pokeapi.GetLocations(u, r.cache)
	if err != nil {
		return fmt.Errorf("error getting locations from pokeAPI: %w", err)
	}

	// update r.nextURL, r.previousURL
	if locations.Previous != nil {
		previousValue := *locations.Previous
		r.previousURL = previousValue
	}
	r.nextURL = locations.Next

	// print locations
	for _, location := range locations.Results {
		fmt.Println(location.Name)
	}

	return nil
}

// fourth command: reverse the map
func (r *REPL) commandMapb(argument ...string) error {
	// construct URL and sanity check
	if r.previousURL == "" || r.nextURL == "https://pokeapi.co/api/v2/location-area/?offset=20&limit=20" {
		fmt.Println("You're on the first page of location results")
		return nil
	}
	locations, err := pokeapi.GetLocations(r.previousURL, r.cache)
	if err != nil {
		return fmt.Errorf("error getting locations from pokeAPI: %w", err)
	}

	// update r.nextURL, r.previousURL
	if locations.Previous != nil {
		previousValue := *locations.Previous
		r.previousURL = previousValue
	}
	r.nextURL = locations.Next

	// print locations
	for _, location := range locations.Results {
		fmt.Println(location.Name)
	}

	return nil
}

// fifth command: list pokemon in a given location-area
func (r *REPL) commandExplore(argument ...string) error {
	// sanity check
	if len(argument) == 0 {
		return fmt.Errorf("no arguments passed")
	}

	// just take the first argument
	a := argument[0]

	// print opening message
	fmt.Printf("Exploring %s...\n", a)

	// construct url
	url := "https://pokeapi.co/api/v2/location-area/" + a

	// get data via pokeapi package
	encounters, err := pokeapi.GetEncounters(url, r.cache)
	if err != nil {
		return fmt.Errorf("error getting encounters from PokeAPI: %w", err)
	}

	// print pokemon
	for _, e := range encounters.PokemonEncounters {
		fmt.Printf("- %s\n", e.Pokemon.Name)
	}

	return nil
}

// sixth command: attempt to catch a given pokemon
func (r *REPL) commandCatch(argument ...string) error {
	// sanity check
	if len(argument) == 0 {
		return fmt.Errorf("no arguments passed")
	}

	// just take first argument
	a := argument[0]

	//print opening message
	fmt.Printf("Throwing a Pokeball at %s...\n", a)

	//construct url
	url := "https://pokeapi.co/api/v2/pokemon/" + a

	// get data via pokeapi package
	pokemon, err := pokeapi.GetPokemon(url, r.cache)
	if err != nil {
		return fmt.Errorf("error getting Pokemon from PokeAPI: %w", err)
	}

	// calculate catch result
	maxExp := 750 // Found a forum post from 2015 or so saying Blissey is highest baseexp at 600 or so ... not sure if correct
	catchRand := rand.Intn(750)
	baseExp := pokemon.BaseExperience
	if baseExp > maxExp {
		return fmt.Errorf("higher exp than assumed maximum %d found", maxExp)
	}

	// print catch result
	if catchRand >= baseExp {
		fmt.Printf("%s was caught!\n", a)

		// extract stats to dict
		extractedStats := make(map[string]int)
		for _, s := range pokemon.Stats {
			extractedStats[s.Stat.Name] = s.BaseStat
		}

		// extract types to slice
		extractedTypes := make([]string, 0, len(pokemon.Types))
		for _, t := range pokemon.Types {
			extractedTypes = append(extractedTypes, t.Type.Name)
		}

		// and add to pokedex
		r.pokemon[a] = Pokemon{
			name:   pokemon.Name,
			id:     pokemon.ID,
			height: pokemon.Height,
			weight: pokemon.Weight,
			stats:  extractedStats,
			types:  extractedTypes,
		}
	} else {
		fmt.Printf("%s escaped!\n", a)
	}

	return nil
}

// main function: open a CLI that loops until interrupt or commandExit() is called
func (r *REPL) ReplCLI() {
	// register commands here
	r.RegisterCommand("exit", "Exit the Pokedex", r.commandExit, "none")
	r.RegisterCommand("help", "Prints this help message", r.commandHelp, "none")
	r.RegisterCommand("map", "View map locations", r.commandMap, "none")
	r.RegisterCommand("mapb", "View previous map locations", r.commandMapb, "none")
	r.RegisterCommand("explore", "View list of wild Pokemon on a given map location", r.commandExplore, "<area-name>")
	r.RegisterCommand("catch", "Attempt to catch a Pokemon", r.commandCatch, "<pokemon-name>")

	// initialize scanner
	s := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")

		// take input and clean it
		s.Scan()
		input := s.Text()
		if ok := len(input) > 0; !ok {
			continue
		}
		cleanedInput := stringutils.CleanInput(input)
		commandInput := cleanedInput[0]
		var argumentInput []string // we'll use this if more than 1 word is input

		// check if an argument was passed
		if len(cleanedInput) != 1 { // this means there were at least two words in the input, i.e. an argument was passed
			argumentInput = cleanedInput[1:]
		}

		//check if command exists (i.e. is present in slice r.commands)
		command, ok := r.commands[commandInput]
		if !ok {
			fmt.Printf("Unknown command: %s\n", commandInput)
			continue
		}

		// try executing the callback corresponding to the command input, raise error if an issue arises
		if err := command.callback(argumentInput...); err != nil {
			fmt.Printf("Error executing command: %v\n", err)
			continue
		}
	}
}
