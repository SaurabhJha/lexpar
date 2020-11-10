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

// RegularExpression represents the string representation of regular expressions. It has methods for
// regular expression operations and compilation.
type RegularExpression string

func (r RegularExpression) isValid() bool {
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

func (r RegularExpression) getMatchingParenIndex() int {
	if r[0] != '(' {
		return -1
	}

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

func (r RegularExpression) trimParenthesis() RegularExpression {
	if r[0] != '(' {
		return r
	}
	if r[len(r)-1] != ')' {
		return r
	}
	return r[1 : len(r)-1]
}

func (r RegularExpression) getFirstOperand() RegularExpression {
	if r[0] != '(' {
		return RegularExpression(r[0])
	}

	return RegularExpression(r[:r.getMatchingParenIndex()+1])
}

func (r RegularExpression) getOperator() regularExpressionOperator {
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

func (r RegularExpression) getSecondOperand() RegularExpression {
	secondOperandIndex := len(r.getFirstOperand()) + r.getOperator().length()
	return RegularExpression(r[secondOperandIndex:])
}

func (r RegularExpression) compile() nondeterministicFiniteAutomata {
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
