package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/SaurabhJha/lexpar/io"
	"github.com/SaurabhJha/lexpar/lexer"
	"github.com/SaurabhJha/lexpar/parser"
)

func main() {
	jsonContent, _ := ioutil.ReadFile("example.json")
	var definitions io.DefinitionsTable
	json.Unmarshal(jsonContent, &definitions)

	var tok lexer.Tokenizer
	tok.Init(definitions.RegularExpressions)

	var pars parser.Parser
	pars.Init(definitions.Grammar)
	for {
		text := io.ReadFromStdin()
		fmt.Printf("Input: %v\n", text)
		tokens := tok.Tokenize(text)
		fmt.Printf("Tokens: %v\n", tokens)
		reductions := pars.Parse(tokens)
		fmt.Printf("Reductions: %v\n", reductions)

		tok.Reset()
		pars.Reset()
	}
}
