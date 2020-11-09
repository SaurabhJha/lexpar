package parser

import (
	"reflect"
	"testing"
)

func TestLrItemNextSymbol(t *testing.T) {
	var g Grammar
	g.Productions = []Production{
		{"expr'", []grammarSymbol{"expr"}},
		{"expr", []grammarSymbol{"expr", "+", "number"}},
		{"expr", []grammarSymbol{"number"}},
	}

	var testData = []struct {
		item     lrItem
		expected grammarSymbol
	}{
		{lrItem{g, g.Productions[0], 1}, ""},
		{lrItem{g, g.Productions[1], 1}, "+"},
		{lrItem{g, g.Productions[2], 0}, "number"},
	}

	for _, test := range testData {
		if got := test.item.getNextSymbol(); got != test.expected {
			t.Errorf("Expected next symbol of %v to be %v, got %v", test.item, test.expected, got)
		}
	}
}

func TestLrItemGetNextItem(t *testing.T) {
	var g Grammar
	g.Productions = []Production{
		{"expr'", []grammarSymbol{"expr"}},
		{"expr", []grammarSymbol{"expr", "+", "number"}},
		{"expr", []grammarSymbol{"number"}},
	}

	var testData = []struct {
		inputItem   lrItem
		inputSymbol grammarSymbol
		expected    lrItem
	}{
		{lrItem{g, g.Productions[1], 0}, "expr", lrItem{g, g.Productions[1], 1}},
		{lrItem{g, g.Productions[2], 0}, "number", lrItem{g, g.Productions[2], 1}},
		{lrItem{g, g.Productions[2], 0}, "expr", lrItem{}},
	}

	for _, test := range testData {
		if got := test.inputItem.getNextItem(test.inputSymbol); !reflect.DeepEqual(got, test.expected) {
			t.Errorf("Expected %v.getNextItem = %v, got %v", test.inputItem, test.expected, got)
		}
	}
}

func TestLrItemEmpty(t *testing.T) {
	var g Grammar
	g.Productions = []Production{
		{"expr'", []grammarSymbol{"expr"}},
		{"expr", []grammarSymbol{"expr", "+", "number"}},
		{"expr", []grammarSymbol{"number"}},
	}

	var testData = []struct {
		inputItem lrItem
		expected  bool
	}{
		{lrItem{g, g.Productions[2], 1}, false},
		{lrItem{}, true},
	}

	for _, test := range testData {
		if got := test.inputItem.empty(); got != test.expected {
			t.Errorf("Expected inputItem.empty() to be %v but got %v", test.expected, got)
		}
	}
}

func TestLrItemSetHas(t *testing.T) {
	var g Grammar
	g.Productions = []Production{
		{"expr'", []grammarSymbol{"expr"}},
		{"expr", []grammarSymbol{"expr", "+", "number"}},
		{"expr", []grammarSymbol{"number"}},
	}

	var testData = []struct {
		itemSet  lrItemSet
		item     lrItem
		expected bool
	}{
		{
			lrItemSet{
				[]lrItem{
					{g, g.Productions[0], 0},
					{g, g.Productions[1], 0},
					{g, g.Productions[2], 0},
				},
			},
			lrItem{g, g.Productions[1], 0},
			true,
		},
		{
			lrItemSet{
				[]lrItem{
					{g, g.Productions[0], 1},
					{g, g.Productions[1], 0},
					{g, g.Productions[2], 2},
				},
			},
			lrItem{g, g.Productions[2], 2},
			true,
		},
		{
			lrItemSet{
				[]lrItem{
					{g, g.Productions[0], 1},
					{g, g.Productions[1], 0},
					{g, g.Productions[2], 2},
				},
			},
			lrItem{g, g.Productions[1], 0},
			true,
		},
		{
			lrItemSet{
				[]lrItem{
					{g, g.Productions[0], 1},
					{g, g.Productions[1], 0},
					{g, g.Productions[2], 2},
				},
			},
			lrItem{g, g.Productions[1], 1},
			false,
		},
	}

	for _, test := range testData {
		if test.itemSet.has(test.item) != test.expected {
			t.Errorf("expected %v.has(%v) to be %v", test.itemSet, test.item, test.expected)
		}
	}
}

