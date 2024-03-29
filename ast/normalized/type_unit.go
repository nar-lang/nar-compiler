package normalized

import (
	"github.com/nar-lang/nar-compiler/ast"
	"github.com/nar-lang/nar-compiler/ast/typed"
	"github.com/nar-lang/nar-compiler/common"
)

type TUnit struct {
	*typeBase
}

func NewTUnit(loc ast.Location) Type {
	return &TUnit{
		typeBase: newTypeBase(loc),
	}
}

func (e *TUnit) annotate(ctx *typed.SolvingContext, params typeParamsMap, source bool, placeholders placeholderMap) (typed.Type, error) {
	return e.setSuccessor(typed.NewTNative(e.location, common.NarBaseBasicsUnit, nil))
}
