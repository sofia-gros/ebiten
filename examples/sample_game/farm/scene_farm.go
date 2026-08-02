package main

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"os"


	"github.com/hajimehoshi/ebiten/v2"

	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/sofia-gros/ebiten/camera"
	"github.com/sofia-gros/ebiten/emit"
	"github.com/sofia-gros/ebiten/physics"
	"github.com/sofia-gros/ebiten/scene"
	"github.com/sofia-gros/ebiten/sound"
	"github.com/sofia-gros/ebiten/ui"
)

type FarmScene struct {
	ctx          *GameContext
	playerBody   physics.Body
	waterTrack   *sound.Track
	grassImg     *ebiten.Image
	houseImg     *ebiten.Image
	tilledImg    *ebiten.Image
	plantsImg    *ebiten.Image
	charSheetImg *ebiten.Image
}

func (s *FarmScene) Init(sCtx *scene.Context) {
	// 1. テクスチャのロード (asset)
	s.loadTextures()

	// 2. プレイヤー物理ボディ作成 (physics)
	s.playerBody = s.ctx.PhysicsMgr.CreateBody(physics.BodyOptions{
		Type:  physics.BodyTypeDynamic,
		X:     s.ctx.PlayerX,
		Y:     s.ctx.PlayerY,
		Shape: physics.BoxShape{Width: 16, Height: 16},
	})

	// 障害物の静的物理ボディ作成 (マップ外周、家、柵)
	s.setupMapPhysics()


	// 3. UI パネルとツール選択ボタンの構築 (ui)
	s.setupUI()

	// 4. 音声・BGM の開始 (sound)
	s.setupAudio()

	// 5. カメラ初期位置
	s.ctx.MainCam.SetPos(s.ctx.PlayerX, s.ctx.PlayerY)
	s.ctx.MiniCam.SetPos(s.ctx.PlayerX, s.ctx.PlayerY)
}

func (s *FarmScene) loadTextures() {
	prefix := "assets/Sprout Lands - Sprites - Basic pack/"

	// 草原・タイルセット
	if img, _, err := ebitenutil.NewImageFromFile(prefix + "Tilesets/Grass.png"); err == nil {
		s.grassImg = img
	}
	if img, _, err := ebitenutil.NewImageFromFile(prefix + "Tilesets/Wooden House.png"); err == nil {
		s.houseImg = img
	}
	if img, _, err := ebitenutil.NewImageFromFile(prefix + "Tilesets/Tilled Dirt.png"); err == nil {
		s.tilledImg = img
	}
	if img, _, err := ebitenutil.NewImageFromFile(prefix + "Objects/Basic Plants.png"); err == nil {
		s.plantsImg = img
	}
	if img, _, err := ebitenutil.NewImageFromFile(prefix + "Characters/Basic Charakter Spritesheet.png"); err == nil {
		s.charSheetImg = img
	}

	// SE 音声ファイルのロード
	seFiles := map[string]string{
		"click":      "assets/Audio/click.wav",
		"coinsCling": "assets/Audio/coinsCling.wav",
		"chestOpen":  "assets/Audio/chestOpen.wav",
		"footstep":   "assets/Audio/footstepDirt1.wav",
	}

	for key, path := range seFiles {
		if b, err := os.ReadFile(path); err == nil {
			s.ctx.SEBytes[key] = b
		}
	}
}

func (s *FarmScene) setupMapPhysics() {
	// マップ外周境界
	s.ctx.PhysicsMgr.CreateBody(physics.BodyOptions{
		Type: physics.BodyTypeStatic, X: 640, Y: 0,
		Shape: physics.BoxShape{Width: 1280, Height: 32},
	})
	s.ctx.PhysicsMgr.CreateBody(physics.BodyOptions{
		Type: physics.BodyTypeStatic, X: 640, Y: 960,
		Shape: physics.BoxShape{Width: 1280, Height: 32},
	})
	s.ctx.PhysicsMgr.CreateBody(physics.BodyOptions{
		Type: physics.BodyTypeStatic, X: 0, Y: 480,
		Shape: physics.BoxShape{Width: 32, Height: 960},
	})
	s.ctx.PhysicsMgr.CreateBody(physics.BodyOptions{
		Type: physics.BodyTypeStatic, X: 1280, Y: 480,
		Shape: physics.BoxShape{Width: 32, Height: 960},
	})

	// 家の物理障害判定 (320, 240)
	s.ctx.PhysicsMgr.CreateBody(physics.BodyOptions{
		Type: physics.BodyTypeStatic, X: 320, Y: 240,
		Shape: physics.BoxShape{Width: 128, Height: 96},
	})
}


