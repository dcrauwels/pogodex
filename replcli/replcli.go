package replcli

import (
	"bufio"
	"fmt"
	"os"

	"github.com/dcrauwels/pogodex/stringutils"
)

// define struct for commands
type cliCommand struct {
	name        string
	description string
	callback    func() error
}

// first command: exit pokedex
func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")

	// any cleanup functions (that might produce errors) go here

	os.Exit(0)
	return nil
}

// second command: print help
func commandHelp(m map[string]cliCommand) error {
	// header fluff
	fmt.Println("Welcome to the Pokedex!\nUsage:\n")

	// sanity check
	if len(m) < 1 {
		return fmt.Errorf("no commands are implemented")
	}

	// loop over all commands
	for _, command := range m {
		fmt.Fprintln("%s: %s", command.name, command.description)
	}
	return nil
}

// main function: open a CLI that loops until interrupt or commandExit() is called
func ReplCLI() {
	commands := map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Prints a help message",
			callback:    commandHelp,
		},
	}

	s := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")
		s.Scan()
		input := s.Text()
		cleanedInput := stringutils.CleanInput(input)
		firstInput := cleanedInput[0]
		if command, ok := commands[firstInput]; ok {
			command.callback()
		}
		fmt.Printf("Your command was: %s\n", firstInput)
	}
}
