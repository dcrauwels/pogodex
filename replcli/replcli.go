package replcli

import (
	"bufio"
	"fmt"
	"os"

	"github.com/dcrauwels/pogodex/internal/pokeapi"
	"github.com/dcrauwels/pogodex/internal/stringutils"
)

// define struct for commands
type cliCommand struct {
	name        string
	description string
	callback    func() error
}

// define struct for REPL and associated function
type REPL struct {
	commands    map[string]cliCommand
	nextURL     string // store what URL we are at when using commandMap
	previousURL string // see above
}

func NewREPL() *REPL {
	return &REPL{
		commands:    make(map[string]cliCommand),
		nextURL:     "",
		previousURL: "",
	}
}

// add a command to the REPL, for use in the ReplCLI function to not have this huge definition at the start
func (r *REPL) RegisterCommand(name string, description string, callback func() error) {
	r.commands[name] = cliCommand{
		name:        name,
		description: description,
		callback:    callback,
	}
}

// first command: 'exit' pokedex
func (r *REPL) commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")

	// any cleanup functions (that might produce errors) go here

	os.Exit(0)
	return nil
}

// second command: print 'help'
func (r *REPL) commandHelp() error {
	// header fluff
	fmt.Println("Welcome to the Pokedex!\nUsage:")

	// sanity check
	if len(r.commands) < 1 {
		return fmt.Errorf("no commands are implemented")
	}

	// loop over all commands
	for _, command := range r.commands {
		fmt.Printf(" %s: %s\n", command.name, command.description)
	}
	return nil
}

// third command: print area names on the 'map'
func (r *REPL) commandMap() error {

	// construct URL
	var u string
	if r.nextURL != "" {
		u = r.nextURL // this means we have a URL stored from previous commandMap() calls
	} else {
		u = "https://pokeapi.co/api/v2/location-area/"
	}

	locations, err := pokeapi.GetLocations(u)
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

func (r *REPL) commandMapb() error {
	// construct URL and sanity check
	if r.previousURL == "" || r.nextURL == "https://pokeapi.co/api/v2/location-area/?offset=20&limit=20" {
		fmt.Println("You're on the first page of location results")
		return nil
	}
	locations, err := pokeapi.GetLocations(r.previousURL)
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

// main function: open a CLI that loops until interrupt or commandExit() is called
func (r *REPL) ReplCLI() error {
	// register commands here
	r.RegisterCommand("exit", "Exit the Pokedex", r.commandExit)
	r.RegisterCommand("help", "Prints this help message", r.commandHelp)
	r.RegisterCommand("map", "View map locations", r.commandMap)
	r.RegisterCommand("mapb", "View previous map locations", r.commandMapb)

	s := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")

		// take input and clean it
		s.Scan()
		input := s.Text()
		if ok := len(input) > 0; !ok {
			//fmt.Println("Please enter a command")
			continue
		}
		cleanedInput := stringutils.CleanInput(input)
		firstInput := cleanedInput[0]

		//check if command in r.commands
		command, ok := r.commands[firstInput]
		if !ok {
			fmt.Printf("Unknown command: %s\n", firstInput)
			continue
		}

		// try command, raise error if an issue arises
		if err := command.callback(); err != nil {
			fmt.Printf("Error executing command: %v\n", err)
			continue
		}

		fmt.Printf("Your command was: %s\n", firstInput)
	}
}
