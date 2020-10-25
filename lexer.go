package main

import "fmt"

type token struct {
	tokenType string
	lexeme    string
}

func matchPrefix(d deterministicFiniteAutomata, s string) string {
	matchingPrefixIndex := -1
	for i, character := range s {
		d.move(transitionLabel(character))
		if d.accepted {
			matchingPrefixIndex = i
		}
		if d.dead {
			break
		}
	}

	if matchingPrefixIndex == -1 {
		return ""
	}
	return s[:matchingPrefixIndex+1]
}

func getMaxMatchedPrefix(table map[string]deterministicFiniteAutomata, s string) (string, string) {
	var maxMatchedPrefix string
	var maxLabel string
	for l, a := range table {
		matched := matchPrefix(a, s)
		if len(matched) > len(maxMatchedPrefix) {
			maxMatchedPrefix = matched
			maxLabel = l
		}
	}
	return maxMatchedPrefix, maxLabel
}

func tokenize(table map[string]deterministicFiniteAutomata, s string) []token {
	tokens := make([]token, 0, 100)
	for len(s) != 0 {
		prefix, label := getMaxMatchedPrefix(table, s)
		if len(prefix) == 0 {
			fmt.Println(fmt.Errorf("Problem"))
			break
		}
		s = s[len(prefix):]
		tokens = append(tokens, token{prefix, label})
	}

	return tokens
}
