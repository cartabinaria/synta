package format

import (
	"sort"

	"github.com/cartabinaria/synta"
)

// Format formats a synta structure into a string which represnts the contents
// of the associated synta file. The definitions are sorted by name. and blank
// lines are let as per internal convention.
func Format(syntaFile synta.Synta) (code string) {
	definitions := []string{}
	for id := range syntaFile.Definitions {
		definitions = append(definitions, string(id))
	}
	sort.Strings(definitions)

	for _, id := range definitions {
		def := syntaFile.Definitions[synta.Identifier(id)]
		for _, comment := range def.Comments {
			code += "; " + comment + "\n"
		}
		code += string(id) + " = " + def.Regexp.String() + "\n\n"
	}

	code += "> "
	expr := formatSegments(syntaFile.Filename.Segments)
	code += expr
	code += "." + string(syntaFile.Filename.Extension) + "\n"

	return
}

func formatSegments(segments []synta.Segment) (expr string) {
	for i, segment := range segments {
		switch segment.Kind {
		case synta.SegmentTypeIdentifier:
			expr += string(*segment.Value)
		case synta.SegmentTypeOptional:
			exp := formatSegments(segment.Subsegments)
			expr += "(-" + exp + ")?"
		}

		if i != len(segments)-1 && segments[i+1].Kind != synta.SegmentTypeOptional {
			expr += "-"
		}
	}
	return
}
