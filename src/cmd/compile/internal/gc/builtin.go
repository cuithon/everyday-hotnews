// Code generated by mkbuiltin.go. DO NOT EDIT.

package gc

import "cmd/compile/internal/types"

var runtimeDecls = [...]struct {
	name string
	tag  int
	typ  int
}{
	{"newobject", funcTag, 4},
	{"panicindex", funcTag, 5},
	{"panicslice", funcTag, 5},
	{"panicdivide", funcTag, 5},
	{"throwinit", funcTag, 5},
	{"panicwrap", funcTag, 5},
	{"gopanic", funcTag, 7},
	{"gorecover", funcTag, 10},
	{"goschedguarded", funcTag, 5},
	{"printbool", funcTag, 12},
	{"printfloat", funcTag, 14},
	{"printint", funcTag, 16},
	{"printhex", funcTag, 18},
	{"printuint", funcTag, 18},
	{"printcomplex", funcTag, 20},
	{"printstring", funcTag, 22},
	{"printpointer", funcTag, 23},
	{"printiface", funcTag, 23},
	{"printeface", funcTag, 23},
	{"printslice", funcTag, 23},
	{"printnl", funcTag, 5},
	{"printsp", funcTag, 5},
	{"printlock", funcTag, 5},
	{"printunlock", funcTag, 5},
	{"concatstring2", funcTag, 26},
	{"concatstring3", funcTag, 27},
	{"concatstring4", funcTag, 28},
	{"concatstring5", funcTag, 29},
	{"concatstrings", funcTag, 31},
	{"cmpstring", funcTag, 33},
	{"eqstring", funcTag, 34},
	{"intstring", funcTag, 37},
	{"slicebytetostring", funcTag, 39},
	{"slicebytetostringtmp", funcTag, 40},
	{"slicerunetostring", funcTag, 43},
	{"stringtoslicebyte", funcTag, 44},
	{"stringtoslicerune", funcTag, 47},
	{"decoderune", funcTag, 48},
	{"slicecopy", funcTag, 50},
	{"slicestringcopy", funcTag, 51},
	{"convI2I", funcTag, 52},
	{"convT2E", funcTag, 53},
	{"convT2E16", funcTag, 53},
	{"convT2E32", funcTag, 53},
	{"convT2E64", funcTag, 53},
	{"convT2Estring", funcTag, 53},
	{"convT2Eslice", funcTag, 53},
	{"convT2Enoptr", funcTag, 53},
	{"convT2I", funcTag, 53},
	{"convT2I16", funcTag, 53},
	{"convT2I32", funcTag, 53},
	{"convT2I64", funcTag, 53},
	{"convT2Istring", funcTag, 53},
	{"convT2Islice", funcTag, 53},
	{"convT2Inoptr", funcTag, 53},
	{"assertE2I", funcTag, 52},
	{"assertE2I2", funcTag, 54},
	{"assertI2I", funcTag, 52},
	{"assertI2I2", funcTag, 54},
	{"panicdottypeE", funcTag, 55},
	{"panicdottypeI", funcTag, 55},
	{"panicnildottype", funcTag, 56},
	{"ifaceeq", funcTag, 59},
	{"efaceeq", funcTag, 59},
	{"makemap", funcTag, 61},
	{"mapaccess1", funcTag, 62},
	{"mapaccess1_fast32", funcTag, 63},
	{"mapaccess1_fast64", funcTag, 63},
	{"mapaccess1_faststr", funcTag, 63},
	{"mapaccess1_fat", funcTag, 64},
	{"mapaccess2", funcTag, 65},
	{"mapaccess2_fast32", funcTag, 66},
	{"mapaccess2_fast64", funcTag, 66},
	{"mapaccess2_faststr", funcTag, 66},
	{"mapaccess2_fat", funcTag, 67},
	{"mapassign", funcTag, 62},
	{"mapassign_fast32", funcTag, 63},
	{"mapassign_fast64", funcTag, 63},
	{"mapassign_faststr", funcTag, 63},
	{"mapiterinit", funcTag, 68},
	{"mapdelete", funcTag, 68},
	{"mapdelete_fast32", funcTag, 69},
	{"mapdelete_fast64", funcTag, 69},
	{"mapdelete_faststr", funcTag, 69},
	{"mapiternext", funcTag, 70},
	{"makechan64", funcTag, 72},
	{"makechan", funcTag, 73},
	{"chanrecv1", funcTag, 75},
	{"chanrecv2", funcTag, 76},
	{"chansend1", funcTag, 78},
	{"closechan", funcTag, 23},
	{"writeBarrier", varTag, 80},
	{"writebarrierptr", funcTag, 81},
	{"typedmemmove", funcTag, 82},
	{"typedmemclr", funcTag, 83},
	{"typedslicecopy", funcTag, 84},
	{"selectnbsend", funcTag, 85},
	{"selectnbrecv", funcTag, 86},
	{"selectnbrecv2", funcTag, 88},
	{"newselect", funcTag, 89},
	{"selectsend", funcTag, 90},
	{"selectrecv", funcTag, 91},
	{"selectdefault", funcTag, 56},
	{"selectgo", funcTag, 92},
	{"block", funcTag, 5},
	{"makeslice", funcTag, 94},
	{"makeslice64", funcTag, 95},
	{"growslice", funcTag, 96},
	{"memmove", funcTag, 97},
	{"memclrNoHeapPointers", funcTag, 98},
	{"memclrHasPointers", funcTag, 98},
	{"memequal", funcTag, 99},
	{"memequal8", funcTag, 100},
	{"memequal16", funcTag, 100},
	{"memequal32", funcTag, 100},
	{"memequal64", funcTag, 100},
	{"memequal128", funcTag, 100},
	{"int64div", funcTag, 101},
	{"uint64div", funcTag, 102},
	{"int64mod", funcTag, 101},
	{"uint64mod", funcTag, 102},
	{"float64toint64", funcTag, 103},
	{"float64touint64", funcTag, 104},
	{"float64touint32", funcTag, 106},
	{"int64tofloat64", funcTag, 107},
	{"uint64tofloat64", funcTag, 108},
	{"uint32tofloat64", funcTag, 109},
	{"complex128div", funcTag, 110},
	{"racefuncenter", funcTag, 111},
	{"racefuncexit", funcTag, 5},
	{"raceread", funcTag, 111},
	{"racewrite", funcTag, 111},
	{"racereadrange", funcTag, 112},
	{"racewriterange", funcTag, 112},
	{"msanread", funcTag, 112},
	{"msanwrite", funcTag, 112},
	{"support_popcnt", varTag, 11},
}

