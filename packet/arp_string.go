// Code generated by "stringer -type=ARPType,ARPOpCode -output arp_string.go"; DO NOT EDIT.

package packet

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[ARPTypeEther-1]
}

const _ARPType_name = "ARPTypeEther"

var _ARPType_index = [...]uint8{0, 12}

func (i ARPType) String() string {
	i -= 1
	if i >= ARPType(len(_ARPType_index)-1) {
		return "ARPType(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	return _ARPType_name[_ARPType_index[i]:_ARPType_index[i+1]]
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[ARPOPCodeRequest-1]
	_ = x[ARPOPCodeReply-2]
}

const _ARPOpCode_name = "ARPOPCodeRequestARPOPCodeReply"

var _ARPOpCode_index = [...]uint8{0, 16, 30}

func (i ARPOpCode) String() string {
	i -= 1
	if i >= ARPOpCode(len(_ARPOpCode_index)-1) {
		return "ARPOpCode(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	return _ARPOpCode_name[_ARPOpCode_index[i]:_ARPOpCode_index[i+1]]
}