// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
 * type check the whole tree of an expression.
 * calculates expression types.
 * evaluates compile time constants.
 * marks variables that escape the local frame.
 * rewrites n->op to be more specific in some cases.
 * sets n->walk to walking function.
 *
 * TODO:
 *	trailing ... section of function calls
 *	statements
 */

#include "go.h"

static void	implicitstar(Node**);
static int	onearg(Node*);
static int	lookdot(Node*, Type*);
static void	typecheckaste(int, Type*, NodeList*);
static int	exportassignok(Type*);
static Type*	lookdot1(Sym *s, Type *t, Type *f);
static int	nokeys(NodeList*);
static void	typecheckcomplit(Node**);
static void	addrescapes(Node*);

static void	checklvalue(Node*, char*);
static void checkassign(Node*);
static void checkassignlist(NodeList*);
static int	islvalue(Node*);

void
typechecklist(NodeList *l, int top)
{
	for(; l; l=l->next)
		typecheck(&l->n, top);
}

/*
 * type check node *np.
 * replaces *np with a new pointer in some cases.
 * returns the final value of *np as a convenience.
 */
Node*
typecheck(Node **np, int top)
{
	int et, op;
	Node *n, *l, *r;
	NodeList *args;
	int lno, ok;
	Type *t;

	n = *np;
	if(n == N)
		return N;
	if(n->typecheck == 1 && n->op != ONAME)	// XXX for test/func4.go
		return n;
	if(n->typecheck == 2)
		fatal("typecheck loop");
	n->typecheck = 2;

	if(n->sym && n->walkdef != 1)
		walkdef(n);

	lno = setlineno(n);

reswitch:
	ok = 0;
	switch(n->op) {
	default:
		// until typecheck is complete, do nothing.
		dump("typecheck", n);
		fatal("typecheck %O", n->op);

	/*
	 * names
	 */
	case OLITERAL:
		ok |= Erv;
		goto ret;

	case ONONAME:
		ok |= Erv;
		goto ret;

	case ONAME:
		if(n->etype != 0) {
			yyerror("must call builtin %S", n->sym);
			goto error;
		}
		ok |= Erv;
		goto ret;

	/*
	 * types (OIND is with exprs)
	 */
	case OTYPE:
		ok |= Etype;
		if(n->type == T)
			goto error;
		break;

	case OTARRAY:
		ok |= Etype;
		t = typ(TARRAY);
		l = n->left;
		r = n->right;
		if(l == nil) {
			t->bound = -1;
		} else {
			typecheck(&l, Erv | Etype);
			switch(l->op) {
			default:
				yyerror("invalid array bound %#N", l);
				goto error;

			case OLITERAL:
				if(consttype(l) == CTINT) {
					t->bound = mpgetfix(l->val.u.xval);
					if(t->bound < 0) {
						yyerror("array bound must be non-negative");
						goto error;
					}
				}
				break;

			case OTYPE:
				if(l->type == T)
					goto error;
				if(l->type->etype != TDDD) {
					yyerror("invalid array bound %T", l->type);
					goto error;
				}
				t->bound = -100;
				break;
			}
		}
		typecheck(&r, Etype);
		if(r->type == T)
			goto error;
		t->type = r->type;
		n->op = OTYPE;
		n->type = t;
		n->left = N;
		n->right = N;
		checkwidth(t);
		break;

	case OTMAP:
		ok |= Etype;
		l = typecheck(&n->left, Etype);
		r = typecheck(&n->right, Etype);
		if(l->type == T || r->type == T)
			goto error;
		n->op = OTYPE;
		n->type = maptype(l->type, r->type);
		n->left = N;
		n->right = N;
		break;

	case OTCHAN:
		ok |= Etype;
		l = typecheck(&n->left, Etype);
		if(l->type == T)
			goto error;
		t = typ(TCHAN);
		t->type = l->type;
		t->chan = n->etype;
		n->op = OTYPE;
		n->type = t;
		n->left = N;
		n->etype = 0;
		break;

	case OTSTRUCT:
		ok |= Etype;
		n->op = OTYPE;
		n->type = dostruct(n->list, TSTRUCT);
		if(n->type == T)
			goto error;
		n->list = nil;
		break;

	case OTINTER:
		ok |= Etype;
		n->op = OTYPE;
		n->type = dostruct(n->list, TINTER);
		if(n->type == T)
			goto error;
		n->type = sortinter(n->type);
		break;

	case OTFUNC:
		ok |= Etype;
		n->op = OTYPE;
		n->type = functype(n->left, n->list, n->rlist);
		if(n->type == T)
			goto error;
		break;

	/*
	 * type or expr
	 */
	case OIND:
		l = typecheck(&n->left, Erv | Etype);
		if((t = l->type) == T)
			goto error;
		if(l->op == OTYPE) {
			ok |= Etype;
			n->op = OTYPE;
			n->type = ptrto(l->type);
			n->left = N;
			goto ret;
		}
		if(!isptr[t->etype]) {
			yyerror("invalid indirect of %+N", n);
			goto error;
		}
		ok |= Erv;
		n->type = t->type;
		goto ret;

	/*
	 * arithmetic exprs
	 */
	case OASOP:
		ok |= Etop;
		l = typecheck(&n->left, Erv);
		checkassign(n->left);
		r = typecheck(&n->right, Erv);
		if(l->type == T || r->type == T)
			goto error;
		op = n->etype;
		goto arith;

	case OADD:
	case OAND:
	case OANDAND:
	case OANDNOT:
	case ODIV:
	case OEQ:
	case OGE:
	case OGT:
	case OLE:
	case OLT:
	case OLSH:
	case ORSH:
	case OMOD:
	case OMUL:
	case ONE:
	case OOR:
	case OOROR:
	case OSUB:
	case OXOR:
		ok |= Erv;
		l = typecheck(&n->left, Erv);
		r = typecheck(&n->right, Erv);
		if(l->type == T || r->type == T)
			goto error;
		op = n->op;
	arith:
		if(op == OLSH || op == ORSH)
			goto shift;
		// ideal mixed with non-ideal
		defaultlit2(&l, &r, 0);
		n->left = l;
		n->right = r;
		t = l->type;
		if(t->etype == TIDEAL)
			t = r->type;
		et = t->etype;
		if(et == TIDEAL)
			et = TINT;
		if(t->etype != TIDEAL && !eqtype(l->type, r->type)) {
		badbinary:
			defaultlit2(&l, &r, 1);
			yyerror("invalid operation: %#N (type %T %#O %T)", n, l->type, op, r->type);
			goto error;
		}
		if(!okfor[op][et])
			goto badbinary;
		// okfor allows any array == array;
		// restrict to slice == nil and nil == slice.
		if(l->type->etype == TARRAY && !isslice(l->type))
			goto badbinary;
		if(r->type->etype == TARRAY && !isslice(r->type))
			goto badbinary;
		if(isslice(l->type) && !isnil(l) && !isnil(r))
			goto badbinary;
		t = l->type;
		if(iscmp[n->op]) {
			t = types[TBOOL];
			evconst(n);
			if(n->op != OLITERAL) {
				defaultlit2(&l, &r, 1);
				n->left = l;
				n->right = r;
			}
		}
		if(et == TSTRING) {
			if(iscmp[n->op]) {
				n->etype = n->op;
				n->op = OCMPSTR;
			} else if(n->op == OASOP)
				n->op = OAPPENDSTR;
			else if(n->op == OADD)
				n->op = OADDSTR;
		}
		if(et == TINTER) {
			n->etype = n->op;
			n->op = OCMPIFACE;
		}
		n->type = t;
		goto ret;

	shift:
		defaultlit(&r, types[TUINT]);
		n->right = r;
		t = r->type;
		if(!isint[t->etype] || issigned[t->etype]) {
			yyerror("invalid operation: %#N (shift count type %T)", n, r->type);
			goto error;
		}
		// no defaultlit for left
		// the outer context gives the type
		n->type = l->type;
		goto ret;

	case OCOM:
	case OMINUS:
	case ONOT:
	case OPLUS:
		ok |= Erv;
		l = typecheck(&n->left, Erv);
		if((t = l->type) == T)
			goto error;
		if(!okfor[n->op][t->etype]) {
			yyerror("invalid operation: %#O %T", n->op, t);
			goto error;
		}
		n->type = t;
		goto ret;

	/*
	 * exprs
	 */
	case OADDR:
		typecheck(&n->left, Erv);
		if(n->left->type == T)
			goto error;
		switch(n->left->op) {
		case OMAPLIT:
		case OSTRUCTLIT:
		case OARRAYLIT:
			break;
		default:
			checklvalue(n->left, "take the address of");
		}
		defaultlit(&n->left, T);
		l = n->left;
		if((t = l->type) == T)
			goto error;
		addrescapes(n->left);
		n->type = ptrto(t);
		goto ret;

	case OCOMPLIT:
		typecheckcomplit(&n);
		if(n->type == T)
			goto error;
		goto ret;

	case ODOT:
		l = typecheck(&n->left, Erv);
		if((t = l->type) == T)
			goto error;
		if(n->right->op != ONAME) {
			yyerror("rhs of . must be a name");	// impossible
			goto error;
		}
		if(isptr[t->etype]) {
			t = t->type;
			if(t == T)
				goto error;
			n->op = ODOTPTR;
		}
		if(!lookdot(n, t)) {
			yyerror("%#N undefined (type %T has no field %S)", n, t, n->right->sym);
			goto error;
		}
		switch(n->op) {
		case ODOTINTER:
		case ODOTMETH:
			ok |= Ecall;
			break;
		default:
			ok |= Erv;
			break;
		}
		goto ret;

	case ODOTTYPE:
		typecheck(&n->left, Erv);
		defaultlit(&n->left, T);
		l = n->left;
		if((t = l->type) == T)
			goto error;
		if(!isinter(t)) {
			yyerror("invalid type assertion: %#N (non-interface type %T on left)", n, t);
			goto error;
		}
		if(n->right != N) {
			typecheck(&n->right, Etype);
			n->type = n->right->type;
			n->right = N;
			if(n->type == T)
				goto error;
		}
		goto ret;

	case OINDEX:
		typecheck(&n->left, Erv);
		defaultlit(&n->left, T);
		implicitstar(&n->left);
		l = n->left;
		typecheck(&n->right, Erv);
		r = n->right;
		if((t = l->type) == T || r->type == T)
			goto error;
		switch(t->etype) {
		default:
			yyerror("invalid operation: %#N (index of type %T)", n, t);
			goto error;

		case TARRAY:
			ok |= Erv;
			defaultlit(&n->right, types[TUINT]);
			n->type = t->type;
			break;

		case TMAP:
			n->etype = 0;
			ok |= Erv;
			defaultlit(&n->right, t->down);
			n->type = t->type;
			n->op = OINDEXMAP;
			break;

		case TSTRING:
			ok |= Erv;
			defaultlit(&n->right, types[TUINT]);
			n->type = types[TUINT8];
			n->op = OINDEXSTR;
			break;
		}
		goto ret;

	case ORECV:
		typecheck(&n->left, Erv);
		defaultlit(&n->left, T);
		l = n->left;
		if((t = l->type) == T)
			goto error;
		if(t->etype != TCHAN) {
			yyerror("invalid operation: %#N (receive from non-chan type %T)", n, t);
			goto error;
		}
		if(!(t->chan & Crecv)) {
			yyerror("invalid operation: %#N (receive from send-only type %T)", n, t);
			goto error;
		}
		n->type = t->type;
		ok |= Erv;
		goto ret;

	case OSEND:
		l = typecheck(&n->left, Erv);
		typecheck(&n->right, Erv);
		defaultlit(&n->left, T);
		l = n->left;
		if((t = l->type) == T)
			goto error;
		if(!(t->chan & Csend)) {
			yyerror("invalid operation: %#N (send to receive-only type %T)", n, t);
			goto error;
		}
		defaultlit(&n->right, t->type);
		r = n->right;
		if((t = r->type) == T)
			goto error;
		// TODO: more aggressive
		n->etype = 0;
		if(top & Erv)
			n->op = OSENDNB;
		ok |= Etop | Erv;
		n->type = types[TBOOL];
		goto ret;

	case OSLICE:
		ok |= Erv;
		typecheck(&n->left, top);
		typecheck(&n->right->left, Erv);
		typecheck(&n->right->right, Erv);
		defaultlit(&n->left, T);
		defaultlit(&n->right->left, types[TUINT]);
		defaultlit(&n->right->right, types[TUINT]);
		implicitstar(&n->left);
		if(n->right->left == N || n->right->right == N) {
			yyerror("missing slice bounds?");
			goto error;
		}
		if((t = n->right->left->type) == T)
			goto error;
		if(!isint[t->etype]) {
			yyerror("invalid slice index %#N (type %T)", n->right->left, t);
			goto error;
		}
		if((t = n->right->right->type) == T)
			goto error;
		if(!isint[t->etype]) {
			yyerror("invalid slice index %#N (type %T)", n->right->right, t);
			goto error;
		}
		l = n->left;
		if((t = l->type) == T)
			goto error;
		// TODO(rsc): 64-bit slice index needs to be checked
		// for overflow in generated code
		if(istype(t, TSTRING)) {
			n->type = t;
			n->op = OSLICESTR;
			goto ret;
		}
		if(isfixedarray(t)) {
			n->type = typ(TARRAY);
			n->type->type = t->type;
			n->type->bound = -1;
			dowidth(n->type);
			n->op = OSLICEARR;
			goto ret;
		}
		if(isslice(t)) {
			n->type = t;
			goto ret;
		}
		yyerror("cannot slice %#N (type %T)", l, t);
		goto error;

	/*
	 * call and call like
	 */
	case OCALL:
		l = n->left;
		if(l->op == ONAME && l->etype != 0) {
			// builtin: OLEN, OCAP, etc.
			n->op = l->etype;
			n->left = n->right;
			n->right = N;
			goto reswitch;
		}
		typecheck(&n->left, Erv | Etype | Ecall);
		defaultlit(&n->left, T);
		l = n->left;
		typechecklist(n->list, Erv);
		if((t = l->type) == T)
			goto error;
		dowidth(t);

		switch(l->op) {
		case OTYPE:
			ok |= Erv;
			// turn CALL(type, arg) into CONV(arg) w/ type
			n->left = N;
			if(onearg(n) < 0)
				goto error;
			n->op = OCONV;
			n->type = l->type;
			goto doconv;

		case ODOTINTER:
			n->op = OCALLINTER;
			break;

		case ODOTMETH:
			n->op = OCALLMETH;
			typecheckaste(OCALL, getthisx(t), list1(l->left));
			break;

		default:
			n->op = OCALLFUNC;
			if(t->etype != TFUNC) {
				yyerror("cannot call non-function %#N (type %T)", l, t);
				goto error;
			}
			break;
		}
		typecheckaste(OCALL, getinargx(t), n->list);
		if(t->outtuple == 0) {
			ok |= Etop;
			goto ret;
		}
		if(t->outtuple == 1) {
			ok |= Erv;
			t = getoutargx(l->type)->type;
			if(t->etype == TFIELD)
				t = t->type;
			n->type = t;
			goto ret;
		}
		// multiple return
		// ok |= Emulti;
		n->type = getoutargx(l->type);
		goto ret;

	case OCAP:
	case OLEN:
		if(onearg(n) < 0)
			goto error;
		typecheck(&n->left, Erv);
		defaultlit(&n->left, T);
		implicitstar(&n->left);
		l = n->left;
		if((t = l->type) == T)
			goto error;
		switch(n->op) {
		case OCAP:
			if(!okforcap[t->etype])
				goto badcall1;
			break;
		case OLEN:
			if(!okforlen[t->etype])
				goto badcall1;
			break;
		}
		// might be constant
		switch(t->etype) {
		case TSTRING:
			if(isconst(l, CTSTR))
				nodconst(n, types[TINT], l->val.u.sval->len);
			break;
		case TARRAY:
			if(t->bound >= 0)
				nodconst(n, types[TINT], t->bound);
			break;
		}
		n->type = types[TINT];
		goto ret;

	case OCLOSED:
	case OCLOSE:
		if(onearg(n) < 0)
			goto error;
		typecheck(&n->left, Erv);
		defaultlit(&n->left, T);
		l = n->left;
		if((t = l->type) == T)
			goto error;
		if(t->etype != TCHAN) {
			yyerror("invalid operation: %#N (non-chan type %T)", n, t);
			goto error;
		}
		if(n->op == OCLOSED) {
			n->type = types[TBOOL];
			ok |= Erv;
		} else
			ok |= Etop;
		goto ret;

	case OCONV:
	doconv:
		typecheck(&n->left, Erv);
		defaultlit(&n->left, n->type);
		if((t = n->left->type) == T || n->type == T)
			goto error;
		n = typecheckconv(n, n->left, n->type, 1);
		if(n->type == T)
			goto error;
		goto ret;

	case OMAKE:
		args = n->list;
		if(args == nil) {
			yyerror("missing argument to make");
			goto error;
		}
		l = args->n;
		args = args->next;
		typecheck(&l, Etype);
		if((t = l->type) == T)
			goto error;

		switch(t->etype) {
		default:
		badmake:
			yyerror("cannot make type %T", t);
			goto error;

		case TARRAY:
			if(!isslice(t))
				goto badmake;
			if(args == nil) {
				yyerror("missing len argument to make(%T)", t);
				goto error;
			}
			l = args->n;
			args = args->next;
			typecheck(&l, Erv);
			defaultlit(&l, types[TUINT]);
			r = N;
			if(args != nil) {
				r = args->n;
				args = args->next;
				typecheck(&r, Erv);
				defaultlit(&r, types[TUINT]);
			}
			if(l->type == T || (r && r->type == T))
				goto error;
			if(!isint[l->type->etype]) {
				yyerror("non-integer len argument to make(%T)", t);
				goto error;
			}
			if(r && !isint[r->type->etype]) {
				yyerror("non-integer cap argument to make(%T)", t);
				goto error;
			}
			if(r == N)
				r = nodintconst(0);
			n->left = l;
			n->right = r;
			n->op = OMAKESLICE;
			break;

		case TMAP:
			if(args != nil) {
				l = args->n;
				args = args->next;
				typecheck(&l, Erv);
				defaultlit(&l, types[TUINT]);
				if(l->type == T)
					goto error;
				if(!isint[l->type->etype]) {
					yyerror("non-integer size argument to make(%T)", t);
					goto error;
				}
				n->left = l;
			} else
				n->left = nodintconst(0);
			n->op = OMAKEMAP;
			break;

		case TCHAN:
			l = N;
			if(args != nil) {
				l = args->n;
				args = args->next;
				typecheck(&l, Erv);
				defaultlit(&l, types[TUINT]);
				if(l->type == T)
					goto error;
				if(!isint[l->type->etype]) {
					yyerror("non-integer buffer argument to make(%T)", t);
					goto error;
				}
				n->left = l;
			} else
				n->left = nodintconst(0);
			n->op = OMAKECHAN;
			break;
		}
		if(args != nil) {
			yyerror("too many arguments to make(%T)", t);
			n->op = OMAKE;
			goto error;
		}
		n->type = t;
		goto ret;

	case ONEW:
		args = n->list;
		if(args == nil) {
			yyerror("missing argument to new");
			goto error;
		}
		l = args->n;
		typecheck(&l, Etype);
		if((t = l->type) == T)
			goto error;
		if(args->next != nil) {
			yyerror("too many arguments to new(%T)", t);
			goto error;
		}
		n->left = l;
		n->type = ptrto(t);
		goto ret;

	case OPANIC:
	case OPANICN:
	case OPRINT:
	case OPRINTN:
		typechecklist(n->list, Erv);
		goto ret;

	/*
	 * statements
	 */
	case OAS:
		typecheck(&n->left, Erv);
		checkassign(n->left);
		typecheck(&n->right, Erv);
		if(n->left->type != T && n->right && n->right->type != T)
			n->right = typecheckconv(nil, n->right, n->left->type, 0);
		goto ret;

	case OAS2:
		typechecklist(n->list, Erv);
		checkassignlist(n->list);
		typechecklist(n->rlist, Erv);
		goto ret;

	case OBREAK:
	case OCONTINUE:
	case ODCL:
	case OEMPTY:
	case OGOTO:
	case OLABEL:
	case OXFALL:
		goto ret;

	case ODEFER:
	case OPROC:
		typecheck(&n->left, Etop);
		goto ret;

	case OFOR:
		typechecklist(n->ninit, Etop);
		typecheck(&n->ntest, Erv);	// TODO Ebool
		typecheck(&n->nincr, Etop);
		typechecklist(n->nbody, Etop);
		goto ret;

	case OIF:
		typechecklist(n->ninit, Etop);
		typecheck(&n->ntest, Erv);	// TODO Ebool
		typechecklist(n->nbody, Etop);
		typechecklist(n->nelse, Etop);
		goto ret;

	case ORETURN:
		typechecklist(n->list, Erv);
		if(curfn->type->outnamed && n->list == nil)
			goto ret;
		typecheckaste(ORETURN, getoutargx(curfn->type), n->list);
		goto ret;

	case OSELECT:
		typechecklist(n->ninit, Etop);
		typecheck(&n->ntest, Erv);
		typechecklist(n->list, Etop);
		goto ret;

	case OSWITCH:
		typechecklist(n->ninit, Etop);
		typecheck(&n->ntest, Erv);
		typechecklist(n->list, Etop);
		goto ret;

	case OTYPECASE:
		typecheck(&n->left, Erv);
		goto ret;

	case OTYPESW:
		typecheck(&n->right, Erv);
		goto ret;

	case OXCASE:
		typechecklist(n->list, Erv);
		typechecklist(n->nbody, Etop);
		goto ret;
	}

ret:
	evconst(n);
	if(n->op == OTYPE && !(top & Etype)) {
		yyerror("type %T is not an expression", n->type);
		goto error;
	}
	if((top & (Erv|Etype)) == Etype && n->op != OTYPE) {
		yyerror("%O is not a type", n->op);
		goto error;
	}
	if((ok & Ecall) && !(top & Ecall)) {
		yyerror("must call method %#N", n);
		goto error;
	}

	/* TODO
	if(n->type == T)
		fatal("typecheck nil type");
	*/
	goto out;

badcall1:
	yyerror("invalid argument %#N (type %T) for %#O", n->left, n->left->type, n->op);
	goto error;

error:
	n->type = T;

out:
	lineno = lno;
	n->typecheck = 1;
	*np = n;
	return n;
}

