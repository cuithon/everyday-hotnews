// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
 * range
 */

#include <u.h>
#include <libc.h>
#include "go.h"

void
typecheckrange(Node *n)
{
	char *why;
	Type *t, *t1, *t2;
	Node *v1, *v2;
	NodeList *ll;

	// delicate little dance.  see typecheckas2
	for(ll=n->list; ll; ll=ll->next)
		if(ll->n->defn != n)
			typecheck(&ll->n, Erv | Easgn);

	typecheck(&n->right, Erv);
	if((t = n->right->type) == T)
		goto out;
	if(isptr[t->etype] && isfixedarray(t->type))
		t = t->type;
	n->type = t;

	switch(t->etype) {
	default:
		yyerror("cannot range over %lN", n->right);
		goto out;

	case TARRAY:
		t1 = types[TINT];
		t2 = t->type;
		break;

	case TMAP:
		t1 = t->down;
		t2 = t->type;
		break;

	case TCHAN:
		if(!(t->chan & Crecv)) {
			yyerror("invalid operation: range %N (receive from send-only type %T)", n->right, n->right->type);
			goto out;
		}
		t1 = t->type;
		t2 = nil;
		if(count(n->list) == 2)
			goto toomany;
		break;

	case TSTRING:
		t1 = types[TINT];
		t2 = runetype;
		break;
	}

	if(count(n->list) > 2) {
	toomany:
		yyerror("too many variables in range");
	}

	v1 = n->list->n;
	v2 = N;
	if(n->list->next)
		v2 = n->list->next->n;

	// this is not only a optimization but also a requirement in the spec.
	// "if the second iteration variable is the blank identifier, the range
	// clause is equivalent to the same clause with only the first variable
	// present."
	if(isblank(v2)) {
		n->list = list1(v1);
		v2 = N;
	}

	if(v1->defn == n)
		v1->type = t1;
	else if(v1->type != T && assignop(t1, v1->type, &why) == 0)
		yyerror("cannot assign type %T to %lN in range%s", t1, v1, why);
	if(v2) {
		if(v2->defn == n)
			v2->type = t2;
		else if(v2->type != T && assignop(t2, v2->type, &why) == 0)
			yyerror("cannot assign type %T to %lN in range%s", t2, v2, why);
	}

out:
	typechecklist(n->nbody, Etop);

	// second half of dance
	n->typecheck = 1;
	for(ll=n->list; ll; ll=ll->next)
		if(ll->n->typecheck == 0)
			typecheck(&ll->n, Erv | Easgn);
}

