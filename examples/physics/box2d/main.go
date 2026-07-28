package main

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"github.com/sofia-gros/ebiten/physics"
	"github.com/sofia-gros/ebiten/physics/adapters/box2d"
)

const (
	GroupPlayer physics.Group = 1
	GroupFloor  physics.Group = 2
)

type Game struct {
	physManager *physics.Manager
	playerBody  physics.Body
	floorBody   physics.Body
}

func NewGame() *Game {
	g := &Game{
		physManager: physics.NewManager(),
	}
	
	// Box2Dエンジンをセット
	g.physManager.SetWorld(box2d.NewWorld())
	// 共通で重力をセット
	g.physManager.SetGravity(0, 100)

	// プレイヤー（落下する四角形）
	g.playerBody = g.physManager.CreateBody(physics.BodyOptions{
		Type: physics.BodyTypeDynamic,
		X:    160,
		Y:    50,
		Shapes: []physics.ShapeDef{
			{Shape: physics.BoxShape{Width: 32, Height: 32}},
		},
		Restitution: 0.8, // 跳ね返り係数
		OnCollisionBegin: func(other physics.Body) {
			if other.Group() == GroupFloor {
				fmt.Printf("[Box2D] 床に着地しました！\n")
			}
		},
	})
	g.playerBody.SetGroup(GroupPlayer)

	// 床
	g.floorBody = g.physManager.CreateBody(physics.BodyOptions{
		Type: physics.BodyTypeStatic,
		X:    160,
		Y:    200,
		Shapes: []physics.ShapeDef{
			{Shape: physics.BoxShape{Width: 200, Height: 20}},
		},
	})
	g.floorBody.SetGroup(GroupFloor)

	return g
}

func (g *Game) Update() error {
	// 物理エンジンの更新（60FPS想定）
	g.physManager.World().Step(1.0 / 60.0)

	// 床の下まで落ちたら上に戻す（ループ処理）
	x, y := g.playerBody.Position()
	if y > 240 {
		g.playerBody.SetPosition(x, 0)
		g.playerBody.SetVelocity(0, 0)
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.physManager.DrawDebug(screen)
	ebitenutil.DebugPrint(screen, "Ebiten Physics: Box2D Engine")
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Ebiten Physics - Box2D Example")
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
