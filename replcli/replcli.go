package replcli

import (
	"bufio"
	"fmt"
	"os"
)

func replCLI() {
	s := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")
		s.Scan()
		input := s.Text()
		fI := stringutils.cleanInput(input)
		fmt.Printf("Your command was: %s", fI[0])
	}
}
