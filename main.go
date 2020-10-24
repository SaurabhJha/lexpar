package main

import "fmt"

func main() {
	for {
		regex := readRegexFromStdio()
		fmt.Println(regex.compile())
	}
}
