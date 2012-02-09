// Copyright 2012 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "a.h"
#include <stdio.h>

/*
 * Helpers for building pkg/runtime.
 */

// mkzversion writes zversion.go:
//
//	package runtime
//	const defaultGoroot = <goroot>
//	const theVersion = <version>
//
void
mkzversion(char *dir, char *file)
{
	Buf b, out;
	
	binit(&b);
	binit(&out);
	
	bwritestr(&out, bprintf(&b,
		"// auto generated by go tool dist\n"
		"\n"
		"package runtime\n"
		"\n"
		"const defaultGoroot = `%s`\n"
		"const theVersion = `%s`\n", goroot_final, goversion));

	writefile(&out, file);
	
	bfree(&b);
	bfree(&out);
}

// mkzgoarch writes zgoarch_$GOARCH.go:
//
//	package runtime
//	const theGoarch = <goarch>
//
void
mkzgoarch(char *dir, char *file)
{
	Buf b, out;
	
	binit(&b);
	binit(&out);
	
	bwritestr(&out, bprintf(&b,
		"// auto generated by go tool dist\n"
		"\n"
		"package runtime\n"
		"\n"
		"const theGoarch = `%s`\n", goarch));

	writefile(&out, file);
	
	bfree(&b);
	bfree(&out);
}

// mkzgoos writes zgoos_$GOOS.go:
//
//	package runtime
//	const theGoos = <goos>
//
void
mkzgoos(char *dir, char *file)
{
	Buf b, out;
	
	binit(&b);
	binit(&out);
	
	bwritestr(&out, bprintf(&b,
		"// auto generated by go tool dist\n"
		"\n"
		"package runtime\n"
		"\n"
		"const theGoos = `%s`\n", goos));

	writefile(&out, file);
	
	bfree(&b);
	bfree(&out);
}

static struct {
	char *goarch;
	char *goos;
	char *hdr;
} zasmhdr[] = {
	{"386", "windows",
		"#define	get_tls(r)	MOVL 0x14(FS), r\n"
		"#define	g(r)	0(r)\n"
		"#define	m(r)	4(r)\n"
	},
	{"386", "plan9",
		"#define	get_tls(r)	MOVL _tos(SB), r \n"
		"#define	g(r)	-8(r)\n"
		"#define	m(r)	-4(r)\n"
	},
	{"386", "linux",
		"// On Linux systems, what we call 0(GS) and 4(GS) for g and m\n"
		"// turn into %gs:-8 and %gs:-4 (using gcc syntax to denote\n"
		"// what the machine sees as opposed to 8l input).\n"
		"// 8l rewrites 0(GS) and 4(GS) into these.\n"
		"//\n"
		"// On Linux Xen, it is not allowed to use %gs:-8 and %gs:-4\n"
		"// directly.  Instead, we have to store %gs:0 into a temporary\n"
		"// register and then use -8(%reg) and -4(%reg).  This kind\n"
		"// of addressing is correct even when not running Xen.\n"
		"//\n"
		"// 8l can rewrite MOVL 0(GS), CX into the appropriate pair\n"
		"// of mov instructions, using CX as the intermediate register\n"
		"// (safe because CX is about to be written to anyway).\n"
		"// But 8l cannot handle other instructions, like storing into 0(GS),\n"
		"// which is where these macros come into play.\n"
		"// get_tls sets up the temporary and then g and r use it.\n"
		"//\n"
		"// The final wrinkle is that get_tls needs to read from %gs:0,\n"
		"// but in 8l input it's called 8(GS), because 8l is going to\n"
		"// subtract 8 from all the offsets, as described above.\n"
		"#define	get_tls(r)	MOVL 8(GS), r\n"
		"#define	g(r)	-8(r)\n"
		"#define	m(r)	-4(r)\n"
	},
	{"386", "",
		"#define	get_tls(r)\n"
		"#define	g(r)	0(GS)\n"
		"#define	m(r)	4(GS)\n"
	},
	
	{"amd64", "windows",
		"#define	get_tls(r) MOVQ 0x28(GS), r\n"
		"#define	g(r) 0(r)\n"
		"#define	m(r) 8(r)\n"
	},
	{"amd64", "",
		"// The offsets 0 and 8 are known to:\n"
		"//	../../cmd/6l/pass.c:/D_GS\n"
		"//	cgo/gcc_linux_amd64.c:/^threadentry\n"
		"//	cgo/gcc_darwin_amd64.c:/^threadentry\n"
		"//\n"
		"#define	get_tls(r)\n"
		"#define	g(r) 0(GS)\n"
		"#define	m(r) 8(GS)\n"
	},
	
	{"arm", "",
	"#define	g	R10\n"
	"#define	m	R9\n"
	"#define	LR	R14\n"
	},
};

