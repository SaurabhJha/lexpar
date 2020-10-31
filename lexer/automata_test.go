package lexer

import (
	"testing"
)

func TestDeterministicFiniteAutomataMove(t *testing.T) {
	var testData = []struct {
		inputRegex regularExpression
		testInput  string
		expected   bool
	}{
		{"dfa", "dfa", true},
		{"(a|b)(c|d)", "bc", true},
		{"(a|b)(c|d)", "bbc", false},
		{"abd*", "aaaabbbbddddd", false},
		{"abd*", "ab", true},
		{"abd*", "abddd", true},
		{"a(b|c)*", "abccbbc", true},
		{"a(b|c)*", "a", true},
		{"a(b|c)*", "abccdfdf", false},
	}

	for _, test := range testData {
		nfa := test.inputRegex.compile()
		dfa := nfa.convertToDfa()
		for _, character := range test.testInput {
			dfa.move(transitionLabel(character))
		}
		if dfa.accepted != test.expected {
			t.Errorf("expected dfa to accept %v", test.testInput)
		}
	}
}
