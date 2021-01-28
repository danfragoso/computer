flash:
	@tinygo flash -target=bluepill .

build:
	@echo "Building..."
	@tinygo build -target=bluepill -o kernel.bin .

clean:
	@rm -f kernel.bin core DAC*
