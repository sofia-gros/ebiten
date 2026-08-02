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
	group      *camera.Group
	mainCam    *camera.Camera
	miniCam    *camera.Camera
	playerX    float64
	playerY    float64
	angle      float64
	monoShader *ebiten.Shader
	gridImg    *ebiten.Image
}

func (g *Game) Init() {
	g.playerX = 320
	g.playerY = 240

	// グレースケール Custom Shader のコンパイル
	shaderCode := []byte(`
		package main

		func Fragment(dstPos vec4, srcPos vec2, color vec4) vec4 {
			c := imageSrc0At(srcPos)
			gray := c.r*0.299 + c.g*0.587 + c.b*0.114
			return vec4(gray, gray, gray, c.a)
		}
	`)
	if s, err := ebiten.NewShader(shaderCode); err == nil {
		g.monoShader = s
	}

	// 1. メインカメラ (ZIndex: 0)
	g.mainCam = camera.New(640, 480)
	g.mainCam.SetZIndex(0)

	// 2. ミニマップカメラ (ZIndex: 10 最前面、右上 160x120 枠)
	g.miniCam = camera.New(640, 480)
	g.miniCam.SetViewport(460, 20, 160, 120)
	g.miniCam.SetZoom(0.3)
	g.miniCam.SetZIndex(10)

	g.group = camera.NewGroup(g.mainCam, g.miniCam)

	// 1000x1000 のグリッド画像を事前生成 (チカチカ防止)
	g.gridImg = ebiten.NewImage(1000, 1000)
	g.gridImg.Fill(color.RGBA{15, 15, 25, 255})
	gridColor := color.RGBA{50, 50, 80, 255}

	for x := 0; x <= 1000; x += 40 {
		for y := 0; y < 1000; y++ {
			g.gridImg.Set(x, y, gridColor)
		}
	}
	for y := 0; y <= 1000; y += 40 {
		for x := 0; x < 1000; x++ {
			g.gridImg.Set(x, y, gridColor)
		}
	}
}



func (g *Game) Update() error {
	dt := 1.0 / 60.0

	// カメラの正面ベクトル・右ベクトルの取得 (相対移動用)
	fwdX, fwdY := g.mainCam.Forward()
	rtX, rtY := g.mainCam.Right()

	speed := 3.0

	// WASD キー: カメラの向きに沿った相対移動
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		g.playerX += fwdX * speed
		g.playerY += fwdY * speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.playerX -= fwdX * speed
		g.playerY -= fwdY * speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		g.playerX -= rtX * speed
		g.playerY -= rtY * speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		g.playerX += rtX * speed
		g.playerY += rtY * speed
	}

	// K キー: 画面振動 Shake
	if inpututil.IsKeyJustPressed(ebiten.KeyK) {
		g.mainCam.Shake(15, 0.5)
	}

	// R / E キー: カメラ回転
	if ebiten.IsKeyPressed(ebiten.KeyR) {
		g.angle += 0.02
		g.mainCam.SetRotation(g.angle)
	}
	if ebiten.IsKeyPressed(ebiten.KeyE) {
		g.angle -= 0.02
		g.mainCam.SetRotation(g.angle)
	}

	// T キー: カスタムシェーダー (モノクロ shader) のトグル
	if inpututil.IsKeyJustPressed(ebiten.KeyT) {
		if g.mainCam.HasShader() {
			g.mainCam.ClearShader()
		} else if g.monoShader != nil {
			g.mainCam.SetShader(g.monoShader)
		}
	}

	// カメラ位置の滑らかな移動追従
	g.mainCam.MoveTo(g.playerX, g.playerY, 0.1)
	g.miniCam.SetPos(g.playerX, g.playerY)

	g.group.Update(dt)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// グループ一括描画 (ZIndex順)
	// Render に渡した描画クロージャ内で、純粋なワールド座標 (0〜1000) のまま普通に DrawImage するだけで
	// カメラの追従・ズーム・回転・画面振動・ビューポート切り抜き・カスタムシェーダーがすべて全自動適用される！
	g.group.Render(screen, func(cam *camera.Camera, target *ebiten.Image) {
		lineColor := color.RGBA{50, 50, 80, 255}

		// 生ワールド空間への直線描画 (cam.Apply や特別なハックは一切不要)
		for x := 0; x <= 1000; x += 40 {
			vectorLine(target, float64(x), 0, float64(x), 1000, lineColor)
		}
		for y := 0; y <= 1000; y += 40 {
			vectorLine(target, 0, float64(y), 1000, float64(y), lineColor)
		}

		// プレイヤー描画 (純粋なワールド座標 g.playerX, g.playerY で描画)
		pImg := ebiten.NewImage(24, 24)
		pImg.Fill(color.RGBA{50, 220, 100, 255})
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(g.playerX-12, g.playerY-12)
		target.DrawImage(pImg, opts)
	})

	ebitenutil.DebugPrint(screen, "[CAMERA DEMO]\nMove: WASD (Relative to rotation) | Rotate: R/E | Shake: K | Shader Toggle: T\nTop-Right: Mini-map (Multi-Camera Viewport)")
}


func vectorLine(target *ebiten.Image, x1, y1, x2, y2 float64, c color.Color) {
	w := int(math.Max(2, math.Abs(x2-x1)))
	h := int(math.Max(2, math.Abs(y2-y1)))
	img := ebiten.NewImage(w, h)
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
