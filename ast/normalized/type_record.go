package normalized

import (
	"github.com/nar-lang/nar-compiler/ast"
	"github.com/nar-lang/nar-compiler/ast/typed"
	"github.com/nar-lang/nar-compiler/common"
)

type TRecord struct {
	*typeBase
	fields map[ast.Identifier]Type
}

func NewTRecord(loc ast.Location, fields map[ast.Identifier]Type) Type {
	return &TRecord{
		typeBase: newTypeBase(loc),
		fields:   fields,
	}
}

func (e *TRecord) annotate(ctx *typed.SolvingContext, params typeParamsMap, source bool, placeholders placeholderMap) (typed.Type, error) {
	fields := map[ast.Identifier]typed.Type{}
	for n, v := range e.fields {
		if v == nil {
			return nil, common.NewErrorOf(e, "record field type is not declared")
		}
		var err error
		fields[n], err = v.annotate(ctx, params, source, placeholders)
		if err != nil {
			return nil, err
		}
	}
	return e.setSuccessor(typed.NewTRecord(e.location, fields, false))
}
