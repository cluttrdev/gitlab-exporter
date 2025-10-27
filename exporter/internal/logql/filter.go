package logql

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strings"
)

type LineFilter []filterExpression

func (f *LineFilter) Match(line []byte) bool {
	for _, expr := range *f {
		if !expr.Match(line) {
			return false
		}
	}
	return true
}

type filterExpression struct {
	Operator string
	Patterns []string
}

func (e filterExpression) Match(line []byte) bool {
	switch e.Operator {
	case "|=":
		for _, p := range e.Patterns {
			if strings.Contains(string(line), p) {
				return true
			}
		}
	case "!=":
		for _, p := range e.Patterns {
			if strings.Contains(string(line), p) {
				return false
			}
		}
	case "|~":
		for _, p := range e.Patterns {
			matched, _ := regexp.Match(p, line)
			if matched {
				return true
			}
		}
		return false
	case "!~":
		for _, p := range e.Patterns {
			matched, _ := regexp.Match(p, line)
			if matched {
				return false
			}
		}
		return true
	}

	return false
}

func ParseLineFilter(s string) (LineFilter, error) {
	var p parser
	return p.ParseString(s)
}

// ===========================================================================

type parser struct {
	l *lexer

	currToken token
	peekToken token
}

func (p *parser) Parse(r io.Reader) ([]filterExpression, error) {
	// create new lexer and ensure currToken and peekToken are set
	p.l = newLexer(r)
	p.nextToken()

	exprs, err := p.parseFilterExpressions()
	if err != nil {
		return nil, err
	}

	return exprs, nil
}

func (p *parser) ParseString(s string) ([]filterExpression, error) {
	return p.Parse(strings.NewReader(s))
}

func (p *parser) nextToken() {
	p.currToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *parser) currTokenIs(t tokenType) bool {
	return p.currToken.Type == t
}

func (p *parser) peekTokenIs(t tokenType) bool {
	return p.peekToken.Type == t
}

func (p *parser) expect(t tokenType) error {
	if p.peekToken.Type != t {
		return fmt.Errorf("expected token type %v, got %v", t, p.peekToken.Type)
	}
	p.nextToken()
	return nil
}

func (p *parser) parseFilterExpressions() ([]filterExpression, error) {
	var exprs []filterExpression
	for p.peekTokenIs(TOKEN_FILTER_OPERATOR) {
		expr, err := p.parseFilterExpression()
		if err != nil {
			return nil, err
		}

		exprs = append(exprs, expr)
	}

	return exprs, nil
}

func (p *parser) parseFilterExpression() (filterExpression, error) {
	if err := p.expect(TOKEN_FILTER_OPERATOR); err != nil {
		return filterExpression{}, err
	}
	expr := filterExpression{
		Operator: p.currToken.Literal,
	}

	if err := p.expect(TOKEN_FILTER_STRING); err != nil {
		return filterExpression{}, err
	}
	expr.Patterns = append(expr.Patterns, p.currToken.Literal)

	for p.peekTokenIs(TOKEN_FILTER_OR) {
		p.nextToken()

		if err := p.expect(TOKEN_FILTER_STRING); err != nil {
			return filterExpression{}, err
		}
		expr.Patterns = append(expr.Patterns, p.currToken.Literal)
	}

	return expr, nil
}

// ========

type tokenType string

const (
	TOKEN_ILLEGAL tokenType = "ILLEGAL"
	TOKEN_EOL               = "EOL"
	TOKEN_WS                = "WHITESPACE"

	TOKEN_FILTER_OPERATOR = "OPERATOR"
	TOKEN_FILTER_STRING   = "STRING"
	TOKEN_FILTER_OR       = "OR"
)

type token struct {
	Type    tokenType
	Literal string
}

func newToken(t tokenType, l string) token {
	return token{
		Type:    t,
		Literal: l,
	}
}

func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\r' || ch == '\n'
}

func isFilterOperator(lit string) bool {
	return lit == "|=" || lit == "|~" || lit == "!=" || lit == "!~"
}

func isStringLiteral(ch rune) bool {
	return ch == '"' || ch == '\'' || ch == '`'
}

// ========

type lexer struct {
	r *bufio.Reader
}

func newLexer(r io.Reader) *lexer {
	return &lexer{
		r: bufio.NewReader(r),
	}
}

func (l *lexer) NextToken() token {
	ch := l.readIgnoreWhitespace()

	switch ch {
	case eol:
		return newToken(TOKEN_EOL, "")
	case '|', '!':
		if lit := string(ch) + string(l.peek()); isFilterOperator(lit) {
			l.read()
			return newToken(TOKEN_FILTER_OPERATOR, lit)
		}
	case 'o':
		if lit := string(ch) + string(l.peek()); lit == "or" {
			l.read()
			return newToken(TOKEN_FILTER_OR, lit)
		}
	default:
		if isStringLiteral(ch) {
			return newToken(TOKEN_FILTER_STRING, l.readString(ch))
		}
	}

	return token{
		Type: TOKEN_ILLEGAL,
	}
}

const eol = rune(0)

func (l *lexer) read() rune {
	ch, _, err := l.r.ReadRune()
	if err != nil {
		return eol
	}
	return ch
}

func (l *lexer) unread() {
	_ = l.r.UnreadRune()
}

func (l *lexer) readIgnoreWhitespace() rune {
	for {
		ch := l.read()
		if !isWhitespace(ch) {
			return ch
		}
	}
}

func (l *lexer) peek() rune {
	bs, err := l.r.Peek(1)
	if err != nil {
		return eol
	}
	return rune(bs[0])
}

func (l *lexer) readString(lit rune) string {
	var buf bytes.Buffer

	for buf.Len() < 255 {
		ch := l.read()
		if ch == eol || ch == lit {
			break
		}

		if ch == '\\' {
			bs, err := l.r.Peek(1)
			if err != nil {
				break
			}
			if bs[0] == byte(lit) || bs[0] == '\\' {
				ch = l.read()
			} else {
				break
			}
		}

		_, _ = buf.WriteRune(ch)
	}

	return buf.String()
}
