package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/sofia-gros/ebiten/scene"
	"github.com/sofia-gros/ebiten/tween"
	"github.com/sofia-gros/ebiten/ui"
)

type TitleScene struct {
	ctx        *GameContext
	titleScale float64
}

func (s *TitleScene) Init(sCtx *scene.Context) {
	// タイトルアニメーション (tween)
	s.ctx.TweenGroup.New(tween.Option{
		Start:    0.8,
		End:      1.1,
		Duration: 1.5,
		Ease:     tween.EaseInOutQuad,
		Yoyo:     true,
		Loop:     -1,
	}).OnUpdate(func(val float64) {
		s.titleScale = val
	}).Play()

	// UI ボタン配置 (ui)
	btnStart := ui.NewButton(ui.ButtonOption{Text: "START GAME", Width: 160, Height: 40})
	btnStart.SetPos(240, 280)
	btnStart.OnClick(func() {
		sCtx.Start(&FarmScene{ctx: s.ctx})
	})

	btnLoad := ui.NewButton(ui.ButtonOption{Text: "LOAD GAME", Width: 160, Height: 40})
	btnLoad.SetPos(240, 340)
	btnLoad.OnClick(func() {
		// セーブデータ読込試行 (save)
		var sd FarmSaveData
		if err := s.ctx.SaveMgr.Slot(1).Load("farm", "json", &sd); err == nil {
			s.ctx.PlayerX = sd.PlayerX
			s.ctx.PlayerY = sd.PlayerY
			s.ctx.Gold = sd.Gold
			if sd.TilledMap != nil {
				s.ctx.TilledTiles = sd.TilledMap
			}
		}
		sCtx.Start(&FarmScene{ctx: s.ctx})
	})

	s.ctx.UIRoot.Clear()
	s.ctx.UIRoot.Add(btnStart)
	s.ctx.UIRoot.Add(btnLoad)
}

func (s *TitleScene) Update(sCtx *scene.Context) error {
	dt := 1.0 / 60.0
	s.ctx.TweenGroup.Update(dt)
	s.ctx.UIRoot.Update()

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		sCtx.Start(&FarmScene{ctx: s.ctx})
	}
	return nil
}

func (s *TitleScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{35, 75, 55, 255})

	// タイトルロゴ描画
	ebitenutil.DebugPrintAt(screen, "==========================================", 160, 120)
	ebitenutil.DebugPrintAt(screen, "        SPROUT LANDS FARM GAME           ", 160, 140)
	ebitenutil.DebugPrintAt(screen, "==========================================", 160, 160)

	ebitenutil.DebugPrintAt(screen, "Press SPACE or Click START to Begin", 190, 220)
	ebitenutil.DebugPrintAt(screen, "Credit: Assets by Cup Nooble (sprout-lands)", 170, 440)

	s.ctx.UIRoot.Draw(screen)
}

func (s *TitleScene) Destroy(sCtx *scene.Context) {
	s.ctx.UIRoot.Clear()
}
