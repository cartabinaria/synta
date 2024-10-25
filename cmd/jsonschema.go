package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"

	"github.com/cartabinaria/synta"
	"github.com/google/subcommands"
	"github.com/invopop/jsonschema"
)

type jsonSchemaCommand struct{}

func (*jsonSchemaCommand) Name() string     { return "jsonschema" }
func (*jsonSchemaCommand) Synopsis() string { return "Return a json schema for Synta." }
func (*jsonSchemaCommand) Usage() string {
	return `jsonschema
  Return a json schema for Synta.
`
}

func (p *jsonSchemaCommand) SetFlags(f *flag.FlagSet) {}

func (p *jsonSchemaCommand) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	res, err := json.Marshal(jsonschema.Reflect(&synta.Synta{}))
	if err != nil {
		fmt.Printf("Error from converting Synta to json schema\n")
		return subcommands.ExitFailure
	}
	fmt.Println(string(res))
	return subcommands.ExitSuccess
}
