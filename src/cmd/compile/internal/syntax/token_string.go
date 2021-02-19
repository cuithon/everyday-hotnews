// Code generated by "stringer -type token -linecomment tokens.go"; DO NOT EDIT.

package syntax

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[_EOF-1]
	_ = x[_Name-2]
	_ = x[_Literal-3]
	_ = x[_Operator-4]
	_ = x[_AssignOp-5]
	_ = x[_IncOp-6]
	_ = x[_Assign-7]
	_ = x[_Define-8]
	_ = x[_Arrow-9]
	_ = x[_Star-10]
	_ = x[_Lparen-11]
	_ = x[_Lbrack-12]
	_ = x[_Lbrace-13]
	_ = x[_Rparen-14]
	_ = x[_Rbrack-15]
	_ = x[_Rbrace-16]
	_ = x[_Comma-17]
	_ = x[_Semi-18]
	_ = x[_Colon-19]
	_ = x[_Dot-20]
	_ = x[_DotDotDot-21]
	_ = x[_Break-22]
	_ = x[_Case-23]
	_ = x[_Chan-24]
	_ = x[_Const-25]
	_ = x[_Continue-26]
	_ = x[_Default-27]
	_ = x[_Defer-28]
	_ = x[_Else-29]
	_ = x[_Fallthrough-30]
	_ = x[_For-31]
	_ = x[_Func-32]
	_ = x[_Go-33]
	_ = x[_Goto-34]
	_ = x[_If-35]
	_ = x[_Import-36]
	_ = x[_Interface-37]
	_ = x[_Map-38]
	_ = x[_Package-39]
	_ = x[_Range-40]
	_ = x[_Return-41]
	_ = x[_Select-42]
	_ = x[_Struct-43]
	_ = x[_Switch-44]
	_ = x[_Type-45]
	_ = x[_Var-46]
	_ = x[tokenCount-47]
}

const _token_name = "EOFnameliteralopop=opop=:=<-*([{)]},;:....breakcasechanconstcontinuedefaultdeferelsefallthroughforfuncgogotoifimportinterfacemappackagerangereturnselectstructswitchtypevar"

var _token_index = [...]uint8{0, 3, 7, 14, 16, 19, 23, 24, 26, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 42, 47, 51, 55, 60, 68, 75, 80, 84, 95, 98, 102, 104, 108, 110, 116, 125, 128, 135, 140, 146, 152, 158, 164, 168, 171, 171}

func (i token) String() string {
	i -= 1
	if i >= token(len(_token_index)-1) {
		return "token(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	return _token_name[_token_index[i]:_token_index[i+1]]
}
