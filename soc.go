package main

func startSOC() {
	emulator := CreateEmulator(bootrom())
	emulator.Run()
}