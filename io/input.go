package io

import (
	"bufio"
	"fmt"
	"os"
)

// ReadFromStdin abstracts away all the handling of reading from STDIN.
func ReadFromStdin() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("> ")
	text, _ := reader.ReadString('\n')
	text = text[:len(text)-1]
	return text
}
