// autogenerated: do not edit!
// generated from gen/*Ops.go
package ssa

const (
	blockInvalid BlockKind = iota

	BlockAMD64EQ
	BlockAMD64NE
	BlockAMD64LT
	BlockAMD64LE
	BlockAMD64GT
	BlockAMD64GE
	BlockAMD64ULT
	BlockAMD64ULE
	BlockAMD64UGT
	BlockAMD64UGE

	BlockExit
	BlockPlain
	BlockIf
	BlockCall
)

var blockString = [...]string{
	blockInvalid: "BlockInvalid",

	BlockAMD64EQ:  "EQ",
	BlockAMD64NE:  "NE",
	BlockAMD64LT:  "LT",
	BlockAMD64LE:  "LE",
	BlockAMD64GT:  "GT",
	BlockAMD64GE:  "GE",
	BlockAMD64ULT: "ULT",
	BlockAMD64ULE: "ULE",
	BlockAMD64UGT: "UGT",
	BlockAMD64UGE: "UGE",

	BlockExit:  "Exit",
	BlockPlain: "Plain",
	BlockIf:    "If",
	BlockCall:  "Call",
}

func (k BlockKind) String() string { return blockString[k] }

const (
	OpInvalid Op = iota

	OpAMD64ADDQ
	OpAMD64ADDQconst
	OpAMD64SUBQ
	OpAMD64SUBQconst
	OpAMD64MULQ
	OpAMD64MULQconst
	OpAMD64SHLQ
	OpAMD64SHLQconst
	OpAMD64NEGQ
	OpAMD64CMPQ
	OpAMD64CMPQconst
	OpAMD64TESTQ
	OpAMD64TESTB
	OpAMD64SETEQ
	OpAMD64SETNE
	OpAMD64SETL
	OpAMD64SETG
	OpAMD64SETGE
	OpAMD64SETB
	OpAMD64MOVQconst
	OpAMD64LEAQ
	OpAMD64LEAQ2
	OpAMD64LEAQ4
	OpAMD64LEAQ8
	OpAMD64LEAQglobal
	OpAMD64MOVBload
	OpAMD64MOVBQZXload
	OpAMD64MOVBQSXload
	OpAMD64MOVQload
	OpAMD64MOVQloadidx8
	OpAMD64MOVBstore
	OpAMD64MOVQstore
	OpAMD64MOVQstoreidx8
	OpAMD64MOVQloadglobal
	OpAMD64MOVQstoreglobal
	OpAMD64REPMOVSB
	OpAMD64ADDL
	OpAMD64InvertFlags

	OpAdd
	OpSub
	OpMul
	OpLsh
	OpRsh
	OpLess
	OpPhi
	OpCopy
	OpConst
	OpArg
	OpGlobal
	OpSP
	OpFP
	OpFunc
	OpLoad
	OpStore
	OpMove
	OpCall
	OpStaticCall
	OpConvert
	OpConvNop
	OpIsNonNil
	OpIsInBounds
	OpArrayIndex
	OpPtrIndex
	OpOffPtr
	OpSliceMake
	OpSlicePtr
	OpSliceLen
	OpSliceCap
	OpStringMake
	OpStringPtr
	OpStringLen
	OpStoreReg8
	OpLoadReg8
	OpFwdRef
)

