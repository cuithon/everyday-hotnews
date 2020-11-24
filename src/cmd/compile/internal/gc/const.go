// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gc

import (
	"cmd/compile/internal/types"
	"cmd/internal/src"
	"fmt"
	"go/constant"
	"math/big"
	"strings"
)

type Val struct {
	// U contains one of:
	// bool     bool when Ctype() == CTBOOL
	// *Mpint   int when Ctype() == CTINT
	// *Mpflt   float when Ctype() == CTFLT
	// *Mpcplx  pair of floats when Ctype() == CTCPLX
	// string   string when Ctype() == CTSTR
	U interface{}
}

func (v Val) Kind() constant.Kind {
	switch v.U.(type) {
	default:
		Fatalf("unexpected Ctype for %T", v.U)
		panic("unreachable")
	case nil:
		return constant.Unknown
	case bool:
		return constant.Bool
	case *Mpint:
		return constant.Int
	case *Mpflt:
		return constant.Float
	case *Mpcplx:
		return constant.Complex
	case string:
		return constant.String
	}
}

func eqval(a, b Val) bool {
	if a.Kind() != b.Kind() {
		return false
	}
	switch x := a.U.(type) {
	default:
		Fatalf("unexpected Ctype for %T", a.U)
		panic("unreachable")
	case bool:
		y := b.U.(bool)
		return x == y
	case *Mpint:
		y := b.U.(*Mpint)
		return x.Cmp(y) == 0
	case *Mpflt:
		y := b.U.(*Mpflt)
		return x.Cmp(y) == 0
	case *Mpcplx:
		y := b.U.(*Mpcplx)
		return x.Real.Cmp(&y.Real) == 0 && x.Imag.Cmp(&y.Imag) == 0
	case string:
		y := b.U.(string)
		return x == y
	}
}

// Interface returns the constant value stored in v as an interface{}.
// It returns int64s for ints and runes, float64s for floats,
// complex128s for complex values, and nil for constant nils.
func (v Val) Interface() interface{} {
	switch x := v.U.(type) {
	default:
		Fatalf("unexpected Interface for %T", v.U)
		panic("unreachable")
	case bool, string:
		return x
	case *Mpint:
		return x.Int64()
	case *Mpflt:
		return x.Float64()
	case *Mpcplx:
		return complex(x.Real.Float64(), x.Imag.Float64())
	}
}

// Int64Val returns n as an int64.
// n must be an integer or rune constant.
func (n *Node) Int64Val() int64 {
	if !Isconst(n, constant.Int) {
		Fatalf("Int64Val(%v)", n)
	}
	return n.Val().U.(*Mpint).Int64()
}

// CanInt64 reports whether it is safe to call Int64Val() on n.
func (n *Node) CanInt64() bool {
	if !Isconst(n, constant.Int) {
		return false
	}

	// if the value inside n cannot be represented as an int64, the
	// return value of Int64 is undefined
	return n.Val().U.(*Mpint).CmpInt64(n.Int64Val()) == 0
}

// BoolVal returns n as a bool.
// n must be a boolean constant.
func (n *Node) BoolVal() bool {
	if !Isconst(n, constant.Bool) {
		Fatalf("BoolVal(%v)", n)
	}
	return n.Val().U.(bool)
}

// StringVal returns the value of a literal string Node as a string.
// n must be a string constant.
func (n *Node) StringVal() string {
	if !Isconst(n, constant.String) {
		Fatalf("StringVal(%v)", n)
	}
	return n.Val().U.(string)
}

// truncate float literal fv to 32-bit or 64-bit precision
// according to type; return truncated value.
func truncfltlit(oldv *Mpflt, t *types.Type) *Mpflt {
	if t == nil {
		return oldv
	}

	if overflow(Val{oldv}, t) {
		// If there was overflow, simply continuing would set the
		// value to Inf which in turn would lead to spurious follow-on
		// errors. Avoid this by returning the existing value.
		return oldv
	}

	fv := newMpflt()

	// convert large precision literal floating
	// into limited precision (float64 or float32)
	switch t.Etype {
	case types.TFLOAT32:
		fv.SetFloat64(oldv.Float32())
	case types.TFLOAT64:
		fv.SetFloat64(oldv.Float64())
	default:
		Fatalf("truncfltlit: unexpected Etype %v", t.Etype)
	}

	return fv
}

