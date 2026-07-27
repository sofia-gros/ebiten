# ebintenphysics

ebintenphysics は Ebitengine (ebiten) 向けの、実践的で堅牢な物理エンジンマネージャーライブラリです。
エンジン間の差異をなくし、同一操作により様々なエンジンに対応したAPIを提供しており、シューティングゲームやRPG、FPSなどの本格的なゲーム開発に最適です。

## 特徴

- **システムに優しい**: 外部の物理エンジンライブラリを呼び出すので、成熟した物理エンジンをゲームに組み込めます。

## インストール

```bash
go get [github.com/sofia-gros/ebiten/physics](https://github.com/sofia-gros/ebiten/physics)
```

## 基本的な使い方

```go
package main

import (
	"fmt"
	"log"

	"[github.com/hajimehoshi/ebiten/v2](https://github.com/hajimehoshi/ebiten/v2)"
	"[github.com/sofia-gros/ebiten/physics](https://github.com/sofia-gros/ebiten/physics)"
)

// --- 1. 物理エンジンライブラリ ---
type Physics struct {
	
}



// --- 1. ゲーム本編シーン ---
type GameScene struct {
	Score int
	pause *PauseScene // ポインタを保持して使い回す
}

func (s *GameScene) Update(ctx *scene.Context) error {
	s.Score += 10 // スコアが加算されていく

	// Pキーでポーズ画面の切り替え（トグル操作）
	if ebiten.IsKeyPressed(ebiten.KeyP) {
		if s.pause == nil {
			// 最初だけ作成して Overlay
			s.pause = &PauseScene{}
			ctx.Overlay(s.pause)
		} else {
			// 2回目以降は非表示状態を解除して最前面へ（Show）
			ctx.Show(s.pause)
			ctx.Stop(s) // ゲーム本編の動きを一時停止
		}
	}

	// ゲームオーバー時にスコアを持ってリザルト画面へ
	if ebiten.IsKeyPressed(ebiten.KeyEnter) {
		result := &ResultScene{FinalScore: s.Score} // データを直接渡して起動
		ctx.Start(result)                           // GameSceneは破棄される
	}
	return nil
}

func (s *GameScene) Draw(screen *ebiten.Image) {}

// --- 2. ポーズシーン ---
type PauseScene struct {
	CursorIdx int // 選択中のメニュー位置（Hideしてもこの状態が保持される）
}

func (s *PauseScene) Update(ctx *scene.Context) error {
	// メニュー選択ロジックなどをここに記述...

	// 再度Pキーが押されたら、自分を非表示にしてゲーム本編を復帰
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		ctx.Hide(s)                       // ポーズ画面を隠す（データは保持）
		ctx.Run(&GameScene{})            // ゲーム本編を再開
	}
	return nil
}

func (s *PauseScene) Draw(screen *ebiten.Image) {}

// --- 3. リザルトシーン（データを受け取る画面） ---
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
	// "FINAL SCORE: X" を描画
}

// --- 省略: タイトル、メインループ（Managerのセットアップなど） ---
```
