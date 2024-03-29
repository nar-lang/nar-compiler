package typed

import (
	"github.com/nar-lang/nar-compiler/ast"
	"github.com/nar-lang/nar-compiler/bytecode"
)

type Statement interface {
	ast.Coder
	Location() ast.Location
	Children() []Statement
}

type bytecoder interface {
	appendBytecode(ops []bytecode.Op, locations []bytecode.Location, binary *bytecode.Binary, hash *bytecode.BinaryHash) ([]bytecode.Op, []bytecode.Location)
}

type localTypesMap map[ast.Identifier]Type

type TypePredecessor interface {
	SetSuccessor(Type) Type
}
