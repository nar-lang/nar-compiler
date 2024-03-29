package parsed

import (
	"github.com/nar-lang/nar-compiler/ast"
	"github.com/nar-lang/nar-compiler/ast/normalized"
)

func NewTUnit(loc ast.Location) Type {
	return &TUnit{
		typeBase: newTypeBase(loc),
	}
}

type TUnit struct {
	*typeBase
}

func (t *TUnit) SemanticTokens() []ast.SemanticToken {
	return []ast.SemanticToken{t.location.ToToken(ast.TokenTypeRegexp)}
}

func (t *TUnit) Iterate(f func(statement Statement)) {
	f(t)
}

func (t *TUnit) normalize(modules map[ast.QualifiedIdentifier]*Module, module *Module, namedTypes namedTypeMap) (normalized.Type, error) {
	return t.setSuccessor(normalized.NewTUnit(t.location))
}

func (t *TUnit) applyArgs(params map[ast.Identifier]Type, loc ast.Location) (Type, error) {
	return t, nil
}
