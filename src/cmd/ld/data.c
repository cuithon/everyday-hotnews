// Inferno utils/8l/asm.c
// http://code.google.com/p/inferno-os/source/browse/utils/8l/asm.c
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

// Data layout and relocation.

#include	"l.h"
#include	"../ld/lib.h"
#include	"../ld/elf.h"
#include	"../ld/macho.h"
#include	"../ld/pe.h"
#include	"../../pkg/runtime/mgc0.h"

void	dynreloc(void);

/*
 * divide-and-conquer list-link
 * sort of LSym* structures.
 * Used for the data block.
 */
int
datcmp(LSym *s1, LSym *s2)
{
	if(s1->type != s2->type)
		return (int)s1->type - (int)s2->type;
	if(s1->size != s2->size) {
		if(s1->size < s2->size)
			return -1;
		return +1;
	}
	return strcmp(s1->name, s2->name);
}

LSym*
listsort(LSym *l, int (*cmp)(LSym*, LSym*), int off)
{
	LSym *l1, *l2, *le;
	#define NEXT(l) (*(LSym**)((char*)(l)+off))

	if(l == 0 || NEXT(l) == 0)
		return l;

	l1 = l;
	l2 = l;
	for(;;) {
		l2 = NEXT(l2);
		if(l2 == 0)
			break;
		l2 = NEXT(l2);
		if(l2 == 0)
			break;
		l1 = NEXT(l1);
	}

	l2 = NEXT(l1);
	NEXT(l1) = 0;
	l1 = listsort(l, cmp, off);
	l2 = listsort(l2, cmp, off);

	/* set up lead element */
	if(cmp(l1, l2) < 0) {
		l = l1;
		l1 = NEXT(l1);
	} else {
		l = l2;
		l2 = NEXT(l2);
	}
	le = l;

	for(;;) {
		if(l1 == 0) {
			while(l2) {
				NEXT(le) = l2;
				le = l2;
				l2 = NEXT(l2);
			}
			NEXT(le) = 0;
			break;
		}
		if(l2 == 0) {
			while(l1) {
				NEXT(le) = l1;
				le = l1;
				l1 = NEXT(l1);
			}
			break;
		}
		if(cmp(l1, l2) < 0) {
			NEXT(le) = l1;
			le = l1;
			l1 = NEXT(l1);
		} else {
			NEXT(le) = l2;
			le = l2;
			l2 = NEXT(l2);
		}
	}
	NEXT(le) = 0;
	return l;
	
	#undef NEXT
}

