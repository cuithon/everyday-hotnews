// Copyright 2009 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gobuild

import (
	"fmt";
	"gobuild";
	"io";
	"path";
	"template";
)

var makefileTemplate =
	"# DO NOT EDIT.  Automatically generated by gobuild.\n"
	"{Args|args} >Makefile\n"
	"\n"
	"D={.section Dir}/{@}{.end}\n"
	"\n"
	"O_arm=5\n"	// TODO(rsc): include something here?
	"O_amd64=6\n"
	"O_386=8\n"
	"OS=568vq\n"
	"\n"
	"O=$(O_$(GOARCH))\n"
	"GC=$(O)g -I{ObjDir}\n"
	"CC=$(O)c -FVw\n"
	"AS=$(O)a\n"
	"AR=6ar\n"
	"\n"
	"default: packages\n"
	"\n"
	"clean:\n"
	"	rm -rf *.[$(OS)] *.a [$(OS)].out {ObjDir}\n"
	"\n"
	"test: packages\n"
	"	gotest\n"
	"\n"
	"coverage: packages\n"
	"	gotest\n"
	"	6cov -g `pwd` | grep -v '_test\\.go:'\n"
	"\n"
	"%.$O: %.go\n"
	"	$(GC) $*.go\n"
	"\n"
	"%.$O: %.c\n"
	"	$(CC) $*.c\n"
	"\n"
	"%.$O: %.s\n"
	"	$(AS) $*.s\n"
	"\n"
	"{.repeated section Phases}\n"
	"O{Phase}=\\\n"
	"{.repeated section ArCmds}\n"
	"{.repeated section Files}\n"
	"	{Name|basename}.$O\\\n"
	"{.end}\n"
	"{.end}\n"
	"\n"
	"{.end}\n"
	"\n"
	"phases:{.repeated section Phases} a{Phase}{.end}\n"
	"{.repeated section Packages}\n"
	"{ObjDir}$D/{Name}.a: phases\n"
	"{.end}\n"
	"\n"
	"{.repeated section Phases}\n"
	"a{Phase}: $(O{Phase})\n"
	"{.repeated section ArCmds}\n"
	"	$(AR) grc {ObjDir}$D/{.section Pkg}{Name}.a{.end}{.repeated section Files} {Name|basename}.$O{.end}\n"
	"{.end}\n"
	"	rm -f $(O{Phase})\n"
	"\n"
	"{.end}\n"
	"\n"
	"newpkg: clean\n"
	"	mkdir -p {ObjDir}$D\n"
	"{.repeated section Packages}\n"
	"	$(AR) grc {ObjDir}$D/{Name}.a\n"
	"{.end}\n"
	"\n"
	"$(O1): newpkg\n"
	"{.repeated section Phases}\n"
	"$(O{Phase|+1}): a{Phase}\n"
	"{.end}\n"
	"\n"
	"nuke: clean\n"
	"	rm -f{.repeated section Packages} $(GOROOT)/pkg$D/{Name}.a{.end}\n"
	"\n"
	"packages:{.repeated section Packages} {ObjDir}$D/{Name}.a{.end}\n"
	"\n"
	"install: packages\n"
	"	test -d $(GOROOT)/pkg && mkdir -p $(GOROOT)/pkg$D\n"
	"{.repeated section Packages}\n"
	"	cp {ObjDir}$D/{Name}.a $(GOROOT)/pkg$D/{Name}.a\n"
	"{.end}\n"

func argsFmt(w io.Write, x interface{}, format string) {
	args := x.([]string);
	fmt.Fprint(w, "#");
	for i, a := range args {
		fmt.Fprint(w, " ", ShellString(a));
	}
}

func basenameFmt(w io.Write, x interface{}, format string) {
	t := fmt.Sprint(x);
	t = t[0:len(t)-len(path.Ext(t))];
	fmt.Fprint(w, MakeString(t));
}

func plus1Fmt(w io.Write, x interface{}, format string) {
	fmt.Fprint(w, x.(int) + 1);
}

func makeFmt(w io.Write, x interface{}, format string) {
	fmt.Fprint(w, MakeString(fmt.Sprint(x)));
}

var makefileMap = template.FormatterMap {
	"": makeFmt,
	"+1": plus1Fmt,
	"args": argsFmt,
	"basename": basenameFmt,
}
