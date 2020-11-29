package parser

import (
	"reflect"
	"testing"

	"github.com/SaurabhJha/lexpar/lexer"
)

func TestGetFirstBody(t *testing.T) {
	var testData = []struct {
		p    Production
		want grammarSymbol
	}{
		{Production{"expr", []grammarSymbol{"expr", "+", "term"}, SemanticRule{"", "", []int{}}}, "expr"},
		{Production{"term", []grammarSymbol{"factor"}, SemanticRule{"", "", []int{}}}, "factor"},
		{Production{"expr'", []grammarSymbol{}, SemanticRule{"", "", []int{}}}, ""},
	}

	for _, test := range testData {
		if got := test.p.getFirstBodySymbol(); got != test.want {
			t.Errorf("p.getFirstBodySymbol() = %v, expected %v", got, test.want)
		}
	}
}

func TestGetProductionsOfSymbol(t *testing.T) {
	var g Grammar
	g.Start = "expr'"
	g.Productions = []Production{
		{"expr'", []grammarSymbol{"expr"}, SemanticRule{"", "", []int{}}},
		{"expr", []grammarSymbol{"expr", "+", "term"}, SemanticRule{"", "", []int{}}},
		{"expr", []grammarSymbol{"term"}, SemanticRule{"", "", []int{}}},
		{"term", []grammarSymbol{"term", "*", "factor"}, SemanticRule{"", "", []int{}}},
		{"term", []grammarSymbol{"factor"}, SemanticRule{"", "", []int{}}},
		{"factor", []grammarSymbol{"number"}, SemanticRule{"", "", []int{}}},
		{"factor", []grammarSymbol{"(", "expr", ")"}, SemanticRule{"", "", []int{}}},
	}

	var testData = []struct {
		input    grammarSymbol
		expected []Production
	}{
		{
			"expr'",
			[]Production{g.Productions[0]},
		},
		{
			"expr",
			[]Production{g.Productions[1], g.Productions[2]},
		},
		{
			"term",
			[]Production{g.Productions[3], g.Productions[4]},
		},
	}

	for _, test := range testData {
		if got := g.getProductionsOfSymbol(test.input); !reflect.DeepEqual(got, test.expected) {
			t.Errorf("Expected productions of %v to be %v but got %v", test.input, test.expected, got)
		}
	}
}

func TestComputeFirstSet(t *testing.T) {
	var g Grammar
	g.Start = "expr'"
	g.Productions = []Production{
		{"expr'", []grammarSymbol{"expr"}, SemanticRule{"", "", []int{}}},
		{"expr", []grammarSymbol{"expr", "+", "term"}, SemanticRule{"", "", []int{}}},
		{"expr", []grammarSymbol{"term"}, SemanticRule{"", "", []int{}}},
		{"term", []grammarSymbol{"term", "*", "factor"}, SemanticRule{"", "", []int{}}},
		{"term", []grammarSymbol{"factor"}, SemanticRule{"", "", []int{}}},
		{"factor", []grammarSymbol{"number"}, SemanticRule{"", "", []int{}}},
		{"factor", []grammarSymbol{"(", "expr", ")"}, SemanticRule{"", "", []int{}}},
	}

	var testData = []struct {
		input    grammarSymbol
		expected setOfSymbols
	}{
		{"number", setOfSymbols{"number": true}},
		{"expr", setOfSymbols{"(": true, "number": true}},
		{"term", setOfSymbols{"(": true, "number": true}},
		{"factor", setOfSymbols{"(": true, "number": true}},
	}

	for _, test := range testData {
		if got := g.computeFirstSet(test.input); !reflect.DeepEqual(got, test.expected) {
			t.Errorf("Expected computeFirstSet(%q) = %v, got %v", test.input, test.expected, got)
		}
	}
}

func TestComputeFollowSet(t *testing.T) {
	var g Grammar
	g.Start = "expr'"
	g.Productions = []Production{
		{"expr'", []grammarSymbol{"expr"}, SemanticRule{"", "", []int{}}},
		{"expr", []grammarSymbol{"expr", "+", "term"}, SemanticRule{"", "", []int{}}},
		{"expr", []grammarSymbol{"term"}, SemanticRule{"", "", []int{}}},
		{"term", []grammarSymbol{"term", "*", "factor"}, SemanticRule{"", "", []int{}}},
		{"term", []grammarSymbol{"factor"}, SemanticRule{"", "", []int{}}},
		{"factor", []grammarSymbol{"number"}, SemanticRule{"", "", []int{}}},
		{"factor", []grammarSymbol{"(", "expr", ")"}, SemanticRule{"", "", []int{}}},
	}

	var testData = []struct {
		input    grammarSymbol
		expected setOfSymbols
	}{
		{"expr", setOfSymbols{"+": true, ")": true, "$": true}},
		{"term", setOfSymbols{"+": true, "*": true, ")": true, "$": true}},
		{"factor", setOfSymbols{"+": true, "*": true, ")": true, "$": true}},
	}

	for _, test := range testData {
		if got := g.computeFollowSet(test.input); !reflect.DeepEqual(got, test.expected) {
			t.Errorf("Expected computeFirstSet(%q) = %v, got %v", test.input, test.expected, got)
		}
	}
}

