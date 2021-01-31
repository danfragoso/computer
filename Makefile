flash:
	@tinygo flash -target=bluepill .

build:
	@echo "Building..."
	@tinygo build -target=bluepill -o kernel.bin .

clean:
	@rm -f kernel.bin core DAC*

asm:
	cd code && rm -f *.o && rm -f *.bin
	cd code && riscv32-unknown-elf-as -o add-addi.o add-addi.s
	cd code && riscv32-unknown-elf-objcopy -O binary add-addi.o add-addi.bin

	cd code && riscv32-unknown-elf-as -o li-add.o li-add.s
	cd code && riscv32-unknown-elf-objcopy -O binary li-add.o li-add.bin