package main

import "fmt"

func main() {
	for {
		regex := readRegexFromStdio()
		automata := regex.compile()
		fmt.Println(automata.convertToDfa())
	}
}
