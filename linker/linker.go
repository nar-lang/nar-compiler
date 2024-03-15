package linker

import (
	"github.com/nar-lang/nar-common/bytecode"
	"github.com/nar-lang/nar-common/logger"
	"github.com/nar-lang/nar-compiler/locator"
)

type Linker interface {
	Link(log *logger.LogWriter, binary *bytecode.Binary, lc locator.Locator, debug bool) error
}
