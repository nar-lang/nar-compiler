package parsed

import (
	"fmt"
	"github.com/nar-lang/nar-compiler/ast"
	"github.com/nar-lang/nar-compiler/ast/normalized"
	"github.com/nar-lang/nar-compiler/common"
	"strings"
)

func NewTNamed(loc ast.Location, name ast.QualifiedIdentifier, args []Type, nameLocation ast.Location) Type {
	return &TNamed{
		typeBase:     newTypeBase(loc),
		name:         name,
		args:         args,
		nameLocation: nameLocation,
	}
}

type TNamed struct {
	*typeBase
	name         ast.QualifiedIdentifier
	args         []Type
	nameLocation ast.Location
}

func (t *TNamed) SemanticTokens() []ast.SemanticToken {
	return []ast.SemanticToken{t.nameLocation.ToToken(ast.TokenTypeType)}
}

func (t *TNamed) Iterate(f func(statement Statement)) {
	f(t)
	for _, arg := range t.args {
		if arg != nil {
			arg.Iterate(f)
		}
	}
}

func (t *TNamed) Find(
	modules map[ast.QualifiedIdentifier]*Module, module *Module,
) (Type, *Module, []ast.FullIdentifier, error) {
	return module.findType(modules, t.name, t.args, t.location)
}

func (t *TNamed) normalize(modules map[ast.QualifiedIdentifier]*Module, module *Module, namedTypes namedTypeMap) (normalized.Type, error) {
	x, _, ids, err := t.Find(modules, module)
	if err != nil {
		return nil, err
	}
	if ids == nil {
		args := ""
		if len(t.args) > 0 {
			args = fmt.Sprintf("[%s]", strings.Join(common.Repeat("_", len(t.args)), ", "))
		}
		return nil, common.NewErrorOf(t, "type `%s%s` not found", t.name, args)
	}
	if len(ids) > 1 {
		return nil, common.NewErrorOf(t,
			"ambiguous type `%s`, it can be one of %s. "+
				"Use import or qualified Name to clarify which one to use",
			t.name, ast.FullIdentifiers(ids).Join(", "))
	}
	if named, ok := x.(*TNamed); ok {
		if named.name == t.name {
			return nil, common.NewErrorOf(named, "type `%s` aliased to itself", t.name)
		}
	}

	nType, err := x.normalize(modules, module, namedTypes)
	if err != nil {
		return nil, err
	}
	return t.setSuccessor(nType)
}

func (t *TNamed) applyArgs(params map[ast.Identifier]Type, loc ast.Location) (Type, error) {
	var args []Type
	for _, arg := range t.args {
		nArg, err := arg.applyArgs(params, loc)
		if err != nil {
			return nil, err
		}
		args = append(args, nArg)
	}
	return NewTNamed(loc, t.name, args, t.nameLocation), nil
}