static void
implicitstar(Node **nn)
{
	Type *t;
	Node *n;

	// insert implicit * if needed
	n = *nn;
	t = n->type;
	if(t == T || !isptr[t->etype])
		return;
	t = t->type;
	if(t == T)
		return;
	if(!isfixedarray(t))
		return;
	n = nod(OIND, n, N);
	typecheck(&n, Erv);
	*nn = n;
}

static int
onearg(Node *n)
{
	if(n->left != N)
		return 0;
	if(n->list == nil) {
		yyerror("missing argument to %#O - %#N", n->op, n);
		return -1;
	}
	n->left = n->list->n;
	if(n->list->next != nil) {
		yyerror("too many arguments to %#O", n->op);
		n->list = nil;
		return -1;
	}
	n->list = nil;
	return 0;
}

static Type*
lookdot1(Sym *s, Type *t, Type *f)
{
	Type *r;

	r = T;
	for(; f!=T; f=f->down) {
		if(f->sym != s)
			continue;
		if(r != T) {
			yyerror("ambiguous DOT reference %T.%S", t, s);
			break;
		}
		r = f;
	}
	return r;
}

static int
lookdot(Node *n, Type *t)
{
	Type *f1, *f2, *tt, *rcvr;
	Sym *s;

	s = n->right->sym;

	f1 = T;
	if(t->etype == TSTRUCT || t->etype == TINTER)
		f1 = lookdot1(s, t, t->type);

	f2 = methtype(n->left->type);
	if(f2 != T)
		f2 = lookdot1(s, f2, f2->method);

	if(f1 != T) {
		if(f2 != T)
			yyerror("ambiguous DOT reference %S as both field and method",
				n->right->sym);
		n->xoffset = f1->width;
		n->type = f1->type;
		if(t->etype == TINTER) {
			if(isptr[n->left->type->etype]) {
				n->left = nod(OIND, n->left, N);	// implicitstar
				typecheck(&n->left, Erv);
			}
			n->op = ODOTINTER;
		}
		return 1;
	}

	if(f2 != T) {
		tt = n->left->type;
		rcvr = getthisx(f2->type)->type->type;
		if(!eqtype(rcvr, tt)) {
			if(rcvr->etype == tptr && eqtype(rcvr->type, tt)) {
				typecheck(&n->left, Erv);
				checklvalue(n->left, "call pointer method on");
				addrescapes(n->left);
				n->left = nod(OADDR, n->left, N);
				typecheck(&n->left, Erv);
			} else if(tt->etype == tptr && eqtype(tt->type, rcvr)) {
				n->left = nod(OIND, n->left, N);
				typecheck(&n->left, Erv);
			} else {
				// method is attached to wrong type?
				fatal("method mismatch: %T for %T", rcvr, tt);
			}
		}
		n->right = methodname(n->right, n->left->type);
		n->xoffset = f2->width;
		n->type = f2->type;
		n->op = ODOTMETH;
		return 1;
	}

	return 0;
}

