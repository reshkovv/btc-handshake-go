package main

import (
	"bytes"
	"encoding/binary"
	"errors"
)

const MAGIC_NUMBER = 0xD9B4BEF9

type BtcMessage struct {
	Magic    uint32
	Command  [12]byte
	Length   uint32
	Checksum uint32
	Payload  []byte
}

func (m *BtcMessage) ToBytes() []byte {
	var buf bytes.Buffer

	// Write Magic
	err := binary.Write(&buf, binary.LittleEndian, m.Magic)
	if err != nil {
		panic(err)
	}
	// Write Command
	_, err = buf.Write(m.Command[:])
	if err != nil {
		panic(err)
	}

	// Write Length
	err = binary.Write(&buf, binary.LittleEndian, m.Length)
	if err != nil {
		panic(err)
	}

	// Write Checksum
	err = binary.Write(&buf, binary.NativeEndian, m.Checksum)
	if err != nil {
		panic(err)
	}

	// Write Payload
	_, err = buf.Write(m.Payload)
	if err != nil {
		panic(err)
	}

	return buf.Bytes()
}

func FromBytes(b []byte) (*BtcMessage, error) {
	// Parse Magic
	magic, err := parseFromBytesLE[uint32](b[:4])
	if err != nil {
		return nil, err
	}
	// Parse Command
	command := [12]byte{}
	copy(command[:], b[4:16])

	// Parse Length
	length, err := parseFromBytesLE[uint32](b[16:20])
	if err != nil {
		return nil, err
	}

	// Parse Checksum
	checksum, err := parseFromBytesLE[uint32](b[20:24])
	if err != nil {
		return nil, err
	}
	return &BtcMessage{
		Magic:    magic,
		Command:  command,
		Length:   length,
		Checksum: checksum,
		Payload:  b[24:],
	}, nil
}

func parseFromBytesLE[T any](b []byte) (T, error) {
	var value T
	err := binary.Read(bytes.NewReader(b), binary.LittleEndian, &value)
	if err != nil {
		return value, errors.New("failed to parse from bytes")
	}
	return value, nil
}
