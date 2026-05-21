package normalized

import (
	"fmt"
	"github.com/nar-lang/nar-compiler/ast"
	"sort"
	"strings"
)

// stringTree renders a normalized AST node as a multi-line tree string.
// The output format mirrors compiler.ast.normalized.tree_print (Lua) byte-for-byte:
// node headers are prefixed with `N`, children rendered with `offset + 1`
// tabs, and one line per node.
func stringTree(s any, offset int) string {
	if s == nil {
		return ""
	}
	switch n := s.(type) {
	case *Module:
		return treeModule(n, offset)
	case *definition:
		return treeDefinition(n, offset)

	case *Access:
		return treeAccess(n, offset)
	case *Apply:
		return treeApply(n, offset)
	case *Call:
		return treeCall(n, offset)
	case *Const:
		return treeConst(n, offset)
	case *Constructor:
		return treeConstructor(n, offset)
	case *Function:
		return treeFunction(n, offset)
	case *Global:
		return treeGlobal(n, offset)
	case *Lambda:
		return treeLambda(n, offset)
	case *Let:
		return treeLet(n, offset)
	case *List:
		return treeList(n, offset)
	case *Local:
		return treeLocal(n, offset)
	case *Record:
		return treeRecord(n, offset)
	case *Select:
		return treeSelect(n, offset)
	case *Tuple:
		return treeTuple(n, offset)
	case *Update:
		return treeUpdate(n, offset)

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
	case *TNative:
		return treeTNative(n, offset)
	case *TParameter:
		return treeTParameter(n, offset)
	case *TPlaceholder:
		return treeTPlaceholder(n, offset)
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

// strconvQuote matches Go's strconv.Quote behavior used by Lua's quoteString
// helper, restricted to ASCII escapes that the Lua side emits.
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

// ---------------------------------------------------------- module-level --

func treeModule(m *Module, offset int) string {
	var b strings.Builder
	b.WriteString(indent(offset))
	b.WriteString(fmt.Sprintf("NModule(name=%s)\n", formatName(string(m.name))))
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
			"NDependency(module=%s, names=%s)\n", formatName(k), formatIdentList(sorted)))
	}
	for _, d := range m.definitions {
		if d != nil {
			b.WriteString(stringTree(d, offset+1))
		}
	}
	return b.String()
}

