package main

type state uint
type transitionLabel string

type graph map[state]map[transitionLabel][]state

func (g *graph) addTransition(start state, end state, input transitionLabel) {
	row, ok := (*g)[start]
	if !ok {
		row = make(map[transitionLabel][]state)
	}
	row[input] = append(row[input], end)
	(*g)[start] = row
}

func (g *graph) incrementStatesBy(increment state) {
	newG := make(graph)
	for start, row := range *g {
		for input, endSlice := range row {
			for _, end := range endSlice {
				newG.addTransition(start+increment, end+increment, input)
			}
		}
	}
	*g = newG
}

func (g *graph) merge(o *graph) {
	for start, row := range *o {
		for input, endSlice := range row {
			for _, end := range endSlice {
				g.addTransition(start, end, input)
			}
		}
	}
}
