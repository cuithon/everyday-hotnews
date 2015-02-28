// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ld

import (
	"bytes"
	"cmd/internal/obj"
	"fmt"
	"os"
	"strings"
	"unicode/utf8"
)

// go-specific code shared across loaders (5l, 6l, 8l).

// replace all "". with pkg.
func expandpkg(t0 string, pkg string) string {
	return strings.Replace(t0, `"".`, pkg+".", -1)
}

// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// go-specific code shared across loaders (5l, 6l, 8l).

// accumulate all type information from .6 files.
// check for inconsistencies.

// TODO:
//	generate debugging section in binary.
//	once the dust settles, try to move some code to
//		libmach, so that other linkers and ar can share.

/*
 *	package import data
 */
type Import struct {
	hash   *Import
	prefix string
	name   string
	def    string
	file   string
}

const (
	NIHASH = 1024
)

var ihash [NIHASH]*Import

var nimport int

func hashstr(name string) int {
	var h uint32
	var cp string

	h = 0
	for cp = name; cp != ""; cp = cp[1:] {
		h = h*1119 + uint32(cp[0])
	}
	h &= 0xffffff
	return int(h)
}

func ilookup(name string) *Import {
	var h int
	var x *Import

	h = hashstr(name) % NIHASH
	for x = ihash[h]; x != nil; x = x.hash {
		if x.name[0] == name[0] && x.name == name {
			return x
		}
	}
	x = new(Import)
	x.name = name
	x.hash = ihash[h]
	ihash[h] = x
	nimport++
	return x
}

func ldpkg(f *Biobuf, pkg string, length int64, filename string, whence int) {
	var bdata []byte
	var data string
	var p0, p1 int
	var name string

	if Debug['g'] != 0 {
		return
	}

	if int64(int(length)) != length {
		fmt.Fprintf(os.Stderr, "%s: too much pkg data in %s\n", os.Args[0], filename)
		if Debug['u'] != 0 {
			Errorexit()
		}
		return
	}

	bdata = make([]byte, length)
	if int64(Bread(f, bdata)) != length {
		fmt.Fprintf(os.Stderr, "%s: short pkg read %s\n", os.Args[0], filename)
		if Debug['u'] != 0 {
			Errorexit()
		}
		return
	}
	data = string(bdata)

	// first \n$$ marks beginning of exports - skip rest of line
	p0 = strings.Index(data, "\n$$")
	if p0 < 0 {
		if Debug['u'] != 0 && whence != ArchiveObj {
			fmt.Fprintf(os.Stderr, "%s: cannot find export data in %s\n", os.Args[0], filename)
			Errorexit()
		}
		return
	}

	p0 += 3
	for p0 < len(data) && data[0] != '\n' {
		p0++
	}

	// second marks end of exports / beginning of local data
	p1 = strings.Index(data[p0:], "\n$$")
	if p1 < 0 {
		fmt.Fprintf(os.Stderr, "%s: cannot find end of exports in %s\n", os.Args[0], filename)
		if Debug['u'] != 0 {
			Errorexit()
		}
		return
	}
	p1 += p0

	for p0 < p1 && (data[p0] == ' ' || data[0] == '\t' || data[0] == '\n') {
		p0++
	}
	if p0 < p1 {
		if !strings.HasPrefix(data[p0:], "package ") {
			fmt.Fprintf(os.Stderr, "%s: bad package section in %s - %.20s\n", os.Args[0], filename, data[p0:])
			if Debug['u'] != 0 {
				Errorexit()
			}
			return
		}

		p0 += 8
		for p0 < p1 && (data[p0] == ' ' || data[p0] == '\t' || data[p0] == '\n') {
			p0++
		}
		name = data[p0:]
		for p0 < p1 && data[p0] != ' ' && data[p0] != '\t' && data[p0] != '\n' {
			p0++
		}
		if Debug['u'] != 0 && whence != ArchiveObj && (p0+6 > p1 || !strings.HasPrefix(data[p0:], " safe\n")) {
			fmt.Fprintf(os.Stderr, "%s: load of unsafe package %s\n", os.Args[0], filename)
			nerrors++
			Errorexit()
		}

		name = name[:p1-p0]
		if p0 < p1 {
			if data[p0] == '\n' {
				p0++
			} else {
				p0++
				for p0 < p1 && data[p0] != '\n' {
					p0++
				}
			}
		}

		if pkg == "main" && name != "main" {
			fmt.Fprintf(os.Stderr, "%s: %s: not package main (package %s)\n", os.Args[0], filename, name)
			nerrors++
			Errorexit()
		}

		loadpkgdata(filename, pkg, data[p0:p1])
	}

	// __.PKGDEF has no cgo section - those are in the C compiler-generated object files.
	if whence == Pkgdef {
		return
	}

	// look for cgo section
	p0 = strings.Index(data[p1:], "\n$$  // cgo")
	if p0 >= 0 {
		p0 += p1
		i := strings.IndexByte(data[p0+1:], '\n')
		if i < 0 {
			fmt.Fprintf(os.Stderr, "%s: found $$ // cgo but no newline in %s\n", os.Args[0], filename)
			if Debug['u'] != 0 {
				Errorexit()
			}
			return
		}
		p0 += 1 + i

		p1 = strings.Index(data[p0:], "\n$$")
		if p1 < 0 {
			p1 = strings.Index(data[p0:], "\n!\n")
		}
		if p1 < 0 {
			fmt.Fprintf(os.Stderr, "%s: cannot find end of // cgo section in %s\n", os.Args[0], filename)
			if Debug['u'] != 0 {
				Errorexit()
			}
			return
		}
		p1 += p0

		loadcgo(filename, pkg, data[p0:p1])
	}
}

