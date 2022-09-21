// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file implements tests for various issues.

package types2_test

import (
	"cmd/compile/internal/syntax"
	"fmt"
	"internal/testenv"
	"sort"
	"strings"
	"testing"

	. "cmd/compile/internal/types2"
)

func mustParse(t *testing.T, src string) *syntax.File {
	f, err := parseSrc("", src)
	if err != nil {
		t.Fatal(err)
	}
	return f
}
func TestIssue5770(t *testing.T) {
	f := mustParse(t, `package p; type S struct{T}`)
	var conf Config
	_, err := conf.Check(f.PkgName.Value, []*syntax.File{f}, nil) // do not crash
	const want = "undefined: T"
	if err == nil || !strings.Contains(err.Error(), want) {
		t.Errorf("got: %v; want: %s", err, want)
	}
}

func TestIssue5849(t *testing.T) {
	src := `
package p
var (
	s uint
	_ = uint8(8)
	_ = uint16(16) << s
	_ = uint32(32 << s)
	_ = uint64(64 << s + s)
	_ = (interface{})("foo")
	_ = (interface{})(nil)
)`
	f := mustParse(t, src)

	var conf Config
	types := make(map[syntax.Expr]TypeAndValue)
	_, err := conf.Check(f.PkgName.Value, []*syntax.File{f}, &Info{Types: types})
	if err != nil {
		t.Fatal(err)
	}

	for x, tv := range types {
		var want Type
		switch x := x.(type) {
		case *syntax.BasicLit:
			switch x.Value {
			case `8`:
				want = Typ[Uint8]
			case `16`:
				want = Typ[Uint16]
			case `32`:
				want = Typ[Uint32]
			case `64`:
				want = Typ[Uint] // because of "+ s", s is of type uint
			case `"foo"`:
				want = Typ[String]
			}
		case *syntax.Name:
			if x.Value == "nil" {
				want = NewInterfaceType(nil, nil) // interface{} (for now, go/types types this as "untyped nil")
			}
		}
		if want != nil && !Identical(tv.Type, want) {
			t.Errorf("got %s; want %s", tv.Type, want)
		}
	}
}

func TestIssue6413(t *testing.T) {
	src := `
package p
func f() int {
	defer f()
	go f()
	return 0
}
`
	f := mustParse(t, src)

	var conf Config
	types := make(map[syntax.Expr]TypeAndValue)
	_, err := conf.Check(f.PkgName.Value, []*syntax.File{f}, &Info{Types: types})
	if err != nil {
		t.Fatal(err)
	}

	want := Typ[Int]
	n := 0
	for x, tv := range types {
		if _, ok := x.(*syntax.CallExpr); ok {
			if tv.Type != want {
				t.Errorf("%s: got %s; want %s", x.Pos(), tv.Type, want)
			}
			n++
		}
	}

	if n != 2 {
		t.Errorf("got %d CallExprs; want 2", n)
	}
}

func TestIssue7245(t *testing.T) {
	src := `
package p
func (T) m() (res bool) { return }
type T struct{} // receiver type after method declaration
`
	f := mustParse(t, src)

	var conf Config
	defs := make(map[*syntax.Name]Object)
	_, err := conf.Check(f.PkgName.Value, []*syntax.File{f}, &Info{Defs: defs})
	if err != nil {
		t.Fatal(err)
	}

	m := f.DeclList[0].(*syntax.FuncDecl)
	res1 := defs[m.Name].(*Func).Type().(*Signature).Results().At(0)
	res2 := defs[m.Type.ResultList[0].Name].(*Var)

	if res1 != res2 {
		t.Errorf("got %s (%p) != %s (%p)", res1, res2, res1, res2)
	}
}

