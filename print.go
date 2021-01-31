package main

import (
	"machine"
)

func print(str string) {
	machine.UART0.Write([]byte(str + "\r"))
}

func printx(str string) {
	machine.UART0.Write([]byte(str))
}

func println(str string) {
	machine.UART0.Write([]byte(str + "\n" + "\r"))
}

func printInit() {
	machine.UART0.SetBaudRate(115200)
	
	print("\n\n")
	println("RV-STM32 Computer Project")
	println("2021 - Danilo Fragoso <danilo.fragoso@gmail.com>")
	println("--------------------------------------------------")
}