package lexer

import "testing"

func TestTokenize(t *testing.T) {
	regularExpressionTable := make(map[string]regularExpression)
	regularExpressionTable["id"] = regularExpression("(a|b|c)(a|b|c|0|1|2|3)*")
	regularExpressionTable["number"] = regularExpression("(0|1|2|3)(0|1|2|3)*")
	regularExpressionTable["+"] = regularExpression("+")
	regularExpressionTable["*"] = regularExpression("*")
	regularExpressionTable["("] = regularExpression("(")
	regularExpressionTable[")"] = regularExpression(")")
	regularExpressionTable["whitespace"] = regularExpression("( )( )*")

	automataTable := compileRegex(regularExpressionTable)

	var testData = []struct {
		input          string
		expectedOutput []token
	}{
		{
			"12 +  1231",
			[]token{
				{"12", "number"},
				{" ", "whitespace"},
				{"+", "+"},
				{"  ", "whitespace"},
				{"1231", "number"},
			},
		},
		{
			"(abca+12)*a",
			[]token{
				{"(", "("},
				{"abca", "id"},
				{"+", "+"},
				{"12", "number"},
				{")", ")"},
				{"*", "*"},
				{"a", "id"},
			},
		},
		{
			"13 + 123 * 1233",
			[]token{
				{"13", "number"},
				{" ", "whitespace"},
				{"+", "+"},
				{" ", "whitespace"},
				{"123", "number"},
				{" ", "whitespace"},
				{"*", "*"},
				{" ", "whitespace"},
				{"1233", "number"},
			},
		},
	}

	for _, test := range testData {
		actualOutput := tokenize(automataTable, test.input)
		if len(actualOutput) != len(test.expectedOutput) {
			t.Fatalf("For input %v, expected %v but got %v", test.input, test.expectedOutput, actualOutput)
		}

		for i, expectedToken := range test.expectedOutput {
			actualToken := actualOutput[i]
			if actualToken.lexeme != expectedToken.lexeme || actualToken.tokenType != expectedToken.tokenType {
				t.Fatalf("For input %v, expected %v but got %v", test.input, test.expectedOutput, actualOutput)
			}
		}
	}
}
