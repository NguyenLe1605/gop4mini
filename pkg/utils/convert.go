// This package contains helper functions for encoding to byte strings:
// - integers
// - IPv4 address strings
// - Ethernet address strings
package utils

import (
	"fmt"
	"math"
	"net"
	"regexp"
)

var macPattern *regexp.Regexp
var ipPattern *regexp.Regexp

func init() {
	macPattern = regexp.MustCompile(`^([\da-fA-F]{2}:){5}([\da-fA-F]{2})$`)
	ipPattern = regexp.MustCompile(`^(\d{1,3}\.){3}(\d{1,3})$`)
}

func MatchesMac(addr string) bool {
	return macPattern.MatchString(addr)
}

func EncodeMac(addr string) ([]byte, error) {
	return net.ParseMAC(addr)
}

func DecodeMac(addr net.HardwareAddr) string {
	return addr.String()
}

func MatchesIPv4(addr string) bool {
	return ipPattern.MatchString(addr)
}

func EncodeIPv4(addr string) []byte {
	return net.ParseIP(addr).To4()
}

func DecodeIPv4(addr net.IP) string {
	return addr.String()
}

func BitWidthToBytes(bitwidth int) int {
	return int(math.Ceil(float64(bitwidth) / 8.0))
}

func EncodeNum(number int64, bitWidth int) ([]byte, error) {
	if number >= 1<<bitWidth {
		return nil, fmt.Errorf("number, %d, does not fit in %d bits", number, bitWidth)
	}
	byteLen := BitWidthToBytes(bitWidth)
	bytes := make([]byte, byteLen)
	for i := 0; i < byteLen; i++ {
		bytes[byteLen-i-1] = byte(number >> (i * 8))
	}

	return bytes, nil
}

func DecodeNum(number []byte) int64 {
	byteLen := len(number)
	var res int64 = 0
	for i := 0; i < byteLen; i++ {
		res |= int64(number[byteLen-i-1]) << (i * 8)
	}

	return res
}

func CastToInt64(value any) (int64, error) {
	switch v := value.(type) {
	case int:
		return int64(v), nil
	case int16:
		return int64(v), nil
	case int32:
		return int64(v), nil
	case int64:
		return int64(v), nil
	case uint:
		return int64(v), nil
	case uint16:
		return int64(v), nil
	case uint32:
		return int64(v), nil
	case uint64:
		return int64(v), nil
	default:
		return int64(0), fmt.Errorf("unsupported integer type")
	}
}

func Encode(x any, bitwdith int) ([]byte, error) {
	// Tries to infer the type of `x` and encode it
	byteLen := BitWidthToBytes(bitwdith)
	var encodedBytes []byte
	var err error
	switch x := x.(type) {
	case string:
		if MatchesMac(x) {
			encodedBytes, err = EncodeMac(x)
			if err != nil {
				return nil, err
			}
		} else if MatchesIPv4(x) {
			encodedBytes = EncodeIPv4(x)
		} else {
			// Assume that the string is already encoded
			encodedBytes = []byte(x)
		}
	case int, int16, int32, int64, uint, uint16, uint32, uint64:
		x64, err := CastToInt64(x)
		if err != nil {
			return nil, err
		}
		encodedBytes, err = EncodeNum(x64, bitwdith)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("encoding objects of %T is not support", x)
	}
	if len(encodedBytes) != byteLen {
		return nil, fmt.Errorf("can not convert into bytes with bitwidth %d", bitwdith)
	}
	return encodedBytes, nil
}
