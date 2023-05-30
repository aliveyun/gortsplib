package formats

import (
	"testing"

	"github.com/pion/rtp"
	"github.com/stretchr/testify/require"
)

func TestMJPEGAttributes(t *testing.T) {
	format := &MJPEG{}
	require.Equal(t, "M-JPEG", format.String())
	require.Equal(t, 90000, format.ClockRate())
	require.Equal(t, true, format.PTSEqualsDTS(&rtp.Packet{}))
}

func TestMJPEGDecEncoder(t *testing.T) {
	format := &MJPEG{}

	b := []byte{
		0xff, 0xd8, 0xff, 0xdb, 0x00, 0x84, 0x00, 0x0d,
		0x09, 0x0a, 0x0b, 0x0a, 0x08, 0x0d, 0x0b, 0x0a,
		0x0b, 0x0e, 0x0e, 0x0d, 0x0f, 0x13, 0x20, 0x15,
		0x13, 0x12, 0x12, 0x13, 0x27, 0x1c, 0x1e, 0x17,
		0x20, 0x2e, 0x29, 0x31, 0x30, 0x2e, 0x29, 0x2d,
		0x2c, 0x33, 0x3a, 0x4a, 0x3e, 0x33, 0x36, 0x46,
		0x37, 0x2c, 0x2d, 0x40, 0x57, 0x41, 0x46, 0x4c,
		0x4e, 0x52, 0x53, 0x52, 0x32, 0x3e, 0x5a, 0x61,
		0x5a, 0x50, 0x60, 0x4a, 0x51, 0x52, 0x4f, 0x01,
		0x0e, 0x0e, 0x0e, 0x13, 0x11, 0x13, 0x26, 0x15,
		0x15, 0x26, 0x4f, 0x35, 0x2d, 0x35, 0x4f, 0x4f,
		0x4f, 0x4f, 0x4f, 0x4f, 0x4f, 0x4f, 0x4f, 0x4f,
		0x4f, 0x4f, 0x4f, 0x4f, 0x4f, 0x4f, 0x4f, 0x4f,
		0x4f, 0x4f, 0x4f, 0x4f, 0x4f, 0x4f, 0x4f, 0x4f,
		0x4f, 0x4f, 0x4f, 0x4f, 0x4f, 0x4f, 0x4f, 0x4f,
		0x4f, 0x4f, 0x4f, 0x4f, 0x4f, 0x4f, 0x4f, 0x4f,
		0x4f, 0x4f, 0x4f, 0x4f, 0x4f, 0x4f, 0x4f, 0x4f,
		0xff, 0xc0, 0x00, 0x11, 0x08, 0x04, 0x38, 0x07,
		0x80, 0x03, 0x00, 0x22, 0x00, 0x01, 0x11, 0x01,
		0x02, 0x11, 0x01, 0xff, 0xc4, 0x00, 0x1f, 0x00,
		0x00, 0x01, 0x05, 0x01, 0x01, 0x01, 0x01, 0x01,
		0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
		0x08, 0x09, 0x0a, 0x0b, 0xff, 0xc4, 0x00, 0xb5,
		0x10, 0x00, 0x02, 0x01, 0x03, 0x03, 0x02, 0x04,
		0x03, 0x05, 0x05, 0x04, 0x04, 0x00, 0x00, 0x01,
		0x7d, 0x01, 0x02, 0x03, 0x00, 0x04, 0x11, 0x05,
		0x12, 0x21, 0x31, 0x41, 0x06, 0x13, 0x51, 0x61,
		0x07, 0x22, 0x71, 0x14, 0x32, 0x81, 0x91, 0xa1,
		0x08, 0x23, 0x42, 0xb1, 0xc1, 0x15, 0x52, 0xd1,
		0xf0, 0x24, 0x33, 0x62, 0x72, 0x82, 0x09, 0x0a,
		0x16, 0x17, 0x18, 0x19, 0x1a, 0x25, 0x26, 0x27,
		0x28, 0x29, 0x2a, 0x34, 0x35, 0x36, 0x37, 0x38,
		0x39, 0x3a, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48,
		0x49, 0x4a, 0x53, 0x54, 0x55, 0x56, 0x57, 0x58,
		0x59, 0x5a, 0x63, 0x64, 0x65, 0x66, 0x67, 0x68,
		0x69, 0x6a, 0x73, 0x74, 0x75, 0x76, 0x77, 0x78,
		0x79, 0x7a, 0x83, 0x84, 0x85, 0x86, 0x87, 0x88,
		0x89, 0x8a, 0x92, 0x93, 0x94, 0x95, 0x96, 0x97,
		0x98, 0x99, 0x9a, 0xa2, 0xa3, 0xa4, 0xa5, 0xa6,
		0xa7, 0xa8, 0xa9, 0xaa, 0xb2, 0xb3, 0xb4, 0xb5,
		0xb6, 0xb7, 0xb8, 0xb9, 0xba, 0xc2, 0xc3, 0xc4,
		0xc5, 0xc6, 0xc7, 0xc8, 0xc9, 0xca, 0xd2, 0xd3,
		0xd4, 0xd5, 0xd6, 0xd7, 0xd8, 0xd9, 0xda, 0xe1,
		0xe2, 0xe3, 0xe4, 0xe5, 0xe6, 0xe7, 0xe8, 0xe9,
		0xea, 0xf1, 0xf2, 0xf3, 0xf4, 0xf5, 0xf6, 0xf7,
		0xf8, 0xf9, 0xfa, 0xff, 0xc4, 0x00, 0x1f, 0x01,
		0x00, 0x03, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01,
		0x01, 0x01, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
		0x08, 0x09, 0x0a, 0x0b, 0xff, 0xc4, 0x00, 0xb5,
		0x11, 0x00, 0x02, 0x01, 0x02, 0x04, 0x04, 0x03,
		0x04, 0x07, 0x05, 0x04, 0x04, 0x00, 0x01, 0x02,
		0x77, 0x00, 0x01, 0x02, 0x03, 0x11, 0x04, 0x05,
		0x21, 0x31, 0x06, 0x12, 0x41, 0x51, 0x07, 0x61,
		0x71, 0x13, 0x22, 0x32, 0x81, 0x08, 0x14, 0x42,
		0x91, 0xa1, 0xb1, 0xc1, 0x09, 0x23, 0x33, 0x52,
		0xf0, 0x15, 0x62, 0x72, 0xd1, 0x0a, 0x16, 0x24,
		0x34, 0xe1, 0x25, 0xf1, 0x17, 0x18, 0x19, 0x1a,
		0x26, 0x27, 0x28, 0x29, 0x2a, 0x35, 0x36, 0x37,
		0x38, 0x39, 0x3a, 0x43, 0x44, 0x45, 0x46, 0x47,
		0x48, 0x49, 0x4a, 0x53, 0x54, 0x55, 0x56, 0x57,
		0x58, 0x59, 0x5a, 0x63, 0x64, 0x65, 0x66, 0x67,
		0x68, 0x69, 0x6a, 0x73, 0x74, 0x75, 0x76, 0x77,
		0x78, 0x79, 0x7a, 0x82, 0x83, 0x84, 0x85, 0x86,
		0x87, 0x88, 0x89, 0x8a, 0x92, 0x93, 0x94, 0x95,
		0x96, 0x97, 0x98, 0x99, 0x9a, 0xa2, 0xa3, 0xa4,
		0xa5, 0xa6, 0xa7, 0xa8, 0xa9, 0xaa, 0xb2, 0xb3,
		0xb4, 0xb5, 0xb6, 0xb7, 0xb8, 0xb9, 0xba, 0xc2,
		0xc3, 0xc4, 0xc5, 0xc6, 0xc7, 0xc8, 0xc9, 0xca,
		0xd2, 0xd3, 0xd4, 0xd5, 0xd6, 0xd7, 0xd8, 0xd9,
		0xda, 0xe2, 0xe3, 0xe4, 0xe5, 0xe6, 0xe7, 0xe8,
		0xe9, 0xea, 0xf2, 0xf3, 0xf4, 0xf5, 0xf6, 0xf7,
		0xf8, 0xf9, 0xfa, 0xff, 0xda, 0x00, 0x0c, 0x03,
		0x00, 0x00, 0x01, 0x11, 0x02, 0x11, 0x00, 0x3f,
		0x00, 0x92, 0x8a, 0x28, 0xaf, 0x54, 0xf2, 0x42,
		0x8a, 0x28, 0xa0, 0x02, 0x96, 0x92, 0x96, 0x80,
		0x0a, 0x4a, 0x75, 0x25, 0x02, 0x12, 0x8a, 0x5a,
		0x28, 0x18, 0x94, 0x52, 0xd1, 0x40, 0x09, 0x45,
		0x2d, 0x14, 0x08, 0x29, 0x69, 0x29, 0x68, 0x00,
		0xa5, 0xa4, 0xa5, 0xa0, 0x02, 0x8a, 0x28, 0xa0,
		0x02, 0x8a, 0x28, 0xa0, 0x04, 0xa5, 0xa2, 0x8a,
		0x00, 0x5a, 0x28, 0xa2, 0x80, 0x0a, 0x28, 0xa2,
		0x80, 0x0a, 0x28, 0xa2, 0x80, 0x12, 0x8a, 0x5a,
		0x28, 0x24, 0x29, 0x69, 0x29, 0x68, 0x00, 0xa2,
		0x8a, 0x28, 0x00, 0xa2, 0x8a, 0x28, 0x00, 0xa2,
		0x8a, 0x28, 0x00, 0xa2, 0x8a, 0x28, 0x00, 0xa2,
		0x8a, 0x28, 0x00, 0xa2, 0x8a, 0x28, 0x00, 0xa2,
		0x8a, 0x28, 0x00, 0xa2, 0x8a, 0x28, 0x00, 0xa2,
		0x8a, 0x28, 0x00, 0xa2, 0x8a, 0x28, 0x00, 0xa2,
		0x8a, 0x28, 0x00, 0xa4, 0xa5, 0xa4, 0xa0, 0x02,
		0x8a, 0x28, 0xa0, 0x02, 0x8a, 0x28, 0xa0, 0x02,
		0x8a, 0x28, 0xa0, 0x02, 0x8a, 0x28, 0xa0, 0x02,
		0x8a, 0x28, 0xa0, 0x02, 0x8a, 0x28, 0xa0, 0x02,
		0x8a, 0x28, 0xa0, 0x02, 0x8a, 0x28, 0xa0, 0x02,
		0x96, 0x92, 0x96, 0x80, 0x0a, 0x28, 0xa2, 0x80,
		0x0a, 0x28, 0xa2, 0x80, 0x0a, 0x28, 0xa2, 0x80,
		0x0a, 0x28, 0xa2, 0x80, 0x0a, 0x28, 0xa2, 0x80,
		0x0a, 0x28, 0xa2, 0x80, 0x0a, 0x28, 0xa2, 0x80,
		0x0a, 0x28, 0xa2, 0x80, 0x0a, 0x28, 0xa2, 0x81,
		0x85, 0x14, 0x51, 0x40, 0x05, 0x14, 0x51, 0x40,
		0x05, 0x14, 0x51, 0x40, 0x05, 0x14, 0x52, 0xd0,
		0x01, 0x45, 0x14, 0x50, 0x01, 0x45, 0x14, 0x50,
		0x01, 0x45, 0x14, 0x50, 0x01, 0x45, 0x2d, 0x14,
		0x00, 0x94, 0xb4, 0x51, 0x40, 0x05, 0x14, 0x52,
		0xd0, 0x02, 0x51, 0x4b, 0x45, 0x00, 0x25, 0x2d,
		0x14, 0x50, 0x01, 0x45, 0x14, 0x50, 0x01, 0x45,
		0x14, 0x50, 0x20, 0xa5, 0xa4, 0xa5, 0xa0, 0x02,
		0x8a, 0x28, 0xa0, 0x02, 0x8a, 0x28, 0xa0, 0x02,
		0x8a, 0x5a, 0x28, 0x18, 0x94, 0xb4, 0x51, 0x40,
		0xc2, 0x8a, 0x28, 0xa0, 0x05, 0xa2, 0x92, 0x9d,
		0x40, 0x05, 0x14, 0x51, 0x48, 0x02, 0x8a, 0x28,
		0xa4, 0x01, 0x4b, 0x49, 0x4b, 0x40, 0x05, 0x14,
		0x51, 0x40, 0x05, 0x14, 0xb4, 0x50, 0x02, 0x51,
		0x4b, 0x45, 0x00, 0x25, 0x2d, 0x14, 0x50, 0x03,
		0xa8, 0xa2, 0x8a, 0x00, 0x28, 0xa2, 0x8a, 0x00,
		0x5a, 0x29, 0x29, 0x68, 0x00, 0xa2, 0x8a, 0x28,
		0x00, 0xa2, 0x96, 0x8a, 0x06, 0x25, 0x14, 0xb4,
		0x50, 0x01, 0x45, 0x14, 0x50, 0x02, 0xd2, 0xd2,
		0x52, 0xd0, 0x20, 0xa2, 0x8a, 0x28, 0x01, 0x68,
		0xa2, 0x8a, 0x40, 0x14, 0x51, 0x45, 0x30, 0x0a,
		0x5a, 0x4a, 0x5a, 0x06, 0x14, 0x51, 0x45, 0x02,
		0x0a, 0x28, 0xa5, 0xa0, 0x62, 0x51, 0x4b, 0x45,
		0x00, 0x2d, 0x14, 0x51, 0x48, 0x02, 0x8a, 0x28,
		0xa0, 0x61, 0x45, 0x14, 0x50, 0x03, 0xa8, 0xa2,
		0x8a, 0x06, 0x2d, 0x14, 0x51, 0x48, 0x02, 0x8a,
		0x28, 0xa4, 0x30, 0xa2, 0x8a, 0x2a, 0x80, 0x28,
		0xa2, 0x8a, 0x00, 0x28, 0xa2, 0x8a, 0x92, 0x45,
		0xa5, 0xa2, 0x96, 0x82, 0x82, 0x8a, 0x28, 0xa0,
		0x02, 0x8a, 0x28, 0xa0, 0x05, 0xa2, 0x8a, 0x29,
		0x80, 0x52, 0xd2, 0x52, 0xd0, 0x01, 0x45, 0x14,
		0x50, 0x01, 0x4e, 0xa2, 0x8a, 0x43, 0x0a, 0x28,
		0xa2, 0x80, 0x0a, 0x28, 0xa4, 0xa4, 0x31, 0x68,
		0xa4, 0xf3, 0x62, 0xff, 0x00, 0x9e, 0xd1, 0x7e,
		0xea, 0x9f, 0xfe, 0xb6, 0xa4, 0x62, 0x52, 0x53,
		0xa9, 0x28, 0x01, 0x28, 0xa2, 0x6f, 0xdd, 0x7f,
		0xaf, 0xa5, 0xa0, 0x62, 0x51, 0x4b, 0x45, 0x00,
		0x25, 0x14, 0xbf, 0xba, 0xff, 0x00, 0xae, 0xb4,
		0xea, 0x60, 0x36, 0x9d, 0x49, 0x49, 0x34, 0xb1,
		0x45, 0xfe, 0xbe, 0x6f, 0x2a, 0x98, 0x0f, 0xa2,
		0xb9, 0xbd, 0x0f, 0x59, 0x97, 0x54, 0xf1, 0x2d,
		0xd7, 0xfc, 0xf9, 0xd7, 0x49, 0x52, 0x30, 0xac,
		0x7d, 0x5b, 0x54, 0xd5, 0x74, 0xbf, 0xdf, 0x41,
		0x67, 0x15, 0xd4, 0x35, 0x63, 0x56, 0xd5, 0x22,
		0xd2, 0xe0, 0xf3, 0xbc, 0x99, 0x65, 0xae, 0x4b,
		0x51, 0xf1, 0x44, 0x52, 0xff, 0x00, 0xc7, 0x97,
		0x9b, 0xe4, 0xd0, 0x69, 0x4c, 0x76, 0xa3, 0xe2,
		0xd9, 0xa5, 0x82, 0x29, 0xb4, 0xbf, 0xdd, 0x79,
		0xbf, 0xeb, 0x22, 0xaa, 0x7e, 0x1e, 0xd5, 0x3e,
		0xc1, 0x3f, 0xfa, 0x9f, 0xfa, 0xe9, 0x25, 0x61,
		0xd6, 0x9e, 0x87, 0xf6, 0xbf, 0xed, 0x68, 0xbc,
		0x8a, 0x82, 0xcf, 0x4e, 0x86, 0x5f, 0x36, 0x0a,
		0x92, 0x9b, 0xff, 0x00, 0x5d, 0xeb, 0x3a, 0x6d,
		0x53, 0xfd, 0x3f, 0xec, 0x76, 0x3f, 0xbd, 0xff,
		0x00, 0xa6, 0xbf, 0xc1, 0x54, 0x66, 0x69, 0x51,
		0x50, 0x4b, 0x2c, 0x51, 0x7e, 0xe7, 0xce, 0xf3,
		0x66, 0xff, 0x00, 0xa6, 0x75, 0x35, 0x00, 0x2d,
		0x65, 0x6a, 0xde, 0x23, 0xd3, 0xec, 0x3f, 0xe5,
		0xb7, 0x9b, 0x34, 0x5f, 0xf2, 0xce, 0xb4, 0x7e,
		0xd5, 0x69, 0x2f, 0xfc, 0xbe, 0x45, 0x5e, 0x7b,
		0xe2, 0x7d, 0x07, 0xec, 0x1e, 0x6e, 0xa5, 0xf6,
		0xc8, 0xa5, 0xf3, 0x6a, 0xc0, 0xdf, 0xa2, 0x9d,
		0x45, 0x6c, 0x72, 0x0d, 0xa5, 0xa5, 0xa2, 0x80,
		0x12, 0x96, 0x8a, 0x28, 0x10, 0x51, 0x4b, 0x45,
		0x00, 0x25, 0x14, 0xb4, 0x50, 0x02, 0x51, 0x4b,
		0x45, 0x03, 0x12, 0x8a, 0x5a, 0x28, 0x10, 0x94,
		0xb4, 0x51, 0x40, 0x05, 0x14, 0x51, 0x40, 0x05,
		0x14, 0xb4, 0x50, 0x02, 0x51, 0x4b, 0x45, 0x00,
		0x14, 0x51, 0x4b, 0x40, 0x84, 0xa2, 0x8a, 0x28,
		0x00, 0xa2, 0x96, 0x8a, 0x00, 0x4a, 0x29, 0x68,
		0xa0, 0x04, 0xa2, 0x96, 0x8a, 0x00, 0x4a, 0x29,
		0x68, 0xa0, 0x41, 0x45, 0x14, 0x50, 0x01, 0x45,
		0x14, 0x50, 0x01, 0x45, 0x14, 0x50, 0x01, 0x45,
		0x14, 0x50, 0x01, 0x45, 0x14, 0x50, 0x01, 0x45,
		0x14, 0x50, 0x01, 0x45, 0x14, 0x50, 0x01, 0x45,
		0x14, 0x50, 0x01, 0x45, 0x14, 0x50, 0x01, 0x45,
		0x14, 0x50, 0x01, 0x45, 0x14, 0x50, 0x01, 0x49,
		0x4b, 0x45, 0x00, 0x25, 0x14, 0xb4, 0x50, 0x01,
		0x45, 0x14, 0x50, 0x30, 0xa4, 0xa5, 0xa2, 0x81,
		0x09, 0x45, 0x2d, 0x14, 0x00, 0x94, 0x52, 0xd1,
		0x40, 0x09, 0x45, 0x2d, 0x14, 0x00, 0x94, 0x52,
		0xd1, 0x40, 0x05, 0x14, 0xb4, 0x94, 0x0c, 0x28,
		0xa2, 0x8a, 0x04, 0x14, 0x51, 0x45, 0x00, 0x14,
		0x51, 0x45, 0x00, 0x14, 0x51, 0x45, 0x00, 0x14,
		0x51, 0x45, 0x00, 0x14, 0x52, 0xd1, 0x40, 0xc4,
		0xa2, 0x96, 0x8a, 0x04, 0x25, 0x14, 0xb4, 0x50,
		0x31, 0x28, 0xa5, 0xa2, 0x80, 0x12, 0x8a, 0x5a,
		0x4a, 0x00, 0x5a, 0x28, 0xa2, 0x80, 0x0a, 0x29,
		0x68, 0xa0, 0x04, 0xa2, 0x96, 0x8a, 0x00, 0x4a,
		0x29, 0x68, 0xa0, 0x04, 0xa2, 0x9d, 0x45, 0x03,
		0x1b, 0x45, 0x3a, 0x8a, 0x00, 0x6d, 0x14, 0xea,
		0x28, 0x10, 0x94, 0x52, 0xd1, 0x40, 0x09, 0x4b,
		0x45, 0x14, 0x08, 0x28, 0xa2, 0x8a, 0x06, 0x14,
		0x52, 0xd1, 0x40, 0x09, 0x45, 0x2d, 0x2d, 0x00,
		0x25, 0x14, 0xb4, 0x50, 0x21, 0x28, 0xa5, 0xa2,
		0x81, 0x89, 0x4b, 0x45, 0x14, 0x00, 0x51, 0x45,
		0x14, 0x00, 0x51, 0x4b, 0x45, 0x00, 0x25, 0x14,
		0xb4, 0x50, 0x01, 0x45, 0x2d, 0x14, 0x0c, 0x28,
		0xa2, 0x8a, 0x00, 0x4a, 0x5a, 0x5a, 0x28, 0x01,
		0x28, 0xa5, 0xa2, 0x90, 0x05, 0x14, 0x51, 0x40,
		0x05, 0x14, 0xb4, 0x52, 0x01, 0x28, 0xa5, 0xa2,
		0x80, 0x0a, 0x28, 0xa7, 0x50, 0x03, 0x68, 0xa7,
		0x51, 0x40, 0x0d, 0xa7, 0x51, 0x45, 0x00, 0x14,
		0x51, 0x4b, 0x40, 0x09, 0x45, 0x2d, 0x2d, 0x00,
		0x25, 0x14, 0xb4, 0x50, 0x02, 0x51, 0x4b, 0x45,
		0x00, 0x14, 0x51, 0x45, 0x00, 0x14, 0x51, 0x45,
		0x00, 0x14, 0xb4, 0x53, 0xa8, 0x18, 0xda, 0x5a,
		0x5a, 0x28, 0x10, 0x94, 0xb4, 0x51, 0x40, 0xc2,
		0x8a, 0x5a, 0x29, 0x00, 0x94, 0x52, 0xd1, 0x40,
		0x05, 0x14, 0x51, 0x4c, 0x02, 0x96, 0x8a, 0x29,
		0x00, 0x52, 0xd2, 0x52, 0xd0, 0x01, 0x45, 0x14,
		0x50, 0x02, 0xd1, 0x45, 0x14, 0x0c, 0x28, 0xa2,
		0x96, 0x80, 0x0a, 0x28, 0xa2, 0x81, 0x85, 0x2d,
		0x14, 0x50, 0x01, 0x45, 0x14, 0xb4, 0x80, 0x4a,
		0x29, 0x68, 0xa0, 0x62, 0xd1, 0x45, 0x14, 0xc0,
		0x28, 0xa5, 0xa2, 0x90, 0x0d, 0xa5, 0xa5, 0xa2,
		0x80, 0x0a, 0x5a, 0x28, 0xa4, 0x01, 0x45, 0x14,
		0xb4, 0x00, 0x94, 0x52, 0xd1, 0x40, 0x05, 0x14,
		0x51, 0x4c, 0x02, 0x8a, 0x29, 0x68, 0x01, 0x28,
		0xa5, 0xa2, 0x80, 0x16, 0x8a, 0x28, 0xa4, 0x30,
		0xa7, 0x55, 0x7b, 0xbb, 0x5f, 0xb7, 0xc1, 0xe4,
		0xff, 0x00, 0xaa, 0xac, 0x49, 0xb4, 0x6d, 0x57,
		0x4b, 0xf3, 0x66, 0xd2, 0xef, 0x3c, 0xd8, 0x69,
		0x17, 0x4c, 0xe9, 0x2b, 0x07, 0xc4, 0x32, 0xcb,
		0x75, 0xa4, 0xff, 0x00, 0xc4, 0xae, 0x68, 0xbf,
		0x75, 0xfe, 0xb2, 0xb9, 0x09, 0xb5, 0x9d, 0x42,
		0x5f, 0x36, 0xcf, 0xce, 0xff, 0x00, 0x5b, 0x56,
		0x3c, 0x3d, 0x6b, 0xfd, 0xa9, 0xab, 0x4b, 0x67,
		0x3f, 0xfc, 0xf3, 0xff, 0x00, 0x96, 0x74, 0x8d,
		0xbd, 0x99, 0x9b, 0x0d, 0xfc, 0xb2, 0xf9, 0xbf,
		0xeb, 0x7f, 0x7b, 0x5a, 0xf0, 0xf8, 0x8e, 0xee,
		0x2d, 0x27, 0xec, 0x70, 0x7f, 0xdf, 0xca, 0x66,
		0xa3, 0xe1, 0x7f, 0xb0, 0x79, 0xbe, 0x46, 0xa5,
		0x17, 0xfd, 0x73, 0x92, 0xb0, 0xe1, 0xa8, 0x0f,
		0x66, 0x76, 0x90, 0xf8, 0xa2, 0x29, 0x6c, 0x3c,
		0x9f, 0xde, 0xf9, 0xde, 0x5f, 0xef, 0x24, 0xaa,
		0x76, 0x9e, 0x2d, 0xbb, 0x8b, 0x49, 0xf2, 0x7f,
		0xe5, 0xb5, 0x62, 0x5a, 0x5f, 0xcb, 0x6b, 0xfe,
		0xa3, 0xfe, 0x5a, 0xd1, 0xe5, 0x45, 0xff, 0x00,
		0x2c, 0x29, 0x0f, 0xd9, 0x97, 0xff, 0x00, 0xb6,
		0x65, 0xba, 0x82, 0x5f, 0xb7, 0x4d, 0xe6, 0xcd,
		0xff, 0x00, 0x2c, 0xe2, 0xae, 0xa7, 0x49, 0xd5,
		0x3f, 0xe2, 0x9a, 0xf3, 0xa7, 0xff, 0x00, 0x5d,
		0x15, 0x70, 0xb4, 0x43, 0xfb, 0xdf, 0xf5, 0xf3,
		0x4b, 0xff, 0x00, 0x5c, 0xe8, 0xf6, 0x83, 0xf6,
		0x67, 0xa1, 0x68, 0x7a, 0xcc, 0x5a, 0xcf, 0xfd,
		0x32, 0x9b, 0xfe, 0x79, 0xd5, 0xeb, 0xbb, 0xab,
		0x4b, 0x08, 0x3f, 0xd3, 0xab, 0xcd, 0x66, 0xba,
		0xfb, 0x57, 0xfd, 0x32, 0x9a, 0x2a, 0x96, 0xd2,
		0xff, 0x00, 0xcd, 0xbf, 0xf3, 0xb5, 0x49, 0xbc,
		0xdf, 0x2b, 0xfe, 0x59, 0xd3, 0x27, 0xd9, 0x9a,
		0x30, 0xdd, 0x7d, 0xab, 0x5e, 0xfb, 0x1d, 0x8c,
		0xde, 0x55, 0x9c, 0xb2, 0x79, 0x95, 0xd3, 0x6a,
		0x3a, 0xa7, 0xd9, 0x75, 0xdb, 0x5b, 0x39, 0xff,
		0x00, 0xd4, 0xcb, 0xff, 0x00, 0x2d, 0x2b, 0xcf,
		0xbe, 0xdf, 0xe5, 0x6a, 0xdf, 0x6c, 0x83, 0xfe,
		0x59, 0x54, 0xda, 0xb6, 0xb3, 0x2e, 0xb3, 0x7f,
		0xe7, 0x7f, 0xaa, 0xff, 0xff, 0xd9,
	}

	enc := format.CreateEncoder()
	pkts, err := enc.Encode(b, 0)
	require.NoError(t, err)
	require.Equal(t, format.PayloadType(), pkts[0].PayloadType)

	dec := format.CreateDecoder()
	var byts []byte
	for _, pkt := range pkts {
		byts, _, _ = dec.Decode(pkt)
	}
	require.Equal(t, b, byts)
}