func loadpkgdata(file string, pkg string, data string) {
	var p string
	var prefix string
	var name string
	var def string
	var x *Import

	file = file
	p = data
	for parsepkgdata(file, pkg, &p, &prefix, &name, &def) > 0 {
		x = ilookup(name)
		if x.prefix == "" {
			x.prefix = prefix
			x.def = def
			x.file = file
		} else if x.prefix != prefix {
			fmt.Fprintf(os.Stderr, "%s: conflicting definitions for %s\n", os.Args[0], name)
			fmt.Fprintf(os.Stderr, "%s:\t%s %s ...\n", x.file, x.prefix, name)
			fmt.Fprintf(os.Stderr, "%s:\t%s %s ...\n", file, prefix, name)
			nerrors++
		} else if x.def != def {
			fmt.Fprintf(os.Stderr, "%s: conflicting definitions for %s\n", os.Args[0], name)
			fmt.Fprintf(os.Stderr, "%s:\t%s %s %s\n", x.file, x.prefix, name, x.def)
			fmt.Fprintf(os.Stderr, "%s:\t%s %s %s\n", file, prefix, name, def)
			nerrors++
		}
	}
}

func parsepkgdata(file string, pkg string, pp *string, prefixp *string, namep *string, defp *string) int {
	var p string
	var prefix string
	var name string
	var def string
	var meth string
	var inquote bool

	// skip white space
	p = *pp

loop:
	for len(p) > 0 && (p[0] == ' ' || p[0] == '\t' || p[0] == '\n') {
		p = p[1:]
	}
	if len(p) == 0 || strings.HasPrefix(p, "$$\n") {
		return 0
	}

	// prefix: (var|type|func|const)
	prefix = p

	if len(p) < 7 {
		return -1
	}
	if strings.HasPrefix(p, "var ") {
		p = p[4:]
	} else if strings.HasPrefix(p, "type ") {
		p = p[5:]
	} else if strings.HasPrefix(p, "func ") {
		p = p[5:]
	} else if strings.HasPrefix(p, "const ") {
		p = p[6:]
	} else if strings.HasPrefix(p, "import ") {
		p = p[7:]
		for len(p) > 0 && p[0] != ' ' {
			p = p[1:]
		}
		p = p[1:]
		name := p
		for len(p) > 0 && p[0] != '\n' {
			p = p[1:]
		}
		if len(p) == 0 {
			fmt.Fprintf(os.Stderr, "%s: %s: confused in import line\n", os.Args[0], file)
			nerrors++
			return -1
		}
		name = name[:len(name)-len(p)]
		p = p[1:]
		imported(pkg, name)
		goto loop
	} else {
		fmt.Fprintf(os.Stderr, "%s: %s: confused in pkg data near <<%.40s>>\n", os.Args[0], file, prefix)
		nerrors++
		return -1
	}

	prefix = prefix[:len(prefix)-len(p)-1]

	// name: a.b followed by space
	name = p

	inquote = false
	for len(p) > 0 {
		if p[0] == ' ' && !inquote {
			break
		}

		if p[0] == '\\' {
			p = p[1:]
		} else if p[0] == '"' {
			inquote = !inquote
		}

		p = p[1:]
	}

	if len(p) == 0 {
		return -1
	}
	name = name[:len(name)-len(p)]
	p = p[1:]

	// def: free form to new line
	def = p

	for len(p) > 0 && p[0] != '\n' {
		p = p[1:]
	}
	if len(p) == 0 {
		return -1
	}
	def = def[:len(def)-len(p)]
	var defbuf *bytes.Buffer
	p = p[1:]

	// include methods on successive lines in def of named type
	for parsemethod(&p, &meth) > 0 {
		if defbuf == nil {
			defbuf = new(bytes.Buffer)
			defbuf.WriteString(def)
		}
		defbuf.WriteString("\n\t")
		defbuf.WriteString(meth)
	}
	if defbuf != nil {
		def = defbuf.String()
	}

	name = expandpkg(name, pkg)
	def = expandpkg(def, pkg)

	// done
	*pp = p

	*prefixp = prefix
	*namep = name
	*defp = def
	return 1
}

