package ast

import (
	"reflect"

	"github.com/SaurabhJha/lexpar/lexer"
	"github.com/SaurabhJha/lexpar/parser"
)

func getProductionNumber(g parser.Grammar, p parser.Production) int {
	for i, production := range g.Productions {
		if reflect.DeepEqual(production, p) {
			return i
		}
	}
	return -1
}

func getSymbolsToCheck(p parser.Production) []string {
	symbolsToCheck := make([]string, 0, 10)
	for _, symbol := range p.Body {
		symbolsToCheck = append(symbolsToCheck, string(symbol))
	}
	return symbolsToCheck
}

// SemanticRule represents a syntax directed definition to construct a tree node.
type SemanticRule struct {
	ForProductions []int
	Type           string
	Children       []int
	RootLabel      string
}

func (rule *SemanticRule) executeCopy(stck *stack, reduction parser.Production) {
	stackOffset := len(*stck) - len(reduction.Body)
	recordToCopy := (*stck)[stackOffset+rule.Children[0]]
	recordToCopy.symbol = string(reduction.Head)
	for range reduction.Body {
		stck.pop()
	}
	stck.push(recordToCopy)
}

func (rule *SemanticRule) executeTree(stck *stack, reduction parser.Production, ast *abstractSyntaxTree) {
	parentNodeValue := ast.getNextNodeValue()
	stackOffset := len(*stck) - len(reduction.Body)
	for _, child := range rule.Children {
		ast.g[parentNodeValue] = append(ast.g[parentNodeValue], (*stck)[child+stackOffset].nodeIndex)
	}
	ast.n[parentNodeValue] = rule.RootLabel
	for range reduction.Body {
		stck.pop()
	}
	stck.push(stackRecord{parentNodeValue, string(reduction.Head)})
}

// A SyntaxDirectedTranslator object provides all the functionality necessary to implement rules defined
// as SemanticRule objects.
type SyntaxDirectedTranslator struct {
	g                 parser.Grammar
	grammarToRulesMap map[int]SemanticRule
}

// Init function of SyntaxDirectedTranslator initialises all the state necessary to implement syntax directed
// translation.
func (sdt *SyntaxDirectedTranslator) Init(g parser.Grammar, rules []SemanticRule) {
	sdt.g = g
	sdt.grammarToRulesMap = make(map[int]SemanticRule)
	for _, rule := range rules {
		for _, productionNumber := range rule.ForProductions {
			sdt.grammarToRulesMap[productionNumber] = rule
		}
	}
}

func (sdt *SyntaxDirectedTranslator) getRule(reduction parser.Production) (SemanticRule, bool) {
	productionNumber := getProductionNumber(sdt.g, reduction)
	rule, ok := sdt.grammarToRulesMap[productionNumber]
	return rule, ok
}

// ConstructAST takes the stream of tokens and the reductions as arguments and constructs and AST using
// SemanticRules.
func (sdt *SyntaxDirectedTranslator) ConstructAST(tokens []lexer.Token, reductions []parser.Production) abstractSyntaxTree {
	currentTokenIndex := 0
	stck := make(stack, 0, 10)
	var ast abstractSyntaxTree
	ast.init()
	for _, reduction := range reductions {
		symbolsToCheck := getSymbolsToCheck(reduction)
		for !stck.hasSymbolsOnTop(symbolsToCheck) {
			nextNodeValue := ast.getNextNodeValue()
			token := tokens[currentTokenIndex]
			ast.addNode(nextNodeValue, token.Lexeme())
			stck.push(stackRecord{nextNodeValue, token.TokenType()})
			currentTokenIndex++
		}
		rule, ok := sdt.getRule(reduction)
		if !ok {
			rule = SemanticRule{[]int{}, "copy", []int{0}, string(reduction.Head)}
		}
		switch rule.Type {
		case "copy":
			rule.executeCopy(&stck, reduction)
		case "tree":
			rule.executeTree(&stck, reduction, &ast)
		}
	}

	ast.currentNodeValue--
	return ast
}