// truncate Real and Imag parts of Mpcplx to 32-bit or 64-bit
// precision, according to type; return truncated value. In case of
// overflow, calls yyerror but does not truncate the input value.
func trunccmplxlit(oldv *Mpcplx, t *types.Type) *Mpcplx {
	if t == nil {
		return oldv
	}

	if overflow(Val{oldv}, t) {
		// If there was overflow, simply continuing would set the
		// value to Inf which in turn would lead to spurious follow-on
		// errors. Avoid this by returning the existing value.
		return oldv
	}

	cv := newMpcmplx()

	switch t.Etype {
	case types.TCOMPLEX64:
		cv.Real.SetFloat64(oldv.Real.Float32())
		cv.Imag.SetFloat64(oldv.Imag.Float32())
	case types.TCOMPLEX128:
		cv.Real.SetFloat64(oldv.Real.Float64())
		cv.Imag.SetFloat64(oldv.Imag.Float64())
	default:
		Fatalf("trunccplxlit: unexpected Etype %v", t.Etype)
	}

	return cv
}

// TODO(mdempsky): Replace these with better APIs.
func convlit(n *Node, t *types.Type) *Node    { return convlit1(n, t, false, nil) }
func defaultlit(n *Node, t *types.Type) *Node { return convlit1(n, t, false, nil) }

// convlit1 converts an untyped expression n to type t. If n already
// has a type, convlit1 has no effect.
//
// For explicit conversions, t must be non-nil, and integer-to-string
// conversions are allowed.
//
// For implicit conversions (e.g., assignments), t may be nil; if so,
// n is converted to its default type.
//
// If there's an error converting n to t, context is used in the error
// message.
func convlit1(n *Node, t *types.Type, explicit bool, context func() string) *Node {
	if explicit && t == nil {
		Fatalf("explicit conversion missing type")
	}
	if t != nil && t.IsUntyped() {
		Fatalf("bad conversion to untyped: %v", t)
	}

	if n == nil || n.Type == nil {
		// Allow sloppy callers.
		return n
	}
	if !n.Type.IsUntyped() {
		// Already typed; nothing to do.
		return n
	}

	if n.Op == OLITERAL || n.Op == ONIL {
		// Can't always set n.Type directly on OLITERAL nodes.
		// See discussion on CL 20813.
		n = n.rawcopy()
	}

	// Nil is technically not a constant, so handle it specially.
	if n.Type.Etype == TNIL {
		if n.Op != ONIL {
			Fatalf("unexpected op: %v (%v)", n, n.Op)
		}
		if t == nil {
			yyerror("use of untyped nil")
			n.SetDiag(true)
			n.Type = nil
			return n
		}

		if !t.HasNil() {
			// Leave for caller to handle.
			return n
		}

		n.Type = t
		return n
	}

	if t == nil || !okforconst[t.Etype] {
		t = defaultType(n.Type)
	}

	switch n.Op {
	default:
		Fatalf("unexpected untyped expression: %v", n)

	case OLITERAL:
		v := convertVal(n.Val(), t, explicit)
		if v.U == nil {
			break
		}
		n.Type = t
		n.SetVal(v)
		return n

	case OPLUS, ONEG, OBITNOT, ONOT, OREAL, OIMAG:
		ot := operandType(n.Op, t)
		if ot == nil {
			n = defaultlit(n, nil)
			break
		}

		n.Left = convlit(n.Left, ot)
		if n.Left.Type == nil {
			n.Type = nil
			return n
		}
		n.Type = t
		return n

	case OADD, OSUB, OMUL, ODIV, OMOD, OOR, OXOR, OAND, OANDNOT, OOROR, OANDAND, OCOMPLEX:
		ot := operandType(n.Op, t)
		if ot == nil {
			n = defaultlit(n, nil)
			break
		}

		n.Left = convlit(n.Left, ot)
		n.Right = convlit(n.Right, ot)
		if n.Left.Type == nil || n.Right.Type == nil {
			n.Type = nil
			return n
		}
		if !types.Identical(n.Left.Type, n.Right.Type) {
			yyerror("invalid operation: %v (mismatched types %v and %v)", n, n.Left.Type, n.Right.Type)
			n.Type = nil
			return n
		}

		n.Type = t
		return n

	case OEQ, ONE, OLT, OLE, OGT, OGE:
		if !t.IsBoolean() {
			break
		}
		n.Type = t
		return n

	case OLSH, ORSH:
		n.Left = convlit1(n.Left, t, explicit, nil)
		n.Type = n.Left.Type
		if n.Type != nil && !n.Type.IsInteger() {
			yyerror("invalid operation: %v (shift of type %v)", n, n.Type)
			n.Type = nil
		}
		return n
	}

	if !n.Diag() {
		if !t.Broke() {
			if explicit {
				yyerror("cannot convert %L to type %v", n, t)
			} else if context != nil {
				yyerror("cannot use %L as type %v in %s", n, t, context())
			} else {
				yyerror("cannot use %L as type %v", n, t)
			}
		}
		n.SetDiag(true)
	}
	n.Type = nil
	return n
}