static int
nokeys(NodeList *l)
{
	for(; l; l=l->next)
		if(l->n->op == OKEY)
			return 0;
	return 1;
}

Node*
typecheckconv(Node *nconv, Node *n, Type *t, int explicit)
{
	int et, op;
	Node *n1;

	op = OCONV;
	et = 0;

	// preexisting error
	if(t == T || t->etype == TFORW)
		return n;

	/*
	 * implicit conversions
	 */

	convlit1(&n, t, explicit);
	if(n->type == T)
		return n;

	if(eqtype(t, n->type)) {
		exportassignok(t);
		op = OCONVNOP;
		if(!explicit || t == n->type)
			return n;
		goto conv;
	}

	// interfaces are not subject to the name restrictions below.
	// accept anything involving interfaces and let walkiface
	// generate a good message.  some messages have to be
	// delayed anyway.
	if(isnilinter(t) || isnilinter(n->type) || isinter(t) || isinter(n->type)) {
		et = ifaceas1(t, n->type, 0);
		op = OCONVIFACE;
		goto conv;
	}

	// otherwise, if concrete types have names, they must match.
	if(!explicit && t->sym && n->type->sym && t != n->type)
		goto badimplicit;

	// channel must not lose directionality
	if(t->etype == TCHAN && n->type->etype == TCHAN) {
		if(t->chan & ~n->type->chan) {
			if(!explicit)
				goto badimplicit;
			goto badexplicit;
		}
		if(eqtype(t->type, n->type->type)) {
			op = OCONVNOP;
			goto conv;
		}
	}

	// array to slice
	if(isslice(t) && isptr[n->type->etype] && isfixedarray(n->type->type)
	&& eqtype(t->type, n->type->type->type)) {
		op = OCONVSLICE;
		goto conv;
	}

	if(!explicit) {
	badimplicit:
		yyerror("cannot use %+N as type %T", n, t);
		n = nod(OCONV, n, N);	// leave type == T
		n->typecheck = 1;
		return n;
	}

	/*
	 * explicit conversions
	 */

	// same representation
	if(cvttype(t, n->type)) {
		if(n->op == OLITERAL) {
			// can convert literal in place
			n1 = nod(OXXX, N, N);
			*n1 = *n;
			n1->type = t;
			return n1;
		}
		op = OCONVNOP;
		goto conv;
	}

	// simple fix-float
	if(isint[t->etype] || isfloat[t->etype])
	if(isint[n->type->etype] || isfloat[n->type->etype]) {
		// evconst(n);	// XXX is this needed?
		goto conv;
	}

	// to string
	if(istype(t, TSTRING)) {
		// integer rune
		if(isint[n->type->etype]) {
			op = ORUNESTR;
			goto conv;
		}

		// *[10]byte -> string?  convert *[10]byte -> []byte
		// in preparation for next step
		if(isptr[n->type->etype] && isfixedarray(n->type->type)) {
			switch(n->type->type->type->etype) {
			case TUINT8:
			case TINT:
				n1 = nod(OCONV, n, N);
				n1->type = typ(TARRAY);
				n1->type->bound = -1;
				n1->type->type = n->type->type->type;
				dowidth(n1->type);
				typecheck(&n1, Erv);
				walkexpr(&n1, nil);
				n = n1;
				break;
			}
		}

		// []byte -> string
		if(isslice(n->type)) {
			switch(n->type->type->etype) {
			case TUINT8:
				op = OARRAYBYTESTR;
				goto conv;
			case TINT:
				op = OARRAYRUNESTR;
				goto conv;
			}
		}
	}

	// convert to unsafe pointer
	if(isptrto(t, TANY)
	&& (isptr[n->type->etype] || n->type->etype == TUINTPTR))
		goto conv;

	// convert from unsafe pointer
	if(isptrto(n->type, TANY)
	&& (isptr[t->etype] || t->etype == TUINTPTR))
		goto conv;

badexplicit:
	yyerror("cannot convert %+N to type %T", n, t);
	nconv->type = T;
	return nconv;

conv:
	if(nconv == nil) {
		nconv = nod(OXXX, n, N);
		nconv->type = t;
		nconv->typecheck = 1;
	}
	nconv->etype = et;
	nconv->op = op;
	return nconv;
}

