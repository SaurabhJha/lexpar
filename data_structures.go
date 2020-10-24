package main

import "reflect"

type stack []byte

func (s *stack) push(b byte) {
	(*s) = append(*s, b)
}

func (s *stack) pop() byte {
	top := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return top
}

func (s *stack) empty() bool {
	return len(*s) == 0
}

type state uint

type transitionLabel string

type nondeterministicGraph map[state]map[transitionLabel][]state

func (g *nondeterministicGraph) addTransition(s state, e state, l transitionLabel) {
	row, ok := (*g)[s]
	if !ok {
		row = make(map[transitionLabel][]state)
	}
	row[l] = append(row[l], e)
	(*g)[s] = row
}

func (g *nondeterministicGraph) incrementStatesBy(increment state) {
	newG := make(nondeterministicGraph)
	for start, row := range *g {
		for input, endSlice := range row {
			for _, end := range endSlice {
				newG.addTransition(start+increment, end+increment, input)
			}
		}
	}
	*g = newG
}

func (g *nondeterministicGraph) merge(o *nondeterministicGraph) {
	for start, row := range *o {
		for input, endSlice := range row {
			for _, end := range endSlice {
				g.addTransition(start, end, input)
			}
		}
	}
}

type setOfStates map[state]bool

func (ss *setOfStates) add(s state) {
	(*ss)[s] = true
}

func (ss *setOfStates) unionWith(os *setOfStates) {
	for s := range *os {
		(*ss).add(s)
	}
}

func (ss *setOfStates) has(s state) bool {
	_, ok := (*ss)[s]
	return ok
}

type setOfTransitionLables map[transitionLabel]bool

func (st *setOfTransitionLables) add(l transitionLabel) {
	(*st)[l] = true
}

type queue []setOfStates

func (q *queue) enqueue(s setOfStates) {
	*q = append(*q, s)
}

func (q *queue) dequeue() setOfStates {
	front := (*q)[0]
	(*q) = (*q)[1:]
	return front
}

func (q *queue) empty() bool {
	return len(*q) == 0
}

type deterministicGraph map[state]map[transitionLabel]state

func (d *deterministicGraph) addTransition(s state, e state, l transitionLabel) {
	if (*d)[s] == nil {
		(*d)[s] = make(map[transitionLabel]state)
	}
	(*d)[s][l] = e
}

type seenStates []setOfStates

func (s *seenStates) add(ss setOfStates) {
	*s = append(*s, ss)
}

func (s *seenStates) has(ss setOfStates) bool {
	for i := range *s {
		if reflect.DeepEqual((*s)[i], ss) {
			return true
		}
	}
	return false
}

func (s *seenStates) getStateNumber(ss setOfStates) state {
	for i := range *s {
		if reflect.DeepEqual((*s)[i], ss) {
			return state(i)
		}
	}
	*s = append(*s, ss)
	return state(len(*s) - 1)
}
