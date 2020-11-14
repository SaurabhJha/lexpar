package lexer

import (
	"math"
)

// nondeterministicFiniteAutomata represent NFAs. They are the intermediate step in regex compilation
// process and cannot be executed
type nondeterministicFiniteAutomata struct {
	start           state
	final           state // Its guaranteed by the way we construct NFAs that we will have only one final state.
	transitionGraph nondeterministicGraph
	closureSets     map[state]setOfStates
}

func (nfa *nondeterministicFiniteAutomata) init(input transitionLabel) {
	nfa.start = 0
	nfa.final = 1

	graph := make(nondeterministicGraph)
	graph.addTransition(0, 1, input)
	nfa.transitionGraph = graph
	nfa.closureSets = make(map[state]setOfStates)
}

func (nfa *nondeterministicFiniteAutomata) incrementStatesBy(increment state) {
	nfa.start += increment
	nfa.final += increment
	nfa.transitionGraph.incrementStatesBy(increment)
}

func (nfa *nondeterministicFiniteAutomata) combineUsingUnion(otherNfa *nondeterministicFiniteAutomata) {
	nfa.incrementStatesBy(1)
	otherNfa.incrementStatesBy(nfa.final + 1)
	nfa.transitionGraph.merge(&otherNfa.transitionGraph)

	oldNStart := nfa.start
	oldOStart := otherNfa.start
	nfa.start = state(math.Min(float64(nfa.start), float64(otherNfa.start))) - 1
	nfa.transitionGraph.addTransition(nfa.start, oldNStart, "")
	nfa.transitionGraph.addTransition(nfa.start, oldOStart, "")

	oldFFinal := nfa.final
	oldOFinal := otherNfa.final
	nfa.final = state(math.Max(float64(nfa.final), float64(otherNfa.final))) + 1
	nfa.transitionGraph.addTransition(oldFFinal, nfa.final, "")
	nfa.transitionGraph.addTransition(oldOFinal, nfa.final, "")
}

func (nfa *nondeterministicFiniteAutomata) combineUsingConcat(otherNfa *nondeterministicFiniteAutomata) {
	otherNfa.incrementStatesBy(nfa.final + 1)
	nfa.transitionGraph.merge(&otherNfa.transitionGraph)
	nfa.transitionGraph.addTransition(nfa.final, otherNfa.start, "")
	nfa.final = otherNfa.final
}

func (nfa *nondeterministicFiniteAutomata) applyStar() {
	nfa.incrementStatesBy(1)
	nfa.transitionGraph.addTransition(nfa.start-1, nfa.start, "")
	nfa.transitionGraph.addTransition(nfa.final, nfa.final+1, "")
	nfa.transitionGraph.addTransition(nfa.final, nfa.start, "")
	nfa.transitionGraph.addTransition(nfa.start-1, nfa.final+1, "")
	nfa.start--
	nfa.final++
}

func (nfa *nondeterministicFiniteAutomata) constructClosureSet(s state) setOfStates {
	if states, ok := nfa.closureSets[s]; ok {
		return states
	}
	states := make(setOfStates)
	states.add(s)
	for _, es := range nfa.transitionGraph[s][""] {
		closureSetEs := nfa.constructClosureSet(es)
		states.unionWith(&closureSetEs)
	}
	nfa.closureSets[s] = states
	return states
}

func (nfa *nondeterministicFiniteAutomata) getOutgoingTransitionLabels(ss setOfStates) setOfTransitionLables {
	labels := make(setOfTransitionLables)
	for s := range ss {
		for tl := range nfa.transitionGraph[s] {
			if tl != "" {
				labels.add(tl)
			}
		}
	}
	return labels
}

func (nfa *nondeterministicFiniteAutomata) getNextDfaState(ss setOfStates, l transitionLabel) setOfStates {
	nextStates := make(setOfStates)
	for s := range ss {
		for _, ns := range nfa.transitionGraph[s][l] {
			nextStates.add(ns)
		}
	}

	nextDfaState := make(setOfStates)
	for s := range nextStates {
		for o := range nfa.constructClosureSet(s) {
			nextDfaState.add(o)
		}
	}

	return nextDfaState
}

type deterministicFiniteAutomata struct {
	start           state
	final           setOfStates
	current         state
	transitionGraph deterministicGraph
	dead            bool
	accepted        bool
}

func (nfa *nondeterministicFiniteAutomata) convertToDfa() deterministicFiniteAutomata {
	dfaGraph := make(deterministicGraph)
	q := make(queue, 0, 100)
	seen := make(seenStates, 0, 100)

	dfaStartState := nfa.constructClosureSet(nfa.start)
	q.enqueue(dfaStartState)
	seen.add(dfaStartState)
	for !q.empty() {
		currentDfaState := q.dequeue()
		for label := range nfa.getOutgoingTransitionLabels(currentDfaState) {
			nextDfaState := nfa.getNextDfaState(currentDfaState, label)
			if !seen.has(nextDfaState) {
				q.enqueue(nextDfaState)
				seen.add(nextDfaState)
			}
			dfaGraph.addTransition(seen.getStateNumber(currentDfaState), seen.getStateNumber(nextDfaState), label)
		}
	}

	finalStates := make(setOfStates)
	for s := range seen {
		if seen[s].has(nfa.final) {
			finalStates.add(state(s))
		}
	}

	return deterministicFiniteAutomata{start: 0, final: finalStates, current: 0, transitionGraph: dfaGraph}
}

func (d *deterministicFiniteAutomata) move(input transitionLabel) {
	if d.dead {
		return
	}

	nextState, ok := d.transitionGraph[d.current][input]
	if !ok {
		d.dead, d.accepted = true, false
		return
	}
	d.current = nextState
	if d.final.has(d.current) {
		d.accepted = true
	}
}

func (d *deterministicFiniteAutomata) reset() {
	d.dead, d.accepted, d.current = false, false, d.start
}
