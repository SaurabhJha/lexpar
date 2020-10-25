package lexer

type regularExpressionOperator int

const (
	union regularExpressionOperator = iota
	concat
	star
)

func (r regularExpressionOperator) length() int {
	switch r {
	case union:
		return 1
	case star:
		return 1
	default:
		return 0
	}
}

type regularExpression string

func (r regularExpression) isValid() bool {
	s := make(stack, 0, 10)
	for _, character := range r {
		switch character {
		case '(':
			s.push('(')
		case ')':
			if s.empty() {
				return false
			}
			s.pop()
		}
	}

	return s.empty()
}

func (r regularExpression) getMatchingParenIndex() int {
	s := make(stack, 0, 10)
	for i, character := range r {
		switch character {
		case '(':
			s.push('(')
		case ')':
			s.pop()
		}
		if s.empty() {
			return i
		}
	}
	return -1
}

func (r regularExpression) trimParenthesis() regularExpression {
	if r[0] != '(' {
		return r
	}
	if r[len(r)-1] != ')' {
		return r
	}
	return r[1 : len(r)-1]
}

func (r regularExpression) getFirstOperand() regularExpression {
	if r[0] != '(' {
		return regularExpression(r[0])
	}

	return regularExpression(r[:r.getMatchingParenIndex()+1])
}

func (r regularExpression) getOperator() regularExpressionOperator {
	operatorIndex := len(r.getFirstOperand())
	switch r[operatorIndex] {
	case '|':
		return union
	case '*':
		return star
	default:
		return concat
	}
}

func (r regularExpression) getSecondOperand() regularExpression {
	secondOperandIndex := len(r.getFirstOperand()) + r.getOperator().length()
	return regularExpression(r[secondOperandIndex:])
}

func (r regularExpression) compile() nondeterministicFiniteAutomata {
	if len(r) == 0 {
		var f nondeterministicFiniteAutomata
		f.init("")
		return f
	}

	if len(r) == 1 {
		var f nondeterministicFiniteAutomata
		f.init(transitionLabel(r[0]))
		return f
	}

	firstOperand := r.getFirstOperand()
	switch r.getOperator() {
	case star:
		firstOperandAutomata := firstOperand.trimParenthesis().compile()
		firstOperandAutomata.applyStar()
		return firstOperandAutomata
	case union:
		firstOperandAutomata := firstOperand.trimParenthesis().compile()
		secondOperand := r.getSecondOperand()
		secondOperandAutomata := secondOperand.trimParenthesis().compile()
		firstOperandAutomata.combineUsingUnion(&secondOperandAutomata)
		return firstOperandAutomata
	default:
		firstOperandAutomata := firstOperand.trimParenthesis().compile()
		secondOperand := r.getSecondOperand()
		secondOperandAutomata := secondOperand.trimParenthesis().compile()
		firstOperandAutomata.combineUsingConcat(&secondOperandAutomata)
		return firstOperandAutomata
	}
}
