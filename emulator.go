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
	PC uint32

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

func CreateEmulator(payloadSize uint32) *Emulator {
	emulator := &Emulator{}

	emulator.CreateMemory(payloadSize)
	emulator.CreateCPU()

	return emulator
}

func (emulator *Emulator) CreateMemory(memorySize uint32) {
	if memorySize == 0 {
		memorySize = bootROM_Size()
		writeBytes([]byte{0x0, 0x0, 0x0}, bootROM())
	}
	
	emulator.Memory = &Memory{
		Size:     memorySize,
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

	for emulator.CPU.PC < emulator.Memory.Size {
		print("PC: " + strconv.FormatUint(uint64(emulator.CPU.PC/4), 10))
		instruction := emulator.CPU.Fetch()
		emulator.CPU.PC += 4
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
	return cpu.Emulator.Memory.Load(cpu.PC)
}

func (cpu *CPU) Execute(inst uint32) {
	OPCODE := (inst & 0b00000000000000000000000001111111)
	RD     := (inst & 0b00000000000000000000111110000000) >> 7
	
	RS1    := (inst & 0b00000000000011111000000000000000) >> 15
	RS2    := (inst & 0b00000001111100000000000000000000) >> 20

	FUNCT3 := (inst & 0b00000000000000000111000000000000) >> 12

	switch OPCODE {
	case 0b0110111: // LUI
		cpu.Registers[RD] = inst & 0b11111111111111111111000000000000
	
	case 0b0010111: // AUIPC
		cpu.Registers[RD] = cpu.PC + (inst & 0b11111111111111111111000000000000)

	case 0b1101111: // JAL
		cpu.Registers[RD] = cpu.PC + 4
		cpu.PC += (inst >> 11) & 0b10000000000000000000000000000000 | // imm[20]
				  (inst)       & 0b00000000000011111111000000000000 | // imm[19:12]
				  (inst >> 9)  & 0b00000000000000000000100000000000 | // imm[11]
				  (inst >> 20) & 0b00000000000000000000011111111110   // imm[10:1]
	
	case 0b0010011: // ADDI, SLTI, SLTIU, XORI, ORI, ANDI, SLLI, SRLI, SRAI
		switch FUNCT3 {
		case 0b000: // ADDI
			cpu.Registers[RD] = cpu.Registers[RS1] + I_Type_IMM(inst)

		case 0b010: // SLTI 
			if cpu.Registers[RS1] < I_Type_IMM(inst) {
				cpu.Registers[RD] = 1
			} else {
				cpu.Registers[RD] = 0
			}
		}

	case 0b00110011: // ADD
		cpu.Registers[RD] = cpu.Registers[RS1] + cpu.Registers[RS2]

	default:
		println("Opcode '0b" + strconv.FormatUint(uint64(OPCODE), 2) + "' not implemented!")
	}

}

func I_Type_IMM(inst uint32) uint32 {
	return (inst & 0b11111111111100000000000000000000) >> 20
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