var opcodeTable = [...]opInfo{
	{name: "OpInvalid"},

	{
		name: "ADDQ",
		reg: regInfo{
			inputs: []regMask{
				4295032831,
				4295032831,
			},
			clobbers: 0,
			outputs: []regMask{
				65519,
			},
		},
	},
	{
		name: "ADDQconst",
		reg: regInfo{
			inputs: []regMask{
				4295032831,
			},
			clobbers: 0,
			outputs: []regMask{
				65519,
			},
		},
	},
	{
		name: "SUBQ",
		reg: regInfo{
			inputs: []regMask{
				4295032831,
				4295032831,
			},
			clobbers: 0,
			outputs: []regMask{
				65519,
			},
		},
	},
	{
		name: "SUBQconst",
		reg: regInfo{
			inputs: []regMask{
				4295032831,
			},
			clobbers: 0,
			outputs: []regMask{
				65519,
			},
		},
	},
	{
		name: "MULQ",
		reg: regInfo{
			inputs: []regMask{
				4295032831,
				4295032831,
			},
			clobbers: 0,
			outputs: []regMask{
				65519,
			},
		},
	},
	{
		name: "MULQconst",
		reg: regInfo{
			inputs: []regMask{
				4295032831,
			},
			clobbers: 0,
			outputs: []regMask{
				65519,
			},
		},
	},
	{
		name: "SHLQ",
		reg: regInfo{
			inputs: []regMask{
				4295032831,
				2,
			},
			clobbers: 0,
			outputs: []regMask{
				65519,
			},
		},
	},
	{
		name: "SHLQconst",
		reg: regInfo{
			inputs: []regMask{
				4295032831,
			},
			clobbers: 0,
			outputs: []regMask{
				65519,
			},
		},
	},
	{
		name: "NEGQ",
		reg: regInfo{
			inputs: []regMask{
				4295032831,
			},
			clobbers: 0,
			outputs: []regMask{
				65519,
			},
		},
	},
	{
		name: "CMPQ",
		reg: regInfo{
			inputs: []regMask{
				4295032831,
				4295032831,
			},
			clobbers: 0,
			outputs: []regMask{
				8589934592,
			},
		},
	},
	{
		name: "CMPQconst",
		reg: regInfo{
			inputs: []regMask{
				4295032831,
			},
			clobbers: 0,
			outputs: []regMask{
				8589934592,
			},
		},
	},
	{
		name: "TESTQ",
		reg: regInfo{
			inputs: []regMask{
				4295032831,
				4295032831,
			},
			clobbers: 0,
			outputs: []regMask{
				8589934592,
			},
		},
	},
	{
		name: "TESTB",
		reg: regInfo{
			inputs: []regMask{
				4295032831,
				4295032831,
			},
			clobbers: 0,
			outputs: []regMask{
				8589934592,
			},
		},
	},
	{
		name: "SETEQ",
		reg: regInfo{
			inputs: []regMask{
				8589934592,
			},
			clobbers: 0,
			outputs: []regMask{
				65519,
			},
		},
	},
	{
		name: "SETNE",
		reg: regInfo{
			inputs: []regMask{
				8589934592,
			},
			clobbers: 0,
			outputs: []regMask{
				65519,
			},
		},
	},
	{
		name: "SETL",
		reg: regInfo{
			inputs: []regMask{
				8589934592,
			},
			clobbers: 0,
			outputs: []regMask{
				65519,
			},
		},
	},
	{
		name: "SETG",
		reg: regInfo{
			inputs: []regMask{
				8589934592,
			},
			clobbers: 0,
			outputs: []regMask{
				65519,
			},
		},
	},
	{
		name: "SETGE",
		reg: regInfo{
			inputs: []regMask{
				8589934592,
			},
			clobbers: 0,
			outputs: []regMask{
				65519,
			},
		},
	},
	{
		name: "SETB",
		reg: regInfo{
			inputs: []regMask{
				8589934592,
			},
			clobbers: 0,
			outputs: []regMask{
				65519,
			},
		},
	},
	{
		name: "MOVQconst",
		reg: regInfo{
			inputs:   []regMask{},
			clobbers: 0,
			outputs: []regMask{
				65519,
			},
		},
	},
	{
		name: "LEAQ",
		reg: regInfo{
			inputs: []regMask{
				4295032831,
				4295032831,
			},
			clobbers: 0,
			outputs: []regMask{
				65519,
			},
		},
	},
	{
		name: "LEAQ2",
		reg: regInfo{
			inputs: []regMask{
				4295032831,
				4295032831,
			},
			clobbers: 0,
			outputs: []regMask{
				65519,
			},
		},
	},
	{
		name: "LEAQ4",
		reg: regInfo{
			inputs: []regMask{
				4295032831,
				4295032831,
			},
			clobbers: 0,
			outputs: []regMask{
				65519,
			},
		},
	},
	{
		name: "LEAQ8",
		reg: regInfo{
			inputs: []regMask{
				4295032831,
				4295032831,
			},
			clobbers: 0,
			outputs: []regMask{
				65519,
			},
		},
	},
	{
		name: "LEAQglobal",
		reg: regInfo{
			inputs:   []regMask{},
			clobbers: 0,
			outputs: []regMask{
				65519,
			},
		},
	},
	{
		name: "MOVBload",
		reg: regInfo{
			inputs: []regMask{
				4295032831,
				0,
			},
			clobbers: 0,
			outputs: []regMask{
				65519,
			},
		},
	},
	{
		name: "MOVBQZXload",
		reg: regInfo{
			inputs: []regMask{
				4295032831,
				0,
			},
			clobbers: 0,
			outputs: []regMask{
				65519,
			},
		},
	},
	{
		name: "MOVBQSXload",
		reg: regInfo{
			inputs: []regMask{
				4295032831,
				0,
			},
			clobbers: 0,
			outputs: []regMask{
				65519,
			},
		},
	},
	{
		name: "MOVQload",
		reg: regInfo{
			inputs: []regMask{
				4295032831,
				0,
			},
			clobbers: 0,
			outputs: []regMask{
				65519,
			},
		},
	},
	{
		name: "MOVQloadidx8",
		reg: regInfo{
			inputs: []regMask{
				4295032831,
				4295032831,
				0,
			},
			clobbers: 0,
			outputs: []regMask{
				65519,
			},
		},
	},
	{
		name: "MOVBstore",
		reg: regInfo{
			inputs: []regMask{
				4295032831,
				4295032831,
				0,
			},
			clobbers: 0,
			outputs:  []regMask{},
		},
	},
	{
		name: "MOVQstore",
		reg: regInfo{
			inputs: []regMask{
				4295032831,
				4295032831,
				0,
			},
			clobbers: 0,
			outputs:  []regMask{},
		},
	},
	{
		name: "MOVQstoreidx8",
		reg: regInfo{
			inputs: []regMask{
				4295032831,
				4295032831,
				4295032831,
				0,
			},
			clobbers: 0,
			outputs:  []regMask{},
		},
	},
	{
		name: "MOVQloadglobal",
		reg: regInfo{
			inputs:   []regMask{},
			clobbers: 0,
			outputs:  []regMask{},
		},
	},
	{
		name: "MOVQstoreglobal",
		reg: regInfo{
			inputs:   []regMask{},
			clobbers: 0,
			outputs:  []regMask{},
		},
	},
	{
		name: "REPMOVSB",
		reg: regInfo{
			inputs: []regMask{
				128,
				64,
				2,
			},
			clobbers: 194,
			outputs:  []regMask{},
		},
	},
	{
		name: "ADDL",
		reg: regInfo{
			inputs: []regMask{
				4295032831,
				4295032831,
			},
			clobbers: 0,
			outputs: []regMask{
				65519,
			},
		},
	},
	{
		name: "InvertFlags",
		reg: regInfo{
			inputs:   []regMask{},
			clobbers: 0,
			outputs:  []regMask{},
		},
	},

	{
		name: "Add",
		reg: regInfo{
			inputs:   []regMask{},
			clobbers: 0,
			outputs:  []regMask{},
		},
		generic: true,
	},
	{
		name: "Sub",
		reg: regInfo{
			inputs:   []regMask{},
			clobbers: 0,
			outputs:  []regMask{},
		},
		generic: true,
	},
	{
		name: "Mul",
		reg: regInfo{
			inputs:   []regMask{},
			clobbers: 0,
			outputs:  []regMask{},
		},
		generic: true,
	},
	{
		name: "Lsh",
		reg: regInfo{
			inputs:   []regMask{},
			clobbers: 0,
			outputs:  []regMask{},
		},
		generic: true,
	},
	{
		name: "Rsh",
		reg: regInfo{
			inputs:   []regMask{},
			clobbers: 0,
			outputs:  []regMask{},
		},
		generic: true,
	},
	{
		name: "Less",
		reg: regInfo{
			inputs:   []regMask{},
			clobbers: 0,
			outputs:  []regMask{},
		},
		generic: true,
	},
	{
		name: "Phi",
		reg: regInfo{
			inputs:   []regMask{},
			clobbers: 0,
			outputs:  []regMask{},
		},
		generic: true,
	},
	{
		name: "Copy",
		reg: regInfo{
			inputs:   []regMask{},
			clobbers: 0,
			outputs:  []regMask{},
		},
		generic: true,
	},
	{
		name: "Const",
		reg: regInfo{
			inputs:   []regMask{},
			clobbers: 0,
			outputs:  []regMask{},
		},
		generic: true,
	},
	{
		name: "Arg",
		reg: regInfo{
			inputs:   []regMask{},
			clobbers: 0,
			outputs:  []regMask{},
		},
		generic: true,
	},
	{
		name: "Global",
		reg: regInfo{
			inputs:   []regMask{},
			clobbers: 0,
			outputs:  []regMask{},
		},
		generic: true,
	},
	{
		name: "SP",
		reg: regInfo{
			inputs:   []regMask{},
			clobbers: 0,
			outputs:  []regMask{},
		},
		generic: true,
	},
	{
		name: "FP",
		reg: regInfo{
			inputs:   []regMask{},
			clobbers: 0,
			outputs:  []regMask{},
		},
		generic: true,
	},
	{
		name: "Func",
		reg: regInfo{
			inputs:   []regMask{},
			clobbers: 0,
			outputs:  []regMask{},
		},
		generic: true,
	},
	{
		name: "Load",
		reg: regInfo{
			inputs:   []regMask{},
			clobbers: 0,
			outputs:  []regMask{},
		},
		generic: true,
	},
	{
		name: "Store",
		reg: regInfo{
			inputs:   []regMask{},
			clobbers: 0,
			outputs:  []regMask{},
		},
		generic: true,
	},
	{
		name: "Move",
		reg: regInfo{
			inputs:   []regMask{},
			clobbers: 0,
			outputs:  []regMask{},
		},
		generic: true,
	},
	{
		name: "Call",
		reg: regInfo{
			inputs:   []regMask{},
			clobbers: 0,
			outputs:  []regMask{},
		},
		generic: true,
	},
	{
		name: "StaticCall",
		reg: regInfo{
			inputs:   []regMask{},
			clobbers: 0,
			outputs:  []regMask{},
		},
		generic: true,
	},
	{
		name: "Convert",
		reg: regInfo{
			inputs:   []regMask{},
			clobbers: 0,
			outputs:  []regMask{},
		},
		generic: true,
	},
	{
		name: "ConvNop",
		reg: regInfo{
			inputs:   []regMask{},
			clobbers: 0,
			outputs:  []regMask{},
		},
		generic: true,
	},
	{
		name: "IsNonNil",
		reg: regInfo{
			inputs:   []regMask{},
			clobbers: 0,
			outputs:  []regMask{},
		},
		generic: true,
	},
	{
		name: "IsInBounds",
		reg: regInfo{
			inputs:   []regMask{},
			clobbers: 0,
			outputs:  []regMask{},
		},
		generic: true,
	},
	{
		name: "ArrayIndex",
		reg: regInfo{
			inputs:   []regMask{},
			clobbers: 0,
			outputs:  []regMask{},
		},
		generic: true,
	},
	{
		name: "PtrIndex",
		reg: regInfo{
			inputs:   []regMask{},
			clobbers: 0,
			outputs:  []regMask{},
		},
		generic: true,
	},
	{
		name: "OffPtr",
		reg: regInfo{
			inputs:   []regMask{},
			clobbers: 0,
			outputs:  []regMask{},
		},
		generic: true,
	},
	{
		name: "SliceMake",
		reg: regInfo{
			inputs:   []regMask{},
			clobbers: 0,
			outputs:  []regMask{},
		},
		generic: true,
	},
	{
		name: "SlicePtr",
		reg: regInfo{
			inputs:   []regMask{},
			clobbers: 0,
			outputs:  []regMask{},
		},
		generic: true,
	},
	{
		name: "SliceLen",
		reg: regInfo{
			inputs:   []regMask{},
			clobbers: 0,
			outputs:  []regMask{},
		},
		generic: true,
	},
	{
		name: "SliceCap",
		reg: regInfo{
			inputs:   []regMask{},
			clobbers: 0,
			outputs:  []regMask{},
		},
		generic: true,
	},
	{
		name: "StringMake",
		reg: regInfo{
			inputs:   []regMask{},
			clobbers: 0,
			outputs:  []regMask{},
		},
		generic: true,
	},
	{
		name: "StringPtr",
		reg: regInfo{
			inputs:   []regMask{},
			clobbers: 0,
			outputs:  []regMask{},
		},
		generic: true,
	},
	{
		name: "StringLen",
		reg: regInfo{
			inputs:   []regMask{},
			clobbers: 0,
			outputs:  []regMask{},
		},
		generic: true,
	},
	{
		name: "StoreReg8",
		reg: regInfo{
			inputs:   []regMask{},
			clobbers: 0,
			outputs:  []regMask{},
		},
		generic: true,
	},
	{
		name: "LoadReg8",
		reg: regInfo{
			inputs:   []regMask{},
			clobbers: 0,
			outputs:  []regMask{},
		},
		generic: true,
	},
	{
		name: "FwdRef",
		reg: regInfo{
			inputs:   []regMask{},
			clobbers: 0,
			outputs:  []regMask{},
		},
		generic: true,
	},
}

func (o Op) String() string { return opcodeTable[o].name }
