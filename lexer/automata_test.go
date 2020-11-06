package lexer

import (
	"reflect"
	"testing"
)

func TestNonDeterministicFiniteAutomataInit(t *testing.T) {
	var nfa nondeterministicFiniteAutomata
	nfa.init("a")

	if nfa.start != 0 {
		t.Errorf("Expected start to be 0 but got %v", nfa.start)
	}
	if nfa.final != 1 {
		t.Errorf("Expected final to be 1 but got %v", nfa.final)
	}
	if got := len(nfa.transitionGraph[nfa.start]); got != 1 {
		t.Errorf("Expected start state to move to only one state but is moving to %v states", got)
	}
	if got := nfa.transitionGraph[nfa.start]["a"][0]; got != nfa.final {
		t.Errorf("Expected start state to move to state %v but is moving to state %v", nfa.final, got)
	}
}

func TestNonDeterministicFiniteAutomataIncrementStates(t *testing.T) {
	var nfa nondeterministicFiniteAutomata
	nfa.init("a")
	nfa.incrementStatesBy(3)

	if nfa.start != 3 {
		t.Errorf("Expected start to be 3 but got %v", nfa.start)
	}
	if nfa.final != 4 {
		t.Errorf("Expected final to be 4 but got %v", nfa.final)
	}
	if got := len(nfa.transitionGraph[nfa.start]); got != 1 {
		t.Errorf("Expected start state to move to only one state but is moving to %v states", got)
	}
	if got := nfa.transitionGraph[nfa.start]["a"][0]; got != nfa.final {
		t.Errorf("Expected start state to move to state %v but is moving to state %v", nfa.final, got)
	}
}

func TestNonDeterministicFiniteAutomataUnion(t *testing.T) {
	var nfa1, nfa2 nondeterministicFiniteAutomata
	nfa1.init("a")
	nfa2.init("b")
	nfa1.combineUsingUnion(&nfa2)

	if nfa1.start != 0 {
		t.Errorf("Expected start state to be 0, got %v", nfa1.start)
	}
	if nfa1.final != 5 {
		t.Errorf("Expected final state to be 5, got %v", nfa1.final)
	}

	var testData = []struct {
		s        state
		t        transitionLabel
		expected []state
	}{
		{0, "", []state{1, 3}},
		{1, "a", []state{2}},
		{3, "b", []state{4}},
		{2, "", []state{5}},
		{4, "", []state{5}},
	}

	for _, test := range testData {
		if got := nfa1.transitionGraph[test.s][test.t]; !reflect.DeepEqual(got, test.expected) {
			t.Errorf("On state %v and input %v, expected %v but got %v",
				test.s, test.t, test.expected, got)
		}
	}
}

func TestNonDeterministicFiniteAutomataConcat(t *testing.T) {
	var nfa1, nfa2 nondeterministicFiniteAutomata
	nfa1.init("a")
	nfa2.init("b")
	nfa1.combineUsingConcat(&nfa2)

	if nfa1.start != 0 {
		t.Errorf("Expected start state to be 0, got %v", nfa1.start)
	}
	if nfa1.final != 3 {
		t.Errorf("Expected final state to be 3, got %v", nfa1.final)
	}

	var testData = []struct {
		s        state
		t        transitionLabel
		expected []state
	}{
		{0, "a", []state{1}},
		{1, "", []state{2}},
		{2, "b", []state{3}},
	}

	for _, test := range testData {
		if got := nfa1.transitionGraph[test.s][test.t]; !reflect.DeepEqual(got, test.expected) {
			t.Errorf("On state %v and input %v, expected %v but got %v",
				test.s, test.t, test.expected, got)
		}
	}
}

func TestNonDeterministicFiniteAutomataStar(t *testing.T) {
	var nfa nondeterministicFiniteAutomata
	nfa.init("a")
	nfa.applyStar()

	if nfa.start != 0 {
		t.Errorf("Expected start state to be 0, got %v", nfa.start)
	}
	if nfa.final != 3 {
		t.Errorf("Expected final state to be 3, got %v", nfa.final)
	}

	var testData = []struct {
		s        state
		t        transitionLabel
		expected []state
	}{
		{0, "", []state{1, 3}},
		{1, "a", []state{2}},
		{2, "", []state{3, 1}},
	}

	for _, test := range testData {
		if got := nfa.transitionGraph[test.s][test.t]; !reflect.DeepEqual(got, test.expected) {
			t.Errorf("On state %v and input %v, expected %v but got %v",
				test.s, test.t, test.expected, got)
		}
	}
}

func TestNonDeterministicFiniteAutomataClosure(t *testing.T) {
	var nfa1, nfa2 nondeterministicFiniteAutomata
	nfa1.init("a")
	nfa2.init("b")
	nfa1.combineUsingUnion(&nfa2)
	nfa1.applyStar()

	var testData = []struct {
		s        state
		expected setOfStates
	}{
		{0, setOfStates{0: true, 1: true, 2: true, 4: true, 7: true}},
		{1, setOfStates{1: true, 2: true, 4: true}},
		{2, setOfStates{2: true}},
		{3, setOfStates{3: true, 6: true, 1: true, 2: true, 4: true, 7: true}},
		{4, setOfStates{4: true}},
		{5, setOfStates{5: true, 6: true, 1: true, 2: true, 4: true, 7: true}},
		{6, setOfStates{6: true, 1: true, 2: true, 4: true, 7: true}},
		{7, setOfStates{7: true}},
	}

	for _, test := range testData {
		if got := nfa1.constructClosureSet(test.s); !reflect.DeepEqual(got, test.expected) {
			t.Errorf("Expected closure of %v to be %v but got %v", test.s, test.expected, got)
		}
	}
}

