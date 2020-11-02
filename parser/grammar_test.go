package parser

import (
	"testing"
)

func TestGetFirstBody(t *testing.T) {
	var testData = []struct {
		p    production
		want grammarSymbol
	}{
		{production{"expr", []grammarSymbol{"expr", "+", "term"}}, "expr"},
		{production{"term", []grammarSymbol{"factor"}}, "factor"},
		{production{"expr'", []grammarSymbol{}}, ""},
	}

	for _, test := range testData {
		if got := test.p.getFirstBodySymbol(); got != test.want {
			t.Errorf("p.getFirstBodySymbol() = %v, expected %v", got, test.want)
		}
	}
}

func TestComputeFirstSet(t *testing.T) {
	var g grammar
	g.start = "expr'"
	g.productions = []production{
		{"expr'", []grammarSymbol{"expr"}},
		{"expr", []grammarSymbol{"expr", "+", "term"}},
		{"expr", []grammarSymbol{"term"}},
		{"term", []grammarSymbol{"term", "*", "factor"}},
		{"term", []grammarSymbol{"factor"}},
		{"factor", []grammarSymbol{"number"}},
		{"factor", []grammarSymbol{"(", "expr", ")"}},
	}

	var testData = []struct {
		input    grammarSymbol
		expected setOfSymbols
	}{
		{"number", setOfSymbols{"number": true}},
		{"expr", setOfSymbols{"(": true, "number": true}},
		{"term", setOfSymbols{"(": true, "number": true}},
		{"factor", setOfSymbols{"(": true, "number": true}},
	}

	for _, test := range testData {
		if got := g.computeFirstSet(test.input); !got.isEqualTo(&test.expected) {
			t.Errorf("Expected computeFirstSet(%q) = %v, got %v", test.input, test.expected, got)
		}
	}
}

func TestComputeFollowSet(t *testing.T) {
	var g grammar
	g.start = "expr'"
	g.productions = []production{
		{"expr'", []grammarSymbol{"expr"}},
		{"expr", []grammarSymbol{"expr", "+", "term"}},
		{"expr", []grammarSymbol{"term"}},
		{"term", []grammarSymbol{"term", "*", "factor"}},
		{"term", []grammarSymbol{"factor"}},
		{"factor", []grammarSymbol{"number"}},
		{"factor", []grammarSymbol{"(", "expr", ")"}},
	}

	var testData = []struct {
		input    grammarSymbol
		expected setOfSymbols
	}{
		{"expr", setOfSymbols{"+": true, ")": true, "$": true}},
		{"term", setOfSymbols{"+": true, "*": true, ")": true, "$": true}},
		{"factor", setOfSymbols{"+": true, "*": true, ")": true, "$": true}},
	}

	for _, test := range testData {
		if got := g.computeFollowSet(test.input); !got.isEqualTo(&test.expected) {
			t.Errorf("Expected computeFirstSet(%q) = %v, got %v", test.input, test.expected, got)
		}
	}
}

func TestGetProductionNumber(t *testing.T) {
	var g grammar
	g.start = "expr'"
	g.productions = []production{
		{"expr'", []grammarSymbol{"expr"}},
		{"expr", []grammarSymbol{"expr", "+", "term"}},
		{"expr", []grammarSymbol{"term"}},
		{"term", []grammarSymbol{"term", "*", "factor"}},
		{"term", []grammarSymbol{"factor"}},
		{"factor", []grammarSymbol{"number"}},
		{"factor", []grammarSymbol{"(", "expr", ")"}},
	}

	var testData = []struct {
		input    production
		expected int
	}{
		{production{"expr", []grammarSymbol{"expr", "+", "term"}}, 1},
		{production{"term", []grammarSymbol{"term", "*", "factor"}}, 3},
		{production{"factor", []grammarSymbol{"number"}}, 5},
		{production{"expr", []grammarSymbol{"factor"}}, -1},
	}

	for _, test := range testData {
		if got := g.getProductionNumber(test.input); got != test.expected {
			t.Errorf("Expected getProductionNumber(%v) to be %v but got %v", test.input, test.expected, got)
		}
	}
}

func TestCompile(t *testing.T) {
	var g grammar
	g.start = "expr'"
	g.productions = []production{
		{"expr'", []grammarSymbol{"expr"}},
		{"expr", []grammarSymbol{"expr", "+", "term"}},
		{"expr", []grammarSymbol{"term"}},
		{"term", []grammarSymbol{"term", "*", "factor"}},
		{"term", []grammarSymbol{"factor"}},
		{"factor", []grammarSymbol{"number"}},
		{"factor", []grammarSymbol{"(", "expr", ")"}},
	}

	var testData = []struct {
		input    []string
		expected bool
	}{
		{[]string{"number", "+", "number", "$"}, true},
		{[]string{"number", "+", "number", "*", "number", "$"}, true},
		{[]string{"(", "number", "+", "number", "*", "number", "$"}, false},
		{[]string{"(", "number", "+", "number", "*", "number", ")", "$"}, true},
		{[]string{"(", "number", "+", "number", ")", "*", "number", ")", "$"}, false},
		{[]string{"(", "number", "+", "number", ")", "*", "(", "number", "*", "number", ")", "$"}, true},
	}

	for _, test := range testData {
		ps := g.compile()
		for _, token := range test.input {
			ps.move(grammarSymbol(token))
		}
		if ps.accepted != test.expected {
			t.Errorf("Expected parser to output %v on input %v", test.expected, test.input)
		}
	}
}
