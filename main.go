package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func main() {
	regexTable := make(map[string]regularExpression)
	jsonContent, _ := ioutil.ReadFile("test.json")
	json.Unmarshal(jsonContent, &regexTable)
	automataTable := make(map[string]deterministicFiniteAutomata)
	for label, regex := range regexTable {
		nfa := regex.compile()
		dfa := nfa.convertToDfa()
		automataTable[label] = dfa
	}
	for {
		fmt.Println(tokenize(automataTable, readFromStdio()))
	}
}
