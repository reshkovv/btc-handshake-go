package main

import (
	"bytes"
	"encoding/binary"
	"math/rand"
	"net"
	"time"
)

type VersionMessage struct {
	ProtocolVersion int32
	Service         uint64
	Timestamp       int64
	AddrRecv        net.TCPAddr
	AddrFrom        net.TCPAddr
	Nonce           uint64
	UserAgent       string
	StartHeight     int32
}

func (m *VersionMessage) ToBytes() []byte {
	var buf bytes.Buffer

	// Write ProtocolVersion
	err := binary.Write(&buf, binary.LittleEndian, m.ProtocolVersion)
	if err != nil {
		panic(err)
	}

	// Write Service
	err = binary.Write(&buf, binary.LittleEndian, m.Service)
	if err != nil {
		panic(err)
	}

	// Write Timestamp
	err = binary.Write(&buf, binary.LittleEndian, m.Timestamp)
	if err != nil {
		panic(err)
	}

	addressRecvBytes := netAddrAsBytes(&m.Service, &m.AddrRecv)
	// Write AddrRecv
	_, err = buf.Write(addressRecvBytes)
	if err != nil {
		panic(err)
	}
	addressFromBytes := netAddrAsBytes(&m.Service, &m.AddrFrom)
	// Write AddrFrom
	_, err = buf.Write(addressFromBytes)
	if err != nil {
		panic(err)
	}

	// Write Nonce
	err = binary.Write(&buf, binary.LittleEndian, m.Nonce)
	if err != nil {
		panic(err)
	}

	// Write UserAgent
	_, err = buf.WriteString(m.UserAgent)
	if err != nil {
		panic(err)
	}
	buf.WriteByte(0x00) // Null terminator

	// Write StartHeight
	err = binary.Write(&buf, binary.LittleEndian, m.StartHeight)
	if err != nil {
		panic(err)
	}

	return buf.Bytes()
}

func netAddrAsBytes(nodeBitmask *uint64, address *net.TCPAddr) []byte {
	var buffer []byte

	// Convert node_bitmask to little-endian bytes and append to buffer
	bitmaskBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bitmaskBytes, *nodeBitmask)
	buffer = append(buffer, bitmaskBytes...)

	// Convert IP address to IPv6-compatible bytes and append to buffer
	ipv6Addr := address.IP.To16()
	buffer = append(buffer, ipv6Addr...)

	// Convert port to big-endian bytes and append to buffer
	portBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(portBytes, uint16(address.Port))
	buffer = append(buffer, portBytes...)

	return buffer
}

func NewVersionMessage(protocolVersion int32, addrRecv net.TCPAddr) VersionMessage {
	timestamp := time.Now().Unix()
	addrFrom := net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 8080}
	nonce := rand.Int()
	userAgent := ""
	startHeight := int32(1)

	return VersionMessage{
		ProtocolVersion: protocolVersion,
		Service:         0,
		Timestamp:       timestamp,
		AddrRecv:        addrRecv,
		AddrFrom:        addrFrom,
		Nonce:           uint64(nonce),
		UserAgent:       userAgent,
		StartHeight:     startHeight,
	}
}
