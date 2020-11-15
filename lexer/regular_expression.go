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

func (r RegularExpression) getCharacters() []string {
	characters := make([]string, 0, 100)
	for i := 0; i < len(r); {
		if r[i] == '/' {
			characters = append(characters, string(r[i])+string(r[i+1]))
			i += 2
		} else {
			characters = append(characters, string(r[i]))
			i++
		}
	}
	return characters
}

func (r RegularExpression) isValid() bool {
	s := make(stack, 0, 10)
	for _, currentCharacter := range r.getCharacters() {
		switch currentCharacter {
		case "(":
			s.push('(')
		case ")":
			if s.empty() {
				return false
			}
			s.pop()
		}
	}

	return s.empty()
}

func (r RegularExpression) getMatchingParenIndex() int {
	characters := r.getCharacters()
	if characters[0] != "(" {
		return -1
	}

	s := make(stack, 0, 10)
	stringLengthSoFar := 0
	for _, character := range characters {
		stringLengthSoFar += len(character)
		switch character {
		case "(":
			s.push('(')
		case ")":
			s.pop()
		}
		if s.empty() {
			return stringLengthSoFar - 1
		}
	}
	return -1
}

func (r RegularExpression) trimParenthesis() RegularExpression {
	characters := r.getCharacters()
	if characters[0] != "(" {
		return r
	}
	if characters[len(characters)-1] != ")" {
		return r
	}
	return r[1 : len(r)-1]
}

func (r RegularExpression) getFirstOperand() RegularExpression {
	characters := r.getCharacters()
	if characters[0] != "(" {
		return RegularExpression(characters[0])
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
	characters := r.getCharacters()
	if len(characters) == 0 {
		var f nondeterministicFiniteAutomata
		f.init("")
		return f
	}

	if len(characters) == 1 {
		var f nondeterministicFiniteAutomata
		character := characters[0]
		switch len(character) {
		case 2:
			f.init(transitionLabel(character[1]))
		case 1:
			f.init(transitionLabel(character[0]))
		}
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
