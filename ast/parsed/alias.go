package parsed

import (
	"github.com/nar-lang/nar-compiler/ast"
	"github.com/nar-lang/nar-compiler/ast/normalized"
	"github.com/nar-lang/nar-compiler/common"
)

type Alias interface {
	Statement
	Name() ast.Identifier
	inferType(moduleName ast.QualifiedIdentifier, args []Type) (Type, ast.FullIdentifier, error)
	Hidden() bool
	aliasType() Type
}

func NewAlias(loc ast.Location, hidden bool, name ast.Identifier, params []ast.Identifier, type_ Type, nameLocation ast.Location) Alias {
	return &alias{
		location:     loc,
		hidden_:      hidden,
		name_:        name,
		params:       params,
		type_:        type_,
		nameLocation: nameLocation,
	}
}

type alias struct {
	location     ast.Location
	hidden_      bool
	name_        ast.Identifier
	params       []ast.Identifier
	type_        Type
	nameLocation ast.Location
}

func (a *alias) SemanticTokens() []ast.SemanticToken {
	return []ast.SemanticToken{a.nameLocation.ToToken(ast.TokenTypeType, ast.TokenModifierDeclaration)}
}

func (a *alias) aliasType() Type {
	return a.type_
}

func (a *alias) Hidden() bool {
	return a.hidden_
}

func (a *alias) inferType(moduleName ast.QualifiedIdentifier, args []Type) (Type, ast.FullIdentifier, error) {
	id := common.MakeFullIdentifier(moduleName, a.name_)
	if a.type_ == nil {
		return NewTNative(a.location, id, args, a.nameLocation), id, nil
	}
	if len(a.params) != len(args) {
		return nil, "", common.NewErrorAt(a.location, "wrong number of type parameters, expected %d, got %d", len(a.params), len(args))
	}
	typeMap := map[ast.Identifier]Type{}
	for i, x := range a.params {
		typeMap[x] = args[i]
	}
	withAppliedArgs, err := a.type_.applyArgs(typeMap, a.location)
	if err != nil {
		return nil, "", err
	}
	return withAppliedArgs, id, nil
}

func (a *alias) Successor() normalized.Statement {
	if a.type_ == nil {
		return nil
	}
	return a.type_.Successor()
}

func (a *alias) Name() ast.Identifier {
	return a.name_
}

func (a *alias) Location() ast.Location {
	return a.location
}

func (a *alias) Iterate(f func(statement Statement)) {
	f(a)
	if a.type_ != nil {
		a.type_.Iterate(f)
	}
}

func (a *alias) _parsed() {}
