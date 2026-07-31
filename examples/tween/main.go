package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/sofia-gros/ebiten/tween"
)

type Game struct {
	group   *tween.Group
	box1X   float64
	box2X   float64
	box3X   float64
	tw1     *tween.Tween
	tw2     *tween.Tween
	tw3     *tween.Tween
	status  string
}

func (g *Game) Init() {
	g.group = tween.NewGroup()

	// 1. Bounce イージング + 往復リピート
	g.tw1 = g.group.New(tween.Option{
		Start:    50.0,
		End:      500.0,
		Duration: 2.0,
		Ease:     tween.EaseOutBounce,
		Yoyo:     true,
		Loop:     -1,
	}).OnUpdate(func(val float64) {
		g.box1X = val
	}).Play()

	// 2. Back イージング
	g.tw2 = g.group.New(tween.Option{
		Start:    50.0,
		End:      500.0,
		Duration: 1.5,
		Ease:     tween.EaseOutBack,
		Yoyo:     true,
		Loop:     -1,
	}).OnUpdate(func(val float64) {
		g.box2X = val
	}).Play()

	// 3. Elastic イージング
	g.tw3 = g.group.New(tween.Option{
		Start:    50.0,
		End:      500.0,
		Duration: 2.5,
		Ease:     tween.EaseOutElastic,
		Yoyo:     true,
		Loop:     -1,
	}).OnUpdate(func(val float64) {
		g.box3X = val
	}).Play()

	g.status = "Playing all tweens in Group"
}

func (g *Game) Update() error {
	dt := 1.0 / 60.0

	// P キー: グループ全体の一括一時停止 / 再開
	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		g.group.PauseAll()
		g.status = "Paused all group tweens"
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		g.group.ResumeAll()
		g.status = "Resumed all group tweens"
	}

	g.group.Update(dt)
	return nil
}


func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, fmt.Sprintf("[TWEEN DEMO]\nPress P: Pause All | Press R: Resume All\nStatus: %s", g.status))

	// Box 1 (Bounce)
	drawBox(screen, g.box1X, 120, color.RGBA{220, 80, 80, 255}, "Bounce")

	// Box 2 (Back)
	drawBox(screen, g.box2X, 220, color.RGBA{80, 220, 80, 255}, "Back")

	// Box 3 (Elastic)
	drawBox(screen, g.box3X, 320, color.RGBA{80, 140, 240, 255}, "Elastic")
}

func drawBox(target *ebiten.Image, x, y float64, c color.Color, label string) {
	img := ebiten.NewImage(40, 40)
	img.Fill(c)
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(x, y)
	target.DrawImage(img, opts)

	ebitenutil.DebugPrintAt(target, label, int(x), int(y-18))
}

func (g *Game) Layout(w, h int) (int, int) {
	return 640, 480
}

func main() {
	g := &Game{}
	g.Init()

	ebiten.SetWindowTitle("tween Demo - Easings, Group, Yoyo, PauseAll, ResumeAll")
	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}
