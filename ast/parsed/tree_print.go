package parsed

import (
	"fmt"
	"github.com/nar-lang/nar-compiler/ast"
	"sort"
	"strings"
)

// stringTree renders a parsed AST node as a multi-line tree string.
// One line per node, indented with `\t * offset`. Children are rendered
// with `offset + 1`. The implementation lives in this file (next to all
// concrete types) so it can read package-private fields directly.
//
// Each parsed AST type exposes a `StringTree(offset int) string` method
// (declared on the various per-category interfaces in this package) that
// forwards here. The dispatcher accepts `any` because the top-level
// container types (Module/Import/Infix) do not implement Statement.
func stringTree(s any, offset int) string {
	if s == nil {
		return ""
	}
	switch n := s.(type) {
	case *Module:
		return treeModule(n, offset)
	case *import_:
		return treeImport(n, offset)
	case *alias:
		return treeAlias(n, offset)
	case *infix:
		return treeInfix(n, offset)
	case *dataType:
		return treeDataType(n, offset)
	case *dataTypeOption:
		// DataTypeOption is normally rendered inline by treeDataType;
		// this branch is kept so that iterating arbitrary subtrees works.
		return treeDataTypeOption(n, offset)
	case *definition:
		return treeDefinition(n, offset)

	case *Access:
		return treeAccess(n, offset)
	case *Accessor:
		return treeAccessor(n, offset)
	case *Apply:
		return treeApply(n, offset)
	case *BinOp:
		return treeBinOp(n, offset)
	case *Call:
		return treeCall(n, offset)
	case *Const:
		return treeConst(n, offset)
	case *Constructor:
		return treeConstructor(n, offset)
	case *Function:
		return treeFunction(n, offset)
	case *If:
		return treeIf(n, offset)
	case *InfixVar:
		return treeInfixVar(n, offset)
	case *Lambda:
		return treeLambda(n, offset)
	case *Let:
		return treeLet(n, offset)
	case *List:
		return treeList(n, offset)
	case *Negate:
		return treeNegate(n, offset)
	case *Record:
		return treeRecord(n, offset)
	case *Select:
		return treeSelect(n, offset)
	case *Tuple:
		return treeTuple(n, offset)
	case *Update:
		return treeUpdate(n, offset)
	case *Var:
		return treeVar(n, offset)

	case *PAlias:
		return treePAlias(n, offset)
	case *PAny:
		return treePAny(n, offset)
	case *PCons:
		return treePCons(n, offset)
	case *PConst:
		return treePConst(n, offset)
	case *PList:
		return treePList(n, offset)
	case *PNamed:
		return treePNamed(n, offset)
	case *POption:
		return treePOption(n, offset)
	case *PRecord:
		return treePRecord(n, offset)
	case *PTuple:
		return treePTuple(n, offset)

	case *TData:
		return treeTData(n, offset)
	case *TFunc:
		return treeTFunc(n, offset)
	case *TNamed:
		return treeTNamed(n, offset)
	case *TNative:
		return treeTNative(n, offset)
	case *TParameter:
		return treeTParameter(n, offset)
	case *TRecord:
		return treeTRecord(n, offset)
	case *TTuple:
		return treeTTuple(n, offset)
	case *TUnit:
		return treeTUnit(n, offset)
	}
	return indent(offset) + fmt.Sprintf("?(%T)\n", s)
}

// ---------------------------------------------------------------- helpers --

func indent(offset int) string {
	return strings.Repeat("\t", offset)
}