void
walkrange(Node *n)
{
	Node *ohv1, *hv1, *hv2;	// hidden (old) val 1, 2
	Node *ha, *hit;	// hidden aggregate, iterator
	Node *hn, *hp;	// hidden len, pointer
	Node *hb;  // hidden bool
	Node *a, *v1, *v2;	// not hidden aggregate, val 1, 2
	Node *fn, *tmp;
	Node *keyname, *valname;
	Node *key, *val;
	NodeList *body, *init;
	Type *th, *t;
	int lno;

	t = n->type;
	init = nil;

	a = n->right;
	lno = setlineno(a);
	if(t->etype == TSTRING && !eqtype(t, types[TSTRING])) {
		a = nod(OCONV, n->right, N);
		a->type = types[TSTRING];
	}

	v1 = n->list->n;
	v2 = N;
	if(n->list->next)
		v2 = n->list->next->n;
	// n->list has no meaning anymore, clear it
	// to avoid erroneous processing by racewalk.
	n->list = nil;
	hv2 = N;

	if(v2 == N && t->etype == TARRAY) {
		// will have just one reference to argument.
		// no need to make a potentially expensive copy.
		ha = a;
	} else {
		ha = temp(a->type);
		init = list(init, nod(OAS, ha, a));
	}

	switch(t->etype) {
	default:
		fatal("walkrange");

	case TARRAY:
		hv1 = temp(types[TINT]);
		hn = temp(types[TINT]);
		hp = nil;

		init = list(init, nod(OAS, hv1, N));
		init = list(init, nod(OAS, hn, nod(OLEN, ha, N)));
		if(v2) {
			hp = temp(ptrto(n->type->type));
			tmp = nod(OINDEX, ha, nodintconst(0));
			tmp->bounded = 1;
			init = list(init, nod(OAS, hp, nod(OADDR, tmp, N)));
		}

		n->ntest = nod(OLT, hv1, hn);
		n->nincr = nod(OAS, hv1, nod(OADD, hv1, nodintconst(1)));
		if(v2 == N)
			body = list1(nod(OAS, v1, hv1));
		else {
			a = nod(OAS2, N, N);
			a->list = list(list1(v1), v2);
			a->rlist = list(list1(hv1), nod(OIND, hp, N));
			body = list1(a);
			
			// Advance pointer as part of increment.
			// We used to advance the pointer before executing the loop body,
			// but doing so would make the pointer point past the end of the
			// array during the final iteration, possibly causing another unrelated
			// piece of memory not to be garbage collected until the loop finished.
			// Advancing during the increment ensures that the pointer p only points
			// pass the end of the array during the final "p++; i++; if(i >= len(x)) break;",
			// after which p is dead, so it cannot confuse the collector.
			tmp = nod(OADD, hp, nodintconst(t->type->width));
			tmp->type = hp->type;
			tmp->typecheck = 1;
			tmp->right->type = types[tptr];
			tmp->right->typecheck = 1;
			a = nod(OAS, hp, tmp);
			typecheck(&a, Etop);
			n->nincr->ninit = list1(a);
		}
		break;

	case TMAP:
		// allocate an iterator state structure on the stack
		th = hiter(t);
		hit = temp(th);
		init = list(init, nod(OAS, hit, N));
		keyname = newname(th->type->sym);  // depends on layout of iterator struct.  See reflect.c:hiter
		valname = newname(th->type->down->sym); // ditto

		fn = syslook("mapiterinit", 1);
		argtype(fn, t->down);
		argtype(fn, t->type);
		argtype(fn, th);
		init = list(init, mkcall1(fn, T, nil, typename(t), ha, nod(OADDR, hit, N)));
		n->ntest = nod(ONE, nod(ODOT, hit, keyname), nodnil());

		fn = syslook("mapiternext", 1);
		argtype(fn, th);
		n->nincr = mkcall1(fn, T, nil, nod(OADDR, hit, N));

		key = nod(ODOT, hit, keyname);
		key = nod(OIND, key, N);
		if(v2 == N) {
			a = nod(OAS, v1, key);
		} else {
			val = nod(ODOT, hit, valname);
			val = nod(OIND, val, N);
			a = nod(OAS2, N, N);
			a->list = list(list1(v1), v2);
			a->rlist = list(list1(key), val);
		}
		body = list1(a);
		break;

	case TCHAN:
		hv1 = temp(t->type);
		if(haspointers(t->type))
			init = list(init, nod(OAS, hv1, N));
		hb = temp(types[TBOOL]);

		n->ntest = nod(ONE, hb, nodbool(0));
		a = nod(OAS2RECV, N, N);
		a->typecheck = 1;
		a->list = list(list1(hv1), hb);
		a->rlist = list1(nod(ORECV, ha, N));
		n->ntest->ninit = list1(a);
		body = list1(nod(OAS, v1, hv1));
		break;

	case TSTRING:
		ohv1 = temp(types[TINT]);

		hv1 = temp(types[TINT]);
		init = list(init, nod(OAS, hv1, N));

		if(v2 == N)
			a = nod(OAS, hv1, mkcall("stringiter", types[TINT], nil, ha, hv1));
		else {
			hv2 = temp(runetype);
			a = nod(OAS2, N, N);
			a->list = list(list1(hv1), hv2);
			fn = syslook("stringiter2", 0);
			a->rlist = list1(mkcall1(fn, getoutargx(fn->type), nil, ha, hv1));
		}
		n->ntest = nod(ONE, hv1, nodintconst(0));
		n->ntest->ninit = list(list1(nod(OAS, ohv1, hv1)), a);

		body = list1(nod(OAS, v1, ohv1));
		if(v2 != N)
			body = list(body, nod(OAS, v2, hv2));
		break;
	}

	n->op = OFOR;
	typechecklist(init, Etop);
	n->ninit = concat(n->ninit, init);
	typechecklist(n->ntest->ninit, Etop);
	typecheck(&n->ntest, Erv);
	typecheck(&n->nincr, Etop);
	typechecklist(body, Etop);
	n->nbody = concat(body, n->nbody);
	walkstmt(&n);
	
	lineno = lno;
}

