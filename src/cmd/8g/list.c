// Derived from Inferno utils/8c/list.c
// http://code.google.com/p/inferno-os/source/browse/utils/8c/list.c
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

#include <u.h>
#include <libc.h>
#include "gg.h"

static	int	sconsize;
void
listinit(void)
{

	fmtinstall('A', Aconv);		// as
	fmtinstall('P', Pconv);		// Prog*
	fmtinstall('D', Dconv);		// Addr*
	fmtinstall('R', Rconv);		// reg
	fmtinstall('Y', Yconv);		// sconst
}

int
Pconv(Fmt *fp)
{
	char str[STRINGSZ];
	Prog *p;
	char scale[40];

	p = va_arg(fp->args, Prog*);
	sconsize = 8;
	scale[0] = '\0';
	if(p->from.scale != 0 && (p->as == AGLOBL || p->as == ATEXT))
		snprint(scale, sizeof scale, "%d,", p->from.scale);
	switch(p->as) {
	default:
		snprint(str, sizeof(str), "%.4d (%L) %-7A %D,%s%D",
			p->loc, p->lineno, p->as, &p->from, scale, &p->to);
		break;

	case ADATA:
		sconsize = p->from.scale;
		snprint(str, sizeof(str), "%.4d (%L) %-7A %D/%d,%D",
			p->loc, p->lineno, p->as, &p->from, sconsize, &p->to);
		break;

	case ATEXT:
		snprint(str, sizeof(str), "%.4d (%L) %-7A %D,%s%lD",
			p->loc, p->lineno, p->as, &p->from, scale, &p->to);
		break;
	}
	return fmtstrcpy(fp, str);
}

int
Dconv(Fmt *fp)
{
	char str[STRINGSZ], s[STRINGSZ];
	Addr *a;
	int i;
	uint32 d1, d2;

	a = va_arg(fp->args, Addr*);
	i = a->type;
	if(i >= D_INDIR) {
		if(a->offset)
			snprint(str, sizeof(str), "%d(%R)", a->offset, i-D_INDIR);
		else
			snprint(str, sizeof(str), "(%R)", i-D_INDIR);
		goto brk;
	}
	switch(i) {

	default:
		if(a->offset)
			snprint(str, sizeof(str), "$%d,%R", a->offset, i);
		else
			snprint(str, sizeof(str), "%R", i);
		break;

	case D_NONE:
		str[0] = 0;
		break;

	case D_BRANCH:
		snprint(str, sizeof(str), "%d", a->u.branch->loc);
		break;

	case D_EXTERN:
		snprint(str, sizeof(str), "%S+%d(SB)", a->sym, a->offset);
		break;

	case D_STATIC:
		snprint(str, sizeof(str), "%S<>+%d(SB)", a->sym, a->offset);
		break;

	case D_AUTO:
		snprint(str, sizeof(str), "%S+%d(SP)", a->sym, a->offset);
		break;

	case D_PARAM:
		snprint(str, sizeof(str), "%S+%d(FP)", a->sym, a->offset);
		break;

	case D_CONST:
		if(fp->flags & FmtLong) {
			d1 = a->offset;
			d2 = a->offset2;
			snprint(str, sizeof(str), "$%lud-%lud", (ulong)d1, (ulong)d2);
			break;
		}
		snprint(str, sizeof(str), "$%d", a->offset);
		break;

	case D_FCONST:
		snprint(str, sizeof(str), "$(%.17e)", a->u.dval);
		break;

	case D_SCONST:
		snprint(str, sizeof(str), "$\"%Y\"", a->u.sval);
		break;

	case D_ADDR:
		a->type = a->index;
		a->index = D_NONE;
		snprint(str, sizeof(str), "$%D", a);
		a->index = a->type;
		a->type = D_ADDR;
		goto conv;
	}
brk:
	if(a->index != D_NONE) {
		snprint(s, sizeof(s), "(%R*%d)", (int)a->index, (int)a->scale);
		strcat(str, s);
	}
conv:
	return fmtstrcpy(fp, str);
}

static	char*	regstr[] =
{
	"AL",		/* [D_AL] */
	"CL",
	"DL",
	"BL",

	"AH",	/* [D_AH] */
	"CH",
	"DH",
	"BH",

	"AX",		/* [D_AX] */
	"CX",
	"DX",
	"BX",
	"SP",
	"BP",
	"SI",
	"DI",

	"F0",		/* [D_F0] */
	"F1",
	"F2",
	"F3",
	"F4",
	"F5",
	"F6",
	"F7",

	"CS",		/* [D_CS] */
	"SS",
	"DS",
	"ES",
	"FS",
	"GS",

	"GDTR",		/* [D_GDTR] */
	"IDTR",		/* [D_IDTR] */
	"LDTR",		/* [D_LDTR] */
	"MSW",		/* [D_MSW] */
	"TASK",		/* [D_TASK] */

	"CR0",		/* [D_CR] */
	"CR1",
	"CR2",
	"CR3",
	"CR4",
	"CR5",
	"CR6",
	"CR7",

	"DR0",		/* [D_DR] */
	"DR1",
	"DR2",
	"DR3",
	"DR4",
	"DR5",
	"DR6",
	"DR7",

	"TR0",		/* [D_TR] */
	"TR1",
	"TR2",
	"TR3",
	"TR4",
	"TR5",
	"TR6",
	"TR7",

	"NONE",		/* [D_NONE] */
};

int
Rconv(Fmt *fp)
{
	char str[STRINGSZ];
	int r;

	r = va_arg(fp->args, int);
	if(r < 0 || r >= nelem(regstr) || regstr[r] == nil) {
		snprint(str, sizeof(str), "BAD_R(%d)", r);
		return fmtstrcpy(fp, str);
	}
	return fmtstrcpy(fp, regstr[r]);
}

int
Aconv(Fmt *fp)
{
	int i;

	i = va_arg(fp->args, int);
	return fmtstrcpy(fp, anames[i]);
}


int
Yconv(Fmt *fp)
{
	int i, c;
	char str[STRINGSZ], *p, *a;

	a = va_arg(fp->args, char*);
	p = str;
	for(i=0; i<sconsize; i++) {
		c = a[i] & 0xff;
		if((c >= 'a' && c <= 'z') ||
		   (c >= 'A' && c <= 'Z') ||
		   (c >= '0' && c <= '9')) {
			*p++ = c;
			continue;
		}
		*p++ = '\\';
		switch(c) {
		default:
			if(c < 040 || c >= 0177)
				break;	/* not portable */
			p[-1] = c;
			continue;
		case 0:
			*p++ = 'z';
			continue;
		case '\\':
		case '"':
			*p++ = c;
			continue;
		case '\n':
			*p++ = 'n';
			continue;
		case '\t':
			*p++ = 't';
			continue;
		}
		*p++ = (c>>6) + '0';
		*p++ = ((c>>3) & 7) + '0';
		*p++ = (c & 7) + '0';
	}
	*p = 0;
	return fmtstrcpy(fp, str);
}
