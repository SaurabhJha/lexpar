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