/*
 * typecheck assignment: type list = type list
 */
static void
typecheckastt(int op, Type *t1, Type *t2)
{
	for(t1=t1->type, t2=t2->type; t1; t1=t1->down, t2=t2->down) {
		if(t2 == nil) {
			yyerror("too few");
			return;
		}
		if(!eqtype(t1->type, t2->type)) {
			yyerror("wrong");
		}
	}
	if(t2 != nil)
		yyerror("too many");
}

/*
 * typecheck assignment: type list = expression list
 */
static void
typecheckaste(int op, Type *tstruct, NodeList *nl)
{
	Type *t, *tl;
	Node *n;

	if(nl != nil && nl->next == nil && nl->n->type != T && nl->n->type->etype == TSTRUCT && nl->n->type->funarg) {
		typecheckastt(op, tstruct, nl->n->type);
		return;
	}

	for(tl=tstruct->type; tl; tl=tl->down) {
		t = tl->type;
		if(isddd(t)) {
			for(; nl; nl=nl->next)
				defaultlit(&nl->n, T);
			return;
		}
		if(nl == nil) {
			yyerror("not enough arguments to %#O", op);
			return;
		}
		n = nl->n;
		if(n->type != T)
			nl->n = typecheckconv(nil, n, t, 0);
		nl = nl->next;
	}
	if(nl != nil) {
		yyerror("too many arguments to %#O", op);
		return;
	}
}

