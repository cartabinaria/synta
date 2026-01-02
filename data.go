package synta

import (
	"regexp"
)

// A regexp which describes an identifier
var IdentifierRegexp = regexp.MustCompile("[a-z]+")

// An Identifier is a lowercase alphabetical string.
// It corresponds to the <id> BNF definition
type Identifier string

// NodeType represents the type of an AST node
type NodeType uint

const (
	NodeTypeDefinition NodeType = iota
	NodeTypeFilename
)

// String returns the string representation of a NodeType
func (nt NodeType) String() string {
	switch nt {
	case NodeTypeDefinition:
		return "Definition"
	case NodeTypeFilename:
		return "Filename"
	default:
		return "Unknown"
	}
}

// Node represents a single statement in the Synta file
type Node struct {
	Type       NodeType
	Identifier Identifier  // used for definitions
	Definition *Definition // used for definitions
	Filename   *Filename   // used for filename node
}

// A Definition is a named regexp along with comments
// to clarify the regexp's purpose.
// It corresponds to the <commdef> BNF definition.
type Definition struct {
	Comments []string
	Regexp   *regexp.Regexp
}

type SegmentType uint

const (
	SegmentTypeIdentifier = iota
	SegmentTypeOptional
)

// String returns the string representation of a SegmentType
func (st SegmentType) String() string {
	switch st {
	case SegmentTypeIdentifier:
		return "Identifier"
	case SegmentTypeOptional:
		return "Optional"
	default:
		return "Unknown"
	}
}

// A Segment is a section of the main filename.
// It corresponds to the <segment> BNF definition.
type Segment struct {
	Kind        SegmentType
	Value       *Identifier
	Subsegments []Segment
}

// Filename represents the filename definition, made up
// of a series of segments and a file extension.
type Filename struct {
	Segments  []Segment
	Extension Identifier
}

// Synta represents the contents of a Synta file.
// It corresponds to the <language> BNF definition.
// The last segment of the Filename is the extension.
type Synta struct {
	Nodes       []Node                    // AST nodes in order
	Definitions map[Identifier]Definition // for quick lookup
	Filename    Filename                  // for quick access
}