func runtimeTypes() []*types.Type {
	var typs [113]*types.Type
	typs[0] = types.Bytetype
	typs[1] = types.NewPtr(typs[0])
	typs[2] = types.Types[TANY]
	typs[3] = types.NewPtr(typs[2])
	typs[4] = functype(nil, []*Node{anonfield(typs[1])}, []*Node{anonfield(typs[3])})
	typs[5] = functype(nil, nil, nil)
	typs[6] = types.Types[TINTER]
	typs[7] = functype(nil, []*Node{anonfield(typs[6])}, nil)
	typs[8] = types.Types[TINT32]
	typs[9] = types.NewPtr(typs[8])
	typs[10] = functype(nil, []*Node{anonfield(typs[9])}, []*Node{anonfield(typs[6])})
	typs[11] = types.Types[TBOOL]
	typs[12] = functype(nil, []*Node{anonfield(typs[11])}, nil)
	typs[13] = types.Types[TFLOAT64]
	typs[14] = functype(nil, []*Node{anonfield(typs[13])}, nil)
	typs[15] = types.Types[TINT64]
	typs[16] = functype(nil, []*Node{anonfield(typs[15])}, nil)
	typs[17] = types.Types[TUINT64]
	typs[18] = functype(nil, []*Node{anonfield(typs[17])}, nil)
	typs[19] = types.Types[TCOMPLEX128]
	typs[20] = functype(nil, []*Node{anonfield(typs[19])}, nil)
	typs[21] = types.Types[TSTRING]
	typs[22] = functype(nil, []*Node{anonfield(typs[21])}, nil)
	typs[23] = functype(nil, []*Node{anonfield(typs[2])}, nil)
	typs[24] = types.NewArray(typs[0], 32)
	typs[25] = types.NewPtr(typs[24])
	typs[26] = functype(nil, []*Node{anonfield(typs[25]), anonfield(typs[21]), anonfield(typs[21])}, []*Node{anonfield(typs[21])})
	typs[27] = functype(nil, []*Node{anonfield(typs[25]), anonfield(typs[21]), anonfield(typs[21]), anonfield(typs[21])}, []*Node{anonfield(typs[21])})
	typs[28] = functype(nil, []*Node{anonfield(typs[25]), anonfield(typs[21]), anonfield(typs[21]), anonfield(typs[21]), anonfield(typs[21])}, []*Node{anonfield(typs[21])})
	typs[29] = functype(nil, []*Node{anonfield(typs[25]), anonfield(typs[21]), anonfield(typs[21]), anonfield(typs[21]), anonfield(typs[21]), anonfield(typs[21])}, []*Node{anonfield(typs[21])})
	typs[30] = types.NewSlice(typs[21])
	typs[31] = functype(nil, []*Node{anonfield(typs[25]), anonfield(typs[30])}, []*Node{anonfield(typs[21])})
	typs[32] = types.Types[TINT]
	typs[33] = functype(nil, []*Node{anonfield(typs[21]), anonfield(typs[21])}, []*Node{anonfield(typs[32])})
	typs[34] = functype(nil, []*Node{anonfield(typs[21]), anonfield(typs[21])}, []*Node{anonfield(typs[11])})
	typs[35] = types.NewArray(typs[0], 4)
	typs[36] = types.NewPtr(typs[35])
	typs[37] = functype(nil, []*Node{anonfield(typs[36]), anonfield(typs[15])}, []*Node{anonfield(typs[21])})
	typs[38] = types.NewSlice(typs[0])
	typs[39] = functype(nil, []*Node{anonfield(typs[25]), anonfield(typs[38])}, []*Node{anonfield(typs[21])})
	typs[40] = functype(nil, []*Node{anonfield(typs[38])}, []*Node{anonfield(typs[21])})
	typs[41] = types.Runetype
	typs[42] = types.NewSlice(typs[41])
	typs[43] = functype(nil, []*Node{anonfield(typs[25]), anonfield(typs[42])}, []*Node{anonfield(typs[21])})
	typs[44] = functype(nil, []*Node{anonfield(typs[25]), anonfield(typs[21])}, []*Node{anonfield(typs[38])})
	typs[45] = types.NewArray(typs[41], 32)
	typs[46] = types.NewPtr(typs[45])
	typs[47] = functype(nil, []*Node{anonfield(typs[46]), anonfield(typs[21])}, []*Node{anonfield(typs[42])})
	typs[48] = functype(nil, []*Node{anonfield(typs[21]), anonfield(typs[32])}, []*Node{anonfield(typs[41]), anonfield(typs[32])})
	typs[49] = types.Types[TUINTPTR]
	typs[50] = functype(nil, []*Node{anonfield(typs[2]), anonfield(typs[2]), anonfield(typs[49])}, []*Node{anonfield(typs[32])})
	typs[51] = functype(nil, []*Node{anonfield(typs[2]), anonfield(typs[2])}, []*Node{anonfield(typs[32])})
	typs[52] = functype(nil, []*Node{anonfield(typs[1]), anonfield(typs[2])}, []*Node{anonfield(typs[2])})
	typs[53] = functype(nil, []*Node{anonfield(typs[1]), anonfield(typs[3])}, []*Node{anonfield(typs[2])})
	typs[54] = functype(nil, []*Node{anonfield(typs[1]), anonfield(typs[2])}, []*Node{anonfield(typs[2]), anonfield(typs[11])})
	typs[55] = functype(nil, []*Node{anonfield(typs[1]), anonfield(typs[1]), anonfield(typs[1])}, nil)
	typs[56] = functype(nil, []*Node{anonfield(typs[1])}, nil)
	typs[57] = types.NewPtr(typs[49])
	typs[58] = types.Types[TUNSAFEPTR]
	typs[59] = functype(nil, []*Node{anonfield(typs[57]), anonfield(typs[58]), anonfield(typs[58])}, []*Node{anonfield(typs[11])})
	typs[60] = types.NewMap(typs[2], typs[2])
	typs[61] = functype(nil, []*Node{anonfield(typs[1]), anonfield(typs[15]), anonfield(typs[3]), anonfield(typs[3])}, []*Node{anonfield(typs[60])})
	typs[62] = functype(nil, []*Node{anonfield(typs[1]), anonfield(typs[60]), anonfield(typs[3])}, []*Node{anonfield(typs[3])})
	typs[63] = functype(nil, []*Node{anonfield(typs[1]), anonfield(typs[60]), anonfield(typs[2])}, []*Node{anonfield(typs[3])})
	typs[64] = functype(nil, []*Node{anonfield(typs[1]), anonfield(typs[60]), anonfield(typs[3]), anonfield(typs[1])}, []*Node{anonfield(typs[3])})
	typs[65] = functype(nil, []*Node{anonfield(typs[1]), anonfield(typs[60]), anonfield(typs[3])}, []*Node{anonfield(typs[3]), anonfield(typs[11])})
	typs[66] = functype(nil, []*Node{anonfield(typs[1]), anonfield(typs[60]), anonfield(typs[2])}, []*Node{anonfield(typs[3]), anonfield(typs[11])})
	typs[67] = functype(nil, []*Node{anonfield(typs[1]), anonfield(typs[60]), anonfield(typs[3]), anonfield(typs[1])}, []*Node{anonfield(typs[3]), anonfield(typs[11])})
	typs[68] = functype(nil, []*Node{anonfield(typs[1]), anonfield(typs[60]), anonfield(typs[3])}, nil)
	typs[69] = functype(nil, []*Node{anonfield(typs[1]), anonfield(typs[60]), anonfield(typs[2])}, nil)
	typs[70] = functype(nil, []*Node{anonfield(typs[3])}, nil)
	typs[71] = types.NewChan(typs[2], types.Cboth)
	typs[72] = functype(nil, []*Node{anonfield(typs[1]), anonfield(typs[15])}, []*Node{anonfield(typs[71])})
	typs[73] = functype(nil, []*Node{anonfield(typs[1]), anonfield(typs[32])}, []*Node{anonfield(typs[71])})
	typs[74] = types.NewChan(typs[2], types.Crecv)
	typs[75] = functype(nil, []*Node{anonfield(typs[74]), anonfield(typs[3])}, nil)
	typs[76] = functype(nil, []*Node{anonfield(typs[74]), anonfield(typs[3])}, []*Node{anonfield(typs[11])})
	typs[77] = types.NewChan(typs[2], types.Csend)
	typs[78] = functype(nil, []*Node{anonfield(typs[77]), anonfield(typs[3])}, nil)
	typs[79] = types.NewArray(typs[0], 3)
	typs[80] = tostruct([]*Node{namedfield("enabled", typs[11]), namedfield("pad", typs[79]), namedfield("needed", typs[11]), namedfield("cgo", typs[11]), namedfield("alignme", typs[17])})
	typs[81] = functype(nil, []*Node{anonfield(typs[3]), anonfield(typs[2])}, nil)
	typs[82] = functype(nil, []*Node{anonfield(typs[1]), anonfield(typs[3]), anonfield(typs[3])}, nil)
	typs[83] = functype(nil, []*Node{anonfield(typs[1]), anonfield(typs[3])}, nil)
	typs[84] = functype(nil, []*Node{anonfield(typs[1]), anonfield(typs[2]), anonfield(typs[2])}, []*Node{anonfield(typs[32])})
	typs[85] = functype(nil, []*Node{anonfield(typs[77]), anonfield(typs[3])}, []*Node{anonfield(typs[11])})
	typs[86] = functype(nil, []*Node{anonfield(typs[3]), anonfield(typs[74])}, []*Node{anonfield(typs[11])})
	typs[87] = types.NewPtr(typs[11])
	typs[88] = functype(nil, []*Node{anonfield(typs[3]), anonfield(typs[87]), anonfield(typs[74])}, []*Node{anonfield(typs[11])})
	typs[89] = functype(nil, []*Node{anonfield(typs[1]), anonfield(typs[15]), anonfield(typs[8])}, nil)
	typs[90] = functype(nil, []*Node{anonfield(typs[1]), anonfield(typs[77]), anonfield(typs[3])}, nil)
	typs[91] = functype(nil, []*Node{anonfield(typs[1]), anonfield(typs[74]), anonfield(typs[3]), anonfield(typs[87])}, nil)
	typs[92] = functype(nil, []*Node{anonfield(typs[1])}, []*Node{anonfield(typs[32])})
	typs[93] = types.NewSlice(typs[2])
	typs[94] = functype(nil, []*Node{anonfield(typs[1]), anonfield(typs[32]), anonfield(typs[32])}, []*Node{anonfield(typs[93])})
	typs[95] = functype(nil, []*Node{anonfield(typs[1]), anonfield(typs[15]), anonfield(typs[15])}, []*Node{anonfield(typs[93])})
	typs[96] = functype(nil, []*Node{anonfield(typs[1]), anonfield(typs[93]), anonfield(typs[32])}, []*Node{anonfield(typs[93])})
	typs[97] = functype(nil, []*Node{anonfield(typs[3]), anonfield(typs[3]), anonfield(typs[49])}, nil)
	typs[98] = functype(nil, []*Node{anonfield(typs[58]), anonfield(typs[49])}, nil)
	typs[99] = functype(nil, []*Node{anonfield(typs[3]), anonfield(typs[3]), anonfield(typs[49])}, []*Node{anonfield(typs[11])})
	typs[100] = functype(nil, []*Node{anonfield(typs[3]), anonfield(typs[3])}, []*Node{anonfield(typs[11])})
	typs[101] = functype(nil, []*Node{anonfield(typs[15]), anonfield(typs[15])}, []*Node{anonfield(typs[15])})
	typs[102] = functype(nil, []*Node{anonfield(typs[17]), anonfield(typs[17])}, []*Node{anonfield(typs[17])})
	typs[103] = functype(nil, []*Node{anonfield(typs[13])}, []*Node{anonfield(typs[15])})
	typs[104] = functype(nil, []*Node{anonfield(typs[13])}, []*Node{anonfield(typs[17])})
	typs[105] = types.Types[TUINT32]
	typs[106] = functype(nil, []*Node{anonfield(typs[13])}, []*Node{anonfield(typs[105])})
	typs[107] = functype(nil, []*Node{anonfield(typs[15])}, []*Node{anonfield(typs[13])})
	typs[108] = functype(nil, []*Node{anonfield(typs[17])}, []*Node{anonfield(typs[13])})
	typs[109] = functype(nil, []*Node{anonfield(typs[105])}, []*Node{anonfield(typs[13])})
	typs[110] = functype(nil, []*Node{anonfield(typs[19]), anonfield(typs[19])}, []*Node{anonfield(typs[19])})
	typs[111] = functype(nil, []*Node{anonfield(typs[49])}, nil)
	typs[112] = functype(nil, []*Node{anonfield(typs[49]), anonfield(typs[49])}, nil)
	return typs[:]
}
