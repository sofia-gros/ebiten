package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sofia-gros/ebiten/asset"
	"github.com/sofia-gros/ebiten/camera"
	"github.com/sofia-gros/ebiten/emit"
	pad "github.com/sofia-gros/ebiten/pad/input"
	"github.com/sofia-gros/ebiten/physics"
	arcade "github.com/sofia-gros/ebiten/physics/adapters/arcade"
	"github.com/sofia-gros/ebiten/save"
	"github.com/sofia-gros/ebiten/scene"
	"github.com/sofia-gros/ebiten/sound"
	"github.com/sofia-gros/ebiten/tween"
	"github.com/sofia-gros/ebiten/ui"
)

// --- イベント定義 (emit) ---

type CropPlantedEvent struct {
	TileX, TileY int
	CropType     string
}

type CropHarvestedEvent struct {
	TileX, TileY int
	GoldEarned   int
}

type GoldChangedEvent struct {
	NewAmount int
}

// --- セーブデータ定義 (save) ---

type FarmSaveData struct {
	PlayerX   float64        `json:"player_x"`
	PlayerY   float64        `json:"player_y"`
	Gold      int            `json:"gold"`
	TilledMap map[string]int `json:"tilled_map"`
}

// --- ツールの列挙 ---
type ToolType int

const (
	ToolHoe ToolType = iota
	ToolWater
	ToolSeed
	ToolHarvest
)

// --- 共通ゲームコンテキスト (GameContext) ---
type GameContext struct {
	// 全10個のライブラリマネージャー
	SceneMgr   *scene.Manager
	AssetMgr   *asset.Manager
	SoundMgr   *sound.Manager
	SaveMgr    *save.Manager
	PhysicsMgr *physics.Manager
	Emitter    *emit.Emitter
	PadMgr     *pad.Input
	TweenGroup *tween.Group
	UIRoot     *ui.Container

	// メイン＆ミニマップカメラ
	CameraGroup *camera.Group
	MainCam     *camera.Camera
	MiniCam     *camera.Camera

	// ゲーム状態
	PlayerX     float64
	PlayerY     float64
	PlayerSpeed float64
	PlayerDir   int // 0: Down, 1: Up, 2: Left, 3: Right
	AnimTimer   float64

	Gold         int
	SelectedTool ToolType

	// 耕された農地タイル (キー: "x,y", 値: 成長段階)
	TilledTiles map[string]int

	// 2D サウンドトラックハンドラ
	BGMTrack *sound.Track
	SEBytes  map[string][]byte

	// UI エレメント
	GoldLabel   *ui.Label
	StatusLabel *ui.Label
}

func NewGameContext() *GameContext {
	ctx := &GameContext{
		PlayerX:      640,
		PlayerY:      480,
		PlayerSpeed:  2.5,
		Gold:         100,
		SelectedTool: ToolHoe,
		TilledTiles:  make(map[string]int),
		SEBytes:      make(map[string][]byte),
	}

	// 1. emit (イベントバス)
	ctx.Emitter = emit.New()

	// 2. physics (Arcade AABB 物理エンジン)
	ctx.PhysicsMgr = physics.NewManager()
	ctx.PhysicsMgr.SetWorld(arcade.NewWorld())


	// 3. sound (オーディオマネージャー)
	ctx.SoundMgr = sound.NewManager(44100)

	// 4. save (セーブデータマネージャー)
	ctx.SaveMgr = save.NewManager("SproutLandsFarm", save.Option{
		Format:     save.FormatJSON,
		EncryptKey: "SproutLandsSecret123",
	})

	// 5. asset (アセットマネージャー)
	ctx.AssetMgr = asset.NewManager()

	// 6. camera (メインカメラ & ミニマップ Viewport カメラ)
	ctx.MainCam = camera.New(640, 480)
	ctx.MainCam.SetZIndex(0)

	ctx.MiniCam = camera.New(640, 480)
	ctx.MiniCam.SetViewport(460, 20, 160, 120)
	ctx.MiniCam.SetZoom(0.2)
	ctx.MiniCam.SetZIndex(10)

	ctx.CameraGroup = camera.NewGroup(ctx.MainCam, ctx.MiniCam)

	// 7. tween (トゥイーンアニメーション)
	ctx.TweenGroup = tween.NewGroup()

	// 8. pad (入力抽象化)
	ctx.PadMgr = pad.NewInput()
	ctx.PadMgr.For(pad.DefaultController).BindKey(1, ebiten.KeySpace)
	ctx.PadMgr.For(pad.DefaultController).BindKey(1, ebiten.KeyZ)
	ctx.PadMgr.For(pad.DefaultController).BindKey(1, ebiten.KeyEnter)

	// 9. ui (GUI UI コンテナ)
	ctx.UIRoot = ui.NewContainer()

	// 10. scene (シーンマネージャー)
	ctx.SceneMgr = scene.NewManager(640, 480)

	return ctx
}

func (ctx *GameContext) SetupEvents() {
	// イベント購読 (emit)
	emit.On[CropPlantedEvent](ctx.Emitter, func(e CropPlantedEvent) {
		ctx.SoundPlaySE("click")
		ctx.MainCam.Shake(4, 0.2)
	})

	emit.On[CropHarvestedEvent](ctx.Emitter, func(e CropHarvestedEvent) {
		ctx.Gold += e.GoldEarned
		emit.Emit(ctx.Emitter, GoldChangedEvent{NewAmount: ctx.Gold})

		ctx.SoundPlaySE("coinsCling")
		ctx.MainCam.Shake(8, 0.3)
	})

	emit.On[GoldChangedEvent](ctx.Emitter, func(e GoldChangedEvent) {
		if ctx.GoldLabel != nil {
			ctx.GoldLabel.SetText(fmt.Sprintf("GOLD: %d G", e.NewAmount))
		}
	})
}

func (ctx *GameContext) SoundPlaySE(name string) {
	if data, ok := ctx.SEBytes[name]; ok && len(data) > 0 {
		ctx.SoundMgr.Play(data, sound.Option{Type: sound.TypeSE, Volume: 0.8})
	}
}