func parsemethod(pp *string, methp *string) int {
	var p string

	// skip white space
	p = *pp

	for len(p) > 0 && (p[0] == ' ' || p[0] == '\t') {
		p = p[1:]
	}
	if len(p) == 0 {
		return 0
	}

	// might be a comment about the method
	if strings.HasPrefix(p, "//") {
		goto useline
	}

	// if it says "func (", it's a method
	if strings.HasPrefix(p, "func (") {
		goto useline
	}
	return 0

	// definition to end of line
useline:
	*methp = p

	for len(p) > 0 && p[0] != '\n' {
		p = p[1:]
	}
	if len(p) == 0 {
		fmt.Fprintf(os.Stderr, "%s: lost end of line in method definition\n", os.Args[0])
		*pp = ""
		return -1
	}

	*methp = (*methp)[:len(*methp)-len(p)]
	*pp = p[1:]
	return 1
}

func loadcgo(file string, pkg string, p string) {
	var next string
	var p0 string
	var q string
	var f []string
	var local string
	var remote string
	var lib string
	var s *LSym

	p0 = ""
	for ; p != ""; p = next {
		if i := strings.Index(p, "\n"); i >= 0 {
			p, next = p[:i], p[i+1:]
		} else {
			next = ""
		}

		p0 = p // save for error message
		f = tokenize(p)
		if len(f) == 0 {
			continue
		}

		if f[0] == "cgo_import_dynamic" {
			if len(f) < 2 || len(f) > 4 {
				goto err
			}

			local = f[1]
			remote = local
			if len(f) > 2 {
				remote = f[2]
			}
			lib = ""
			if len(f) > 3 {
				lib = f[3]
			}

			if Debug['d'] != 0 {
				fmt.Fprintf(os.Stderr, "%s: %s: cannot use dynamic imports with -d flag\n", os.Args[0], file)
				nerrors++
				return
			}

			if local == "_" && remote == "_" {
				// allow #pragma dynimport _ _ "foo.so"
				// to force a link of foo.so.
				havedynamic = 1

				Thearch.Adddynlib(lib)
				continue
			}

			local = expandpkg(local, pkg)
			q = ""
			if i := strings.Index(remote, "#"); i >= 0 {
				remote, q = remote[:i], remote[i+1:]
			}
			s = Linklookup(Ctxt, local, 0)
			if local != f[1] {
			}
			if s.Type == 0 || s.Type == SXREF || s.Type == SHOSTOBJ {
				s.Dynimplib = lib
				s.Extname = remote
				s.Dynimpvers = q
				if s.Type != SHOSTOBJ {
					s.Type = SDYNIMPORT
				}
				havedynamic = 1
			}

			continue
		}

		if f[0] == "cgo_import_static" {
			if len(f) != 2 {
				goto err
			}
			local = f[1]
			s = Linklookup(Ctxt, local, 0)
			s.Type = SHOSTOBJ
			s.Size = 0
			continue
		}

		if f[0] == "cgo_export_static" || f[0] == "cgo_export_dynamic" {
			// TODO: Remove once we know Windows is okay.
			if f[0] == "cgo_export_static" && HEADTYPE == Hwindows {
				continue
			}

			if len(f) < 2 || len(f) > 3 {
				goto err
			}
			local = f[1]
			if len(f) > 2 {
				remote = f[2]
			} else {
				remote = local
			}
			local = expandpkg(local, pkg)
			s = Linklookup(Ctxt, local, 0)

			if Flag_shared != 0 && s == Linklookup(Ctxt, "main", 0) {
				continue
			}

			// export overrides import, for openbsd/cgo.
			// see issue 4878.
			if s.Dynimplib != "" {
				s.Dynimplib = ""
				s.Extname = ""
				s.Dynimpvers = ""
				s.Type = 0
			}

			if s.Cgoexport == 0 {
				s.Extname = remote
				dynexp = append(dynexp, s)
			} else if s.Extname != remote {
				fmt.Fprintf(os.Stderr, "%s: conflicting cgo_export directives: %s as %s and %s\n", os.Args[0], s.Name, s.Extname, remote)
				nerrors++
				return
			}

			if f[0] == "cgo_export_static" {
				s.Cgoexport |= CgoExportStatic
			} else {
				s.Cgoexport |= CgoExportDynamic
			}
			if local != f[1] {
			}
			continue
		}

		if f[0] == "cgo_dynamic_linker" {
			if len(f) != 2 {
				goto err
			}

			if Debug['I'] == 0 {
				if interpreter != "" && interpreter != f[1] {
					fmt.Fprintf(os.Stderr, "%s: conflict dynlinker: %s and %s\n", os.Args[0], interpreter, f[1])
					nerrors++
					return
				}

				interpreter = f[1]
			}

			continue
		}

		if f[0] == "cgo_ldflag" {
			if len(f) != 2 {
				goto err
			}
			ldflag = append(ldflag, f[1])
			continue
		}
	}

	return

err:
	fmt.Fprintf(os.Stderr, "%s: %s: invalid dynimport line: %s\n", os.Args[0], file, p0)
	nerrors++
}

