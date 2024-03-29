package linker

import (
	"bufio"
	"bytes"
	"github.com/nar-lang/nar-compiler/bytecode"
	"github.com/nar-lang/nar-compiler/locator"
	"github.com/nar-lang/nar-compiler/logger"
	"os"
	"path/filepath"
)

func NewDllLinker(outFilePath string) Linker {
	return &dllLinker{outFilePath: outFilePath}
}

type dllLinker struct {
	outFilePath string
}

func (d dllLinker) Link(log *logger.LogWriter, binary *bytecode.Binary, lc locator.Locator, debug bool) error {
	var err error
	binary.Entry, err = lc.EntryPoint()
	if err != nil {
		return err
	}
	outDir := filepath.Dir(d.outFilePath)
	_ = os.Mkdir(outDir, 0755)

	buf := bytes.NewBuffer(nil)
	w := bufio.NewWriter(buf)
	err = binary.Write(w, debug)
	if err != nil {
		return err
	}
	err = w.Flush()
	if err != nil {
		return err
	}
	err = os.WriteFile(filepath.Join(d.outFilePath), buf.Bytes(), 0755)
	if err != nil {
		return err
	}

	for pkgName := range binary.Packages {
		pkg, ok, err := lc.FindPackage(string(pkgName))
		if err != nil {
			return err
		}
		if ok {
			paths, err := pkg.NativeFilePaths("dll")
			if err != nil {
				return err
			}
			for _, path := range paths {
				f, err := os.ReadFile(path)
				if err != nil {
					return err
				}
				targetPath := filepath.Join(outDir, filepath.Base(path))
				err = os.WriteFile(targetPath, f, 0755)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
