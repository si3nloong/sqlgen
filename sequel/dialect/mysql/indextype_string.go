// Code generated by "stringer --type indexType --linecomment"; DO NOT EDIT.

package mysql

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[bTree-0]
	_ = x[fullText-1]
	_ = x[unique-2]
	_ = x[spatial-3]
}

const _indexType_name = "BTREEFULLTEXTUNIQUESPATIAL"

var _indexType_index = [...]uint8{0, 5, 13, 19, 26}

func (i indexType) String() string {
	if i >= indexType(len(_indexType_index)-1) {
		return "indexType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _indexType_name[_indexType_index[i]:_indexType_index[i+1]]
}
