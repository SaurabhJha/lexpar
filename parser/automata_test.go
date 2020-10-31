package parser

import "testing"

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
