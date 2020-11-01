package parser

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

func (p production) equals(p2 production) bool {
	if p.head != p2.head {
		return false
	}

	if len(p.body) != len(p2.body) {
		return false
	}

	for i, s1 := range p.body {
		s2 := p2.body[i]
		if s1 != s2 {
			return false
		}
	}

	return true
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

func (g grammar) computeFirstSet(s grammarSymbol) setOfSymbols {
	if g.isTerminal(s) == true {
		firstSet := setOfSymbols{s: true}
		return firstSet
	}

	firstSet := make(setOfSymbols)
	for _, production := range g.productions {
		if production.head == s && production.getFirstBodySymbol() != "" && production.getFirstBodySymbol() != production.head {
			firstSetOfBody := g.computeFirstSet(production.getFirstBodySymbol())
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
		for i, sym := range p.body {
			if sym == s {
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

func (g grammar) getProductionsOfSymbol(s grammarSymbol) []production {
	productions := make([]production, 0, 10)
	for _, p := range g.productions {
		if p.head == s {
			productions = append(productions, p)
		}
	}

	return productions
}

func (g grammar) getProductionNumber(p production) int {
	for i := range g.productions {
		if g.productions[i].equals(p) {
			return i
		}
	}
	return -1
}

func (g grammar) equals(g2 grammar) bool {
	if g.start != g2.start {
		return false
	}

	if len(g.productions) != len(g2.productions) {
		return false
	}

	for i, p1 := range g.productions {
		p2 := g2.productions[i]
		if !p1.equals(p2) {
			return false
		}
	}

	return true
}

func (g grammar) compile() parsingTable {
	table := make(parsingTable)
	startProduction := g.getProductionsOfSymbol(g.start)[0]
	startItem := lrItem{g, startProduction, 0}
	startItemSet := startItem.computeClosureSet()
	q := make(queueOfItemSets, 0, 10)
	q.enqueue(startItemSet)
	seen := make(seenLrItemSets, 0, 100)

	for !q.empty() {
		currentItemSet := q.dequeue()

		// Add shift moves.
		for symbol := range currentItemSet.getNextSymbols() {
			nextItemSet := currentItemSet.getNextItemSet(symbol)
			if !seen.has(nextItemSet) {
				seen.add(nextItemSet)
				q.enqueue(nextItemSet)
			}
			table.addShiftMove(seen.getStateNumber(currentItemSet), seen.getStateNumber(nextItemSet), symbol)
		}

		// Add reduce moves.
		for _, item := range currentItemSet.itemSet {
			if item.getNextSymbol() == "" {
				productionNumber := item.g.getProductionNumber(item.p)
				for symbol := range item.g.computeFollowSet(item.p.head) {
					table.addReduceMove(seen.getStateNumber(currentItemSet), productionNumber, symbol)
				}
			}
		}
	}

	return table
}