/*
 * do the export rules allow writing to this type?
 * cannot be implicitly assigning to any type with
 * an unavailable field.
 */
static int
exportassignok(Type *t)
{
	Type *f;
	Sym *s;

	if(t == T)
		return 1;
	switch(t->etype) {
	default:
		// most types can't contain others; they're all fine.
		break;
	case TSTRUCT:
		for(f=t->type; f; f=f->down) {
			if(f->etype != TFIELD)
				fatal("structas: not field");
			s = f->sym;
			// s == nil doesn't happen for embedded fields (they get the type symbol).
			// it only happens for fields in a ... struct.
			if(s != nil && !exportname(s->name) && strcmp(package, s->package) != 0) {
				yyerror("implicit assignment of %T field '%s'", t, s->name);
				return 0;
			}
			if(!exportassignok(f->type))
				return 0;
		}
		break;

	case TARRAY:
		if(t->bound < 0)	// slices are pointers; that's fine
			break;
		if(!exportassignok(t->type))
			return 0;
		break;
	}
	return 1;
}


/*
 * type check composite
 */

static void
fielddup(Node *n, Node *hash[], ulong nhash)
{
	uint h;
	char *s;
	Node *a;

	if(n->op != ONAME)
		fatal("fielddup: not ONAME");
	s = n->sym->name;
	h = stringhash(s)%nhash;
	for(a=hash[h]; a!=N; a=a->ntest) {
		if(strcmp(a->sym->name, s) == 0) {
			yyerror("duplicate field name in struct literal: %s", s);
			return;
		}
	}
	n->ntest = hash[h];
	hash[h] = n;
}