func operandType(op Op, t *types.Type) *types.Type {
	switch op {
	case OCOMPLEX:
		if t.IsComplex() {
			return floatForComplex(t)
		}
	case OREAL, OIMAG:
		if t.IsFloat() {
			return complexForFloat(t)
		}
	default:
		if okfor[op][t.Etype] {
			return t
		}
	}
	return nil
}

// convertVal converts v into a representation appropriate for t. If
// no such representation exists, it returns Val{} instead.
//
// If explicit is true, then conversions from integer to string are
// also allowed.
func convertVal(v Val, t *types.Type, explicit bool) Val {
	switch ct := v.Kind(); ct {
	case constant.Bool:
		if t.IsBoolean() {
			return v
		}

	case constant.String:
		if t.IsString() {
			return v
		}

	case constant.Int:
		if explicit && t.IsString() {
			return tostr(v)
		}
		fallthrough
	case constant.Float, constant.Complex:
		switch {
		case t.IsInteger():
			v = toint(v)
			overflow(v, t)
			return v
		case t.IsFloat():
			v = toflt(v)
			v = Val{truncfltlit(v.U.(*Mpflt), t)}
			return v
		case t.IsComplex():
			v = tocplx(v)
			v = Val{trunccmplxlit(v.U.(*Mpcplx), t)}
			return v
		}
	}

	return Val{}
}

func tocplx(v Val) Val {
	switch u := v.U.(type) {
	case *Mpint:
		c := newMpcmplx()
		c.Real.SetInt(u)
		c.Imag.SetFloat64(0.0)
		v.U = c

	case *Mpflt:
		c := newMpcmplx()
		c.Real.Set(u)
		c.Imag.SetFloat64(0.0)
		v.U = c
	}

	return v
}

func toflt(v Val) Val {
	switch u := v.U.(type) {
	case *Mpint:
		f := newMpflt()
		f.SetInt(u)
		v.U = f

	case *Mpcplx:
		f := newMpflt()
		f.Set(&u.Real)
		if u.Imag.CmpFloat64(0) != 0 {
			yyerror("constant %v truncated to real", u.GoString())
		}
		v.U = f
	}

	return v
}

func toint(v Val) Val {
	switch u := v.U.(type) {
	case *Mpint:

	case *Mpflt:
		i := new(Mpint)
		if !i.SetFloat(u) {
			if i.checkOverflow(0) {
				yyerror("integer too large")
			} else {
				// The value of u cannot be represented as an integer;
				// so we need to print an error message.
				// Unfortunately some float values cannot be
				// reasonably formatted for inclusion in an error
				// message (example: 1 + 1e-100), so first we try to
				// format the float; if the truncation resulted in
				// something that looks like an integer we omit the
				// value from the error message.
				// (See issue #11371).
				var t big.Float
				t.Parse(u.GoString(), 10)
				if t.IsInt() {
					yyerror("constant truncated to integer")
				} else {
					yyerror("constant %v truncated to integer", u.GoString())
				}
			}
		}
		v.U = i

	case *Mpcplx:
		i := new(Mpint)
		if !i.SetFloat(&u.Real) || u.Imag.CmpFloat64(0) != 0 {
			yyerror("constant %v truncated to integer", u.GoString())
		}

		v.U = i
	}

	return v
}

