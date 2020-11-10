package parser

import "github.com/SaurabhJha/lexpar/lexer"

// Parser is the data structure used to export all the functionality that can be expected
// from an LR parser
type Parser struct {
	p parser
}

// Init of Parser sets up all the state required by the parser to start processing terminals.
func (P *Parser) Init(g Grammar) {
	P.p = g.compile()
}

// Parse takes as input a slice of tokens and parses them.
func (P *Parser) Parse(tokens []lexer.Token) []Production {
	tokenTypes := make([]grammarSymbol, 0, 100)
	for _, token := range tokens {
		tokenTypes = append(tokenTypes, grammarSymbol(token.TokenType()))
	}
	tokenTypes = append(tokenTypes, "$")

	reductions := make([]Production, 0, 100)
	P.p.parse(tokenTypes)
	reductions = P.p.reductions
	return reductions
}

// Reset resets parser state back to its initial state where it can parse more tokens.
func (P *Parser) Reset() {
	P.p.reset()
}
