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
