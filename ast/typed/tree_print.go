package typed

import (
	"fmt"
	"github.com/nar-lang/nar-compiler/ast"
	"sort"
	"strings"
)

// stringTree renders a typed AST node as a multi-line tree string.
// The output format mirrors compiler.ast.typed.tree_print (Lua) byte-for-byte:
// tab-indented, one line per node, headers `Ty*`/`T*`.
//
// Each expression/pattern node emits its inferred `type_` as a trailing
// `Type` child for full visibility into inference. Back-pointers
// (Constructor.dataType, Global/POption.definition, Local/Update.target) are
// NOT recursed — only identifying fields appear on the header.
//
// TUnbound is rendered inline; its `predecessor` chain is NOT followed
// (it can be cyclic). Type subtrees are cycle-guarded by table identity so
// recursive ADTs do not blow the stack.
func stringTree(s any, offset int, visited map[any]bool) string {
	if s == nil {
		return ""
	}
	switch n := s.(type) {
	case *Module:
		return treeModule(n, offset, visited)
	case *Definition:
		return treeDefinition(n, offset, visited)

	case *Access:
		return treeAccess(n, offset, visited)
	case *Apply:
		return treeApply(n, offset, visited)
	case *Call:
		return treeCall(n, offset, visited)
	case *Const:
		return treeConst(n, offset, visited)
	case *Constructor:
		return treeConstructor(n, offset, visited)
	case *Global:
		return treeGlobal(n, offset, visited)
	case *Let:
		return treeLet(n, offset, visited)
	case *List:
		return treeList(n, offset, visited)
	case *Local:
		return treeLocal(n, offset, visited)
	case *Record:
		return treeRecord(n, offset, visited)
	case *Select:
		return treeSelect(n, offset, visited)
	case *Tuple:
		return treeTuple(n, offset, visited)
	case *Update:
		return treeUpdate(n, offset, visited)

	case *PAlias:
		return treePAlias(n, offset, visited)
	case *PAny:
		return treePAny(n, offset, visited)
	case *PCons:
		return treePCons(n, offset, visited)
	case *PConst:
		return treePConst(n, offset, visited)
	case *PList:
		return treePList(n, offset, visited)
	case *PNamed:
		return treePNamed(n, offset, visited)
	case *POption:
		return treePOption(n, offset, visited)
	case *PRecord:
		return treePRecord(n, offset, visited)
	case *PTuple:
		return treePTuple(n, offset, visited)

	case *TData:
		return treeTypeGuarded(n, offset, visited, treeTData)
	case *TFunc:
		return treeTypeGuarded(n, offset, visited, treeTFunc)
	case *TNative:
		return treeTypeGuarded(n, offset, visited, treeTNative)
	case *TRecord:
		return treeTypeGuarded(n, offset, visited, treeTRecord)
	case *TTuple:
		return treeTypeGuarded(n, offset, visited, treeTTuple)
	case *TUnbound:
		return treeTUnbound(n, offset)
	}
	return indent(offset) + fmt.Sprintf("?(%T)\n", s)
}

// treeTypeGuarded enforces the per-traversal cycle guard for any type node.
func treeTypeGuarded[T any](n *T, offset int, visited map[any]bool,
	render func(*T, int, map[any]bool) string) string {
	key := any(n)
	if visited[key] {
		return indent(offset) + fmt.Sprintf("<recursive %s>\n", typeNameOf(n))
	}
	visited[key] = true
	out := render(n, offset, visited)
	delete(visited, key)
	return out
}

// typeNameOf returns the short typed-AST kind name (e.g. "TData", "TFunc")
// used in the cycle-guard marker, matching the Lua side which prints
// `node.kind`.
func typeNameOf(n any) string {
	switch n.(type) {
	case *TData:
		return "TData"
	case *TFunc:
		return "TFunc"
	case *TNative:
		return "TNative"
	case *TRecord:
		return "TRecord"
	case *TTuple:
		return "TTuple"
	case *TUnbound:
		return "TUnbound"
	}
	return fmt.Sprintf("%T", n)
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
		return fmt.Sprintf("CString(%s)", strconvQuote(c.Value))
	case ast.CChar:
		return fmt.Sprintf("CChar(%s)", strconvQuote(string(c.Value)))
	case ast.CUnit:
		return "CUnit"
	}
	return fmt.Sprintf("%T", v)
}

