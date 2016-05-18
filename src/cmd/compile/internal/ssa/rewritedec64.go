// autogenerated from gen/dec64.rules: do not edit!
// generated with: cd gen; go run *.go

package ssa

import "math"

var _ = math.MinInt8 // in case not otherwise used
func rewriteValuedec64(v *Value, config *Config) bool {
	switch v.Op {
	case OpAdd64:
		return rewriteValuedec64_OpAdd64(v, config)
	case OpAnd64:
		return rewriteValuedec64_OpAnd64(v, config)
	case OpArg:
		return rewriteValuedec64_OpArg(v, config)
	case OpCom64:
		return rewriteValuedec64_OpCom64(v, config)
	case OpConst64:
		return rewriteValuedec64_OpConst64(v, config)
	case OpEq64:
		return rewriteValuedec64_OpEq64(v, config)
	case OpGeq64:
		return rewriteValuedec64_OpGeq64(v, config)
	case OpGeq64U:
		return rewriteValuedec64_OpGeq64U(v, config)
	case OpGreater64:
		return rewriteValuedec64_OpGreater64(v, config)
	case OpGreater64U:
		return rewriteValuedec64_OpGreater64U(v, config)
	case OpInt64Hi:
		return rewriteValuedec64_OpInt64Hi(v, config)
	case OpInt64Lo:
		return rewriteValuedec64_OpInt64Lo(v, config)
	case OpLeq64:
		return rewriteValuedec64_OpLeq64(v, config)
	case OpLeq64U:
		return rewriteValuedec64_OpLeq64U(v, config)
	case OpLess64:
		return rewriteValuedec64_OpLess64(v, config)
	case OpLess64U:
		return rewriteValuedec64_OpLess64U(v, config)
	case OpLoad:
		return rewriteValuedec64_OpLoad(v, config)
	case OpLsh16x64:
		return rewriteValuedec64_OpLsh16x64(v, config)
	case OpLsh32x64:
		return rewriteValuedec64_OpLsh32x64(v, config)
	case OpLsh8x64:
		return rewriteValuedec64_OpLsh8x64(v, config)
	case OpMul64:
		return rewriteValuedec64_OpMul64(v, config)
	case OpNeg64:
		return rewriteValuedec64_OpNeg64(v, config)
	case OpNeq64:
		return rewriteValuedec64_OpNeq64(v, config)
	case OpOr64:
		return rewriteValuedec64_OpOr64(v, config)
	case OpRsh16Ux64:
		return rewriteValuedec64_OpRsh16Ux64(v, config)
	case OpRsh16x64:
		return rewriteValuedec64_OpRsh16x64(v, config)
	case OpRsh32Ux64:
		return rewriteValuedec64_OpRsh32Ux64(v, config)
	case OpRsh32x64:
		return rewriteValuedec64_OpRsh32x64(v, config)
	case OpRsh8Ux64:
		return rewriteValuedec64_OpRsh8Ux64(v, config)
	case OpRsh8x64:
		return rewriteValuedec64_OpRsh8x64(v, config)
	case OpSignExt16to64:
		return rewriteValuedec64_OpSignExt16to64(v, config)
	case OpSignExt32to64:
		return rewriteValuedec64_OpSignExt32to64(v, config)
	case OpSignExt8to64:
		return rewriteValuedec64_OpSignExt8to64(v, config)
	case OpStore:
		return rewriteValuedec64_OpStore(v, config)
	case OpSub64:
		return rewriteValuedec64_OpSub64(v, config)
	case OpTrunc64to16:
		return rewriteValuedec64_OpTrunc64to16(v, config)
	case OpTrunc64to32:
		return rewriteValuedec64_OpTrunc64to32(v, config)
	case OpTrunc64to8:
		return rewriteValuedec64_OpTrunc64to8(v, config)
	case OpXor64:
		return rewriteValuedec64_OpXor64(v, config)
	case OpZeroExt16to64:
		return rewriteValuedec64_OpZeroExt16to64(v, config)
	case OpZeroExt32to64:
		return rewriteValuedec64_OpZeroExt32to64(v, config)
	case OpZeroExt8to64:
		return rewriteValuedec64_OpZeroExt8to64(v, config)
	}
	return false
}
func rewriteValuedec64_OpAdd64(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (Add64 x y)
	// cond:
	// result: (Int64Make 		(Add32withcarry <config.fe.TypeInt32()> 			(Int64Hi x) 			(Int64Hi y) 			(Select0 <TypeFlags> (Add32carry (Int64Lo x) (Int64Lo y)))) 		(Select1 <config.fe.TypeUInt32()> (Add32carry (Int64Lo x) (Int64Lo y))))
	for {
		x := v.Args[0]
		y := v.Args[1]
		v.reset(OpInt64Make)
		v0 := b.NewValue0(v.Line, OpAdd32withcarry, config.fe.TypeInt32())
		v1 := b.NewValue0(v.Line, OpInt64Hi, config.fe.TypeUInt32())
		v1.AddArg(x)
		v0.AddArg(v1)
		v2 := b.NewValue0(v.Line, OpInt64Hi, config.fe.TypeUInt32())
		v2.AddArg(y)
		v0.AddArg(v2)
		v3 := b.NewValue0(v.Line, OpSelect0, TypeFlags)
		v4 := b.NewValue0(v.Line, OpAdd32carry, MakeTuple(TypeFlags, config.fe.TypeUInt32()))
		v5 := b.NewValue0(v.Line, OpInt64Lo, config.fe.TypeUInt32())
		v5.AddArg(x)
		v4.AddArg(v5)
		v6 := b.NewValue0(v.Line, OpInt64Lo, config.fe.TypeUInt32())
		v6.AddArg(y)
		v4.AddArg(v6)
		v3.AddArg(v4)
		v0.AddArg(v3)
		v.AddArg(v0)
		v7 := b.NewValue0(v.Line, OpSelect1, config.fe.TypeUInt32())
		v8 := b.NewValue0(v.Line, OpAdd32carry, MakeTuple(TypeFlags, config.fe.TypeUInt32()))
		v9 := b.NewValue0(v.Line, OpInt64Lo, config.fe.TypeUInt32())
		v9.AddArg(x)
		v8.AddArg(v9)
		v10 := b.NewValue0(v.Line, OpInt64Lo, config.fe.TypeUInt32())
		v10.AddArg(y)
		v8.AddArg(v10)
		v7.AddArg(v8)
		v.AddArg(v7)
		return true
	}
}
func rewriteValuedec64_OpAnd64(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (And64 x y)
	// cond:
	// result: (Int64Make 		(And32 <config.fe.TypeUInt32()> (Int64Hi x) (Int64Hi y)) 		(And32 <config.fe.TypeUInt32()> (Int64Lo x) (Int64Lo y)))
	for {
		x := v.Args[0]
		y := v.Args[1]
		v.reset(OpInt64Make)
		v0 := b.NewValue0(v.Line, OpAnd32, config.fe.TypeUInt32())
		v1 := b.NewValue0(v.Line, OpInt64Hi, config.fe.TypeUInt32())
		v1.AddArg(x)
		v0.AddArg(v1)
		v2 := b.NewValue0(v.Line, OpInt64Hi, config.fe.TypeUInt32())
		v2.AddArg(y)
		v0.AddArg(v2)
		v.AddArg(v0)
		v3 := b.NewValue0(v.Line, OpAnd32, config.fe.TypeUInt32())
		v4 := b.NewValue0(v.Line, OpInt64Lo, config.fe.TypeUInt32())
		v4.AddArg(x)
		v3.AddArg(v4)
		v5 := b.NewValue0(v.Line, OpInt64Lo, config.fe.TypeUInt32())
		v5.AddArg(y)
		v3.AddArg(v5)
		v.AddArg(v3)
		return true
	}
}
func rewriteValuedec64_OpArg(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (Arg {n} [off])
	// cond: is64BitInt(v.Type) && v.Type.IsSigned()
	// result: (Int64Make     (Arg <config.fe.TypeInt32()> {n} [off+4])     (Arg <config.fe.TypeUInt32()> {n} [off]))
	for {
		n := v.Aux
		off := v.AuxInt
		if !(is64BitInt(v.Type) && v.Type.IsSigned()) {
			break
		}
		v.reset(OpInt64Make)
		v0 := b.NewValue0(v.Line, OpArg, config.fe.TypeInt32())
		v0.Aux = n
		v0.AuxInt = off + 4
		v.AddArg(v0)
		v1 := b.NewValue0(v.Line, OpArg, config.fe.TypeUInt32())
		v1.Aux = n
		v1.AuxInt = off
		v.AddArg(v1)
		return true
	}
	// match: (Arg {n} [off])
	// cond: is64BitInt(v.Type) && !v.Type.IsSigned()
	// result: (Int64Make     (Arg <config.fe.TypeUInt32()> {n} [off+4])     (Arg <config.fe.TypeUInt32()> {n} [off]))
	for {
		n := v.Aux
		off := v.AuxInt
		if !(is64BitInt(v.Type) && !v.Type.IsSigned()) {
			break
		}
		v.reset(OpInt64Make)
		v0 := b.NewValue0(v.Line, OpArg, config.fe.TypeUInt32())
		v0.Aux = n
		v0.AuxInt = off + 4
		v.AddArg(v0)
		v1 := b.NewValue0(v.Line, OpArg, config.fe.TypeUInt32())
		v1.Aux = n
		v1.AuxInt = off
		v.AddArg(v1)
		return true
	}
	return false
}
func rewriteValuedec64_OpCom64(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (Com64 x)
	// cond:
	// result: (Int64Make 		(Com32 <config.fe.TypeUInt32()> (Int64Hi x)) 		(Com32 <config.fe.TypeUInt32()> (Int64Lo x)))
	for {
		x := v.Args[0]
		v.reset(OpInt64Make)
		v0 := b.NewValue0(v.Line, OpCom32, config.fe.TypeUInt32())
		v1 := b.NewValue0(v.Line, OpInt64Hi, config.fe.TypeUInt32())
		v1.AddArg(x)
		v0.AddArg(v1)
		v.AddArg(v0)
		v2 := b.NewValue0(v.Line, OpCom32, config.fe.TypeUInt32())
		v3 := b.NewValue0(v.Line, OpInt64Lo, config.fe.TypeUInt32())
		v3.AddArg(x)
		v2.AddArg(v3)
		v.AddArg(v2)
		return true
	}
}
func rewriteValuedec64_OpConst64(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (Const64 <t> [c])
	// cond: t.IsSigned()
	// result: (Int64Make (Const32 <config.fe.TypeInt32()> [c>>32]) (Const32 <config.fe.TypeUInt32()> [c&0xffffffff]))
	for {
		t := v.Type
		c := v.AuxInt
		if !(t.IsSigned()) {
			break
		}
		v.reset(OpInt64Make)
		v0 := b.NewValue0(v.Line, OpConst32, config.fe.TypeInt32())
		v0.AuxInt = c >> 32
		v.AddArg(v0)
		v1 := b.NewValue0(v.Line, OpConst32, config.fe.TypeUInt32())
		v1.AuxInt = c & 0xffffffff
		v.AddArg(v1)
		return true
	}
	// match: (Const64 <t> [c])
	// cond: !t.IsSigned()
	// result: (Int64Make (Const32 <config.fe.TypeUInt32()> [c>>32]) (Const32 <config.fe.TypeUInt32()> [c&0xffffffff]))
	for {
		t := v.Type
		c := v.AuxInt
		if !(!t.IsSigned()) {
			break
		}
		v.reset(OpInt64Make)
		v0 := b.NewValue0(v.Line, OpConst32, config.fe.TypeUInt32())
		v0.AuxInt = c >> 32
		v.AddArg(v0)
		v1 := b.NewValue0(v.Line, OpConst32, config.fe.TypeUInt32())
		v1.AuxInt = c & 0xffffffff
		v.AddArg(v1)
		return true
	}
	return false
}
func rewriteValuedec64_OpEq64(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (Eq64 x y)
	// cond:
	// result: (AndB 		(Eq32 (Int64Hi x) (Int64Hi y)) 		(Eq32 (Int64Lo x) (Int64Lo y)))
	for {
		x := v.Args[0]
		y := v.Args[1]
		v.reset(OpAndB)
		v0 := b.NewValue0(v.Line, OpEq32, config.fe.TypeBool())
		v1 := b.NewValue0(v.Line, OpInt64Hi, config.fe.TypeUInt32())
		v1.AddArg(x)
		v0.AddArg(v1)
		v2 := b.NewValue0(v.Line, OpInt64Hi, config.fe.TypeUInt32())
		v2.AddArg(y)
		v0.AddArg(v2)
		v.AddArg(v0)
		v3 := b.NewValue0(v.Line, OpEq32, config.fe.TypeBool())
		v4 := b.NewValue0(v.Line, OpInt64Lo, config.fe.TypeUInt32())
		v4.AddArg(x)
		v3.AddArg(v4)
		v5 := b.NewValue0(v.Line, OpInt64Lo, config.fe.TypeUInt32())
		v5.AddArg(y)
		v3.AddArg(v5)
		v.AddArg(v3)
		return true
	}
}
func rewriteValuedec64_OpGeq64(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (Geq64 x y)
	// cond:
	// result: (OrB 		(Greater32 (Int64Hi x) (Int64Hi y)) 		(AndB 			(Eq32 (Int64Hi x) (Int64Hi y)) 			(Geq32U (Int64Lo x) (Int64Lo y))))
	for {
		x := v.Args[0]
		y := v.Args[1]
		v.reset(OpOrB)
		v0 := b.NewValue0(v.Line, OpGreater32, config.fe.TypeBool())
		v1 := b.NewValue0(v.Line, OpInt64Hi, config.fe.TypeUInt32())
		v1.AddArg(x)
		v0.AddArg(v1)
		v2 := b.NewValue0(v.Line, OpInt64Hi, config.fe.TypeUInt32())
		v2.AddArg(y)
		v0.AddArg(v2)
		v.AddArg(v0)
		v3 := b.NewValue0(v.Line, OpAndB, config.fe.TypeBool())
		v4 := b.NewValue0(v.Line, OpEq32, config.fe.TypeBool())
		v5 := b.NewValue0(v.Line, OpInt64Hi, config.fe.TypeUInt32())
		v5.AddArg(x)
		v4.AddArg(v5)
		v6 := b.NewValue0(v.Line, OpInt64Hi, config.fe.TypeUInt32())
		v6.AddArg(y)
		v4.AddArg(v6)
		v3.AddArg(v4)
		v7 := b.NewValue0(v.Line, OpGeq32U, config.fe.TypeBool())
		v8 := b.NewValue0(v.Line, OpInt64Lo, config.fe.TypeUInt32())
		v8.AddArg(x)
		v7.AddArg(v8)
		v9 := b.NewValue0(v.Line, OpInt64Lo, config.fe.TypeUInt32())
		v9.AddArg(y)
		v7.AddArg(v9)
		v3.AddArg(v7)
		v.AddArg(v3)
		return true
	}
}
func rewriteValuedec64_OpGeq64U(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (Geq64U x y)
	// cond:
	// result: (OrB 		(Greater32U (Int64Hi x) (Int64Hi y)) 		(AndB 			(Eq32 (Int64Hi x) (Int64Hi y)) 			(Geq32U (Int64Lo x) (Int64Lo y))))
	for {
		x := v.Args[0]
		y := v.Args[1]
		v.reset(OpOrB)
		v0 := b.NewValue0(v.Line, OpGreater32U, config.fe.TypeBool())
		v1 := b.NewValue0(v.Line, OpInt64Hi, config.fe.TypeUInt32())
		v1.AddArg(x)
		v0.AddArg(v1)
		v2 := b.NewValue0(v.Line, OpInt64Hi, config.fe.TypeUInt32())
		v2.AddArg(y)
		v0.AddArg(v2)
		v.AddArg(v0)
		v3 := b.NewValue0(v.Line, OpAndB, config.fe.TypeBool())
		v4 := b.NewValue0(v.Line, OpEq32, config.fe.TypeBool())
		v5 := b.NewValue0(v.Line, OpInt64Hi, config.fe.TypeUInt32())
		v5.AddArg(x)
		v4.AddArg(v5)
		v6 := b.NewValue0(v.Line, OpInt64Hi, config.fe.TypeUInt32())
		v6.AddArg(y)
		v4.AddArg(v6)
		v3.AddArg(v4)
		v7 := b.NewValue0(v.Line, OpGeq32U, config.fe.TypeBool())
		v8 := b.NewValue0(v.Line, OpInt64Lo, config.fe.TypeUInt32())
		v8.AddArg(x)
		v7.AddArg(v8)
		v9 := b.NewValue0(v.Line, OpInt64Lo, config.fe.TypeUInt32())
		v9.AddArg(y)
		v7.AddArg(v9)
		v3.AddArg(v7)
		v.AddArg(v3)
		return true
	}
}
func rewriteValuedec64_OpGreater64(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (Greater64 x y)
	// cond:
	// result: (OrB 		(Greater32 (Int64Hi x) (Int64Hi y)) 		(AndB 			(Eq32 (Int64Hi x) (Int64Hi y)) 			(Greater32U (Int64Lo x) (Int64Lo y))))
	for {
		x := v.Args[0]
		y := v.Args[1]
		v.reset(OpOrB)
		v0 := b.NewValue0(v.Line, OpGreater32, config.fe.TypeBool())
		v1 := b.NewValue0(v.Line, OpInt64Hi, config.fe.TypeUInt32())
		v1.AddArg(x)
		v0.AddArg(v1)
		v2 := b.NewValue0(v.Line, OpInt64Hi, config.fe.TypeUInt32())
		v2.AddArg(y)
		v0.AddArg(v2)
		v.AddArg(v0)
		v3 := b.NewValue0(v.Line, OpAndB, config.fe.TypeBool())
		v4 := b.NewValue0(v.Line, OpEq32, config.fe.TypeBool())
		v5 := b.NewValue0(v.Line, OpInt64Hi, config.fe.TypeUInt32())
		v5.AddArg(x)
		v4.AddArg(v5)
		v6 := b.NewValue0(v.Line, OpInt64Hi, config.fe.TypeUInt32())
		v6.AddArg(y)
		v4.AddArg(v6)
		v3.AddArg(v4)
		v7 := b.NewValue0(v.Line, OpGreater32U, config.fe.TypeBool())
		v8 := b.NewValue0(v.Line, OpInt64Lo, config.fe.TypeUInt32())
		v8.AddArg(x)
		v7.AddArg(v8)
		v9 := b.NewValue0(v.Line, OpInt64Lo, config.fe.TypeUInt32())
		v9.AddArg(y)
		v7.AddArg(v9)
		v3.AddArg(v7)
		v.AddArg(v3)
		return true
	}
}
func rewriteValuedec64_OpGreater64U(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (Greater64U x y)
	// cond:
	// result: (OrB 		(Greater32U (Int64Hi x) (Int64Hi y)) 		(AndB 			(Eq32 (Int64Hi x) (Int64Hi y)) 			(Greater32U (Int64Lo x) (Int64Lo y))))
	for {
		x := v.Args[0]
		y := v.Args[1]
		v.reset(OpOrB)
		v0 := b.NewValue0(v.Line, OpGreater32U, config.fe.TypeBool())
		v1 := b.NewValue0(v.Line, OpInt64Hi, config.fe.TypeUInt32())
		v1.AddArg(x)
		v0.AddArg(v1)
		v2 := b.NewValue0(v.Line, OpInt64Hi, config.fe.TypeUInt32())
		v2.AddArg(y)
		v0.AddArg(v2)
		v.AddArg(v0)
		v3 := b.NewValue0(v.Line, OpAndB, config.fe.TypeBool())
		v4 := b.NewValue0(v.Line, OpEq32, config.fe.TypeBool())
		v5 := b.NewValue0(v.Line, OpInt64Hi, config.fe.TypeUInt32())
		v5.AddArg(x)
		v4.AddArg(v5)
		v6 := b.NewValue0(v.Line, OpInt64Hi, config.fe.TypeUInt32())
		v6.AddArg(y)
		v4.AddArg(v6)
		v3.AddArg(v4)
		v7 := b.NewValue0(v.Line, OpGreater32U, config.fe.TypeBool())
		v8 := b.NewValue0(v.Line, OpInt64Lo, config.fe.TypeUInt32())
		v8.AddArg(x)
		v7.AddArg(v8)
		v9 := b.NewValue0(v.Line, OpInt64Lo, config.fe.TypeUInt32())
		v9.AddArg(y)
		v7.AddArg(v9)
		v3.AddArg(v7)
		v.AddArg(v3)
		return true
	}
}
func rewriteValuedec64_OpInt64Hi(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (Int64Hi (Int64Make hi _))
	// cond:
	// result: hi
	for {
		v_0 := v.Args[0]
		if v_0.Op != OpInt64Make {
			break
		}
		hi := v_0.Args[0]
		v.reset(OpCopy)
		v.Type = hi.Type
		v.AddArg(hi)
		return true
	}
	return false
}
func rewriteValuedec64_OpInt64Lo(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (Int64Lo (Int64Make _ lo))
	// cond:
	// result: lo
	for {
		v_0 := v.Args[0]
		if v_0.Op != OpInt64Make {
			break
		}
		lo := v_0.Args[1]
		v.reset(OpCopy)
		v.Type = lo.Type
		v.AddArg(lo)
		return true
	}
	return false
}
func rewriteValuedec64_OpLeq64(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (Leq64 x y)
	// cond:
	// result: (OrB 		(Less32 (Int64Hi x) (Int64Hi y)) 		(AndB 			(Eq32 (Int64Hi x) (Int64Hi y)) 			(Leq32U (Int64Lo x) (Int64Lo y))))
	for {
		x := v.Args[0]
		y := v.Args[1]
		v.reset(OpOrB)
		v0 := b.NewValue0(v.Line, OpLess32, config.fe.TypeBool())
		v1 := b.NewValue0(v.Line, OpInt64Hi, config.fe.TypeUInt32())
		v1.AddArg(x)
		v0.AddArg(v1)
		v2 := b.NewValue0(v.Line, OpInt64Hi, config.fe.TypeUInt32())
		v2.AddArg(y)
		v0.AddArg(v2)
		v.AddArg(v0)
		v3 := b.NewValue0(v.Line, OpAndB, config.fe.TypeBool())
		v4 := b.NewValue0(v.Line, OpEq32, config.fe.TypeBool())
		v5 := b.NewValue0(v.Line, OpInt64Hi, config.fe.TypeUInt32())
		v5.AddArg(x)
		v4.AddArg(v5)
		v6 := b.NewValue0(v.Line, OpInt64Hi, config.fe.TypeUInt32())
		v6.AddArg(y)
		v4.AddArg(v6)
		v3.AddArg(v4)
		v7 := b.NewValue0(v.Line, OpLeq32U, config.fe.TypeBool())
		v8 := b.NewValue0(v.Line, OpInt64Lo, config.fe.TypeUInt32())
		v8.AddArg(x)
		v7.AddArg(v8)
		v9 := b.NewValue0(v.Line, OpInt64Lo, config.fe.TypeUInt32())
		v9.AddArg(y)
		v7.AddArg(v9)
		v3.AddArg(v7)
		v.AddArg(v3)
		return true
	}
}
func rewriteValuedec64_OpLeq64U(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (Leq64U x y)
	// cond:
	// result: (OrB 		(Less32U (Int64Hi x) (Int64Hi y)) 		(AndB 			(Eq32 (Int64Hi x) (Int64Hi y)) 			(Leq32U (Int64Lo x) (Int64Lo y))))
	for {
		x := v.Args[0]
		y := v.Args[1]
		v.reset(OpOrB)
		v0 := b.NewValue0(v.Line, OpLess32U, config.fe.TypeBool())
		v1 := b.NewValue0(v.Line, OpInt64Hi, config.fe.TypeUInt32())
		v1.AddArg(x)
		v0.AddArg(v1)
		v2 := b.NewValue0(v.Line, OpInt64Hi, config.fe.TypeUInt32())
		v2.AddArg(y)
		v0.AddArg(v2)
		v.AddArg(v0)
		v3 := b.NewValue0(v.Line, OpAndB, config.fe.TypeBool())
		v4 := b.NewValue0(v.Line, OpEq32, config.fe.TypeBool())
		v5 := b.NewValue0(v.Line, OpInt64Hi, config.fe.TypeUInt32())
		v5.AddArg(x)
		v4.AddArg(v5)
		v6 := b.NewValue0(v.Line, OpInt64Hi, config.fe.TypeUInt32())
		v6.AddArg(y)
		v4.AddArg(v6)
		v3.AddArg(v4)
		v7 := b.NewValue0(v.Line, OpLeq32U, config.fe.TypeBool())
		v8 := b.NewValue0(v.Line, OpInt64Lo, config.fe.TypeUInt32())
		v8.AddArg(x)
		v7.AddArg(v8)
		v9 := b.NewValue0(v.Line, OpInt64Lo, config.fe.TypeUInt32())
		v9.AddArg(y)
		v7.AddArg(v9)
		v3.AddArg(v7)
		v.AddArg(v3)
		return true
	}
}
func rewriteValuedec64_OpLess64(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (Less64 x y)
	// cond:
	// result: (OrB 		(Less32 (Int64Hi x) (Int64Hi y)) 		(AndB 			(Eq32 (Int64Hi x) (Int64Hi y)) 			(Less32U (Int64Lo x) (Int64Lo y))))
	for {
		x := v.Args[0]
		y := v.Args[1]
		v.reset(OpOrB)
		v0 := b.NewValue0(v.Line, OpLess32, config.fe.TypeBool())
		v1 := b.NewValue0(v.Line, OpInt64Hi, config.fe.TypeUInt32())
		v1.AddArg(x)
		v0.AddArg(v1)
		v2 := b.NewValue0(v.Line, OpInt64Hi, config.fe.TypeUInt32())
		v2.AddArg(y)
		v0.AddArg(v2)
		v.AddArg(v0)
		v3 := b.NewValue0(v.Line, OpAndB, config.fe.TypeBool())
		v4 := b.NewValue0(v.Line, OpEq32, config.fe.TypeBool())
		v5 := b.NewValue0(v.Line, OpInt64Hi, config.fe.TypeUInt32())
		v5.AddArg(x)
		v4.AddArg(v5)
		v6 := b.NewValue0(v.Line, OpInt64Hi, config.fe.TypeUInt32())
		v6.AddArg(y)
		v4.AddArg(v6)
		v3.AddArg(v4)
		v7 := b.NewValue0(v.Line, OpLess32U, config.fe.TypeBool())
		v8 := b.NewValue0(v.Line, OpInt64Lo, config.fe.TypeUInt32())
		v8.AddArg(x)
		v7.AddArg(v8)
		v9 := b.NewValue0(v.Line, OpInt64Lo, config.fe.TypeUInt32())
		v9.AddArg(y)
		v7.AddArg(v9)
		v3.AddArg(v7)
		v.AddArg(v3)
		return true
	}
}
func rewriteValuedec64_OpLess64U(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (Less64U x y)
	// cond:
	// result: (OrB 		(Less32U (Int64Hi x) (Int64Hi y)) 		(AndB 			(Eq32 (Int64Hi x) (Int64Hi y)) 			(Less32U (Int64Lo x) (Int64Lo y))))
	for {
		x := v.Args[0]
		y := v.Args[1]
		v.reset(OpOrB)
		v0 := b.NewValue0(v.Line, OpLess32U, config.fe.TypeBool())
		v1 := b.NewValue0(v.Line, OpInt64Hi, config.fe.TypeUInt32())
		v1.AddArg(x)
		v0.AddArg(v1)
		v2 := b.NewValue0(v.Line, OpInt64Hi, config.fe.TypeUInt32())
		v2.AddArg(y)
		v0.AddArg(v2)
		v.AddArg(v0)
		v3 := b.NewValue0(v.Line, OpAndB, config.fe.TypeBool())
		v4 := b.NewValue0(v.Line, OpEq32, config.fe.TypeBool())
		v5 := b.NewValue0(v.Line, OpInt64Hi, config.fe.TypeUInt32())
		v5.AddArg(x)
		v4.AddArg(v5)
		v6 := b.NewValue0(v.Line, OpInt64Hi, config.fe.TypeUInt32())
		v6.AddArg(y)
		v4.AddArg(v6)
		v3.AddArg(v4)
		v7 := b.NewValue0(v.Line, OpLess32U, config.fe.TypeBool())
		v8 := b.NewValue0(v.Line, OpInt64Lo, config.fe.TypeUInt32())
		v8.AddArg(x)
		v7.AddArg(v8)
		v9 := b.NewValue0(v.Line, OpInt64Lo, config.fe.TypeUInt32())
		v9.AddArg(y)
		v7.AddArg(v9)
		v3.AddArg(v7)
		v.AddArg(v3)
		return true
	}
}
func rewriteValuedec64_OpLoad(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (Load <t> ptr mem)
	// cond: is64BitInt(t) && t.IsSigned()
	// result: (Int64Make 		(Load <config.fe.TypeInt32()> (OffPtr <config.fe.TypeInt32().PtrTo()> [4] ptr) mem) 		(Load <config.fe.TypeUInt32()> ptr mem))
	for {
		t := v.Type
		ptr := v.Args[0]
		mem := v.Args[1]
		if !(is64BitInt(t) && t.IsSigned()) {
			break
		}
		v.reset(OpInt64Make)
		v0 := b.NewValue0(v.Line, OpLoad, config.fe.TypeInt32())
		v1 := b.NewValue0(v.Line, OpOffPtr, config.fe.TypeInt32().PtrTo())
		v1.AuxInt = 4
		v1.AddArg(ptr)
		v0.AddArg(v1)
		v0.AddArg(mem)
		v.AddArg(v0)
		v2 := b.NewValue0(v.Line, OpLoad, config.fe.TypeUInt32())
		v2.AddArg(ptr)
		v2.AddArg(mem)
		v.AddArg(v2)
		return true
	}
	// match: (Load <t> ptr mem)
	// cond: is64BitInt(t) && !t.IsSigned()
	// result: (Int64Make 		(Load <config.fe.TypeUInt32()> (OffPtr <config.fe.TypeUInt32().PtrTo()> [4] ptr) mem) 		(Load <config.fe.TypeUInt32()> ptr mem))
	for {
		t := v.Type
		ptr := v.Args[0]
		mem := v.Args[1]
		if !(is64BitInt(t) && !t.IsSigned()) {
			break
		}
		v.reset(OpInt64Make)
		v0 := b.NewValue0(v.Line, OpLoad, config.fe.TypeUInt32())
		v1 := b.NewValue0(v.Line, OpOffPtr, config.fe.TypeUInt32().PtrTo())
		v1.AuxInt = 4
		v1.AddArg(ptr)
		v0.AddArg(v1)
		v0.AddArg(mem)
		v.AddArg(v0)
		v2 := b.NewValue0(v.Line, OpLoad, config.fe.TypeUInt32())
		v2.AddArg(ptr)
		v2.AddArg(mem)
		v.AddArg(v2)
		return true
	}
	return false
}
func rewriteValuedec64_OpLsh16x64(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (Lsh16x64 _ (Int64Make (Const32 [c]) _))
	// cond: c != 0
	// result: (Const32 [0])
	for {
		v_1 := v.Args[1]
		if v_1.Op != OpInt64Make {
			break
		}
		v_1_0 := v_1.Args[0]
		if v_1_0.Op != OpConst32 {
			break
		}
		c := v_1_0.AuxInt
		if !(c != 0) {
			break
		}
		v.reset(OpConst32)
		v.AuxInt = 0
		return true
	}
	// match: (Lsh16x64 x (Int64Make (Const32 [0]) lo))
	// cond:
	// result: (Lsh16x32 x lo)
	for {
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpInt64Make {
			break
		}
		v_1_0 := v_1.Args[0]
		if v_1_0.Op != OpConst32 {
			break
		}
		if v_1_0.AuxInt != 0 {
			break
		}
		lo := v_1.Args[1]
		v.reset(OpLsh16x32)
		v.AddArg(x)
		v.AddArg(lo)
		return true
	}
	return false
}
func rewriteValuedec64_OpLsh32x64(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (Lsh32x64 _ (Int64Make (Const32 [c]) _))
	// cond: c != 0
	// result: (Const32 [0])
	for {
		v_1 := v.Args[1]
		if v_1.Op != OpInt64Make {
			break
		}
		v_1_0 := v_1.Args[0]
		if v_1_0.Op != OpConst32 {
			break
		}
		c := v_1_0.AuxInt
		if !(c != 0) {
			break
		}
		v.reset(OpConst32)
		v.AuxInt = 0
		return true
	}
	// match: (Lsh32x64 x (Int64Make (Const32 [0]) lo))
	// cond:
	// result: (Lsh32x32 x lo)
	for {
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpInt64Make {
			break
		}
		v_1_0 := v_1.Args[0]
		if v_1_0.Op != OpConst32 {
			break
		}
		if v_1_0.AuxInt != 0 {
			break
		}
		lo := v_1.Args[1]
		v.reset(OpLsh32x32)
		v.AddArg(x)
		v.AddArg(lo)
		return true
	}
	return false
}
func rewriteValuedec64_OpLsh8x64(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (Lsh8x64 _ (Int64Make (Const32 [c]) _))
	// cond: c != 0
	// result: (Const32 [0])
	for {
		v_1 := v.Args[1]
		if v_1.Op != OpInt64Make {
			break
		}
		v_1_0 := v_1.Args[0]
		if v_1_0.Op != OpConst32 {
			break
		}
		c := v_1_0.AuxInt
		if !(c != 0) {
			break
		}
		v.reset(OpConst32)
		v.AuxInt = 0
		return true
	}
	// match: (Lsh8x64 x (Int64Make (Const32 [0]) lo))
	// cond:
	// result: (Lsh8x32 x lo)
	for {
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpInt64Make {
			break
		}
		v_1_0 := v_1.Args[0]
		if v_1_0.Op != OpConst32 {
			break
		}
		if v_1_0.AuxInt != 0 {
			break
		}
		lo := v_1.Args[1]
		v.reset(OpLsh8x32)
		v.AddArg(x)
		v.AddArg(lo)
		return true
	}
	return false
}
func rewriteValuedec64_OpMul64(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (Mul64 x y)
	// cond:
	// result: (Int64Make 		(Add32 <config.fe.TypeUInt32()> 			(Mul32 <config.fe.TypeUInt32()> (Int64Lo x) (Int64Hi y)) 			(Add32 <config.fe.TypeUInt32()> 				(Mul32 <config.fe.TypeUInt32()> (Int64Hi x) (Int64Lo y)) 				(Select0 <config.fe.TypeUInt32()> (Mul32uhilo (Int64Lo x) (Int64Lo y))))) 		(Select1 <config.fe.TypeUInt32()> (Mul32uhilo (Int64Lo x) (Int64Lo y))))
	for {
		x := v.Args[0]
		y := v.Args[1]
		v.reset(OpInt64Make)
		v0 := b.NewValue0(v.Line, OpAdd32, config.fe.TypeUInt32())
		v1 := b.NewValue0(v.Line, OpMul32, config.fe.TypeUInt32())
		v2 := b.NewValue0(v.Line, OpInt64Lo, config.fe.TypeUInt32())
		v2.AddArg(x)
		v1.AddArg(v2)
		v3 := b.NewValue0(v.Line, OpInt64Hi, config.fe.TypeUInt32())
		v3.AddArg(y)
		v1.AddArg(v3)
		v0.AddArg(v1)
		v4 := b.NewValue0(v.Line, OpAdd32, config.fe.TypeUInt32())
		v5 := b.NewValue0(v.Line, OpMul32, config.fe.TypeUInt32())
		v6 := b.NewValue0(v.Line, OpInt64Hi, config.fe.TypeUInt32())
		v6.AddArg(x)
		v5.AddArg(v6)
		v7 := b.NewValue0(v.Line, OpInt64Lo, config.fe.TypeUInt32())
		v7.AddArg(y)
		v5.AddArg(v7)
		v4.AddArg(v5)
		v8 := b.NewValue0(v.Line, OpSelect0, config.fe.TypeUInt32())
		v9 := b.NewValue0(v.Line, OpMul32uhilo, MakeTuple(config.fe.TypeUInt32(), config.fe.TypeUInt32()))
		v10 := b.NewValue0(v.Line, OpInt64Lo, config.fe.TypeUInt32())
		v10.AddArg(x)
		v9.AddArg(v10)
		v11 := b.NewValue0(v.Line, OpInt64Lo, config.fe.TypeUInt32())
		v11.AddArg(y)
		v9.AddArg(v11)
		v8.AddArg(v9)
		v4.AddArg(v8)
		v0.AddArg(v4)
		v.AddArg(v0)
		v12 := b.NewValue0(v.Line, OpSelect1, config.fe.TypeUInt32())
		v13 := b.NewValue0(v.Line, OpMul32uhilo, MakeTuple(config.fe.TypeUInt32(), config.fe.TypeUInt32()))
		v14 := b.NewValue0(v.Line, OpInt64Lo, config.fe.TypeUInt32())
		v14.AddArg(x)
		v13.AddArg(v14)
		v15 := b.NewValue0(v.Line, OpInt64Lo, config.fe.TypeUInt32())
		v15.AddArg(y)
		v13.AddArg(v15)
		v12.AddArg(v13)
		v.AddArg(v12)
		return true
	}
}
func rewriteValuedec64_OpNeg64(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (Neg64 <t> x)
	// cond:
	// result: (Sub64 (Const64 <t> [0]) x)
	for {
		t := v.Type
		x := v.Args[0]
		v.reset(OpSub64)
		v0 := b.NewValue0(v.Line, OpConst64, t)
		v0.AuxInt = 0
		v.AddArg(v0)
		v.AddArg(x)
		return true
	}
}
func rewriteValuedec64_OpNeq64(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (Neq64 x y)
	// cond:
	// result: (OrB 		(Neq32 (Int64Hi x) (Int64Hi y)) 		(Neq32 (Int64Lo x) (Int64Lo y)))
	for {
		x := v.Args[0]
		y := v.Args[1]
		v.reset(OpOrB)
		v0 := b.NewValue0(v.Line, OpNeq32, config.fe.TypeBool())
		v1 := b.NewValue0(v.Line, OpInt64Hi, config.fe.TypeUInt32())
		v1.AddArg(x)
		v0.AddArg(v1)
		v2 := b.NewValue0(v.Line, OpInt64Hi, config.fe.TypeUInt32())
		v2.AddArg(y)
		v0.AddArg(v2)
		v.AddArg(v0)
		v3 := b.NewValue0(v.Line, OpNeq32, config.fe.TypeBool())
		v4 := b.NewValue0(v.Line, OpInt64Lo, config.fe.TypeUInt32())
		v4.AddArg(x)
		v3.AddArg(v4)
		v5 := b.NewValue0(v.Line, OpInt64Lo, config.fe.TypeUInt32())
		v5.AddArg(y)
		v3.AddArg(v5)
		v.AddArg(v3)
		return true
	}
}
func rewriteValuedec64_OpOr64(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (Or64 x y)
	// cond:
	// result: (Int64Make 		(Or32 <config.fe.TypeUInt32()> (Int64Hi x) (Int64Hi y)) 		(Or32 <config.fe.TypeUInt32()> (Int64Lo x) (Int64Lo y)))
	for {
		x := v.Args[0]
		y := v.Args[1]
		v.reset(OpInt64Make)
		v0 := b.NewValue0(v.Line, OpOr32, config.fe.TypeUInt32())
		v1 := b.NewValue0(v.Line, OpInt64Hi, config.fe.TypeUInt32())
		v1.AddArg(x)
		v0.AddArg(v1)
		v2 := b.NewValue0(v.Line, OpInt64Hi, config.fe.TypeUInt32())
		v2.AddArg(y)
		v0.AddArg(v2)
		v.AddArg(v0)
		v3 := b.NewValue0(v.Line, OpOr32, config.fe.TypeUInt32())
		v4 := b.NewValue0(v.Line, OpInt64Lo, config.fe.TypeUInt32())
		v4.AddArg(x)
		v3.AddArg(v4)
		v5 := b.NewValue0(v.Line, OpInt64Lo, config.fe.TypeUInt32())
		v5.AddArg(y)
		v3.AddArg(v5)
		v.AddArg(v3)
		return true
	}
}
func rewriteValuedec64_OpRsh16Ux64(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (Rsh16Ux64 _ (Int64Make (Const32 [c]) _))
	// cond: c != 0
	// result: (Const32 [0])
	for {
		v_1 := v.Args[1]
		if v_1.Op != OpInt64Make {
			break
		}
		v_1_0 := v_1.Args[0]
		if v_1_0.Op != OpConst32 {
			break
		}
		c := v_1_0.AuxInt
		if !(c != 0) {
			break
		}
		v.reset(OpConst32)
		v.AuxInt = 0
		return true
	}
	// match: (Rsh16Ux64 x (Int64Make (Const32 [0]) lo))
	// cond:
	// result: (Rsh16Ux32 x lo)
	for {
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpInt64Make {
			break
		}
		v_1_0 := v_1.Args[0]
		if v_1_0.Op != OpConst32 {
			break
		}
		if v_1_0.AuxInt != 0 {
			break
		}
		lo := v_1.Args[1]
		v.reset(OpRsh16Ux32)
		v.AddArg(x)
		v.AddArg(lo)
		return true
	}
	return false
}
func rewriteValuedec64_OpRsh16x64(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (Rsh16x64 x (Int64Make (Const32 [c]) _))
	// cond: c != 0
	// result: (Signmask (SignExt16to32 x))
	for {
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpInt64Make {
			break
		}
		v_1_0 := v_1.Args[0]
		if v_1_0.Op != OpConst32 {
			break
		}
		c := v_1_0.AuxInt
		if !(c != 0) {
			break
		}
		v.reset(OpSignmask)
		v0 := b.NewValue0(v.Line, OpSignExt16to32, config.fe.TypeInt32())
		v0.AddArg(x)
		v.AddArg(v0)
		return true
	}
	// match: (Rsh16x64 x (Int64Make (Const32 [0]) lo))
	// cond:
	// result: (Rsh16x32 x lo)
	for {
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpInt64Make {
			break
		}
		v_1_0 := v_1.Args[0]
		if v_1_0.Op != OpConst32 {
			break
		}
		if v_1_0.AuxInt != 0 {
			break
		}
		lo := v_1.Args[1]
		v.reset(OpRsh16x32)
		v.AddArg(x)
		v.AddArg(lo)
		return true
	}
	return false
}
func rewriteValuedec64_OpRsh32Ux64(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (Rsh32Ux64 _ (Int64Make (Const32 [c]) _))
	// cond: c != 0
	// result: (Const32 [0])
	for {
		v_1 := v.Args[1]
		if v_1.Op != OpInt64Make {
			break
		}
		v_1_0 := v_1.Args[0]
		if v_1_0.Op != OpConst32 {
			break
		}
		c := v_1_0.AuxInt
		if !(c != 0) {
			break
		}
		v.reset(OpConst32)
		v.AuxInt = 0
		return true
	}
	// match: (Rsh32Ux64 x (Int64Make (Const32 [0]) lo))
	// cond:
	// result: (Rsh32Ux32 x lo)
	for {
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpInt64Make {
			break
		}
		v_1_0 := v_1.Args[0]
		if v_1_0.Op != OpConst32 {
			break
		}
		if v_1_0.AuxInt != 0 {
			break
		}
		lo := v_1.Args[1]
		v.reset(OpRsh32Ux32)
		v.AddArg(x)
		v.AddArg(lo)
		return true
	}
	return false
}
func rewriteValuedec64_OpRsh32x64(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (Rsh32x64 x (Int64Make (Const32 [c]) _))
	// cond: c != 0
	// result: (Signmask x)
	for {
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpInt64Make {
			break
		}
		v_1_0 := v_1.Args[0]
		if v_1_0.Op != OpConst32 {
			break
		}
		c := v_1_0.AuxInt
		if !(c != 0) {
			break
		}
		v.reset(OpSignmask)
		v.AddArg(x)
		return true
	}
	// match: (Rsh32x64 x (Int64Make (Const32 [0]) lo))
	// cond:
	// result: (Rsh32x32 x lo)
	for {
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpInt64Make {
			break
		}
		v_1_0 := v_1.Args[0]
		if v_1_0.Op != OpConst32 {
			break
		}
		if v_1_0.AuxInt != 0 {
			break
		}
		lo := v_1.Args[1]
		v.reset(OpRsh32x32)
		v.AddArg(x)
		v.AddArg(lo)
		return true
	}
	return false
}
func rewriteValuedec64_OpRsh8Ux64(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (Rsh8Ux64 _ (Int64Make (Const32 [c]) _))
	// cond: c != 0
	// result: (Const32 [0])
	for {
		v_1 := v.Args[1]
		if v_1.Op != OpInt64Make {
			break
		}
		v_1_0 := v_1.Args[0]
		if v_1_0.Op != OpConst32 {
			break
		}
		c := v_1_0.AuxInt
		if !(c != 0) {
			break
		}
		v.reset(OpConst32)
		v.AuxInt = 0
		return true
	}
	// match: (Rsh8Ux64 x (Int64Make (Const32 [0]) lo))
	// cond:
	// result: (Rsh8Ux32 x lo)
	for {
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpInt64Make {
			break
		}
		v_1_0 := v_1.Args[0]
		if v_1_0.Op != OpConst32 {
			break
		}
		if v_1_0.AuxInt != 0 {
			break
		}
		lo := v_1.Args[1]
		v.reset(OpRsh8Ux32)
		v.AddArg(x)
		v.AddArg(lo)
		return true
	}
	return false
}
func rewriteValuedec64_OpRsh8x64(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (Rsh8x64 x (Int64Make (Const32 [c]) _))
	// cond: c != 0
	// result: (Signmask (SignExt8to32 x))
	for {
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpInt64Make {
			break
		}
		v_1_0 := v_1.Args[0]
		if v_1_0.Op != OpConst32 {
			break
		}
		c := v_1_0.AuxInt
		if !(c != 0) {
			break
		}
		v.reset(OpSignmask)
		v0 := b.NewValue0(v.Line, OpSignExt8to32, config.fe.TypeInt32())
		v0.AddArg(x)
		v.AddArg(v0)
		return true
	}
	// match: (Rsh8x64 x (Int64Make (Const32 [0]) lo))
	// cond:
	// result: (Rsh8x32 x lo)
	for {
		x := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpInt64Make {
			break
		}
		v_1_0 := v_1.Args[0]
		if v_1_0.Op != OpConst32 {
			break
		}
		if v_1_0.AuxInt != 0 {
			break
		}
		lo := v_1.Args[1]
		v.reset(OpRsh8x32)
		v.AddArg(x)
		v.AddArg(lo)
		return true
	}
	return false
}
func rewriteValuedec64_OpSignExt16to64(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (SignExt16to64 x)
	// cond:
	// result: (SignExt32to64 (SignExt16to32 x))
	for {
		x := v.Args[0]
		v.reset(OpSignExt32to64)
		v0 := b.NewValue0(v.Line, OpSignExt16to32, config.fe.TypeInt32())
		v0.AddArg(x)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuedec64_OpSignExt32to64(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (SignExt32to64 x)
	// cond:
	// result: (Int64Make (Signmask x) x)
	for {
		x := v.Args[0]
		v.reset(OpInt64Make)
		v0 := b.NewValue0(v.Line, OpSignmask, config.fe.TypeInt32())
		v0.AddArg(x)
		v.AddArg(v0)
		v.AddArg(x)
		return true
	}
}
func rewriteValuedec64_OpSignExt8to64(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (SignExt8to64 x)
	// cond:
	// result: (SignExt32to64 (SignExt8to32 x))
	for {
		x := v.Args[0]
		v.reset(OpSignExt32to64)
		v0 := b.NewValue0(v.Line, OpSignExt8to32, config.fe.TypeInt32())
		v0.AddArg(x)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuedec64_OpStore(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (Store [8] dst (Int64Make hi lo) mem)
	// cond:
	// result: (Store [4] 		(OffPtr <hi.Type.PtrTo()> [4] dst) 		hi 		(Store [4] dst lo mem))
	for {
		if v.AuxInt != 8 {
			break
		}
		dst := v.Args[0]
		v_1 := v.Args[1]
		if v_1.Op != OpInt64Make {
			break
		}
		hi := v_1.Args[0]
		lo := v_1.Args[1]
		mem := v.Args[2]
		v.reset(OpStore)
		v.AuxInt = 4
		v0 := b.NewValue0(v.Line, OpOffPtr, hi.Type.PtrTo())
		v0.AuxInt = 4
		v0.AddArg(dst)
		v.AddArg(v0)
		v.AddArg(hi)
		v1 := b.NewValue0(v.Line, OpStore, TypeMem)
		v1.AuxInt = 4
		v1.AddArg(dst)
		v1.AddArg(lo)
		v1.AddArg(mem)
		v.AddArg(v1)
		return true
	}
	return false
}
func rewriteValuedec64_OpSub64(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (Sub64 x y)
	// cond:
	// result: (Int64Make 		(Sub32withcarry <config.fe.TypeInt32()> 			(Int64Hi x) 			(Int64Hi y) 			(Select0 <TypeFlags> (Sub32carry (Int64Lo x) (Int64Lo y)))) 		(Select1 <config.fe.TypeUInt32()> (Sub32carry (Int64Lo x) (Int64Lo y))))
	for {
		x := v.Args[0]
		y := v.Args[1]
		v.reset(OpInt64Make)
		v0 := b.NewValue0(v.Line, OpSub32withcarry, config.fe.TypeInt32())
		v1 := b.NewValue0(v.Line, OpInt64Hi, config.fe.TypeUInt32())
		v1.AddArg(x)
		v0.AddArg(v1)
		v2 := b.NewValue0(v.Line, OpInt64Hi, config.fe.TypeUInt32())
		v2.AddArg(y)
		v0.AddArg(v2)
		v3 := b.NewValue0(v.Line, OpSelect0, TypeFlags)
		v4 := b.NewValue0(v.Line, OpSub32carry, MakeTuple(TypeFlags, config.fe.TypeUInt32()))
		v5 := b.NewValue0(v.Line, OpInt64Lo, config.fe.TypeUInt32())
		v5.AddArg(x)
		v4.AddArg(v5)
		v6 := b.NewValue0(v.Line, OpInt64Lo, config.fe.TypeUInt32())
		v6.AddArg(y)
		v4.AddArg(v6)
		v3.AddArg(v4)
		v0.AddArg(v3)
		v.AddArg(v0)
		v7 := b.NewValue0(v.Line, OpSelect1, config.fe.TypeUInt32())
		v8 := b.NewValue0(v.Line, OpSub32carry, MakeTuple(TypeFlags, config.fe.TypeUInt32()))
		v9 := b.NewValue0(v.Line, OpInt64Lo, config.fe.TypeUInt32())
		v9.AddArg(x)
		v8.AddArg(v9)
		v10 := b.NewValue0(v.Line, OpInt64Lo, config.fe.TypeUInt32())
		v10.AddArg(y)
		v8.AddArg(v10)
		v7.AddArg(v8)
		v.AddArg(v7)
		return true
	}
}
func rewriteValuedec64_OpTrunc64to16(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (Trunc64to16 (Int64Make _ lo))
	// cond:
	// result: (Trunc32to16 lo)
	for {
		v_0 := v.Args[0]
		if v_0.Op != OpInt64Make {
			break
		}
		lo := v_0.Args[1]
		v.reset(OpTrunc32to16)
		v.AddArg(lo)
		return true
	}
	return false
}
func rewriteValuedec64_OpTrunc64to32(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (Trunc64to32 (Int64Make _ lo))
	// cond:
	// result: lo
	for {
		v_0 := v.Args[0]
		if v_0.Op != OpInt64Make {
			break
		}
		lo := v_0.Args[1]
		v.reset(OpCopy)
		v.Type = lo.Type
		v.AddArg(lo)
		return true
	}
	return false
}
func rewriteValuedec64_OpTrunc64to8(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (Trunc64to8 (Int64Make _ lo))
	// cond:
	// result: (Trunc32to8 lo)
	for {
		v_0 := v.Args[0]
		if v_0.Op != OpInt64Make {
			break
		}
		lo := v_0.Args[1]
		v.reset(OpTrunc32to8)
		v.AddArg(lo)
		return true
	}
	return false
}
func rewriteValuedec64_OpXor64(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (Xor64 x y)
	// cond:
	// result: (Int64Make 		(Xor32 <config.fe.TypeUInt32()> (Int64Hi x) (Int64Hi y)) 		(Xor32 <config.fe.TypeUInt32()> (Int64Lo x) (Int64Lo y)))
	for {
		x := v.Args[0]
		y := v.Args[1]
		v.reset(OpInt64Make)
		v0 := b.NewValue0(v.Line, OpXor32, config.fe.TypeUInt32())
		v1 := b.NewValue0(v.Line, OpInt64Hi, config.fe.TypeUInt32())
		v1.AddArg(x)
		v0.AddArg(v1)
		v2 := b.NewValue0(v.Line, OpInt64Hi, config.fe.TypeUInt32())
		v2.AddArg(y)
		v0.AddArg(v2)
		v.AddArg(v0)
		v3 := b.NewValue0(v.Line, OpXor32, config.fe.TypeUInt32())
		v4 := b.NewValue0(v.Line, OpInt64Lo, config.fe.TypeUInt32())
		v4.AddArg(x)
		v3.AddArg(v4)
		v5 := b.NewValue0(v.Line, OpInt64Lo, config.fe.TypeUInt32())
		v5.AddArg(y)
		v3.AddArg(v5)
		v.AddArg(v3)
		return true
	}
}
func rewriteValuedec64_OpZeroExt16to64(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (ZeroExt16to64 x)
	// cond:
	// result: (ZeroExt32to64 (ZeroExt16to32 x))
	for {
		x := v.Args[0]
		v.reset(OpZeroExt32to64)
		v0 := b.NewValue0(v.Line, OpZeroExt16to32, config.fe.TypeUInt32())
		v0.AddArg(x)
		v.AddArg(v0)
		return true
	}
}
func rewriteValuedec64_OpZeroExt32to64(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (ZeroExt32to64 x)
	// cond:
	// result: (Int64Make (Const32 <config.fe.TypeUInt32()> [0]) x)
	for {
		x := v.Args[0]
		v.reset(OpInt64Make)
		v0 := b.NewValue0(v.Line, OpConst32, config.fe.TypeUInt32())
		v0.AuxInt = 0
		v.AddArg(v0)
		v.AddArg(x)
		return true
	}
}
func rewriteValuedec64_OpZeroExt8to64(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (ZeroExt8to64 x)
	// cond:
	// result: (ZeroExt32to64 (ZeroExt8to32 x))
	for {
		x := v.Args[0]
		v.reset(OpZeroExt32to64)
		v0 := b.NewValue0(v.Line, OpZeroExt8to32, config.fe.TypeUInt32())
		v0.AddArg(x)
		v.AddArg(v0)
		return true
	}
}
func rewriteBlockdec64(b *Block) bool {
	switch b.Kind {
	}
	return false
}
