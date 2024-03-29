package parsed

import (
	"github.com/nar-lang/nar-compiler/ast"
	"github.com/nar-lang/nar-compiler/ast/normalized"
)

type Infix interface {
	name() ast.InfixIdentifier
	hasLowerPrecedenceThan(fn Infix) bool
	alias() ast.Identifier
	hidden() bool
	Name() ast.InfixIdentifier
	Location() ast.Location
}

func NewInfix(
	loc ast.Location, hidden bool, name ast.InfixIdentifier, associativity Associativity,
	precedence int, aliasLoc ast.Location, alias ast.Identifier,
) Infix {
	return &infix{
		location:      loc,
		hidden_:       hidden,
		name_:         name,
		associativity: associativity,
		precedence:    precedence,
		aliasLocation: aliasLoc,
		alias_:        alias,
	}
}

type infix struct {
	location      ast.Location
	hidden_       bool
	name_         ast.InfixIdentifier
	associativity Associativity
	precedence    int
	aliasLocation ast.Location
	alias_        ast.Identifier
	successor     normalized.Statement
}

func (i *infix) Location() ast.Location {
	return i.location
}

func (i *infix) Name() ast.InfixIdentifier {
	return i.name_
}

func (i *infix) hidden() bool {
	return i.hidden_
}

func (i *infix) alias() ast.Identifier {
	return i.alias_
}

func (i *infix) name() ast.InfixIdentifier {
	return i.name_
}

func (i *infix) hasLowerPrecedenceThan(other Infix) bool {
	i2 := other.(*infix)
	return i2.precedence > i.precedence ||
		(i2.precedence == i.precedence && i.associativity == Left)
}

type Associativity int

const (
	Left  Associativity = -1
	None                = 0
	Right               = 1
)
