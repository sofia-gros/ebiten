# ebitentween

[English](./README_en.md)

ebitentween は Ebitengine および汎用 Go 開発向けの、安全・高機能かつ直感的なイージング＆アニメーション補間ライブラリです。
ポインタ直接書き換えの危険性を排除した `OnUpdate(func(val))` / `OnRun(func(progress))` コールバック方式で、オブジェクトやプロパティの構造を問わず数値アニメーションを適用できます。

---

## 特徴

- **安全で柔軟な値の受取 (`OnUpdate` / `OnRun`)**:
  - `OnUpdate(func(val float64))`: 補間された数値を受け取る (`player.X = val` や `cam.SetZoom(val)` など自由な反映が可能)。
  - `OnRun(func(progress float64))`: イージング適用前の進捗率 (`0.0 〜 1.0`) を受け取る。
- **定義と明示的再生の分離 (`tween.New` & `.Play()`)**:
  - `tween.New(&Option{...})` で事前にアニメーション設定を定義。
  - `.Play()` メソッドを呼んだ（またはチェーンした）瞬間に自動的にマネージャーへ登録されて再生開始。
- **インスタンスごとの動的コントロール (`Pause`, `Resume`, `Restart`, `Stop`)**:
  - `tw.Pause()`, `tw.Resume()`, `tw.Restart()`, `tw.Stop()`, `tw.Progress()` により、リアルタイムなゲームイベントに応じた再生制御が可能。
- **カテゴリ別のグループ一括管理 (`Group`)**:
  - `uiGroup := tween.NewGroup()` を使い、UI用・演出用・エネミー用などのグループ単位で `uiGroup.PauseAll()`, `uiGroup.Clear()` の一括停止・再生・削除が可能。
- **自動進行のグローバルデフォルトマネージャー (`DefaultManager`)**:
  - `tween.Update(dt)` をゲームの Update ループに1行書くだけで、特定のグループに属さない全 Tween の時間更新・自動破棄を制御。
- **豊富な標準イージング関数**:
  - Linear, Quad, Cubic, Quart, Quint, Sine, Expo, Circ, Elastic, Back, Bounce (In / Out / InOut)

---

## インストール

```bash
go get github.com/sofia-gros/ebiten/tween
```

---

## 使い方

### クイックスタート

1回限りのアニメーションや、ワンショットで数値を補間再生させるシンプルな記述方法です。


```go
package main

import (
	"fmt"
	"github.com/sofia-gros/ebiten/tween"
)

type Game struct {
	playerX float64
	tw      *tween.Tween
}

func (g *Game) Init() {
	// 定義 ＋ .Play() チェーンで即座に再生開始！
	g.tw = tween.New(&tween.Option{
		Start:    0.0,
		End:      500.0,
		Duration: 1.5,
		Ease:     tween.EaseOutBounce,
		Yoyo:     true, // 往復再生
		Loop:     -1,   // 無限リピート
	}).OnUpdate(func(val float64) {
		g.playerX = val
	}).Play()
}

func (g *Game) Update() error {
	dt := 1.0 / 60.0

	// グローバルマネージャーの全 Tween を時間進行 (Update ループ内で呼ぶだけ)
	tween.Update(dt)

	// ポーズボタンが押されたら一時停止
	if isPaused {
		g.tw.Pause()
	} else {
		g.tw.Resume()
	}

	return nil
}
```

---

### 全機能の使い方

UI用のアニメーション群などを `Group` にまとめて一括管理し、`PauseAll` や `Clear` を行う全機能の使い方です。


```go
package main

import (
	"fmt"
	"github.com/sofia-gros/ebiten/tween"
)

type Game struct {
	uiGroup *tween.Group
	uiAlpha float64
	panelX  float64
}

func (g *Game) Init() {
	// UI 専用のグループを作成
	g.uiGroup = tween.NewGroup()

	// 1. フェードインアニメーション
	g.uiGroup.New(tween.Option{
		Start:    0.0,
		End:      1.0,
		Duration: 0.5,
		Ease:     tween.EaseInOutSine,
	}).OnUpdate(func(val float64) {
		g.uiAlpha = val
	}).Play()

	// 2. パネルのスライドインアニメーション
	g.uiGroup.New(tween.Option{
		Start:    -200.0,
		End:      100.0,
		Duration: 0.8,
		Ease:     tween.EaseOutBack,
	}).OnUpdate(func(val float64) {
		g.panelX = val
	}).OnComplete(func() {
		fmt.Println("UI Slide-in Finished!")
	}).Play()
}

func (g *Game) Update() error {
	dt := 1.0 / 60.0

	// グループ個別の時間を更新
	g.uiGroup.Update(dt)

	return nil
}

func (g *Game) OnOpenDialog() {
	// ダイアログが開いたので UI グループのアニメーションを一時停止
	g.uiGroup.PauseAll()
}

func (g *Game) OnCloseScene() {
	// シーン終了時に UI グループの全アニメーションを一括破棄
	g.uiGroup.Clear()
}
```

---

## 主要 API リファレンス

### `tween.New(opts...)` / `tween.FromTo(start, end, duration)`
- `Ease(easeFunc)`: イージング関数をセット。
- `Delay(sec)`: アニメーション開始前の遅延時間を設定。
- `Loop(count)`: リピート回数を設定 (-1 で無限)。
- `Yoyo(bool)`: 往復リピートの切り替え。
- `OnUpdate(fn)`: 補間された数値を受け取るコールバックを登録。
- `OnRun(fn)`: 生の進捗率 (0.0 〜 1.0) を受け取るコールバックを登録。
- `OnComplete(fn)`: アニメーション完了時のコールバックを登録。
- `Play()`: 再生を開始し、マネージャーに自動登録。
- `Pause()` / `Resume()` / `Restart()` / `Stop()`: 動的操作。
- `Progress()`: 現在の進捗率を取得。

### `tween.Group`
- `NewGroup()`: 新しいグループを作成。
- `New(opts...)` / `FromTo(...)`: このグループに属する Tween を作成。
- `Update(dt)`: グループ内の全 Tween を更新。
- `PauseAll()` / `ResumeAll()` / `RestartAll()` / `Clear()`: グループ内の一括操作。

---

## ライセンス

MIT License
