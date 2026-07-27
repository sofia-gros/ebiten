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

	// プレイヤー（落下する四角形）
	g.playerBody = g.physManager.CreateBody(physics.BodyOptions{
		Type: physics.BodyTypeDynamic,
		X:    160,
		Y:    50,
		Density: 1.0, // Box2DはDensity(密度)がないと動的ボディが回転・落下等で不具合を起こすことがある
		Friction: 0.3,
		Restitution: 0.5, // 0.5の反発力で少しバウンドする
		Shapes: []physics.ShapeDef{
			{Shape: physics.BoxShape{Width: 32, Height: 32}},
		},
		OnCollisionBegin: func(other physics.Body) {
			if other.Group() == GroupFloor {
				fmt.Println("[Box2D] 床にバウンドしました！")
			}
		},
	})
	g.playerBody.SetGroup(GroupPlayer)
	
	// Box2DはNewWorld時に重力0の環境を作っているため、手動で下向きの力をかけるか、速度を与える
	// 今回は下向きの速度を与える
	g.playerBody.SetVelocity(0, 100)

	// 床（静的オブジェクト）
	g.floorBody = g.physManager.CreateBody(physics.BodyOptions{
		Type: physics.BodyTypeStatic,
		X:    160,
		Y:    200,
		Friction: 0.5,
		Shapes: []physics.ShapeDef{
			{Shape: physics.BoxShape{Width: 200, Height: 20}},
		},
	})
	g.floorBody.SetGroup(GroupFloor)

	return g
}

func (g *Game) Update() error {
	// Box2Dのシミュレーションを進める
	g.physManager.Update(1.0 / 60.0)
	
	// Box2Dは摩擦などで止まることがあるので、必要に応じて力を加える
	// 今回は画面外に落ちた時のリセット処理のみ
	x, y := g.playerBody.Position()
	if y > 240 {
		g.playerBody.SetPosition(x, 0)
		g.playerBody.SetVelocity(0, 100)
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// デバッグ用の枠線描画
	g.physManager.DrawDebug(screen)
	
	ebitenutil.DebugPrint(screen, "Box2D Physics Example")
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 320, 240
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Ebiten Physics - Box2D Example")
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
