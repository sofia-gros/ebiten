# PSVita (ebiten/psvita) Build Makefile
# Target: Go -> WASM -> w2c2 -> VitaSDK -> VPK

TARGET = game
GO_SRC = .
WASM_BUILD = build/$(TARGET).wasm
C_BUILD_DIR = build/w2c2_out

.PHONY: all clean wasm w2c2 vpk

all: vpk

wasm:
	@mkdir -p build
	GOOS=wasip1 GOARCH=wasm tinygo build -o $(WASM_BUILD) $(GO_SRC)

w2c2: wasm
	@mkdir -p $(C_BUILD_DIR)
	w2c2 $(WASM_BUILD) $(C_BUILD_DIR)/game_wasm.c

vpk: w2c2
	@echo "Compiling for PSVita with arm-vita-eabi-gcc..."
	arm-vita-eabi-gcc -O2 -Ibridge -I$(C_BUILD_DIR) \
		bridge/vita_bridge.c $(C_BUILD_DIR)/*.c \
		-lSceTouch_stub -lSceMotion_stub -lScePower_stub -lSceCtrl_stub \
		-o build/$(TARGET).velf
	vita-make-fself build/$(TARGET).velf build/eboot.bin
	vita-pack-vpk -s param.sfo -b build/eboot.bin build/$(TARGET).vpk
	@echo "Build Complete: build/$(TARGET).vpk"

clean:
	rm -rf build
