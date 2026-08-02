#!/bin/bash
set -e

cd "$(dirname "$0")"

BUILD_DIR="build"
W2C_OUT="$BUILD_DIR/w2c_gen"
TOOLS_DIR="$PWD/.tools"

echo ""
echo "╔══════════════════════════════════════════════════╗"
echo "║  🎮 Ebitengine PS Vita Build (WSL + wasm-opt)    ║"
echo "╚══════════════════════════════════════════════════╝"
echo ""

mkdir -p "$BUILD_DIR" "$W2C_OUT" "$TOOLS_DIR"

# ─── 1. Tool Downloads (cached - only runs once) ───────────────────

# 1a. VitaSDK
if [ ! -d "$TOOLS_DIR/vitasdk/bin" ]; then
    echo "📦 Downloading VitaSDK for Linux..."
    VITASDK_URL=$(curl -s https://api.github.com/repos/vitasdk/autobuilds/releases | grep browser_download_url | grep linux | head -n 1 | cut -d '"' -f 4)
    curl -sL -o "$TOOLS_DIR/vitasdk.tar.bz2" "$VITASDK_URL"
    echo "  📂 Extracting VitaSDK..."
    mkdir -p "$TOOLS_DIR/vitasdk"
    tar -xjf "$TOOLS_DIR/vitasdk.tar.bz2" -C "$TOOLS_DIR/vitasdk" --strip-components=1
    rm "$TOOLS_DIR/vitasdk.tar.bz2"
    echo "  ✅ VitaSDK installed"
fi
export VITASDK="$TOOLS_DIR/vitasdk"
export PATH="$VITASDK/bin:$PATH"

# 1b. WABT (wasm2c)
if [ ! -d "$TOOLS_DIR/wabt/bin" ]; then
    echo "📦 Downloading WABT (wasm2c) for Linux..."
    curl -sL -o "$TOOLS_DIR/wabt.tar.gz" https://github.com/WebAssembly/wabt/releases/download/1.0.34/wabt-1.0.34-ubuntu.tar.gz
    echo "  📂 Extracting WABT..."
    mkdir -p "$TOOLS_DIR/wabt"
    tar -xzf "$TOOLS_DIR/wabt.tar.gz" -C "$TOOLS_DIR/wabt" --strip-components=1
    rm "$TOOLS_DIR/wabt.tar.gz"
    echo "  ✅ WABT installed"
fi
WASM2C="$TOOLS_DIR/wabt/bin/wasm2c"

# 1c. Go (Standard Go)
if [ ! -d "$TOOLS_DIR/go/bin" ]; then
    echo "📦 Downloading standard Go for Linux..."
    curl -sL -o "$TOOLS_DIR/go.tar.gz" https://go.dev/dl/go1.22.5.linux-amd64.tar.gz
    echo "  📂 Extracting Go..."
    mkdir -p "$TOOLS_DIR/go"
    tar -xzf "$TOOLS_DIR/go.tar.gz" -C "$TOOLS_DIR/go" --strip-components=1
    rm "$TOOLS_DIR/go.tar.gz"
    echo "  ✅ Go installed"
fi
export PATH="$TOOLS_DIR/go/bin:$PATH"
export GOROOT="$TOOLS_DIR/go"

# 1d. Binaryen (wasm-opt)
if [ ! -f "$TOOLS_DIR/binaryen/bin/wasm-opt" ]; then
    echo "📦 Downloading Binaryen (wasm-opt) for Linux..."
    curl -sL -o "$TOOLS_DIR/binaryen.tar.gz" https://github.com/WebAssembly/binaryen/releases/download/version_116/binaryen-version_116-x86_64-linux.tar.gz
    echo "  📂 Extracting Binaryen..."
    mkdir -p "$TOOLS_DIR/binaryen"
    tar -xzf "$TOOLS_DIR/binaryen.tar.gz" -C "$TOOLS_DIR/binaryen" --strip-components=1
    rm "$TOOLS_DIR/binaryen.tar.gz"
    echo "  ✅ Binaryen installed"
fi
WASM_OPT="$TOOLS_DIR/binaryen/bin/wasm-opt"

echo ""

# ─── 2. Go → WebAssembly (standard Go) + wasm-opt ─────────────────

STEP_START=$(date +%s)
echo "⚙️  [Step 1/5] Compiling Go → WebAssembly..."
GOOS=js GOARCH=wasm go build -ldflags="-s -w" -o "$BUILD_DIR/game.wasm" .
WASM_RAW=$(stat -c%s "$BUILD_DIR/game.wasm" 2>/dev/null || echo "0")
echo "  📊 Raw WASM: $((WASM_RAW / 1024))KB"

echo -n "  🔧 Optimizing WASM with wasm-opt... "
OPT_START=$(date +%s)
"$WASM_OPT" -O3 --all-features --strip-debug --strip-producers \
    --dce --remove-unused-names --remove-unused-module-elements \
    "$BUILD_DIR/game.wasm" -o "$BUILD_DIR/game_opt.wasm" &
OPT_PID=$!
SPINNER=('⠋' '⠙' '⠹' '⠸' '⠼' '⠴' '⠦' '⠧' '⠇' '⠏')
SI=0
while kill -0 $OPT_PID 2>/dev/null; do
    NOW=$(date +%s)
    printf "\r  🔧 Optimizing WASM with wasm-opt... %s %ds" "${SPINNER[$SI]}" $((NOW - OPT_START))
    SI=$(( (SI + 1) % 10 ))
    sleep 0.2
done
wait $OPT_PID
printf "\n"

mv "$BUILD_DIR/game_opt.wasm" "$BUILD_DIR/game.wasm"
WASM_OPT_SIZE=$(stat -c%s "$BUILD_DIR/game.wasm" 2>/dev/null || echo "0")
REDUCTION=$(( (WASM_RAW - WASM_OPT_SIZE) * 100 / WASM_RAW ))
STEP_END=$(date +%s)
echo "  ✅ Optimized WASM: $((WASM_OPT_SIZE / 1024))KB (${REDUCTION}% smaller) ($((STEP_END - STEP_START))s)"
echo ""

# ─── 3. WebAssembly → C (via wasm2c) ──────────────────────────────

STEP_START=$(date +%s)
echo "⚙️  [Step 2/5] Transpiling WASM → C (wasm2c --num-outputs=64)..."
rm -rf "$W2C_OUT"
mkdir -p "$W2C_OUT"
"$WASM2C" "$BUILD_DIR/game.wasm" --num-outputs=64 -o "$W2C_OUT/game_wasm.c"

# Copy WABT runtime files
cp "$TOOLS_DIR/wabt"/include/wasm-rt*.h "$W2C_OUT/"
cp "$TOOLS_DIR/wabt"/share/wabt/wasm2c/wasm-rt*.h "$W2C_OUT/" 2>/dev/null || true
cp "$TOOLS_DIR/wabt"/share/wabt/wasm2c/wasm-rt*.c "$W2C_OUT/" 2>/dev/null || true
cp "$TOOLS_DIR/wabt"/share/wabt/wasm2c/wasm-rt*.inc "$W2C_OUT/" 2>/dev/null || true

# Patch for PS Vita (no POSIX signals, no mmap)
sed -i 's/sigsetjmp(buf, 1)/setjmp(buf)/g' "$W2C_OUT/wasm-rt.h"
sed -i 's/siglongjmp(buf, val)/longjmp(buf, val)/g' "$W2C_OUT/wasm-rt.h"
for f in "$W2C_OUT"/wasm-rt*.c; do
    sed -i 's|#include <sys/mman.h>|#if WASM_RT_USE_MMAP\n#include <sys/mman.h>\n#endif|g' "$f"
done

# Count generated C output size
TOTAL_C_SIZE=0
for f in "$W2C_OUT"/game_wasm*.c; do
    [ -f "$f" ] && TOTAL_C_SIZE=$((TOTAL_C_SIZE + $(stat -c%s "$f")))
done
TOTAL_C_MB=$((TOTAL_C_SIZE / 1048576))
STEP_END=$(date +%s)
echo "  ✅ Generated C code: ${TOTAL_C_MB}MB ($((STEP_END - STEP_START))s)"
echo ""

# ─── 4. Generate param.sfo ─────────────────────────────────────────

PARAM_SFO="$BUILD_DIR/param.sfo"
if [ ! -f "$PARAM_SFO" ]; then
    "$VITASDK/bin/vita-mksfoex" -s TITLE_ID="EBITEN001" "PSVita Ebiten Game" "$PARAM_SFO"
fi

# ─── 5. Compile C → ARM (VitaSDK GCC, parallel) ───────────────────

NPROC=$(nproc)
echo "⚙️  [Step 3/5] Compiling C → ARM (GCC, ${NPROC} cores parallel)..."

BRIDGE_C="../../psvita/bridge/vita_bridge.c"
BRIDGE_INCLUDE="../../psvita/bridge"
GCC="$VITASDK/bin/arm-vita-eabi-gcc"
GCC_FLAGS="-O0 -c -pipe -DWASM_RT_USE_MMAP=0 -I$BRIDGE_INCLUDE -I$W2C_OUT"

# Collect all C files
ALL_C=("$BRIDGE_C")
for f in "$W2C_OUT"/wasm-rt*.c; do [ -f "$f" ] && ALL_C+=("$f"); done
for f in "$W2C_OUT"/game_wasm_*.c; do [ -f "$f" ] && ALL_C+=("$f"); done

TOTAL=${#ALL_C[@]}
BUILD_START=$(date +%s)

# Calculate total bytes for accurate ETA
TOTAL_BYTES=0
declare -A FILE_SIZES
for src in "${ALL_C[@]}"; do
    sz=$(stat -c%s "$src" 2>/dev/null || echo "1")
    FILE_SIZES["$src"]=$sz
    TOTAL_BYTES=$((TOTAL_BYTES + sz))
done

# Progress bar function (byte-weighted ETA)
COMPILED_BYTES=0
show_bar() {
    local done=$1 total=$2 now=$(date +%s)
    local elapsed=$((now - BUILD_START))
    local pct=$((done * 100 / total))
    local w=30 filled=$((pct * 30 / 100)) empty=$((30 - pct * 30 / 100))
    local bar=""
    for ((i=0; i<filled; i++)); do bar+="█"; done
    for ((i=0; i<empty; i++)); do bar+="░"; done
    local eta="--:--"
    if [ "$COMPILED_BYTES" -gt 0 ] && [ "$COMPILED_BYTES" -lt "$TOTAL_BYTES" ]; then
        local bps=$((COMPILED_BYTES / (elapsed + 1)))
        local rem_bytes=$((TOTAL_BYTES - COMPILED_BYTES))
        local rem_s=$((rem_bytes / (bps + 1)))
        eta=$(printf "%dm%02ds" $((rem_s/60)) $((rem_s%60)))
    elif [ "$done" -eq "$total" ]; then
        eta="done!"
    fi
    printf "\r  %s %3d%% [%d/%d] ⏱️%dm%02ds ETA:%s  " \
        "$bar" "$pct" "$done" "$total" $((elapsed/60)) $((elapsed%60)) "$eta"
}

# Launch parallel jobs
PIDS=()
SRC_NAMES=()
SRC_SIZES_ORDERED=()
OBJ_FILES=()

for src in "${ALL_C[@]}"; do
    bn=$(basename "$src")
    obj="$BUILD_DIR/${bn%.c}.o"
    OBJ_FILES+=("$obj")

    $GCC $GCC_FLAGS "$src" -o "$obj" 2>"$BUILD_DIR/${bn%.c}.log" &
    PIDS+=($!)
    SRC_NAMES+=("$bn")
    SRC_SIZES_ORDERED+=(${FILE_SIZES["$src"]})

    # Limit to NPROC parallel jobs
    while [ $(jobs -r | wc -l) -ge $NPROC ]; do
        DONE=0
        COMPILED_BYTES=0
        for i in "${!PIDS[@]}"; do
            if ! kill -0 "${PIDS[$i]}" 2>/dev/null; then
                DONE=$((DONE + 1))
                COMPILED_BYTES=$((COMPILED_BYTES + ${SRC_SIZES_ORDERED[$i]}))
            fi
        done
        show_bar "$DONE" "$TOTAL"
        sleep 1
    done
done

# Wait for all remaining
while true; do
    DONE=0
    COMPILED_BYTES=0
    for i in "${!PIDS[@]}"; do
        if ! kill -0 "${PIDS[$i]}" 2>/dev/null; then
            DONE=$((DONE + 1))
            COMPILED_BYTES=$((COMPILED_BYTES + ${SRC_SIZES_ORDERED[$i]}))
        fi
    done
    show_bar "$DONE" "$TOTAL"
    [ "$DONE" -eq "$TOTAL" ] && break
    sleep 1
done
echo ""

# Check failures
FAILED=0
for i in "${!PIDS[@]}"; do
    wait "${PIDS[$i]}" || {
        echo "  ❌ FAILED: ${SRC_NAMES[$i]}"
        cat "$BUILD_DIR/${SRC_NAMES[$i]%.c}.log" 2>/dev/null
        FAILED=$((FAILED + 1))
    }
done
[ "$FAILED" -gt 0 ] && { echo "  ❌ $FAILED file(s) failed."; exit 1; }

BUILD_END=$(date +%s)
echo "  ✅ Compiled $TOTAL files in $((BUILD_END - BUILD_START))s"
echo ""

# ─── 6. Link → ELF → VELF → FSELF → VPK ──────────────────────────

STEP_START=$(date +%s)
echo "⚙️  [Step 4/5] Linking → ELF..."
"$GCC" -O0 -Wl,-q "${OBJ_FILES[@]}" \
    -lSceTouch_stub -lSceMotion_stub -lScePower_stub -lSceCtrl_stub \
    -o "$BUILD_DIR/game.elf"

echo "⚙️  [Step 5/5] Packaging → VPK..."
"$VITASDK/bin/vita-elf-create" "$BUILD_DIR/game.elf" "$BUILD_DIR/game.velf"
"$VITASDK/bin/vita-make-fself" -c "$BUILD_DIR/game.velf" "$BUILD_DIR/eboot.bin"
"$VITASDK/bin/vita-pack-vpk" -s "$PARAM_SFO" -b "$BUILD_DIR/eboot.bin" "$BUILD_DIR/game.vpk"
STEP_END=$(date +%s)

VPK_SIZE=$(stat -c%s "$BUILD_DIR/game.vpk" 2>/dev/null || echo "0")
VPK_MB=$((VPK_SIZE / 1048576))

echo ""
echo "╔══════════════════════════════════════════════════╗"
echo "║  ✅ BUILD COMPLETE                               ║"
echo "║  📦 build/game.vpk  (${VPK_MB}MB)                ║"
echo "╚══════════════════════════════════════════════════╝"
echo ""
