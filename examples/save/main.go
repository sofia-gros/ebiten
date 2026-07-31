package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/sofia-gros/ebiten/save"
)

type PlayerData struct {
	Name  string `json:"name"`
	Level int    `json:"level"`
	Gold  int    `json:"gold"`
}

type Game struct {
	m          *save.Manager
	loadedData PlayerData
	statusMsg  string
}

func (g *Game) Init() {
	g.m = save.NewManager("save_demo", save.Option{
		Dir: save.DirCurrent,
	})
	g.statusMsg = "Press 1: Save Slot 1 | Press 2: Save New Slot (Auto Fill) | Press L: Load Slot 1 | Press C: Clear All"
}

func (g *Game) Update() error {
	data := PlayerData{Name: "Hero", Level: 10, Gold: 500}

	// 1 キー: Slot 1 保存
	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		err := g.m.Slot(1).Save("player", "json", data)
		if err == nil {
			g.statusMsg = "Saved into Slot 1 (player_1.json)"
		}
	}

	// 2 キー: SaveNewSlot (欠番自動検出保存)
	if inpututil.IsKeyJustPressed(ebiten.Key2) {
		slotNum, err := g.m.SaveNewSlot("player", "json", data)
		if err == nil {
			g.statusMsg = fmt.Sprintf("Auto-filled & Saved into Slot %d (player_%d.json)", slotNum, slotNum)
		}
	}

	// L キー: Slot 1 ロード
	if inpututil.IsKeyJustPressed(ebiten.KeyL) {
		var loaded PlayerData
		err := g.m.Slot(1).Load("player", "json", &loaded)
		if err == nil {
			g.loadedData = loaded
			g.statusMsg = fmt.Sprintf("Loaded Slot 1: Name=%s, Level=%d, Gold=%d", loaded.Name, loaded.Level, loaded.Gold)
		} else {
			g.statusMsg = "Slot 1 load failed (file does not exist)"
		}
	}

	// C キー: 暗号化保存デモ (secret.dat)
	if inpututil.IsKeyJustPressed(ebiten.KeyC) {
		err := g.m.Save("secret", "dat", data, save.Option{
			Format:     save.FormatBinary,
			EncryptKey: "my_secret_key",
			Checksum:   true,
		})
		if err == nil {
			g.statusMsg = "Encrypted binary saved to secret.dat with AES & Checksum!"
		}
	}

	return nil
}


func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, fmt.Sprintf("[SAVE DEMO]\n\nControls:\n[1] Save Slot 1\n[2] Save New Slot (Auto-fill)\n[L] Load Slot 1\n[C] Save Encrypted Binary (AES-GCM)\n\nStatus: %s", g.statusMsg))
}

func (g *Game) Layout(w, h int) (int, int) {
	return 640, 480
}

func main() {
	g := &Game{}
	g.Init()

	ebiten.SetWindowTitle("save Demo - Slot, SaveNewSlot, Encryption")
	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}
