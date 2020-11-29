package io

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/SaurabhJha/lexpar/lexer"
)

// ReadFromStdin abstracts away all the handling of reading from STDIN.
func ReadFromStdin() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("> ")
	text, _ := reader.ReadString('\n')
	text = text[:len(text)-1]
	return text
}

// GetCommandType parses a command and gets its command type.
func GetCommandType(input string) string {
	inputSlice := strings.Split(input, " ")
	switch inputSlice[0] {
	case "quit":
		return "quit"
	case "setRegex":
		return "setRegex"
	case "addProduction":
		return "addProduction"
	case "setStartSymbol":
		return "setStartSymbol"
	case "persist":
		return "persist"
	case "print":
		return "print"
	default:
		return "eval"
	}
}

// ExecuteRegexCommand assumes that the command type is "setRegex" and returns token type - regular expression
// pair
func ExecuteRegexCommand(command string, definitions *DefinitionsTable) {
	commandSlice := strings.SplitN(command, " ", 3)
	regexType, regex := commandSlice[1], commandSlice[2]
	definitions.RegularExpressions[regexType] = lexer.RegularExpression(regex)
}

// Persist takes current defintiions and persist it to disk, overwriting current contents.
func Persist(definitions *DefinitionsTable) {
	definitionsJSON, _ := json.MarshalIndent(*definitions, "", "	")
	file, err := os.OpenFile("example.json", os.O_WRONLY, 0755)
	fmt.Println(definitionsJSON, err)
	_, err1 := file.Write(definitionsJSON)
	fmt.Println(err1)
	file.Close()
}

// Print just prints out the definitions data structure
func Print(definitions *DefinitionsTable) {
	fmt.Println("Regular expressions")
	for regexName, regex := range definitions.RegularExpressions {
		fmt.Printf("  %s: %s\n", regexName, regex)
	}
	fmt.Println("Grammar")
	fmt.Println("  Start symbol: ", definitions.Grammar.Start)
	for _, production := range definitions.Grammar.Productions {
		fmt.Printf("  %s -> %s\n", production.Head, production.Body)
	}
}
