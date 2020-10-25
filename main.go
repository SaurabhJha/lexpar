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
		i := readFromStdio()
		var maxMatchedPrefix string
		var maxLabel string
		for l, a := range automataTable {
			matched := matchPrefix(a, i)
			if len(matched) > len(maxMatchedPrefix) {
				maxMatchedPrefix = matched
				maxLabel = l
			}
		}
		fmt.Println(maxLabel, maxMatchedPrefix)
	}
}
