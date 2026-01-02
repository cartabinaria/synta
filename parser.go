package synta

import (
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"
)

// ParseSynta attempts to parse a file's contents into a Synta internal
// representation. If an error is encountered, the parsing is aborted and the
// error returned.
func ParseSynta(contents string) (Synta, error) {
	return ParseSyntaFromReader(strings.NewReader(contents))
}

// ParseSyntaFromReader parses Synta from an io.Reader, reading lines on demand.
func ParseSyntaFromReader(r io.Reader) (Synta, error) {
	lexer := NewLexer(r)
	p := &parser{
		lexer: lexer,
	}

	// Read first token
	if err := p.advance(); err != nil {
		return Synta{}, err
	}

	return p.parseFile()
}

func getRequiredIdentifiers(segments []Segment) (requiredIdentifiers []Identifier) {
	for _, seg := range segments {
		if seg.Kind == SegmentTypeOptional {
			requiredIdentifiers = append(requiredIdentifiers, getRequiredIdentifiers(seg.Subsegments)...)
		} else {
			requiredIdentifiers = append(requiredIdentifiers, *seg.Value)
		}
	}
	return
}

// MustSynta parses the contents and panics if an error occurs.
// This is useful for testing or when the input is known to be valid.
func MustSynta(contents string) Synta {
	s, err := ParseSynta(contents)
	if err != nil {
		panic(err)
	}

	return s
}

// parser represents the unified parser for the entire file
type parser struct {
	lexer        *Lexer
	currentToken Token
}

// advance moves to the next token
func (p *parser) advance() error {
	token, err := p.lexer.NextToken()
	if err != nil {
		return err
	}
	p.currentToken = token
	return nil
}

// parseFile parses the entire file from start to end
func (p *parser) parseFile() (Synta, error) {
	var s Synta
	s.Definitions = map[Identifier]Definition{}

	if p.currentToken.Type == TokenEOF {
		return Synta{}, errors.New("empty file provided")
	}

	// Parse tokens into AST nodes
	for p.currentToken.Type != TokenEOF {
		var node Node
		var err error

		switch p.currentToken.Type {
		case TokenComment:
			// Comment tokens are collected as part of definitions
			node, err = p.parseDefinitionNode()
		case TokenIdentifier:
			// Definition starts with identifier
			node, err = p.parseDefinitionNode()
		case TokenFilenamePrefix:
			// Filename declaration
			node, err = p.parseFilenameNode()
		default:
			return Synta{}, fmt.Errorf("unexpected token at line %d: %s", p.currentToken.Line, p.currentToken.Type)
		}

		if err != nil {
			return Synta{}, err
		}

		s.Nodes = append(s.Nodes, node)

		// Populate quick-access structures
		switch node.Type {
		case NodeTypeDefinition:
			if _, ok := s.Definitions[node.Identifier]; ok {
				return Synta{}, fmt.Errorf("definition for `%s` is provided twice", node.Identifier)
			}
			s.Definitions[node.Identifier] = *node.Definition
		case NodeTypeFilename:
			if s.Filename.Extension != "" {
				return Synta{}, errors.New("multiple filename declarations found")
			}
			s.Filename = *node.Filename
		}
	}

	// Validate that we have a filename
	if s.Filename.Extension == "" {
		return Synta{}, errors.New("missing filename declaration")
	}

	// Validate that all required identifiers are defined
	requiredIdentifiers := getRequiredIdentifiers(s.Filename.Segments)
	requiredIdentifiers = append(requiredIdentifiers, s.Filename.Extension)
	for _, id := range requiredIdentifiers {
		if _, ok := s.Definitions[id]; !ok {
			return Synta{}, fmt.Errorf("missing definition for `%s`", id)
		}
	}

	return s, nil
}

// parseFilenameNode parses a filename declaration from tokens
func (p *parser) parseFilenameNode() (Node, error) {
	if p.currentToken.Type != TokenFilenamePrefix {
		return Node{}, fmt.Errorf("expected filename prefix, got %s", p.currentToken.Type)
	}

	// Advance past '>'
	if err := p.advance(); err != nil {
		return Node{}, err
	}

	// Parse segments
	segments, err := p.parseSegments()
	if err != nil {
		return Node{}, err
	}

	// Expect dot
	if p.currentToken.Type != TokenDot {
		return Node{}, fmt.Errorf("expected '.' before extension, got %s", p.currentToken.Type)
	}
	if err := p.advance(); err != nil {
		return Node{}, err
	}

	// Parse extension
	ext, err := p.parseIdentifier()
	if err != nil {
		return Node{}, err
	}

	return Node{
		Type: NodeTypeFilename,
		Filename: &Filename{
			Segments:  segments,
			Extension: ext,
		},
	}, nil
}

