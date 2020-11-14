package ast

type stackRecord struct {
	nodeIndex int
	symbol    string
}

type stack []stackRecord

func (s *stack) push(record stackRecord) {
	*s = append(*s, record)
}

func (s *stack) pop() stackRecord {
	top := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return top
}

func (s *stack) peek(reverseIdx int) stackRecord {
	return (*s)[len(*s)-reverseIdx-1]
}

func (s *stack) hasSymbolsOnTop(symbols []string) bool {
	if len(*s) < len(symbols) {
		return false
	}

	for i, symbol := range symbols {
		peekIndex := len(symbols) - i - 1
		record := s.peek(peekIndex)
		if record.symbol != symbol {
			return false
		}
	}

	return true
}

type graph map[int][]int

type nodeMap map[int]string

type abstractSyntaxTree struct {
	g                graph
	n                nodeMap
	currentNodeValue int
}

func (ast *abstractSyntaxTree) init() {
	(*ast).g = make(graph)
	(*ast).n = make(nodeMap)
	(*ast).currentNodeValue = 0
}

func (ast *abstractSyntaxTree) addNode(nodeValue int, nodeContent string) {
	(*ast).n[nodeValue] = nodeContent
}

func (ast *abstractSyntaxTree) getNextNodeValue() int {
	value := ast.currentNodeValue
	ast.currentNodeValue++
	return value
}