void
relocsym(LSym *s)
{
	Reloc *r;
	LSym *rs;
	Prog p;
	int32 i, off, siz, fl;
	vlong o;
	uchar *cast;

	ctxt->cursym = s;
	memset(&p, 0, sizeof p);
	for(r=s->r; r<s->r+s->nr; r++) {
		r->done = 1;
		off = r->off;
		siz = r->siz;
		if(off < 0 || off+siz > s->np) {
			diag("%s: invalid relocation %d+%d not in [%d,%d)", s->name, off, siz, 0, s->np);
			continue;
		}
		if(r->sym != S && (r->sym->type & SMASK == 0 || r->sym->type & SMASK == SXREF)) {
			diag("%s: not defined", r->sym->name);
			continue;
		}
		if(r->type >= 256)
			continue;

		// Solaris needs the ability to reference dynimport symbols.
		if(HEADTYPE != Hsolaris && r->sym != S && r->sym->type == SDYNIMPORT)
			diag("unhandled relocation for %s (type %d rtype %d)", r->sym->name, r->sym->type, r->type);
		if(r->sym != S && r->sym->type != STLSBSS && !r->sym->reachable)
			diag("unreachable sym in relocation: %s %s", s->name, r->sym->name);

		switch(r->type) {
		default:
			o = 0;
			if(archreloc(r, s, &o) < 0)
				diag("unknown reloc %d", r->type);
			break;
		case D_TLS:
			if(linkmode == LinkInternal && iself && thechar == '5') {
				// On ELF ARM, the thread pointer is 8 bytes before
				// the start of the thread-local data block, so add 8
				// to the actual TLS offset (r->sym->value).
				// This 8 seems to be a fundamental constant of
				// ELF on ARM (or maybe Glibc on ARM); it is not
				// related to the fact that our own TLS storage happens
				// to take up 8 bytes.
				o = 8 + r->sym->value;
				break;
			}
			r->done = 0;
			o = 0;
			if(thechar != '6')
				o = r->add;
			break;
		case D_ADDR:
			if(linkmode == LinkExternal && r->sym->type != SCONST) {
				r->done = 0;

				// set up addend for eventual relocation via outer symbol.
				rs = r->sym;
				r->xadd = r->add;
				while(rs->outer != nil) {
					r->xadd += symaddr(rs) - symaddr(rs->outer);
					rs = rs->outer;
				}
				if(rs->type != SHOSTOBJ && rs->type != SDYNIMPORT && rs->sect == nil)
					diag("missing section for %s", rs->name);
				r->xsym = rs;

				o = r->xadd;
				if(iself) {
					if(thechar == '6')
						o = 0;
				} else if(HEADTYPE == Hdarwin) {
					if(rs->type != SHOSTOBJ)
						o += symaddr(rs);
				} else {
					diag("unhandled pcrel relocation for %s", headstring);
				}
				break;
			}
			o = symaddr(r->sym) + r->add;
			break;
		case D_PCREL:
			// r->sym can be null when CALL $(constant) is transformed from absolute PC to relative PC call.
			if(linkmode == LinkExternal && r->sym && r->sym->type != SCONST && r->sym->sect != ctxt->cursym->sect) {
				r->done = 0;

				// set up addend for eventual relocation via outer symbol.
				rs = r->sym;
				r->xadd = r->add;
				while(rs->outer != nil) {
					r->xadd += symaddr(rs) - symaddr(rs->outer);
					rs = rs->outer;
				}
				r->xadd -= r->siz; // relative to address after the relocated chunk
				if(rs->type != SHOSTOBJ && rs->type != SDYNIMPORT && rs->sect == nil)
					diag("missing section for %s", rs->name);
				r->xsym = rs;

				o = r->xadd;
				if(iself) {
					if(thechar == '6')
						o = 0;
				} else if(HEADTYPE == Hdarwin) {
					if(rs->type != SHOSTOBJ)
						o += symaddr(rs) - rs->sect->vaddr;
					o -= r->off; // WTF?
				} else {
					diag("unhandled pcrel relocation for %s", headstring);
				}
				break;
			}
			o = 0;
			if(r->sym)
				o += symaddr(r->sym);
			// NOTE: The (int32) cast on the next line works around a bug in Plan 9's 8c
			// compiler. The expression s->value + r->off + r->siz is int32 + int32 +
			// uchar, and Plan 9 8c incorrectly treats the expression as type uint32
			// instead of int32, causing incorrect values when sign extended for adding
			// to o. The bug only occurs on Plan 9, because this C program is compiled by
			// the standard host compiler (gcc on most other systems).
			o += r->add - (s->value + r->off + (int32)r->siz);
			break;
		case D_SIZE:
			o = r->sym->size + r->add;
			break;
		}
//print("relocate %s %#llux (%#llux+%#llux, size %d) => %s %#llux +%#llx [%llx]\n", s->name, (uvlong)(s->value+off), (uvlong)s->value, (uvlong)r->off, r->siz, r->sym ? r->sym->name : "<nil>", (uvlong)symaddr(r->sym), (vlong)r->add, (vlong)o);
		switch(siz) {
		default:
			ctxt->cursym = s;
			diag("bad reloc size %#ux for %s", siz, r->sym->name);
		case 4:
			if(r->type == D_PCREL) {
				if(o != (int32)o)
					diag("pc-relative relocation address is too big: %#llx", o);
			} else {
				if(o != (int32)o && o != (uint32)o)
					diag("non-pc-relative relocation address is too big: %#llux", o);
			}
			fl = o;
			cast = (uchar*)&fl;
			for(i=0; i<4; i++)
				s->p[off+i] = cast[inuxi4[i]];
			break;
		case 8:
			cast = (uchar*)&o;
			for(i=0; i<8; i++)
				s->p[off+i] = cast[inuxi8[i]];
			break;
		}
	}
}

void
reloc(void)
{
	LSym *s;

	if(debug['v'])
		Bprint(&bso, "%5.2f reloc\n", cputime());
	Bflush(&bso);

	for(s=ctxt->textp; s!=S; s=s->next)
		relocsym(s);
	for(s=datap; s!=S; s=s->next)
		relocsym(s);
}

void
dynrelocsym(LSym *s)
{
	Reloc *r;

	if(HEADTYPE == Hwindows) {
		LSym *rel, *targ;

		rel = linklookup(ctxt, ".rel", 0);
		if(s == rel)
			return;
		for(r=s->r; r<s->r+s->nr; r++) {
			targ = r->sym;
			if(!targ->reachable)
				diag("internal inconsistency: dynamic symbol %s is not reachable.", targ->name);
			if(r->sym->plt == -2 && r->sym->got != -2) { // make dynimport JMP table for PE object files.
				targ->plt = rel->size;
				r->sym = rel;
				r->add = targ->plt;

				// jmp *addr
				if(thechar == '8') {
					adduint8(ctxt, rel, 0xff);
					adduint8(ctxt, rel, 0x25);
					addaddr(ctxt, rel, targ);
					adduint8(ctxt, rel, 0x90);
					adduint8(ctxt, rel, 0x90);
				} else {
					adduint8(ctxt, rel, 0xff);
					adduint8(ctxt, rel, 0x24);
					adduint8(ctxt, rel, 0x25);
					addaddrplus4(ctxt, rel, targ, 0);
					adduint8(ctxt, rel, 0x90);
				}
			} else if(r->sym->plt >= 0) {
				r->sym = rel;
				r->add = targ->plt;
			}
		}
		return;
	}

	for(r=s->r; r<s->r+s->nr; r++) {
		if(r->sym != S && r->sym->type == SDYNIMPORT || r->type >= 256) {
			if(r->sym != S && !r->sym->reachable)
				diag("internal inconsistency: dynamic symbol %s is not reachable.", r->sym->name);
			adddynrel(s, r);
		}
	}
}

