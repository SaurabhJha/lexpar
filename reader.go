package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func readFromStdio() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(">> ")
	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(text)
	return text
}
