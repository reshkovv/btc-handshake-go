package main

import (
	"crypto/sha256"
	"encoding/binary"
	"log"
	"net"
)

func main() {
	address := net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 8333}
	version_msg := NewVersionMessage(60002, address)

	payload := version_msg.ToBytes()

	commandBytes := [12]byte{}
	copy(commandBytes[:], "version")

	h1 := sha256.Sum256(payload)
	h2 := sha256.Sum256(h1[:])
	checksum := binary.LittleEndian.Uint32(h2[:4])

	btc_msg := BtcMessage{
		Magic:    MAGIC_NUMBER,
		Command:  commandBytes,
		Length:   uint32(len(payload)),
		Checksum: checksum,
		Payload:  payload,
	}

	// Connect to the Bitcoin Core node
	conn, err := net.Dial("tcp", address.String())
	if err != nil {
		log.Println("Failed to connect to the node:", err)
		return
	}
	defer conn.Close()

	// Send the "version" message
	_, err = conn.Write(btc_msg.ToBytes())
	if err != nil {
		log.Println("Failed to send the version message:", err)
		return
	}
	log.Println("Version message sent")

	// Receive the "version" message from the node
	log.Println("Waiting for verack answer...")
	mes, err := readMessage(conn)
	if err != nil {
		log.Println("Failed to receive the version message:", err)
		return
	}

	btc_msg_resp, err := FromBytes(mes)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("Verack answer received")
	if string(btc_msg_resp.Command[:7]) != "version" {
		log.Fatal("Verack answer is not correct:", string(btc_msg_resp.Command[:]))
	}

	log.Println("Handshake completed successfully")
}

func readMessage(conn net.Conn) ([]byte, error) {
	recBuff := make([]byte, 1024)
	_, err := conn.Read(recBuff)
	if err != nil {
		return nil, err
	}
	return recBuff, nil
}