func (s *FarmScene) setupUI() {
	s.ctx.UIRoot.Clear()

	// 所持金 Label (ui)
	s.ctx.GoldLabel = ui.NewLabel(fmt.Sprintf("GOLD: %d G", s.ctx.Gold))
	s.ctx.GoldLabel.SetPos(20, 20)
	s.ctx.UIRoot.Add(s.ctx.GoldLabel)

	// ステータス Label
	s.ctx.StatusLabel = ui.NewLabel("Tool: HOE (1:Hoe 2:Water 3:Seed 4:Harvest)")
	s.ctx.StatusLabel.SetPos(20, 45)
	s.ctx.UIRoot.Add(s.ctx.StatusLabel)


	// ツール切替ボタン群 (VBox)
	vbox := ui.NewVBox()
	vbox.SetPos(20, 80)
	vbox.SetSpacing(10)

	btnHoe := ui.NewButton(ui.ButtonOption{Text: "1: Hoe", Width: 90, Height: 30})
	btnHoe.OnClick(func() { s.selectTool(ToolHoe) })

	btnWater := ui.NewButton(ui.ButtonOption{Text: "2: Water", Width: 90, Height: 30})
	btnWater.OnClick(func() { s.selectTool(ToolWater) })

	btnSeed := ui.NewButton(ui.ButtonOption{Text: "3: Seed", Width: 90, Height: 30})
	btnSeed.OnClick(func() { s.selectTool(ToolSeed) })

	btnHarvest := ui.NewButton(ui.ButtonOption{Text: "4: Harvest", Width: 90, Height: 30})
	btnHarvest.OnClick(func() { s.selectTool(ToolHarvest) })

	vbox.Add(btnHoe)
	vbox.Add(btnWater)
	vbox.Add(btnSeed)
	vbox.Add(btnHarvest)

	s.ctx.UIRoot.Add(vbox)
}


func (s *FarmScene) selectTool(t ToolType) {
	s.ctx.SelectedTool = t
	names := []string{"HOE (クワ)", "WATER (ジョウロ)", "SEED (タネ)", "HARVEST (カマ)"}
	if s.ctx.StatusLabel != nil {
		s.ctx.StatusLabel.SetText("Tool: " + names[t])
	}
	s.ctx.SoundPlaySE("click")
}

func (s *FarmScene) setupAudio() {
	// BGM ループ再生 (sound)
	if bgmBytes, err := os.ReadFile("assets/blossom.wav"); err == nil {
		s.ctx.BGMTrack = s.ctx.SoundMgr.Play(bgmBytes, sound.Option{Type: sound.TypeBGM, Loop: true, Volume: 0.5})
	}

	// 2D ポジショナルサウンド (川の環境音 at 800, 480)
	if waterSE, ok := s.ctx.SEBytes["chestOpen"]; ok {
		s.waterTrack = s.ctx.SoundMgr.PlayAt(waterSE, 800, 480, s.ctx.PlayerX, s.ctx.PlayerY, 500.0, sound.Option{Type: sound.TypeEnv, Loop: true, Volume: 0.4})
	}
}

