package lexer

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func loadRegexFromDisk() map[string]regularExpression {
	regexTable := make(map[string]regularExpression)
	jsonContent, _ := ioutil.ReadFile("test.json")
	json.Unmarshal(jsonContent, &regexTable)
	return regexTable
}

func compileRegex(r map[string]regularExpression) map[string]deterministicFiniteAutomata {
	automataTable := make(map[string]deterministicFiniteAutomata)
	for label, regex := range r {
		nfa := regex.compile()
		dfa := nfa.convertToDfa()
		automataTable[label] = dfa
	}
	return automataTable
}

func readFromStdio() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(">> ")
	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(text)
	return text
}

// StartRegexInterpreter starts the top-level CLI regex interpreter
func StartRegexInterpreter() {
	regexTable := loadRegexFromDisk()
	automataTable := compileRegex(regexTable)
	for {
		fmt.Println(tokenize(automataTable, readFromStdio()))
	}
}