static void
keydup(Node *n, Node *hash[], ulong nhash)
{
	uint h;
	ulong b;
	double d;
	int i;
	Node *a;
	Node cmp;
	char *s;

	evconst(n);
	if(n->op != OLITERAL)
		return;	// we dont check variables

	switch(n->val.ctype) {
	default:	// unknown, bool, nil
		b = 23;
		break;
	case CTINT:
		b = mpgetfix(n->val.u.xval);
		break;
	case CTFLT:
		d = mpgetflt(n->val.u.fval);
		s = (char*)&d;
		b = 0;
		for(i=sizeof(d); i>0; i--)
			b = b*PRIME1 + *s++;
		break;
	case CTSTR:
		b = 0;
		s = n->val.u.sval->s;
		for(i=n->val.u.sval->len; i>0; i--)
			b = b*PRIME1 + *s++;
		break;
	}

	h = b%nhash;
	memset(&cmp, 0, sizeof(cmp));
	for(a=hash[h]; a!=N; a=a->ntest) {
		cmp.op = OEQ;
		cmp.left = n;
		cmp.right = a;
		evconst(&cmp);
		b = cmp.val.u.bval;
		if(b) {
			// too lazy to print the literal
			yyerror("duplicate key in map literal");
			return;
		}
	}
	n->ntest = hash[h];
	hash[h] = n;
}

