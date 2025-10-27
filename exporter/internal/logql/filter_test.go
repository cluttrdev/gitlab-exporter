package logql

import (
	"strings"
	"testing"
)

func TestNextToken(t *testing.T) {
	input := "|= \"error\" != 'warning' or 'warn' |~ `^info` !~ `debug$`"

	tests := []struct {
		expectedType    tokenType
		expectedLiteral string
	}{
		{TOKEN_FILTER_OPERATOR, "|="},
		{TOKEN_FILTER_STRING, "error"},
		{TOKEN_FILTER_OPERATOR, "!="},
		{TOKEN_FILTER_STRING, "warning"},
		{TOKEN_FILTER_OR, "or"},
		{TOKEN_FILTER_STRING, "warn"},
		{TOKEN_FILTER_OPERATOR, "|~"},
		{TOKEN_FILTER_STRING, "^info"},
		{TOKEN_FILTER_OPERATOR, "!~"},
		{TOKEN_FILTER_STRING, "debug$"},
	}

	lex := newLexer(strings.NewReader(input))

	for i, tt := range tests {
		tok := lex.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - wrong token type. expected=%q, got %q", i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - wrong literal. expected=%q, got %q", i, tt.expectedLiteral, tok.Literal)
		}
	}

	if tok := lex.NextToken(); tok.Type != TOKEN_EOL {
		t.Fatalf("expected EOL token, got %q: %q", tok.Type, tok.Literal)
	}
}

func TestParse(t *testing.T) {
	input := "|= \"error\" != 'warning' or 'warn' |~ `^info` !~ `debug$`"

	expectedFileterExpressions := []filterExpression{
		{Operator: "|=", Patterns: []string{"error"}},
		{Operator: "!=", Patterns: []string{"warning", "warn"}},
		{Operator: "|~", Patterns: []string{"^info"}},
		{Operator: "!~", Patterns: []string{"debug$"}},
	}

	var parser parser

	exprs, err := parser.ParseString(input)
	if err != nil {
		t.Fatal(err)
	}

	if len(exprs) != len(expectedFileterExpressions) {
		t.Fatalf("expected %d filter expressions, got %d", len(expectedFileterExpressions), len(exprs))
	}

	for i, expectedFileterExpression := range expectedFileterExpressions {
		if exprs[i].Operator != expectedFileterExpression.Operator {
			t.Errorf("expected filter expression operator %q, got %q", expectedFileterExpression.Operator, exprs[i].Operator)
		}

		if len(exprs[i].Patterns) != len(expectedFileterExpressions[i].Patterns) {
			t.Fatalf("expected %d filter expression patterns, got %d", len(expectedFileterExpressions[i].Patterns), len(exprs[i].Patterns))
		}

		for j, expectedPattern := range expectedFileterExpression.Patterns {
			if exprs[i].Patterns[j] != expectedPattern {
				t.Errorf("expected filter expression pattern %q, got %q", expectedPattern, exprs[i].Patterns[j])
			}
		}
	}
}
