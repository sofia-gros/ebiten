package psvita

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

// Gyro は 6 軸ジャイロセンサーの角速度ベクトル (rad/s) を取得します。
func Gyro() Vector3 {
	return getGyroImpl()
}

// Accelerometer は 3 軸加速度センサーの加速度ベクトル (G) を取得します。
func Accelerometer() Vector3 {
	return getAccelerometerImpl()
}

func getGyroImpl() Vector3 {
	// PC上でのエミュレーション（矢印キーやマウス移動に基づく模擬ジャイロ）
	var x, y float64
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		y -= 0.5
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		y += 0.5
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		x -= 0.5
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		x += 0.5
	}
	return Vector3{X: x, Y: y, Z: 0.0}
}

func getAccelerometerImpl() Vector3 {
	// 擬似重力ベクトル (静止状態 1.0G) + 入力傾き
	g := getGyroImpl()
	accX := math.Max(-1.0, math.Min(1.0, g.Y))
	accY := math.Max(-1.0, math.Min(1.0, g.X+1.0)) // 標準傾き
	return Vector3{X: accX, Y: accY, Z: -0.8}
}
