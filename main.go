package main

func main() {
	printInit()
	
	setupRAM()
	
	//Loads data from serial to RAM and return the numeber of bytes loaded.
	payloadSize := loadFromSerial()

	emulator := CreateEmulator(payloadSize)
	emulator.Run()	
}