func doesoverflow(v Val, t *types.Type) bool {
	switch u := v.U.(type) {
	case *Mpint:
		if !t.IsInteger() {
			Fatalf("overflow: %v integer constant", t)
		}
		return u.Cmp(minintval[t.Etype]) < 0 || u.Cmp(maxintval[t.Etype]) > 0

	case *Mpflt:
		if !t.IsFloat() {
			Fatalf("overflow: %v floating-point constant", t)
		}
		return u.Cmp(minfltval[t.Etype]) <= 0 || u.Cmp(maxfltval[t.Etype]) >= 0

	case *Mpcplx:
		if !t.IsComplex() {
			Fatalf("overflow: %v complex constant", t)
		}
		return u.Real.Cmp(minfltval[t.Etype]) <= 0 || u.Real.Cmp(maxfltval[t.Etype]) >= 0 ||
			u.Imag.Cmp(minfltval[t.Etype]) <= 0 || u.Imag.Cmp(maxfltval[t.Etype]) >= 0
	}

	return false
}

func overflow(v Val, t *types.Type) bool {
	// v has already been converted
	// to appropriate form for t.
	if t == nil || t.Etype == TIDEAL {
		return false
	}

	// Only uintptrs may be converted to pointers, which cannot overflow.
	if t.IsPtr() || t.IsUnsafePtr() {
		return false
	}

	if doesoverflow(v, t) {
		yyerror("constant %v overflows %v", v, t)
		return true
	}

	return false

}

func tostr(v Val) Val {
	switch u := v.U.(type) {
	case *Mpint:
		var r rune = 0xFFFD
		if u.Cmp(minintval[TINT32]) >= 0 && u.Cmp(maxintval[TINT32]) <= 0 {
			r = rune(u.Int64())
		}
		v.U = string(r)
	}

	return v
}

func consttype(n *Node) constant.Kind {
	if n == nil || n.Op != OLITERAL {
		return constant.Unknown
	}
	return n.Val().Kind()
}

func Isconst(n *Node, ct constant.Kind) bool {
	return consttype(n) == ct
}

// evalConst returns a constant-evaluated expression equivalent to n.
// If n is not a constant, evalConst returns n.
// Otherwise, evalConst returns a new OLITERAL with the same value as n,
// and with .Orig pointing back to n.
func evalConst(n *Node) *Node {
	nl, nr := n.Left, n.Right

	// Pick off just the opcodes that can be constant evaluated.
	switch op := n.Op; op {
	case OPLUS, ONEG, OBITNOT, ONOT:
		if nl.Op == OLITERAL {
			return origConst(n, unaryOp(op, nl.Val(), n.Type))
		}

	case OADD, OSUB, OMUL, ODIV, OMOD, OOR, OXOR, OAND, OANDNOT, OOROR, OANDAND:
		if nl.Op == OLITERAL && nr.Op == OLITERAL {
			return origConst(n, binaryOp(nl.Val(), op, nr.Val()))
		}

	case OEQ, ONE, OLT, OLE, OGT, OGE:
		if nl.Op == OLITERAL && nr.Op == OLITERAL {
			return origBoolConst(n, compareOp(nl.Val(), op, nr.Val()))
		}

	case OLSH, ORSH:
		if nl.Op == OLITERAL && nr.Op == OLITERAL {
			return origConst(n, shiftOp(nl.Val(), op, nr.Val()))
		}

	case OCONV, ORUNESTR:
		if okforconst[n.Type.Etype] && nl.Op == OLITERAL {
			return origConst(n, convertVal(nl.Val(), n.Type, true))
		}

	case OCONVNOP:
		if okforconst[n.Type.Etype] && nl.Op == OLITERAL {
			// set so n.Orig gets OCONV instead of OCONVNOP
			n.Op = OCONV
			return origConst(n, nl.Val())
		}

	case OADDSTR:
		// Merge adjacent constants in the argument list.
		s := n.List.Slice()
		need := 0
		for i := 0; i < len(s); i++ {
			if i == 0 || !Isconst(s[i-1], constant.String) || !Isconst(s[i], constant.String) {
				// Can't merge s[i] into s[i-1]; need a slot in the list.
				need++
			}
		}
		if need == len(s) {
			return n
		}
		if need == 1 {
			var strs []string
			for _, c := range s {
				strs = append(strs, c.StringVal())
			}
			return origConst(n, Val{U: strings.Join(strs, "")})
		}
		newList := make([]*Node, 0, need)
		for i := 0; i < len(s); i++ {
			if Isconst(s[i], constant.String) && i+1 < len(s) && Isconst(s[i+1], constant.String) {
				// merge from i up to but not including i2
				var strs []string
				i2 := i
				for i2 < len(s) && Isconst(s[i2], constant.String) {
					strs = append(strs, s[i2].StringVal())
					i2++
				}

				nl := origConst(s[i], Val{U: strings.Join(strs, "")})
				nl.Orig = nl // it's bigger than just s[i]
				newList = append(newList, nl)
				i = i2 - 1
			} else {
				newList = append(newList, s[i])
			}
		}

		n = n.copy()
		n.List.Set(newList)
		return n

	case OCAP, OLEN:
		switch nl.Type.Etype {
		case TSTRING:
			if Isconst(nl, constant.String) {
				return origIntConst(n, int64(len(nl.StringVal())))
			}
		case TARRAY:
			if !hascallchan(nl) {
				return origIntConst(n, nl.Type.NumElem())
			}
		}

	case OALIGNOF, OOFFSETOF, OSIZEOF:
		return origIntConst(n, evalunsafe(n))

	case OREAL, OIMAG:
		if nl.Op == OLITERAL {
			var re, im *Mpflt
			switch u := nl.Val().U.(type) {
			case *Mpint:
				re = newMpflt()
				re.SetInt(u)
				// im = 0
			case *Mpflt:
				re = u
				// im = 0
			case *Mpcplx:
				re = &u.Real
				im = &u.Imag
			default:
				Fatalf("impossible")
			}
			if n.Op == OIMAG {
				if im == nil {
					im = newMpflt()
				}
				re = im
			}
			return origConst(n, Val{re})
		}

	case OCOMPLEX:
		if nl.Op == OLITERAL && nr.Op == OLITERAL {
			// make it a complex literal
			c := newMpcmplx()
			c.Real.Set(toflt(nl.Val()).U.(*Mpflt))
			c.Imag.Set(toflt(nr.Val()).U.(*Mpflt))
			return origConst(n, Val{c})
		}
	}

	return n
}