var markq *LSym

var emarkq *LSym

func mark1(s *LSym, parent *LSym) {
	if s == nil || s.Reachable {
		return
	}
	if strings.HasPrefix(s.Name, "go.weak.") {
		return
	}
	s.Reachable = true
	s.Reachparent = parent
	if markq == nil {
		markq = s
	} else {
		emarkq.Queue = s
	}
	emarkq = s
}

func mark(s *LSym) {
	mark1(s, nil)
}

func markflood() {
	var a *Auto
	var s *LSym
	var i int

	for s = markq; s != nil; s = s.Queue {
		if s.Type == STEXT {
			if Debug['v'] > 1 {
				fmt.Fprintf(&Bso, "marktext %s\n", s.Name)
			}
			for a = s.Autom; a != nil; a = a.Link {
				mark1(a.Gotype, s)
			}
		}

		for i = 0; i < len(s.R); i++ {
			mark1(s.R[i].Sym, s)
		}
		if s.Pcln != nil {
			for i = 0; i < s.Pcln.Nfuncdata; i++ {
				mark1(s.Pcln.Funcdata[i], s)
			}
		}

		mark1(s.Gotype, s)
		mark1(s.Sub, s)
		mark1(s.Outer, s)
	}
}

var markextra = []string{
	"runtime.morestack",
	"runtime.morestackx",
	"runtime.morestack00",
	"runtime.morestack10",
	"runtime.morestack01",
	"runtime.morestack11",
	"runtime.morestack8",
	"runtime.morestack16",
	"runtime.morestack24",
	"runtime.morestack32",
	"runtime.morestack40",
	"runtime.morestack48",
	// on arm, lock in the div/mod helpers too
	"_div",
	"_divu",
	"_mod",
	"_modu",
}

func deadcode() {
	var i int
	var s *LSym
	var last *LSym
	var p *LSym
	var fmt_ string

	if Debug['v'] != 0 {
		fmt.Fprintf(&Bso, "%5.2f deadcode\n", obj.Cputime())
	}

	mark(Linklookup(Ctxt, INITENTRY, 0))
	for i = 0; i < len(markextra); i++ {
		mark(Linklookup(Ctxt, markextra[i], 0))
	}

	for i = 0; i < len(dynexp); i++ {
		mark(dynexp[i])
	}

	markflood()

	// keep each beginning with 'typelink.' if the symbol it points at is being kept.
	for s = Ctxt.Allsym; s != nil; s = s.Allsym {
		if strings.HasPrefix(s.Name, "go.typelink.") {
			s.Reachable = len(s.R) == 1 && s.R[0].Sym.Reachable
		}
	}

	// remove dead text but keep file information (z symbols).
	last = nil

	for s = Ctxt.Textp; s != nil; s = s.Next {
		if !s.Reachable {
			continue
		}

		// NOTE: Removing s from old textp and adding to new, shorter textp.
		if last == nil {
			Ctxt.Textp = s
		} else {
			last.Next = s
		}
		last = s
	}

	if last == nil {
		Ctxt.Textp = nil
	} else {
		last.Next = nil
	}

	for s = Ctxt.Allsym; s != nil; s = s.Allsym {
		if strings.HasPrefix(s.Name, "go.weak.") {
			s.Special = 1 // do not lay out in data segment
			s.Reachable = true
			s.Hide = 1
		}
	}

	// record field tracking references
	fmt_ = ""

	for s = Ctxt.Allsym; s != nil; s = s.Allsym {
		if strings.HasPrefix(s.Name, "go.track.") {
			s.Special = 1 // do not lay out in data segment
			s.Hide = 1
			if s.Reachable {
				fmt_ += fmt.Sprintf("%s", s.Name[9:])
				for p = s.Reachparent; p != nil; p = p.Reachparent {
					fmt_ += fmt.Sprintf("\t%s", p.Name)
				}
				fmt_ += fmt.Sprintf("\n")
			}

			s.Type = SCONST
			s.Value = 0
		}
	}

	if tracksym == "" {
		return
	}
	s = Linklookup(Ctxt, tracksym, 0)
	if !s.Reachable {
		return
	}
	addstrdata(tracksym, fmt_)
}

