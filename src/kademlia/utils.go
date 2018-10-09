package kademlia

import (
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"encoding/hex"
	"io"
	"net"
	"strconv"
)

func IPToLong(ip string) uint32 {
	var long uint32
	binary.Read(bytes.NewBuffer(net.ParseIP(ip).To4()), binary.BigEndian, &long)
	return long
}

func IPToStr(ipInt uint32) string {
	var ipLong int64 = int64(ipInt)

	b0 := strconv.FormatInt((ipLong>>24)&0xff, 10)
	b1 := strconv.FormatInt((ipLong>>16)&0xff, 10)
	b2 := strconv.FormatInt((ipLong>>8)&0xff, 10)
	b3 := strconv.FormatInt((ipLong & 0xff), 10)
	return b0 + "." + b1 + "." + b2 + "." + b3
}

func HashBytes(data []byte) string {
	h := sha1.New()
	io.WriteString(h, string(data))
	return hex.EncodeToString(h.Sum(nil))
}
