package parser

import "reflect"

type grammarSymbol string

// SemanticRule is a syntax directed definition (SDD) associated with a grammar.
type SemanticRule struct {
	Type      string
	RootLabel string
	Children  []int
}

func (rule SemanticRule) isEmpty() bool {
	return reflect.DeepEqual(
		rule,
		SemanticRule{
			"",
			"",
			nil,
		},
	)
}

// Production is a grammar production in Backus-Naur form.
type Production struct {
	Head grammarSymbol
	Body []grammarSymbol
	Rule SemanticRule
}

func (p Production) getFirstBodySymbol() grammarSymbol {
	if len(p.Body) == 0 {
		return ""
	}

	return p.Body[0]
}

// Grammar is a context-free grammar which is a list of productions and a start symbol.
type Grammar struct {
	Productions []Production
	Start       grammarSymbol
}

func (g Grammar) isTerminal(s grammarSymbol) bool {
	for _, production := range g.Productions {
		if s == production.Head {
			return false
		}
	}
	return true
}

func (g Grammar) getProductionsOfSymbol(s grammarSymbol) []Production {
	productions := make([]Production, 0, 10)
	for _, p := range g.Productions {
		if p.Head == s {
			productions = append(productions, p)
		}
	}

	return productions
}

func (g Grammar) computeFirstSet(s grammarSymbol) setOfSymbols {
	if g.isTerminal(s) {
		return setOfSymbols{s: true}
	}

	firstSet := make(setOfSymbols)
	for _, production := range g.getProductionsOfSymbol(s) {
		firstBodySymbol := production.getFirstBodySymbol()
		if firstBodySymbol != "" && firstBodySymbol != production.Head {
			firstSetOfBody := g.computeFirstSet(firstBodySymbol)
			firstSet.unionWith(&firstSetOfBody)
		}
	}

	return firstSet
}

func (g Grammar) computeFollowSet(s grammarSymbol) setOfSymbols {
	if g.isTerminal(s) {
		return setOfSymbols{}
	}

	followSet := setOfSymbols{}
	if s == g.Start {
		followSet.add("$")
	}
	for _, p := range g.Productions {
		for i, bodySymbol := range p.Body {
			if bodySymbol == s {
				if i == len(p.Body)-1 {
					followSetOfHead := g.computeFollowSet(p.Head)
					followSet.unionWith(&followSetOfHead)
				} else {
					firstSetOfNextSym := g.computeFirstSet(p.Body[i+1])
					followSet.unionWith(&firstSetOfNextSym)
				}
			}
		}
	}
	return followSet
}

func (g Grammar) getProductionNumber(p Production) int {
	for i := range g.Productions {
		if reflect.DeepEqual(g.Productions[i], p) {
			return i
		}
	}
	return -1
}

func (g Grammar) compile() parser {
	table := make(parsingTable)
	startProduction := g.getProductionsOfSymbol(g.Start)[0]
	startItem := lrItem{g, startProduction, 0, map[grammarSymbol]bool{"$": true}}
	startItemSet := startItem.computeClosureSet()
	q := make(queueOfItemSets, 0, 10)
	q.enqueue(startItemSet)
	seen := make(seenLrItemSets, 0, 100)
	seen.add(startItemSet)

	for !q.empty() {
		currentItemSet := q.dequeue()
		currentState := seen.getStateNumber(currentItemSet)

		// Add shift moves.
		for symbol := range currentItemSet.getNextSymbols() {
			nextItemSet := currentItemSet.getNextItemSet(symbol)
			if !seen.has(nextItemSet) {
				q.enqueue(nextItemSet)
				seen.add(nextItemSet)
			}
			nextState := seen.getStateNumber(nextItemSet)
			table.addShiftMove(currentState, nextState, symbol)
		}

		// Add reduce and accept moves.
		for _, item := range currentItemSet.itemSet {
			if item.getNextSymbol() == "" {
				productionNumber := item.g.getProductionNumber(item.p)
				for symbol := range item.followSet {
					if productionNumber == 0 && symbol == "$" {
						table.addAcceptMove(currentState)
					} else {
						table.addReduceMove(currentState, productionNumber, symbol)
					}
				}
			}
		}
	}

	var ps parser
	ps.init(table, g)
	return ps
}
