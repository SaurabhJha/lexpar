package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/SaurabhJha/lexpar/ast"
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

	var sdt ast.SyntaxDirectedTranslator
	sdt.Init(definitions.Grammar, definitions.SemanticRules)
	for {
		text := io.ReadFromStdin()
		tokens := tok.Tokenize(text)
		reductions := pars.Parse(tokens)
		ast := sdt.ConstructAST(tokens, reductions)
		fmt.Println(ast)
		tok.Reset()
		pars.Reset()
	}
}