func match(x, y Val) (Val, Val) {
	switch {
	case x.Kind() == constant.Complex || y.Kind() == constant.Complex:
		return tocplx(x), tocplx(y)
	case x.Kind() == constant.Float || y.Kind() == constant.Float:
		return toflt(x), toflt(y)
	}

	// Mixed int/rune are fine.
	return x, y
}

func compareOp(x Val, op Op, y Val) bool {
	x, y = match(x, y)

	switch x.Kind() {
	case constant.Bool:
		x, y := x.U.(bool), y.U.(bool)
		switch op {
		case OEQ:
			return x == y
		case ONE:
			return x != y
		}

	case constant.Int:
		x, y := x.U.(*Mpint), y.U.(*Mpint)
		return cmpZero(x.Cmp(y), op)

	case constant.Float:
		x, y := x.U.(*Mpflt), y.U.(*Mpflt)
		return cmpZero(x.Cmp(y), op)

	case constant.Complex:
		x, y := x.U.(*Mpcplx), y.U.(*Mpcplx)
		eq := x.Real.Cmp(&y.Real) == 0 && x.Imag.Cmp(&y.Imag) == 0
		switch op {
		case OEQ:
			return eq
		case ONE:
			return !eq
		}

	case constant.String:
		x, y := x.U.(string), y.U.(string)
		switch op {
		case OEQ:
			return x == y
		case ONE:
			return x != y
		case OLT:
			return x < y
		case OLE:
			return x <= y
		case OGT:
			return x > y
		case OGE:
			return x >= y
		}
	}

	Fatalf("compareOp: bad comparison: %v %v %v", x, op, y)
	panic("unreachable")
}

func cmpZero(x int, op Op) bool {
	switch op {
	case OEQ:
		return x == 0
	case ONE:
		return x != 0
	case OLT:
		return x < 0
	case OLE:
		return x <= 0
	case OGT:
		return x > 0
	case OGE:
		return x >= 0
	}

	Fatalf("cmpZero: want comparison operator, got %v", op)
	panic("unreachable")
}

