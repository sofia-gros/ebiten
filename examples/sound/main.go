package main

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/sofia-gros/ebiten/sound"
)

type Game struct {
	m          *sound.Manager
	seBytes    []byte
	bgmBytes   []byte
	bgmTrack   *sound.Track
	posTrack   *sound.Track
	targetX    float64
	statusMsg  string
}

func (g *Game) Init() {
	g.m = sound.NewManager(44100)
	g.targetX = 100

	// ダミー音声を動的生成 (正弦波 440Hz / 880Hz BGM & SE)
	g.seBytes = generateSineWave(880, 0.2)
	g.bgmBytes = generateSineWave(440, 2.0)

	g.statusMsg = "Press 1: Play SE | Press 2: Play Loop BGM | Press 3: CrossFade | Press 4: Stop Fade | Press 5: Play 2D Positional"
}

func (g *Game) Update() error {
	// 1 キー: SE 再生
	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		g.m.Play(g.seBytes, sound.Option{Type: sound.TypeSE})
		g.statusMsg = "Played SE!"
	}

	// 2 キー: BGM 再生
	if inpututil.IsKeyJustPressed(ebiten.Key2) {
		g.bgmTrack = g.m.Play(g.bgmBytes, sound.Option{Type: sound.TypeBGM, Loop: true, Volume: 0.7})
		g.statusMsg = "Playing BGM (Loop)..."
	}

	// 3 キー: CrossFade
	if inpututil.IsKeyJustPressed(ebiten.Key3) {
		newBgm := generateSineWave(554, 2.0)
		g.m.CrossFade(newBgm, 1.5, sound.Option{Type: sound.TypeBGM, Loop: true})
		g.statusMsg = "CrossFading to new BGM (1.5s)..."
	}

	// 4 キー: フェードアウト停止
	if inpututil.IsKeyJustPressed(ebiten.Key4) {
		g.m.StopType(sound.TypeBGM, 1.0)
		g.statusMsg = "Fading out BGM (1.0s)..."
	}

	// 5 キー: 2D ポジショナルサウンド再生 (最大距離: 400.0)
	if inpututil.IsKeyJustPressed(ebiten.Key5) {
		g.posTrack = g.m.PlayAt(g.seBytes, g.targetX, 240, 320, 240, 400.0, sound.Option{Type: sound.TypeSE, Loop: true})
		g.statusMsg = "Playing 2D Positional Sound! Move Target with Arrow Keys."
	}

	// 左右キーで音源位置移動 (2D Positional)
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.targetX -= 4
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.targetX += 4
	}

	if g.posTrack != nil {
		g.posTrack.SetPosition(g.targetX, 240, 320, 240, 400.0)
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, fmt.Sprintf("[SOUND DEMO]\nControls:\n[1] Play SE\n[2] Play BGM\n[3] CrossFade BGM\n[4] Fade Out Stop\n[5] 2D Positional Audio (Move with Left/Right Arrows)\n\nStatus: %s", g.statusMsg))

	if g.posTrack != nil {
		// 聴き手 (プレイヤー x:320)
		pImg := ebiten.NewImage(16, 16)
		pImg.Fill(color.RGBA{100, 250, 100, 255})
		pOpts := &ebiten.DrawImageOptions{}
		pOpts.GeoM.Translate(320-8, 240-8)
		screen.DrawImage(pImg, pOpts)
		ebitenutil.DebugPrintAt(screen, "Listener", 300, 260)

		// 音源 (Target x:targetX)
		tImg := ebiten.NewImage(16, 16)
		tImg.Fill(color.RGBA{250, 100, 100, 255})
		tOpts := &ebiten.DrawImageOptions{}
		tOpts.GeoM.Translate(g.targetX-8, 240-8)
		screen.DrawImage(tImg, tOpts)
		ebitenutil.DebugPrintAt(screen, "Sound Source", int(g.targetX)-30, 260)
	}
}

func (g *Game) Layout(w, h int) (int, int) {
	return 640, 480
}

func generateSineWave(freq float64, durationSec float64) []byte {
	sampleRate := 44100
	numSamples := int(float64(sampleRate) * durationSec)
	bytes := make([]byte, numSamples*4) // 16bit stereo = 4 bytes per sample

	for i := 0; i < numSamples; i++ {
		t := float64(i) / float64(sampleRate)
		val := int16(math.Sin(2*math.Pi*freq*t) * 16000)

		// L channel
		bytes[i*4] = byte(val)
		bytes[i*4+1] = byte(val >> 8)
		// R channel
		bytes[i*4+2] = byte(val)
		bytes[i*4+3] = byte(val >> 8)
	}
	return bytes
}

func main() {
	g := &Game{}
	g.Init()

	ebiten.SetWindowTitle("sound Demo - Play, CrossFade, PlayAt Positional Audio")
	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}