func strconvFloat(f float64) string {
	return fmt.Sprintf("%g", f)
}

func strconvQuote(s string) string {
	var b strings.Builder
	b.WriteByte('"')
	for i := 0; i < len(s); i++ {
		ch := s[i]
		switch ch {
		case '\\':
			b.WriteString("\\\\")
		case '"':
			b.WriteString("\\\"")
		case '\a':
			b.WriteString("\\a")
		case '\b':
			b.WriteString("\\b")
		case '\f':
			b.WriteString("\\f")
		case '\n':
			b.WriteString("\\n")
		case '\r':
			b.WriteString("\\r")
		case '\t':
			b.WriteString("\\t")
		case '\v':
			b.WriteString("\\v")
		default:
			if ch < 0x20 || ch == 0x7F {
				b.WriteString(fmt.Sprintf("\\x%02x", ch))
			} else {
				b.WriteByte(ch)
			}
		}
	}
	b.WriteByte('"')
	return b.String()
}

func appendType(b *strings.Builder, t Type, offset int, visited map[any]bool) {
	if t == nil {
		return
	}
	b.WriteString(indent(offset) + "Type\n")
	b.WriteString(stringTree(t, offset+1, visited))
}

func appendDeclaredType(b *strings.Builder, t Type, offset int, visited map[any]bool) {
	if t == nil {
		return
	}
	b.WriteString(indent(offset) + "DeclaredType\n")
	b.WriteString(stringTree(t, offset+1, visited))
}

// ---------------------------------------------------------- module-level --

func treeModule(m *Module, offset int, visited map[any]bool) string {
	var b strings.Builder
	b.WriteString(indent(offset))
	b.WriteString(fmt.Sprintf("TyModule(name=%s)\n", formatName(string(m.name))))
	depKeys := make([]string, 0, len(m.dependencies))
	for k := range m.dependencies {
		depKeys = append(depKeys, string(k))
	}
	sort.Strings(depKeys)
	for _, k := range depKeys {
		idents := m.dependencies[ast.QualifiedIdentifier(k)]
		sorted := make([]ast.Identifier, len(idents))
		copy(sorted, idents)
		sort.Slice(sorted, func(i, j int) bool { return sorted[i] < sorted[j] })
		b.WriteString(indent(offset+1) + fmt.Sprintf(
			"TyDependency(module=%s, names=%s)\n", formatName(k), formatIdentList(sorted)))
	}
	for _, d := range m.definitions {
		if d != nil {
			b.WriteString(stringTree(d, offset+1, visited))
		}
	}
	return b.String()
}

func treeDefinition(d *Definition, offset int, visited map[any]bool) string {
	var b strings.Builder
	b.WriteString(indent(offset))
	b.WriteString(fmt.Sprintf("TyDefinition(name=%s, hidden=%s)\n",
		formatName(string(d.name)), formatBool(d.hidden)))
	for _, p := range d.params {
		if p != nil {
			b.WriteString(stringTree(p, offset+1, visited))
		}
	}
	if d.body != nil {
		b.WriteString(stringTree(d.body, offset+1, visited))
	}
	appendDeclaredType(&b, d.declaredType, offset+1, visited)
	appendType(&b, d.type_, offset+1, visited)
	return b.String()
}

// -------------------------------------------------------------- expressions --

func treeAccess(e *Access, offset int, visited map[any]bool) string {
	var b strings.Builder
	b.WriteString(indent(offset))
	b.WriteString(fmt.Sprintf("TyAccess(fieldName=%s)\n", formatName(string(e.fieldName))))
	if e.record != nil {
		b.WriteString(stringTree(e.record, offset+1, visited))
	}
	appendType(&b, e.type_, offset+1, visited)
	return b.String()
}