func TestLrItemSetEquals(t *testing.T) {
	var g Grammar
	g.Productions = []Production{
		{"expr'", []grammarSymbol{"expr"}},
		{"expr", []grammarSymbol{"expr", "+", "number"}},
		{"expr", []grammarSymbol{"number"}},
	}

	var testData = []struct {
		thisItemSet  lrItemSet
		otherItemSet lrItemSet
		expected     bool
	}{
		{
			lrItemSet{
				[]lrItem{
					{g, g.Productions[0], 0},
					{g, g.Productions[2], 0},
				},
			},
			lrItemSet{
				[]lrItem{
					{g, g.Productions[0], 0},
					{g, g.Productions[2], 0},
				},
			},
			true,
		},
		{
			lrItemSet{
				[]lrItem{
					{g, g.Productions[0], 0},
					{g, g.Productions[2], 0},
				},
			},
			lrItemSet{
				[]lrItem{
					{g, g.Productions[0], 0},
					{g, g.Productions[1], 0},
				},
			},
			false,
		},
		{
			lrItemSet{
				[]lrItem{
					{g, g.Productions[0], 0},
					{g, g.Productions[1], 0},
					{g, g.Productions[2], 0},
				},
			},
			lrItemSet{
				[]lrItem{
					{g, g.Productions[0], 0},
					{g, g.Productions[2], 0},
					{g, g.Productions[1], 0},
				},
			},
			true,
		},
	}

	for _, test := range testData {
		if got := test.thisItemSet.equals(&test.otherItemSet); got != test.expected {
			t.Errorf("Expected equals to be %v, but got %v", test.expected, got)
		}
	}
}

func TestLrItemSetAdd(t *testing.T) {
	var g Grammar
	g.Productions = []Production{
		{"expr'", []grammarSymbol{"expr"}},
		{"expr", []grammarSymbol{"expr", "+", "number"}},
		{"expr", []grammarSymbol{"number"}},
	}

	var testData = []struct {
		itemSet  lrItemSet
		item     lrItem
		expected lrItemSet
	}{
		{
			lrItemSet{
				[]lrItem{
					{g, g.Productions[0], 0},
					{g, g.Productions[2], 0},
				},
			},
			lrItem{g, g.Productions[1], 0},
			lrItemSet{
				[]lrItem{
					{g, g.Productions[0], 0},
					{g, g.Productions[2], 0},
					{g, g.Productions[1], 0},
				},
			},
		},
		{
			lrItemSet{
				[]lrItem{
					{g, g.Productions[0], 0},
					{g, g.Productions[1], 0},
					{g, g.Productions[2], 0},
				},
			},
			lrItem{g, g.Productions[1], 0},
			lrItemSet{
				[]lrItem{
					{g, g.Productions[0], 0},
					{g, g.Productions[2], 0},
					{g, g.Productions[1], 0},
				},
			},
		},
	}

	for _, test := range testData {
		if test.itemSet.add(test.item); !test.itemSet.equals(&test.expected) {
			t.Errorf("Expected adding %v to result in %v, got %v", test.item, test.itemSet, test.expected)
		}
	}
}

func TestComputeLrItemSetNextSymbols(t *testing.T) {
	var g Grammar
	g.Productions = []Production{
		{"expr'", []grammarSymbol{"expr"}},
		{"expr", []grammarSymbol{"expr", "+", "number"}},
		{"expr", []grammarSymbol{"number"}},
	}

	var testData = []struct {
		itemSet  lrItemSet
		expected setOfSymbols
	}{
		{
			lrItemSet{
				[]lrItem{
					{g, g.Productions[0], 1},
					{g, g.Productions[2], 1},
				},
			},
			setOfSymbols{},
		},
		{
			lrItemSet{
				[]lrItem{
					{g, g.Productions[0], 0},
					{g, g.Productions[1], 0},
					{g, g.Productions[2], 0},
				},
			},
			setOfSymbols{"expr": true, "number": true},
		},
	}

	for _, test := range testData {
		if got := test.itemSet.getNextSymbols(); !reflect.DeepEqual(got, test.expected) {
			t.Errorf("Expected symbols out of item set %v to be %v, got %v", test.itemSet, test.expected, got)
		}
	}
}

func TestComputeLrItemSetMergeWith(t *testing.T) {
	var g Grammar
	g.Productions = []Production{
		{"expr'", []grammarSymbol{"expr"}},
		{"expr", []grammarSymbol{"expr", "+", "number"}},
		{"expr", []grammarSymbol{"number"}},
	}

	var testData = []struct {
		itemSet      lrItemSet
		otherItemSet lrItemSet
		expected     lrItemSet
	}{
		{
			lrItemSet{
				[]lrItem{
					{g, g.Productions[0], 1},
					{g, g.Productions[2], 1},
				},
			},
			lrItemSet{
				[]lrItem{
					{g, g.Productions[0], 0},
					{g, g.Productions[2], 1},
				},
			},
			lrItemSet{
				[]lrItem{
					{g, g.Productions[0], 0},
					{g, g.Productions[0], 1},
					{g, g.Productions[2], 1},
				},
			},
		},
	}

	for _, test := range testData {
		if test.itemSet.mergeWith(&test.otherItemSet); !test.itemSet.equals(&test.expected) {
			t.Errorf("Expected %v to equal %v", test.itemSet, test.expected)
		}
	}
}