void
dynreloc(void)
{
	LSym *s;

	// -d suppresses dynamic loader format, so we may as well not
	// compute these sections or mark their symbols as reachable.
	if(debug['d'] && HEADTYPE != Hwindows)
		return;
	if(debug['v'])
		Bprint(&bso, "%5.2f reloc\n", cputime());
	Bflush(&bso);

	for(s=ctxt->textp; s!=S; s=s->next)
		dynrelocsym(s);
	for(s=datap; s!=S; s=s->next)
		dynrelocsym(s);
	if(iself)
		elfdynhash();
}

static void
blk(LSym *start, int32 addr, int32 size)
{
	LSym *sym;
	int32 eaddr;
	uchar *p, *ep;

	for(sym = start; sym != nil; sym = sym->next)
		if(!(sym->type&SSUB) && sym->value >= addr)
			break;

	eaddr = addr+size;
	for(; sym != nil; sym = sym->next) {
		if(sym->type&SSUB)
			continue;
		if(sym->value >= eaddr)
			break;
		if(sym->value < addr) {
			diag("phase error: addr=%#llx but sym=%#llx type=%d", (vlong)addr, (vlong)sym->value, sym->type);
			errorexit();
		}
		ctxt->cursym = sym;
		for(; addr < sym->value; addr++)
			cput(0);
		p = sym->p;
		ep = p + sym->np;
		while(p < ep)
			cput(*p++);
		addr += sym->np;
		for(; addr < sym->value+sym->size; addr++)
			cput(0);
		if(addr != sym->value+sym->size) {
			diag("phase error: addr=%#llx value+size=%#llx", (vlong)addr, (vlong)sym->value+sym->size);
			errorexit();
		}
	}

	for(; addr < eaddr; addr++)
		cput(0);
	cflush();
}

void
codeblk(int32 addr, int32 size)
{
	LSym *sym;
	int32 eaddr, n;
	uchar *q;

	if(debug['a'])
		Bprint(&bso, "codeblk [%#x,%#x) at offset %#llx\n", addr, addr+size, cpos());

	blk(ctxt->textp, addr, size);

	/* again for printing */
	if(!debug['a'])
		return;

	for(sym = ctxt->textp; sym != nil; sym = sym->next) {
		if(!sym->reachable)
			continue;
		if(sym->value >= addr)
			break;
	}

	eaddr = addr + size;
	for(; sym != nil; sym = sym->next) {
		if(!sym->reachable)
			continue;
		if(sym->value >= eaddr)
			break;

		if(addr < sym->value) {
			Bprint(&bso, "%-20s %.8llux|", "_", (vlong)addr);
			for(; addr < sym->value; addr++)
				Bprint(&bso, " %.2ux", 0);
			Bprint(&bso, "\n");
		}

		Bprint(&bso, "%.6llux\t%-20s\n", (vlong)addr, sym->name);
		n = sym->size;
		q = sym->p;

		while(n >= 16) {
			Bprint(&bso, "%.6ux\t%-20.16I\n", addr, q);
			addr += 16;
			q += 16;
			n -= 16;
		}
		if(n > 0)
			Bprint(&bso, "%.6ux\t%-20.*I\n", addr, (int)n, q);
		addr += n;
	}

	if(addr < eaddr) {
		Bprint(&bso, "%-20s %.8llux|", "_", (vlong)addr);
		for(; addr < eaddr; addr++)
			Bprint(&bso, " %.2ux", 0);
	}
	Bflush(&bso);
}