func binaryOp(x Val, op Op, y Val) Val {
	x, y = match(x, y)

Outer:
	switch x.Kind() {
	case constant.Bool:
		x, y := x.U.(bool), y.U.(bool)
		switch op {
		case OANDAND:
			return Val{U: x && y}
		case OOROR:
			return Val{U: x || y}
		}

	case constant.Int:
		x, y := x.U.(*Mpint), y.U.(*Mpint)

		u := new(Mpint)
		u.Set(x)
		switch op {
		case OADD:
			u.Add(y)
		case OSUB:
			u.Sub(y)
		case OMUL:
			u.Mul(y)
		case ODIV:
			if y.CmpInt64(0) == 0 {
				yyerror("division by zero")
				return Val{}
			}
			u.Quo(y)
		case OMOD:
			if y.CmpInt64(0) == 0 {
				yyerror("division by zero")
				return Val{}
			}
			u.Rem(y)
		case OOR:
			u.Or(y)
		case OAND:
			u.And(y)
		case OANDNOT:
			u.AndNot(y)
		case OXOR:
			u.Xor(y)
		default:
			break Outer
		}
		return Val{U: u}

	case constant.Float:
		x, y := x.U.(*Mpflt), y.U.(*Mpflt)

		u := newMpflt()
		u.Set(x)
		switch op {
		case OADD:
			u.Add(y)
		case OSUB:
			u.Sub(y)
		case OMUL:
			u.Mul(y)
		case ODIV:
			if y.CmpFloat64(0) == 0 {
				yyerror("division by zero")
				return Val{}
			}
			u.Quo(y)
		default:
			break Outer
		}
		return Val{U: u}

	case constant.Complex:
		x, y := x.U.(*Mpcplx), y.U.(*Mpcplx)

		u := newMpcmplx()
		u.Real.Set(&x.Real)
		u.Imag.Set(&x.Imag)
		switch op {
		case OADD:
			u.Real.Add(&y.Real)
			u.Imag.Add(&y.Imag)
		case OSUB:
			u.Real.Sub(&y.Real)
			u.Imag.Sub(&y.Imag)
		case OMUL:
			u.Mul(y)
		case ODIV:
			if !u.Div(y) {
				yyerror("complex division by zero")
				return Val{}
			}
		default:
			break Outer
		}
		return Val{U: u}
	}

	Fatalf("binaryOp: bad operation: %v %v %v", x, op, y)
	panic("unreachable")
}

func unaryOp(op Op, x Val, t *types.Type) Val {
	switch op {
	case OPLUS:
		switch x.Kind() {
		case constant.Int, constant.Float, constant.Complex:
			return x
		}

	case ONEG:
		switch x.Kind() {
		case constant.Int:
			x := x.U.(*Mpint)
			u := new(Mpint)
			u.Set(x)
			u.Neg()
			return Val{U: u}

		case constant.Float:
			x := x.U.(*Mpflt)
			u := newMpflt()
			u.Set(x)
			u.Neg()
			return Val{U: u}

		case constant.Complex:
			x := x.U.(*Mpcplx)
			u := newMpcmplx()
			u.Real.Set(&x.Real)
			u.Imag.Set(&x.Imag)
			u.Real.Neg()
			u.Imag.Neg()
			return Val{U: u}
		}

	case OBITNOT:
		switch x.Kind() {
		case constant.Int:
			x := x.U.(*Mpint)

			u := new(Mpint)
			if t.IsSigned() || t.IsUntyped() {
				// Signed values change sign.
				u.SetInt64(-1)
			} else {
				// Unsigned values invert their bits.
				u.Set(maxintval[t.Etype])
			}
			u.Xor(x)
			return Val{U: u}
		}

	case ONOT:
		return Val{U: !x.U.(bool)}
	}

	Fatalf("unaryOp: bad operation: %v %v", op, x)
	panic("unreachable")
}

func shiftOp(x Val, op Op, y Val) Val {
	x = toint(x)
	y = toint(y)

	u := new(Mpint)
	u.Set(x.U.(*Mpint))
	switch op {
	case OLSH:
		u.Lsh(y.U.(*Mpint))
	case ORSH:
		u.Rsh(y.U.(*Mpint))
	default:
		Fatalf("shiftOp: bad operator: %v", op)
		panic("unreachable")
	}
	return Val{U: u}
}

// origConst returns an OLITERAL with orig n and value v.
func origConst(n *Node, v Val) *Node {
	// If constant folding was attempted (we were called)
	// but it produced an invalid constant value,
	// mark n as broken and give up.
	if v.U == nil {
		n.Type = nil
		return n
	}

	orig := n
	n = nod(OLITERAL, nil, nil)
	n.Orig = orig
	n.Pos = orig.Pos
	n.Type = orig.Type
	n.SetVal(v)

	// Check range.
	lno := setlineno(n)
	overflow(v, n.Type)
	lineno = lno

	if !n.Type.IsUntyped() {
		switch v.Kind() {
		// Truncate precision for non-ideal float.
		case constant.Float:
			n.SetVal(Val{truncfltlit(v.U.(*Mpflt), n.Type)})
		// Truncate precision for non-ideal complex.
		case constant.Complex:
			n.SetVal(Val{trunccmplxlit(v.U.(*Mpcplx), n.Type)})
		}
	}
	return n
}

