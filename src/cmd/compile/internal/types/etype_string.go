// Code generated by "stringer -type EType -trimprefix T"; DO NOT EDIT.

package types

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Txxx-0]
	_ = x[TINT8-1]
	_ = x[TUINT8-2]
	_ = x[TINT16-3]
	_ = x[TUINT16-4]
	_ = x[TINT32-5]
	_ = x[TUINT32-6]
	_ = x[TINT64-7]
	_ = x[TUINT64-8]
	_ = x[TINT-9]
	_ = x[TUINT-10]
	_ = x[TUINTPTR-11]
	_ = x[TCOMPLEX64-12]
	_ = x[TCOMPLEX128-13]
	_ = x[TFLOAT32-14]
	_ = x[TFLOAT64-15]
	_ = x[TBOOL-16]
	_ = x[TPTR-17]
	_ = x[TFUNC-18]
	_ = x[TSLICE-19]
	_ = x[TARRAY-20]
	_ = x[TSTRUCT-21]
	_ = x[TCHAN-22]
	_ = x[TMAP-23]
	_ = x[TINTER-24]
	_ = x[TFORW-25]
	_ = x[TANY-26]
	_ = x[TSTRING-27]
	_ = x[TUNSAFEPTR-28]
	_ = x[TIDEAL-29]
	_ = x[TNIL-30]
	_ = x[TBLANK-31]
	_ = x[TFUNCARGS-32]
	_ = x[TCHANARGS-33]
	_ = x[TSSA-34]
	_ = x[TTUPLE-35]
	_ = x[NTYPE-36]
}

const _EType_name = "xxxINT8UINT8INT16UINT16INT32UINT32INT64UINT64INTUINTUINTPTRCOMPLEX64COMPLEX128FLOAT32FLOAT64BOOLPTRFUNCSLICEARRAYSTRUCTCHANMAPINTERFORWANYSTRINGUNSAFEPTRIDEALNILBLANKFUNCARGSCHANARGSSSATUPLENTYPE"

var _EType_index = [...]uint8{0, 3, 7, 12, 17, 23, 28, 34, 39, 45, 48, 52, 59, 68, 78, 85, 92, 96, 99, 103, 108, 113, 119, 123, 126, 131, 135, 138, 144, 153, 158, 161, 166, 174, 182, 185, 190, 195}

func (i EType) String() string {
	if i >= EType(len(_EType_index)-1) {
		return "EType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _EType_name[_EType_index[i]:_EType_index[i+1]]
}