void
datblk(int32 addr, int32 size)
{
	LSym *sym;
	int32 i, eaddr;
	uchar *p, *ep;
	char *typ, *rsname;
	Reloc *r;

	if(debug['a'])
		Bprint(&bso, "datblk [%#x,%#x) at offset %#llx\n", addr, addr+size, cpos());

	blk(datap, addr, size);

	/* again for printing */
	if(!debug['a'])
		return;

	for(sym = datap; sym != nil; sym = sym->next)
		if(sym->value >= addr)
			break;

	eaddr = addr + size;
	for(; sym != nil; sym = sym->next) {
		if(sym->value >= eaddr)
			break;
		if(addr < sym->value) {
			Bprint(&bso, "\t%.8ux| 00 ...\n", addr);
			addr = sym->value;
		}
		Bprint(&bso, "%s\n\t%.8ux|", sym->name, (uint)addr);
		p = sym->p;
		ep = p + sym->np;
		while(p < ep) {
			if(p > sym->p && (int)(p-sym->p)%16 == 0)
				Bprint(&bso, "\n\t%.8ux|", (uint)(addr+(p-sym->p)));
			Bprint(&bso, " %.2ux", *p++);
		}
		addr += sym->np;
		for(; addr < sym->value+sym->size; addr++)
			Bprint(&bso, " %.2ux", 0);
		Bprint(&bso, "\n");
		
		if(linkmode == LinkExternal) {
			for(i=0; i<sym->nr; i++) {
				r = &sym->r[i];
				rsname = "";
				if(r->sym)
					rsname = r->sym->name;
				typ = "?";
				switch(r->type) {
				case D_ADDR:
					typ = "addr";
					break;
				case D_PCREL:
					typ = "pcrel";
					break;
				}
				Bprint(&bso, "\treloc %.8ux/%d %s %s+%#llx [%#llx]\n",
					(uint)(sym->value+r->off), r->siz, typ, rsname, (vlong)r->add, (vlong)(r->sym->value+r->add));
			}
		}				
	}

	if(addr < eaddr)
		Bprint(&bso, "\t%.8ux| 00 ...\n", (uint)addr);
	Bprint(&bso, "\t%.8ux|\n", (uint)eaddr);
}

void
strnput(char *s, int n)
{
	for(; n > 0 && *s; s++) {
		cput(*s);
		n--;
	}
	while(n > 0) {
		cput(0);
		n--;
	}
}

void
addstrdata(char *name, char *value)
{
	LSym *s, *sp;
	char *p;

	p = smprint("%s.str", name);
	sp = linklookup(ctxt, p, 0);
	free(p);
	addstring(sp, value);

	s = linklookup(ctxt, name, 0);
	s->size = 0;
	s->dupok = 1;
	addaddr(ctxt, s, sp);
	adduint32(ctxt, s, strlen(value));
	if(PtrSize == 8)
		adduint32(ctxt, s, 0);  // round struct to pointer width

	// in case reachability has already been computed
	sp->reachable = s->reachable;
}

vlong
addstring(LSym *s, char *str)
{
	int n;
	int32 r;

	if(s->type == 0)
		s->type = SNOPTRDATA;
	s->reachable = 1;
	r = s->size;
	n = strlen(str)+1;
	if(strcmp(s->name, ".shstrtab") == 0)
		elfsetstring(str, r);
	symgrow(ctxt, s, r+n);
	memmove(s->p+r, str, n);
	s->size += n;
	return r;
}

void
dosymtype(void)
{
	LSym *s;

	for(s = ctxt->allsym; s != nil; s = s->allsym) {
		if(s->np > 0) {
			if(s->type == SBSS)
				s->type = SDATA;
			if(s->type == SNOPTRBSS)
				s->type = SNOPTRDATA;
		}
	}
}

static int32
symalign(LSym *s)
{
	int32 align;

	if(s->align != 0)
		return s->align;

	align = MaxAlign;
	while(align > s->size && align > 1)
		align >>= 1;
	if(align < s->align)
		align = s->align;
	return align;
}
	
static vlong
aligndatsize(vlong datsize, LSym *s)
{
	return rnd(datsize, symalign(s));
}

// maxalign returns the maximum required alignment for
// the list of symbols s; the list stops when s->type exceeds type.
static int32
maxalign(LSym *s, int type)
{
	int32 align, max;
	
	max = 0;
	for(; s != S && s->type <= type; s = s->next) {
		align = symalign(s);
		if(max < align)
			max = align;
	}
	return max;
}

static void
gcaddsym(LSym *gc, LSym *s, vlong off)
{
	vlong a;
	LSym *gotype;

	if(s->size < PtrSize)
		return;
	if(strcmp(s->name, ".string") == 0)
		return;

	gotype = s->gotype;
	if(gotype != nil) {
		//print("gcaddsym:    %s    %d    %s\n", s->name, s->size, gotype->name);
		adduintxx(ctxt, gc, GC_CALL, PtrSize);
		adduintxx(ctxt, gc, off, PtrSize);
		addpcrelplus(ctxt, gc, decodetype_gc(gotype), 3*PtrSize+4);
		if(PtrSize == 8)
			adduintxx(ctxt, gc, 0, 4);
	} else {
		//print("gcaddsym:    %s    %d    <unknown type>\n", s->name, s->size);
		for(a = -off&(PtrSize-1); a+PtrSize<=s->size; a+=PtrSize) {
			adduintxx(ctxt, gc, GC_APTR, PtrSize);
			adduintxx(ctxt, gc, off+a, PtrSize);
		}
	}
}

