package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/sofia-gros/ebiten/camera"
)



type Game struct {
	group    *camera.Group
	mainCam  *camera.Camera
	miniCam  *camera.Camera
	playerX  float64
	playerY  float64
	angle    float64
}

func (g *Game) Init() {
	g.playerX = 320
	g.playerY = 240

	// 1. メインカメラ (ZIndex: 0)
	g.mainCam = camera.New(640, 480)
	g.mainCam.SetZIndex(0)

	// 2. ミニマップカメラ (ZIndex: 10 最前面、右上 160x120 枠)
	g.miniCam = camera.New(640, 480)
	g.miniCam.SetViewport(460, 20, 160, 120)
	g.miniCam.SetZoom(0.3)
	g.miniCam.SetZIndex(10)

	g.group = camera.NewGroup(g.mainCam, g.miniCam)
}


func (g *Game) Update() error {
	dt := 1.0 / 60.0

	// 矢印キーでプレイヤー移動
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.playerX -= 3
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.playerX += 3
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.playerY -= 3
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		g.playerY += 3
	}

	// S キー: 画面振動 Shake
	if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		g.mainCam.Shake(15, 0.5)
	}


	// R キー: カメラ回転
	if ebiten.IsKeyPressed(ebiten.KeyR) {
		g.angle += 0.02
		g.mainCam.SetRotation(g.angle)
	}

	// カメラ位置の滑らかな移動追従
	g.mainCam.MoveTo(g.playerX, g.playerY, 0.1)
	g.miniCam.SetPos(g.playerX, g.playerY)

	g.group.Update(dt)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// グループ一括描画 (ZIndex順)
	g.group.Render(screen, func(cam *camera.Camera, target *ebiten.Image) {
		// グリッド線描画
		for x := 0; x < 1000; x += 40 {
			vectorLine(target, float64(x), 0, float64(x), 1000, color.RGBA{40, 40, 60, 255})
		}
		for y := 0; y < 1000; y += 40 {
			vectorLine(target, 0, float64(y), 1000, float64(y), color.RGBA{40, 40, 60, 255})
		}

		// プレイヤー描画 (緑の円)
		pImg := ebiten.NewImage(24, 24)
		pImg.Fill(color.RGBA{50, 220, 100, 255})
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(g.playerX-12, g.playerY-12)
		target.DrawImage(pImg, opts)
	})

	ebitenutil.DebugPrint(screen, "[CAMERA DEMO]\nMove: Arrow Keys | S: Shake Screen | R: Rotate Camera\nTop-Right: Mini-map (Multi-Camera Viewport)")
}

func vectorLine(target *ebiten.Image, x1, y1, x2, y2 float64, c color.Color) {
	img := ebiten.NewImage(int(math.Max(1, math.Abs(x2-x1))), int(math.Max(1, math.Abs(y2-y1))))
	img.Fill(c)
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(x1, y1)
	target.DrawImage(img, opts)
}

func (g *Game) Layout(w, h int) (int, int) {
	return 640, 480
}

func main() {
	g := &Game{}
	g.Init()

	ebiten.SetWindowTitle("camera Demo - MoveTo, Shake, Rotation, Multi-Camera Viewport")
	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}
