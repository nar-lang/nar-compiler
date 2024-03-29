package parsed

import (
	"github.com/nar-lang/nar-compiler/ast"
	"github.com/nar-lang/nar-compiler/ast/normalized"
)

func NewTData(loc ast.Location, name ast.FullIdentifier, args []Type, options []*DataOption, nameLocation ast.Location) Type {
	return &TData{
		typeBase:     newTypeBase(loc),
		name:         name,
		args:         args,
		options:      options,
		nameLocation: nameLocation,
	}
}

type TData struct {
	*typeBase
	name         ast.FullIdentifier
	args         []Type
	options      []*DataOption
	nameLocation ast.Location
}

func (t *TData) SemanticTokens() []ast.SemanticToken {
	var tokens []ast.SemanticToken
	tokens = append(tokens, t.nameLocation.ToToken(ast.TokenTypeEnum))
	for _, opt := range t.options {
		tokens = append(tokens, opt.nameLocation.ToToken(ast.TokenTypeEnumMember))
	}
	return tokens
}

func (t *TData) Iterate(f func(statement Statement)) {
	f(t)
	for _, arg := range t.args {
		if arg != nil {
			arg.Iterate(f)
		}
	}
	for _, opt := range t.options {
		for _, value := range opt.values {
			if value != nil {
				value.Iterate(f)
			}
		}
	}
}

func (t *TData) normalize(modules map[ast.QualifiedIdentifier]*Module, module *Module, namedTypes namedTypeMap) (normalized.Type, error) {
	if namedTypes == nil {
		namedTypes = namedTypeMap{}
	}
	if placeholder, cached := namedTypes[t.name]; cached {
		return placeholder, nil
	}
	namedTypes[t.name] = normalized.NewTPlaceholder(t.name).(*normalized.TPlaceholder)

	var args []normalized.Type
	for _, arg := range t.args {
		nArg, err := arg.normalize(modules, module, namedTypes)
		if err != nil {
			return nil, err
		}
		args = append(args, nArg)
	}
	var options []*normalized.DataOption
	for _, option := range t.options {
		var values []normalized.Type
		for _, value := range option.values {
			nValue, err := value.normalize(modules, module, namedTypes)
			if err != nil {
				return nil, err
			}
			values = append(values, nValue)
		}
		options = append(options, normalized.NewDataOption(option.name, option.hidden, values))
	}
	return t.setSuccessor(normalized.NewTData(t.location, t.name, args, options))
}

func (t *TData) applyArgs(params map[ast.Identifier]Type, loc ast.Location) (Type, error) {
	var args []Type
	for _, arg := range t.args {
		nArg, err := arg.applyArgs(params, loc)
		if err != nil {
			return nil, err
		}
		args = append(args, nArg)
	}
	return NewTData(loc, t.name, args, t.options, t.nameLocation), nil
}

type DataOption struct {
	name         ast.Identifier
	hidden       bool
	values       []Type
	nameLocation ast.Location
}

func NewDataOption(name ast.Identifier, hidden bool, values []Type, nameLocation ast.Location) *DataOption {
	return &DataOption{
		name:         name,
		hidden:       hidden,
		values:       values,
		nameLocation: nameLocation,
	}
}
