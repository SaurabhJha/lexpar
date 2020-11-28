package lexer

import "fmt"

// Token represents one "word" of a program text. A sequence of tokens are output by a tokenizer.
type Token struct {
	TokenType string
	Lexeme    string
}

// A Tokenizer object breaks up strings using a collection of regular expressions.
type Tokenizer struct {
	automata map[string]deterministicFiniteAutomata
}

// Init sets up all the state required for Tokenizer to start processing strings.
func (t *Tokenizer) Init(regexJSON map[string]RegularExpression) {
	t.automata = make(map[string]deterministicFiniteAutomata)
	for regexID, regex := range regexJSON {
		if !regex.isValid() {
			panic(fmt.Sprintf("Regex '%v' is invalid, aborting", regex))
		}
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

// Reset method of tokenizer resets the tokenizer back to its initial state so that it can parse new
// strings.
func (t *Tokenizer) Reset() {
	for _, automata := range t.automata {
		automata.reset()
	}
}
