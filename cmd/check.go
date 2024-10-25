package main

import (
	"context"
	"flag"
	"fmt"
	"strings"

	"github.com/cartabinaria/synta"
	"github.com/google/subcommands"
)

type checkCommand struct{}

func (*checkCommand) Name() string     { return "check" }
func (*checkCommand) Synopsis() string { return "Checks if a synta file has a corrent syntax." }
func (*checkCommand) Usage() string {
	return `check <file>:
  Checks if a synta file has a corrent syntax.
`
}

func (p *checkCommand) SetFlags(f *flag.FlagSet) {}

func (p *checkCommand) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	syntaFilePtr, status := parseFile(p, f)
	if status != subcommands.ExitSuccess {
		return status
	}

	syntaFile := *syntaFilePtr
	clearedSyntaFile := synta.Clear(syntaFile)

	if len(syntaFile.Definitions) > len(clearedSyntaFile.Definitions) {
		unusedDefintions := []string{}
		for k := range syntaFile.Definitions {
			if _, notRemoved := clearedSyntaFile.Definitions[k]; !notRemoved {
				unusedDefintions = append(unusedDefintions, string(k))
			}
		}

		fmt.Printf("Your Synta file contains unused definitions: %s\n", strings.Join(unusedDefintions, ", "))
		return subcommands.ExitFailure
	}

	return subcommands.ExitSuccess
}
