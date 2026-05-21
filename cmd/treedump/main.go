// Command treedump walks one or more directories looking for `.nar` source
// files. For every file it parses the source using the package-level Parse
// entrypoint and writes the parsed AST tree, rendered via
// (*parsed.Module).StringTree(0), to `<file>.tree.go.txt`.
//
// With the `--normalized` flag, all `.nar` files under the given roots are
// parsed first, then normalized (with module dependencies fully resolved),
// and the normalized AST tree is written to `<file>.tree.normalized.go.txt`.
//
// With the `--typed` flag, normalized modules are additionally lowered into
// the typed AST via `(*normalized.Module).Annotate(...)` and the resulting
// typed tree is written to `<file>.tree.typed.go.txt`.
//
// With the `--checked` flag, typed modules additionally run CheckTypes()
// (Hindley-Milner unification) and CheckPatterns() (exhaustiveness +
// redundancy). The fully-typed tree is written to
// `<file>.tree.checked.go.txt`.
//
// Usage:
//
//	go run ./cmd/treedump [--normalized | --typed | --checked] <root> [<root> ...]
package main

import (
	"fmt"
	nar_compiler "github.com/nar-lang/nar-compiler"
	"github.com/nar-lang/nar-compiler/ast"
	"github.com/nar-lang/nar-compiler/ast/normalized"
	"github.com/nar-lang/nar-compiler/ast/parsed"
	"github.com/nar-lang/nar-compiler/ast/typed"
	"github.com/nar-lang/nar-compiler/bytecode"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func main() {
	args := os.Args[1:]
	normalizedMode := false
	typedMode := false
	checkedMode := false
	bytecodeMode := false
	roots := []string{}
	for _, a := range args {
		if a == "--normalized" {
			normalizedMode = true
		} else if a == "--typed" {
			typedMode = true
			normalizedMode = true
		} else if a == "--checked" {
			checkedMode = true
			typedMode = true
			normalizedMode = true
		} else if a == "--bytecode" {
			bytecodeMode = true
			checkedMode = true
			typedMode = true
			normalizedMode = true
		} else {
			roots = append(roots, a)
		}
	}
	if len(roots) == 0 {
		fmt.Fprintln(os.Stderr, "usage: treedump [--normalized | --typed | --checked | --bytecode] <root> [<root> ...]")
		os.Exit(2)
	}

	if normalizedMode {
		runNormalized(roots, typedMode, checkedMode, bytecodeMode)
		return
	}
	runParsed(roots)
}

func runParsed(roots []string) {
	totalOk, totalErr := 0, 0
	for _, root := range roots {
		err := filepath.WalkDir(root, func(path string, d os.DirEntry, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}
			if d.IsDir() || !strings.HasSuffix(path, ".nar") {
				return nil
			}

			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			module, errs := nar_compiler.Parse(path, []rune(string(data)))
			if len(errs) > 0 {
				totalErr++
				fmt.Fprintf(os.Stderr, "parse %s: %v\n", path, errs[0])
				return nil
			}
			out := path + ".tree.go.txt"
			if err := os.WriteFile(out, []byte(module.StringTree(0)), 0o644); err != nil {
				return err
			}
			totalOk++
			return nil
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "walk %s: %v\n", root, err)
			os.Exit(1)
		}
	}
	fmt.Printf("treedump: %d ok, %d failed\n", totalOk, totalErr)
}

