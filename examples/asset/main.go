package main

import (
	"embed"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/sofia-gros/ebiten/asset"
)

//go:embed assets/*
var assetsFS embed.FS

type ItemConfig struct {
	Name  string `json:"name"`
	Power int    `json:"power"`
}

type Game struct {
	m         *asset.Manager
	statusMsg string
}

func (g *Game) Init() {
	g.m = asset.NewManager()
	g.m.SetFS(assetsFS)

	// 1. Image 登録
	g.m.AddImage("hero_icon", "assets/hero.png")

	// 2. Sprite 登録 (コマ割り 32x32)
	g.m.AddSprite("player_sheet", "assets/player_sheet.png", asset.Option{
		FrameWidth:  32,
		FrameHeight: 32,
	})

	// 3. Data (JSON) 登録
	g.m.AddData("item_data", "assets/item.json")

	// 一括ロード実行
	g.m.Load()

	// JSON を構造体へ展開
	item, err := asset.DataAs[ItemConfig](g.m, "item_data")
	if err == nil {
		g.statusMsg = fmt.Sprintf("Asset Manager Loaded! Item JSON: %s (Power %d)", item.Name, item.Power)
	} else {
		g.statusMsg = "DataAs Error: " + err.Error()
	}
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, fmt.Sprintf("[ASSET DEMO]\nStatus: %s", g.statusMsg))

	// 単体画像の描画
	if img, err := g.m.Image("hero_icon"); err == nil && img != nil {
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(50, 100)
		screen.DrawImage(img, opts)
		ebitenutil.DebugPrintAt(screen, "Image: hero_icon", 50, 140)
	}

	// スプライトシートの特定のコマ (Frame 0) の描画
	if sheet, err := g.m.Sprite("player_sheet"); err == nil && sheet != nil {
		if frame := sheet.Frame(0); frame != nil {
			opts := &ebiten.DrawImageOptions{}
			opts.GeoM.Translate(250, 100)
			screen.DrawImage(frame, opts)
			ebitenutil.DebugPrintAt(screen, "SpriteSheet: Frame 0", 250, 140)
		}
	}
}


func (g *Game) Layout(w, h int) (int, int) {
	return 640, 480
}

func main() {
	g := &Game{}
	g.Init()

	ebiten.SetWindowTitle("asset Demo - AddImage, AddSprite, AddData, DataAs")
	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}
