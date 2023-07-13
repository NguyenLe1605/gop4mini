package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncodeAndDecodeMacAddress(t *testing.T) {
	mac := "aa:bb:cc:dd:ee:ff"
	require.True(t, MatchesMac(mac))
	encMac, err := EncodeMac(mac)
	require.NoError(t, err)
	require.Equal(t, encMac, []byte("\xaa\xbb\xcc\xdd\xee\xff"))

	decMac := DecodeMac(encMac)
	require.Equal(t, mac, decMac)

	encodeMac, err := Encode(mac, 6 * 8)
	require.NoError(t, err)
	require.Equal(t, encMac, encodeMac)
}

func TestEncodeAndDecodeIPv4Address(t *testing.T) {
	ip := "10.0.0.1"
	require.True(t, MatchesIPv4(ip))

	encIp := EncodeIPv4(ip)
	require.Equal(t, encIp, []byte("\x0a\x00\x00\x01"))

	decIp := DecodeIPv4(encIp)
	require.Equal(t, decIp, ip)

	encodeIP, err := Encode(ip, 4 * 8)
	require.NoError(t, err)
	require.Equal(t, encIp, encodeIP)
}

func TestMatchIPv4(t *testing.T) {
	require.False(t, MatchesIPv4("10.0.0.1.5"))
	require.False(t, MatchesIPv4("1000.0.0.1"))
	require.False(t, MatchesIPv4("10001"))
}

func TestEncodeNum(t *testing.T) {
	var num int64 = 1337
	byteLen := 5
	encNum, err:= EncodeNum(int64(num), byteLen * 8)
	require.NoError(t, err)
	require.Equal(t, encNum, []byte("\x00\x00\x00\x05\x39"))

	decNum := DecodeNum(encNum)
	require.Equal(t, decNum, int64(num))

	encodeNum, err := Encode(num, byteLen * 8)
	require.NoError(t, err)
	require.Equal(t, encodeNum, encNum)

	num = 256
	byteLen = 2
	_, err = Encode(num, byteLen * 2)
	require.Errorf(t, err, "can not convert into bytes with bitwidth %d",byteLen * 8)
}