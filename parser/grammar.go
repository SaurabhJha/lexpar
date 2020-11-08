package parser

import "reflect"

type grammarSymbol string

type production struct {
	head grammarSymbol
	body []grammarSymbol
}

func (p production) getFirstBodySymbol() grammarSymbol {
	if len(p.body) == 0 {
		return ""
	}

	return p.body[0]
}

type grammar struct {
	productions []production
	start       grammarSymbol
}

func (g grammar) isTerminal(s grammarSymbol) bool {
	for _, production := range g.productions {
		if s == production.head {
			return false
		}
	}
	return true
}

func (g grammar) getProductionsOfSymbol(s grammarSymbol) []production {
	productions := make([]production, 0, 10)
	for _, p := range g.productions {
		if p.head == s {
			productions = append(productions, p)
		}
	}

	return productions
}

func (g grammar) computeFirstSet(s grammarSymbol) setOfSymbols {
	if g.isTerminal(s) {
		return setOfSymbols{s: true}
	}

	firstSet := make(setOfSymbols)
	for _, production := range g.getProductionsOfSymbol(s) {
		firstBodySymbol := production.getFirstBodySymbol()
		if firstBodySymbol != "" && firstBodySymbol != production.head {
			firstSetOfBody := g.computeFirstSet(firstBodySymbol)
			firstSet.unionWith(&firstSetOfBody)
		}
	}

	return firstSet
}

func (g grammar) computeFollowSet(s grammarSymbol) setOfSymbols {
	if g.isTerminal(s) {
		return setOfSymbols{}
	}

	followSet := setOfSymbols{}
	if s == g.start {
		followSet.add("$")
	}
	for _, p := range g.productions {
		for i, bodySymbol := range p.body {
			if bodySymbol == s {
				if i == len(p.body)-1 {
					followSetOfHead := g.computeFollowSet(p.head)
					followSet.unionWith(&followSetOfHead)
				} else {
					firstSetOfNextSym := g.computeFirstSet(p.body[i+1])
					followSet.unionWith(&firstSetOfNextSym)
				}
			}
		}
	}
	return followSet
}

func (g grammar) getProductionNumber(p production) int {
	for i := range g.productions {
		if reflect.DeepEqual(g.productions[i], p) {
			return i
		}
	}
	return -1
}

func (g grammar) compile() parser {
	table := make(parsingTable)
	startProduction := g.getProductionsOfSymbol(g.start)[0]
	startItem := lrItem{g, startProduction, 0}
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
				for symbol := range item.g.computeFollowSet(item.p.head) {
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