func treeApply(e *Apply, offset int, visited map[any]bool) string {
	var b strings.Builder
	b.WriteString(indent(offset) + "TyApply()\n")
	if e.func_ != nil {
		b.WriteString(stringTree(e.func_, offset+1, visited))
	}
	for _, a := range e.args {
		if a != nil {
			b.WriteString(stringTree(a, offset+1, visited))
		}
	}
	appendType(&b, e.type_, offset+1, visited)
	return b.String()
}

func treeCall(e *Call, offset int, visited map[any]bool) string {
	var b strings.Builder
	b.WriteString(indent(offset))
	b.WriteString(fmt.Sprintf("TyCall(name=%s)\n", formatName(string(e.name))))
	for _, a := range e.args {
		if a != nil {
			b.WriteString(stringTree(a, offset+1, visited))
		}
	}
	appendType(&b, e.type_, offset+1, visited)
	return b.String()
}

func treeConst(e *Const, offset int, visited map[any]bool) string {
	var b strings.Builder
	b.WriteString(indent(offset))
	b.WriteString(fmt.Sprintf("TyConst(value=%s)\n", formatConst(e.value)))
	appendType(&b, e.type_, offset+1, visited)
	return b.String()
}

func treeConstructor(e *Constructor, offset int, visited map[any]bool) string {
	var b strings.Builder
	b.WriteString(indent(offset))
	b.WriteString(fmt.Sprintf(
		"TyConstructor(dataName=%s, optionName=%s)\n",
		formatName(string(e.dataName)), formatName(string(e.optionName))))
	for _, a := range e.args {
		if a != nil {
			b.WriteString(stringTree(a, offset+1, visited))
		}
	}
	appendType(&b, e.type_, offset+1, visited)
	return b.String()
}

func treeGlobal(e *Global, offset int, visited map[any]bool) string {
	var b strings.Builder
	b.WriteString(indent(offset))
	b.WriteString(fmt.Sprintf(
		"TyGlobal(moduleName=%s, definitionName=%s)\n",
		formatName(string(e.moduleName)), formatName(string(e.definitionName))))
	appendType(&b, e.type_, offset+1, visited)
	return b.String()
}

func treeLet(e *Let, offset int, visited map[any]bool) string {
	var b strings.Builder
	b.WriteString(indent(offset) + "TyLet()\n")
	if e.pattern != nil {
		b.WriteString(stringTree(e.pattern, offset+1, visited))
	}
	if e.value != nil {
		b.WriteString(stringTree(e.value, offset+1, visited))
	}
	if e.body != nil {
		b.WriteString(stringTree(e.body, offset+1, visited))
	}
	appendType(&b, e.type_, offset+1, visited)
	return b.String()
}

func treeList(e *List, offset int, visited map[any]bool) string {
	var b strings.Builder
	b.WriteString(indent(offset) + "TyList()\n")
	for _, item := range e.items {
		if item != nil {
			b.WriteString(stringTree(item, offset+1, visited))
		}
	}
	appendType(&b, e.type_, offset+1, visited)
	return b.String()
}

func treeLocal(e *Local, offset int, visited map[any]bool) string {
	var b strings.Builder
	b.WriteString(indent(offset))
	b.WriteString(fmt.Sprintf("TyLocal(name=%s)\n", formatName(string(e.name))))
	appendType(&b, e.type_, offset+1, visited)
	return b.String()
}

func treeRecord(e *Record, offset int, visited map[any]bool) string {
	var b strings.Builder
	b.WriteString(indent(offset) + "TyRecord()\n")
	for _, field := range e.fields {
		if field == nil {
			continue
		}
		b.WriteString(indent(offset+1) + fmt.Sprintf(
			"TyRecordField(name=%s)\n", formatName(string(field.name))))
		if field.value != nil {
			b.WriteString(stringTree(field.value, offset+2, visited))
		}
	}
	appendType(&b, e.type_, offset+1, visited)
	return b.String()
}

