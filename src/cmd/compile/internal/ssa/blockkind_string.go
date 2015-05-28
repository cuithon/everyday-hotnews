// generated by stringer -type=BlockKind; DO NOT EDIT

package ssa

import "fmt"

const (
	_BlockKind_name_0 = "blockInvalid"
	_BlockKind_name_1 = "blockGenericStartBlockExitBlockPlainBlockIfBlockCall"
	_BlockKind_name_2 = "blockAMD64StartBlockEQBlockNEBlockLTBlockLEBlockGTBlockGEBlockULTBlockULEBlockUGTBlockUGE"
)

var (
	_BlockKind_index_0 = [...]uint8{0, 12}
	_BlockKind_index_1 = [...]uint8{0, 17, 26, 36, 43, 52}
	_BlockKind_index_2 = [...]uint8{0, 15, 22, 29, 36, 43, 50, 57, 65, 73, 81, 89}
)

func (i BlockKind) String() string {
	switch {
	case i == 0:
		return _BlockKind_name_0
	case 101 <= i && i <= 105:
		i -= 101
		return _BlockKind_name_1[_BlockKind_index_1[i]:_BlockKind_index_1[i+1]]
	case 201 <= i && i <= 211:
		i -= 201
		return _BlockKind_name_2[_BlockKind_index_2[i]:_BlockKind_index_2[i+1]]
	default:
		return fmt.Sprintf("BlockKind(%d)", i)
	}
}
