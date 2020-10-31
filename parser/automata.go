package parser

type parserActionType int

const (
	shift parserActionType = iota
	reduce
)

type parserAction struct {
	actionType parserActionType
	number     int
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

func (ls *lrItemSet) getOutgoingSymbols() setOfSymbols {
	outgoingSymbols := make(setOfSymbols)
	for _, l := range ls.itemSet {
		if l.getNextSymbol() != "" {
			outgoingSymbols.add(l.getNextSymbol())
		}
	}
	return outgoingSymbols
}
