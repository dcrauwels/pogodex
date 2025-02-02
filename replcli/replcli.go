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

// define struct for REPL and associated function
type REPL struct {
	commands map[string]cliCommand
}

func NewREPL() *REPL {
	return &REPL{
		commands: make(map[string]cliCommand),
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

// first command: exit pokedex
func (r *REPL) commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")

	// any cleanup functions (that might produce errors) go here

	os.Exit(0)
	return nil
}

// second command: print help
func (r *REPL) commandHelp() error {
	// header fluff
	fmt.Println("Welcome to the Pokedex!\nUsage:")

	// sanity check
	if len(r.commands) < 1 {
		return fmt.Errorf("no commands are implemented")
	}

	// loop over all commands
	for _, command := range r.commands {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}
	return nil
}

// main function: open a CLI that loops until interrupt or commandExit() is called
func (r *REPL) ReplCLI() error {
	// register commands here
	r.RegisterCommand("exit", "Exit the Pokedex", r.commandExit)
	r.RegisterCommand("help", "Prints this help message", r.commandHelp)

	s := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")

		// take input and clean it
		s.Scan()
		input := s.Text()
		cleanedInput := stringutils.CleanInput(input)
		firstInput := cleanedInput[0]

		//check if command in r.commands
		command, ok := r.commands[firstInput]
		if !ok {
			fmt.Errorf("Unknown command: %s", firstInput)
		}

		// try command, raise error if an issue arises
		if err := command.callback(); err != nil {
			fmt.Errorf("Error executing command: %w", err)
		}

		command.callback()
		fmt.Printf("Your command was: %s\n", firstInput)
	}
}
