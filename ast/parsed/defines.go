package parsed

import (
	"github.com/nar-lang/nar-compiler/ast"
	"github.com/nar-lang/nar-compiler/ast/normalized"
)

type Statement interface {
	Location() ast.Location
	Iterate(f func(statement Statement))
	Successor() normalized.Statement
	_parsed()
	SemanticTokens() []ast.SemanticToken
}

type namedTypeMap map[ast.FullIdentifier]*normalized.TPlaceholder
