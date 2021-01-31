package main

import (
	"github.com/jacobsa/go-serial/serial"
	"encoding/binary"
	"encoding/hex"
	"io/ioutil"
	"fmt"
	"time"
	"os"
)

func main() {
	payload, err := ioutil.ReadFile(os.Args[1])
    if err != nil {
        fmt.Print(err)
		os.Exit(1)
	}
	
	port, err := serial.Open(serial.OpenOptions{
		PortName: "/dev/ttyUSB0",
		BaudRate: 115200,
		DataBits: 8,
		StopBits: 1,
		MinimumReadSize: 4,
	})

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	byteCount := make([]byte, 4, 4)
	binary.BigEndian.PutUint32(byteCount, uint32(len(payload)))

	port.Write([]byte{0x21})
	port.Write(byteCount)

	byteCounter := 0

	fmt.Printf("Loading %d bytes...", len(payload))
	time.Sleep(time.Second)

	payloadBuffer := make([]byte, 4, 4)
	for byteCounter < len(payload) {
		payloadBuffer = []byte{0, 0, 0, 0}
		for i, _ := range payloadBuffer {
			if byteCounter+i < len(payload) {
				payloadBuffer[i] = payload[byteCounter+i]
			}
		}

		// payloadBuffer = []byte{
		// 	payload[byteCounter],
		// 	payload[byteCounter+1],
		// 	payload[byteCounter+2],
		// 	payload[byteCounter+3],
		// }


		port.Write(payloadBuffer)
		fmt.Print("DATA_0x" + hex.EncodeToString(payloadBuffer))

		buff := make([]byte, 4, 4)
		port.Read(buff)
		fmt.Printf("; ADDR_0x" + hex.EncodeToString(buff)+"\n")

		byteCounter+=4
	}
}