func assertRepresents(t *types.Type, v Val) {
	if !represents(t, v) {
		Fatalf("%v does not represent %v", t, v)
	}
}

func represents(t *types.Type, v Val) bool {
	if !t.IsUntyped() {
		// TODO(mdempsky): Stricter handling of typed types.
		return true
	}

	vt := idealType(v.Kind())
	return t == vt || (t == types.UntypedRune && vt == types.UntypedInt)
}

func origBoolConst(n *Node, v bool) *Node {
	return origConst(n, Val{U: v})
}

func origIntConst(n *Node, v int64) *Node {
	u := new(Mpint)
	u.SetInt64(v)
	return origConst(n, Val{u})
}

// nodlit returns a new untyped constant with value v.
func nodlit(v Val) *Node {
	n := nod(OLITERAL, nil, nil)
	n.Type = idealType(v.Kind())
	n.SetVal(v)
	return n
}

func idealType(ct constant.Kind) *types.Type {
	switch ct {
	case constant.String:
		return types.UntypedString
	case constant.Bool:
		return types.UntypedBool
	case constant.Int:
		return types.UntypedInt
	case constant.Float:
		return types.UntypedFloat
	case constant.Complex:
		return types.UntypedComplex
	}
	Fatalf("unexpected Ctype: %v", ct)
	return nil
}

// defaultlit on both nodes simultaneously;
// if they're both ideal going in they better
// get the same type going out.
// force means must assign concrete (non-ideal) type.
// The results of defaultlit2 MUST be assigned back to l and r, e.g.
// 	n.Left, n.Right = defaultlit2(n.Left, n.Right, force)
func defaultlit2(l *Node, r *Node, force bool) (*Node, *Node) {
	if l.Type == nil || r.Type == nil {
		return l, r
	}
	if !l.Type.IsUntyped() {
		r = convlit(r, l.Type)
		return l, r
	}

	if !r.Type.IsUntyped() {
		l = convlit(l, r.Type)
		return l, r
	}

	if !force {
		return l, r
	}

	// Can't mix bool with non-bool, string with non-string, or nil with anything (untyped).
	if l.Type.IsBoolean() != r.Type.IsBoolean() {
		return l, r
	}
	if l.Type.IsString() != r.Type.IsString() {
		return l, r
	}
	if l.isNil() || r.isNil() {
		return l, r
	}

	t := defaultType(mixUntyped(l.Type, r.Type))
	l = convlit(l, t)
	r = convlit(r, t)
	return l, r
}

func mixUntyped(t1, t2 *types.Type) *types.Type {
	if t1 == t2 {
		return t1
	}

	rank := func(t *types.Type) int {
		switch t {
		case types.UntypedInt:
			return 0
		case types.UntypedRune:
			return 1
		case types.UntypedFloat:
			return 2
		case types.UntypedComplex:
			return 3
		}
		Fatalf("bad type %v", t)
		panic("unreachable")
	}

	if rank(t2) > rank(t1) {
		return t2
	}
	return t1
}

func defaultType(t *types.Type) *types.Type {
	if !t.IsUntyped() || t.Etype == TNIL {
		return t
	}

	switch t {
	case types.UntypedBool:
		return types.Types[TBOOL]
	case types.UntypedString:
		return types.Types[TSTRING]
	case types.UntypedInt:
		return types.Types[TINT]
	case types.UntypedRune:
		return types.Runetype
	case types.UntypedFloat:
		return types.Types[TFLOAT64]
	case types.UntypedComplex:
		return types.Types[TCOMPLEX128]
	}

	Fatalf("bad type %v", t)
	return nil
}

func smallintconst(n *Node) bool {
	if n.Op == OLITERAL && Isconst(n, constant.Int) && n.Type != nil {
		switch simtype[n.Type.Etype] {
		case TINT8,
			TUINT8,
			TINT16,
			TUINT16,
			TINT32,
			TUINT32,
			TBOOL:
			return true

		case TIDEAL, TINT64, TUINT64, TPTR:
			v, ok := n.Val().U.(*Mpint)
			if ok && v.Cmp(minintval[TINT32]) >= 0 && v.Cmp(maxintval[TINT32]) <= 0 {
				return true
			}
		}
	}

	return false
}

