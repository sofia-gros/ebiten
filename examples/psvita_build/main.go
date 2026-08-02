

package main

import (
	"fmt"
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"github.com/sofia-gros/ebiten/camera"
	"github.com/sofia-gros/ebiten/physics"
	"github.com/sofia-gros/ebiten/physics/adapters/arcade"

	"ebiten/psvita"
)

type DemoGame struct {
	physManager *physics.Manager
	playerBody  physics.Body
	floorBody   physics.Body
	cam         *camera.Camera
	touchFront  []psvita.TouchPoint
	touchBack   []psvita.TouchPoint
	gyro        psvita.Vector3
}

func NewDemoGame() *DemoGame {
	pm := physics.NewManager()
	pm.SetWorld(arcade.NewWorld())
	pm.SetGravity(0, 9.8*20)

	player := pm.CreateBody(physics.BodyOptions{
		Type: physics.BodyTypeDynamic,
		X:    psvita.ScreenWidth / 2,
		Y:    100,
		Shapes: []physics.ShapeDef{
			{Shape: physics.BoxShape{Width: 32, Height: 32}},
		},
		Restitution: 0.8,
	})

	floor := pm.CreateBody(physics.BodyOptions{
		Type: physics.BodyTypeStatic,
		X:    psvita.ScreenWidth / 2,
		Y:    psvita.ScreenHeight - 20,
		Shapes: []physics.ShapeDef{
			{Shape: physics.BoxShape{Width: psvita.ScreenWidth, Height: 40}},
		},
	})
	_ = floor

	cam := camera.New(float64(psvita.ScreenWidth), float64(psvita.ScreenHeight))

	return &DemoGame{
		physManager: pm,
		playerBody:  player,
		cam:         cam,
	}
}

func (g *DemoGame) Update() error {
	// 1. PSVita 入力・センサ・システム情報の取得
	g.touchFront = psvita.TouchFront()
	g.touchBack = psvita.TouchBack()
	g.gyro = psvita.Gyro()

	vx, vy := g.playerBody.Velocity()

	// 2. PC操作: WASDキー / 方向キーでの直接移動およびジャンプ
	moveSpeed := 180.0
	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyLeft) || psvita.IsButtonPressed(psvita.ButtonLeft) {
		vx = -moveSpeed
	} else if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyRight) || psvita.IsButtonPressed(psvita.ButtonRight) {
		vx = moveSpeed
	} else {
		// 【摩擦処理】キー入力がない場合は急速に水平速度を減衰（摩擦）させて滑り続けを解消
		vx *= 0.82
		if math.Abs(vx) < 0.1 {
			vx = 0
		}
	}

	// ジャンプ操作 (Wキー / Spaceキー / Kキー / Crossボタン / 背面タッチ)
	if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeySpace) || ebiten.IsKeyPressed(ebiten.KeyK) ||
		psvita.IsButtonPressed(psvita.ButtonCross) || len(g.touchBack) > 0 {
		// 接地時または上昇開始時のジャンプ
		if math.Abs(vy) < 5.0 {
			vy = -320
		}
	}

	g.playerBody.SetVelocity(vx, vy)

	// 3. 6軸ジャイロの傾きで物理重力方向を変更
	g.physManager.SetGravity(g.gyro.Y*200, 9.8*30+g.gyro.X*200)

	// 4. 物理ワールドの更新
	g.physManager.World().Step(1.0 / 60.0)

	// 5. カメラをボールに追従
	px, py := g.playerBody.Position()
	g.cam.SetPos(px, py)

	return nil
}

func (g *DemoGame) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{15, 20, 35, 255})

	// 物理デバッグ描画
	g.physManager.DrawDebug(screen)

	// 前面タッチ座標の描画
	for _, t := range g.touchFront {
		ebitenutil.DrawRect(screen, float64(t.X-10), float64(t.Y-10), 20, 20, color.RGBA{0, 255, 150, 200})
	}

	// 背面タッチ座標の描画
	for _, t := range g.touchBack {
		ebitenutil.DrawRect(screen, float64(t.X-15), float64(t.Y-15), 30, 30, color.RGBA{255, 200, 0, 180})
	}

	// PSVita 全情報の詳細画面表示
	accel := psvita.Accelerometer()
	mem, _ := psvita.GetMemoryInfo()

	// 押下ボタン一覧の構築
	pressedButtons := ""
	buttons := []struct {
		name string
		b    psvita.Button
	}{
		{"Select", psvita.ButtonSelect}, {"Start", psvita.ButtonStart},
		{"Up", psvita.ButtonUp}, {"Right", psvita.ButtonRight}, {"Down", psvita.ButtonDown}, {"Left", psvita.ButtonLeft},
		{"L1", psvita.ButtonL1}, {"R1", psvita.ButtonR1},
		{"Triangle", psvita.ButtonTriangle}, {"Circle", psvita.ButtonCircle}, {"Cross", psvita.ButtonCross}, {"Square", psvita.ButtonSquare},
	}
	for _, btn := range buttons {
		if psvita.IsButtonPressed(btn.b) {
			pressedButtons += btn.name + " "
		}
	}
	if pressedButtons == "" {
		pressedButtons = "None"
	}

	info := fmt.Sprintf(
		"=== PSVita Full Status & Hardware Info ===\n"+
			"[System/Power] Battery: %d%% | Charging: %v | PowerMode: %v\n"+
			"[Core/Memory]  CPU Frequency: 444MHz (Boost) | Free Main Mem: %d MB / %d MB\n"+
			"[6-Axis Gyro]  X: %.2f | Y: %.2f | Z: %.2f\n"+
			"[3-Axis Accel] X: %.2f | Y: %.2f | Z: %.2f\n"+
			"[Touch Panel]  Front Touches: %d | Back Touches: %d\n"+
			"[Buttons]      Pressed: %s\n"+
			"--------------------------------------------------\n"+
			"[Controls] Move: [W/A/S/D] or [Arrow Keys] | Jump: [Space/W/Cross] | RightClick: RearTouch",
		psvita.BatteryLevel(), psvita.IsCharging(), psvita.CurrentPowerMode(),
		mem.FreeMainMemoryBytes/(1024*1024), mem.TotalMainMemoryBytes/(1024*1024),
		g.gyro.X, g.gyro.Y, g.gyro.Z,
		accel.X, accel.Y, accel.Z,
		len(g.touchFront), len(g.touchBack),
		pressedButtons,
	)
	ebitenutil.DebugPrint(screen, info)
}

func (g *DemoGame) Layout(w, h int) (int, int) {
	return psvita.ScreenWidth, psvita.ScreenHeight
}

func main() {
	// Core層: CPUクロックを444MHzにブースト
	if err := psvita.SetCPUClock(psvita.ClockFrequencyMax); err != nil {
		log.Println("SetCPUClock Warning:", err)
	}

	ebiten.SetWindowSize(psvita.ScreenWidth, psvita.ScreenHeight)
	ebiten.SetWindowTitle("PSVita Ebiten Physics Demo (psvita_build)")

	if err := ebiten.RunGame(NewDemoGame()); err != nil {
		log.Fatal(err)
	}
}
