package parser

import "testing"

func TestLrItemEmpty(t *testing.T) {
	var g grammar
	g.productions = []production{
		{"expr'", []grammarSymbol{"expr"}},
		{"expr", []grammarSymbol{"expr", "+", "number"}},
		{"expr", []grammarSymbol{"number"}},
	}

	var testData = []struct {
		inputItem lrItem
		expected  bool
	}{
		{lrItem{g, g.productions[2], 1}, false},
		{lrItem{}, true},
	}

	for _, test := range testData {
		if got := test.inputItem.empty(); got != test.expected {
			t.Errorf("Expected inputItem.empty() to be %v but got %v", test.expected, got)
		}
	}
}

func TestLrItemGetNextItem(t *testing.T) {
	var g grammar
	g.productions = []production{
		{"expr'", []grammarSymbol{"expr"}},
		{"expr", []grammarSymbol{"expr", "+", "number"}},
		{"expr", []grammarSymbol{"number"}},
	}

	var testData = []struct {
		inputItem   lrItem
		inputSymbol grammarSymbol
		expected    lrItem
	}{
		{lrItem{g, g.productions[1], 0}, "expr", lrItem{g, g.productions[1], 1}},
		{lrItem{g, g.productions[2], 0}, "number", lrItem{g, g.productions[2], 1}},
		{lrItem{g, g.productions[2], 0}, "expr", lrItem{}},
	}

	for _, test := range testData {
		if got := test.inputItem.getNextItem(test.inputSymbol); !got.equals(test.expected) {
			t.Errorf("Expected %v.getNextItem = %v, got %v", test.inputItem, test.expected, got)
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

func TestItemSetGetNextItemSet(t *testing.T) {
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
		inputItemSet lrItemSet
		inputSymbol  grammarSymbol
		expected     lrItemSet
	}{
		{
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
			"expr",
			lrItemSet{
				itemSet: []lrItem{
					{g, g.productions[0], 1},
					{g, g.productions[1], 1},
				},
			},
		},
		{
			lrItemSet{
				itemSet: []lrItem{
					{g, g.productions[0], 1},
					{g, g.productions[1], 1},
				},
			},
			"+",
			lrItemSet{
				itemSet: []lrItem{
					{g, g.productions[1], 2},
					{g, g.productions[3], 0},
					{g, g.productions[4], 0},
					{g, g.productions[5], 0},
					{g, g.productions[6], 0},
				},
			},
		},
		{
			lrItemSet{
				itemSet: []lrItem{
					{g, g.productions[6], 0},
				},
			},
			"(",
			lrItemSet{
				itemSet: []lrItem{
					{g, g.productions[6], 1},
					{g, g.productions[1], 0},
					{g, g.productions[2], 0},
					{g, g.productions[3], 0},
					{g, g.productions[4], 0},
					{g, g.productions[5], 0},
					{g, g.productions[6], 0},
				},
			},
		},
	}

	for _, test := range testData {
		if got := test.inputItemSet.getNextItemSet(test.inputSymbol); !got.equals(&test.expected) {
			t.Errorf("Expected getNextItemSet(%v) to equal %v but got %v", test.inputSymbol, test.expected.itemSet, got.itemSet)
		}
	}
}
