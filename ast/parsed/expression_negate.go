package parsed

import (
	"github.com/nar-lang/nar-compiler/ast"
	"github.com/nar-lang/nar-compiler/ast/normalized"
	"github.com/nar-lang/nar-compiler/common"
)

func NewNegate(location ast.Location, nested Expression) Expression {
	return &Negate{
		expressionBase: newExpressionBase(location),
		nested:         nested,
	}
}

type Negate struct {
	*expressionBase
	nested Expression
}

func (e *Negate) SemanticTokens() []ast.SemanticToken {
	return nil
}

func (e *Negate) Iterate(f func(statement Statement)) {
	f(e)
	if e.nested != nil {
		e.nested.Iterate(f)
	}
}

func (e *Negate) normalize(
	locals map[ast.Identifier]normalized.Pattern,
	modules map[ast.QualifiedIdentifier]*Module,
	module *Module,
	normalizedModule *normalized.Module,
) (normalized.Expression, error) {
	nested, err := e.nested.normalize(locals, modules, module, normalizedModule)
	if err != nil {
		return nil, err
	}
	return e.setSuccessor(normalized.NewApply(
		e.location,
		normalized.NewGlobal(e.location, common.NarBaseMathName, common.NarNegName),
		[]normalized.Expression{nested},
	))
}
