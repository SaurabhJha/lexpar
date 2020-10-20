package main

import (
	"fmt"
)

func main() {
	var automata finiteAutomata
	var otherAutomata finiteAutomata
	automata.initialize("a")
	otherAutomata.initialize("b")
	automata.combineUsingConcat(&otherAutomata)
	automata.applyStar()
	fmt.Printf("%+v\n", automata)
}
