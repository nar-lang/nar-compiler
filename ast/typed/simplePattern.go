package typed

import (
	"github.com/nar-lang/nar-compiler/ast"
	"github.com/nar-lang/nar-compiler/common"
	"strings"
)

type simplePattern interface {
	_simplePattern()
	String() string
}

type simpleAnything struct{}

func (simpleAnything) _simplePattern() {}

func (simpleAnything) String() string {
	return "_"
}

type simpleLiteral struct {
	Literal ast.ConstValue
}

func (simpleLiteral) _simplePattern() {}

func (p simpleLiteral) String() string {
	return p.Literal.Code("")
}

type simpleConstructor struct {
	Union *TData
	Name  ast.DataOptionIdentifier
	Args  []simplePattern
}

func (simpleConstructor) _simplePattern() {}

func (c simpleConstructor) String() string {
	sb := strings.Builder{}
	sb.WriteString(string(c.Name))
	if len(c.Args) > 0 {
		sb.WriteString("(")
		for i, a := range c.Args {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(a.String())
		}
		sb.WriteString(")")

	}
	return sb.String()
}

func (c simpleConstructor) Option() (*DataOption, error) {
	for _, o := range c.Union.options {
		if o.name == c.Name {
			return o, nil
		}
	}
	return nil, common.NewCompilerError("option not found")
}