func doweak() {
	var s *LSym
	var t *LSym

	// resolve weak references only if
	// target symbol will be in binary anyway.
	for s = Ctxt.Allsym; s != nil; s = s.Allsym {
		if strings.HasPrefix(s.Name, "go.weak.") {
			t = Linkrlookup(Ctxt, s.Name[8:], int(s.Version))
			if t != nil && t.Type != 0 && t.Reachable {
				s.Value = t.Value
				s.Type = t.Type
				s.Outer = t
			} else {
				s.Type = SCONST
				s.Value = 0
			}

			continue
		}
	}
}

func addexport() {
	var i int

	if HEADTYPE == Hdarwin {
		return
	}

	for i = 0; i < len(dynexp); i++ {
		Thearch.Adddynsym(Ctxt, dynexp[i])
	}
}

/* %Z from gc, for quoting import paths */
func Zconv(s string, flag int) string {
	// NOTE: Keep in sync with gc Zconv.
	var n int
	var fp string
	for i := 0; i < len(s); i += n {
		var r rune
		r, n = utf8.DecodeRuneInString(s[i:])
		switch r {
		case utf8.RuneError:
			if n == 1 {
				fp += fmt.Sprintf("\\x%02x", s[i])
				break
			}
			fallthrough

			// fall through
		default:
			if r < ' ' {
				fp += fmt.Sprintf("\\x%02x", r)
				break
			}

			fp += string(r)

		case '\t':
			fp += "\\t"

		case '\n':
			fp += "\\n"

		case '"',
			'\\':
			fp += `\` + string(r)

		case 0xFEFF: // BOM, basically disallowed in source code
			fp += "\\uFEFF"
		}
	}

	return fp
}

type Pkg struct {
	mark    uint8
	checked uint8
	next    *Pkg
	path_   string
	impby   []*Pkg
	all     *Pkg
}

var phash [1024]*Pkg

var pkgall *Pkg

func getpkg(path_ string) *Pkg {
	var p *Pkg
	var h int

	h = hashstr(path_) % len(phash)
	for p = phash[h]; p != nil; p = p.next {
		if p.path_ == path_ {
			return p
		}
	}
	p = new(Pkg)
	p.path_ = path_
	p.next = phash[h]
	phash[h] = p
	p.all = pkgall
	pkgall = p
	return p
}

func imported(pkg string, import_ string) {
	var p *Pkg
	var i *Pkg

	// everyone imports runtime, even runtime.
	if import_ == "\"runtime\"" {
		return
	}

	pkg = fmt.Sprintf("\"%v\"", Zconv(pkg, 0)) // turn pkg path into quoted form, freed below
	p = getpkg(pkg)
	i = getpkg(import_)
	i.impby = append(i.impby, p)
}

func cycle(p *Pkg) *Pkg {
	var i int
	var bad *Pkg

	if p.checked != 0 {
		return nil
	}

	if p.mark != 0 {
		nerrors++
		fmt.Printf("import cycle:\n")
		fmt.Printf("\t%s\n", p.path_)
		return p
	}

	p.mark = 1
	for i = 0; i < len(p.impby); i++ {
		bad = cycle(p.impby[i])
		if bad != nil {
			p.mark = 0
			p.checked = 1
			fmt.Printf("\timports %s\n", p.path_)
			if bad == p {
				return nil
			}
			return bad
		}
	}

	p.checked = 1
	p.mark = 0
	return nil
}

func importcycles() {
	var p *Pkg

	for p = pkgall; p != nil; p = p.all {
		cycle(p)
	}
}

func setlinkmode(arg string) {
	if arg == "internal" {
		Linkmode = LinkInternal
	} else if arg == "external" {
		Linkmode = LinkExternal
	} else if arg == "auto" {
		Linkmode = LinkAuto
	} else {
		fmt.Fprintf(os.Stderr, "unknown link mode -linkmode %s\n", arg)
		Errorexit()
	}
}
