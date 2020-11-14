package ast

import (
	"testing"
)

func TestStackPeek(t *testing.T) {
	testData := []struct {
		stck     stack
		idx      int
		expected stackRecord
	}{
		{
			stack{
				stackRecord{0, "expr"},
				stackRecord{1, "+"},
				stackRecord{2, "term"},
			},
			0,
			stackRecord{2, "term"},
		},
		{
			stack{
				stackRecord{0, "expr"},
				stackRecord{1, "+"},
				stackRecord{2, "term"},
			},
			1,
			stackRecord{1, "+"},
		},
		{
			stack{
				stackRecord{0, "expr"},
				stackRecord{1, "+"},
				stackRecord{2, "term"},
			},
			2,
			stackRecord{0, "expr"},
		},
		{
			stack{
				stackRecord{0, "factor"},
			},
			0,
			stackRecord{0, "factor"},
		},
	}

	for _, test := range testData {
		if got := test.stck.peek(test.idx); got != test.expected {
			t.Errorf("Expected peek to get %v, got %v", test.expected, got)
		}
	}
}

func TestStackSymbolsOnTop(t *testing.T) {
	testData := []struct {
		stck     stack
		symbols  []string
		expected bool
	}{
		{
			stack{
				stackRecord{0, "expr"},
				stackRecord{1, "+"},
				stackRecord{2, "term"},
			},
			[]string{"term"},
			true,
		},
		{
			stack{
				stackRecord{0, "expr"},
				stackRecord{1, "+"},
				stackRecord{2, "term"},
			},
			[]string{"expr", "+"},
			false,
		},
		{
			stack{
				stackRecord{0, "expr"},
				stackRecord{1, "+"},
				stackRecord{2, "term"},
			},
			[]string{"expr", "+", "term"},
			true,
		},
		{
			stack{
				stackRecord{0, "factor"},
				stackRecord{1, "+"},
				stackRecord{2, "number"},
			},
			[]string{"number"},
			true,
		},
		{
			stack{
				stackRecord{0, "number"},
			},
			[]string{"factor", "+", "number"},
			false,
		},
	}

	for _, test := range testData {
		if got := test.stck.hasSymbolsOnTop(test.symbols); got != test.expected {
			t.Errorf("For %v.hasSymbolsOnTop(%v), expected %v, got %v", test.stck, test.symbols, test.expected, got)
		}
	}
}