static void
indexdup(Node *n, Node *hash[], ulong nhash)
{
	uint h;
	Node *a;
	ulong b, c;

	if(n->op != OLITERAL)
		fatal("indexdup: not OLITERAL");

	b = mpgetfix(n->val.u.xval);
	h = b%nhash;
	for(a=hash[h]; a!=N; a=a->ntest) {
		c = mpgetfix(a->val.u.xval);
		if(b == c) {
			yyerror("duplicate index in array literal: %ld", b);
			return;
		}
	}
	n->ntest = hash[h];
	hash[h] = n;
}

static void
typecheckcomplit(Node **np)
{
	int bad, i, len, nerr;
	Node *l, *n, *hash[101];
	NodeList *ll;
	Type *t, *f;

	n = *np;

	memset(hash, 0, sizeof hash);

	// TODO: dup detection
	l = typecheck(&n->right /* sic */, Etype /* TODO | Edotarray */);
	if((t = l->type) == T)
		goto error;
	nerr = nerrors;
	switch(t->etype) {
	default:
		yyerror("invalid type for composite literal: %T", t);
		n->type = T;
		break;

	case TARRAY:
		len = 0;
		i = 0;
		for(ll=n->list; ll; ll=ll->next) {
			l = ll->n;
			if(l->op == OKEY) {
				typecheck(&l->left, Erv);
				evconst(l->left);
				i = nonnegconst(l->left);
				if(i < 0) {
					yyerror("array index must be non-negative integer constant");
					i = -(1<<30);	// stay negative for a while
				}
				typecheck(&l->right, Erv);
				defaultlit(&l->right, t->type);
				l->right = typecheckconv(nil, l->right, t->type, 0);
			} else {
				typecheck(&ll->n, Erv);
				defaultlit(&ll->n, t->type);
				ll->n = typecheckconv(nil, ll->n, t->type, 0);
				ll->n = nod(OKEY, nodintconst(i), ll->n);
				ll->n->left->type = types[TINT];
				ll->n->left->typecheck = 1;
			}
			if(i >= 0)
				indexdup(ll->n->left, hash, nelem(hash));
			i++;
			if(i > len) {
				len = i;
				if(t->bound >= 0 && len > t->bound) {
					setlineno(l);
					yyerror("array index out of bounds");
					t->bound = -1;	// no more errors
				}
			}
		}
		if(t->bound == -100)
			t->bound = len;
		if(t->bound < 0)
			n->right = nodintconst(len);
		n->op = OARRAYLIT;
		break;

	case TMAP:
		for(ll=n->list; ll; ll=ll->next) {
			l = ll->n;
			if(l->op != OKEY) {
				typecheck(&ll->n, Erv);
				yyerror("missing key in map literal");
				continue;
			}
			typecheck(&l->left, Erv);
			typecheck(&l->right, Erv);
			defaultlit(&l->left, t->down);
			defaultlit(&l->right, t->type);
			l->left = typecheckconv(nil, l->left, t->down, 0);
			l->right = typecheckconv(nil, l->right, t->type, 0);
			keydup(l->left, hash, nelem(hash));
		}
		n->op = OMAPLIT;
		break;

	case TSTRUCT:
		bad = 0;
		if(n->list != nil && nokeys(n->list)) {
			// simple list of variables
			f = t->type;
			for(ll=n->list; ll; ll=ll->next) {
				typecheck(&ll->n, Erv);
				if(f == nil) {
					if(!bad++)
						yyerror("too many values in struct initializer");
					continue;
				}
				ll->n = typecheckconv(nil, ll->n, f->type, 0);
				ll->n = nod(OKEY, newname(f->sym), ll->n);
				ll->n->left->typecheck = 1;
				f = f->down;
			}
		} else {
			// keyed list
			for(ll=n->list; ll; ll=ll->next) {
				l = ll->n;
				if(l->op != OKEY) {
					if(!bad++)
						yyerror("mixture of field:value and value initializers");
					typecheck(&ll->n, Erv);
					continue;
				}
				if(l->left->sym == S) {
					yyerror("invalid field name %#N in struct initializer", l->left);
					typecheck(&l->right, Erv);
					continue;
				}
				l->left->typecheck = 1;
				f = lookdot1(l->left->sym, t, t->type);
				typecheck(&l->right, Erv);
				if(f == nil)
					continue;
				fielddup(newname(f->sym), hash, nelem(hash));
				l->right = typecheckconv(nil, l->right, f->type, 0);
			}
		}
		n->op = OSTRUCTLIT;
		break;
	}
	if(nerr != nerrors)
		goto error;
	n->type = t;

	*np = n;
	return;

error:
	n->type = T;
	*np = n;
}