func (s *FarmScene) Update(sCtx *scene.Context) error {
	dt := 1.0 / 60.0

	// 1. 移動入力 (pad & WASD 相対移動)
	fwdX, fwdY := s.ctx.MainCam.Forward()
	rtX, rtY := s.ctx.MainCam.Right()

	dx, dy := 0.0, 0.0
	if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyUp) {
		dx += fwdX
		dy += fwdY
		s.ctx.PlayerDir = 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyDown) {
		dx -= fwdX
		dy -= fwdY
		s.ctx.PlayerDir = 0
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyLeft) {
		dx -= rtX
		dy -= rtY
		s.ctx.PlayerDir = 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyRight) {
		dx += rtX
		dy += rtY
		s.ctx.PlayerDir = 3
	}

	// 物理エンジンの速度更新 (physics)
	if dx != 0 || dy != 0 {
		norm := math.Sqrt(dx*dx + dy*dy)
		vx := (dx / norm) * s.ctx.PlayerSpeed * 60.0
		vy := (dy / norm) * s.ctx.PlayerSpeed * 60.0
		if s.playerBody != nil {
			s.playerBody.SetVelocity(vx, vy)
		}
		s.ctx.AnimTimer += dt * 8
	} else {
		if s.playerBody != nil {
			s.playerBody.SetVelocity(0, 0)
		}
	}

	// 物理シミュレーションステップ進行
	s.ctx.PhysicsMgr.Update(dt)

	// プレイヤー位置更新
	if s.playerBody != nil {
		px, py := s.playerBody.Position()
		s.ctx.PlayerX = px
		s.ctx.PlayerY = py
	}


	// 2D 音声ポジション追従 (sound)
	if s.waterTrack != nil {
		s.waterTrack.SetPosition(800, 480, s.ctx.PlayerX, s.ctx.PlayerY, 500.0)
	}

	// 2. アクション入力 (Space / Z / Enter) -> 農地アクション
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyZ) {
		s.handleFarmAction()
	}

	// ツール数字キー切替 (1, 2, 3, 4)
	if inpututil.IsKeyJustPressed(ebiten.Key1) { s.selectTool(ToolHoe) }
	if inpututil.IsKeyJustPressed(ebiten.Key2) { s.selectTool(ToolWater) }
	if inpututil.IsKeyJustPressed(ebiten.Key3) { s.selectTool(ToolSeed) }
	if inpututil.IsKeyJustPressed(ebiten.Key4) { s.selectTool(ToolHarvest) }

	// K キー: セーブ, L キー: ロード (save)
	if inpututil.IsKeyJustPressed(ebiten.KeyK) {
		sd := FarmSaveData{
			PlayerX:   s.ctx.PlayerX,
			PlayerY:   s.ctx.PlayerY,
			Gold:      s.ctx.Gold,
			TilledMap: s.ctx.TilledTiles,
		}
		if err := s.ctx.SaveMgr.Slot(1).Save("farm", "json", sd); err == nil {
			s.ctx.SoundPlaySE("chestOpen")
			s.ctx.MainCam.Shake(10, 0.4)
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyL) {
		var sd FarmSaveData
		if err := s.ctx.SaveMgr.Slot(1).Load("farm", "json", &sd); err == nil {
			s.ctx.PlayerX = sd.PlayerX
			s.ctx.PlayerY = sd.PlayerY
			s.ctx.Gold = sd.Gold
			if sd.TilledMap != nil {
				s.ctx.TilledTiles = sd.TilledMap
			}
			if s.playerBody != nil {
				s.playerBody.SetPosition(sd.PlayerX, sd.PlayerY)
			}
			emit.Emit(s.ctx.Emitter, GoldChangedEvent{NewAmount: s.ctx.Gold})
			s.ctx.SoundPlaySE("coinsCling")
		}
	}



	// 3. カメラ追従 (camera)
	s.ctx.MainCam.MoveTo(s.ctx.PlayerX, s.ctx.PlayerY, 0.1)
	s.ctx.MiniCam.SetPos(s.ctx.PlayerX, s.ctx.PlayerY)
	s.ctx.CameraGroup.Update(dt)

	// 4. UI と Tween の更新 (ui & tween)
	s.ctx.TweenGroup.Update(dt)
	s.ctx.UIRoot.Update()

	return nil
}

func (s *FarmScene) handleFarmAction() {
	// プレイヤーの前方 24px のグリッド位置を計算
	tx := int((s.ctx.PlayerX) / 32)
	ty := int((s.ctx.PlayerY) / 32)
	key := fmt.Sprintf("%d,%d", tx, ty)

	currentStage, exists := s.ctx.TilledTiles[key]

	switch s.ctx.SelectedTool {
	case ToolHoe: // 耕す
		if !exists {
			s.ctx.TilledTiles[key] = 0
			s.ctx.SoundPlaySE("footstep")
			s.ctx.MainCam.Shake(5, 0.2)
		}

	case ToolWater: // 水やり
		if exists && currentStage == 0 {
			s.ctx.TilledTiles[key] = 1
			s.ctx.SoundPlaySE("click")
		}

	case ToolSeed: // 種まき (費用 10G)
		if exists && currentStage == 1 && s.ctx.Gold >= 10 {
			s.ctx.Gold -= 10
			s.ctx.TilledTiles[key] = 2
			emit.Emit(s.ctx.Emitter, GoldChangedEvent{NewAmount: s.ctx.Gold})
			emit.Emit(s.ctx.Emitter, CropPlantedEvent{TileX: tx, TileY: ty, CropType: "Tomato"})
		}

	case ToolHarvest: // 収穫 (報酬 +25G)
		if exists && currentStage == 2 {
			delete(s.ctx.TilledTiles, key)
			emit.Emit(s.ctx.Emitter, CropHarvestedEvent{TileX: tx, TileY: ty, GoldEarned: 25})
		}
	}
}

