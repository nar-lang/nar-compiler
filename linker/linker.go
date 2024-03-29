package linker

import (
	"github.com/nar-lang/nar-compiler/bytecode"
	"github.com/nar-lang/nar-compiler/locator"
	"github.com/nar-lang/nar-compiler/logger"
)

type Linker interface {
	Link(log *logger.LogWriter, binary *bytecode.Binary, lc locator.Locator, debug bool) error
}