// This tests that uses of existing vars on the LHS of an assignment
// are Uses, not Defs; and also that the (illegal) use of a non-var on
// the LHS of an assignment is a Use nonetheless.
func TestIssue7827(t *testing.T) {
	const src = `
package p
func _() {
	const w = 1        // defs w
        x, y := 2, 3       // defs x, y
        w, x, z := 4, 5, 6 // uses w, x, defs z; error: cannot assign to w
        _, _, _ = x, y, z  // uses x, y, z
}
`
	f := mustParse(t, src)

	const want = `L3 defs func p._()
L4 defs const w untyped int
L5 defs var x int
L5 defs var y int
L6 defs var z int
L6 uses const w untyped int
L6 uses var x int
L7 uses var x int
L7 uses var y int
L7 uses var z int`

	// don't abort at the first error
	conf := Config{Error: func(err error) { t.Log(err) }}
	defs := make(map[*syntax.Name]Object)
	uses := make(map[*syntax.Name]Object)
	_, err := conf.Check(f.PkgName.Value, []*syntax.File{f}, &Info{Defs: defs, Uses: uses})
	if s := fmt.Sprint(err); !strings.HasSuffix(s, "cannot assign to w") {
		t.Errorf("Check: unexpected error: %s", s)
	}

	var facts []string
	for id, obj := range defs {
		if obj != nil {
			fact := fmt.Sprintf("L%d defs %s", id.Pos().Line(), obj)
			facts = append(facts, fact)
		}
	}
	for id, obj := range uses {
		fact := fmt.Sprintf("L%d uses %s", id.Pos().Line(), obj)
		facts = append(facts, fact)
	}
	sort.Strings(facts)

	got := strings.Join(facts, "\n")
	if got != want {
		t.Errorf("Unexpected defs/uses\ngot:\n%s\nwant:\n%s", got, want)
	}
}

// This tests that the package associated with the types2.Object.Pkg method
// is the type's package independent of the order in which the imports are
// listed in the sources src1, src2 below.
// The actual issue is in go/internal/gcimporter which has a corresponding
// test; we leave this test here to verify correct behavior at the go/types
// level.
func TestIssue13898(t *testing.T) {
	testenv.MustHaveGoBuild(t)

	const src0 = `
package main

import "go/types"

func main() {
	var info types.Info
	for _, obj := range info.Uses {
		_ = obj.Pkg()
	}
}
`
	// like src0, but also imports go/importer
	const src1 = `
package main

import (
	"go/types"
	_ "go/importer"
)

func main() {
	var info types.Info
	for _, obj := range info.Uses {
		_ = obj.Pkg()
	}
}
`
	// like src1 but with different import order
	// (used to fail with this issue)
	const src2 = `
package main

import (
	_ "go/importer"
	"go/types"
)

func main() {
	var info types.Info
	for _, obj := range info.Uses {
		_ = obj.Pkg()
	}
}
`
	f := func(test, src string) {
		f := mustParse(t, src)
		conf := Config{Importer: defaultImporter()}
		info := Info{Uses: make(map[*syntax.Name]Object)}
		_, err := conf.Check("main", []*syntax.File{f}, &info)
		if err != nil {
			t.Fatal(err)
		}

		var pkg *Package
		count := 0
		for id, obj := range info.Uses {
			if id.Value == "Pkg" {
				pkg = obj.Pkg()
				count++
			}
		}
		if count != 1 {
			t.Fatalf("%s: got %d entries named Pkg; want 1", test, count)
		}
		if pkg.Name() != "types" {
			t.Fatalf("%s: got %v; want package types2", test, pkg)
		}
	}

	f("src0", src0)
	f("src1", src1)
	f("src2", src2)
}

func TestIssue22525(t *testing.T) {
	f := mustParse(t, `package p; func f() { var a, b, c, d, e int }`)

	got := "\n"
	conf := Config{Error: func(err error) { got += err.Error() + "\n" }}
	conf.Check(f.PkgName.Value, []*syntax.File{f}, nil) // do not crash
	want := `
:1:27: a declared but not used
:1:30: b declared but not used
:1:33: c declared but not used
:1:36: d declared but not used
:1:39: e declared but not used
`
	if got != want {
		t.Errorf("got: %swant: %s", got, want)
	}
}

