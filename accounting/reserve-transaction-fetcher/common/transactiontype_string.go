// Code generated by "stringer -type=TransactionType -linecomment"; DO NOT EDIT.

package common

import "strconv"

const _TransactionType_name = "erc20normalinternal"

var _TransactionType_index = [...]uint8{0, 5, 11, 19}

func (i TransactionType) String() string {
	if i < 0 || i >= TransactionType(len(_TransactionType_index)-1) {
		return "TransactionType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _TransactionType_name[_TransactionType_index[i]:_TransactionType_index[i+1]]
}
