package h265

import (
	"fmt"
)

// NALUType is the type of a NALU.
type NALUType uint8
//enum { NAL_VPS = 32, NAL_SPS = 33, NAL_PPS = 34, NAL_AUD = 35, NAL_PREFIX_SEI = 39, };
// NALU types.
const (
	NAL_VPS                        NALUType = 32
	NAL_SPS                         NALUType = 33
	NAL_PPS                        NALUType = 34
	NAL_AUD                        NALUType = 35
	NAL_PREFIX_SEI                           NALUType = 36

)

var naluTypelabels = map[NALUType]string{
	NAL_VPS:                        "vps",
	NAL_SPS:                "sps",
	NAL_PPS:                "pps",
	NAL_AUD:                "aud",
	NAL_PREFIX_SEI:          "sei",

}

// String implements fmt.Stringer.
func (nt NALUType) String() string {
	if l, ok := naluTypelabels[nt]; ok {
		return l
	}
	return fmt.Sprintf("unknown (%d)", nt)
}
