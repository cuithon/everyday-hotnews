// Code generated by "stringer -type state"; DO NOT EDIT.

package template

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[stateText-0]
	_ = x[stateTag-1]
	_ = x[stateAttrName-2]
	_ = x[stateAfterName-3]
	_ = x[stateBeforeValue-4]
	_ = x[stateHTMLCmt-5]
	_ = x[stateRCDATA-6]
	_ = x[stateAttr-7]
	_ = x[stateURL-8]
	_ = x[stateSrcset-9]
	_ = x[stateJS-10]
	_ = x[stateJSDqStr-11]
	_ = x[stateJSSqStr-12]
	_ = x[stateJSBqStr-13]
	_ = x[stateJSRegexp-14]
	_ = x[stateJSBlockCmt-15]
	_ = x[stateJSLineCmt-16]
	_ = x[stateCSS-17]
	_ = x[stateCSSDqStr-18]
	_ = x[stateCSSSqStr-19]
	_ = x[stateCSSDqURL-20]
	_ = x[stateCSSSqURL-21]
	_ = x[stateCSSURL-22]
	_ = x[stateCSSBlockCmt-23]
	_ = x[stateCSSLineCmt-24]
	_ = x[stateError-25]
	_ = x[stateDead-26]
}

const _state_name = "stateTextstateTagstateAttrNamestateAfterNamestateBeforeValuestateHTMLCmtstateRCDATAstateAttrstateURLstateSrcsetstateJSstateJSDqStrstateJSSqStrstateJSBqStrstateJSRegexpstateJSBlockCmtstateJSLineCmtstateCSSstateCSSDqStrstateCSSSqStrstateCSSDqURLstateCSSSqURLstateCSSURLstateCSSBlockCmtstateCSSLineCmtstateErrorstateDead"

var _state_index = [...]uint16{0, 9, 17, 30, 44, 60, 72, 83, 92, 100, 111, 118, 130, 142, 154, 167, 182, 196, 204, 217, 230, 243, 256, 267, 283, 298, 308, 317}

func (i state) String() string {
	if i >= state(len(_state_index)-1) {
		return "state(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _state_name[_state_index[i]:_state_index[i+1]]
}
