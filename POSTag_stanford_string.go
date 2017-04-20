// +build stanfordtags

// Code generated by "stringer -type=POSTag -output=POSTag_stanford_string.go"; DO NOT EDIT

package lingo

import "fmt"

const _POSTag_name = "XUNKNOWN_TAGROOT_TAGCCCDDTEXFWINJJJJRJJSLSMDNNNNSNNPNNPSPDTPOSPRPPPRPRBRBRRBSRPSYMTOUHVBVBDVBGVBNVBPVBZWDTWPPWPWRBCOMMAFULLSTOPOPENQUOTECLOSEQUOTECOLONDOLLARHASHSIGNLEFTBRACERIGHTBRACEHYPHAFXADDNFPGWXXMAXTAG"

var _POSTag_index = [...]uint8{0, 1, 12, 20, 22, 24, 26, 28, 30, 32, 34, 37, 40, 42, 44, 46, 49, 52, 56, 59, 62, 65, 69, 71, 74, 77, 79, 82, 84, 86, 88, 91, 94, 97, 100, 103, 106, 108, 111, 114, 119, 127, 136, 146, 151, 157, 165, 174, 184, 188, 191, 194, 197, 199, 201, 207}

func (i POSTag) String() string {
	if i >= POSTag(len(_POSTag_index)-1) {
		return fmt.Sprintf("POSTag(%d)", i)
	}
	return _POSTag_name[_POSTag_index[i]:_POSTag_index[i+1]]
}
