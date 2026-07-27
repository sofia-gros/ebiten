package main

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"github.com/sofia-gros/ebiten/physics"
	"github.com/sofia-gros/ebiten/physics/adapters/arcade"
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
	
	// Arcadeエンジンをセット
	g.physManager.SetWorld(arcade.NewWorld())

	// プレイヤー（落下する四角形）
	g.playerBody = g.physManager.CreateBody(physics.BodyOptions{
		Type: physics.BodyTypeDynamic,
		X:    160,
		Y:    50,
		Shapes: []physics.ShapeDef{
			{Shape: physics.BoxShape{Width: 32, Height: 32}},
		},
		OnCollisionBegin: func(other physics.Body) {
			if other.Group() == GroupFloor {
				fmt.Println("[Arcade] 床に着地しました！")
			}
		},
	})
	g.playerBody.SetGroup(GroupPlayer)
	
	// 下向きの速度（重力の代わり。Arcadeは自力で速度を制御する）
	g.playerBody.SetVelocity(0, 100)

	// 床（静的オブジェクト）
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
	// Arcadeでは動的ボディは設定されたVelocityに従って進む
	g.physManager.Update(1.0 / 60.0)
	
	// 画面外に落ちたら上に戻す
	x, y := g.playerBody.Position()
	if y > 240 {
		g.playerBody.SetPosition(x, 0)
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// デバッグ用の枠線描画
	g.physManager.DrawDebug(screen)
	
	ebitenutil.DebugPrint(screen, "Arcade Physics Example")
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 320, 240
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Ebiten Physics - Arcade Example")
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