func runNormalized(roots []string, typedMode bool, checkedMode bool, bytecodeMode bool) {
	// 1) Find all .nar files
	var files []string
	for _, root := range roots {
		err := filepath.WalkDir(root, func(path string, d os.DirEntry, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}
			if d.IsDir() || !strings.HasSuffix(path, ".nar") {
				return nil
			}
			files = append(files, path)
			return nil
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "walk %s: %v\n", root, err)
			os.Exit(1)
		}
	}
	sort.Strings(files)

	// 2) Parse all
	parsedModules := map[ast.QualifiedIdentifier]*parsed.Module{}
	moduleFile := map[ast.QualifiedIdentifier]string{}
	totalErr := 0
	for _, path := range files {
		data, err := os.ReadFile(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "read %s: %v\n", path, err)
			totalErr++
			continue
		}
		m, errs := nar_compiler.Parse(path, []rune(string(data)))
		if len(errs) > 0 {
			fmt.Fprintf(os.Stderr, "parse %s: %v\n", path, errs[0])
			totalErr++
			continue
		}
		if m == nil {
			continue
		}
		// All modules are treated as belonging to a single virtual package
		// so that isReferenced returns true between any pair.
		m.SetPackageName(ast.PackageIdentifier("treedump"))
		m.SetReferencedPackages(map[ast.PackageIdentifier]struct{}{
			"treedump": {},
		})
		parsedModules[m.Name()] = m
		moduleFile[m.Name()] = path
	}

	// 3) Generate (lower data types, expand imports)
	names := make([]ast.QualifiedIdentifier, 0, len(parsedModules))
	for n := range parsedModules {
		names = append(names, n)
	}
	sort.Slice(names, func(i, j int) bool { return string(names[i]) < string(names[j]) })

	for _, n := range names {
		errs := parsedModules[n].Generate(parsedModules)
		for _, e := range errs {
			fmt.Fprintf(os.Stderr, "generate %s: %v\n", n, e)
			totalErr++
		}
	}

	// 4) Normalize
	normalizedModules := map[ast.QualifiedIdentifier]*normalized.Module{}
	for _, n := range names {
		errs := parsedModules[n].Normalize(parsedModules, normalizedModules)
		for _, e := range errs {
			fmt.Fprintf(os.Stderr, "normalize %s: %v\n", n, e)
			totalErr++
		}
	}

	// 5) Write outputs
	totalOk := 0
	for _, n := range names {
		nm, ok := normalizedModules[n]
		if !ok || nm == nil {
			continue
		}
		path, ok := moduleFile[n]
		if !ok {
			continue
		}
		out := path + ".tree.normalized.go.txt"
		if err := os.WriteFile(out, []byte(nm.StringTree(0)), 0o644); err != nil {
			fmt.Fprintf(os.Stderr, "write %s: %v\n", out, err)
			totalErr++
			continue
		}
		totalOk++
	}
	fmt.Printf("treedump (normalized): %d ok, %d failed\n", totalOk, totalErr)

	if !typedMode {
		return
	}

	// 6) Annotate (normalized -> typed)
	typedModules := map[ast.QualifiedIdentifier]*typed.Module{}
	for _, n := range names {
		nm, ok := normalizedModules[n]
		if !ok || nm == nil {
			continue
		}
		errs := nm.Annotate(normalizedModules, typedModules)
		for _, e := range errs {
			fmt.Fprintf(os.Stderr, "annotate %s: %v\n", n, e)
			totalErr++
		}
	}

	// 7) Write typed outputs
	typedOk := 0
	for _, n := range names {
		tm, ok := typedModules[n]
		if !ok || tm == nil {
			continue
		}
		path, ok := moduleFile[n]
		if !ok {
			continue
		}
		out := path + ".tree.typed.go.txt"
		if err := os.WriteFile(out, []byte(tm.StringTree(0)), 0o644); err != nil {
			fmt.Fprintf(os.Stderr, "write %s: %v\n", out, err)
			totalErr++
			continue
		}
		typedOk++
	}
	fmt.Printf("treedump (typed): %d ok, %d failed\n", typedOk, totalErr)

	if !checkedMode {
		return
	}

	// 8) CheckTypes on every typed module.
	for _, n := range names {
		tm, ok := typedModules[n]
		if !ok || tm == nil {
			continue
		}
		errs := tm.CheckTypes()
		for _, e := range errs {
			fmt.Fprintf(os.Stderr, "checkTypes %s: %v\n", n, e)
			totalErr++
		}
	}

	// 9) CheckPatterns on every typed module.
	for _, n := range names {
		tm, ok := typedModules[n]
		if !ok || tm == nil {
			continue
		}
		errs := tm.CheckPatterns()
		for _, e := range errs {
			fmt.Fprintf(os.Stderr, "checkPatterns %s: %v\n", n, e)
			totalErr++
		}
	}

	// 10) Write checked outputs.
	checkedOk := 0
	for _, n := range names {
		tm, ok := typedModules[n]
		if !ok || tm == nil {
			continue
		}
		path, ok := moduleFile[n]
		if !ok {
			continue
		}
		out := path + ".tree.checked.go.txt"
		if err := os.WriteFile(out, []byte(tm.StringTree(0)), 0o644); err != nil {
			fmt.Fprintf(os.Stderr, "write %s: %v\n", out, err)
			totalErr++
			continue
		}
		checkedOk++
	}
	fmt.Printf("treedump (checked): %d ok, %d failed\n", checkedOk, totalErr)

	if !bytecodeMode {
		return
	}

	// 11) For each root, build a single Binary blob from all typed modules
	// under that root.
	filesByRoot := map[string]map[string]struct{}{}
	for _, root := range roots {
		filesByRoot[root] = map[string]struct{}{}
		_ = filepath.WalkDir(root, func(path string, d os.DirEntry, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}
			if d.IsDir() || !strings.HasSuffix(path, ".nar") {
				return nil
			}
			filesByRoot[root][path] = struct{}{}
			return nil
		})
	}

	bytecodeOk := 0
	for _, root := range roots {
		var rootNames []ast.QualifiedIdentifier
		for _, n := range names {
			path, ok := moduleFile[n]
			if !ok {
				continue
			}
			if _, inRoot := filesByRoot[root][path]; inRoot {
				rootNames = append(rootNames, n)
			}
		}
		if len(rootNames) == 0 {
			continue
		}
		sort.Slice(rootNames, func(i, j int) bool { return string(rootNames[i]) < string(rootNames[j]) })

		bin := bytecode.NewBinary()
		hash := bytecode.NewBinaryHash()
		composeOk := true
		for _, n := range rootNames {
			tm := typedModules[n]
			if tm == nil {
				fmt.Fprintf(os.Stderr, "compose %s: typed module missing\n", n)
				totalErr++
				composeOk = false
				break
			}
			if err := tm.Compose(typedModules, true, bin, hash); err != nil {
				fmt.Fprintf(os.Stderr, "compose %s: %v\n", n, err)
				totalErr++
				composeOk = false
				break
			}
		}
		if composeOk {
			outPath := filepath.Join(root, "binary.go.bin")
			f, err := os.Create(outPath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "write %s: %v\n", outPath, err)
				totalErr++
				continue
			}
			if err := bin.Write(f, true); err != nil {
				fmt.Fprintf(os.Stderr, "write %s: %v\n", outPath, err)
				totalErr++
			}
			_ = f.Close()
			bytecodeOk++
		}
	}
	fmt.Printf("treedump (bytecode): %d ok, %d failed\n", bytecodeOk, totalErr)
}
