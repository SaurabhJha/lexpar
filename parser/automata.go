package parser

import (
	"fmt"
	"reflect"

	"github.com/SaurabhJha/lexpar/lexer"
)

type lrItem struct {
	g   Grammar
	p   Production
	pos int
}

func (l lrItem) getNextSymbol() grammarSymbol {
	if l.pos == len(l.p.Body) {
		return ""
	}
	return l.p.Body[l.pos]
}

func (l lrItem) getNextItem(s grammarSymbol) lrItem {
	if l.pos >= len(l.p.Body) {
		return lrItem{}
	}
	if l.p.Body[l.pos] != s {
		return lrItem{}
	}
	return lrItem{l.g, l.p, l.pos + 1}
}

func (l lrItem) empty() bool {
	return reflect.DeepEqual(
		l, lrItem{
			Grammar{nil, ""},
			Production{"", nil, SemanticRule{"", "", nil}},
			0,
		},
	)
}

type lrItemSet struct {
	itemSet []lrItem
}

func (ls *lrItemSet) has(l lrItem) bool {
	for _, item := range ls.itemSet {
		if reflect.DeepEqual(item, l) {
			return true
		}
	}
	return false
}

func (ls *lrItemSet) equals(ols *lrItemSet) bool {
	for _, item := range ls.itemSet {
		if !ols.has(item) {
			return false
		}
	}

	for _, item := range ols.itemSet {
		if !ls.has(item) {
			return false
		}
	}

	return true
}

func (ls *lrItemSet) add(l lrItem) {
	if !ls.has(l) {
		(*ls).itemSet = append((*ls).itemSet, l)
	}
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
			// We have not encountered this non terminal before. Add its productions to the item set
			// and to the queue. We add them to queue so that items reachable from those can be added
			// later.
			for _, p := range l.g.getProductionsOfSymbol(nextSymbol) {
				nextItem := lrItem{l.g, p, 0}
				if !ls.has(nextItem) {
					q.enqueue(nextItem)
					ls.add(nextItem)
				}
			}
			seenNonTerminals.add(nextSymbol)
		}
	}

	return ls
}

func (ls *lrItemSet) getNextItemSet(s grammarSymbol) lrItemSet {
	var nextLs lrItemSet
	for _, l := range ls.itemSet {
		if nextItem := l.getNextItem(s); !nextItem.empty() {
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

	if _, ok := (*p)[s]["$"]; ok {
		panic(fmt.Sprintf("Cannot add accept move, already an action on %v and %v", s, "$"))
	}

	(*p)[s]["$"] = parserAction{accept, 0}
}

type parser struct {
	g        Grammar
	table    parsingTable
	pStack   parserStack
	dead     bool
	accepted bool
	ast      SyntaxGraph
	gStack   graphStack
}

func (ps *parser) init(t parsingTable, g Grammar) {
	stack := make(parserStack, 0, 10)
	stack.push(0)
	*ps = parser{g, t, stack, false, false, SyntaxGraph{}, graphStack{}}
}

func (ps *parser) move(token lexer.Token) {
	tokenType := grammarSymbol(token.TokenType)
	if ps.dead {
		return
	}

	if _, ok := ps.table[ps.pStack.top()]; !ok {
		ps.dead = true
		return
	}

	if _, ok := ps.table[ps.pStack.top()][tokenType]; !ok {
		ps.dead = true
		return
	}

	// Make as many reduce moves as possible on current input symbol and execute production SDD on each
	// reduction.
	for ps.table[ps.pStack.top()][tokenType].actionType == reduce {
		prodNumber := ps.table[ps.pStack.top()][tokenType].number
		prod := ps.g.Productions[prodNumber]

		// SLR reduction
		for range prod.Body {
			ps.pStack.pop()
		}
		nextParserAction := ps.table[ps.pStack.top()][prod.Head]
		nextState := state(nextParserAction.number)
		ps.pStack.push(nextState)

		// SDD execution
		rule := prod.Rule
		if !rule.isEmpty() {
			stackContents := make([]int, 0, 5)
			for range prod.Body {
				stackContents = append(stackContents, ps.gStack.pop())
			}
			switch rule.Type {
			case "tree":
				rootNodeIndex := ps.ast.createNewNode(rule.RootLabel)
				for childIdx := len(rule.Children) - 1; childIdx >= 0; childIdx-- {
					childNodeIndex := stackContents[rule.Children[childIdx]]
					ps.ast.addEdge(rootNodeIndex, childNodeIndex)
				}
				ps.gStack.push(rootNodeIndex)
			case "copy":
				ps.gStack.push(stackContents[rule.Children[0]])
			}
		}
	}

	switch nextParserAction := ps.table[ps.pStack.top()][tokenType]; nextParserAction.actionType {
	case accept:
		ps.accepted = true
	case shift:
		nextState := state(ps.table[ps.pStack.top()][tokenType].number)
		ps.pStack.push(nextState)
		newNode := ps.ast.createNewNode(token.Lexeme)
		ps.gStack.push(newNode)
	}
}

func (ps *parser) parse(tokens []lexer.Token) SyntaxGraph {
	for _, token := range tokens {
		ps.move(token)
	}
	return ps.ast
}

func (ps *parser) reset() {
	ps.pStack = make(parserStack, 0, 10)
	ps.pStack.push(0)
	ps.dead = false
	ps.accepted = false
	ps.ast = SyntaxGraph{}
}
