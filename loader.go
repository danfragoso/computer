package main

import (
	"encoding/binary"
	"machine"
	"strconv"
	"time"
)

func loadFromSerial() uint32 {
	println("Waiting for loader... \n\rPress any key other than '!' to boot from bootrom")
	
	rByte, _ := machine.UART0.ReadByte()
	for rByte == 0x00 {
		rByte, _ = machine.UART0.ReadByte()
	}

	if rByte != 0x21 {
		return 0
	}

	byteBuffer := make([]byte, 4, 4)
	byteCounter := 0

	for byteCounter < 4 {
		if machine.UART0.Buffered() > 0 {
			byteBuffer[byteCounter], _ = machine.UART0.ReadByte()
			byteCounter++
		}
	}

	realAddress := make([]byte, 4, 4)
	var payloadSize, bytesWritten uint32
	payloadSize = binary.BigEndian.Uint32(byteBuffer)

	time.Sleep(time.Millisecond * 1000)

	byteCounter = 0
	for bytesWritten < payloadSize {
		for byteCounter < 4 {
			if machine.UART0.Buffered() > 0 {
				time.Sleep(time.Millisecond)
				byteBuffer[byteCounter], _ = machine.UART0.ReadByte()
				
				byteCounter++
				bytesWritten++
			}
		}
		
		byteCounter = 0

		writeBytes(realAddress[1:], byteBuffer)
		machine.UART0.Write(realAddress)
		binary.BigEndian.PutUint32(realAddress, bytesWritten)		
	}
	
	println("Loaded "+strconv.Itoa(int(payloadSize))+" bytes... \n\rPress '!' to boot from RAM")
	
	rByte, _ = machine.UART0.ReadByte()
	for rByte == 0x00 {
		rByte, _ = machine.UART0.ReadByte()

		if rByte != 0x00 && rByte != 0x21 {
			println("Loaded "+strconv.Itoa(int(payloadSize))+" bytes... \n\rPress '!' to boot from RAM")
			rByte = 0x00
		}
	}

	return payloadSize
}