package lexer

import (
	"reflect"
	"testing"
)

func TestTokenizerMatchingPrefix(t *testing.T) {
	regexTable := map[string]RegularExpression{
		"id":     "(a|b|c)(a|b|c|0|1|2)*",
		"number": "(1|2)(0|1|2|3|4)*",
		"+":      "+",
	}
	var tokenizer Tokenizer
	tokenizer.Init(regexTable)

	testData := []struct {
		id       string
		input    string
		expected string
	}{
		{"id", "abc+12", "abc"},
		{"number", "abc+12", ""},
		{"id", "123+abc", ""},
		{"number", "123+abc", "123"},
	}

	for _, test := range testData {
		if got := tokenizer.getMatchingPrefix(test.id, test.input); got != test.expected {
			t.Errorf("Matching prefix expected %v, got %v", test.expected, got)
		}
	}
}

func TestTokenizerMaxMatchingPrefix(t *testing.T) {
	regexTable := map[string]RegularExpression{
		"id": "(a|b|c)(a|b|c|0|1|2)*",
		"=":  "=",
		"==": "==",
	}
	var tokenizer Tokenizer
	tokenizer.Init(regexTable)

	testData := []struct {
		input          string
		expectedID     string
		expectedLexeme string
	}{
		{"abc121", "id", "abc121"},
		{"abc+", "id", "abc"},
		{"==123", "==", "=="},
	}

	for _, test := range testData {
		if gotID, gotLexeme := tokenizer.getMaxMatchingPrefix(test.input); gotID != test.expectedID || gotLexeme != test.expectedLexeme {
			t.Errorf("Max matching prefix expected %v %v, got %v %v",
				test.expectedID, test.expectedLexeme, gotID, gotLexeme)
		}
	}
}

func TestTokenizerTokenize(t *testing.T) {
	regexTable := map[string]RegularExpression{
		"id":     "(a|b|c)(a|b|c|0|1|2)*",
		"+":      "+",
		"=":      "=",
		"==":     "==",
		"number": "(1|2|3)(0|1|2|3)*",
	}
	var tokenizer Tokenizer
	tokenizer.Init(regexTable)

	testData := []struct {
		input    string
		expected []Token
	}{
		{
			"123+23",
			[]Token{
				{"number", "123"},
				{"+", "+"},
				{"number", "23"},
			},
		},
		{
			"abc==123",
			[]Token{
				{"id", "abc"},
				{"==", "=="},
				{"number", "123"},
			},
		},
		{

			"**123",
			[]Token{},
		},
	}

	for _, test := range testData {
		if got := tokenizer.Tokenize(test.input); !reflect.DeepEqual(got, test.expected) {
			t.Errorf("Tokenization on input %v expected %v, got %v", test.input, test.expected, got)
		}
	}
}
