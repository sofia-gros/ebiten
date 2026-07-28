package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/sofia-gros/ebiten/pad/input"
	"github.com/sofia-gros/ebiten/physics"
	"github.com/sofia-gros/ebiten/physics/adapters/arcade"
	"github.com/sofia-gros/ebiten/scene"
)

const (
	GroupPlayer physics.Group = 1
	GroupWall   physics.Group = 2
)

type GameScene struct {
	in          *input.Input
	physManager *physics.Manager
	playerBody  physics.Body
	hitMessage  string
	hitTimer    int
}

func NewGameScene(in *input.Input) *GameScene {
	s := &GameScene{
		in:          in,
		physManager: physics.NewManager(),
	}

	// Arcadeエンジンをセットし、重力を0に（トップダウンRPGなので）
	s.physManager.SetWorld(arcade.NewWorld())
	s.physManager.SetGravity(0, 0)

	// プレイヤー（四角形）
	s.playerBody = s.physManager.CreateBody(physics.BodyOptions{
		Type: physics.BodyTypeDynamic,
		X:    320,
		Y:    240,
		Shapes: []physics.ShapeDef{
			{Shape: physics.BoxShape{Width: 32, Height: 32}},
		},
		Restitution: 0.0, // 跳ね返らない
		OnCollisionBegin: func(other physics.Body) {
			if other.Group() == GroupWall {
				s.hitMessage = "Ouch! Wall!"
				s.hitTimer = 60
			}
		},
	})
	s.playerBody.SetGroup(GroupPlayer)

	// 壁（画面の枠）
	walls := []physics.BodyOptions{
		{Type: physics.BodyTypeStatic, X: 320, Y: 10, Shapes: []physics.ShapeDef{{Shape: physics.BoxShape{Width: 600, Height: 20}}}},  // 上
		{Type: physics.BodyTypeStatic, X: 320, Y: 470, Shapes: []physics.ShapeDef{{Shape: physics.BoxShape{Width: 600, Height: 20}}}}, // 下
		{Type: physics.BodyTypeStatic, X: 10, Y: 240, Shapes: []physics.ShapeDef{{Shape: physics.BoxShape{Width: 20, Height: 480}}}},  // 左
		{Type: physics.BodyTypeStatic, X: 630, Y: 240, Shapes: []physics.ShapeDef{{Shape: physics.BoxShape{Width: 20, Height: 480}}}}, // 右
		// 中央の障害物
		{Type: physics.BodyTypeStatic, X: 200, Y: 150, Shapes: []physics.ShapeDef{{Shape: physics.BoxShape{Width: 80, Height: 80}}}},
		{Type: physics.BodyTypeStatic, X: 450, Y: 350, Shapes: []physics.ShapeDef{{Shape: physics.BoxShape{Width: 100, Height: 40}}}},
	}

	for _, w := range walls {
		body := s.physManager.CreateBody(w)
		body.SetGroup(GroupWall)
	}

	return s
}

func (s *GameScene) Update(ctx *scene.Context) error {
	// 移動入力の取得
	vx, vy := 0.0, 0.0
	speed := 200.0 // ピクセル/秒

	if state, ok := s.in.GetActionState(ActionMove); ok {
		// 斜め移動時の速度制限（正規化）
		vx, vy = state.X, state.Y
		len := math.Sqrt(vx*vx + vy*vy)
		if len > 0 {
			vx = (vx / len) * speed * state.Strength
			vy = (vy / len) * speed * state.Strength
		}
	}

	// プレイヤーの速度を更新
	s.playerBody.SetVelocity(vx, vy)

	// アクション入力
	if s.in.JustPressed(ActionHit) {
		s.hitMessage = "Attack!"
		s.hitTimer = 30
	}

	// メッセージタイマー
	if s.hitTimer > 0 {
		s.hitTimer--
	} else {
		s.hitMessage = ""
	}

	// 物理エンジンの更新（60FPS想定）
	s.physManager.World().Step(1.0 / 60.0)

	return nil
}

func (s *GameScene) Draw(screen *ebiten.Image) {
	// 背景塗りつぶし
	screen.Fill(color.RGBA{50, 150, 50, 255})

	// 物理オブジェクト（プレイヤーと壁）のデバッグ描画
	s.physManager.DrawDebug(screen)

	// メッセージ表示
	ebitenutil.DebugPrintAt(screen, "Game Scene - Move with Stick/WASD, Hit with Button/Space", 20, 20)
	if s.hitMessage != "" {
		x, y := s.playerBody.Position()
		ebitenutil.DebugPrintAt(screen, s.hitMessage, int(x)-20, int(y)-40)
	}
}
