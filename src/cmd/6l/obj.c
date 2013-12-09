// Inferno utils/6l/obj.c
// http://code.google.com/p/inferno-os/source/browse/utils/6l/obj.c
//
//	Copyright © 1994-1999 Lucent Technologies Inc.  All rights reserved.
//	Portions Copyright © 1995-1997 C H Forsyth (forsyth@terzarima.net)
//	Portions Copyright © 1997-1999 Vita Nuova Limited
//	Portions Copyright © 2000-2007 Vita Nuova Holdings Limited (www.vitanuova.com)
//	Portions Copyright © 2004,2006 Bruce Ellis
//	Portions Copyright © 2005-2007 C H Forsyth (forsyth@terzarima.net)
//	Revisions Copyright © 2000-2007 Lucent Technologies Inc. and others
//	Portions Copyright © 2009 The Go Authors.  All rights reserved.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.  IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

// Reading object files.

#include	"l.h"
#include	"../ld/lib.h"
#include	"../ld/elf.h"
#include	"../ld/macho.h"
#include	"../ld/dwarf.h"
#include	"../ld/pe.h"
#include	<ar.h>

char*	thestring 	= "amd64";
LinkArch*	thelinkarch = &linkamd64;

void
archinit(void)
{
	// getgoextlinkenabled is based on GO_EXTLINK_ENABLED when
	// Go was built; see ../../make.bash.
	if(linkmode == LinkAuto && strcmp(getgoextlinkenabled(), "0") == 0)
		linkmode = LinkInternal;

	if(flag_shared)
		linkmode = LinkExternal;

	switch(HEADTYPE) {
	default:
		if(linkmode == LinkAuto)
			linkmode = LinkInternal;
		if(linkmode == LinkExternal && strcmp(getgoextlinkenabled(), "1") != 0)
			sysfatal("cannot use -linkmode=external with -H %s", headstr(HEADTYPE));
		break;
	case Hdarwin:
	case Hdragonfly:
	case Hfreebsd:
	case Hlinux:
	case Hnetbsd:
	case Hopenbsd:
		break;
	}
	ctxt->linkmode = linkmode;

	switch(HEADTYPE) {
	default:
		diag("unknown -H option");
		errorexit();
	case Hplan9:		/* plan 9 */
		HEADR = 32L + 8L;
		if(INITTEXT == -1)
			INITTEXT = 0x200000+HEADR;
		if(INITDAT == -1)
			INITDAT = 0;
		if(INITRND == -1)
			INITRND = 0x200000;
		break;
	case Helf:		/* elf32 executable */
		HEADR = rnd(52L+3*32L, 16);
		if(INITTEXT == -1)
			INITTEXT = 0x80110000L;
		if(INITDAT == -1)
			INITDAT = 0;
		if(INITRND == -1)
			INITRND = 4096;
		break;
	case Hdarwin:		/* apple MACH */
		/*
		 * OS X system constant - offset from 0(GS) to our TLS.
		 * Explained in ../../pkg/runtime/cgo/gcc_darwin_amd64.c.
		 */
		ctxt->tlsoffset = 0x8a0;
		machoinit();
		HEADR = INITIAL_MACHO_HEADR;
		if(INITRND == -1)
			INITRND = 4096;
		if(INITTEXT == -1)
			INITTEXT = 4096+HEADR;
		if(INITDAT == -1)
			INITDAT = 0;
		break;
	case Hlinux:		/* elf64 executable */
	case Hfreebsd:		/* freebsd */
	case Hnetbsd:		/* netbsd */
	case Hopenbsd:		/* openbsd */
	case Hdragonfly:	/* dragonfly */
		/*
		 * ELF uses TLS offset negative from FS.
		 * Translate 0(FS) and 8(FS) into -16(FS) and -8(FS).
		 * Also known to ../../pkg/runtime/sys_linux_amd64.s
		 * and ../../pkg/runtime/cgo/gcc_linux_amd64.c.
		 */
		ctxt->tlsoffset = -16;
		elfinit();
		HEADR = ELFRESERVE;
		if(INITTEXT == -1)
			INITTEXT = (1<<22)+HEADR;
		if(INITDAT == -1)
			INITDAT = 0;
		if(INITRND == -1)
			INITRND = 4096;
		break;
	case Hwindows:		/* PE executable */
		peinit();
		HEADR = PEFILEHEADR;
		if(INITTEXT == -1)
			INITTEXT = PEBASE+PESECTHEADR;
		if(INITDAT == -1)
			INITDAT = 0;
		if(INITRND == -1)
			INITRND = PESECTALIGN;
		break;
	}
	if(INITDAT != 0 && INITRND != 0)
		print("warning: -D0x%llux is ignored because of -R0x%ux\n",
			INITDAT, INITRND);
}