void
growdatsize(vlong *datsizep, LSym *s)
{
	vlong datsize;
	
	datsize = *datsizep;
	if(s->size < 0)
		diag("negative size (datsize = %lld, s->size = %lld)", datsize, s->size);
	if(datsize + s->size < datsize)
		diag("symbol too large (datsize = %lld, s->size = %lld)", datsize, s->size);
	*datsizep = datsize + s->size;
}

void
dodata(void)
{
	int32 n;
	vlong datsize;
	Section *sect;
	Segment *segro;
	LSym *s, *last, **l;
	LSym *gcdata1, *gcbss1;

	if(debug['v'])
		Bprint(&bso, "%5.2f dodata\n", cputime());
	Bflush(&bso);

	gcdata1 = linklookup(ctxt, "gcdata", 0);
	gcbss1 = linklookup(ctxt, "gcbss", 0);

	// size of .data and .bss section. the zero value is later replaced by the actual size of the section.
	adduintxx(ctxt, gcdata1, 0, PtrSize);
	adduintxx(ctxt, gcbss1, 0, PtrSize);

	last = nil;
	datap = nil;

	for(s=ctxt->allsym; s!=S; s=s->allsym) {
		if(!s->reachable || s->special)
			continue;
		if(STEXT < s->type && s->type < SXREF) {
			if(last == nil)
				datap = s;
			else
				last->next = s;
			s->next = nil;
			last = s;
		}
	}

	for(s = datap; s != nil; s = s->next) {
		if(s->np > s->size)
			diag("%s: initialize bounds (%lld < %d)",
				s->name, (vlong)s->size, s->np);
	}


	/*
	 * now that we have the datap list, but before we start
	 * to assign addresses, record all the necessary
	 * dynamic relocations.  these will grow the relocation
	 * symbol, which is itself data.
	 *
	 * on darwin, we need the symbol table numbers for dynreloc.
	 */
	if(HEADTYPE == Hdarwin)
		machosymorder();
	dynreloc();

	/* some symbols may no longer belong in datap (Mach-O) */
	for(l=&datap; (s=*l) != nil; ) {
		if(s->type <= STEXT || SXREF <= s->type)
			*l = s->next;
		else
			l = &s->next;
	}
	*l = nil;

	datap = listsort(datap, datcmp, offsetof(LSym, next));

	/*
	 * allocate sections.  list is sorted by type,
	 * so we can just walk it for each piece we want to emit.
	 * segdata is processed before segtext, because we need
	 * to see all symbols in the .data and .bss sections in order
	 * to generate garbage collection information.
	 */

	/* begin segdata */

	/* skip symbols belonging to segtext */
	s = datap;
	for(; s != nil && s->type < SELFSECT; s = s->next)
		;

	/* writable ELF sections */
	datsize = 0;
	for(; s != nil && s->type < SNOPTRDATA; s = s->next) {
		sect = addsection(&segdata, s->name, 06);
		sect->align = symalign(s);
		datsize = rnd(datsize, sect->align);
		sect->vaddr = datsize;
		s->sect = sect;
		s->type = SDATA;
		s->value = datsize - sect->vaddr;
		growdatsize(&datsize, s);
		sect->len = datsize - sect->vaddr;
	}

	/* pointer-free data */
	sect = addsection(&segdata, ".noptrdata", 06);
	sect->align = maxalign(s, SINITARR-1);
	datsize = rnd(datsize, sect->align);
	sect->vaddr = datsize;
	linklookup(ctxt, "noptrdata", 0)->sect = sect;
	linklookup(ctxt, "enoptrdata", 0)->sect = sect;
	for(; s != nil && s->type < SINITARR; s = s->next) {
		datsize = aligndatsize(datsize, s);
		s->sect = sect;
		s->type = SDATA;
		s->value = datsize - sect->vaddr;
		growdatsize(&datsize, s);
	}
	sect->len = datsize - sect->vaddr;

	/* shared library initializer */
	if(flag_shared) {
		sect = addsection(&segdata, ".init_array", 06);
		sect->align = maxalign(s, SINITARR);
		datsize = rnd(datsize, sect->align);
		sect->vaddr = datsize;
		for(; s != nil && s->type == SINITARR; s = s->next) {
			datsize = aligndatsize(datsize, s);
			s->sect = sect;
			s->value = datsize - sect->vaddr;
			growdatsize(&datsize, s);
		}
		sect->len = datsize - sect->vaddr;
	}

	/* data */
	sect = addsection(&segdata, ".data", 06);
	sect->align = maxalign(s, SBSS-1);
	datsize = rnd(datsize, sect->align);
	sect->vaddr = datsize;
	linklookup(ctxt, "data", 0)->sect = sect;
	linklookup(ctxt, "edata", 0)->sect = sect;
	for(; s != nil && s->type < SBSS; s = s->next) {
		if(s->type == SINITARR) {
			ctxt->cursym = s;
			diag("unexpected symbol type %d", s->type);
		}
		s->sect = sect;
		s->type = SDATA;
		datsize = aligndatsize(datsize, s);
		s->value = datsize - sect->vaddr;
		gcaddsym(gcdata1, s, datsize - sect->vaddr);  // gc
		growdatsize(&datsize, s);
	}
	sect->len = datsize - sect->vaddr;

	adduintxx(ctxt, gcdata1, GC_END, PtrSize);
	setuintxx(ctxt, gcdata1, 0, sect->len, PtrSize);

	/* bss */
	sect = addsection(&segdata, ".bss", 06);
	sect->align = maxalign(s, SNOPTRBSS-1);
	datsize = rnd(datsize, sect->align);
	sect->vaddr = datsize;
	linklookup(ctxt, "bss", 0)->sect = sect;
	linklookup(ctxt, "ebss", 0)->sect = sect;
	for(; s != nil && s->type < SNOPTRBSS; s = s->next) {
		s->sect = sect;
		datsize = aligndatsize(datsize, s);
		s->value = datsize - sect->vaddr;
		gcaddsym(gcbss1, s, datsize - sect->vaddr);  // gc
		growdatsize(&datsize, s);
	}
	sect->len = datsize - sect->vaddr;

	adduintxx(ctxt, gcbss1, GC_END, PtrSize);
	setuintxx(ctxt, gcbss1, 0, sect->len, PtrSize);

	/* pointer-free bss */
	sect = addsection(&segdata, ".noptrbss", 06);
	sect->align = maxalign(s, SNOPTRBSS);
	datsize = rnd(datsize, sect->align);
	sect->vaddr = datsize;
	linklookup(ctxt, "noptrbss", 0)->sect = sect;
	linklookup(ctxt, "enoptrbss", 0)->sect = sect;
	for(; s != nil && s->type == SNOPTRBSS; s = s->next) {
		datsize = aligndatsize(datsize, s);
		s->sect = sect;
		s->value = datsize - sect->vaddr;
		growdatsize(&datsize, s);
	}
	sect->len = datsize - sect->vaddr;
	linklookup(ctxt, "end", 0)->sect = sect;

	// 6g uses 4-byte relocation offsets, so the entire segment must fit in 32 bits.
	if(datsize != (uint32)datsize) {
		diag("data or bss segment too large");
	}
	
	if(iself && linkmode == LinkExternal && s != nil && s->type == STLSBSS && HEADTYPE != Hopenbsd) {
		sect = addsection(&segdata, ".tbss", 06);
		sect->align = PtrSize;
		sect->vaddr = 0;
		datsize = 0;
		for(; s != nil && s->type == STLSBSS; s = s->next) {
			datsize = aligndatsize(datsize, s);
			s->sect = sect;
			s->value = datsize - sect->vaddr;
			growdatsize(&datsize, s);
		}
		sect->len = datsize;
	} else {
		// Might be internal linking but still using cgo.
		// In that case, the only possible STLSBSS symbol is tlsgm.
		// Give it offset 0, because it's the only thing here.
		if(s != nil && s->type == STLSBSS && strcmp(s->name, "runtime.tlsgm") == 0) {
			s->value = 0;
			s = s->next;
		}
	}
	
	if(s != nil) {
		ctxt->cursym = nil;
		diag("unexpected symbol type %d for %s", s->type, s->name);
	}

	/*
	 * We finished data, begin read-only data.
	 * Not all systems support a separate read-only non-executable data section.
	 * ELF systems do.
	 * OS X and Plan 9 do not.
	 * Windows PE may, but if so we have not implemented it.
	 * And if we're using external linking mode, the point is moot,
	 * since it's not our decision; that code expects the sections in
	 * segtext.
	 */
	if(iself && linkmode == LinkInternal)
		segro = &segrodata;
	else
		segro = &segtext;

	s = datap;
	
	datsize = 0;
	
	/* read-only executable ELF, Mach-O sections */
	for(; s != nil && s->type < STYPE; s = s->next) {
		sect = addsection(&segtext, s->name, 04);
		sect->align = symalign(s);
		datsize = rnd(datsize, sect->align);
		sect->vaddr = datsize;
		s->sect = sect;
		s->type = SRODATA;
		s->value = datsize - sect->vaddr;
		growdatsize(&datsize, s);
		sect->len = datsize - sect->vaddr;
	}

	/* read-only data */
	sect = addsection(segro, ".rodata", 04);
	sect->align = maxalign(s, STYPELINK-1);
	datsize = rnd(datsize, sect->align);
	sect->vaddr = 0;
	linklookup(ctxt, "rodata", 0)->sect = sect;
	linklookup(ctxt, "erodata", 0)->sect = sect;
	for(; s != nil && s->type < STYPELINK; s = s->next) {
		datsize = aligndatsize(datsize, s);
		s->sect = sect;
		s->type = SRODATA;
		s->value = datsize - sect->vaddr;
		growdatsize(&datsize, s);
	}
	sect->len = datsize - sect->vaddr;

	/* typelink */
	sect = addsection(segro, ".typelink", 04);
	sect->align = maxalign(s, STYPELINK);
	datsize = rnd(datsize, sect->align);
	sect->vaddr = datsize;
	linklookup(ctxt, "typelink", 0)->sect = sect;
	linklookup(ctxt, "etypelink", 0)->sect = sect;
	for(; s != nil && s->type == STYPELINK; s = s->next) {
		datsize = aligndatsize(datsize, s);
		s->sect = sect;
		s->type = SRODATA;
		s->value = datsize - sect->vaddr;
		growdatsize(&datsize, s);
	}
	sect->len = datsize - sect->vaddr;

	/* gosymtab */
	sect = addsection(segro, ".gosymtab", 04);
	sect->align = maxalign(s, SPCLNTAB-1);
	datsize = rnd(datsize, sect->align);
	sect->vaddr = datsize;
	linklookup(ctxt, "symtab", 0)->sect = sect;
	linklookup(ctxt, "esymtab", 0)->sect = sect;
	for(; s != nil && s->type < SPCLNTAB; s = s->next) {
		datsize = aligndatsize(datsize, s);
		s->sect = sect;
		s->type = SRODATA;
		s->value = datsize - sect->vaddr;
		growdatsize(&datsize, s);
	}
	sect->len = datsize - sect->vaddr;

	/* gopclntab */
	sect = addsection(segro, ".gopclntab", 04);
	sect->align = maxalign(s, SELFROSECT-1);
	datsize = rnd(datsize, sect->align);
	sect->vaddr = datsize;
	linklookup(ctxt, "pclntab", 0)->sect = sect;
	linklookup(ctxt, "epclntab", 0)->sect = sect;
	for(; s != nil && s->type < SELFROSECT; s = s->next) {
		datsize = aligndatsize(datsize, s);
		s->sect = sect;
		s->type = SRODATA;
		s->value = datsize - sect->vaddr;
		growdatsize(&datsize, s);
	}
	sect->len = datsize - sect->vaddr;

	/* read-only ELF, Mach-O sections */
	for(; s != nil && s->type < SELFSECT; s = s->next) {
		sect = addsection(segro, s->name, 04);
		sect->align = symalign(s);
		datsize = rnd(datsize, sect->align);
		sect->vaddr = datsize;
		s->sect = sect;
		s->type = SRODATA;
		s->value = datsize - sect->vaddr;
		growdatsize(&datsize, s);
		sect->len = datsize - sect->vaddr;
	}

	// 6g uses 4-byte relocation offsets, so the entire segment must fit in 32 bits.
	if(datsize != (uint32)datsize) {
		diag("read-only data segment too large");
	}
	
	/* number the sections */
	n = 1;
	for(sect = segtext.sect; sect != nil; sect = sect->next)
		sect->extnum = n++;
	for(sect = segrodata.sect; sect != nil; sect = sect->next)
		sect->extnum = n++;
	for(sect = segdata.sect; sect != nil; sect = sect->next)
		sect->extnum = n++;
}