func formatBool(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

func formatName(s string) string {
	if s == "" {
		return "<nil>"
	}
	return s
}

func formatIdentPtr(p *ast.Identifier) string {
	if p == nil {
		return "<nil>"
	}
	return formatName(string(*p))
}

func formatStringList(list []string) string {
	if len(list) == 0 {
		return "[]"
	}
	return "[" + strings.Join(list, ", ") + "]"
}

func formatIdentList(list []ast.Identifier) string {
	if len(list) == 0 {
		return "[]"
	}
	parts := make([]string, len(list))
	for i, x := range list {
		parts[i] = string(x)
	}
	return "[" + strings.Join(parts, ", ") + "]"
}

func formatAssoc(a Associativity) string {
	switch a {
	case Left:
		return "LEFT"
	case Right:
		return "RIGHT"
	}
	return "NONE"
}

func formatConst(v ast.ConstValue) string {
	if v == nil {
		return "nil"
	}
	switch c := v.(type) {
	case ast.CInt:
		return fmt.Sprintf("CInt(%d)", c.Value)
	case ast.CFloat:
		return fmt.Sprintf("CFloat(%s)", strconvFloat(c.Value))
	case ast.CString:
		return fmt.Sprintf("CString(%q)", c.Value)
	case ast.CChar:
		return fmt.Sprintf("CChar(%q)", string(c.Value))
	case ast.CUnit:
		return "CUnit"
	}
	return fmt.Sprintf("%T", v)
}

// strconvFloat formats a float using the same syntax as Lua's tostring
// for parity-friendly output (no trailing zeros for integers, otherwise
// the shortest decimal representation).
func strconvFloat(f float64) string {
	return fmt.Sprintf("%g", f)
}

// ---------------------------------------------------------- module-level --

func treeModule(m *Module, offset int) string {
	var b strings.Builder
	b.WriteString(indent(offset))
	b.WriteString(fmt.Sprintf("Module(name=%s)\n",
		formatName(string(m.name))))
	for _, x := range m.imports {
		b.WriteString(x.(*import_).StringTree(offset + 1))
	}
	for _, x := range m.dataTypes {
		b.WriteString(x.(*dataType).StringTree(offset + 1))
	}
	for _, x := range m.aliases {
		b.WriteString(x.(*alias).StringTree(offset + 1))
	}
	for _, x := range m.infixFns {
		b.WriteString(x.(*infix).StringTree(offset + 1))
	}
	for _, x := range m.definitions {
		b.WriteString(x.(*definition).StringTree(offset + 1))
	}
	return b.String()
}

func treeImport(i *import_, offset int) string {
	return indent(offset) + fmt.Sprintf(
		"Import(module=%s, alias=%s, exposingAll=%s, exposing=%s)\n",
		formatName(string(i.module_)),
		formatIdentPtr(i.alias),
		formatBool(i.exposingAll),
		formatStringList(i.exposing),
	)
}

func treeAlias(a *alias, offset int) string {
	var b strings.Builder
	b.WriteString(indent(offset))
	b.WriteString(fmt.Sprintf("Alias(name=%s, hidden=%s, params=%s)\n",
		formatName(string(a.name_)), formatBool(a.hidden_), formatIdentList(a.params)))
	if a.type_ != nil {
		b.WriteString(stringTree(a.type_, offset+1))
	}
	return b.String()
}

func treeInfix(i *infix, offset int) string {
	return indent(offset) + fmt.Sprintf(
		"Infix(name=%s, hidden=%s, assoc=%s, precedence=%d, alias=%s)\n",
		formatName(string(i.name_)), formatBool(i.hidden_),
		formatAssoc(i.associativity), i.precedence,
		formatName(string(i.alias_)),
	)
}

func treeDataType(d *dataType, offset int) string {
	var b strings.Builder
	b.WriteString(indent(offset))
	b.WriteString(fmt.Sprintf("DataType(name=%s, hidden=%s, params=%s)\n",
		formatName(string(d.name)), formatBool(d.hidden), formatIdentList(d.params)))
	for _, opt := range d.options {
		o := opt.(*dataTypeOption)
		b.WriteString(indent(offset + 1))
		b.WriteString(fmt.Sprintf("DataTypeOption(name=%s, hidden=%s)\n",
			formatName(string(o.name)), formatBool(o.hidden)))
		for _, val := range o.values {
			if val == nil {
				continue
			}
			b.WriteString(indent(offset + 2))
			b.WriteString(fmt.Sprintf("DataTypeValue(name=%s)\n", formatName(string(val.name))))
			if val.type_ != nil {
				b.WriteString(val.type_.StringTree(offset + 3))
			}
		}
	}
	return b.String()
}

func treeDataTypeOption(o *dataTypeOption, offset int) string {
	var b strings.Builder
	b.WriteString(indent(offset))
	b.WriteString(fmt.Sprintf("DataTypeOption(name=%s, hidden=%s)\n",
		formatName(string(o.name)), formatBool(o.hidden)))
	for _, val := range o.values {
		if val == nil {
			continue
		}
		b.WriteString(indent(offset + 1))
		b.WriteString(fmt.Sprintf("DataTypeValue(name=%s)\n", formatName(string(val.name))))
		if val.type_ != nil {
			b.WriteString(stringTree(val.type_, offset+2))
		}
	}
	return b.String()
}

func treeDefinition(d *definition, offset int) string {
	var b strings.Builder
	b.WriteString(indent(offset))
	b.WriteString(fmt.Sprintf("Definition(name=%s, hidden=%s)\n",
		formatName(string(d.name_)), formatBool(d.hidden_)))
	for _, p := range d.params {
		if p != nil {
			b.WriteString(stringTree(p, offset+1))
		}
	}
	if d.body != nil {
		b.WriteString(stringTree(d.body, offset+1))
	}
	if d.declaredType != nil {
		b.WriteString(stringTree(d.declaredType, offset+1))
	}
	return b.String()
}

// -------------------------------------------------------------- expressions --

func treeAccess(e *Access, offset int) string {
	var b strings.Builder
	b.WriteString(indent(offset))
	b.WriteString(fmt.Sprintf("Access(fieldName=%s)\n", formatName(string(e.fieldName))))
	if e.record != nil {
		b.WriteString(stringTree(e.record, offset+1))
	}
	return b.String()
}

func treeAccessor(e *Accessor, offset int) string {
	return indent(offset) + fmt.Sprintf("Accessor(fieldName=%s)\n", formatName(string(e.fieldName)))
}

func treeApply(e *Apply, offset int) string {
	var b strings.Builder
	b.WriteString(indent(offset) + "Apply()\n")
	if e.func_ != nil {
		b.WriteString(stringTree(e.func_, offset+1))
	}
	for _, a := range e.args {
		if a != nil {
			b.WriteString(stringTree(a, offset+1))
		}
	}
	return b.String()
}

func treeBinOp(e *BinOp, offset int) string {
	var b strings.Builder
	b.WriteString(indent(offset))
	b.WriteString(fmt.Sprintf("BinOp(inParentheses=%s)\n", formatBool(e.inParentheses)))
	for _, item := range e.items {
		if item == nil {
			continue
		}
		if item.operand != nil {
			b.WriteString(indent(offset+1) + "BinOpItem(operand)\n")
			b.WriteString(stringTree(item.operand, offset+2))
		} else {
			b.WriteString(indent(offset+1) + fmt.Sprintf(
				"BinOpItem(infix=%s)\n", formatName(string(item.infix))))
		}
	}
	return b.String()
}

func treeCall(e *Call, offset int) string {
	var b strings.Builder
	b.WriteString(indent(offset))
	b.WriteString(fmt.Sprintf("Call(name=%s)\n", formatName(string(e.name))))
	for _, a := range e.args {
		if a != nil {
			b.WriteString(stringTree(a, offset+1))
		}
	}
	return b.String()
}

func treeConst(e *Const, offset int) string {
	return indent(offset) + fmt.Sprintf("Const(value=%s)\n", formatConst(e.value))
}

func treeConstructor(e *Constructor, offset int) string {
	var b strings.Builder
	b.WriteString(indent(offset))
	b.WriteString(fmt.Sprintf(
		"Constructor(moduleName=%s, dataName=%s, optionName=%s)\n",
		formatName(string(e.moduleName)), formatName(string(e.dataName)), formatName(string(e.optionName)),
	))
	for _, a := range e.args {
		if a != nil {
			b.WriteString(stringTree(a, offset+1))
		}
	}
	return b.String()
}

func treeFunction(e *Function, offset int) string {
	var b strings.Builder
	b.WriteString(indent(offset))
	b.WriteString(fmt.Sprintf("Function(name=%s)\n", formatName(string(e.name))))
	for _, p := range e.params {
		if p != nil {
			b.WriteString(stringTree(p, offset+1))
		}
	}
	if e.body != nil {
		b.WriteString(stringTree(e.body, offset+1))
	}
	if e.declaredType != nil {
		b.WriteString(stringTree(e.declaredType, offset+1))
	}
	if e.nested != nil {
		b.WriteString(stringTree(e.nested, offset+1))
	}
	return b.String()
}

func treeIf(e *If, offset int) string {
	var b strings.Builder
	b.WriteString(indent(offset) + "If()\n")
	if e.condition != nil {
		b.WriteString(stringTree(e.condition, offset+1))
	}
	if e.positive != nil {
		b.WriteString(stringTree(e.positive, offset+1))
	}
	if e.negative != nil {
		b.WriteString(stringTree(e.negative, offset+1))
	}
	return b.String()
}

func treeInfixVar(e *InfixVar, offset int) string {
	return indent(offset) + fmt.Sprintf("InfixVar(infix=%s)\n", formatName(string(e.infix)))
}

func treeLambda(e *Lambda, offset int) string {
	var b strings.Builder
	b.WriteString(indent(offset) + "Lambda()\n")
	for _, p := range e.params {
		if p != nil {
			b.WriteString(stringTree(p, offset+1))
		}
	}
	if e.return_ != nil {
		b.WriteString(stringTree(e.return_, offset+1))
	}
	if e.body != nil {
		b.WriteString(stringTree(e.body, offset+1))
	}
	return b.String()
}

func treeLet(e *Let, offset int) string {
	var b strings.Builder
	b.WriteString(indent(offset) + "Let()\n")
	if e.pattern != nil {
		b.WriteString(stringTree(e.pattern, offset+1))
	}
	if e.value != nil {
		b.WriteString(stringTree(e.value, offset+1))
	}
	if e.nested != nil {
		b.WriteString(stringTree(e.nested, offset+1))
	}
	return b.String()
}

func treeList(e *List, offset int) string {
	var b strings.Builder
	b.WriteString(indent(offset) + "List()\n")
	for _, item := range e.items {
		if item != nil {
			b.WriteString(stringTree(item, offset+1))
		}
	}
	return b.String()
}

func treeNegate(e *Negate, offset int) string {
	var b strings.Builder
	b.WriteString(indent(offset) + "Negate()\n")
	if e.nested != nil {
		b.WriteString(stringTree(e.nested, offset+1))
	}
	return b.String()
}

func treeRecord(e *Record, offset int) string {
	var b strings.Builder
	b.WriteString(indent(offset) + "Record()\n")
	for _, field := range e.fields {
		if field == nil {
			continue
		}
		b.WriteString(indent(offset+1) + fmt.Sprintf(
			"RecordField(name=%s)\n", formatName(string(field.name))))
		if field.value != nil {
			b.WriteString(stringTree(field.value, offset+2))
		}
	}
	return b.String()
}

func treeSelect(e *Select, offset int) string {
	var b strings.Builder
	b.WriteString(indent(offset) + "Select()\n")
	if e.condition != nil {
		b.WriteString(stringTree(e.condition, offset+1))
	}
	for _, c := range e.cases {
		if c == nil {
			continue
		}
		b.WriteString(indent(offset+1) + "SelectCase()\n")
		if c.pattern != nil {
			b.WriteString(stringTree(c.pattern, offset+2))
		}
		if c.body != nil {
			b.WriteString(stringTree(c.body, offset+2))
		}
	}
	return b.String()
}

func treeTuple(e *Tuple, offset int) string {
	var b strings.Builder
	b.WriteString(indent(offset) + "Tuple()\n")
	for _, item := range e.items {
		if item != nil {
			b.WriteString(stringTree(item, offset+1))
		}
	}
	return b.String()
}

func treeUpdate(e *Update, offset int) string {
	var b strings.Builder
	b.WriteString(indent(offset))
	b.WriteString(fmt.Sprintf("Update(recordName=%s)\n", formatName(string(e.recordName))))
	for _, field := range e.fields {
		if field == nil {
			continue
		}
		b.WriteString(indent(offset+1) + fmt.Sprintf(
			"RecordField(name=%s)\n", formatName(string(field.name))))
		if field.value != nil {
			b.WriteString(stringTree(field.value, offset+2))
		}
	}
	return b.String()
}

func treeVar(e *Var, offset int) string {
	return indent(offset) + fmt.Sprintf("Var(name=%s)\n", formatName(string(e.name)))
}

// ---------------------------------------------------------------- patterns --

func treePAlias(p *PAlias, offset int) string {
	var b strings.Builder
	b.WriteString(indent(offset))
	b.WriteString(fmt.Sprintf("PAlias(alias=%s)\n", formatName(string(p.alias))))
	if p.nested != nil {
		b.WriteString(stringTree(p.nested, offset+1))
	}
	if p.declaredType != nil {
		b.WriteString(stringTree(p.declaredType, offset+1))
	}
	return b.String()
}

func treePAny(p *PAny, offset int) string {
	var b strings.Builder
	b.WriteString(indent(offset) + "PAny()\n")
	if p.declaredType != nil {
		b.WriteString(stringTree(p.declaredType, offset+1))
	}
	return b.String()
}

func treePCons(p *PCons, offset int) string {
	var b strings.Builder
	b.WriteString(indent(offset) + "PCons()\n")
	if p.head != nil {
		b.WriteString(stringTree(p.head, offset+1))
	}
	if p.tail != nil {
		b.WriteString(stringTree(p.tail, offset+1))
	}
	if p.declaredType != nil {
		b.WriteString(stringTree(p.declaredType, offset+1))
	}
	return b.String()
}

func treePConst(p *PConst, offset int) string {
	var b strings.Builder
	b.WriteString(indent(offset))
	b.WriteString(fmt.Sprintf("PConst(value=%s)\n", formatConst(p.value)))
	if p.declaredType != nil {
		b.WriteString(stringTree(p.declaredType, offset+1))
	}
	return b.String()
}

func treePList(p *PList, offset int) string {
	var b strings.Builder
	b.WriteString(indent(offset) + "PList()\n")
	for _, item := range p.items {
		if item != nil {
			b.WriteString(stringTree(item, offset+1))
		}
	}
	if p.declaredType != nil {
		b.WriteString(stringTree(p.declaredType, offset+1))
	}
	return b.String()
}

func treePNamed(p *PNamed, offset int) string {
	var b strings.Builder
	b.WriteString(indent(offset))
	b.WriteString(fmt.Sprintf("PNamed(name=%s)\n", formatName(string(p.name))))
	if p.declaredType != nil {
		b.WriteString(stringTree(p.declaredType, offset+1))
	}
	return b.String()
}

func treePOption(p *POption, offset int) string {
	var b strings.Builder
	b.WriteString(indent(offset))
	b.WriteString(fmt.Sprintf("POption(name=%s)\n", formatName(string(p.name))))
	for _, a := range p.values {
		if a != nil {
			b.WriteString(stringTree(a, offset+1))
		}
	}
	if p.declaredType != nil {
		b.WriteString(stringTree(p.declaredType, offset+1))
	}
	return b.String()
}

func treePRecord(p *PRecord, offset int) string {
	var b strings.Builder
	b.WriteString(indent(offset) + "PRecord()\n")
	for _, field := range p.fields {
		if field == nil {
			continue
		}
		b.WriteString(indent(offset+1) + fmt.Sprintf(
			"PRecordField(name=%s)\n", formatName(string(field.name))))
	}
	if p.declaredType != nil {
		b.WriteString(stringTree(p.declaredType, offset+1))
	}
	return b.String()
}

func treePTuple(p *PTuple, offset int) string {
	var b strings.Builder
	b.WriteString(indent(offset) + "PTuple()\n")
	for _, item := range p.items {
		if item != nil {
			b.WriteString(stringTree(item, offset+1))
		}
	}
	if p.declaredType != nil {
		b.WriteString(stringTree(p.declaredType, offset+1))
	}
	return b.String()
}

// ------------------------------------------------------------------- types --

func treeTData(t *TData, offset int) string {
	var b strings.Builder
	b.WriteString(indent(offset))
	b.WriteString(fmt.Sprintf("TData(name=%s)\n", formatName(string(t.name))))
	for _, a := range t.args {
		if a != nil {
			b.WriteString(stringTree(a, offset+1))
		}
	}
	for _, opt := range t.options {
		if opt == nil {
			continue
		}
		b.WriteString(indent(offset+1) + fmt.Sprintf(
			"DataOption(name=%s, hidden=%s)\n",
			formatName(string(opt.name)), formatBool(opt.hidden)))
		for _, v := range opt.values {
			if v != nil {
				b.WriteString(stringTree(v, offset+2))
			}
		}
	}
	return b.String()
}

func treeTFunc(t *TFunc, offset int) string {
	var b strings.Builder
	b.WriteString(indent(offset) + "TFunc()\n")
	for _, p := range t.params {
		if p != nil {
			b.WriteString(stringTree(p, offset+1))
		}
	}
	if t.return_ != nil {
		b.WriteString(stringTree(t.return_, offset+1))
	}
	return b.String()
}

func treeTNamed(t *TNamed, offset int) string {
	var b strings.Builder
	b.WriteString(indent(offset))
	b.WriteString(fmt.Sprintf("TNamed(name=%s)\n", formatName(string(t.name))))
	for _, a := range t.args {
		if a != nil {
			b.WriteString(stringTree(a, offset+1))
		}
	}
	return b.String()
}

func treeTNative(t *TNative, offset int) string {
	var b strings.Builder
	b.WriteString(indent(offset))
	b.WriteString(fmt.Sprintf("TNative(name=%s)\n", formatName(string(t.name))))
	for _, a := range t.args {
		if a != nil {
			b.WriteString(stringTree(a, offset+1))
		}
	}
	return b.String()
}

func treeTParameter(t *TParameter, offset int) string {
	return indent(offset) + fmt.Sprintf("TParameter(name=%s)\n", formatName(string(t.name)))
}

func treeTRecord(t *TRecord, offset int) string {
	var b strings.Builder
	b.WriteString(indent(offset) + "TRecord()\n")
	// Sort keys for stable output across runs (Go map order is random).
	keys := make([]string, 0, len(t.fields))
	for k := range t.fields {
		keys = append(keys, string(k))
	}
	sort.Strings(keys)
	for _, k := range keys {
		b.WriteString(indent(offset+1) + fmt.Sprintf(
			"TRecordField(name=%s)\n", formatName(k)))
		if v := t.fields[ast.Identifier(k)]; v != nil {
			b.WriteString(stringTree(v, offset+2))
		}
	}
	return b.String()
}

func treeTTuple(t *TTuple, offset int) string {
	var b strings.Builder
	b.WriteString(indent(offset) + "TTuple()\n")
	for _, item := range t.items {
		if item != nil {
			b.WriteString(stringTree(item, offset+1))
		}
	}
	return b.String()
}

func treeTUnit(t *TUnit, offset int) string {
	return indent(offset) + "TUnit()\n"
}

// ----------------------------------------------- per-type forwarder methods --
//
// Each concrete parsed AST type satisfies Statement, which now requires a
// StringTree(offset int) string method. The implementations all delegate to
// the package-private stringTree dispatcher above so that adding a new node
// type only requires updating one switch.

func (m *Module) StringTree(offset int) string             { return stringTree(m, offset) }
func (i *import_) StringTree(offset int) string            { return stringTree(i, offset) }
func (a *alias) StringTree(offset int) string              { return stringTree(a, offset) }
func (i *infix) StringTree(offset int) string              { return stringTree(i, offset) }
func (d dataType) StringTree(offset int) string            { return stringTree(&d, offset) }
func (d *dataTypeOption) StringTree(offset int) string     { return stringTree(d, offset) }
func (d *definition) StringTree(offset int) string         { return stringTree(d, offset) }

func (e *Access) StringTree(offset int) string             { return stringTree(e, offset) }
func (e *Accessor) StringTree(offset int) string           { return stringTree(e, offset) }
func (e *Apply) StringTree(offset int) string              { return stringTree(e, offset) }
func (e *BinOp) StringTree(offset int) string              { return stringTree(e, offset) }
func (e *Call) StringTree(offset int) string               { return stringTree(e, offset) }
func (e *Const) StringTree(offset int) string              { return stringTree(e, offset) }
func (e *Constructor) StringTree(offset int) string        { return stringTree(e, offset) }
func (e *Function) StringTree(offset int) string           { return stringTree(e, offset) }
func (e *If) StringTree(offset int) string                 { return stringTree(e, offset) }
func (e *InfixVar) StringTree(offset int) string           { return stringTree(e, offset) }
func (e *Lambda) StringTree(offset int) string             { return stringTree(e, offset) }
func (e *Let) StringTree(offset int) string                { return stringTree(e, offset) }
func (e *List) StringTree(offset int) string               { return stringTree(e, offset) }
func (e *Negate) StringTree(offset int) string             { return stringTree(e, offset) }
func (e *Record) StringTree(offset int) string             { return stringTree(e, offset) }
func (e *Select) StringTree(offset int) string             { return stringTree(e, offset) }
func (e *Tuple) StringTree(offset int) string              { return stringTree(e, offset) }
func (e *Update) StringTree(offset int) string             { return stringTree(e, offset) }
func (e *Var) StringTree(offset int) string                { return stringTree(e, offset) }

func (p *PAlias) StringTree(offset int) string             { return stringTree(p, offset) }
func (p *PAny) StringTree(offset int) string               { return stringTree(p, offset) }
func (p *PCons) StringTree(offset int) string              { return stringTree(p, offset) }
func (p *PConst) StringTree(offset int) string             { return stringTree(p, offset) }
func (p *PList) StringTree(offset int) string              { return stringTree(p, offset) }
func (p *PNamed) StringTree(offset int) string             { return stringTree(p, offset) }
func (p *POption) StringTree(offset int) string            { return stringTree(p, offset) }
func (p *PRecord) StringTree(offset int) string            { return stringTree(p, offset) }
func (p *PTuple) StringTree(offset int) string             { return stringTree(p, offset) }

func (t *TData) StringTree(offset int) string              { return stringTree(t, offset) }
func (t *TFunc) StringTree(offset int) string              { return stringTree(t, offset) }
func (t *TNamed) StringTree(offset int) string             { return stringTree(t, offset) }
func (t *TNative) StringTree(offset int) string            { return stringTree(t, offset) }
func (t *TParameter) StringTree(offset int) string         { return stringTree(t, offset) }
func (t *TRecord) StringTree(offset int) string            { return stringTree(t, offset) }
func (t *TTuple) StringTree(offset int) string             { return stringTree(t, offset) }
func (t *TUnit) StringTree(offset int) string              { return stringTree(t, offset) }
