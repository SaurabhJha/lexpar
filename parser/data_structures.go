package parser

import (
	"reflect"
)

type setOfSymbols map[grammarSymbol]bool

func (ss *setOfSymbols) add(s grammarSymbol) {
	(*ss)[s] = true
}

func (ss *setOfSymbols) unionWith(os *setOfSymbols) {
	for s := range *os {
		(*ss).add(s)
	}
}

func (ss *setOfSymbols) has(s grammarSymbol) bool {
	return (*ss)[s]
}

func (ss *setOfSymbols) hasSubset(os *setOfSymbols) bool {
	for s := range *os {
		if !(*ss).has(s) {
			return false
		}
	}
	return true
}

type queueOfItems []lrItem

func (q *queueOfItems) enqueue(l lrItem) {
	*q = append(*q, l)
}

func (q *queueOfItems) dequeue() lrItem {
	front := (*q)[0]
	(*q) = (*q)[1:]
	return front
}

func (q *queueOfItems) empty() bool {
	return len(*q) == 0
}

type queueOfItemSets []lrItemSet

func (q *queueOfItemSets) enqueue(ls lrItemSet) {
	*q = append(*q, ls)
}

func (q *queueOfItemSets) dequeue() lrItemSet {
	front := (*q)[0]
	(*q) = (*q)[1:]
	return front
}

func (q *queueOfItemSets) empty() bool {
	return len(*q) == 0
}

type seenLrItemSets []lrItemSet

func (lss *seenLrItemSets) add(ls lrItemSet) {
	*lss = append(*lss, ls)
}

func (lss *seenLrItemSets) has(ls lrItemSet) bool {
	for _, sls := range *lss {
		if reflect.DeepEqual(sls, ls) {
			return true
		}
	}
	return false
}

func (lss *seenLrItemSets) getStateNumber(ls lrItemSet) state {
	for i, sls := range *lss {
		if reflect.DeepEqual(sls, ls) {
			return state(i)
		}
	}
	lss.add(ls)
	return state(len(*lss) - 1)
}

type parserStack []state

func (ss *parserStack) push(s state) {
	(*ss) = append(*ss, s)
}

func (ss *parserStack) top() state {
	return (*ss)[len(*ss)-1]
}

func (ss *parserStack) pop() state {
	top := ss.top()
	*ss = (*ss)[:len(*ss)-1]
	return top
}

func (ss *parserStack) empty() bool {
	return len(*ss) == 0
}

// SyntaxGraph is a data structure representation of a program text. It is produced by
// Parser.
type SyntaxGraph struct {
	Graph     map[int][]int
	NodeLabel []string
	Root      int
}

func (ast *SyntaxGraph) createNewNode(lexeme string) int {
	ast.NodeLabel = append(ast.NodeLabel, lexeme)
	return len(ast.NodeLabel) - 1
}

func (ast *SyntaxGraph) addEdge(start int, end int) {
	if ast.Graph == nil {
		ast.Graph = make(map[int][]int)
	}
	ast.Graph[start] = append(ast.Graph[start], end)
}

type graphStack []int

func (gs *graphStack) push(n int) {
	(*gs) = append(*gs, n)
}

func (gs *graphStack) pop() int {
	top := gs.top()
	*gs = (*gs)[:len(*gs)-1]
	return top
}

func (gs *graphStack) top() int {
	return (*gs)[len(*gs)-1]
}