// mkzasm writes zasm_$GOOS_$GOARCH.h,
// which contains struct offsets for use by
// assembly files.  It also writes a copy to the work space
// under the name zasm_GOOS_GOARCH.h (no expansion).
// 
void
mkzasm(char *dir, char *file)
{
	int i, n;
	char *aggr, *p;
	Buf in, b, out;
	Vec argv, lines, fields;

	binit(&in);
	binit(&b);
	binit(&out);
	vinit(&argv);
	vinit(&lines);
	vinit(&fields);
	
	bwritestr(&out, "// auto generated by go tool dist\n\n");
	for(i=0; i<nelem(zasmhdr); i++) {
		if(hasprefix(goarch, zasmhdr[i].goarch) && hasprefix(goos, zasmhdr[i].goos)) {
			bwritestr(&out, zasmhdr[i].hdr);
			goto ok;
		}
	}
	fatal("unknown $GOOS/$GOARCH in mkzasm");
ok:

	// Run 6c -DGOOS_goos -DGOARCH_goarch -Iworkdir -a proc.c
	// to get acid [sic] output.
	vreset(&argv);
	vadd(&argv, bpathf(&b, "%s/bin/tool/%sc", goroot, gochar));
	vadd(&argv, bprintf(&b, "-DGOOS_%s", goos));
	vadd(&argv, bprintf(&b, "-DGOARCH_%s", goarch));
	vadd(&argv, bprintf(&b, "-I%s", workdir));
	vadd(&argv, "-a");
	vadd(&argv, "proc.c");
	runv(&in, dir, CheckExit, &argv);
	
	// Convert input like
	//	aggr G
	//	{
	//		Gobuf 24 sched;
	//		'Y' 48 stack0;
	//	}
	// into output like
	//	#define g_sched 24
	//	#define g_stack0 48
	//
	aggr = nil;
	splitlines(&lines, bstr(&in));
	for(i=0; i<lines.len; i++) {
		splitfields(&fields, lines.p[i]);
		if(fields.len == 2 && streq(fields.p[0], "aggr")) {
			if(streq(fields.p[1], "G"))
				aggr = "g";
			else if(streq(fields.p[1], "M"))
				aggr = "m";
			else if(streq(fields.p[1], "Gobuf"))
				aggr = "gobuf";
			else if(streq(fields.p[1], "WinCall"))
				aggr = "wincall";
		}
		if(hasprefix(lines.p[i], "}"))
			aggr = nil;
		if(aggr && hasprefix(lines.p[i], "\t") && fields.len >= 2) {
			n = fields.len;
			p = fields.p[n-1];
			if(p[xstrlen(p)-1] == ';')
				p[xstrlen(p)-1] = '\0';
			bwritestr(&out, bprintf(&b, "#define %s_%s %s\n", aggr, fields.p[n-1], fields.p[n-2]));
		}
	}
	
	// Write both to file and to workdir/zasm_GOOS_GOARCH.h.
	writefile(&out, file);
	writefile(&out, bprintf(&b, "%s/zasm_GOOS_GOARCH.h", workdir));

	bfree(&in);
	bfree(&b);
	bfree(&out);
	vfree(&argv);
	vfree(&lines);
	vfree(&fields);
}

static char *runtimedefs[] = {
	"proc.c",
	"iface.c",
	"hashmap.c",
	"chan.c",
};

// mkzruntimedefs writes zruntime_defs_$GOOS_$GOARCH.h,
// which contains Go struct definitions equivalent to the C ones.
// Mostly we just write the output of 6c -q to the file.
// However, we run it on multiple files, so we have to delete
// the duplicated definitions, and we don't care about the funcs
// and consts, so we delete those too.
// 
void
mkzruntimedefs(char *dir, char *file)
{
	int i, skip;
	char *p;
	Buf in, b, out;
	Vec argv, lines, fields, seen;
	
	binit(&in);
	binit(&b);
	binit(&out);
	vinit(&argv);
	vinit(&lines);
	vinit(&fields);
	vinit(&seen);
	
	bwritestr(&out, "// auto generated by go tool dist\n"
		"\n"
		"package runtime\n"
		"import \"unsafe\"\n"
		"var _ unsafe.Pointer\n"
		"\n"
	);

	
	// Run 6c -DGOOS_goos -DGOARCH_goarch -Iworkdir -q
	// on each of the runtimedefs C files.
	vadd(&argv, bpathf(&b, "%s/bin/tool/%sc", goroot, gochar));
	vadd(&argv, bprintf(&b, "-DGOOS_%s", goos));
	vadd(&argv, bprintf(&b, "-DGOARCH_%s", goarch));
	vadd(&argv, bprintf(&b, "-I%s", workdir));
	vadd(&argv, "-q");
	vadd(&argv, "");
	p = argv.p[argv.len-1];
	for(i=0; i<nelem(runtimedefs); i++) {
		argv.p[argv.len-1] = runtimedefs[i];
		runv(&b, dir, CheckExit, &argv);
		bwriteb(&in, &b);
	}
	argv.p[argv.len-1] = p;
		
	// Process the aggregate output.
	skip = 0;
	splitlines(&lines, bstr(&in));
	for(i=0; i<lines.len; i++) {
		p = lines.p[i];
		// Drop comment, func, and const lines.
		if(hasprefix(p, "//") || hasprefix(p, "const") || hasprefix(p, "func"))
			continue;
		
		// Note beginning of type or var decl, which can be multiline.
		// Remove duplicates.  The linear check of seen here makes the
		// whole processing quadratic in aggregate, but there are only
		// about 100 declarations, so this is okay (and simple).
		if(hasprefix(p, "type ") || hasprefix(p, "var ")) {
			splitfields(&fields, p);
			if(fields.len < 2)
				continue;
			if(find(fields.p[1], seen.p, seen.len) >= 0) {
				if(streq(fields.p[fields.len-1], "{"))
					skip = 1;  // skip until }
				continue;
			}
			vadd(&seen, fields.p[1]);
		}
		if(skip) {
			if(hasprefix(p, "}"))
				skip = 0;
			continue;
		}
		
		bwritestr(&out, p);
	}
	
	writefile(&out, file);

	bfree(&in);
	bfree(&b);
	bfree(&out);
	vfree(&argv);
	vfree(&lines);
	vfree(&fields);
	vfree(&seen);
}
