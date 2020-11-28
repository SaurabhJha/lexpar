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
	text := "12+123*4"

	tokens := tok.Tokenize(text)
	tree := pars.Parse(tokens)
	fmt.Println(tree)
	tok.Reset()
	// for {
	// 	text := io.ReadFromStdin()
	// 	tokens := tok.Tokenize(text)
	// 	tree := pars.Parse(tokens)
	// 	fmt.Println(tree)
	// 	tok.Reset()
	// 	pars.Reset()
	// }
}
