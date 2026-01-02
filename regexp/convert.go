package regexp

import (
	"fmt"
	"regexp"

	"github.com/cartabinaria/synta"
)

func Convert(synta synta.Synta) (expr *regexp.Regexp, err error) {
	finalString, err := convertWithoutExtensionString(synta.Definitions, synta.Filename.Segments)
	if err != nil {
		return
	}

	finalString += "\\.(" + synta.Definitions[synta.Filename.Extension].Regexp.String() + ")"
	expr, err = regexp.Compile("^" + finalString + "$")

	// Simplify when we use regexp/syntax
	// if err == nil {
	//     expr = expr.Simplify()
	// }
	return
}

func convertWithoutExtensionString(definitions map[synta.Identifier]synta.Definition, segments []synta.Segment) (expr string, err error) {
	for i, segment := range segments {
		definition := synta.Definition{}

		switch segment.Kind {
		case synta.SegmentTypeIdentifier:
			def, isPresent := definitions[*segment.Value]
			if !isPresent {
				err = fmt.Errorf("missing definition for %s", *segment.Value)
				return
			}

			definition = def
			expr += "(" + definition.Regexp.String() + ")"
		case synta.SegmentTypeOptional:
			exp, e := convertWithoutExtensionString(definitions, segment.Subsegments)
			if e != nil {
				err = e
				return
			}
			expr += "(-" + exp + ")?"
		}

		if i != len(segments)-1 && segments[i+1].Kind != synta.SegmentTypeOptional {
			expr += "-"
		}
	}
	return
}

func ConvertWithoutExtension(synta synta.Synta) (expr *regexp.Regexp, err error) {
	exp, err := convertWithoutExtensionString(synta.Definitions, synta.Filename.Segments)
	if err != nil {
		return
	}
	expr, err = regexp.Compile("^" + exp + "$")
	return
}
