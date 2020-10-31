package lexer

import "testing"

func TestRegularExpressionOperatorLength(t *testing.T) {
	var testData = []struct {
		input    regularExpressionOperator
		expected int
	}{
		{union, 1},
		{concat, 0},
		{star, 1},
	}

	for _, test := range testData {
		if got := test.input.length(); got != test.expected {
			t.Errorf("expected length to be %v, got %v", test.expected, got)
		}
	}
}

func TestRegularExpressionIsValid(t *testing.T) {
	var testData = []struct {
		input    regularExpression
		expected bool
	}{
		{")asdfdf(", false},
		{"((sdaffd)", false},
		{"(((dsaf)|sdf)sdf)", true},
		{"sadfa", true},
	}

	for _, test := range testData {
		if got := test.input.isValid(); got != test.expected {
			t.Errorf("Expected %v.isValid() = %v, got %v", test.input, test.expected, got)
		}
	}
}

func TestRegularExpressionGetMatchingParenthesis(t *testing.T) {
	var testData = []struct {
		input    regularExpression
		expected int
	}{
		{"((sdaffd))dasfdaf", 9},
		{"(((dsaf)|sdf)sdf)(sdafdfadsf)", 16},
		{"sadfa", -1},
	}

	for _, test := range testData {
		if got := test.input.getMatchingParenIndex(); got != test.expected {
			t.Errorf("Expected %v.getMatchingParenthesisIndex() = %v, got %v", test.input, test.expected, got)
		}
	}
}

func TestRegularExpressionTrimParenthesis(t *testing.T) {
	var testData = []struct {
		input    regularExpression
		expected regularExpression
	}{
		{"(sadf|asdf)", "sadf|asdf"},
		{"sdfasf", "sdfasf"},
		{"(sdfsdf)*", "(sdfsdf)*"},
	}

	for _, test := range testData {
		if got := test.input.trimParenthesis(); got != test.expected {
			t.Errorf("Expected %v.trimParenthesis() = %v, got %v", test.input, test.expected, got)
		}
	}
}

func TestRegularExpressionGetOperandsAndOperator(t *testing.T) {
	var testData = []struct {
		input                 regularExpression
		expectedFirstOperand  regularExpression
		expectedOperator      regularExpressionOperator
		expectedSecondOperand regularExpression
	}{
		{"s|f", "s", union, "f"},
		{"sdfasf", "s", concat, "dfasf"},
		{"(a|f)(c|d)*", "(a|f)", concat, "(c|d)*"},
		{"(sdfsdf)*", "(sdfsdf)", star, ""},
	}

	for _, test := range testData {
		if got := test.input.getFirstOperand(); got != test.expectedFirstOperand {
			t.Errorf("Expected %v.getFirstOperand() = %v, got %v", test.input, test.expectedFirstOperand, got)
		}
	}
	for _, test := range testData {
		if got := test.input.getOperator(); got != test.expectedOperator {
			t.Errorf("Expected %v.getOperator() = %v, got %v", test.input, test.expectedOperator, got)
		}
	}
	for _, test := range testData {
		if got := test.input.getSecondOperand(); got != test.expectedSecondOperand {
			t.Errorf("Expected %v.getSecondOperand() = %v, got %v", test.input, test.expectedSecondOperand, got)
		}
	}
}
