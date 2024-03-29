package normalized

import (
	"github.com/nar-lang/nar-compiler/ast"
	"github.com/nar-lang/nar-compiler/ast/typed"
	"github.com/nar-lang/nar-compiler/common"
)

type TTuple struct {
	*typeBase
	items []Type
}

func NewTTuple(loc ast.Location, items []Type) Type {
	return &TTuple{
		typeBase: newTypeBase(loc),
		items:    items,
	}
}

func (e *TTuple) annotate(ctx *typed.SolvingContext, params typeParamsMap, source bool, placeholders placeholderMap) (typed.Type, error) {
	items, err := common.MapError(func(t Type) (typed.Type, error) {
		if t == nil {
			return nil, common.NewErrorOf(e, "tuple item type is not declared")
		}
		return t.annotate(ctx, params, source, placeholders)
	}, e.items)
	if err != nil {
		return nil, err
	}
	return e.setSuccessor(typed.NewTTuple(e.location, items))
}
