package normalized

import (
	"github.com/nar-lang/nar-compiler/ast"
	"github.com/nar-lang/nar-compiler/ast/typed"
)

type Statement interface {
	Location() ast.Location
	Successor() typed.Statement
}

type WithSuccessor interface {
	SetSuccessor(s Expression)
}

type placeholderMap map[ast.FullIdentifier]typed.Type

type typeParamsMap map[ast.Identifier]typed.Type

var LastDefinitionId = uint64(0)

var lastLambdaId = uint64(0)
