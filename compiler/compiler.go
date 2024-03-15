package compiler

import (
	"fmt"
	"github.com/nar-lang/nar-common/ast"
	"github.com/nar-lang/nar-common/ast/normalized"
	"github.com/nar-lang/nar-common/ast/parsed"
	"github.com/nar-lang/nar-common/ast/typed"
	"github.com/nar-lang/nar-common/bytecode"
	"github.com/nar-lang/nar-common/common"
	"github.com/nar-lang/nar-common/logger"
	"github.com/nar-lang/nar-compiler"
	"github.com/nar-lang/nar-compiler/linker"
	"github.com/nar-lang/nar-compiler/locator"
)

func Compile(log *logger.LogWriter, lc locator.Locator, link linker.Linker, debug bool) *bytecode.Binary {
	parsedModules := map[ast.QualifiedIdentifier]*parsed.Module{}
	normalizedModules := map[ast.QualifiedIdentifier]*normalized.Module{}
	typedModules := map[ast.QualifiedIdentifier]*typed.Module{}
	bin, _ := CompileEx(log, lc, link, debug, parsedModules, normalizedModules, typedModules)
	return bin
}

func CompileEx(
	log *logger.LogWriter, lc locator.Locator, link linker.Linker, debug bool,
	parsedModules map[ast.QualifiedIdentifier]*parsed.Module,
	normalizedModules map[ast.QualifiedIdentifier]*normalized.Module,
	typedModules map[ast.QualifiedIdentifier]*typed.Module,
) (*bytecode.Binary, []ast.QualifiedIdentifier) {

	bin := bytecode.NewBinary()
	hash := bytecode.NewBinaryHash()

	packages, err := lc.Packages()
	if err != nil {
		log.Err(err)
		return bin, nil
	}

	for _, pkg := range packages {
		bin.Packages[bytecode.QualifiedIdentifier(pkg.Info().Name)] = int32(pkg.Info().Version)
	}

	affectedModuleNames := nar_compiler.Compile(
		log,
		packages,
		parsedModules,
		normalizedModules,
		typedModules)

	if len(log.Errors()) == 0 {
		for _, name := range affectedModuleNames {
			m, ok := typedModules[name]
			if !ok {
				log.Err(common.NewSystemError(fmt.Errorf("module '%s' not found", name)))
				continue
			}
			if err := m.Compose(typedModules, debug, bin, hash); err != nil {
				log.Err(err)
			}
		}
	}

	if !log.Err() {
		if link != nil {
			err := link.Link(log, bin, lc, debug)
			if err != nil {
				log.Err(err)
			}
		}
	}
	return bin, affectedModuleNames
}

func Version() int {
	return common.CompilerVersion
}
