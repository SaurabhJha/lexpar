package parser

import (
	"reflect"
	"testing"
)

func TestGetFirstBody(t *testing.T) {
	var testData = []struct {
		p    Production
		want grammarSymbol
	}{
		{Production{"expr", []grammarSymbol{"expr", "+", "term"}}, "expr"},
		{Production{"term", []grammarSymbol{"factor"}}, "factor"},
		{Production{"expr'", []grammarSymbol{}}, ""},
	}

	for _, test := range testData {
		if got := test.p.getFirstBodySymbol(); got != test.want {
			t.Errorf("p.getFirstBodySymbol() = %v, expected %v", got, test.want)
		}
	}
}

func TestGetProductionsOfSymbol(t *testing.T) {
	var g Grammar
	g.Start = "expr'"
	g.Productions = []Production{
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
		expected []Production
	}{
		{
			"expr'",
			[]Production{g.Productions[0]},
		},
		{
			"expr",
			[]Production{g.Productions[1], g.Productions[2]},
		},
		{
			"term",
			[]Production{g.Productions[3], g.Productions[4]},
		},
	}

	for _, test := range testData {
		if got := g.getProductionsOfSymbol(test.input); !reflect.DeepEqual(got, test.expected) {
			t.Errorf("Expected productions of %v to be %v but got %v", test.input, test.expected, got)
		}
	}
}

func TestComputeFirstSet(t *testing.T) {
	var g Grammar
	g.Start = "expr'"
	g.Productions = []Production{
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
		if got := g.computeFirstSet(test.input); !reflect.DeepEqual(got, test.expected) {
			t.Errorf("Expected computeFirstSet(%q) = %v, got %v", test.input, test.expected, got)
		}
	}
}

func TestComputeFollowSet(t *testing.T) {
	var g Grammar
	g.Start = "expr'"
	g.Productions = []Production{
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
		if got := g.computeFollowSet(test.input); !reflect.DeepEqual(got, test.expected) {
			t.Errorf("Expected computeFirstSet(%q) = %v, got %v", test.input, test.expected, got)
		}
	}
}

func TestGetProductionNumber(t *testing.T) {
	var g Grammar
	g.Start = "expr'"
	g.Productions = []Production{
		{"expr'", []grammarSymbol{"expr"}},
		{"expr", []grammarSymbol{"expr", "+", "term"}},
		{"expr", []grammarSymbol{"term"}},
		{"term", []grammarSymbol{"term", "*", "factor"}},
		{"term", []grammarSymbol{"factor"}},
		{"factor", []grammarSymbol{"number"}},
		{"factor", []grammarSymbol{"(", "expr", ")"}},
	}

	var testData = []struct {
		input    Production
		expected int
	}{
		{Production{"expr", []grammarSymbol{"expr", "+", "term"}}, 1},
		{Production{"term", []grammarSymbol{"term", "*", "factor"}}, 3},
		{Production{"factor", []grammarSymbol{"number"}}, 5},
		{Production{"expr", []grammarSymbol{"factor"}}, -1},
	}

	for _, test := range testData {
		if got := g.getProductionNumber(test.input); got != test.expected {
			t.Errorf("Expected getProductionNumber(%v) to be %v but got %v", test.input, test.expected, got)
		}
	}
}

func TestCompile(t *testing.T) {
	var g Grammar
	g.Start = "expr'"
	g.Productions = []Production{
		{"expr'", []grammarSymbol{"expr"}},
		{"expr", []grammarSymbol{"expr", "+", "term"}},
		{"expr", []grammarSymbol{"term"}},
		{"term", []grammarSymbol{"term", "*", "factor"}},
		{"term", []grammarSymbol{"factor"}},
		{"factor", []grammarSymbol{"number"}},
		{"factor", []grammarSymbol{"(", "expr", ")"}},
	}

	var testData = []struct {
		input    []grammarSymbol
		expected bool
	}{
		{[]grammarSymbol{"number", "+", "number", "$"}, true},
		{[]grammarSymbol{"number", "+", "number", "*", "number", "$"}, true},
		{[]grammarSymbol{"(", "number", "+", "number", "*", "number", "$"}, false},
		{[]grammarSymbol{"(", "number", "+", "number", "*", "number", ")", "$"}, true},
		{[]grammarSymbol{"(", "number", "+", "number", ")", "*", "number", ")", "$"}, false},
		{[]grammarSymbol{"(", "number", "+", "number", ")", "*", "(", "number", "*", "number", ")", "$"}, true},
	}

	for _, test := range testData {
		ps := g.compile()
		ps.parse(test.input)
		if ps.accepted != test.expected {
			t.Errorf("Expected parser to output %v on input %v, got %v",
				test.expected, test.input, ps.accepted)
		}
	}
}
