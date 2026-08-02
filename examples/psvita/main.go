package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"ebiten/camera"
	"ebiten/physics"
	"ebiten/physics/adapters/arcade"
	"ebiten/psvita"
)

type DemoGame struct {
	world      physics.World
	ball       physics.Body
	cam        *camera.Camera
	touchFront []psvita.TouchPoint
	touchBack  []psvita.TouchPoint
	gyro       psvita.Vector3
}

func NewDemoGame() *DemoGame {
	// Arcade 物理エンジンの初期化 (重力 9.8)
	world := arcade.NewWorld(0, 9.8*50)

	// 円形物理ボディの作成
	ball := world.CreateBody(physics.BodyOptions{
		Shape:    physics.ShapeCircle,
		Radius:   20,
		Mass:     1.0,
		Restitution: 0.8, // 跳ね返り
	})
	ball.SetPos(physics.Vec2{X: psvita.ScreenWidth / 2, Y: 100})

	// 画面底面に静的床オブジェクトを作成
	floor := world.CreateBody(physics.BodyOptions{
		IsStatic: true,
		Width:    psvita.ScreenWidth,
		Height:   40,
	})
	floor.SetPos(physics.Vec2{X: psvita.ScreenWidth / 2, Y: psvita.ScreenHeight - 20})

	cam := camera.NewCamera(psvita.ScreenWidth, psvita.ScreenHeight)

	return &DemoGame{
		world: world,
		ball:  ball,
		cam:   cam,
	}
}

func (g *DemoGame) Update() error {
	// 1. PSVita 入力・センサ・システム情報の取得
	g.touchFront = psvita.TouchFront()
	g.touchBack = psvita.TouchBack()
	g.gyro = psvita.Gyro()

	// 2. 背面タッチパットでボールをジャンプさせる
	if len(g.touchBack) > 0 || psvita.IsButtonPressed(psvita.ButtonCross) {
		g.ball.ApplyImpulse(physics.Vec2{X: 0, Y: -300})
	}

	// 3. 6軸ジャイロの傾きで物理重力方向を変更
	g.world.SetGravity(physics.Vec2{
		X: g.gyro.Y * 200,
		Y: 9.8*50 + g.gyro.X*200,
	})

	// 4. 物理ワールドの更新
	g.world.Update(1.0 / 60.0)

	// 5. カメラをボールに滑らかに追従
	g.cam.SetPos(g.ball.Pos().X, g.ball.Pos().Y)

	return nil
}

func (g *DemoGame) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{15, 20, 35, 255})

	// 床の描画
	ebitenutil.DrawRect(screen, 0, float64(psvita.ScreenHeight-40), psvita.ScreenWidth, 40, color.RGBA{60, 80, 120, 255})

	// ボール物理ボディの描画
	bpos := g.ball.Pos()
	ebitenutil.DrawRect(screen, bpos.X-20, bpos.Y-20, 40, 40, color.RGBA{255, 100, 100, 255})

	// 前面タッチ座標の描画
	for _, t := range g.touchFront {
		ebitenutil.DrawRect(screen, float64(t.X-10), float64(t.Y-10), 20, 20, color.RGBA{0, 255, 150, 200})
	}

	// 背面タッチ座標の描画
	for _, t := range g.touchBack {
		ebitenutil.DrawRect(screen, float64(t.X-15), float64(t.Y-15), 30, 30, color.RGBA{255, 200, 0, 180})
	}

	// 情報テキスト表示
	info := fmt.Sprintf(
		"=== PSVita + Ebitengine Physics Demo ===\n"+
			"Battery: %d%% | Charging: %v\n"+
			"Gyro (Tilt): X=%.2f Y=%.2f\n"+
			"Front Touches: %d | Back Touches (Jump): %d\n"+
			"Controls: [Cross/BackTouch]: Jump | [Arrow Keys]: Tilt Gyro",
		psvita.BatteryLevel(), psvita.IsCharging(),
		g.gyro.X, g.gyro.Y,
		len(g.touchFront), len(g.touchBack),
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
	ebiten.SetWindowTitle("PSVita Ebitengine Physics Demo")

	if err := ebiten.RunGame(NewDemoGame()); err != nil {
		log.Fatal(err)
	}
}
