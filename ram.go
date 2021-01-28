package main

import (
	"machine"
)

// RAM ASSIGNED GPIO PINS
const RAM_CE = machine.PB9	
const RAM_VCC = machine.PA0

// LY68L6400 SPI OP_CODES
const RAM_RSTEN_CMD = 0x66
const RAM_RST_CMD   = 0x99
const RAM_WRITE_CMD = 0x02
const RAM_READ_CMD  = 0x03

func setupRAM() {
	print("Setting up RAM...")

	RAM_CE.Configure(machine.PinConfig{Mode: machine.PinOutput})
	RAM_Standby() //CE Active low

	RAM_VCC.Configure(machine.PinConfig{Mode: machine.PinOutput})
	RAM_VCC.Set(true) //Enable RAM VCC

	machine.SPI1.Configure(machine.SPIConfig{
		LSBFirst: false,
		Frequency: 8000000,
		SCK: machine.SPI0_SCK_PIN, SDO: machine.SPI0_SDO_PIN, SDI: machine.SPI0_SDI_PIN,
	})

	RAM_Wakeup() 
	machine.SPI1.Tx([]byte{RAM_RSTEN_CMD, RAM_RST_CMD}, nil) //Send REST SEQ [0x66, 0x99]
	RAM_Standby()

	writeBytes([]byte{0xff, 0xff, 0xff}, []byte{0x12, 0x34, 0x56, 0x78})
	bytesRead := readBytes([]byte{0xff, 0xff, 0xff}, 4)
	
	if bytesRead[0] == 0x12 && bytesRead[1] == 0x34 && bytesRead[2] == 0x56 && bytesRead[3] == 0x78 {
		println("RAM Ok!          ")
	} else {
		println("RAM RW FATAL ERR; RW NO MATCH;")
	}
}
	
func writeBytes(address []byte, data []byte) {
	RAM_Wakeup()
	machine.SPI1.Tx(append([]byte{RAM_WRITE_CMD, address[0], address[1], address[2]}, data...), nil)
	RAM_Standby()
}

func readBytes(address []byte, bytesCount int) []byte {
	RAM_Wakeup()
	readInst := append([]byte{RAM_READ_CMD, address[0], address[1], address[2]}, make([]byte, bytesCount, bytesCount)...)
	machine.SPI1.Tx(readInst, readInst)	
	RAM_Standby()

	return readInst[4:]
}

func RAM_Standby() {
	RAM_CE.High()
}

func RAM_Wakeup() {
	RAM_CE.Low()
}
