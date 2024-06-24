package typed

import (
	"github.com/nar-lang/nar-compiler/ast"
	"github.com/nar-lang/nar-compiler/bytecode"
)

type Type interface {
	Statement
	bytecoder
	_type()
	EqualsTo(other Type, req map[ast.FullIdentifier]struct{}) bool
	merge(other Type, loc ast.Location) (Equations, error)
	mapTo(subst map[uint64]Type) (Type, error)
	makeUnique(ctx *SolvingContext, ubMap map[uint64]uint64) Type
}

type typeBase struct {
	location ast.Location
}

func newTypeBase(loc ast.Location) *typeBase {
	return &typeBase{
		location: loc,
	}
}

func (t *typeBase) _type() {}

func (t *typeBase) Location() ast.Location {
	return t.location
}

func (t *typeBase) appendBytecode(ops []bytecode.Op, locations []bytecode.Location, binary *bytecode.Binary, hash *bytecode.BinaryHash) ([]bytecode.Op, []bytecode.Location) {
	return nil, nil
}
