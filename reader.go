package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func readRegexFromStdio() regularExpression {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(">> ")
	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(text)
	return regularExpression(text)
}
