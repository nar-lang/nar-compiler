package normalized

import (
	"github.com/nar-lang/nar-compiler/ast"
	"github.com/nar-lang/nar-compiler/ast/typed"
	"github.com/nar-lang/nar-compiler/common"
)

type PTuple struct {
	*patternBase
	items []Pattern
}

func (e *PTuple) extractLocals(locals map[ast.Identifier]Pattern) {
	for _, v := range e.items {
		v.extractLocals(locals)
	}
}

func NewPTuple(loc ast.Location, declaredType Type, items []Pattern) Pattern {
	return &PTuple{
		patternBase: newPatternBase(loc, declaredType),
		items:       items,
	}
}

func (e *PTuple) annotate(ctx *typed.SolvingContext, typeParams typeParamsMap, modules map[ast.QualifiedIdentifier]*Module, typedModules map[ast.QualifiedIdentifier]*typed.Module, moduleName ast.QualifiedIdentifier, typeMapSource bool, stack []*typed.Definition) (typed.Pattern, error) {
	items, err := common.MapError(func(x Pattern) (typed.Pattern, error) {
		return x.annotate(ctx, typeParams, modules, typedModules, moduleName, typeMapSource, stack)
	}, e.items)
	if err != nil {
		return nil, err
	}
	annotatedDeclaredType, err := annotateTypeSafe(ctx, e.declaredType, typeParams, typeMapSource)
	if err != nil {
		return nil, err
	}
	tuple, err := typed.NewPTuple(ctx, e.location, annotatedDeclaredType, items)
	if err != nil {
		return nil, err
	}
	return e.setSuccessor(tuple)
}

func (e *PTuple) Items() []Pattern {
	return e.items
}
