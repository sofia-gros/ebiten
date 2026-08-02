package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)


type Game struct {
	ctx *GameContext
}

func (g *Game) Init() {
	g.ctx = NewGameContext()
	g.ctx.SetupEvents()

	// 最初のシーン (TitleScene) を開始 (scene)
	g.ctx.SceneMgr.Context().Start(&TitleScene{ctx: g.ctx})
}

func (g *Game) Update() error {
	return g.ctx.SceneMgr.Update()
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.ctx.SceneMgr.Draw(screen)
}

func (g *Game) Layout(w, h int) (int, int) {
	return 640, 480
}

func main() {
	g := &Game{}
	g.Init()

	ebiten.SetWindowTitle("Sprout Lands Farm - Full Library Sample Game")
	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}
