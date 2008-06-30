// Copyright 2009 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "runtime.h"

typedef	struct	Link	Link;
typedef	struct	Hmap	Hmap;
typedef	struct	Alg	Alg;

struct	Alg
{
	uint64	(*hash)(uint32, void*);
	uint32	(*equal)(uint32, void*, void*);
	void	(*print)(uint32, void*);
	void	(*copy)(uint32, void*, void*);
};

struct	Link
{
	Link*	link;
	byte	data[8];
};

struct	Hmap
{
	uint32	len;		// must be first
	uint32	keysize;
	uint32	valsize;
	uint32	hint;
	Alg*	keyalg;
	Alg*	valalg;
	uint32	valoffset;
	uint32	ko;
	uint32	vo;
	uint32	po;
	Link*	link;
};

static uint64
memhash(uint32 s, void *a)
{
	prints("memhash\n");
	return 0x12345;
}

static uint32
memequal(uint32 s, void *a, void *b)
{
	byte *ba, *bb;
	uint32 i;

	ba = a;
	bb = b;
	for(i=0; i<s; i++)
		if(ba[i] != bb[i])
			return 0;
	return 1;
}

static void
memprint(uint32 s, void *a)
{
	uint64 v;

	v = 0xbadb00b;
	switch(s) {
	case 1:
		v = *(uint8*)a;
		break;
	case 2:
		v = *(uint16*)a;
		break;
	case 4:
		v = *(uint32*)a;
		break;
	case 8:
		v = *(uint64*)a;
		break;
	}
	sys·printint(v);
}

static void
memcopy(uint32 s, void *a, void *b)
{
	byte *ba, *bb;
	uint32 i;

	ba = a;
	bb = b;
	if(bb == nil) {
		for(i=0; i<s; i++)
			ba[i] = 0;
		return;
	}
	for(i=0; i<s; i++)
		ba[i] = bb[i];
}

static uint64
stringhash(uint32 s, string *a)
{
	prints("stringhash\n");
	return 0x12345;
}

static uint32
stringequal(uint32 s, string *a, string *b)
{
	return cmpstring(*a, *b) == 0;
}

static void
stringprint(uint32 s, string *a)
{
	sys·printstring(*a);
}

static void
stringcopy(uint32 s, string *a, string *b)
{
	if(b == nil) {
		*a = nil;
		return;
	}
	*a = *b;
}

static uint64
pointerhash(uint32 s, void **a)
{
	prints("pointerhash\n");
	return 0x12345;
}

static uint32
pointerequal(uint32 s, void **a, void **b)
{
	prints("pointerequal\n");
	return 0;
}

static void
pointerprint(uint32 s, void **a)
{
	prints("pointerprint\n");
}

static void
pointercopy(uint32 s, void **a, void **b)
{
	if(b == nil) {
		*a = nil;
		return;
	}
	*a = *b;
}

static uint32
rnd(uint32 n, uint32 m)
{
	uint32 r;

	r = n % m;
	if(r)
		n += m-r;
	return n;
}

static	Alg
algarray[] =
{
	{	&memhash,	&memequal,	&memprint,	&memcopy	},
	{	&stringhash,	&stringequal,	&stringprint,	&stringcopy	},
	{	&pointerhash,	&pointerequal,	&pointerprint,	&pointercopy	},
};

// newmap(keysize uint32, valsize uint32,
//	keyalg uint32, valalg uint32,
//	hint uint32) (hmap *map[any]any);
void
sys·newmap(uint32 keysize, uint32 valsize,
	uint32 keyalg, uint32 valalg, uint32 hint,
	Hmap* ret)
{
	Hmap *m;

	if(keyalg >= 2 ||
	   valalg >= 3) {
		prints("0<=");
		sys·printint(keyalg);
		prints("<");
		sys·printint(nelem(algarray));
		prints("\n0<=");
		sys·printint(valalg);
		prints("<");
		sys·printint(nelem(algarray));
		prints("\n");

		throw("sys·newmap: key/val algorithm out of range");
	}

	m = mal(sizeof(*m));

	m->len = 0;
	m->keysize = keysize;
	m->valsize = valsize;
	m->keyalg = &algarray[keyalg];
	m->valalg = &algarray[valalg];
	m->hint = hint;

	// these calculations are compiler dependent
	m->valoffset = rnd(keysize, valsize);
	m->ko = rnd(sizeof(m), keysize);
	m->vo = rnd(m->ko+keysize, valsize);
	m->po = rnd(m->vo+valsize, 1);

	ret = m;
	FLUSH(&ret);

	if(debug) {
		prints("newmap: map=");
		sys·printpointer(m);
		prints("; keysize=");
		sys·printint(keysize);
		prints("; valsize=");
		sys·printint(valsize);
		prints("; keyalg=");
		sys·printint(keyalg);
		prints("; valalg=");
		sys·printint(valalg);
		prints("; valoffset=");
		sys·printint(m->valoffset);
		prints("; ko=");
		sys·printint(m->ko);
		prints("; vo=");
		sys·printint(m->vo);
		prints("; po=");
		sys·printint(m->po);
		prints("\n");
	}
}