func TestComputeLrItemClosureSet(t *testing.T) {
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

	// The production order is jumbled up in some tests to test equality.
	var testData = []struct {
		input    lrItem
		expected lrItemSet
	}{
		{
			lrItem{g, g.Productions[0], 0},
			lrItemSet{
				itemSet: []lrItem{
					{g, g.Productions[0], 0},
					{g, g.Productions[1], 0},
					{g, g.Productions[2], 0},
					{g, g.Productions[4], 0},
					{g, g.Productions[3], 0},
					{g, g.Productions[6], 0},
					{g, g.Productions[5], 0},
				},
			},
		},
		{
			lrItem{g, g.Productions[3], 0},
			lrItemSet{
				itemSet: []lrItem{
					{g, g.Productions[6], 0},
					{g, g.Productions[4], 0},
					{g, g.Productions[5], 0},
					{g, g.Productions[3], 0},
				},
			},
		},
		{
			lrItem{g, g.Productions[4], 0},
			lrItemSet{
				itemSet: []lrItem{
					{g, g.Productions[4], 0},
					{g, g.Productions[6], 0},
					{g, g.Productions[5], 0},
				},
			},
		},
		{
			lrItem{g, g.Productions[1], 1},
			lrItemSet{
				itemSet: []lrItem{
					{g, g.Productions[1], 1},
				},
			},
		},
		{
			lrItem{g, g.Productions[3], 1},
			lrItemSet{
				itemSet: []lrItem{
					{g, g.Productions[3], 1},
				},
			},
		},
		{
			lrItem{g, g.Productions[3], 3},
			lrItemSet{
				itemSet: []lrItem{
					{g, g.Productions[3], 3},
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

func TestComputeLrItemSetNext(t *testing.T) {
	var g Grammar
	g.Productions = []Production{
		{"expr'", []grammarSymbol{"expr"}},
		{"expr", []grammarSymbol{"expr", "+", "number"}},
		{"expr", []grammarSymbol{"number"}},
	}

	var testData = []struct {
		itemSet  lrItemSet
		symbol   grammarSymbol
		expected lrItemSet
	}{
		{
			lrItemSet{
				[]lrItem{
					{g, g.Productions[0], 0},
					{g, g.Productions[1], 0},
					{g, g.Productions[2], 0},
				},
			},
			"expr",
			lrItemSet{
				[]lrItem{
					{g, g.Productions[0], 1},
					{g, g.Productions[1], 1},
				},
			},
		},
		{
			lrItemSet{
				[]lrItem{
					{g, g.Productions[1], 1},
				},
			},
			"expr",
			lrItemSet{
				[]lrItem{},
			},
		},
	}

	for _, test := range testData {
		if got := test.itemSet.getNextItemSet(test.symbol); !got.equals(&test.expected) {
			t.Errorf(
				"Expected %v.nextItemSet(%v) to be %v, got %v", test.itemSet, test.symbol, test.expected, got)
		}
	}
}

func TestParsingTableAddShift(t *testing.T) {
	table := make(parsingTable)
	table.addShiftMove(0, 1, "a")
	table.addShiftMove(0, 2, "b")

	if got := table[0]["a"]; !reflect.DeepEqual(got, parserAction{shift, 1}) {
		t.Errorf("On %v and %v, expected %v, got %v", 0, "a", parserAction{shift, 1}, got)
	}
	if got := table[0]["b"]; !reflect.DeepEqual(got, parserAction{shift, 2}) {
		t.Errorf("On %v and %v, expected %v, got %v", 0, "a", parserAction{shift, 2}, got)
	}
}

func TestParsingTableAddReduce(t *testing.T) {
	table := make(parsingTable)
	table.addReduceMove(0, 1, "a")
	table.addReduceMove(0, 2, "b")

	if got := table[0]["a"]; !reflect.DeepEqual(got, parserAction{reduce, 1}) {
		t.Errorf("On %v and %v, expected %v, got %v", 0, "a", parserAction{reduce, 1}, got)
	}
	if got := table[0]["b"]; !reflect.DeepEqual(got, parserAction{reduce, 2}) {
		t.Errorf("On %v and %v, expected %v, got %v", 0, "a", parserAction{reduce, 2}, got)
	}
}

func TestParsingTableAddAccept(t *testing.T) {
	table := make(parsingTable)
	table.addAcceptMove(3)

	if got := table[3]["$"]; !reflect.DeepEqual(got, parserAction{accept, 0}) {
		t.Errorf("On %v and %v, expected %v, got %v", 3, "$", parserAction{accept, 0}, got)
	}
}