// indexconst checks if Node n contains a constant expression
// representable as a non-negative int and returns its value.
// If n is not a constant expression, not representable as an
// integer, or negative, it returns -1. If n is too large, it
// returns -2.
func indexconst(n *Node) int64 {
	if n.Op != OLITERAL {
		return -1
	}

	v := toint(n.Val()) // toint returns argument unchanged if not representable as an *Mpint
	vi, ok := v.U.(*Mpint)
	if !ok || vi.CmpInt64(0) < 0 {
		return -1
	}
	if vi.Cmp(maxintval[TINT]) > 0 {
		return -2
	}

	return vi.Int64()
}

// isGoConst reports whether n is a Go language constant (as opposed to a
// compile-time constant).
//
// Expressions derived from nil, like string([]byte(nil)), while they
// may be known at compile time, are not Go language constants.
func (n *Node) isGoConst() bool {
	return n.Op == OLITERAL
}

func hascallchan(n *Node) bool {
	if n == nil {
		return false
	}
	switch n.Op {
	case OAPPEND,
		OCALL,
		OCALLFUNC,
		OCALLINTER,
		OCALLMETH,
		OCAP,
		OCLOSE,
		OCOMPLEX,
		OCOPY,
		ODELETE,
		OIMAG,
		OLEN,
		OMAKE,
		ONEW,
		OPANIC,
		OPRINT,
		OPRINTN,
		OREAL,
		ORECOVER,
		ORECV:
		return true
	}

	if hascallchan(n.Left) || hascallchan(n.Right) {
		return true
	}
	for _, n1 := range n.List.Slice() {
		if hascallchan(n1) {
			return true
		}
	}
	for _, n2 := range n.Rlist.Slice() {
		if hascallchan(n2) {
			return true
		}
	}

	return false
}

// A constSet represents a set of Go constant expressions.
type constSet struct {
	m map[constSetKey]src.XPos
}

type constSetKey struct {
	typ *types.Type
	val interface{}
}

// add adds constant expression n to s. If a constant expression of
// equal value and identical type has already been added, then add
// reports an error about the duplicate value.
//
// pos provides position information for where expression n occurred
// (in case n does not have its own position information). what and
// where are used in the error message.
//
// n must not be an untyped constant.
func (s *constSet) add(pos src.XPos, n *Node, what, where string) {
	if n.Op == OCONVIFACE && n.Implicit() {
		n = n.Left
	}

	if !n.isGoConst() {
		return
	}
	if n.Type.IsUntyped() {
		Fatalf("%v is untyped", n)
	}

	// Consts are only duplicates if they have the same value and
	// identical types.
	//
	// In general, we have to use types.Identical to test type
	// identity, because == gives false negatives for anonymous
	// types and the byte/uint8 and rune/int32 builtin type
	// aliases.  However, this is not a problem here, because
	// constant expressions are always untyped or have a named
	// type, and we explicitly handle the builtin type aliases
	// below.
	//
	// This approach may need to be revisited though if we fix
	// #21866 by treating all type aliases like byte/uint8 and
	// rune/int32.

	typ := n.Type
	switch typ {
	case types.Bytetype:
		typ = types.Types[TUINT8]
	case types.Runetype:
		typ = types.Types[TINT32]
	}
	k := constSetKey{typ, n.Val().Interface()}

	if hasUniquePos(n) {
		pos = n.Pos
	}

	if s.m == nil {
		s.m = make(map[constSetKey]src.XPos)
	}

	if prevPos, isDup := s.m[k]; isDup {
		yyerrorl(pos, "duplicate %s %s in %s\n\tprevious %s at %v",
			what, nodeAndVal(n), where,
			what, linestr(prevPos))
	} else {
		s.m[k] = pos
	}
}

// nodeAndVal reports both an expression and its constant value, if
// the latter is non-obvious.
//
// TODO(mdempsky): This could probably be a fmt.go flag.
func nodeAndVal(n *Node) string {
	show := n.String()
	val := n.Val().Interface()
	if s := fmt.Sprintf("%#v", val); show != s {
		show += " (value " + s + ")"
	}
	return show
}
