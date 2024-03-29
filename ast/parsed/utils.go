package parsed

import (
	"github.com/nar-lang/nar-compiler/ast"
	"github.com/nar-lang/nar-compiler/common"
)

func newAmbiguousInfixError(ids []ast.FullIdentifier, name ast.InfixIdentifier, loc ast.Location) error {
	if len(ids) == 0 {
		return common.NewErrorAt(loc, "infix definition `%s` not found", name)
	} else {
		return common.NewErrorAt(loc,
			"ambiguous infix identifier `%s`, it can be one of %s. "+
				"Use import to clarify which one to use",
			name, ast.FullIdentifiers(ids).Join(", "))
	}
}

func newAmbiguousDefinitionError(ids []ast.FullIdentifier, name ast.QualifiedIdentifier, loc ast.Location) error {
	if len(ids) == 0 {
		return common.NewErrorAt(loc, "definition `%s` not found", name)
	} else {
		return common.NewErrorAt(loc,
			"ambiguous identifier `%s`, it can be one of %s. "+
				"Use import or qualified identifer to clarify which one to use",
			name, ast.FullIdentifiers(ids).Join(", "))
	}
}
