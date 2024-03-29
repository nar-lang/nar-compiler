package compiler

import (
	"fmt"
	"github.com/nar-lang/nar-compiler"
	"github.com/nar-lang/nar-compiler/ast"
	"github.com/nar-lang/nar-compiler/ast/normalized"
	"github.com/nar-lang/nar-compiler/ast/parsed"
	"github.com/nar-lang/nar-compiler/ast/typed"
	"github.com/nar-lang/nar-compiler/bytecode"
	"github.com/nar-lang/nar-compiler/common"
	"github.com/nar-lang/nar-compiler/linker"
	"github.com/nar-lang/nar-compiler/locator"
	"github.com/nar-lang/nar-compiler/logger"
)

const Version uint32 = 100

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
	bin.CompilerVersion = Version
	hash := bytecode.NewBinaryHash()

	packages, err := lc.Packages()
	if err != nil {
		log.Err(err)
		return bin, nil
	}

	for _, pkg := range packages {
		bin.Packages[bytecode.QualifiedIdentifier(pkg.Info().Name)] = uint32(pkg.Info().Version)
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