func treeSelect(e *Select, offset int, visited map[any]bool) string {
	var b strings.Builder
	b.WriteString(indent(offset) + "TySelect()\n")
	if e.condition != nil {
		b.WriteString(stringTree(e.condition, offset+1, visited))
	}
	for _, c := range e.cases {
		if c == nil {
			continue
		}
		b.WriteString(indent(offset+1) + "TySelectCase()\n")
		if c.pattern != nil {
			b.WriteString(stringTree(c.pattern, offset+2, visited))
		}
		if c.expression != nil {
			b.WriteString(stringTree(c.expression, offset+2, visited))
		}
	}
	appendType(&b, e.type_, offset+1, visited)
	return b.String()
}

func treeTuple(e *Tuple, offset int, visited map[any]bool) string {
	var b strings.Builder
	b.WriteString(indent(offset) + "TyTuple()\n")
	for _, item := range e.items {
		if item != nil {
			b.WriteString(stringTree(item, offset+1, visited))
		}
	}
	appendType(&b, e.type_, offset+1, visited)
	return b.String()
}

func treeUpdate(e *Update, offset int, visited map[any]bool) string {
	var b strings.Builder
	moduleStr := "<nil>"
	if e.moduleName != "" {
		moduleStr = string(e.moduleName)
	}
	b.WriteString(indent(offset))
	b.WriteString(fmt.Sprintf(
		"TyUpdate(moduleName=%s, recordName=%s)\n",
		moduleStr, formatName(string(e.recordName))))
	for _, field := range e.fields {
		if field == nil {
			continue
		}
		b.WriteString(indent(offset+1) + fmt.Sprintf(
			"TyRecordField(name=%s)\n", formatName(string(field.name))))
		if field.value != nil {
			b.WriteString(stringTree(field.value, offset+2, visited))
		}
	}
	appendType(&b, e.type_, offset+1, visited)
	return b.String()
}

// ---------------------------------------------------------------- patterns --

func treePAlias(p *PAlias, offset int, visited map[any]bool) string {
	var b strings.Builder
	b.WriteString(indent(offset))
	b.WriteString(fmt.Sprintf("TyPAlias(alias=%s)\n", formatName(string(p.alias))))
	if p.nested != nil {
		b.WriteString(stringTree(p.nested, offset+1, visited))
	}
	appendDeclaredType(&b, p.declaredType, offset+1, visited)
	appendType(&b, p.type_, offset+1, visited)
	return b.String()
}

func treePAny(p *PAny, offset int, visited map[any]bool) string {
	var b strings.Builder
	b.WriteString(indent(offset) + "TyPAny()\n")
	appendDeclaredType(&b, p.declaredType, offset+1, visited)
	appendType(&b, p.type_, offset+1, visited)
	return b.String()
}

func treePCons(p *PCons, offset int, visited map[any]bool) string {
	var b strings.Builder
	b.WriteString(indent(offset) + "TyPCons()\n")
	if p.head != nil {
		b.WriteString(stringTree(p.head, offset+1, visited))
	}
	if p.tail != nil {
		b.WriteString(stringTree(p.tail, offset+1, visited))
	}
	appendDeclaredType(&b, p.declaredType, offset+1, visited)
	appendType(&b, p.type_, offset+1, visited)
	return b.String()
}

func treePConst(p *PConst, offset int, visited map[any]bool) string {
	var b strings.Builder
	b.WriteString(indent(offset))
	b.WriteString(fmt.Sprintf("TyPConst(value=%s)\n", formatConst(p.value)))
	appendDeclaredType(&b, p.declaredType, offset+1, visited)
	appendType(&b, p.type_, offset+1, visited)
	return b.String()
}

func treePList(p *PList, offset int, visited map[any]bool) string {
	var b strings.Builder
	b.WriteString(indent(offset) + "TyPList()\n")
	for _, item := range p.items {
		if item != nil {
			b.WriteString(stringTree(item, offset+1, visited))
		}
	}
	appendDeclaredType(&b, p.declaredType, offset+1, visited)
	appendType(&b, p.type_, offset+1, visited)
	return b.String()
}

