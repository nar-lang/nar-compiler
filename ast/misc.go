package ast

import (
	"strings"
)

type Identifier string

type QualifiedIdentifier string

type PackageIdentifier string

type InfixIdentifier string

type FullIdentifier string

func (f FullIdentifier) Module() QualifiedIdentifier {
	return QualifiedIdentifier(f[:strings.LastIndex(string(f), ".")])
}

type DataOptionIdentifier string

type Coder interface {
	Code(currentModule QualifiedIdentifier) string
}

type FullIdentifiers []FullIdentifier

func (fx FullIdentifiers) Join(sep string) string {
	sb := strings.Builder{}
	for i, f := range fx {
		if i > 0 {
			sb.WriteString(sep)
		}
		sb.WriteString(string(f))
	}
	return sep
}
