// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore

// This file implements a custom generator to create various go/types
// source files from the corresponding types2 files.

package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	srcDir = "cmd/compile/internal/types2"
	dstDir = "go/types"
)

var fset = token.NewFileSet()

func main() {
	flag.Parse()

	// process provided filenames, if any
	if flag.NArg() > 0 {
		for _, filename := range flag.Args() {
			fmt.Println("generating", filename)
			generate(filename, filemap[filename])
		}
		return
	}

	// otherwise process per filemap below
	for filename, action := range filemap {
		generate(filename, action)
	}
}

func generate(filename string, action action) {
	// parse src
	srcFilename := filepath.FromSlash(runtime.GOROOT() + "/src/" + srcDir + "/" + filename)
	file, err := parser.ParseFile(fset, srcFilename, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	// fix package name
	file.Name.Name = strings.ReplaceAll(file.Name.Name, "types2", "types")

	// rewrite AST as needed
	if action != nil {
		action(file)
	}

	// format AST
	var buf bytes.Buffer
	buf.WriteString("// Code generated by \"go run generator.go\"; DO NOT EDIT.\n\n")
	if err := format.Node(&buf, fset, file); err != nil {
		log.Fatal(err)
	}

	// write dst
	dstFilename := filepath.FromSlash(runtime.GOROOT() + "/src/" + dstDir + "/" + filename)
	if err := os.WriteFile(dstFilename, buf.Bytes(), 0o644); err != nil {
		log.Fatal(err)
	}
}

type action func(in *ast.File)

var filemap = map[string]action{
	"array.go":            nil,
	"basic.go":            nil,
	"chan.go":             nil,
	"context.go":          nil,
	"context_test.go":     nil,
	"gccgosizes.go":       nil,
	"instantiate_test.go": func(f *ast.File) { renameImportPath(f, `"cmd/compile/internal/types2"`, `"go/types"`) },
	"lookup.go":           nil,
	"main_test.go":        nil,
	"map.go":              nil,
	"named.go":            func(f *ast.File) { fixTokenPos(f); fixTraceSel(f) },
	"object.go":           func(f *ast.File) { fixTokenPos(f); renameIdent(f, "NewTypeNameLazy", "_NewTypeNameLazy") },
	"objset.go":           nil,
	"package.go":          nil,
	"pointer.go":          nil,
	"predicates.go":       nil,
	"scope.go": func(f *ast.File) {
		fixTokenPos(f)
		renameIdent(f, "Squash", "squash")
		renameIdent(f, "InsertLazy", "_InsertLazy")
	},
	"selection.go":     nil,
	"sizes.go":         func(f *ast.File) { renameIdent(f, "IsSyncAtomicAlign64", "isSyncAtomicAlign64") },
	"slice.go":         nil,
	"subst.go":         func(f *ast.File) { fixTokenPos(f); fixTraceSel(f) },
	"termlist.go":      nil,
	"termlist_test.go": nil,
	"tuple.go":         nil,
	"typelists.go":     nil,
	"typeparam.go":     nil,
	"typeterm_test.go": nil,
	"typeterm.go":      nil,
	"under.go":         nil,
	"unify.go":         fixSprintf,
	"universe.go":      fixGlobalTypVarDecl,
	"validtype.go":     nil,
}

// TODO(gri) We should be able to make these rewriters more configurable/composable.
//           For now this is a good starting point.

// renameIdent renames an identifier.
// Note: This doesn't change the use of the identifier in comments.
func renameIdent(f *ast.File, from, to string) {
	ast.Inspect(f, func(n ast.Node) bool {
		switch n := n.(type) {
		case *ast.Ident:
			if n.Name == from {
				n.Name = to
			}
			return false
		}
		return true
	})
}

// renameImportPath renames an import path.
func renameImportPath(f *ast.File, from, to string) {
	ast.Inspect(f, func(n ast.Node) bool {
		switch n := n.(type) {
		case *ast.ImportSpec:
			if n.Path.Kind == token.STRING && n.Path.Value == from {
				n.Path.Value = to
				return false
			}
		}
		return true
	})
}

// fixTokenPos changes imports of "cmd/compile/internal/syntax" to "go/token"
// and uses of syntax.Pos to token.Pos.
func fixTokenPos(f *ast.File) {
	ast.Inspect(f, func(n ast.Node) bool {
		switch n := n.(type) {
		case *ast.ImportSpec:
			// rewrite import path "cmd/compile/internal/syntax" to "go/token"
			if n.Path.Kind == token.STRING && n.Path.Value == `"cmd/compile/internal/syntax"` {
				n.Path.Value = `"go/token"`
				return false
			}
		case *ast.SelectorExpr:
			// rewrite syntax.Pos to token.Pos
			if x, _ := n.X.(*ast.Ident); x != nil && x.Name == "syntax" && n.Sel.Name == "Pos" {
				x.Name = "token"
				return false
			}
		case *ast.CallExpr:
			// rewrite x.IsKnown() to x.IsValid()
			if fun, _ := n.Fun.(*ast.SelectorExpr); fun != nil && fun.Sel.Name == "IsKnown" && len(n.Args) == 0 {
				fun.Sel.Name = "IsValid"
				return false
			}
		}
		return true
	})
}

// fixTraceSel renames uses of x.Trace to x.trace, where x for any x with a Trace field.
func fixTraceSel(f *ast.File) {
	ast.Inspect(f, func(n ast.Node) bool {
		switch n := n.(type) {
		case *ast.SelectorExpr:
			// rewrite x.Trace to x.trace (for Config.Trace)
			if n.Sel.Name == "Trace" {
				n.Sel.Name = "trace"
				return false
			}
		}
		return true
	})
}

// fixGlobalTypVarDecl changes the global Typ variable from an array to a slice
// (in types2 we use an array for efficiency, in go/types it's a slice and we
// cannot change that).
func fixGlobalTypVarDecl(f *ast.File) {
	ast.Inspect(f, func(n ast.Node) bool {
		switch n := n.(type) {
		case *ast.ValueSpec:
			// rewrite type Typ = [...]Type{...} to type Typ = []Type{...}
			if len(n.Names) == 1 && n.Names[0].Name == "Typ" && len(n.Values) == 1 {
				n.Values[0].(*ast.CompositeLit).Type.(*ast.ArrayType).Len = nil
				return false
			}
		}
		return true
	})
}

// fixSprintf adds an extra nil argument for the *token.FileSet parameter in sprintf calls.
func fixSprintf(f *ast.File) {
	ast.Inspect(f, func(n ast.Node) bool {
		switch n := n.(type) {
		case *ast.CallExpr:
			if fun, _ := n.Fun.(*ast.Ident); fun != nil && fun.Name == "sprintf" && len(n.Args) >= 4 /* ... args */ {
				n.Args = insert(n.Args, 1, newIdent(n.Args[1].Pos(), "nil"))
				return false
			}
		}
		return true
	})
}

// newIdent returns a new identifier with the given position and name.
func newIdent(pos token.Pos, name string) *ast.Ident {
	id := ast.NewIdent(name)
	id.NamePos = pos
	return id
}

// insert inserts x at list[at] and moves the remaining elements up.
func insert(list []ast.Expr, at int, x ast.Expr) []ast.Expr {
	list = append(list, nil)
	copy(list[at+1:], list[at:])
	list[at] = x
	return list
}
