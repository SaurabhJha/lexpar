package lexer

// Tokenizer is the data structure used to export all the functionality that can
// be expected from a tokenizer or lexer.
type Tokenizer struct {
	automata map[string]deterministicFiniteAutomata
}

// Init sets up all the state required for Tokenizer to start processing strings.
func (t *Tokenizer) Init(regexJSON map[string]string) {
	t.automata = make(map[string]deterministicFiniteAutomata)
	for regexID, regexString := range regexJSON {
		regex := regularExpression(regexString)
		nfa := regex.compile()
		dfa := nfa.convertToDfa()
		t.automata[regexID] = dfa
	}
}

func (t *Tokenizer) getMatchingPrefix(regexID string, input string) string {
	dfa := t.automata[regexID]
	acceptedAt := -1
	for pos, character := range input {
		label := transitionLabel(string(character))
		dfa.move(label)
		if dfa.dead {
			break
		}
		if dfa.accepted {
			acceptedAt = pos
		}
	}
	dfa.reset()
	return input[:acceptedAt+1]
}

func (t *Tokenizer) getMaxMatchingPrefix(input string) (string, string) {
	var maxPrefix string
	var maxRegexID string
	for id := range t.automata {
		prefix := t.getMatchingPrefix(id, input)
		if len(prefix) > len(maxPrefix) {
			maxPrefix = prefix
			maxRegexID = id
		}
	}
	return maxRegexID, maxPrefix
}

// Token represents the output of the lexer.
type Token struct {
	tokenType string
	lexeme    string
}

// TokenType reads off the type of the token
func (token *Token) TokenType() string {
	return token.tokenType
}

// Lexeme reads off the matching string of the token
func (token *Token) Lexeme() string {
	return token.lexeme
}

// Tokenize returns an array of tokens given an input string.
func (t *Tokenizer) Tokenize(input string) []Token {
	tokens := make([]Token, 0, 100)

	remainingInput := input

	for len(remainingInput) != 0 {
		nextTokenType, nextLexeme := t.getMaxMatchingPrefix(remainingInput)
		if len(nextLexeme) == 0 {
			break
		}
		tokens = append(tokens, Token{nextTokenType, nextLexeme})
		remainingInput = remainingInput[len(nextLexeme):]
	}

	return tokens
}
