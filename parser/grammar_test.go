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

func TestComputeLrItemClosureSet(t *testing.T) {
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

	// The production order is jumbled up in some tests to test equality.
	var testData = []struct {
		input    lrItem
		expected lrItemSet
	}{
		{
			lrItem{g, g.productions[0], 0},
			lrItemSet{
				itemSet: []lrItem{
					{g, g.productions[0], 0},
					{g, g.productions[1], 0},
					{g, g.productions[2], 0},
					{g, g.productions[4], 0},
					{g, g.productions[3], 0},
					{g, g.productions[6], 0},
					{g, g.productions[5], 0},
				},
			},
		},
		{
			lrItem{g, g.productions[3], 0},
			lrItemSet{
				itemSet: []lrItem{
					{g, g.productions[6], 0},
					{g, g.productions[4], 0},
					{g, g.productions[5], 0},
					{g, g.productions[3], 0},
				},
			},
		},
		{
			lrItem{g, g.productions[4], 0},
			lrItemSet{
				itemSet: []lrItem{
					{g, g.productions[4], 0},
					{g, g.productions[6], 0},
					{g, g.productions[5], 0},
				},
			},
		},
		{
			lrItem{g, g.productions[1], 1},
			lrItemSet{
				itemSet: []lrItem{
					{g, g.productions[1], 1},
				},
			},
		},
		{
			lrItem{g, g.productions[3], 1},
			lrItemSet{
				itemSet: []lrItem{
					{g, g.productions[3], 1},
				},
			},
		},
		{
			lrItem{g, g.productions[3], 3},
			lrItemSet{
				itemSet: []lrItem{
					{g, g.productions[3], 3},
				},
			},
		},
	}

	for _, test := range testData {
		if got := test.input.computeClosureSet(); !got.equals(&test.expected) {
			t.Errorf("Expected computeClosureSet() = %v, got %v", test.expected, got)
		}
	}
}
