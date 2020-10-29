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

type lrItemSet struct {
	itemSet []lrItem
}

func (ls *lrItemSet) add(l lrItem) {
	(*ls).itemSet = append((*ls).itemSet, l)
}

func (ls *lrItemSet) has(l lrItem) bool {
	for _, item := range ls.itemSet {
		if item.equals(l) {
			return true
		}
	}
	return false
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
