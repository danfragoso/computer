package main

func bootrom() []byte {
	return []byte{0x93, 0xE, 0x50, 0x0, 0x13, 0xF, 0x50, 0x2, 0xB3, 0xF, 0xDF, 0x1, }
}