package parser

type setOfSymbols map[grammarSymbol]bool

func (ss *setOfSymbols) add(s grammarSymbol) {
	(*ss)[s] = true
}

func (ss *setOfSymbols) unionWith(os *setOfSymbols) {
	for s := range *os {
		(*ss).add(s)
	}
}

func (ss *setOfSymbols) hasSubset(os *setOfSymbols) bool {
	for s := range *os {
		if _, ok := (*ss)[s]; ok == false {
			return false
		}
	}
	return true
}

func (ss *setOfSymbols) isEqualTo(os *setOfSymbols) bool {
	return ss.hasSubset(os) && os.hasSubset(ss)
}

func (ss *setOfSymbols) has(s grammarSymbol) bool {
	return (*ss)[s]
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
		if sls.equals(&ls) {
			return true
		}
	}
	return false
}

func (lss *seenLrItemSets) getStateNumber(ls lrItemSet) state {
	for i, sls := range *lss {
		if sls.equals(&ls) {
			return state(i)
		}
	}
	lss.add(ls)
	return state(len(*lss) - 1)
}