/*
 * the address of n has been taken and might be used after
 * the current function returns.  mark any local vars
 * as needing to move to the heap.
 */
static void
addrescapes(Node *n)
{
	char buf[100];
	switch(n->op) {
	default:
		// probably a type error already.
		// dump("addrescapes", n);
		break;

	case ONAME:
		if(n->noescape)
			break;
		switch(n->class) {
		case PPARAMOUT:
			yyerror("cannot take address of out parameter %s", n->sym->name);
			break;
		case PAUTO:
		case PPARAM:
			// if func param, need separate temporary
			// to hold heap pointer.
			if(n->class == PPARAM) {
				// expression to refer to stack copy
				n->stackparam = nod(OPARAM, n, N);
				n->stackparam->type = n->type;
				n->stackparam->addable = 1;
				n->stackparam->xoffset = n->xoffset;
			}

			n->class |= PHEAP;
			n->addable = 0;
			n->ullman = 2;
			n->alloc = callnew(n->type);
			n->xoffset = 0;

			// create stack variable to hold pointer to heap
			n->heapaddr = nod(0, N, N);
			tempname(n->heapaddr, ptrto(n->type));
			snprint(buf, sizeof buf, "&%S", n->sym);
			n->heapaddr->sym = lookup(buf);
			break;
		}
		break;

	case OIND:
	case ODOTPTR:
		break;

	case ODOT:
	case OINDEX:
		// ODOTPTR has already been introduced,
		// so these are the non-pointer ODOT and OINDEX.
		// In &x[0], if x is a slice, then x does not
		// escape--the pointer inside x does, but that
		// is always a heap pointer anyway.
		if(!isslice(n->left->type))
			addrescapes(n->left);
		break;
	}
}

static int
islvalue(Node *n)
{
	switch(n->op) {
	case OINDEX:
	case OIND:
	case ODOTPTR:
		return 1;
	case ODOT:
		return islvalue(n->left);
	case ONAME:
		if(n->class == PFUNC)
			return 0;
		return 1;
	}
	return 0;
}

static void
checklvalue(Node *n, char *verb)
{
	if(!islvalue(n))
		yyerror("cannot %s %#N", verb, n);
}

static void
checkassign(Node *n)
{
	if(islvalue(n))
		return;
	if(n->op == OINDEXMAP) {
		n->etype = 1;
		return;
	}
	yyerror("cannot assign to %#N", n);
}

static void
checkassignlist(NodeList *l)
{
	for(; l; l=l->next)
		checkassign(l->n);
}
