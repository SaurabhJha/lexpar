package io

import "github.com/SaurabhJha/lexpar/parser"

type regularExpressions map[string]string

// DefinitionsTable is used to marshal input json into a data structure
type DefinitionsTable struct {
	RegularExpressions map[string]string `json:"regularExpressions"`
	Grammar            parser.Grammar    `json:"grammar"`
}
