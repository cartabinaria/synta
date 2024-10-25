package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/cartabinaria/synta/json"
	"github.com/google/subcommands"
)

type jsonCommand struct{}

func (*jsonCommand) Name() string     { return "json" }
func (*jsonCommand) Synopsis() string { return "Return a json for the .synta file provided." }
func (*jsonCommand) Usage() string {
	return `json <file>:
	Return a json for the .synta file provided.
`
}

func (p *jsonCommand) SetFlags(f *flag.FlagSet) {}

func (p *jsonCommand) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	syntaFilePtr, code := parseFile(p, f)
	if code != subcommands.ExitSuccess {
		return code
	}

	res, err := json.ToJson(*syntaFilePtr)
	if err != nil {
		fmt.Printf("Error while converting from .synta to .json: %v\n", err)
		return subcommands.ExitFailure
	}
	fmt.Println(string(res))
	return subcommands.ExitSuccess
}
