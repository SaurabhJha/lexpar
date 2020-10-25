package main

type nondeterministicFiniteAutomata struct {
	start           state
	final           state
	current         state
	transitionGraph nondeterministicGraph
}

func (n *nondeterministicFiniteAutomata) init(input transitionLabel) {
	n.start = 0
	n.final = 1
	n.current = 0

	g := make(nondeterministicGraph)
	g.addTransition(0, 1, input)
	n.transitionGraph = g
}

func (n *nondeterministicFiniteAutomata) incrementStatesBy(increment state) {
	n.start += increment
	n.final += increment
	n.current += increment
	n.transitionGraph.incrementStatesBy(increment)
}

func (n *nondeterministicFiniteAutomata) combineUsingUnion(o *nondeterministicFiniteAutomata) {
	n.incrementStatesBy(1)
	o.incrementStatesBy(n.final + 1)
	graphF := n.transitionGraph
	graphO := o.transitionGraph
	graphF.merge(&graphO)
	n.transitionGraph = graphF

	oldFStart := n.start
	oldOStart := o.start
	var newStart state
	if n.start < o.start {
		newStart = state(n.start - 1)
	} else {
		newStart = state(o.start - 1)
	}
	n.start = newStart
	n.transitionGraph.addTransition(n.start, oldFStart, "")
	n.transitionGraph.addTransition(n.start, oldOStart, "")

	oldFFinal := n.final
	oldOFinal := o.final
	var newFinal state
	if n.final < o.final {
		newFinal = state(o.final + 1)
	} else {
		newFinal = state(n.final + 1)
	}
	n.final = newFinal
	n.transitionGraph.addTransition(oldFFinal, n.final, "")
	n.transitionGraph.addTransition(oldOFinal, n.final, "")

	n.current = n.start
}

func (n *nondeterministicFiniteAutomata) combineUsingConcat(o *nondeterministicFiniteAutomata) {
	oldFFinal := n.final
	o.incrementStatesBy(oldFFinal + 1)
	graphF := n.transitionGraph
	graphO := o.transitionGraph
	graphF.merge(&graphO)
	n.transitionGraph = graphF
	n.transitionGraph.addTransition(oldFFinal, o.start, "")
	n.final = o.final
	n.current = n.start
}

func (n *nondeterministicFiniteAutomata) applyStar() {
	n.incrementStatesBy(1)
	newStart := n.start - 1
	newFinal := n.final + 1
	n.transitionGraph.addTransition(newStart, n.start, "")
	n.transitionGraph.addTransition(n.final, newFinal, "")
	n.transitionGraph.addTransition(n.final, n.start, "")
	n.transitionGraph.addTransition(newStart, newFinal, "")
	n.start = newStart
	n.final = newFinal
	n.current = n.start
}

func (n *nondeterministicFiniteAutomata) constructClosureSet(s state) setOfStates {
	states := make(setOfStates)
	states.add(s)
	for _, es := range n.transitionGraph[s][""] {
		closureSetEs := n.constructClosureSet(es)
		states.unionWith(&closureSetEs)
	}
	return states
}

func (n *nondeterministicFiniteAutomata) getOutgoingTransitionLabels(ss setOfStates) setOfTransitionLables {
	labels := make(setOfTransitionLables)
	for s := range ss {
		for tl := range n.transitionGraph[s] {
			if len(tl) != 0 {
				labels.add(tl)
			}
		}
	}
	return labels
}

func (n *nondeterministicFiniteAutomata) getNextDfaState(ss setOfStates, l transitionLabel) setOfStates {
	nextStates := make(setOfStates)
	for s := range ss {
		for _, ns := range n.transitionGraph[s][l] {
			nextStates.add(ns)
		}
	}

	nextDfaState := make(setOfStates)
	for s := range nextStates {
		for o := range n.constructClosureSet(s) {
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

func (n *nondeterministicFiniteAutomata) convertToDfa() deterministicFiniteAutomata {
	dfaGraph := make(deterministicGraph)
	q := make(queue, 0, 100)
	seen := make(seenStates, 0, 100)

	dfaStartState := n.constructClosureSet(n.start)
	q.enqueue(dfaStartState)
	seen.add(dfaStartState)
	for !q.empty() {
		currentDfaState := q.dequeue()
		for label := range n.getOutgoingTransitionLabels(currentDfaState) {
			nextDfaState := n.getNextDfaState(currentDfaState, label)
			if !seen.has(nextDfaState) {
				q.enqueue(nextDfaState)
				seen.add(nextDfaState)
			}
			dfaGraph.addTransition(seen.getStateNumber(currentDfaState), seen.getStateNumber(nextDfaState), label)
		}
	}

	finalStates := make(setOfStates)
	for s := range seen {
		if seen[s].has(n.final) {
			finalStates.add(state(s))
		}
	}

	dfa := new(deterministicFiniteAutomata)
	dfa.start = 0
	dfa.final = finalStates
	dfa.current = dfa.start
	dfa.transitionGraph = dfaGraph
	return *dfa
}

func (d *deterministicFiniteAutomata) move(input transitionLabel) {
	nextState, ok := d.transitionGraph[d.current][input]
	if !ok {
		d.dead = true
		d.accepted = false
		return
	}
	d.current = nextState
	if d.final.has(d.current) {
		d.accepted = true
	}
}

func (d *deterministicFiniteAutomata) reset() {
	d.dead = false
	d.accepted = false
	d.current = d.start
}
