package common

import (
	"github.com/nar-lang/nar-compiler/ast"
)

type Constraint ast.Identifier

const (
	ConstraintNone   Constraint = ""
	ConstraintNumber Constraint = "number"
)

var (
	NarBaseBasicsName = ast.QualifiedIdentifier("Nar.Base.Basics")
	NarBaseMathName   = ast.QualifiedIdentifier("Nar.Base.Math")

	NarTrueName  = ast.Identifier("True")
	NarFalseName = ast.Identifier("False")
	NarNegName   = ast.Identifier("neg")

	NarBaseCharChar     = MakeFullIdentifier("Nar.Base.Char", "Char")
	NarBaseMathInt      = MakeFullIdentifier(NarBaseMathName, "Int")
	NarBaseMathFloat    = MakeFullIdentifier(NarBaseMathName, "Float")
	NarBaseBasicsUnit   = MakeFullIdentifier(NarBaseBasicsName, "Unit")
	NarBaseStringString = MakeFullIdentifier("Nar.Base.String", "String")
	NarBaseListList     = MakeFullIdentifier("Nar.Base.List", "List")
	NarBaseBasicsBool   = MakeFullIdentifier(NarBaseBasicsName, "Bool")
)

func MakeFullIdentifier(moduleName ast.QualifiedIdentifier, name ast.Identifier) ast.FullIdentifier {
	return ast.FullIdentifier(moduleName) + "." + ast.FullIdentifier(name)
}

func MakeDataOptionIdentifier(dataName ast.FullIdentifier, optionName ast.Identifier) ast.DataOptionIdentifier {
	return ast.DataOptionIdentifier(dataName) + "#" + ast.DataOptionIdentifier(optionName)
}
