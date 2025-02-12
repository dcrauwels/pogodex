package replcli

import (
	"bufio"
	"fmt"
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

// define struct for REPL and associated function
type REPL struct {
	commands    map[string]cliCommand
	nextURL     string // store what URL we are at when using commandMap
	previousURL string // see above
	cache       *pokecache.Cache
}

func NewREPL(interval int) *REPL {

	return &REPL{
		commands:    make(map[string]cliCommand),
		nextURL:     "",
		previousURL: "",
		cache:       pokecache.NewCache(interval),
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

// main function: open a CLI that loops until interrupt or commandExit() is called
func (r *REPL) ReplCLI() {
	// register commands here
	r.RegisterCommand("exit", "Exit the Pokedex", r.commandExit, "none")
	r.RegisterCommand("help", "Prints this help message", r.commandHelp, "none")
	r.RegisterCommand("map", "View map locations", r.commandMap, "none")
	r.RegisterCommand("mapb", "View previous map locations", r.commandMapb, "none")
	r.RegisterCommand("explore", "View list of wild Pokemon on a given map location", r.commandExplore, "<area-name>")

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
