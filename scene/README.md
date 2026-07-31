# ebitenscene

[English](./README_en.md)

ebitenscene は Ebitengine 向けのシーン管理・画面遷移ライブラリです。
ゲームの各画面（タイトル、ステージ、メニューなど）を独立したコンポーネントとして分離し、シーンの上に一時停止画面やダイアログを重ねる「オーバーレイスタック」機能を提供します。

---

## 特徴

- **独立したシーンインターフェース (`Scene`)**:
  - 各画面を `Init`, `Update`, `Draw`, `Destroy` を持つ独立したオブジェクトとして管理。
- **直感的なスタックベースの制御 (`Switch`, `Push`, `Pop`)**:
  - `Switch`: 完全に新しい画面へ全入れ替え遷移。
  - `Push`: 現在のゲーム画面の上にポーズ画面やダイアログを重ねて表示 (Overlay)。
  - `Pop`: オーバーレイを閉じて下のゲーム画面へ復帰。
- **シンプルな更新・描画委譲**:
  - ゲーム本体の `Update` と `Draw` に `sceneManager.Update()` と `sceneManager.Draw(screen)` を書くだけでアクティブなシーン群を管理。

---

## インストール

```bash
go get github.com/sofia-gros/ebiten/scene
```

---

## 使い方

### クイックスタート

タイトル画面からゲームメイン画面へ一方向で切り替える基本的なコード例です。


```go
package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sofia-gros/ebiten/scene"
)

// --- タイトルシーン ---
type TitleScene struct {
	sm *scene.Manager
}

func (s *TitleScene) Init() {}
func (s *TitleScene) Update() error {
	// スペースキーでゲームメイン画面へ切り替え
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		s.sm.Switch(&GameScene{sm: s.sm})
	}
	return nil
}
func (s *TitleScene) Draw(screen *ebiten.Image) {}
func (s *TitleScene) Destroy() {}

// --- ゲームメインシーン ---
type GameScene struct {
	sm *scene.Manager
}

func (s *GameScene) Init() {}
func (s *GameScene) Update() error { return nil }
func (s *GameScene) Draw(screen *ebiten.Image) {}
func (s *GameScene) Destroy() {}

// --- エントリーポイント ---
type Game struct {
	sm *scene.Manager
}

func (g *Game) Update() error { return g.sm.Update() }
func (g *Game) Draw(screen *ebiten.Image) { g.sm.Draw(screen) }
func (g *Game) Layout(w, h int) (int, int) { return 640, 480 }

func main() {
	sm := scene.NewManager()
	g := &Game{sm: sm}
	sm.Switch(&TitleScene{sm: sm})

	ebiten.RunGame(g)
}
```

---

### 全機能の使い方

ゲームプレイ中にポーズ画面を上に重ねて表示し、元の画面を破棄せずに一時停止・復帰させる全機能の使い方です。


```go
package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sofia-gros/ebiten/scene"
)

// --- ポーズオーバーレイシーン ---
type PauseScene struct {
	sm *scene.Manager
}

func (s *PauseScene) Init() {}
func (s *PauseScene) Update() error {
	// ESC キーでポーズ画面を閉じて下のゲームシーンへ復帰
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		s.sm.Pop()
	}
	return nil
}
func (s *PauseScene) Draw(screen *ebiten.Image) {
	// 下のゲーム画面が背景に見えた状態で、上に「PAUSE」メニューを描画
}
func (s *PauseScene) Destroy() {}

// --- ゲームメインシーン ---
type GameScene struct {
	sm *scene.Manager
}

func (s *GameScene) Init() {}
func (s *GameScene) Update() error {
	// P キーでポーズ画面を重ね表示 (Push)
	if ebiten.IsKeyPressed(ebiten.KeyP) {
		s.sm.Push(&PauseScene{sm: s.sm})
	}
	return nil
}
func (s *GameScene) Draw(screen *ebiten.Image) {}
func (s *GameScene) Destroy() {}
```

---

## 主要 API リファレンス

### `scene.Manager`
- `NewManager()`: シーンマネージャーを作成。
- `Switch(newScene Scene)`: 現在の全シーンを破棄し、新しいシーンへ一括全入れ替え切り替え。
- `Push(overlayScene Scene)`: 現在の画面の上にポーズメニュー等のオーバーレイ画面を重畳表示。
- `Pop()`: 最前面のオーバーレイ画面を破棄し、下の画面に制御を復帰。
- `Clear()`: 全シーンの完全クリア。
- `Update() error`: 最前面のアクティブシーン（または設定されたスタック）を更新。
- `Draw(screen *ebiten.Image)`: スタックされているシーン群を下から順に一括描画。

---

## ライセンス

MIT License