func TestIssue25627(t *testing.T) {
	const prefix = `package p; import "unsafe"; type P *struct{}; type I interface{}; type T `
	// The src strings (without prefix) are constructed such that the number of semicolons
	// plus one corresponds to the number of fields expected in the respective struct.
	for _, src := range []string{
		`struct { x Missing }`,
		`struct { Missing }`,
		`struct { *Missing }`,
		`struct { unsafe.Pointer }`,
		`struct { P }`,
		`struct { *I }`,
		`struct { a int; b Missing; *Missing }`,
	} {
		f := mustParse(t, prefix+src)

		conf := Config{Importer: defaultImporter(), Error: func(err error) {}}
		info := &Info{Types: make(map[syntax.Expr]TypeAndValue)}
		_, err := conf.Check(f.PkgName.Value, []*syntax.File{f}, info)
		if err != nil {
			if _, ok := err.(Error); !ok {
				t.Fatal(err)
			}
		}

		syntax.Crawl(f, func(n syntax.Node) bool {
			if decl, _ := n.(*syntax.TypeDecl); decl != nil {
				if tv, ok := info.Types[decl.Type]; ok && decl.Name.Value == "T" {
					want := strings.Count(src, ";") + 1
					if got := tv.Type.(*Struct).NumFields(); got != want {
						t.Errorf("%s: got %d fields; want %d", src, got, want)
					}
				}
			}
			return false
		})
	}
}

func TestIssue28005(t *testing.T) {
	// method names must match defining interface name for this test
	// (see last comment in this function)
	sources := [...]string{
		"package p; type A interface{ A() }",
		"package p; type B interface{ B() }",
		"package p; type X interface{ A; B }",
	}

	// compute original file ASTs
	var orig [len(sources)]*syntax.File
	for i, src := range sources {
		orig[i] = mustParse(t, src)
	}

	// run the test for all order permutations of the incoming files
	for _, perm := range [][len(sources)]int{
		{0, 1, 2},
		{0, 2, 1},
		{1, 0, 2},
		{1, 2, 0},
		{2, 0, 1},
		{2, 1, 0},
	} {
		// create file order permutation
		files := make([]*syntax.File, len(sources))
		for i := range perm {
			files[i] = orig[perm[i]]
		}

		// type-check package with given file order permutation
		var conf Config
		info := &Info{Defs: make(map[*syntax.Name]Object)}
		_, err := conf.Check("", files, info)
		if err != nil {
			t.Fatal(err)
		}

		// look for interface object X
		var obj Object
		for name, def := range info.Defs {
			if name.Value == "X" {
				obj = def
				break
			}
		}
		if obj == nil {
			t.Fatal("object X not found")
		}
		iface := obj.Type().Underlying().(*Interface) // object X must be an interface

		// Each iface method m is embedded; and m's receiver base type name
		// must match the method's name per the choice in the source file.
		for i := 0; i < iface.NumMethods(); i++ {
			m := iface.Method(i)
			recvName := m.Type().(*Signature).Recv().Type().(*Named).Obj().Name()
			if recvName != m.Name() {
				t.Errorf("perm %v: got recv %s; want %s", perm, recvName, m.Name())
			}
		}
	}
}

func TestIssue28282(t *testing.T) {
	// create type interface { error }
	et := Universe.Lookup("error").Type()
	it := NewInterfaceType(nil, []Type{et})
	// verify that after completing the interface, the embedded method remains unchanged
	// (interfaces are "completed" lazily now, so the completion happens implicitly when
	// accessing Method(0))
	want := et.Underlying().(*Interface).Method(0)
	got := it.Method(0)
	if got != want {
		t.Fatalf("%s.Method(0): got %q (%p); want %q (%p)", it, got, got, want, want)
	}
	// verify that lookup finds the same method in both interfaces (redundant check)
	obj, _, _ := LookupFieldOrMethod(et, false, nil, "Error")
	if obj != want {
		t.Fatalf("%s.Lookup: got %q (%p); want %q (%p)", et, obj, obj, want, want)
	}
	obj, _, _ = LookupFieldOrMethod(it, false, nil, "Error")
	if obj != want {
		t.Fatalf("%s.Lookup: got %q (%p); want %q (%p)", it, obj, obj, want, want)
	}
}

