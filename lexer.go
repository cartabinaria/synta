package synta

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// TokenType represents the type of a token
type TokenType int

const (
	TokenEOF            TokenType = iota
	TokenComment                  // ; comment
	TokenEquals                   // =
	TokenFilenamePrefix           // >
	TokenIdentifier               // sequence of lowercase letters
	TokenRegexpPattern            // regexp pattern after =
	TokenDash                     // -
	TokenDot                      // .
	TokenLParen                   // (
	TokenRParen                   // )
	TokenQuestion                 // ?
)

// String returns the string representation of a TokenType
func (tt TokenType) String() string {
	switch tt {
	case TokenEOF:
		return "EOF"
	case TokenComment:
		return "COMMENT"
	case TokenEquals:
		return "EQUALS"
	case TokenFilenamePrefix:
		return "FILENAME_PREFIX"
	case TokenIdentifier:
		return "IDENTIFIER"
	case TokenRegexpPattern:
		return "REGEXP"
	case TokenDash:
		return "DASH"
	case TokenDot:
		return "DOT"
	case TokenLParen:
		return "LPAREN"
	case TokenRParen:
		return "RPAREN"
	case TokenQuestion:
		return "QUESTION"
	default:
		return "UNKNOWN"
	}
}

// Token represents a lexical token
type Token struct {
	Type  TokenType
	Value string
	Pos   int // position in current context
	Line  int // line number in file
}

// Lexer tokenizes input from a reader
type Lexer struct {
	scanner       *bufio.Scanner
	currentLine   string
	lineNum       int
	pos           int  // position within current line
	inFilename    bool // true when parsing filename segments
	pendingTokens []Token
}

// NewLexer creates a new lexer for the given input
func NewLexer(r io.Reader) *Lexer {
	return &Lexer{
		scanner: bufio.NewScanner(r),
		lineNum: 0,
	}
}

// NextToken returns the next token from the input
func (l *Lexer) NextToken() (Token, error) {
	// Return pending tokens first
	if len(l.pendingTokens) > 0 {
		token := l.pendingTokens[0]
		l.pendingTokens = l.pendingTokens[1:]
		return token, nil
	}

	// Read next non-empty line if needed
	if l.pos >= len(l.currentLine) {
		if !l.readNextLine() {
			return Token{Type: TokenEOF, Line: l.lineNum}, nil
		}
	}

	line := l.currentLine
	l.pos = len(line) // consume entire line

	// Check line type and produce appropriate tokens
	if len(line) > 0 && line[0] == ';' {
		// Comment line
		return Token{
			Type:  TokenComment,
			Value: strings.TrimSpace(line[1:]),
			Line:  l.lineNum,
		}, nil
	}

	if len(line) >= 2 && line[:2] == "> " {
		// Filename declaration - tokenize segments
		return l.tokenizeFilename(line[2:])
	}

	// Definition line: identifier = regexp
	return l.tokenizeDefinition(line)
}

// readNextLine reads the next non-blank line
func (l *Lexer) readNextLine() bool {
	for l.scanner.Scan() {
		l.lineNum++
		line := strings.TrimSpace(l.scanner.Text())
		if line != "" {
			l.currentLine = line
			l.pos = 0
			return true
		}
	}
	return false
}

// tokenizeDefinition tokenizes a definition line
func (l *Lexer) tokenizeDefinition(line string) (Token, error) {
	parts := strings.SplitN(line, " = ", 2)
	if len(parts) != 2 {
		return Token{}, fmt.Errorf("invalid definition format at line %d: %s", l.lineNum, line)
	}

	id := parts[0]
	pattern := parts[1]

	if !IdentifierRegexp.Match([]byte(id)) {
		return Token{}, fmt.Errorf("invalid identifier at line %d: %s", l.lineNum, id)
	}

	// Queue tokens: Identifier, Equals, RegexpPattern
	l.pendingTokens = []Token{
		{Type: TokenIdentifier, Value: id, Line: l.lineNum},
		{Type: TokenEquals, Value: "=", Line: l.lineNum},
		{Type: TokenRegexpPattern, Value: pattern, Line: l.lineNum},
	}

	// Return first token
	token := l.pendingTokens[0]
	l.pendingTokens = l.pendingTokens[1:]
	return token, nil
}

// tokenizeFilename tokenizes a filename declaration
func (l *Lexer) tokenizeFilename(segments string) (Token, error) {
	// First token is the filename prefix
	l.inFilename = true

	// Parse segments and create tokens
	tokens, err := l.parseFilenameSegments(segments)
	if err != nil {
		return Token{}, err
	}

	l.pendingTokens = tokens
	l.inFilename = false

	// Prepend filename prefix token
	prefixToken := Token{Type: TokenFilenamePrefix, Value: ">", Line: l.lineNum}
	if len(l.pendingTokens) > 0 {
		token := l.pendingTokens[0]
		l.pendingTokens = l.pendingTokens[1:]
		l.pendingTokens = append([]Token{token}, l.pendingTokens...)
	}
	return prefixToken, nil
}

// parseFilenameSegments parses filename segments into tokens
func (l *Lexer) parseFilenameSegments(input string) ([]Token, error) {
	var tokens []Token
	pos := 0

	for pos < len(input) {
		c := input[pos]

		switch c {
		case '-':
			tokens = append(tokens, Token{Type: TokenDash, Value: "-", Pos: pos, Line: l.lineNum})
			pos++
		case '.':
			tokens = append(tokens, Token{Type: TokenDot, Value: ".", Pos: pos, Line: l.lineNum})
			pos++
		case '(':
			tokens = append(tokens, Token{Type: TokenLParen, Value: "(", Pos: pos, Line: l.lineNum})
			pos++
		case ')':
			tokens = append(tokens, Token{Type: TokenRParen, Value: ")", Pos: pos, Line: l.lineNum})
			pos++
		case '?':
			tokens = append(tokens, Token{Type: TokenQuestion, Value: "?", Pos: pos, Line: l.lineNum})
			pos++
		default:
			if isLetter(c) {
				start := pos
				for pos < len(input) && isLetter(input[pos]) {
					pos++
				}
				tokens = append(tokens, Token{
					Type:  TokenIdentifier,
					Value: input[start:pos],
					Pos:   start,
					Line:  l.lineNum,
				})
			} else {
				return nil, fmt.Errorf("unexpected character '%c' at position %d in filename", c, pos)
			}
		}
	}

	return tokens, nil
}

// isLetter checks if a character is a lowercase letter
func isLetter(c byte) bool {
	return c >= 'a' && c <= 'z'
}
