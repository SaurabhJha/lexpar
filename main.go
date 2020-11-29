package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

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
		commandType := io.GetCommandType(text)
		switch commandType {
		case "quit":
			os.Exit(0)
		case "setRegex":
			io.ExecuteRegexCommand(text, &definitions)
		case "persist":
			io.Persist(&definitions)
		case "print":
			io.Print(&definitions)
		default:
			tokens := tok.Tokenize(text)
			tree := pars.Parse(tokens)
			fmt.Println(tree)
			tok.Reset()
			pars.Reset()
		}
	}
}
