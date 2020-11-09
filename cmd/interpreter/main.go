package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/SaurabhJha/lexpar/io"
	"github.com/SaurabhJha/lexpar/lexer"
	"github.com/SaurabhJha/lexpar/parser"
)

func main() {
	jsonContent, _ := ioutil.ReadFile("../../example.json")
	var definitions io.DefinitionsTable
	json.Unmarshal(jsonContent, &definitions)

	for {
		var tokenizer lexer.Tokenizer
		tokenizer.Init(definitions.RegularExpressions)
		var lrParser parser.Parser
		lrParser.Init(definitions.Grammar)
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("> ")
		text, _ := reader.ReadString('\n')
		text = text[:len(text)-1]
		fmt.Printf("Input: %v\n", text)
		// text := "123+345"
		tokens := tokenizer.Tokenize(text)
		fmt.Printf("Tokens: %v\n", tokens)
		reductions := lrParser.Parse(tokens)
		fmt.Printf("Reductions: %v\n", reductions)
	}
}
