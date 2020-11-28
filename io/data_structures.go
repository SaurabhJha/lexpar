package io

import (
	"github.com/SaurabhJha/lexpar/lexer"
	"github.com/SaurabhJha/lexpar/parser"
)

// DefinitionsTable is used to marshal input json into a data structure
type DefinitionsTable struct {
	RegularExpressions map[string]lexer.RegularExpression `json:"regularExpressions"`
	Grammar            parser.Grammar                     `json:"grammar"`
}
