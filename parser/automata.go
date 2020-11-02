package parser

import "fmt"

type lrItem struct {
	g   grammar
	p   production
	pos int
}

func (l lrItem) getNextSymbol() grammarSymbol {
	if l.pos == len(l.p.body) {
		return ""
	}
	return l.p.body[l.pos]
}

func (l lrItem) getNextItem(s grammarSymbol) lrItem {
	var newL lrItem

	if len(l.p.body) == l.pos {
		return lrItem{}
	}

	if l.p.body[l.pos] != s {
		return lrItem{}
	}

	newL.g = l.g
	newL.p = l.p
	newL.pos = l.pos + 1
	return newL
}

func (l lrItem) empty() bool {
	return l.g.start == "" && l.g.productions == nil && l.p.head == "" && l.p.body == nil
}

func (l lrItem) equals(l2 lrItem) bool {
	if !l.g.equals(l2.g) {
		return false
	}

	if !l.p.equals(l2.p) {
		return false
	}

	if l.pos != l2.pos {
		return false
	}

	return true
}

func (l lrItem) computeClosureSet() lrItemSet {
	ls := lrItemSet{[]lrItem{l}}

	if l.g.isTerminal(l.getNextSymbol()) {
		return ls
	}

	q := make(queueOfItems, 0, 10)
	q.enqueue(l)
	seenNonTerminals := make(setOfSymbols)
	for !q.empty() {
		currentItem := q.dequeue()
		nextSymbol := currentItem.getNextSymbol()
		if !l.g.isTerminal(nextSymbol) && !seenNonTerminals.has(nextSymbol) {
			seenNonTerminals.add(nextSymbol)
			for _, p := range l.g.getProductionsOfSymbol(nextSymbol) {
				nextItem := lrItem{l.g, p, 0}
				if !ls.has(nextItem) {
					q.enqueue(nextItem)
					ls.add(nextItem)
				}
			}
		}
	}

	return ls
}

type lrItemSet struct {
	itemSet []lrItem
}

func (ls *lrItemSet) has(l lrItem) bool {
	for _, item := range ls.itemSet {
		if item.equals(l) {
			return true
		}
	}
	return false
}

func (ls *lrItemSet) add(l lrItem) {
	if !ls.has(l) {
		(*ls).itemSet = append((*ls).itemSet, l)
	}
}

func (ls *lrItemSet) equals(ls2 *lrItemSet) bool {
	if len(ls.itemSet) != len(ls2.itemSet) {
		return false
	}

	for _, l := range ls.itemSet {
		if !ls2.has(l) {
			return false
		}
	}

	for _, l := range ls2.itemSet {
		if !ls.has(l) {
			return false
		}
	}

	return true
}

func (ls *lrItemSet) getNextSymbols() setOfSymbols {
	s := make(setOfSymbols)

	for _, l := range ls.itemSet {
		if l.getNextSymbol() != "" {
			s.add(l.getNextSymbol())
		}
	}

	return s
}

func (ls *lrItemSet) mergeWith(otherLs *lrItemSet) {
	for _, l := range otherLs.itemSet {
		if !ls.has(l) {
			ls.add(l)
		}
	}
}

func (ls *lrItemSet) getNextItemSet(s grammarSymbol) lrItemSet {
	var nextLs lrItemSet
	for _, l := range ls.itemSet {
		nextItem := l.getNextItem(s)
		if !nextItem.empty() {
			nextItemSet := nextItem.computeClosureSet()
			nextLs.mergeWith(&nextItemSet)
		}
	}

	return nextLs
}

type state uint

type parserActionType int

const (
	shift parserActionType = iota
	reduce
	accept
)

type parserAction struct {
	actionType parserActionType
	number     int
}

type parsingTable map[state]map[grammarSymbol]parserAction

func (p *parsingTable) addShiftMove(s state, e state, gs grammarSymbol) {
	if (*p)[s] == nil {
		(*p)[s] = make(map[grammarSymbol]parserAction)
	}

	if existingAction, ok := (*p)[s][gs]; ok {
		switch existingAction.actionType {
		case shift:
			panic(fmt.Sprintf("Shift-shift conflict on state %v and input %v", s, gs))
		case reduce:
			panic(fmt.Sprintf("Shift-reduce conflict on state %v and input %v", s, gs))
		case accept:
			panic(fmt.Sprintf("Accept-shift conflict on state %v and input %v", s, gs))
		}
	}

	(*p)[s][gs] = parserAction{shift, int(e)}
}

func (p *parsingTable) addReduceMove(s state, productionNumber int, gs grammarSymbol) {
	if (*p)[s] == nil {
		(*p)[s] = make(map[grammarSymbol]parserAction)
	}

	if existingAction, ok := (*p)[s][gs]; ok {
		switch existingAction.actionType {
		case shift:
			panic(fmt.Sprintf("Shift-reduce conflict on state %v and input %v", s, gs))
		case reduce:
			panic(fmt.Sprintf("Reduce-reduce conflict on state %v and input %v", s, gs))
		case accept:
			panic(fmt.Sprintf("Accept-reduce conflict on state %v and input %v", s, gs))
		}
	}

	(*p)[s][gs] = parserAction{reduce, productionNumber}
}

func (p *parsingTable) addAcceptMove(s state) {
	if (*p)[s] == nil {
		(*p)[s] = make(map[grammarSymbol]parserAction)
	}

	(*p)[s]["$"] = parserAction{accept, 0}
}

type parser struct {
	g        grammar
	table    parsingTable
	stack    parserStack
	dead     bool
	accepted bool
}

func (ps *parser) init(t parsingTable, g grammar) {
	ps.g = g
	ps.table = t
	ps.stack = make(parserStack, 0, 10)
	ps.stack.push(0)
	ps.dead = false
	ps.accepted = false
}

func (ps *parser) move(input grammarSymbol) {
	if ps.dead {
		return
	}

	if _, ok := ps.table[ps.stack.top()]; !ok {
		ps.dead = true
		return
	}

	if _, ok := ps.table[ps.stack.top()][input]; !ok {
		ps.dead = true
		return
	}

	for ps.table[ps.stack.top()][input].actionType == reduce {
		prodNumber := ps.table[ps.stack.top()][input].number
		prod := ps.g.getProductionByNumber(prodNumber)
		for range prod.body {
			ps.stack.pop()
		}
		nextState := state(ps.table[ps.stack.top()][prod.head].number)
		ps.stack.push(nextState)
	}

	switch ps.table[ps.stack.top()][input].actionType {
	case accept:
		ps.accepted = true
	case shift:
		nextState := state(ps.table[ps.stack.top()][input].number)
		ps.stack.push(nextState)
	}
}
