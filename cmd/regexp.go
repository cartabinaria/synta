package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/cartabinaria/synta/regexp"
	"github.com/google/subcommands"
)

type regexpCommand struct {
}

func (*regexpCommand) Name() string     { return "regexp" }
func (*regexpCommand) Synopsis() string { return "Convert synta file into a regular expression" }
func (*regexpCommand) Usage() string {
	return `regexp <file>:
  Convert synta file into a regular expression.
`
}

func (p *regexpCommand) SetFlags(f *flag.FlagSet) {
}

func (p *regexpCommand) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	syntaFilePtr, status := parseFile(p, f)
	if status != subcommands.ExitSuccess {
		return status
	}
	r, err := regexp.Convert(*syntaFilePtr)
	if err != nil {
		fmt.Printf("Could not convert to regexp: %v\n", err)
		return subcommands.ExitFailure
	}

	fmt.Printf("%s\n", r.String())

	return subcommands.ExitSuccess
}