func treePNamed(p *PNamed, offset int, visited map[any]bool) string {
	var b strings.Builder
	b.WriteString(indent(offset))
	b.WriteString(fmt.Sprintf("TyPNamed(name=%s)\n", formatName(string(p.name))))
	appendDeclaredType(&b, p.declaredType, offset+1, visited)
	appendType(&b, p.type_, offset+1, visited)
	return b.String()
}

func treePOption(p *POption, offset int, visited map[any]bool) string {
	defName := "<nil>"
	if p.definition != nil {
		defName = formatName(string(p.definition.name))
	}
	var b strings.Builder
	b.WriteString(indent(offset))
	b.WriteString(fmt.Sprintf("TyPOption(definitionName=%s)\n", defName))
	for _, v := range p.args {
		if v != nil {
			b.WriteString(stringTree(v, offset+1, visited))
		}
	}
	appendDeclaredType(&b, p.declaredType, offset+1, visited)
	appendType(&b, p.type_, offset+1, visited)
	return b.String()
}

func treePRecord(p *PRecord, offset int, visited map[any]bool) string {
	var b strings.Builder
	b.WriteString(indent(offset) + "TyPRecord()\n")
	for _, field := range p.fields {
		if field == nil {
			continue
		}
		b.WriteString(indent(offset+1) + fmt.Sprintf(
			"TyPRecordField(name=%s)\n", formatName(string(field.name))))
	}
	appendDeclaredType(&b, p.declaredType, offset+1, visited)
	appendType(&b, p.type_, offset+1, visited)
	return b.String()
}

func treePTuple(p *PTuple, offset int, visited map[any]bool) string {
	var b strings.Builder
	b.WriteString(indent(offset) + "TyPTuple()\n")
	for _, item := range p.items {
		if item != nil {
			b.WriteString(stringTree(item, offset+1, visited))
		}
	}
	appendDeclaredType(&b, p.declaredType, offset+1, visited)
	appendType(&b, p.type_, offset+1, visited)
	return b.String()
}

// ------------------------------------------------------------------- types --

func treeTData(t *TData, offset int, visited map[any]bool) string {
	var b strings.Builder
	b.WriteString(indent(offset))
	b.WriteString(fmt.Sprintf("TData(name=%s)\n", formatName(string(t.name))))
	for _, a := range t.args {
		if a != nil {
			b.WriteString(stringTree(a, offset+1, visited))
		}
	}
	for _, opt := range t.options {
		if opt == nil {
			continue
		}
		b.WriteString(indent(offset+1) + fmt.Sprintf(
			"DataOption(name=%s)\n", formatName(string(opt.name))))
		for _, v := range opt.values {
			if v != nil {
				b.WriteString(stringTree(v, offset+2, visited))
			}
		}
	}
	return b.String()
}

func treeTFunc(t *TFunc, offset int, visited map[any]bool) string {
	var b strings.Builder
	b.WriteString(indent(offset) + "TFunc()\n")
	for _, p := range t.params {
		if p != nil {
			b.WriteString(stringTree(p, offset+1, visited))
		}
	}
	if t.return_ != nil {
		b.WriteString(stringTree(t.return_, offset+1, visited))
	}
	return b.String()
}

func treeTNative(t *TNative, offset int, visited map[any]bool) string {
	var b strings.Builder
	b.WriteString(indent(offset))
	b.WriteString(fmt.Sprintf("TNative(name=%s)\n", formatName(string(t.name))))
	for _, a := range t.args {
		if a != nil {
			b.WriteString(stringTree(a, offset+1, visited))
		}
	}
	return b.String()
}

func treeTRecord(t *TRecord, offset int, visited map[any]bool) string {
	var b strings.Builder
	b.WriteString(indent(offset))
	b.WriteString(fmt.Sprintf(
		"TRecord(mayHaveMoreFields=%s)\n", formatBool(t.mayHaveMoreFields)))
	keys := make([]string, 0, len(t.fields))
	for k := range t.fields {
		keys = append(keys, string(k))
	}
	sort.Strings(keys)
	for _, k := range keys {
		b.WriteString(indent(offset+1) + fmt.Sprintf(
			"TRecordField(name=%s)\n", formatName(k)))
		if v := t.fields[ast.Identifier(k)]; v != nil {
			b.WriteString(stringTree(v, offset+2, visited))
		}
	}
	return b.String()
}

