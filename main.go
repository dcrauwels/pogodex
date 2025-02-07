package main

import (
	"github.com/dcrauwels/pogodex/replcli"
)

func main() {
	r := replcli.NewREPL(5)
	r.ReplCLI()
}
