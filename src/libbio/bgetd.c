/*
http://code.google.com/p/inferno-os/source/browse/libbio/bgetd.c

	Copyright © 1994-1999 Lucent Technologies Inc.  All rights reserved.
	Revisions Copyright © 2000-2007 Vita Nuova Holdings Limited (www.vitanuova.com).  All rights reserved.

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.  IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

#include	<u.h>
#include	<libc.h>
#include <bio.h>

struct	bgetd
{
	Biobuf*	b;
	int		eof;
};

static int
Bgetdf(void *vp)
{
	int c;
	struct bgetd *bg = vp;

	c = Bgetc(bg->b);
	if(c == Beof)
		bg->eof = 1;
	return c;
}

int
Bgetd(Biobuf *bp, double *dp)
{
	double d;
	struct bgetd b;

	b.b = bp;
	b.eof = 0;
	d = fmtcharstod(Bgetdf, &b);
	if(b.eof)
		return -1;
	Bungetc(bp);
	*dp = d;
	return 1;
}