// parseDefinitionNode parses a definition from tokens
func (p *parser) parseDefinitionNode() (Node, error) {
	var def Definition

	// Collect comment tokens
	for p.currentToken.Type == TokenComment {
		def.Comments = append(def.Comments, p.currentToken.Value)
		if err := p.advance(); err != nil {
			return Node{}, err
		}
	}

	// Expect identifier
	if p.currentToken.Type != TokenIdentifier {
		return Node{}, fmt.Errorf("expected identifier at line %d, got %s", p.currentToken.Line, p.currentToken.Type)
	}
	id := Identifier(p.currentToken.Value)

	// Advance past identifier
	if err := p.advance(); err != nil {
		return Node{}, err
	}

	// Expect equals
	if p.currentToken.Type != TokenEquals {
		return Node{}, fmt.Errorf("expected '=' at line %d, got %s", p.currentToken.Line, p.currentToken.Type)
	}
	if err := p.advance(); err != nil {
		return Node{}, err
	}

	// Expect regexp pattern
	if p.currentToken.Type != TokenRegexpPattern {
		return Node{}, fmt.Errorf("expected regexp pattern at line %d, got %s", p.currentToken.Line, p.currentToken.Type)
	}

	var err error
	def.Regexp, err = regexp.Compile(p.currentToken.Value)
	if err != nil {
		return Node{}, fmt.Errorf("invalid regexp at line %d: %w", p.currentToken.Line, err)
	}

	// Advance past regexp
	if err := p.advance(); err != nil {
		return Node{}, err
	}

	return Node{
		Type:       NodeTypeDefinition,
		Identifier: id,
		Definition: &def,
	}, nil
}

// expect checks if the current token matches the expected type and advances
func (p *parser) expect(expected TokenType) error {
	if p.currentToken.Type != expected {
		return fmt.Errorf("expected %s, got %s at line %d", expected, p.currentToken.Type, p.currentToken.Line)
	}
	return p.advance()
}

// parseIdentifier parses an identifier token
func (p *parser) parseIdentifier() (Identifier, error) {
	if p.currentToken.Type != TokenIdentifier {
		return "", fmt.Errorf("expected identifier, got %s", p.currentToken.Type)
	}
	id := Identifier(p.currentToken.Value)
	if err := p.advance(); err != nil {
		return "", err
	}
	return id, nil
}

// parseSegment parses a single segment (identifier or optional)
func (p *parser) parseSegment() (Segment, error) {
	if p.currentToken.Type == TokenLParen {
		return p.parseOptional()
	}
	id, err := p.parseIdentifier()
	if err != nil {
		return Segment{}, err
	}
	return Segment{
		Kind:  SegmentTypeIdentifier,
		Value: &id,
	}, nil
}

// parseOptional parses an optional segment: (-<segments>)?
func (p *parser) parseOptional() (Segment, error) {
	if err := p.expect(TokenLParen); err != nil {
		return Segment{}, err
	}
	if err := p.expect(TokenDash); err != nil {
		return Segment{}, err
	}

	// Parse segments within the optional group using parseSegments
	// This will handle dash separators and nested optionals
	subsegments, err := p.parseSegments()
	if err != nil {
		return Segment{}, err
	}

	if err := p.expect(TokenRParen); err != nil {
		return Segment{}, err
	}
	if err := p.expect(TokenQuestion); err != nil {
		return Segment{}, err
	}

	return Segment{
		Kind:        SegmentTypeOptional,
		Value:       nil,
		Subsegments: subsegments,
	}, nil
}

// parseSegments parses a sequence of segments separated by '-'
// Optionals can appear without a dash separator
func (p *parser) parseSegments() ([]Segment, error) {
	var segments []Segment

	for {
		seg, err := p.parseSegment()
		if err != nil {
			return nil, err
		}
		segments = append(segments, seg)

		if p.currentToken.Type == TokenDot {
			// reached extension separator
			break
		} else if p.currentToken.Type == TokenRParen {
			// end of optional group (handled by caller)
			break
		} else if p.currentToken.Type == TokenEOF {
			return nil, fmt.Errorf("unexpected end of file while parsing segments")
		} else if p.currentToken.Type == TokenDash {
			if err := p.advance(); err != nil {
				return nil, err
			}
			// continue parsing more segments after dash
		} else if p.currentToken.Type == TokenLParen {
			// optional can follow directly without dash
			continue
		} else {
			return nil, fmt.Errorf("expected '.', '-', ')', or '(', got %s", p.currentToken.Type)
		}
	}

	return segments, nil
}
