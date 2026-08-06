# ebitenanimation

[English](./README_en.md)

ebitenanimation は Ebitengine 向けの直感的かつ高機能な 2D スプライトアニメーション制御ライブラリです。
手動でのコマ定義コード記述に加えて、**Aseprite や TexturePacker などの汎用 JSON フォーマットからの全自動アニメーション構築** に対応しています。

外部依存がなく `ebiten/v2` のみに直接依存しているため、単体で軽量にご利用いただけます。

---

## 主な機能

- **一括・個別制御の二重構造**:
  - `Manager`: ゲーム内の全 `Animator` を一括更新・一括ポーズ・スローモーション/早送り管理。
  - `Animator`: 各オブジェクト（プレイヤー、敵、アイテム等）ごとの個別再生・一時停止・リセット・速度変更。
- **汎用 JSON 全自動読込 (`CreateAnimatorFromJSON`)**:
  - Aseprite や TexturePacker が出力した JSON (`frameTags`, 表示時間 `duration`) からアニメーション群を一発で全自動生成。
- **柔軟な再生オプション & 動的切り替え**:
  - ループ (`Loop`)、逆再生 (`Reverse`)、往復再生 (`PingPong`)、個別コマ速度設定。
  - オプション構造体での初期設定と、`SetLoop(bool)`, `SetReverse(bool)` などの動的メソッドの両対応。
- **イベントコールバック**:
  - 再生完了時 (`OnComplete`)、特定コマ到達時 (`OnFrame`) のイベント発火。

---

## 使い方

```go
package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sofia-gros/ebiten/animation"
)

type Game struct {
	animMgr   *animation.Manager
	player    *animation.Animator
	chest     *animation.Animator
}

func (g *Game) Init(playerSheet *ebiten.Image, chestImg *ebiten.Image, chestJSON []byte) {
	// 1. マネージャーの作成
	g.animMgr = animation.NewManager()

	// --- パターン A: コードでの手動定義 ---
	clipIdle := animation.NewClip("idle", playerSheet, animation.ClipOptions{
		FPS:  4,
		Loop: true,
	})
	clipAttack := animation.NewClip("attack", playerSheet, animation.ClipOptions{
		FPS:  12,
		Loop: false,
	})

	// 攻撃ヒットコマ (コマ2) でイベント発火
	clipAttack.OnFrame(2, func() {
		println("ヒット判定発火！")
	})

	// 攻撃完了時に idle に戻る
	clipAttack.OnComplete(func() {
		g.player.Play("idle")
	})

	g.player = g.animMgr.CreateAnimator(clipIdle)
	g.player.AddClip(clipAttack)

	// --- パターン B: 汎用 JSON (Aseprite / TexturePacker等) から一発全自動生成 ---
	g.chest, _ = g.animMgr.CreateAnimatorFromJSON(chestImg, chestJSON)
}

func (g *Game) Update() error {
	dt := 1.0 / 60.0

	// 個別アニメーション操作
	if ebiten.IsKeyJustPressed(ebiten.KeyZ) {
		g.player.Play("attack")
	}

	// スローモーション演出 (Manager で一括タイムスケール変更)
	if ebiten.IsKeyPressed(ebiten.KeyShift) {
		g.animMgr.SetSpeed(0.2)
	} else {
		g.animMgr.SetSpeed(1.0)
	}

	// 全アニメーションを一括更新
	g.animMgr.Update(dt)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// 現在のコマ画像を描画
	if img := g.player.CurrentFrame(); img != nil {
		screen.DrawImage(img, nil)
	}
}
```