func (s *FarmScene) Draw(screen *ebiten.Image) {
	// カメラグループ一括描画 (ZIndex 順: メイン -> ミニマップ)
	s.ctx.CameraGroup.Render(screen, func(cam *camera.Camera, target *ebiten.Image) {
		// 1. 背景・川の描画
		target.Fill(color.RGBA{70, 150, 180, 255})

		// 2. 草原島 (600x600)
		island := ebiten.NewImage(800, 800)
		island.Fill(color.RGBA{135, 195, 95, 255})
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(240, 140)
		target.DrawImage(island, opts)

		// 3. 家の描画 (320, 240)
		if s.houseImg != nil {
			hOpts := &ebiten.DrawImageOptions{}
			hOpts.GeoM.Translate(250, 180)
			target.DrawImage(s.houseImg, hOpts)
		} else {
			house := ebiten.NewImage(128, 96)
			house.Fill(color.RGBA{160, 100, 60, 255})
			hOpts := &ebiten.DrawImageOptions{}
			hOpts.GeoM.Translate(250, 180)
			target.DrawImage(house, hOpts)
		}

		// 4. 農地タイル & 作物の描画
		for tileKey, stage := range s.ctx.TilledTiles {
			var tx, ty int
			fmt.Sscanf(tileKey, "%d,%d", &tx, &ty)
			wx := float64(tx * 32)
			wy := float64(ty * 32)

			// 耕した土
			if s.tilledImg != nil {
				subTilled := s.tilledImg.SubImage(image.Rect(0, 0, 16, 16)).(*ebiten.Image)
				dOpts := &ebiten.DrawImageOptions{}
				dOpts.GeoM.Scale(2.0, 2.0)
				dOpts.GeoM.Translate(wx, wy)
				target.DrawImage(subTilled, dOpts)
			} else {
				dirt := ebiten.NewImage(32, 32)
				dirt.Fill(color.RGBA{140, 90, 50, 255})
				dOpts := &ebiten.DrawImageOptions{}
				dOpts.GeoM.Translate(wx, wy)
				target.DrawImage(dirt, dOpts)
			}

			// 作物 (種/芽/トマト)
			if stage >= 2 && s.plantsImg != nil {
				// Basic Plants.png からトマトの各成長段階を抽出
				cropSub := s.plantsImg.SubImage(image.Rect(16, 0, 32, 16)).(*ebiten.Image)
				cOpts := &ebiten.DrawImageOptions{}
				cOpts.GeoM.Scale(2.0, 2.0)
				cOpts.GeoM.Translate(wx, wy)
				target.DrawImage(cropSub, cOpts)
			}
		}

		// 5. プレイヤーキャラクター描画 (Sprout Lands Character スプライトシート)
		if s.charSheetImg != nil {
			// キャラクターコマ割り (方向別)
			frameX := (int(s.ctx.AnimTimer) % 2) * 48
			frameY := s.ctx.PlayerDir * 48
			subChar := s.charSheetImg.SubImage(image.Rect(frameX, frameY, frameX+48, frameY+48)).(*ebiten.Image)
			pOpts := &ebiten.DrawImageOptions{}
			pOpts.GeoM.Translate(s.ctx.PlayerX-24, s.ctx.PlayerY-24)
			target.DrawImage(subChar, pOpts)
		} else {
			pImg := ebiten.NewImage(24, 24)
			pImg.Fill(color.RGBA{240, 240, 240, 255})
			pOpts := &ebiten.DrawImageOptions{}
			pOpts.GeoM.Translate(s.ctx.PlayerX-12, s.ctx.PlayerY-12)
			target.DrawImage(pImg, pOpts)
		}

	})

	// 6. UI コンテナのオーバーレイ描画
	s.ctx.UIRoot.Draw(screen)

	// ガイドテキスト
	ebitenutil.DebugPrintAt(screen, "[Controls] WASD: Move | SPACE: Tool Action | K: Save | L: Load", 140, 455)
}

func (s *FarmScene) Destroy(sCtx *scene.Context) {
	if s.waterTrack != nil {
		s.waterTrack.Stop()
	}
	s.ctx.UIRoot.Clear()
}
