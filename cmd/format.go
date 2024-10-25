package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/cartabinaria/synta"
	"github.com/cartabinaria/synta/format"
	"github.com/google/subcommands"
)

type formatCommand struct {
	clear bool
	write bool
}

func (*formatCommand) Name() string     { return "format" }
func (*formatCommand) Synopsis() string { return "Format a synta file in the standard style." }
func (*formatCommand) Usage() string {
	return `format [-clear] [-write] <file>:
  Format a synta file in the standard style.
`
}

func (p *formatCommand) SetFlags(f *flag.FlagSet) {
	f.BoolVar(&p.clear, "clear", false, "Remove unused definitions")
	f.BoolVar(&p.write, "write", false, "Write changes to the input file")
}

func (p *formatCommand) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	syntaFilePtr, status := parseFile(p, f)
	if status != subcommands.ExitSuccess {
		return status
	}

	syntaFile := *syntaFilePtr
	if p.clear {
		syntaFile = synta.Clear(syntaFile)
	}

	formatted := format.Format(syntaFile)
	if p.write {
		ioutil.WriteFile(f.Arg(0), []byte(formatted), 0664)
	} else {
		fmt.Printf("%s", formatted)
	}
	return subcommands.ExitSuccess
}
