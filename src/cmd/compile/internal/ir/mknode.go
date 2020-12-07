// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"bytes"
	"fmt"
	"go/format"
	"go/types"
	"io/ioutil"
	"log"
	"strings"

	"golang.org/x/tools/go/packages"
)

func main() {
	cfg := &packages.Config{
		Mode: packages.NeedSyntax | packages.NeedTypes,
	}
	pkgs, err := packages.Load(cfg, "cmd/compile/internal/ir")
	if err != nil {
		log.Fatal(err)
	}

	pkg := pkgs[0].Types
	scope := pkg.Scope()

	lookup := func(name string) *types.Named {
		return scope.Lookup(name).(*types.TypeName).Type().(*types.Named)
	}

	nodeType := lookup("Node")
	ntypeType := lookup("Ntype")
	nodesType := lookup("Nodes")
	ptrFieldType := types.NewPointer(lookup("Field"))
	slicePtrFieldType := types.NewSlice(ptrFieldType)
	ptrNameType := types.NewPointer(lookup("Name"))

	var buf bytes.Buffer
	fmt.Fprintln(&buf, "// Code generated by mknode.go. DO NOT EDIT.")
	fmt.Fprintln(&buf)
	fmt.Fprintln(&buf, "package ir")
	fmt.Fprintln(&buf)
	fmt.Fprintln(&buf, `import "fmt"`)

	for _, name := range scope.Names() {
		obj, ok := scope.Lookup(name).(*types.TypeName)
		if !ok {
			continue
		}

		typName := obj.Name()
		typ, ok := obj.Type().(*types.Named).Underlying().(*types.Struct)
		if !ok {
			continue
		}

		if strings.HasPrefix(typName, "mini") || !hasMiniNode(typ) {
			continue
		}

		fmt.Fprintf(&buf, "\n")
		fmt.Fprintf(&buf, "func (n *%s) String() string { return fmt.Sprint(n) }\n", name)
		fmt.Fprintf(&buf, "func (n *%s) Format(s fmt.State, verb rune) { FmtNode(n, s, verb) }\n", name)

		fmt.Fprintf(&buf, "func (n *%s) copy() Node { c := *n\n", name)
		forNodeFields(typName, typ, func(name string, is func(types.Type) bool) {
			switch {
			case is(nodesType):
				fmt.Fprintf(&buf, "c.%s = c.%s.Copy()\n", name, name)
			case is(ptrFieldType):
				fmt.Fprintf(&buf, "if c.%s != nil { c.%s = c.%s.copy() }\n", name, name, name)
			case is(slicePtrFieldType):
				fmt.Fprintf(&buf, "c.%s = copyFields(c.%s)\n", name, name)
			}
		})
		fmt.Fprintf(&buf, "return &c }\n")

		fmt.Fprintf(&buf, "func (n *%s) doChildren(do func(Node) error) error { var err error\n", name)
		forNodeFields(typName, typ, func(name string, is func(types.Type) bool) {
			switch {
			case is(ptrNameType):
				fmt.Fprintf(&buf, "if n.%s != nil { err = maybeDo(n.%s, err, do) }\n", name, name)
			case is(nodeType), is(ntypeType):
				fmt.Fprintf(&buf, "err = maybeDo(n.%s, err, do)\n", name)
			case is(nodesType):
				fmt.Fprintf(&buf, "err = maybeDoList(n.%s, err, do)\n", name)
			case is(ptrFieldType):
				fmt.Fprintf(&buf, "err = maybeDoField(n.%s, err, do)\n", name)
			case is(slicePtrFieldType):
				fmt.Fprintf(&buf, "err = maybeDoFields(n.%s, err, do)\n", name)
			}
		})
		fmt.Fprintf(&buf, "return err }\n")

		fmt.Fprintf(&buf, "func (n *%s) editChildren(edit func(Node) Node) {\n", name)
		forNodeFields(typName, typ, func(name string, is func(types.Type) bool) {
			switch {
			case is(ptrNameType):
				fmt.Fprintf(&buf, "if n.%s != nil { n.%s = edit(n.%s).(*Name) }\n", name, name, name)
			case is(nodeType):
				fmt.Fprintf(&buf, "n.%s = maybeEdit(n.%s, edit)\n", name, name)
			case is(ntypeType):
				fmt.Fprintf(&buf, "n.%s = toNtype(maybeEdit(n.%s, edit))\n", name, name)
			case is(nodesType):
				fmt.Fprintf(&buf, "editList(n.%s, edit)\n", name)
			case is(ptrFieldType):
				fmt.Fprintf(&buf, "editField(n.%s, edit)\n", name)
			case is(slicePtrFieldType):
				fmt.Fprintf(&buf, "editFields(n.%s, edit)\n", name)
			}
		})
		fmt.Fprintf(&buf, "}\n")
	}

	out, err := format.Source(buf.Bytes())
	if err != nil {
		// write out mangled source so we can see the bug.
		out = buf.Bytes()
	}

	err = ioutil.WriteFile("node_gen.go", out, 0666)
	if err != nil {
		log.Fatal(err)
	}
}

func forNodeFields(typName string, typ *types.Struct, f func(name string, is func(types.Type) bool)) {
	for i, n := 0, typ.NumFields(); i < n; i++ {
		v := typ.Field(i)
		if v.Embedded() {
			if typ, ok := v.Type().Underlying().(*types.Struct); ok {
				forNodeFields(typName, typ, f)
				continue
			}
		}
		switch typName {
		case "Func":
			if strings.ToLower(strings.TrimSuffix(v.Name(), "_")) != "body" {
				continue
			}
		case "Name":
			continue
		}
		switch v.Name() {
		case "orig":
			continue
		}
		switch typName + "." + v.Name() {
		case "AddStringExpr.Alloc":
			continue
		}
		f(v.Name(), func(t types.Type) bool { return types.Identical(t, v.Type()) })
	}
}

func hasMiniNode(typ *types.Struct) bool {
	for i, n := 0, typ.NumFields(); i < n; i++ {
		v := typ.Field(i)
		if v.Name() == "miniNode" {
			return true
		}
		if v.Embedded() {
			if typ, ok := v.Type().Underlying().(*types.Struct); ok && hasMiniNode(typ) {
				return true
			}
		}
	}
	return false
}
