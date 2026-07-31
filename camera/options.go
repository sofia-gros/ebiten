package camera

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

// Viewport はカメラの画面出力領域（ミニマップ表示や画面分割等）を定義する構造体です。
type Viewport struct {
	X      float64 // 画面上の表示開始 X 座標(px)
	Y      float64 // 画面上の表示開始 Y 座標(px)
	Width  float64 // 表示幅(px)
	Height float64 // 表示高さ(px)
}

// Bounds はカメラ移動を制限するワールド境界領域を表す構造体です。
type Bounds struct {
	MinX float64
	MinY float64
	MaxX float64
	MaxY float64
}

// ShaderOption はカメラにセットする Custom Shader のパラメータを保持する構造体です。
type ShaderOption struct {
	Uniforms map[string]any
	Images   [4]*ebiten.Image
}

// rectToRectangle は image.Rectangle へ変換する内部ヘルパーです。
func rectToRectangle(x, y, w, h int) image.Rectangle {
	return image.Rect(x, y, x+w, y+h)
}
