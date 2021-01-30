package main

import (
	"encoding/binary"
	"strconv"
)

type Emulator struct {
	CPU    *CPU
	Memory *Memory
}

type CPU struct {
	Registers      [32]uint32
	ProgramCounter uint32

	Emulator *Emulator
}

type Memory struct {
	Size uint32
	Emulator *Emulator
}

func (mem *Memory) Load(address uint32) uint32 {
	return binary.LittleEndian.Uint32(readBytes([]byte{
        byte(0xff & (address >> 16)),
		byte(0xff & (address >> 8)),
		byte(0xff & address),
	}, 4))
}

func CreateEmulator(initData []byte) *Emulator {
	emulator := &Emulator{}

	emulator.CreateMemory(initData)
	emulator.CreateCPU()

	return emulator
}

func (emulator *Emulator) CreateMemory(initData []byte) {
	writeBytes([]byte{0x0, 0x0, 0x0}, initData)

	emulator.Memory = &Memory{
		Size:     uint32(len(initData)),
		Emulator: emulator,
	}
}

func (emulator *Emulator) CreateCPU() {
	emulator.CPU = &CPU{
		Emulator: emulator,
	}

	emulator.CPU.Registers[2] = emulator.Memory.Size
	emulator.CPU.Registers[0] = 0
}

func (emulator *Emulator) Run() {
	println("Starting execution")
	println("----------------------")

	for emulator.CPU.ProgramCounter < emulator.Memory.Size {
		print("PC: " + strconv.FormatUint(uint64(emulator.CPU.ProgramCounter/4), 10))
		instruction := emulator.CPU.Fetch()
		emulator.CPU.ProgramCounter += 4
		emulator.CPU.Execute(instruction)
	}

	print("\n")
	println("----------------------")
	println("Execution ended")
	emulator.DumpRegisters()
}

func (emulator *Emulator) DumpRegisters() {
	printx("\nREGISTERS {")

	for idx, register := range emulator.CPU.Registers {
		if idx == 0 || idx == 4 || idx == 8 || idx == 12 ||
			idx == 16 || idx == 20 || idx == 24 || idx == 28 {
			println("")
		}

		printx("x" + padInt(idx) + "=0x" + padUintHex(register) + "; ")
	}

	println("\n\r}")
}

func (cpu *CPU) Fetch() uint32 {
	return cpu.Emulator.Memory.Load(cpu.ProgramCounter)
}

func (cpu *CPU) Execute(inst uint32) {
	OPCODE := inst & 0x0000007f

	RD := (inst & 0x00000f80) >> 7
	RS1 := (inst & 0x000f8000) >> 15
	RS2 := (inst & 0x01f00000) >> 20

	switch OPCODE {
	case 0x13: // ADDI
		IMM := uint32((inst & 0xfff00000) >> 20)
		cpu.Registers[RD] = cpu.Registers[RS1] + IMM

	case 0x33: // ADD
		cpu.Registers[RD] = cpu.Registers[RS1] + cpu.Registers[RS2]

	default:
		println("Opcode '0x" + strconv.FormatUint(uint64(OPCODE), 16) + "' not implemented!")
	}

}

func padInt(n int) string {
	var pad string
	if n < 10 {
		pad = "0"
	}

	return pad + strconv.Itoa(n)
}

func padUintHex(n uint32) string {
	var pad string
	if n < 16 {
		pad = "0"
	}

	return pad + strconv.FormatUint(uint64(n), 16)
}