func TestIssue29029(t *testing.T) {
	f1 := mustParse(t, `package p; type A interface { M() }`)
	f2 := mustParse(t, `package p; var B interface { A }`)

	// printInfo prints the *Func definitions recorded in info, one *Func per line.
	printInfo := func(info *Info) string {
		var buf strings.Builder
		for _, obj := range info.Defs {
			if fn, ok := obj.(*Func); ok {
				fmt.Fprintln(&buf, fn)
			}
		}
		return buf.String()
	}

	// The *Func (method) definitions for package p must be the same
	// independent on whether f1 and f2 are type-checked together, or
	// incrementally.

	// type-check together
	var conf Config
	info := &Info{Defs: make(map[*syntax.Name]Object)}
	check := NewChecker(&conf, NewPackage("", "p"), info)
	if err := check.Files([]*syntax.File{f1, f2}); err != nil {
		t.Fatal(err)
	}
	want := printInfo(info)

	// type-check incrementally
	info = &Info{Defs: make(map[*syntax.Name]Object)}
	check = NewChecker(&conf, NewPackage("", "p"), info)
	if err := check.Files([]*syntax.File{f1}); err != nil {
		t.Fatal(err)
	}
	if err := check.Files([]*syntax.File{f2}); err != nil {
		t.Fatal(err)
	}
	got := printInfo(info)

	if got != want {
		t.Errorf("\ngot : %swant: %s", got, want)
	}
}

func TestIssue34151(t *testing.T) {
	const asrc = `package a; type I interface{ M() }; type T struct { F interface { I } }`
	const bsrc = `package b; import "a"; type T struct { F interface { a.I } }; var _ = a.T(T{})`

	a, err := pkgFor("a", asrc, nil)
	if err != nil {
		t.Fatalf("package %s failed to typecheck: %v", a.Name(), err)
	}

	bast := mustParse(t, bsrc)
	conf := Config{Importer: importHelper{pkg: a}}
	b, err := conf.Check(bast.PkgName.Value, []*syntax.File{bast}, nil)
	if err != nil {
		t.Errorf("package %s failed to typecheck: %v", b.Name(), err)
	}
}

type importHelper struct {
	pkg      *Package
	fallback Importer
}

func (h importHelper) Import(path string) (*Package, error) {
	if path == h.pkg.Path() {
		return h.pkg, nil
	}
	if h.fallback == nil {
		return nil, fmt.Errorf("got package path %q; want %q", path, h.pkg.Path())
	}
	return h.fallback.Import(path)
}

// TestIssue34921 verifies that we don't update an imported type's underlying
// type when resolving an underlying type. Specifically, when determining the
// underlying type of b.T (which is the underlying type of a.T, which is int)
// we must not set the underlying type of a.T again since that would lead to
// a race condition if package b is imported elsewhere, in a package that is
// concurrently type-checked.
func TestIssue34921(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Error(r)
		}
	}()

	var sources = []string{
		`package a; type T int`,
		`package b; import "a"; type T a.T`,
	}

	var pkg *Package
	for _, src := range sources {
		f := mustParse(t, src)
		conf := Config{Importer: importHelper{pkg: pkg}}
		res, err := conf.Check(f.PkgName.Value, []*syntax.File{f}, nil)
		if err != nil {
			t.Errorf("%q failed to typecheck: %v", src, err)
		}
		pkg = res // res is imported by the next package in this test
	}
}

func TestIssue43088(t *testing.T) {
	// type T1 struct {
	//         _ T2
	// }
	//
	// type T2 struct {
	//         _ struct {
	//                 _ T2
	//         }
	// }
	n1 := NewTypeName(syntax.Pos{}, nil, "T1", nil)
	T1 := NewNamed(n1, nil, nil)
	n2 := NewTypeName(syntax.Pos{}, nil, "T2", nil)
	T2 := NewNamed(n2, nil, nil)
	s1 := NewStruct([]*Var{NewField(syntax.Pos{}, nil, "_", T2, false)}, nil)
	T1.SetUnderlying(s1)
	s2 := NewStruct([]*Var{NewField(syntax.Pos{}, nil, "_", T2, false)}, nil)
	s3 := NewStruct([]*Var{NewField(syntax.Pos{}, nil, "_", s2, false)}, nil)
	T2.SetUnderlying(s3)

	// These calls must terminate (no endless recursion).
	Comparable(T1)
	Comparable(T2)
}

