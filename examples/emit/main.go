package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/sofia-gros/ebiten/emit"
)

type DamageEvent struct {
	PlayerID int
	Damage   int
}

type ItemEvent struct {
	ItemName string
}

type Game struct {
	emitter *emit.Emitter
	log     []string
}

func (g *Game) Init() {
	g.emitter = emit.New()

	// 同期リスナー (On)
	emit.On(g.emitter, func(ev DamageEvent) {
		msg := fmt.Sprintf("[Emit Direct] Player %d took %d damage!", ev.PlayerID, ev.Damage)
		g.log = append(g.log, msg)
	})

	// 1回限定リスナー (Once)
	emit.Once(g.emitter, func(ev ItemEvent) {
		msg := fmt.Sprintf("[Once] First item collected: %s", ev.ItemName)
		g.log = append(g.log, msg)
	})

	// 継続リスナー (On)
	emit.On(g.emitter, func(ev ItemEvent) {
		msg := fmt.Sprintf("[Queued -> Flushed] Item: %s", ev.ItemName)
		g.log = append(g.log, msg)
	})
}

func (g *Game) Update() error {
	// D キー: 即時同期発行 (Emit)
	if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		emit.Emit(g.emitter, DamageEvent{PlayerID: 1, Damage: 15})
	}

	// I キー: 遅延キュー発行 (Queue)
	if inpututil.IsKeyJustPressed(ebiten.KeyI) {
		emit.Queue(g.emitter, ItemEvent{ItemName: "Health Potion"})
	}

	// F キー: キューの一括フラッシュ (Flush)
	if inpututil.IsKeyJustPressed(ebiten.KeyF) {
		g.emitter.Flush()
	}


	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "[EMIT DEMO]\nPress D: Emit DamageEvent Immediately\nPress I: Queue ItemEvent\nPress F: Flush Queued Events\n\nLogs:")
	y := 100
	for i, msg := range g.log {
		if i >= 12 {
			break
		}
		ebitenutil.DebugPrintAt(screen, msg, 10, y)
		y += 20
	}
}

func (g *Game) Layout(w, h int) (int, int) {
	return 640, 480
}

func main() {
	g := &Game{}
	g.Init()

	ebiten.SetWindowTitle("emit Demo - Emit, Queue, Flush, Once")
	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}