// assign addresses to text
void
textaddress(void)
{
	uvlong va;
	Section *sect;
	LSym *sym, *sub;

	addsection(&segtext, ".text", 05);

	// Assign PCs in text segment.
	// Could parallelize, by assigning to text
	// and then letting threads copy down, but probably not worth it.
	sect = segtext.sect;
	sect->align = FuncAlign;
	linklookup(ctxt, "text", 0)->sect = sect;
	linklookup(ctxt, "etext", 0)->sect = sect;
	va = INITTEXT;
	sect->vaddr = va;
	for(sym = ctxt->textp; sym != nil; sym = sym->next) {
		sym->sect = sect;
		if(sym->type & SSUB)
			continue;
		if(sym->align != 0)
			va = rnd(va, sym->align);
		sym->value = 0;
		for(sub = sym; sub != S; sub = sub->sub)
			sub->value += va;
		if(sym->size == 0 && sym->sub != S)
			ctxt->cursym = sym;
		va += sym->size;
	}
	sect->len = va - sect->vaddr;
}

// assign addresses
void
address(void)
{
	Section *s, *text, *data, *rodata, *symtab, *pclntab, *noptr, *bss, *noptrbss;
	Section *typelink;
	LSym *sym, *sub;
	uvlong va;
	vlong vlen;

	va = INITTEXT;
	segtext.rwx = 05;
	segtext.vaddr = va;
	segtext.fileoff = HEADR;
	for(s=segtext.sect; s != nil; s=s->next) {
//print("%s at %#llux + %#llux\n", s->name, va, (vlong)s->len);
		va = rnd(va, s->align);
		s->vaddr = va;
		va += s->len;
	}
	segtext.len = va - INITTEXT;
	segtext.filelen = segtext.len;

	if(segrodata.sect != nil) {
		// align to page boundary so as not to mix
		// rodata and executable text.
		va = rnd(va, INITRND);

		segrodata.rwx = 04;
		segrodata.vaddr = va;
		segrodata.fileoff = va - segtext.vaddr + segtext.fileoff;
		segrodata.filelen = 0;
		for(s=segrodata.sect; s != nil; s=s->next) {
			va = rnd(va, s->align);
			s->vaddr = va;
			va += s->len;
		}
		segrodata.len = va - segrodata.vaddr;
		segrodata.filelen = segrodata.len;
	}

	va = rnd(va, INITRND);
	segdata.rwx = 06;
	segdata.vaddr = va;
	segdata.fileoff = va - segtext.vaddr + segtext.fileoff;
	segdata.filelen = 0;
	if(HEADTYPE == Hwindows)
		segdata.fileoff = segtext.fileoff + rnd(segtext.len, PEFILEALIGN);
	if(HEADTYPE == Hplan9)
		segdata.fileoff = segtext.fileoff + segtext.filelen;
	data = nil;
	noptr = nil;
	bss = nil;
	noptrbss = nil;
	for(s=segdata.sect; s != nil; s=s->next) {
		vlen = s->len;
		if(s->next)
			vlen = s->next->vaddr - s->vaddr;
		s->vaddr = va;
		va += vlen;
		segdata.len = va - segdata.vaddr;
		if(strcmp(s->name, ".data") == 0)
			data = s;
		if(strcmp(s->name, ".noptrdata") == 0)
			noptr = s;
		if(strcmp(s->name, ".bss") == 0)
			bss = s;
		if(strcmp(s->name, ".noptrbss") == 0)
			noptrbss = s;
	}
	segdata.filelen = bss->vaddr - segdata.vaddr;

	text = segtext.sect;
	if(segrodata.sect)
		rodata = segrodata.sect;
	else
		rodata = text->next;
	typelink = rodata->next;
	symtab = typelink->next;
	pclntab = symtab->next;

	for(sym = datap; sym != nil; sym = sym->next) {
		ctxt->cursym = sym;
		if(sym->sect != nil)
			sym->value += sym->sect->vaddr;
		for(sub = sym->sub; sub != nil; sub = sub->sub)
			sub->value += sym->value;
	}

	xdefine("text", STEXT, text->vaddr);
	xdefine("etext", STEXT, text->vaddr + text->len);
	xdefine("rodata", SRODATA, rodata->vaddr);
	xdefine("erodata", SRODATA, rodata->vaddr + rodata->len);
	xdefine("typelink", SRODATA, typelink->vaddr);
	xdefine("etypelink", SRODATA, typelink->vaddr + typelink->len);

	sym = linklookup(ctxt, "gcdata", 0);
	xdefine("egcdata", SRODATA, symaddr(sym) + sym->size);
	linklookup(ctxt, "egcdata", 0)->sect = sym->sect;

	sym = linklookup(ctxt, "gcbss", 0);
	xdefine("egcbss", SRODATA, symaddr(sym) + sym->size);
	linklookup(ctxt, "egcbss", 0)->sect = sym->sect;

	xdefine("symtab", SRODATA, symtab->vaddr);
	xdefine("esymtab", SRODATA, symtab->vaddr + symtab->len);
	xdefine("pclntab", SRODATA, pclntab->vaddr);
	xdefine("epclntab", SRODATA, pclntab->vaddr + pclntab->len);
	xdefine("noptrdata", SNOPTRDATA, noptr->vaddr);
	xdefine("enoptrdata", SNOPTRDATA, noptr->vaddr + noptr->len);
	xdefine("bss", SBSS, bss->vaddr);
	xdefine("ebss", SBSS, bss->vaddr + bss->len);
	xdefine("data", SDATA, data->vaddr);
	xdefine("edata", SDATA, data->vaddr + data->len);
	xdefine("noptrbss", SNOPTRBSS, noptrbss->vaddr);
	xdefine("enoptrbss", SNOPTRBSS, noptrbss->vaddr + noptrbss->len);
	xdefine("end", SBSS, segdata.vaddr + segdata.len);
}