func treeTTuple(t *TTuple, offset int, visited map[any]bool) string {
	var b strings.Builder
	b.WriteString(indent(offset) + "TTuple()\n")
	for _, item := range t.items {
		if item != nil {
			b.WriteString(stringTree(item, offset+1, visited))
		}
	}
	return b.String()
}

func treeTUnbound(t *TUnbound, offset int) string {
	constraintName := ""
	if t.constraint != "" {
		constraintName = string(t.constraint)
	}
	return indent(offset) + fmt.Sprintf(
		"TUnbound(index=%d, constraint=%s, name=%s)\n",
		t.index, formatName(constraintName), formatName(string(t.givenName)))
}

// ----------------------------------------------- per-type forwarder methods --

func (m *Module) StringTree(offset int) string      { return stringTree(m, offset, map[any]bool{}) }
func (d *Definition) StringTree(offset int) string  { return stringTree(d, offset, map[any]bool{}) }

func (e *Access) StringTree(offset int) string      { return stringTree(e, offset, map[any]bool{}) }
func (e *Apply) StringTree(offset int) string       { return stringTree(e, offset, map[any]bool{}) }
func (e *Call) StringTree(offset int) string        { return stringTree(e, offset, map[any]bool{}) }
func (e *Const) StringTree(offset int) string       { return stringTree(e, offset, map[any]bool{}) }
func (e *Constructor) StringTree(offset int) string { return stringTree(e, offset, map[any]bool{}) }
func (e *Global) StringTree(offset int) string      { return stringTree(e, offset, map[any]bool{}) }
func (e *Let) StringTree(offset int) string         { return stringTree(e, offset, map[any]bool{}) }
func (e *List) StringTree(offset int) string        { return stringTree(e, offset, map[any]bool{}) }
func (e *Local) StringTree(offset int) string       { return stringTree(e, offset, map[any]bool{}) }
func (e *Record) StringTree(offset int) string      { return stringTree(e, offset, map[any]bool{}) }
func (e *Select) StringTree(offset int) string      { return stringTree(e, offset, map[any]bool{}) }
func (e *Tuple) StringTree(offset int) string       { return stringTree(e, offset, map[any]bool{}) }
func (e *Update) StringTree(offset int) string      { return stringTree(e, offset, map[any]bool{}) }

func (p *PAlias) StringTree(offset int) string  { return stringTree(p, offset, map[any]bool{}) }
func (p *PAny) StringTree(offset int) string    { return stringTree(p, offset, map[any]bool{}) }
func (p *PCons) StringTree(offset int) string   { return stringTree(p, offset, map[any]bool{}) }
func (p *PConst) StringTree(offset int) string  { return stringTree(p, offset, map[any]bool{}) }
func (p *PList) StringTree(offset int) string   { return stringTree(p, offset, map[any]bool{}) }
func (p *PNamed) StringTree(offset int) string  { return stringTree(p, offset, map[any]bool{}) }
func (p *POption) StringTree(offset int) string { return stringTree(p, offset, map[any]bool{}) }
func (p *PRecord) StringTree(offset int) string { return stringTree(p, offset, map[any]bool{}) }
func (p *PTuple) StringTree(offset int) string  { return stringTree(p, offset, map[any]bool{}) }

func (t *TData) StringTree(offset int) string    { return stringTree(t, offset, map[any]bool{}) }
func (t *TFunc) StringTree(offset int) string    { return stringTree(t, offset, map[any]bool{}) }
func (t *TNative) StringTree(offset int) string  { return stringTree(t, offset, map[any]bool{}) }
func (t *TRecord) StringTree(offset int) string  { return stringTree(t, offset, map[any]bool{}) }
func (t *TTuple) StringTree(offset int) string   { return stringTree(t, offset, map[any]bool{}) }
func (t *TUnbound) StringTree(offset int) string { return stringTree(t, offset, map[any]bool{}) }