func TestIssue44515(t *testing.T) {
	typ := Unsafe.Scope().Lookup("Pointer").Type()

	got := TypeString(typ, nil)
	want := "unsafe.Pointer"
	if got != want {
		t.Errorf("got %q; want %q", got, want)
	}

	qf := func(pkg *Package) string {
		if pkg == Unsafe {
			return "foo"
		}
		return ""
	}
	got = TypeString(typ, qf)
	want = "foo.Pointer"
	if got != want {
		t.Errorf("got %q; want %q", got, want)
	}
}

func TestIssue43124(t *testing.T) {
	// All involved packages have the same name (template). Error messages should
	// disambiguate between text/template and html/template by printing the full
	// path.
	const (
		asrc = `package a; import "text/template"; func F(template.Template) {}; func G(int) {}`
		bsrc = `package b; import ("a"; "html/template"); func _() { a.F(template.Template{}) }`
		csrc = `package c; import ("a"; "html/template"); func _() { a.G(template.Template{}) }`
	)

	a, err := pkgFor("a", asrc, nil)
	if err != nil {
		t.Fatalf("package a failed to typecheck: %v", err)
	}
	conf := Config{Importer: importHelper{pkg: a, fallback: defaultImporter()}}

	// Packages should be fully qualified when there is ambiguity within the
	// error string itself.
	bast := mustParse(t, bsrc)
	_, err = conf.Check(bast.PkgName.Value, []*syntax.File{bast}, nil)
	if err == nil {
		t.Fatal("package b had no errors")
	}
	if !strings.Contains(err.Error(), "text/template") || !strings.Contains(err.Error(), "html/template") {
		t.Errorf("type checking error for b does not disambiguate package template: %q", err)
	}

	// ...and also when there is any ambiguity in reachable packages.
	cast := mustParse(t, csrc)
	_, err = conf.Check(cast.PkgName.Value, []*syntax.File{cast}, nil)
	if err == nil {
		t.Fatal("package c had no errors")
	}
	if !strings.Contains(err.Error(), "html/template") {
		t.Errorf("type checking error for c does not disambiguate package template: %q", err)
	}
}

func TestIssue50646(t *testing.T) {
	anyType := Universe.Lookup("any").Type()
	comparableType := Universe.Lookup("comparable").Type()

	if !Comparable(anyType) {
		t.Errorf("any is not a comparable type")
	}
	if !Comparable(comparableType) {
		t.Errorf("comparable is not a comparable type")
	}

	if Implements(anyType, comparableType.Underlying().(*Interface)) {
		t.Errorf("any implements comparable")
	}
	if !Implements(comparableType, anyType.(*Interface)) {
		t.Errorf("comparable does not implement any")
	}

	if AssignableTo(anyType, comparableType) {
		t.Errorf("any assignable to comparable")
	}
	if !AssignableTo(comparableType, anyType) {
		t.Errorf("comparable not assignable to any")
	}
}

func TestIssue55030(t *testing.T) {
	// makeSig makes the signature func(typ...)
	makeSig := func(typ Type) {
		par := NewVar(nopos, nil, "", typ)
		params := NewTuple(par)
		NewSignatureType(nil, nil, nil, params, nil, true)
	}

	// makeSig must not panic for the following (example) types:
	// []int
	makeSig(NewSlice(Typ[Int]))

	// string
	makeSig(Typ[String])

	// P where P's core type is string
	{
		P := NewTypeName(nopos, nil, "P", nil) // [P string]
		makeSig(NewTypeParam(P, NewInterfaceType(nil, []Type{Typ[String]})))
	}

	// P where P's core type is an (unnamed) slice
	{
		P := NewTypeName(nopos, nil, "P", nil) // [P []int]
		makeSig(NewTypeParam(P, NewInterfaceType(nil, []Type{NewSlice(Typ[Int])})))
	}

	// P where P's core type is bytestring (i.e., string or []byte)
	{
		t1 := NewTerm(true, Typ[String])          // ~string
		t2 := NewTerm(false, NewSlice(Typ[Byte])) // []byte
		u := NewUnion([]*Term{t1, t2})            // ~string | []byte
		P := NewTypeName(nopos, nil, "P", nil)    // [P ~string | []byte]
		makeSig(NewTypeParam(P, NewInterfaceType(nil, []Type{u})))
	}
}
