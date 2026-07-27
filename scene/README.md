# ebitenscene

ebitenscene は Ebitengine (ebiten) 向けの、実践的で堅牢なシーンマネージャーライブラリです。
無駄なインスタンス生成（New/Destroy）によるメモリ負荷やバグを防ぐ `Hide/Show` 機能や、シーン間で安全にデータを引き渡す仕組みを備えており、シューティングゲームやRPGなどの本格的なゲーム開発に最適です。

## 特徴

- **型安全なシーン間データ伝送**: 次のシーンを立ち上げる際、インスタンスを直接渡せるため、スコアデータなどの引き渡しがコンパイル時に検証可能で安全に行えます。
- **メモリに優しい `Hide / Show`**: シーンを破棄せず、状態（スコアやメニューの選択位置など）を保持したまま描画と更新だけを停止・再開できます。ポーズ画面の連打などにも非常に強固です。
- **直感的なメソッド群**: `Start`, `Overlay`, `Hide`, `Show`, `Stop`, `Run`, `Destroy` の組み合わせだけで、あらゆる画面状態を完全に制御します。
- **遅延評価による安全な遷移**: シーン内のUpdateループ中に遷移関数を呼んでも即座にスタックが壊れることはなく、次のフレームで安全に切り替わります。

## アップデート情報

Version 1.0.0 2026/07/27
・初期リリース
・Managerベースでの型安全なシーン管理アーキテクチャを採用

## インストール

```bash
go get github.com/sofia-gros/ebiten/scene
```

## 基本的な使い方

```go
package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/sofia-gros/ebiten/scene"
)

// --- 1. タイトルシーン ---
type TitleScene struct{}

func (s *TitleScene) Update(ctx *scene.Context) error {
	if ebiten.IsKeyPressed(ebiten.KeyEnter) {
		// ゲーム本編へ遷移（TitleScene は破棄される）
		ctx.Start(&GameScene{})
	}
	return nil
}

func (s *TitleScene) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "TITLE SCENE\nPress ENTER to start")
}

// --- 2. ゲーム本編シーン ---
type GameScene struct {
	Score int
	pause *PauseScene // ポインタを保持して使い回す
}

func (s *GameScene) Update(ctx *scene.Context) error {
	s.Score += 1 // スコアが加算されていく

	// Pキーでポーズ画面の切り替え
	if ebiten.IsKeyPressed(ebiten.KeyP) {
		if s.pause == nil {
			// 最初だけ作成して Overlay (上に重ねる)
			s.pause = &PauseScene{}
			ctx.Overlay(s.pause)
		} else {
			// 2回目以降は非表示状態を解除して最前面へ
			ctx.Show(s.pause)
		}
		ctx.Stop(s) // ゲーム本編の Update を一時停止
	}

	// Xキーでゲームオーバー、リザルト画面へ
	if ebiten.IsKeyPressed(ebiten.KeyX) {
		// データを直接渡して起動できるため、型安全！
		result := &ResultScene{FinalScore: s.Score}
		ctx.Start(result)
	}
	return nil
}

func (s *GameScene) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "GAME SCENE\nPress P to Pause, X to Game Over")
}

// --- 3. ポーズシーン ---
type PauseScene struct{}

func (s *PauseScene) Update(ctx *scene.Context) error {
	// スペースキーでゲーム本編へ復帰
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		ctx.Hide(s) // ポーズ画面を隠す（インスタンスは保持）
		
		// 型ベースで GameScene を探し出して Update を再開させる
		ctx.Run(&GameScene{})
	}
	return nil
}

func (s *PauseScene) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "\n\nPAUSE SCENE\nPress SPACE to resume")
}

// --- 4. リザルトシーン ---
type ResultScene struct {
	FinalScore int // GameSceneから流れてきたデータ
}

func (s *ResultScene) Update(ctx *scene.Context) error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		ctx.Start(&TitleScene{}) // タイトルへ戻る
	}
	return nil
}

func (s *ResultScene) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "RESULT SCENE\nPress ESC to return Title")
}

// --- 5. メインループ（Managerのセットアップ） ---
func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("ebitenscene Example")

	// マネージャーの作成
	manager := scene.NewManager(640, 480)
	
	// アプリケーション起動時の最初のシーンをセット
	manager.Context().Start(&TitleScene{})

	// ebiten.RunGame にマネージャーを直接渡す
	if err := ebiten.RunGame(manager); err != nil {
		log.Fatal(err)
	}
}
```

## ライブラリの構成

- `Manager`: Ebitengineのメインループ（`ebiten.Game`）に統合されるシーンの統括管理者です。シーンのスタックと描画順、遅延コマンドの処理を行います。
- `Context`: 各シーンの `Update` に渡されるオブジェクトです。シーン側からマネージャーへの画面遷移の指示を行います。
- `Scene`: 各画面が実装すべきインターフェース（`Update`, `Draw`）です。

### Context の主なメソッド一覧

`Context` を通じて呼び出したシーン操作は即座に反映されるのではなく、**そのフレームの終了時（正確には次回のUpdate開始時）** に安全に適用されます。

| メソッド | 概要 | Updateの状態 | Drawの状態 | 主な用途 |
| :--- | :--- | :---: | :---: | :--- |
| `Start(s)` | 既存のすべてのシーンを破棄し、`s` を新しいルートとして起動します | `true` | `true` | タイトルからゲーム本編へ移行するなど、完全な画面の切り替え |
| `Overlay(s)` | 現在のシーンを維持したまま、その上に `s` を重ねて起動します | `true` | `true` | ポーズ画面やダイアログなどを一時的に上に表示する |
| `Hide(s)` | 指定したシーン `s` を「非表示状態」にします（データは保持） | `false` | `false` | Overlayしたポーズ画面を一時的に隠して元のゲームに戻る |
| `Show(s)` | 非表示だったシーン `s` を再開させ、最前面へ移動します | `true` | `true` | Hideで隠していたポーズ画面を再度呼び出す |
| `Stop(s)` | 指定したシーン `s` の「動き（Update）」だけを止めます | `false` | 変化なし | ポーズ画面を出している間、背景のゲームの動きを一時停止する |
| `Run(s)` | Stopで止めていたシーン `s` の「動き（Update）」を再開させます | `true` | 変化なし | ポーズを解除した際、背景のゲームの動きを再開する |
| `Destroy(s)` | 指定したシーン `s` をスタックから完全に削除（破棄）します | - | - | 不要になった特定のレイヤー（シーン）だけを消去する |

## ライセンス

MIT License