func treeDefinition(d *definition, offset int) string {
	var b strings.Builder
	b.WriteString(indent(offset))
	b.WriteString(fmt.Sprintf("NDefinition(name=%s, hidden=%s)\n",
		formatName(string(d.name_)), formatBool(d.hidden)))
	for _, p := range d.params_ {
		if p != nil {
			b.WriteString(stringTree(p, offset+1))
		}
	}
	if d.body_ != nil {
		b.WriteString(stringTree(d.body_, offset+1))
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
	b.WriteString(fmt.Sprintf("NAccess(fieldName=%s)\n", formatName(string(e.fieldName))))
	if e.record != nil {
		b.WriteString(stringTree(e.record, offset+1))
	}
	return b.String()
}

func treeApply(e *Apply, offset int) string {
	var b strings.Builder
	b.WriteString(indent(offset) + "NApply()\n")
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

func treeCall(e *Call, offset int) string {
	var b strings.Builder
	b.WriteString(indent(offset))
	b.WriteString(fmt.Sprintf("NCall(name=%s)\n", formatName(string(e.name))))
	for _, a := range e.args {
		if a != nil {
			b.WriteString(stringTree(a, offset+1))
		}
	}
	return b.String()
}

func treeConst(e *Const, offset int) string {
	return indent(offset) + fmt.Sprintf("NConst(value=%s)\n", formatConst(e.value))
}

func treeConstructor(e *Constructor, offset int) string {
	var b strings.Builder
	b.WriteString(indent(offset))
	b.WriteString(fmt.Sprintf(
		"NConstructor(moduleName=%s, dataName=%s, optionName=%s)\n",
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
	b.WriteString(fmt.Sprintf("NFunction(name=%s)\n", formatName(string(e.name))))
	for _, p := range e.params {
		if p != nil {
			b.WriteString(stringTree(p, offset+1))
		}
	}
	if e.body != nil {
		b.WriteString(stringTree(e.body, offset+1))
	}
	if e.fnType != nil {
		b.WriteString(stringTree(e.fnType, offset+1))
	}
	if e.nested != nil {
		b.WriteString(stringTree(e.nested, offset+1))
	}
	return b.String()
}

func treeGlobal(e *Global, offset int) string {
	return indent(offset) + fmt.Sprintf(
		"NGlobal(moduleName=%s, definitionName=%s)\n",
		formatName(string(e.moduleName)), formatName(string(e.definitionName)))
}

func treeLambda(e *Lambda, offset int) string {
	var b strings.Builder
	b.WriteString(indent(offset) + "NLambda()\n")
	for _, p := range e.params {
		if p != nil {
			b.WriteString(stringTree(p, offset+1))
		}
	}
	if e.body != nil {
		b.WriteString(stringTree(e.body, offset+1))
	}
	return b.String()
}

func treeLet(e *Let, offset int) string {
	var b strings.Builder
	b.WriteString(indent(offset) + "NLet()\n")
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
	b.WriteString(indent(offset) + "NList()\n")
	for _, item := range e.items {
		if item != nil {
			b.WriteString(stringTree(item, offset+1))
		}
	}
	return b.String()
}

func treeLocal(e *Local, offset int) string {
	return indent(offset) + fmt.Sprintf("NLocal(name=%s)\n", formatName(string(e.name)))
}

func treeRecord(e *Record, offset int) string {
	var b strings.Builder
	b.WriteString(indent(offset) + "NRecord()\n")
	for _, field := range e.fields {
		if field == nil {
			continue
		}
		b.WriteString(indent(offset+1) + fmt.Sprintf(
			"NRecordField(name=%s)\n", formatName(string(field.name))))
		if field.value != nil {
			b.WriteString(stringTree(field.value, offset+2))
		}
	}
	return b.String()
}

func treeSelect(e *Select, offset int) string {
	var b strings.Builder
	b.WriteString(indent(offset) + "NSelect()\n")
	if e.condition != nil {
		b.WriteString(stringTree(e.condition, offset+1))
	}
	for _, c := range e.cases {
		if c == nil {
			continue
		}
		b.WriteString(indent(offset+1) + "NSelectCase()\n")
		if c.pattern != nil {
			b.WriteString(stringTree(c.pattern, offset+2))
		}
		if c.expression != nil {
			b.WriteString(stringTree(c.expression, offset+2))
		}
	}
	return b.String()
}

func treeTuple(e *Tuple, offset int) string {
	var b strings.Builder
	b.WriteString(indent(offset) + "NTuple()\n")
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
	moduleStr := "<nil>"
	if e.moduleName != "" {
		moduleStr = string(e.moduleName)
	}
	b.WriteString(fmt.Sprintf(
		"NUpdate(moduleName=%s, recordName=%s)\n",
		moduleStr, formatName(string(e.recordName))))
	for _, field := range e.fields {
		if field == nil {
			continue
		}
		b.WriteString(indent(offset+1) + fmt.Sprintf(
			"NRecordField(name=%s)\n", formatName(string(field.name))))
		if field.value != nil {
			b.WriteString(stringTree(field.value, offset+2))
		}
	}
	return b.String()
}

// ---------------------------------------------------------------- patterns --

func treePAlias(p *PAlias, offset int) string {
	var b strings.Builder
	b.WriteString(indent(offset))
	b.WriteString(fmt.Sprintf("NPAlias(alias=%s)\n", formatName(string(p.alias))))
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
	b.WriteString(indent(offset) + "NPAny()\n")
	if p.declaredType != nil {
		b.WriteString(stringTree(p.declaredType, offset+1))
	}
	return b.String()
}

func treePCons(p *PCons, offset int) string {
	var b strings.Builder
	b.WriteString(indent(offset) + "NPCons()\n")
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
	b.WriteString(fmt.Sprintf("NPConst(value=%s)\n", formatConst(p.value)))
	if p.declaredType != nil {
		b.WriteString(stringTree(p.declaredType, offset+1))
	}
	return b.String()
}

func treePList(p *PList, offset int) string {
	var b strings.Builder
	b.WriteString(indent(offset) + "NPList()\n")
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
	b.WriteString(fmt.Sprintf("NPNamed(name=%s)\n", formatName(string(p.name))))
	if p.declaredType != nil {
		b.WriteString(stringTree(p.declaredType, offset+1))
	}
	return b.String()
}

func treePOption(p *POption, offset int) string {
	var b strings.Builder
	b.WriteString(indent(offset))
	b.WriteString(fmt.Sprintf(
		"NPOption(moduleName=%s, definitionName=%s)\n",
		formatName(string(p.moduleName)), formatName(string(p.definitionName))))
	for _, v := range p.values {
		if v != nil {
			b.WriteString(stringTree(v, offset+1))
		}
	}
	if p.declaredType != nil {
		b.WriteString(stringTree(p.declaredType, offset+1))
	}
	return b.String()
}

func treePRecord(p *PRecord, offset int) string {
	var b strings.Builder
	b.WriteString(indent(offset) + "NPRecord()\n")
	for _, field := range p.fields {
		if field == nil {
			continue
		}
		b.WriteString(indent(offset+1) + fmt.Sprintf(
			"NPRecordField(name=%s)\n", formatName(string(field.name))))
	}
	if p.declaredType != nil {
		b.WriteString(stringTree(p.declaredType, offset+1))
	}
	return b.String()
}

func treePTuple(p *PTuple, offset int) string {
	var b strings.Builder
	b.WriteString(indent(offset) + "NPTuple()\n")
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
	b.WriteString(fmt.Sprintf("NTData(name=%s)\n", formatName(string(t.name))))
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
			"NDataOption(name=%s, hidden=%s)\n",
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
	b.WriteString(indent(offset) + "NTFunc()\n")
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

func treeTNative(t *TNative, offset int) string {
	var b strings.Builder
	b.WriteString(indent(offset))
	b.WriteString(fmt.Sprintf("NTNative(name=%s)\n", formatName(string(t.name))))
	for _, a := range t.args {
		if a != nil {
			b.WriteString(stringTree(a, offset+1))
		}
	}
	return b.String()
}

func treeTParameter(t *TParameter, offset int) string {
	return indent(offset) + fmt.Sprintf("NTParameter(name=%s)\n", formatName(string(t.name)))
}

func treeTPlaceholder(t *TPlaceholder, offset int) string {
	return indent(offset) + fmt.Sprintf("NTPlaceholder(name=%s)\n", formatName(string(t.name)))
}

func treeTRecord(t *TRecord, offset int) string {
	var b strings.Builder
	b.WriteString(indent(offset) + "NTRecord()\n")
	keys := make([]string, 0, len(t.fields))
	for k := range t.fields {
		keys = append(keys, string(k))
	}
	sort.Strings(keys)
	for _, k := range keys {
		b.WriteString(indent(offset+1) + fmt.Sprintf(
			"NTRecordField(name=%s)\n", formatName(k)))
		if v := t.fields[ast.Identifier(k)]; v != nil {
			b.WriteString(stringTree(v, offset+2))
		}
	}
	return b.String()
}

func treeTTuple(t *TTuple, offset int) string {
	var b strings.Builder
	b.WriteString(indent(offset) + "NTTuple()\n")
	for _, item := range t.items {
		if item != nil {
			b.WriteString(stringTree(item, offset+1))
		}
	}
	return b.String()
}

func treeTUnit(t *TUnit, offset int) string {
	return indent(offset) + "NTUnit()\n"
}

// ----------------------------------------------- per-type forwarder methods --

func (m *Module) StringTree(offset int) string       { return stringTree(m, offset) }
func (d *definition) StringTree(offset int) string   { return stringTree(d, offset) }

func (e *Access) StringTree(offset int) string       { return stringTree(e, offset) }
func (e *Apply) StringTree(offset int) string        { return stringTree(e, offset) }
func (e *Call) StringTree(offset int) string         { return stringTree(e, offset) }
func (e *Const) StringTree(offset int) string        { return stringTree(e, offset) }
func (e *Constructor) StringTree(offset int) string  { return stringTree(e, offset) }
func (e *Function) StringTree(offset int) string     { return stringTree(e, offset) }
func (e *Global) StringTree(offset int) string       { return stringTree(e, offset) }
func (e *Lambda) StringTree(offset int) string       { return stringTree(e, offset) }
func (e *Let) StringTree(offset int) string          { return stringTree(e, offset) }
func (e *List) StringTree(offset int) string         { return stringTree(e, offset) }
func (e *Local) StringTree(offset int) string        { return stringTree(e, offset) }
func (e *Record) StringTree(offset int) string       { return stringTree(e, offset) }
func (e *Select) StringTree(offset int) string       { return stringTree(e, offset) }
func (e *Tuple) StringTree(offset int) string        { return stringTree(e, offset) }
func (e *Update) StringTree(offset int) string       { return stringTree(e, offset) }

func (p *PAlias) StringTree(offset int) string       { return stringTree(p, offset) }
func (p *PAny) StringTree(offset int) string         { return stringTree(p, offset) }
func (p *PCons) StringTree(offset int) string        { return stringTree(p, offset) }
func (p *PConst) StringTree(offset int) string       { return stringTree(p, offset) }
func (p *PList) StringTree(offset int) string        { return stringTree(p, offset) }
func (p *PNamed) StringTree(offset int) string       { return stringTree(p, offset) }
func (p *POption) StringTree(offset int) string      { return stringTree(p, offset) }
func (p *PRecord) StringTree(offset int) string      { return stringTree(p, offset) }
func (p *PTuple) StringTree(offset int) string       { return stringTree(p, offset) }

func (t *TData) StringTree(offset int) string        { return stringTree(t, offset) }
func (t *TFunc) StringTree(offset int) string        { return stringTree(t, offset) }
func (t *TNative) StringTree(offset int) string      { return stringTree(t, offset) }
func (t *TParameter) StringTree(offset int) string   { return stringTree(t, offset) }
func (t *TPlaceholder) StringTree(offset int) string { return stringTree(t, offset) }
func (t *TRecord) StringTree(offset int) string      { return stringTree(t, offset) }
func (t *TTuple) StringTree(offset int) string       { return stringTree(t, offset) }
func (t *TUnit) StringTree(offset int) string        { return stringTree(t, offset) }