// mapaccess1(hmap *map[any]any, key any) (val any);
void
sys·mapaccess1(Hmap *m, ...)
{
	Link *l;
	byte *ak, *av;

	ak = (byte*)&m + m->ko;
	av = (byte*)&m + m->vo;

	for(l=m->link; l!=nil; l=l->link) {
		if(m->keyalg->equal(m->keysize, ak, l->data)) {
			m->valalg->copy(m->valsize, av, l->data+m->valoffset);
			goto out;
		}
	}

	m->valalg->copy(m->valsize, av, 0);

out:
	if(debug) {
		prints("sys·mapaccess1: map=");
		sys·printpointer(m);
		prints("; key=");
		m->keyalg->print(m->keysize, ak);
		prints("; val=");
		m->valalg->print(m->valsize, av);
		prints("\n");
	}
}

// mapaccess2(hmap *map[any]any, key any) (val any, pres bool);
void
sys·mapaccess2(Hmap *m, ...)
{
	Link *l;
	byte *ak, *av, *ap;

	ak = (byte*)&m + m->ko;
	av = (byte*)&m + m->vo;
	ap = (byte*)&m + m->po;

	for(l=m->link; l!=nil; l=l->link) {
		if(m->keyalg->equal(m->keysize, ak, l->data)) {
			*ap = true;
			m->valalg->copy(m->valsize, av, l->data+m->valoffset);
			goto out;
		}
	}

	*ap = false;
	m->valalg->copy(m->valsize, av, nil);

out:
	if(debug) {
		prints("sys·mapaccess2: map=");
		sys·printpointer(m);
		prints("; key=");
		m->keyalg->print(m->keysize, ak);
		prints("; val=");
		m->valalg->print(m->valsize, av);
		prints("; pres=");
		sys·printbool(*ap);
		prints("\n");
	}
}

static void
sys·mapassign(Hmap *m, byte *ak, byte *av)
{
	Link *l;

	// mapassign(hmap *map[any]any, key any, val any);

	for(l=m->link; l!=nil; l=l->link) {
		if(m->keyalg->equal(m->keysize, ak, l->data))
			goto out;
	}

	l = mal((sizeof(*l)-8) + m->keysize + m->valsize);
	l->link = m->link;
	m->link = l;
	m->keyalg->copy(m->keysize, l->data, ak);
	m->len++;

out:
	m->valalg->copy(m->valsize, l->data+m->valoffset, av);

	if(debug) {
		prints("mapassign: map=");
		sys·printpointer(m);
		prints("; key=");
		m->keyalg->print(m->keysize, ak);
		prints("; val=");
		m->valalg->print(m->valsize, av);
		prints("\n");
	}
}

// mapassign1(hmap *map[any]any, key any, val any);
void
sys·mapassign1(Hmap *m, ...)
{
	Link **ll;
	byte *ak, *av;

	ak = (byte*)&m + m->ko;
	av = (byte*)&m + m->vo;

	sys·mapassign(m, ak, av);
}

// mapassign2(hmap *map[any]any, key any, val any, pres bool);
void
sys·mapassign2(Hmap *m, ...)
{
	Link **ll;
	byte *ak, *av, *ap;

	ak = (byte*)&m + m->ko;
	av = (byte*)&m + m->vo;
	ap = (byte*)&m + m->po;

	if(*ap == true) {
		// assign
		sys·mapassign(m, ak, av);
		return;
	}

	// delete
	for(ll=&m->link; (*ll)!=nil; ll=&(*ll)->link) {
		if(m->keyalg->equal(m->keysize, ak, (*ll)->data)) {
			m->valalg->copy(m->valsize, (*ll)->data+m->valoffset, nil);
			(*ll) = (*ll)->link;
			m->len--;
			if(debug) {
				prints("mapdelete (found): map=");
				sys·printpointer(m);
				prints("; key=");
				m->keyalg->print(m->keysize, ak);
				prints("\n");
			}
			return;
		}
	}

	if(debug) {
		prints("mapdelete (not found): map=");
		sys·printpointer(m);
		prints("; key=");
		m->keyalg->print(m->keysize, ak);
		prints(" *** not found\n");
	}
}
