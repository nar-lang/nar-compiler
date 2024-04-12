package normalized

import (
	"github.com/nar-lang/nar-compiler/ast"
	"github.com/nar-lang/nar-compiler/ast/typed"
)

type Pattern interface {
	Statement
	extractLocals(locals map[ast.Identifier]Pattern)
	annotate(ctx *typed.SolvingContext, typeParams typeParamsMap, modules map[ast.QualifiedIdentifier]*Module, typedModules map[ast.QualifiedIdentifier]*typed.Module, moduleName ast.QualifiedIdentifier, typeMapSource bool, stack []*typed.Definition) (typed.Pattern, error)
}

type patternBase struct {
	location     ast.Location
	declaredType Type
	successor    typed.Pattern
}

func newPatternBase(loc ast.Location, declaredType Type) *patternBase {
	return &patternBase{location: loc, declaredType: declaredType}
}

func (p *patternBase) Location() ast.Location {
	return p.location
}

func (p *patternBase) Successor() typed.Statement {
	if p.successor == nil {
		return nil
	}
	return p.successor
}

func (p *patternBase) setSuccessor(typedPattern typed.Pattern) (typed.Pattern, error) {
	p.successor = typedPattern
	return typedPattern, nil
}
