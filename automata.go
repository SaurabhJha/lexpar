package main

type finiteAutomata struct {
	start           state
	final           state
	current         state
	transitionGraph graph
}

func (f *finiteAutomata) initialize(input transitionLabel) {
	f.start = 0
	f.final = 1
	f.current = 0

	g := make(graph)
	g.addTransition(0, 1, input)
	f.transitionGraph = g
}

func (f *finiteAutomata) incrementStatesBy(increment state) {
	f.start += increment
	f.final += increment
	f.current += increment
	f.transitionGraph.incrementStatesBy(increment)
}

func (f *finiteAutomata) combineUsingUnion(o *finiteAutomata) {
	f.incrementStatesBy(1)
	o.incrementStatesBy(f.final + 1)
	graphF := f.transitionGraph
	graphO := o.transitionGraph
	graphF.merge(&graphO)
	f.transitionGraph = graphF

	oldFStart := f.start
	oldOStart := o.start
	var newStart state
	if f.start < o.start {
		newStart = state(f.start - 1)
	} else {
		newStart = state(o.start - 1)
	}
	f.start = newStart
	f.transitionGraph.addTransition(f.start, oldFStart, "")
	f.transitionGraph.addTransition(f.start, oldOStart, "")

	oldFFinal := f.final
	oldOFinal := o.final
	var newFinal state
	if f.final < o.final {
		newFinal = state(o.final + 1)
	} else {
		newFinal = state(f.final + 1)
	}
	f.final = newFinal
	f.transitionGraph.addTransition(oldFFinal, f.final, "")
	f.transitionGraph.addTransition(oldOFinal, f.final, "")

	f.current = f.start
}

func (f *finiteAutomata) combineUsingConcat(o *finiteAutomata) {
	oldFFinal := f.final
	o.incrementStatesBy(oldFFinal + 1)
	graphF := f.transitionGraph
	graphO := o.transitionGraph
	graphF.merge(&graphO)
	f.transitionGraph = graphF
	f.transitionGraph.addTransition(oldFFinal, o.start, "")
	f.final = o.final
	f.current = f.start
}

func (f *finiteAutomata) applyStar() {
	f.incrementStatesBy(1)
	newStart := f.start - 1
	newFinal := f.final + 1
	f.transitionGraph.addTransition(newStart, f.start, "")
	f.transitionGraph.addTransition(f.final, newFinal, "")
	f.transitionGraph.addTransition(f.final, f.start, "")
	f.transitionGraph.addTransition(newStart, newFinal, "")
	f.start = newStart
	f.final = newFinal
	f.current = f.start
}