func TestGetProductionNumber(t *testing.T) {
	var g Grammar
	g.Start = "expr'"
	g.Productions = []Production{
		{"expr'", []grammarSymbol{"expr"}, SemanticRule{"", "", []int{}}},
		{"expr", []grammarSymbol{"expr", "+", "term"}, SemanticRule{"", "", []int{}}},
		{"expr", []grammarSymbol{"term"}, SemanticRule{"", "", []int{}}},
		{"term", []grammarSymbol{"term", "*", "factor"}, SemanticRule{"", "", []int{}}},
		{"term", []grammarSymbol{"factor"}, SemanticRule{"", "", []int{}}},
		{"factor", []grammarSymbol{"number"}, SemanticRule{"", "", []int{}}},
		{"factor", []grammarSymbol{"(", "expr", ")"}, SemanticRule{"", "", []int{}}},
	}

	var testData = []struct {
		input    Production
		expected int
	}{
		{Production{"expr", []grammarSymbol{"expr", "+", "term"}, SemanticRule{"", "", []int{}}}, 1},
		{Production{"term", []grammarSymbol{"term", "*", "factor"}, SemanticRule{"", "", []int{}}}, 3},
		{Production{"factor", []grammarSymbol{"number"}, SemanticRule{"", "", []int{}}}, 5},
		{Production{"expr", []grammarSymbol{"factor"}, SemanticRule{"", "", []int{}}}, -1},
	}

	for _, test := range testData {
		if got := g.getProductionNumber(test.input); got != test.expected {
			t.Errorf("Expected getProductionNumber(%v) to be %v but got %v", test.input, test.expected, got)
		}
	}
}

func TestCompile(t *testing.T) {
	var g Grammar
	g.Start = "expr'"
	g.Productions = []Production{
		{"expr'", []grammarSymbol{"expr"}, SemanticRule{"", "", nil}},
		{"expr", []grammarSymbol{"expr", "+", "term"}, SemanticRule{"tree", "+", []int{0, 2}}},
		{"expr", []grammarSymbol{"term"}, SemanticRule{"", "", nil}},
		{"term", []grammarSymbol{"term", "*", "factor"}, SemanticRule{"tree", "*", []int{0, 2}}},
		{"term", []grammarSymbol{"factor"}, SemanticRule{"", "", nil}},
		{"factor", []grammarSymbol{"number"}, SemanticRule{"", "", nil}},
		{"factor", []grammarSymbol{"(", "expr", ")"}, SemanticRule{"copy", "", []int{1}}},
	}

	var testData = []struct {
		input    []lexer.Token
		expected bool
	}{
		{
			[]lexer.Token{
				{TokenType: "number", Lexeme: "12"},
				{TokenType: "+", Lexeme: "+"},
				{TokenType: "number", Lexeme: "8"},
				{TokenType: "*", Lexeme: "*"},
				{TokenType: "number", Lexeme: "75"},
				{TokenType: "$", Lexeme: "$"},
			},
			true,
		},
		{
			[]lexer.Token{
				{TokenType: "number", Lexeme: "12"},
				{TokenType: "+", Lexeme: "+"},
				{TokenType: "number", Lexeme: "8"},
				{TokenType: "*", Lexeme: "*"},
				{TokenType: "number", Lexeme: "75"},
			},
			false,
		},
		{
			[]lexer.Token{
				{TokenType: "(", Lexeme: "("},
				{TokenType: "number", Lexeme: "12"},
				{TokenType: "+", Lexeme: "+"},
				{TokenType: "number", Lexeme: "8"},
				{TokenType: "*", Lexeme: "*"},
				{TokenType: "number", Lexeme: "75"},
				{TokenType: "$", Lexeme: "$"},
			},
			false,
		},
		{
			[]lexer.Token{
				{TokenType: "(", Lexeme: "("},
				{TokenType: "number", Lexeme: "12"},
				{TokenType: "+", Lexeme: "+"},
				{TokenType: "number", Lexeme: "8"},
				{TokenType: "*", Lexeme: "*"},
				{TokenType: "number", Lexeme: "75"},
				{TokenType: ")", Lexeme: ")"},
				{TokenType: "$", Lexeme: "$"},
			},
			true,
		},
		{
			[]lexer.Token{
				{TokenType: "(", Lexeme: "("},
				{TokenType: "number", Lexeme: "12"},
				{TokenType: "+", Lexeme: "+"},
				{TokenType: "number", Lexeme: "8"},
				{TokenType: ")", Lexeme: ")"},
				{TokenType: "*", Lexeme: "*"},
				{TokenType: "number", Lexeme: "75"},
				{TokenType: ")", Lexeme: ")"},
				{TokenType: "$", Lexeme: "$"},
			},
			false,
		},
		{
			[]lexer.Token{
				{TokenType: "(", Lexeme: "("},
				{TokenType: "number", Lexeme: "12"},
				{TokenType: "+", Lexeme: "+"},
				{TokenType: "number", Lexeme: "8"},
				{TokenType: ")", Lexeme: ")"},
				{TokenType: "*", Lexeme: "*"},
				{TokenType: "number", Lexeme: "75"},
				{TokenType: "*", Lexeme: "*"},
				{TokenType: "number", Lexeme: "75"},
				{TokenType: "$", Lexeme: "$"},
			},
			true,
		},
	}

	ps := g.compile()
	for _, test := range testData {
		ps.parse(test.input)
		if ps.accepted != test.expected {
			t.Errorf("Expected parser to output %v on input %v, got %v",
				test.expected, test.input, ps.accepted)
		}
		ps.reset()
	}
}

func TestCompile1(t *testing.T) {
	var g Grammar
	g.Start = "S'"
	g.Productions = []Production{
		{"S'", []grammarSymbol{"S"}, SemanticRule{}},
		{"S", []grammarSymbol{"C", "C"}, SemanticRule{}},
		{"C", []grammarSymbol{"c", "C"}, SemanticRule{}},
		{"C", []grammarSymbol{"d"}, SemanticRule{}},
	}
	g.compile() // It compiles without any conflicts
}
