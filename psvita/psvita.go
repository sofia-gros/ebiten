// Package psvita は、PSVita の画面解像度、各種ハードウェア API（タッチ・センサー・システム情報）、
// および CPU/GPU クロック制御などの Core 層アクセスを提供する Ebitengine 用拡張ライブラリです。
package psvita

import "github.com/hajimehoshi/ebiten/v2"

const (
	// ScreenWidth は PSVita の画面幅 (960 ピクセル) です。
	ScreenWidth = 960

	// ScreenHeight は PSVita の画面高さ (544 ピクセル) です。
	ScreenHeight = 544
)

// Vector3 は 3 軸ベクトルを表します。
type Vector3 struct {
	X float64
	Y float64
	Z float64
}

// Button は PSVita の物理ボタンおよび拡張ボタンを表します。
type Button int

const (
	ButtonSelect Button = iota
	ButtonStart
	ButtonUp
	ButtonRight
	ButtonDown
	ButtonLeft
	ButtonL1
	ButtonR1
	ButtonTriangle
	ButtonCircle
	ButtonCross
	ButtonSquare
	ButtonPS
	ButtonRearTouchPad
)

// IsVita は現在の実行環境が PSVita 実機 (または WASM ビルド) かどうかを返します。
func IsVita() bool {
	return isVitaTarget()
}

// IsButtonPressed は指定された PSVita ボタンが押されているかを判定します。
func IsButtonPressed(button Button) bool {
	return isButtonPressedImpl(button)
}

// MapButtonToStandardGamepad は PSVita のボタンを Ebitengine の StandardGamepadButton に変換します。
func MapButtonToStandardGamepad(button Button) ebiten.StandardGamepadButton {
	switch button {
	case ButtonSelect:
		return ebiten.StandardGamepadButtonCenterLeft
	case ButtonStart:
		return ebiten.StandardGamepadButtonCenterRight
	case ButtonUp:
		return ebiten.StandardGamepadButtonLeftTop
	case ButtonRight:
		return ebiten.StandardGamepadButtonLeftRight
	case ButtonDown:
		return ebiten.StandardGamepadButtonLeftBottom
	case ButtonLeft:
		return ebiten.StandardGamepadButtonLeftLeft
	case ButtonL1:
		return ebiten.StandardGamepadButtonFrontTopLeft
	case ButtonR1:
		return ebiten.StandardGamepadButtonFrontTopRight
	case ButtonTriangle:
		return ebiten.StandardGamepadButtonRightTop
	case ButtonCircle:
		return ebiten.StandardGamepadButtonRightRight
	case ButtonCross:
		return ebiten.StandardGamepadButtonRightBottom
	case ButtonSquare:
		return ebiten.StandardGamepadButtonRightLeft
	default:
		return ebiten.StandardGamepadButtonMax
	}
}
