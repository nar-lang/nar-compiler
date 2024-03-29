package parsed

import (
	"github.com/nar-lang/nar-compiler/ast"
	"github.com/nar-lang/nar-compiler/ast/normalized"
)

func NewInfixVar(location ast.Location, infix ast.InfixIdentifier) Expression {
	return &InfixVar{
		expressionBase: newExpressionBase(location),
		infix:          infix,
	}
}

type InfixVar struct {
	*expressionBase
	infix ast.InfixIdentifier
}

func (e *InfixVar) SemanticTokens() []ast.SemanticToken {
	return []ast.SemanticToken{e.location.ToToken(ast.TokenTypeFunction)}
}

func (e *InfixVar) Iterate(f func(statement Statement)) {
	f(e)
}

func (e *InfixVar) normalize(
	locals map[ast.Identifier]normalized.Pattern,
	modules map[ast.QualifiedIdentifier]*Module,
	module *Module,
	normalizedModule *normalized.Module,
) (normalized.Expression, error) {
	if i, m, ids := module.findInfixFn(modules, e.infix); len(ids) != 1 {
		return nil, newAmbiguousInfixError(ids, e.infix, e.location)
	} else if d, _, ids := m.findDefinitionAndAddDependency(nil, ast.QualifiedIdentifier(i.alias()), normalizedModule); len(ids) != 1 {
		return nil, newAmbiguousDefinitionError(ids, ast.QualifiedIdentifier(i.alias()), e.location)
	} else {
		return e.setSuccessor(normalized.NewGlobal(e.location, m.name, d.Name()))
	}
}
