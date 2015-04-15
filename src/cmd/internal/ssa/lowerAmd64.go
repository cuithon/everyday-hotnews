// autogenerated from rulegen/lower_amd64.rules: do not edit!
// generated with: go run rulegen/rulegen.go rulegen/lower_amd64.rules lowerAmd64 lowerAmd64.go
package ssa

func lowerAmd64(v *Value) bool {
	switch v.Op {
	case OpADDCQ:
		// match: (ADDCQ [c] (LEAQ8 [d] x y))
		// cond:
		// result: (LEAQ8 [c.(int64)+d.(int64)] x y)
		{
			c := v.Aux
			if v.Args[0].Op != OpLEAQ8 {
				goto end16348939e556e99e8447227ecb986f01
			}
			d := v.Args[0].Aux
			x := v.Args[0].Args[0]
			y := v.Args[0].Args[1]
			v.Op = OpLEAQ8
			v.Aux = nil
			v.resetArgs()
			v.Aux = c.(int64) + d.(int64)
			v.AddArg(x)
			v.AddArg(y)
			return true
		}
		goto end16348939e556e99e8447227ecb986f01
	end16348939e556e99e8447227ecb986f01:
		;
		// match: (ADDCQ [off1] (FPAddr [off2]))
		// cond:
		// result: (FPAddr [off1.(int64)+off2.(int64)])
		{
			off1 := v.Aux
			if v.Args[0].Op != OpFPAddr {
				goto end28e093ab0618066e6b2609db7aaf309b
			}
			off2 := v.Args[0].Aux
			v.Op = OpFPAddr
			v.Aux = nil
			v.resetArgs()
			v.Aux = off1.(int64) + off2.(int64)
			return true
		}
		goto end28e093ab0618066e6b2609db7aaf309b
	end28e093ab0618066e6b2609db7aaf309b:
		;
		// match: (ADDCQ [off1] (SPAddr [off2]))
		// cond:
		// result: (SPAddr [off1.(int64)+off2.(int64)])
		{
			off1 := v.Aux
			if v.Args[0].Op != OpSPAddr {
				goto endd0c27c62d150b88168075c5ba113d1fa
			}
			off2 := v.Args[0].Aux
			v.Op = OpSPAddr
			v.Aux = nil
			v.resetArgs()
			v.Aux = off1.(int64) + off2.(int64)
			return true
		}
		goto endd0c27c62d150b88168075c5ba113d1fa
	endd0c27c62d150b88168075c5ba113d1fa:
		;
	case OpADDQ:
		// match: (ADDQ x (Const [c]))
		// cond:
		// result: (ADDCQ [c] x)
		{
			x := v.Args[0]
			if v.Args[1].Op != OpConst {
				goto endef6908cfdf56e102cc327a3ddc14393d
			}
			c := v.Args[1].Aux
			v.Op = OpADDCQ
			v.Aux = nil
			v.resetArgs()
			v.Aux = c
			v.AddArg(x)
			return true
		}
		goto endef6908cfdf56e102cc327a3ddc14393d
	endef6908cfdf56e102cc327a3ddc14393d:
		;
		// match: (ADDQ (Const [c]) x)
		// cond:
		// result: (ADDCQ [c] x)
		{
			if v.Args[0].Op != OpConst {
				goto endb54a32cf3147f424f08b46db62c69b23
			}
			c := v.Args[0].Aux
			x := v.Args[1]
			v.Op = OpADDCQ
			v.Aux = nil
			v.resetArgs()
			v.Aux = c
			v.AddArg(x)
			return true
		}
		goto endb54a32cf3147f424f08b46db62c69b23
	endb54a32cf3147f424f08b46db62c69b23:
		;
		// match: (ADDQ x (SHLCQ [shift] y))
		// cond: shift.(int64) == 3
		// result: (LEAQ8 [int64(0)] x y)
		{
			x := v.Args[0]
			if v.Args[1].Op != OpSHLCQ {
				goto end7fa0d837edd248748cef516853fd9475
			}
			shift := v.Args[1].Aux
			y := v.Args[1].Args[0]
			if !(shift.(int64) == 3) {
				goto end7fa0d837edd248748cef516853fd9475
			}
			v.Op = OpLEAQ8
			v.Aux = nil
			v.resetArgs()
			v.Aux = int64(0)
			v.AddArg(x)
			v.AddArg(y)
			return true
		}
		goto end7fa0d837edd248748cef516853fd9475
	end7fa0d837edd248748cef516853fd9475:
		;
	case OpAdd:
		// match: (Add <t> x y)
		// cond: (is64BitInt(t) || isPtr(t))
		// result: (ADDQ x y)
		{
			t := v.Type
			x := v.Args[0]
			y := v.Args[1]
			if !(is64BitInt(t) || isPtr(t)) {
				goto endf031c523d7dd08e4b8e7010a94cd94c9
			}
			v.Op = OpADDQ
			v.Aux = nil
			v.resetArgs()
			v.AddArg(x)
			v.AddArg(y)
			return true
		}
		goto endf031c523d7dd08e4b8e7010a94cd94c9
	endf031c523d7dd08e4b8e7010a94cd94c9:
		;
		// match: (Add <t> x y)
		// cond: is32BitInt(t)
		// result: (ADDL x y)
		{
			t := v.Type
			x := v.Args[0]
			y := v.Args[1]
			if !(is32BitInt(t)) {
				goto end35a02a1587264e40cf1055856ff8445a
			}
			v.Op = OpADDL
			v.Aux = nil
			v.resetArgs()
			v.AddArg(x)
			v.AddArg(y)
			return true
		}
		goto end35a02a1587264e40cf1055856ff8445a
	end35a02a1587264e40cf1055856ff8445a:
		;
	case OpCMPQ:
		// match: (CMPQ x (Const [c]))
		// cond:
		// result: (CMPCQ x [c])
		{
			x := v.Args[0]
			if v.Args[1].Op != OpConst {
				goto end1770a40e4253d9f669559a360514613e
			}
			c := v.Args[1].Aux
			v.Op = OpCMPCQ
			v.Aux = nil
			v.resetArgs()
			v.AddArg(x)
			v.Aux = c
			return true
		}
		goto end1770a40e4253d9f669559a360514613e
	end1770a40e4253d9f669559a360514613e:
		;
		// match: (CMPQ (Const [c]) x)
		// cond:
		// result: (InvertFlags (CMPCQ <TypeFlags> x [c]))
		{
			if v.Args[0].Op != OpConst {
				goto enda4e64c7eaeda16c1c0db9dac409cd126
			}
			c := v.Args[0].Aux
			x := v.Args[1]
			v.Op = OpInvertFlags
			v.Aux = nil
			v.resetArgs()
			v0 := v.Block.NewValue(OpCMPCQ, TypeInvalid, nil)
			v0.Type = TypeFlags
			v0.AddArg(x)
			v0.Aux = c
			v.AddArg(v0)
			return true
		}
		goto enda4e64c7eaeda16c1c0db9dac409cd126
	enda4e64c7eaeda16c1c0db9dac409cd126:
		;
	case OpCheckBound:
		// match: (CheckBound idx len)
		// cond:
		// result: (SETB (CMPQ <TypeFlags> idx len))
		{
			idx := v.Args[0]
			len := v.Args[1]
			v.Op = OpSETB
			v.Aux = nil
			v.resetArgs()
			v0 := v.Block.NewValue(OpCMPQ, TypeInvalid, nil)
			v0.Type = TypeFlags
			v0.AddArg(idx)
			v0.AddArg(len)
			v.AddArg(v0)
			return true
		}
		goto end249426f6f996d45a62f89a591311a954
	end249426f6f996d45a62f89a591311a954:
		;
	case OpCheckNil:
		// match: (CheckNil p)
		// cond:
		// result: (SETNE (TESTQ <TypeFlags> p p))
		{
			p := v.Args[0]
			v.Op = OpSETNE
			v.Aux = nil
			v.resetArgs()
			v0 := v.Block.NewValue(OpTESTQ, TypeInvalid, nil)
			v0.Type = TypeFlags
			v0.AddArg(p)
			v0.AddArg(p)
			v.AddArg(v0)
			return true
		}
		goto end90d3057824f74ef953074e473aa0b282
	end90d3057824f74ef953074e473aa0b282:
		;
	case OpLess:
		// match: (Less x y)
		// cond: is64BitInt(v.Args[0].Type) && isSigned(v.Args[0].Type)
		// result: (SETL (CMPQ <TypeFlags> x y))
		{
			x := v.Args[0]
			y := v.Args[1]
			if !(is64BitInt(v.Args[0].Type) && isSigned(v.Args[0].Type)) {
				goto endcecf13a952d4c6c2383561c7d68a3cf9
			}
			v.Op = OpSETL
			v.Aux = nil
			v.resetArgs()
			v0 := v.Block.NewValue(OpCMPQ, TypeInvalid, nil)
			v0.Type = TypeFlags
			v0.AddArg(x)
			v0.AddArg(y)
			v.AddArg(v0)
			return true
		}
		goto endcecf13a952d4c6c2383561c7d68a3cf9
	endcecf13a952d4c6c2383561c7d68a3cf9:
		;
	case OpLoad:
		// match: (Load <t> ptr mem)
		// cond: (is64BitInt(t) || isPtr(t))
		// result: (MOVQload [int64(0)] ptr mem)
		{
			t := v.Type
			ptr := v.Args[0]
			mem := v.Args[1]
			if !(is64BitInt(t) || isPtr(t)) {
				goto end581ce5a20901df1b8143448ba031685b
			}
			v.Op = OpMOVQload
			v.Aux = nil
			v.resetArgs()
			v.Aux = int64(0)
			v.AddArg(ptr)
			v.AddArg(mem)
			return true
		}
		goto end581ce5a20901df1b8143448ba031685b
	end581ce5a20901df1b8143448ba031685b:
		;
	case OpMOVQload:
		// match: (MOVQload [off1] (FPAddr [off2]) mem)
		// cond:
		// result: (MOVQloadFP [off1.(int64)+off2.(int64)] mem)
		{
			off1 := v.Aux
			if v.Args[0].Op != OpFPAddr {
				goto endce972b1aa84b56447978c43def87fa57
			}
			off2 := v.Args[0].Aux
			mem := v.Args[1]
			v.Op = OpMOVQloadFP
			v.Aux = nil
			v.resetArgs()
			v.Aux = off1.(int64) + off2.(int64)
			v.AddArg(mem)
			return true
		}
		goto endce972b1aa84b56447978c43def87fa57
	endce972b1aa84b56447978c43def87fa57:
		;
		// match: (MOVQload [off1] (SPAddr [off2]) mem)
		// cond:
		// result: (MOVQloadSP [off1.(int64)+off2.(int64)] mem)
		{
			off1 := v.Aux
			if v.Args[0].Op != OpSPAddr {
				goto end3d8628a6536350a123be81240b8a1376
			}
			off2 := v.Args[0].Aux
			mem := v.Args[1]
			v.Op = OpMOVQloadSP
			v.Aux = nil
			v.resetArgs()
			v.Aux = off1.(int64) + off2.(int64)
			v.AddArg(mem)
			return true
		}
		goto end3d8628a6536350a123be81240b8a1376
	end3d8628a6536350a123be81240b8a1376:
		;
		// match: (MOVQload [off1] (ADDCQ [off2] ptr) mem)
		// cond:
		// result: (MOVQload [off1.(int64)+off2.(int64)] ptr mem)
		{
			off1 := v.Aux
			if v.Args[0].Op != OpADDCQ {
				goto enda68a39292ba2a05b3436191cb0bb0516
			}
			off2 := v.Args[0].Aux
			ptr := v.Args[0].Args[0]
			mem := v.Args[1]
			v.Op = OpMOVQload
			v.Aux = nil
			v.resetArgs()
			v.Aux = off1.(int64) + off2.(int64)
			v.AddArg(ptr)
			v.AddArg(mem)
			return true
		}
		goto enda68a39292ba2a05b3436191cb0bb0516
	enda68a39292ba2a05b3436191cb0bb0516:
		;
		// match: (MOVQload [off1] (LEAQ8 [off2] ptr idx) mem)
		// cond:
		// result: (MOVQload8 [off1.(int64)+off2.(int64)] ptr idx mem)
		{
			off1 := v.Aux
			if v.Args[0].Op != OpLEAQ8 {
				goto end35060118a284c93323ab3fb827156638
			}
			off2 := v.Args[0].Aux
			ptr := v.Args[0].Args[0]
			idx := v.Args[0].Args[1]
			mem := v.Args[1]
			v.Op = OpMOVQload8
			v.Aux = nil
			v.resetArgs()
			v.Aux = off1.(int64) + off2.(int64)
			v.AddArg(ptr)
			v.AddArg(idx)
			v.AddArg(mem)
			return true
		}
		goto end35060118a284c93323ab3fb827156638
	end35060118a284c93323ab3fb827156638:
		;
	case OpMOVQstore:
		// match: (MOVQstore [off1] (FPAddr [off2]) val mem)
		// cond:
		// result: (MOVQstoreFP [off1.(int64)+off2.(int64)] val mem)
		{
			off1 := v.Aux
			if v.Args[0].Op != OpFPAddr {
				goto end0a2a81a20558dfc93790aecb1e9cc81a
			}
			off2 := v.Args[0].Aux
			val := v.Args[1]
			mem := v.Args[2]
			v.Op = OpMOVQstoreFP
			v.Aux = nil
			v.resetArgs()
			v.Aux = off1.(int64) + off2.(int64)
			v.AddArg(val)
			v.AddArg(mem)
			return true
		}
		goto end0a2a81a20558dfc93790aecb1e9cc81a
	end0a2a81a20558dfc93790aecb1e9cc81a:
		;
		// match: (MOVQstore [off1] (SPAddr [off2]) val mem)
		// cond:
		// result: (MOVQstoreSP [off1.(int64)+off2.(int64)] val mem)
		{
			off1 := v.Aux
			if v.Args[0].Op != OpSPAddr {
				goto end1cb5b7e766f018270fa434c6f46f607f
			}
			off2 := v.Args[0].Aux
			val := v.Args[1]
			mem := v.Args[2]
			v.Op = OpMOVQstoreSP
			v.Aux = nil
			v.resetArgs()
			v.Aux = off1.(int64) + off2.(int64)
			v.AddArg(val)
			v.AddArg(mem)
			return true
		}
		goto end1cb5b7e766f018270fa434c6f46f607f
	end1cb5b7e766f018270fa434c6f46f607f:
		;
		// match: (MOVQstore [off1] (ADDCQ [off2] ptr) val mem)
		// cond:
		// result: (MOVQstore [off1.(int64)+off2.(int64)] ptr val mem)
		{
			off1 := v.Aux
			if v.Args[0].Op != OpADDCQ {
				goto end271e3052de832e22b1f07576af2854de
			}
			off2 := v.Args[0].Aux
			ptr := v.Args[0].Args[0]
			val := v.Args[1]
			mem := v.Args[2]
			v.Op = OpMOVQstore
			v.Aux = nil
			v.resetArgs()
			v.Aux = off1.(int64) + off2.(int64)
			v.AddArg(ptr)
			v.AddArg(val)
			v.AddArg(mem)
			return true
		}
		goto end271e3052de832e22b1f07576af2854de
	end271e3052de832e22b1f07576af2854de:
		;
		// match: (MOVQstore [off1] (LEAQ8 [off2] ptr idx) val mem)
		// cond:
		// result: (MOVQstore8 [off1.(int64)+off2.(int64)] ptr idx val mem)
		{
			off1 := v.Aux
			if v.Args[0].Op != OpLEAQ8 {
				goto endb5cba0ee3ba21d2bd8e5aa163d2b984e
			}
			off2 := v.Args[0].Aux
			ptr := v.Args[0].Args[0]
			idx := v.Args[0].Args[1]
			val := v.Args[1]
			mem := v.Args[2]
			v.Op = OpMOVQstore8
			v.Aux = nil
			v.resetArgs()
			v.Aux = off1.(int64) + off2.(int64)
			v.AddArg(ptr)
			v.AddArg(idx)
			v.AddArg(val)
			v.AddArg(mem)
			return true
		}
		goto endb5cba0ee3ba21d2bd8e5aa163d2b984e
	endb5cba0ee3ba21d2bd8e5aa163d2b984e:
		;
	case OpMULCQ:
		// match: (MULCQ [c] x)
		// cond: c.(int64) == 8
		// result: (SHLCQ [int64(3)] x)
		{
			c := v.Aux
			x := v.Args[0]
			if !(c.(int64) == 8) {
				goto end90a1c055d9658aecacce5e101c1848b4
			}
			v.Op = OpSHLCQ
			v.Aux = nil
			v.resetArgs()
			v.Aux = int64(3)
			v.AddArg(x)
			return true
		}
		goto end90a1c055d9658aecacce5e101c1848b4
	end90a1c055d9658aecacce5e101c1848b4:
		;
	case OpMULQ:
		// match: (MULQ x (Const [c]))
		// cond:
		// result: (MULCQ [c] x)
		{
			x := v.Args[0]
			if v.Args[1].Op != OpConst {
				goto endc427f4838d2e83c00cc097b20bd20a37
			}
			c := v.Args[1].Aux
			v.Op = OpMULCQ
			v.Aux = nil
			v.resetArgs()
			v.Aux = c
			v.AddArg(x)
			return true
		}
		goto endc427f4838d2e83c00cc097b20bd20a37
	endc427f4838d2e83c00cc097b20bd20a37:
		;
		// match: (MULQ (Const [c]) x)
		// cond:
		// result: (MULCQ [c] x)
		{
			if v.Args[0].Op != OpConst {
				goto endd70de938e71150d1c9e8173c2a5b2d95
			}
			c := v.Args[0].Aux
			x := v.Args[1]
			v.Op = OpMULCQ
			v.Aux = nil
			v.resetArgs()
			v.Aux = c
			v.AddArg(x)
			return true
		}
		goto endd70de938e71150d1c9e8173c2a5b2d95
	endd70de938e71150d1c9e8173c2a5b2d95:
		;
	case OpMul:
		// match: (Mul <t> x y)
		// cond: is64BitInt(t)
		// result: (MULQ x y)
		{
			t := v.Type
			x := v.Args[0]
			y := v.Args[1]
			if !(is64BitInt(t)) {
				goto endfab0d598f376ecba45a22587d50f7aff
			}
			v.Op = OpMULQ
			v.Aux = nil
			v.resetArgs()
			v.AddArg(x)
			v.AddArg(y)
			return true
		}
		goto endfab0d598f376ecba45a22587d50f7aff
	endfab0d598f376ecba45a22587d50f7aff:
		;
	case OpSETL:
		// match: (SETL (InvertFlags x))
		// cond:
		// result: (SETGE x)
		{
			if v.Args[0].Op != OpInvertFlags {
				goto end456c7681d48305698c1ef462d244bdc6
			}
			x := v.Args[0].Args[0]
			v.Op = OpSETGE
			v.Aux = nil
			v.resetArgs()
			v.AddArg(x)
			return true
		}
		goto end456c7681d48305698c1ef462d244bdc6
	end456c7681d48305698c1ef462d244bdc6:
		;
	case OpSUBQ:
		// match: (SUBQ x (Const [c]))
		// cond:
		// result: (SUBCQ x [c])
		{
			x := v.Args[0]
			if v.Args[1].Op != OpConst {
				goto endb31e242f283867de4722665a5796008c
			}
			c := v.Args[1].Aux
			v.Op = OpSUBCQ
			v.Aux = nil
			v.resetArgs()
			v.AddArg(x)
			v.Aux = c
			return true
		}
		goto endb31e242f283867de4722665a5796008c
	endb31e242f283867de4722665a5796008c:
		;
		// match: (SUBQ <t> (Const [c]) x)
		// cond:
		// result: (NEGQ (SUBCQ <t> x [c]))
		{
			t := v.Type
			if v.Args[0].Op != OpConst {
				goto end569cc755877d1f89a701378bec05c08d
			}
			c := v.Args[0].Aux
			x := v.Args[1]
			v.Op = OpNEGQ
			v.Aux = nil
			v.resetArgs()
			v0 := v.Block.NewValue(OpSUBCQ, TypeInvalid, nil)
			v0.Type = t
			v0.AddArg(x)
			v0.Aux = c
			v.AddArg(v0)
			return true
		}
		goto end569cc755877d1f89a701378bec05c08d
	end569cc755877d1f89a701378bec05c08d:
		;
	case OpStore:
		// match: (Store ptr val mem)
		// cond: (is64BitInt(val.Type) || isPtr(val.Type))
		// result: (MOVQstore [int64(0)] ptr val mem)
		{
			ptr := v.Args[0]
			val := v.Args[1]
			mem := v.Args[2]
			if !(is64BitInt(val.Type) || isPtr(val.Type)) {
				goto end9680b43f504bc06f9fab000823ce471a
			}
			v.Op = OpMOVQstore
			v.Aux = nil
			v.resetArgs()
			v.Aux = int64(0)
			v.AddArg(ptr)
			v.AddArg(val)
			v.AddArg(mem)
			return true
		}
		goto end9680b43f504bc06f9fab000823ce471a
	end9680b43f504bc06f9fab000823ce471a:
		;
	case OpSub:
		// match: (Sub <t> x y)
		// cond: is64BitInt(t)
		// result: (SUBQ x y)
		{
			t := v.Type
			x := v.Args[0]
			y := v.Args[1]
			if !(is64BitInt(t)) {
				goto ende6ef29f885a8ecf3058212bb95917323
			}
			v.Op = OpSUBQ
			v.Aux = nil
			v.resetArgs()
			v.AddArg(x)
			v.AddArg(y)
			return true
		}
		goto ende6ef29f885a8ecf3058212bb95917323
	ende6ef29f885a8ecf3058212bb95917323:
	}
	return false
}
