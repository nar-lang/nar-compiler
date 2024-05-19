package typed

import (
	"fmt"
	"github.com/nar-lang/nar-compiler/ast"
	"github.com/nar-lang/nar-compiler/bytecode"
	"github.com/nar-lang/nar-compiler/common"
)

type Update struct {
	*expressionBase
	recordName ast.Identifier
	target     Pattern
	moduleName ast.QualifiedIdentifier
	definition *Definition
	fields     []*RecordField
}

func NewUpdateGlobal(
	ctx *SolvingContext, loc ast.Location,
	moduleName ast.QualifiedIdentifier, definitionName ast.Identifier,
	targetDef *Definition, fields []*RecordField,
) Expression {
	return ctx.annotateExpression(&Update{
		expressionBase: newExpressionBase(loc),
		moduleName:     moduleName,
		recordName:     definitionName,
		definition:     targetDef,
		fields:         fields,
	})
}

func (e *Update) checkPatterns() error {
	for _, field := range e.fields {
		if err := field.value.checkPatterns(); err != nil {
			return err
		}
	}
	return nil
}

func NewUpdateLocal(
	ctx *SolvingContext, loc ast.Location,
	recordName ast.Identifier, target Pattern, fields []*RecordField,
) Expression {
	return ctx.annotateExpression(&Update{
		expressionBase: newExpressionBase(loc),
		recordName:     recordName,
		fields:         fields,
		target:         target,
	})
}

func (e *Update) mapTypes(subst map[uint64]Type) error {
	var err error
	e.type_, err = e.type_.mapTo(subst)
	if err != nil {
		return err
	}

	for _, f := range e.fields {
		f.type_, err = f.type_.mapTo(subst)
		if err != nil {
			return err
		}
		err = f.value.mapTypes(subst)
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *Update) Children() []Statement {
	return append(e.expressionBase.Children(), common.Map(func(x *RecordField) Statement { return x.value }, e.fields)...)
}

func (e *Update) Code(currentModule ast.QualifiedIdentifier) string {
	name := string(e.recordName)
	if e.moduleName != "" && currentModule != e.moduleName {
		name = string(e.moduleName) + "." + name
	}
	return fmt.Sprintf("{%s | %s}", e.recordName, common.Fold(
		func(x *RecordField, s string) string {
			if s != "" {
				s += ", "
			}
			return s + string(x.name) + " = " + x.value.Code(currentModule)
		}, "", e.fields))
}

func (e *Update) appendEquations(eqs Equations, loc *ast.Location, localDefs localTypesMap, ctx *SolvingContext, stack []*Definition) (Equations, error) {
	var err error
	fieldTypes := map[ast.Identifier]Type{}
	for _, f := range e.fields {
		fieldTypes[f.name] = f.type_
	}

	eqs = append(eqs, NewEquation(e, e.type_, NewTRecord(e.location, fieldTypes, true)))

	for _, f := range e.fields {
		l := loc
		if l == nil {
			l = &f.location
		}
		eqs = append(eqs, NewEquation(e, f.type_, f.value.Type()))
	}

	for _, f := range e.fields {
		eqs, err = f.value.appendEquations(eqs, loc, localDefs, ctx, stack)
		if err != nil {
			return nil, err
		}
	}

	if e.moduleName != "" {
		if e.definition == nil {
			return nil, common.NewErrorOf(e, "definition `%s` not found", common.MakeFullIdentifier(e.moduleName, e.recordName))
		}
		defType, err := e.definition.uniqueType(ctx, stack)
		if err != nil {
			return nil, err
		}

		eqs = append(eqs, NewEquation(e, e.type_, defType))
	}

	return eqs, nil
}

func (e *Update) appendBytecode(ops []bytecode.Op, locations []bytecode.Location, binary *bytecode.Binary, hash *bytecode.BinaryHash) ([]bytecode.Op, []bytecode.Location) {
	if e.moduleName != "" {
		id := common.MakeFullIdentifier(e.moduleName, e.recordName)
		ops, locations = bytecode.AppendLoadGlobal(hash.FuncsMap[bytecode.FullIdentifier(id)], e.location.Bytecode(), ops, locations)
	} else {
		ops, locations = bytecode.AppendLoadLocal(string(e.recordName), e.location.Bytecode(), ops, locations, binary, hash)
	}

	for _, f := range e.fields {
		ops, locations = f.value.appendBytecode(ops, locations, binary, hash)
		ops, locations = bytecode.AppendUpdate(string(f.name), f.location.Bytecode(), ops, locations, binary, hash)
	}

	return ops, locations
}
