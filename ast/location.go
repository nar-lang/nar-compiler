package ast

import (
	"fmt"
	"github.com/nar-lang/nar-compiler/bytecode"
)

type Location struct {
	filePath    string
	fileContent []rune
	start       uint32
	end         uint32
}

func NewLocation(filePath string, content []rune, start uint32, end uint32) Location {
	return Location{
		filePath:    filePath,
		fileContent: content,
		start:       start,
		end:         end,
	}
}

func NewLocationCursor(filePath string, content []rune, start uint32) Location {
	return NewLocation(filePath, content, start, start)
}

func NewLocationSrc(filePath string, content []rune, line uint32, column uint32) Location {
	pos := uint32(0)
	lineCount := uint32(0)
	columnCount := uint32(0)
	for i := uint32(0); i < uint32(len(content)); i++ {
		if lineCount == line && columnCount == column {
			pos = i
		}

		if '\n' == content[i] {
			lineCount++
			columnCount = 0
		} else {
			columnCount++
		}
	}
	return NewLocationCursor(filePath, content, pos)
}

func (loc Location) EqualsTo(other Location) bool {
	return loc.filePath == other.filePath && loc.start == other.start && loc.end == other.end
}

func (loc Location) IsEmpty() bool {
	return loc.filePath == ""
}

func (loc Location) CursorString() string {
	if loc.IsEmpty() {
		return ""
	}
	line, col, _, _ := loc.GetLineAndColumn()
	return fmt.Sprintf("%s:%d:%d", loc.filePath, line, col)
}

func (loc Location) GetLineAndColumn() (startLine, startColumn, endLine, endColumn int) {
	line := 1
	column := 1

	for i := uint32(0); i < uint32(len(loc.fileContent)); i++ {
		if i == loc.start {
			startLine = line
			startColumn = column
		}
		if i == loc.end {
			endLine = line
			endColumn = column
		}

		if '\n' == loc.fileContent[i] {
			line++
			column = 1
		} else {
			column++
		}
	}
	return
}

func (loc Location) ToToken(type_ SemanticTokenType, modifiers ...SemanticTokenModifier) SemanticToken {
	l, c, _, _ := loc.GetLineAndColumn()
	mod := SemanticTokenModifier(0)
	for _, m := range modifiers {
		mod |= m
	}
	return SemanticToken{
		Line:      uint32(l - 1),
		Char:      uint32(c - 1),
		Length:    loc.Size(),
		Type:      type_,
		Modifiers: mod,
	}
}

func (loc Location) FilePath() string {
	return loc.filePath
}

func (loc Location) Text() string {
	return string(loc.fileContent[loc.start:loc.end])
}

func (loc Location) Contains(cursor Location) bool {
	return loc.start <= cursor.start && cursor.end <= loc.end
}

func (loc Location) Start() uint32 {
	return loc.start
}

func (loc Location) End() uint32 {
	return loc.end
}

func (loc Location) Size() uint32 {
	return loc.end - loc.start
}

func (loc Location) FileContent() []rune {
	return loc.fileContent
}

func (loc Location) Bytecode() bytecode.Location {
	l, c, _, _ := loc.GetLineAndColumn()
	return bytecode.Location{Line: uint32(l), Column: uint32(c)}
}
