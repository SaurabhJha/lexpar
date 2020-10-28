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
