# ebitenui

[English](./README_en.md)

ebitenui は Ebitengine 向けの 2D ゲーム UI ライブラリです。
他のライブラリに一切依存せず 100% 単体で動作しつつ、`camera` ライブラリ（画面固定バインド）や `pad` ライブラリ（`input.Action` 数値によるフォーカス＆決定操作）とオプショナル接続が可能です。

---

## 特徴

- **豊富な UI コンポーネント群**:
  - `NineSlice`: 角を歪ませず綺麗に拡大縮小する 9分割枠パネル画像。
  - `Button`: 通常 (`Normal`), ホバー (`Hover`), 押下 (`Pressed`), 無効 (`Disabled`) の状態別画像切替 ＋ アイコン画像・テキストの自由合成。
  - `TextInput`: 文字入力、カーソル点滅、バックスペース、`OnSubmit` 送信。
  - `ScrollBox`: 枠内コンテンツの切り抜き & マウスホイール/ドラッグ表示。
  - `Slider`, `CheckBox`, `Label`, `Panel`, `VBox`, `HBox`, `Container`
- **100% 同等の Option 構造体 ＋ ゲッター/セッター**:
  - メソッドで設定できる全項目 (`SetPos`, `Pos`, `SetSize`, `Size`, `SetText`, `Text`, `SetGrayscale`, `SetNormalImage` 等) が `Option` 構造体でも一括指定可能。
- **明示的なモノクロ化オプション (`SetGrayscale`)**:
  - 自動グレーアウト処理は行わず、ユーザー指定画像を素直に描画。明示的に白黒モノクロ化したい場合のみ `SetGrayscale(true)` で適用可能。
- **シンプルな要素アクセス (`GetAll`, `Get`, `Remove`)**:
  - `uiRoot.GetAll()` で全要素取得、`uiRoot.Get(index)` で目的の UI 要素を直接取得。
- **オプショナル接続 (`BindCamera`, `BindEbitenInput`, `BindPadInput`)**:
  - カメラが揺れても UI を画面固定に保つ `BindCamera(cam)`。
  - `pad` ライブラリの `input.Action` (数値) をバインドしてゲームパッドフォーカス操作する `BindPadInput`。

---

## インストール

```bash
go get github.com/sofia-gros/ebiten/ui
```

---

## 使い方

### クイックスタート

`pad` や `camera` なしで、Ebitengine の標準入力・画面描画で UI を作成するシンプルな記述方法です。


```go
package main

import (
	"fmt"
	"image/color"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sofia-gros/ebiten/ui"
)

type Game struct {
	uiRoot *ui.Container
}

func (g *Game) Init() {
	g.uiRoot = ui.NewContainer()
	g.uiRoot.SetPos(50, 50)

	// 縦並びレイアウト (VBox)
	vbox := ui.NewVBox()
	vbox.SetSpacing(15)

	// 1. ラベル
	titleLabel := ui.NewLabel("メインメニュー", ui.LabelOption{FontSize: 20})

	// 2. ボタン
	btn := ui.NewButton(ui.ButtonOption{
		Text:   "ゲーム開始",
		Width:  200,
		Height: 45,
	})
	btn.OnClick(func() {
		fmt.Println("ゲーム開始が押されました！")
	})

	vbox.Add(titleLabel)
	vbox.Add(btn)
	g.uiRoot.Add(vbox)
}

func (g *Game) Update() error {
	// 単体での入力・ホバー状態更新
	g.uiRoot.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// 画面へ UI 描画
	g.uiRoot.Draw(screen)
}
```

---

### 全機能の使い方

文字入力ボックス、スクロール枠、9スライスパネル枠、カメラ画面固定、`pad` ライブラリとの入力バインドを行う全機能の使い方です。


```go
package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sofia-gros/ebiten/camera"
	"github.com/sofia-gros/ebiten/pad/input"
	"github.com/sofia-gros/ebiten/ui"
)

const (
	ActionUp     input.Action = 1
	ActionDown   input.Action = 2
	ActionSubmit input.Action = 3
)

type Game struct {
	uiRoot *ui.Container
	cam    *camera.Camera
}

func (g *Game) Init() {
	g.uiRoot = ui.NewContainer()
	g.cam = camera.New(640, 480)

	// 1. 文字入力ボックス (TextInput)
	nameInput := ui.NewTextInput(ui.TextInputOption{
		Placeholder: "名前を入力...",
		Width:       220,
		Height:      35,
	})
	nameInput.OnSubmit(func(text string) {
		fmt.Println("名前設定:", text)
	})

	// 2. スクロール枠 (ScrollBox: 200x150 領域)
	scrollBox := ui.NewScrollBox(200, 150)
	scrollBox.Add(nameInput)

	g.uiRoot.Add(scrollBox)

	// --- オプション接続 ---
	// A. カメラが揺れても UI は画面固定追従！
	g.uiRoot.BindCamera(g.cam)

	// B. pad ライブラリの Action (数値) をバインドしてゲームパッド操作対応！
	in := input.NewInput()
	in.BindKey(ActionUp, ebiten.KeyArrowUp)
	in.BindKey(ActionDown, ebiten.KeyArrowDown)
	in.BindKey(ActionSubmit, ebiten.KeyEnter)

	g.uiRoot.BindPadInput(in, ui.PadInputMapping{
		UpAction:     ActionUp,
		DownAction:   ActionDown,
		SubmitAction: ActionSubmit,
	})
}

func (g *Game) Update() error {
	g.uiRoot.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.uiRoot.Draw(screen)
}
```

---

## 主要 API リファレンス

### 共通 `ui.Element` メソッド
- `SetPos(x, y)` / `Pos() (x, y)`: 位置設定・取得。
- `SetSize(w, h)` / `Size() (w, h)`: サイズ設定・取得。
- `SetVisible(bool)` / `Visible() bool`: 表示・非表示切り替え。
- `SetEnabled(bool)` / `Enabled() bool`: 有効・無効切り替え。
- `SetGrayscale(bool)` / `IsGrayscale() bool`: モノクロ描画の明示的オン/オフ。

### コンポーネント別
- **`ui.Button`**: `SetText`, `Text`, `SetNormalImage`, `SetHoverImage`, `SetPressedImage`, `SetDisabledImage`, `SetIconImage`, `OnClick(fn)`
- **`ui.TextInput`**: `SetText`, `Text`, `OnSubmit(fn)`
- **`ui.Slider`**: `SetValue`, `Value`, `OnChange(fn)`
- **`ui.CheckBox`**: `SetChecked`, `Checked`, `OnChange(fn)`
- **`ui.Container` / `ui.VBox`**: `GetAll()`, `Get(index)`, `Remove(elem)`, `SetAllVisible(bool)`, `BindCamera(cam)`, `BindPadInput(in, mapping)`, `BindEbitenInput(mapping)`

---

## ライセンス

MIT License