func TestNonDeterministicFiniteAutomataTransitionSymbols(t *testing.T) {
	var nfa1, nfa2 nondeterministicFiniteAutomata
	nfa1.init("a")
	nfa2.init("b")
	nfa1.combineUsingUnion(&nfa2)
	nfa1.applyStar()

	var testData = []struct {
		ss       setOfStates
		expected setOfTransitionLables
	}{
		{setOfStates{1: true}, setOfTransitionLables{}},
		{setOfStates{2: true}, setOfTransitionLables{"a": true}},
		{setOfStates{4: true}, setOfTransitionLables{"b": true}},
		{setOfStates{1: true, 2: true, 4: true}, setOfTransitionLables{"a": true, "b": true}},
	}

	for _, test := range testData {
		if got := nfa1.getOutgoingTransitionLabels(test.ss); !reflect.DeepEqual(got, test.expected) {
			t.Errorf("Expected outgoing edges from %v with labels %v but got %v",
				test.ss, test.expected, got)
		}
	}
}

func TestNonDeterministicFiniteAutomataNextDfaState(t *testing.T) {
	var nfa1, nfa2 nondeterministicFiniteAutomata
	nfa1.init("a")
	nfa2.init("b")
	nfa1.combineUsingUnion(&nfa2)
	nfa1.applyStar()

	var testData = []struct {
		currentSs setOfStates
		input     transitionLabel
		expected  setOfStates
	}{
		{
			setOfStates{0: true, 1: true, 2: true, 4: true, 7: true},
			"a",
			setOfStates{3: true, 6: true, 1: true, 2: true, 4: true, 7: true},
		},
		{
			setOfStates{0: true, 1: true, 2: true, 4: true, 7: true},
			"b",
			setOfStates{5: true, 6: true, 1: true, 2: true, 4: true, 7: true},
		},
		{
			setOfStates{1: true, 2: true, 4: true},
			"a",
			setOfStates{3: true, 6: true, 1: true, 2: true, 4: true, 7: true},
		},
		{
			setOfStates{1: true, 2: true, 4: true},
			"b",
			setOfStates{5: true, 6: true, 1: true, 2: true, 4: true, 7: true},
		},
		{
			setOfStates{3: true, 6: true, 1: true, 2: true, 4: true, 7: true},
			"a",
			setOfStates{3: true, 6: true, 1: true, 2: true, 4: true, 7: true},
		},
		{
			setOfStates{3: true, 6: true, 1: true, 2: true, 4: true, 7: true},
			"b",
			setOfStates{5: true, 6: true, 1: true, 2: true, 4: true, 7: true},
		},
		{
			setOfStates{5: true, 6: true, 1: true, 2: true, 4: true, 7: true},
			"a",
			setOfStates{3: true, 6: true, 1: true, 2: true, 4: true, 7: true},
		},
		{
			setOfStates{5: true, 6: true, 1: true, 2: true, 4: true, 7: true},
			"b",
			setOfStates{5: true, 6: true, 1: true, 2: true, 4: true, 7: true},
		},
	}

	for _, test := range testData {
		if got := nfa1.getNextDfaState(test.currentSs, test.input); !reflect.DeepEqual(got, test.expected) {
			t.Errorf("Expected next DFA state on state %v and input %v to be %v but got %v",
				test.currentSs, test.input, test.expected, got)
		}
	}
}

func TestNonDeterministicFiniteAutomataConvertToDfa(t *testing.T) {
	var nfa1, nfa2 nondeterministicFiniteAutomata
	nfa1.init("a")
	nfa2.init("b")
	nfa1.combineUsingUnion(&nfa2)
	nfa1.applyStar()

	dfa := nfa1.convertToDfa()
	// Need to do this because the order of getting transition symbols is non deterministic
	stateOnA := dfa.transitionGraph[0]["a"]
	stateOnB := dfa.transitionGraph[0]["b"]

	var testData = []struct {
		startState       state
		input            transitionLabel
		expectedEndState state
	}{
		{0, "a", stateOnA},
		{0, "b", stateOnB},
		{1, "a", stateOnA},
		{1, "b", stateOnB},
		{2, "a", stateOnA},
		{2, "b", stateOnB},
	}
	for _, test := range testData {
		if got := dfa.transitionGraph[test.startState][test.input]; got != test.expectedEndState {
			t.Errorf(
				"Expected next state on %v, %v to be %v, got %v", test.startState, test.input, test.expectedEndState, got)
		}
	}
	if dfa.start != 0 {
		t.Errorf("Expected start to be %v, got %v", 0, dfa.start)
	}
	if !dfa.final.has(1) && !dfa.final.has(2) {
		t.Errorf("Expected final state to be %v but got %v", setOfStates{1: true, 2: true}, dfa.final)
	}
}

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
