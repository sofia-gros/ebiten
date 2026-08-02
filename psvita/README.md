# psvita

PlayStation Vita 用の各種 API（前面・背面タッチ、6軸ジャイロ・加速度センサー、システム情報）および Core 層アクセス（CPU/GPU クロック制御）を提供する Ebitengine 用拡張ライブラリです。

`pspdev-go` 方式（**Go -> WASM -> w2c2 -> VitaSDK -> VPK**）のパイプラインによる実機バイナリ (.vpk) のコンパイルに対応し、PC/Web環境での開発時は各種入力のエミュレーション（フォールバック）を標準でサポートします。

---

## 提供機能・API

### 1. タッチ入力 (`psvita.TouchFront`, `psvita.TouchBack`)
- 前面マルチタッチパネルおよび背面タッチパッドの独立した座標・識別ID・タッチ圧力を取得します。

### 2. 6軸モーションセンサー (`psvita.Gyro`, `psvita.Accelerometer`)
- 3軸ジャイロスコープの角速度 (rad/s) および 3軸加速度 (G) を取得します。

### 3. システム・バッテリー情報 (`psvita.BatteryLevel`, `psvita.IsCharging`)
- バッテリー残量 (%) や充電状態、省電力モードの設定・取得に対応します。

### 4. Core 層アクセス (`psvita.SetCPUClock`, `psvita.GetMemoryInfo`)
- CPU/GPU クロック周波数（333MHz 標準 / 444MHz ブースト）の動的変更やメモリ割り当て状況を参照できます。

---

## 使い方・サンプルコード

本リポジトリ内の `physics` や `camera` ライブラリと組み合わせて使用できます。

```go
package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"ebiten/psvita"
)

type Game struct{}

func (g *Game) Update() error {
	// 前面・背面タッチの取得
	frontTouches := psvita.TouchFront()
	backTouches := psvita.TouchBack()

	// 6軸ジャイロの取得
	gyro := psvita.Gyro()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{20, 20, 30, 255})
	ebitenutil.DebugPrint(screen, fmt.Sprintf("Battery: %d%%", psvita.BatteryLevel()))
}

func (g *Game) Layout(w, h int) (int, int) {
	return psvita.ScreenWidth, psvita.ScreenHeight
}

func main() {
	// Core層: CPUクロックを444MHzにブースト
	psvita.SetCPUClock(psvita.ClockFrequencyMax)

	ebiten.SetWindowSize(psvita.ScreenWidth, psvita.ScreenHeight)
	ebiten.RunGame(&Game{})
}
```

---

## PSVita 向けコンパイル手順 (.vpk の作成)

### 必要ツール
1. [TinyGo](https://tinygo.org/)
2. [w2c2](https://github.com/wasmerio/w2c2) (WASM -> C Transpiler)
3. [vitasdk](https://vitasdk.org/) (ARM-vita-eabi-gcc ツールチェーン)

### ビルド実行
```bash
cd psvita
make
```

`build/game.vpk` が作成されますので、VitaShell 等経由で PSVita に転送してインストールしてください。
