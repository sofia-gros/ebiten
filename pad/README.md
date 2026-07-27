# ebitenpad

ebitenpad は Ebitengine (ebiten) 向けの、アクションベースの入力管理ライブラリです。
キーボード、ゲームパッド、そしてバーチャルパッド（タッチ操作）を一つの論理的な「アクション」に統合して扱うことができます。

## 特徴

- **アクションベースの抽象化**: ゲームロジックをデバイス（キーボード、ゲームパッド、タッチ）から切り離します。
- **バーチャルパッド内蔵**: スティックやボタンなどのUIを簡単に作成し、マルチタッチで操作可能です。
- **シンプルなAPI**: `Pressed`, `JustPressed`, `JustReleased` などの直感的なメソッドを提供します。
- **デバイス統合**: 1つのアクションに対して複数のデバイス入力をバインドできます。

## アップデート情報

v1.1.0 2026/07/06
- 1画面分割対戦のように、プレイヤーごとに入力を分ける機能を追加
- Strengthが実際の入力値を使用するように変更

## インストール

```bash
go get github.com/sofia-gros/ebiten/pad
```

## 基本的な使い方

```go
package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/sofia-gros/ebiten/pad/input"
)

const (
	ActionJump input.Action = 1
	ActionMove input.Action = 2
)

type Game struct {
	in *input.Input
}

func (g *Game) Update() error {
	g.in.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	vx, vy := 0.0, 0.0
	strength := 0.0
	if state, ok := g.in.GetActionState(ActionMove); ok {
		vx, vy = state.X, state.Y
		strength = state.Strength
	}

	// バーチャル UI の描画
	g.in.Virtual().Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 480
}

func main() {
	in := input.NewInput()

	// Bind Keyboard
	in.BindKey(ActionJump, ebiten.KeySpace)
	in.BindKeyAxis(ActionMove, ebiten.KeyA, ebiten.KeyD, ebiten.KeyW, ebiten.KeyS)

	// Bind Gamepad
	in.BindGamepadButton(ActionJump, ebiten.StandardGamepadButtonRightBottom)
	in.BindGamepadAxis(ActionMove, 0, 1)

	// Virtual Padの設定
	vpad := in.Virtual()
	jumpBtn := vpad.AddButton().SetPosition(550, 400).SetRadius(40)
	moveStick := vpad.AddStick().SetPosition(100, 380).SetRadius(60)

	in.BindButton(ActionJump, jumpBtn)
	in.BindStick(ActionMove, moveStick)

	g := &Game{in: in}

	ebiten.SetWindowTitle("ebitenpad WASM Example")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
```

## ライブラリの構成

- `input`: メインの入力マネージャー。アクションのバインドと状態管理を行います。
- `virtual`: バーチャルスティックやボタンなどのUIコンポーネント。
- `examples/wasm`: WebAssembly で動作するデモコード。

## 複数プレイヤー対応 (`Controller` / `For`)

1画面分割対戦のように、プレイヤーごとに入力を分けたい場合は `Controller` 型と `For()` を使います。

```go
const (
    P1 input.Controller = 0 // gamepad 0 にも対応
    P2 input.Controller = 1 // gamepad 1 にも対応
)

const (
    ActionJump input.Action = 1
    ActionMove input.Action = 2
)

// セットアップ
in.For(P1).BindKey(ActionJump, ebiten.KeySpace)
in.For(P1).BindKeyAxis(ActionMove, ebiten.KeyA, ebiten.KeyD, ebiten.KeyW, ebiten.KeyS)
in.For(P1).BindGamepadButton(ActionJump, ebiten.StandardGamepadButtonRightBottom) // gamepad 0

in.For(P2).BindKey(ActionJump, ebiten.KeyEnter)
in.For(P2).BindKeyAxis(ActionMove, ebiten.KeyLeft, ebiten.KeyRight, ebiten.KeyUp, ebiten.KeyDown)
in.For(P2).BindGamepadButton(ActionJump, ebiten.StandardGamepadButtonRightBottom) // gamepad 1

// Update は1回だけ
in.Update()

// 状態取得 (P1/P2 は完全に独立。同時入力しても上書きなし)
if state, ok := g.in.For(P1).GetActionState(ActionMove); ok {
    p1vx, p1vy = state.X, state.Y
}
if state, ok := g.in.For(P2).GetActionState(ActionMove); ok {
    p2vx, p2vy = state.X, state.Y
}
```

1人プレイの場合は `For()` を省略するだけで従来通りに動きます。`Controller` の値はそのまま gamepad のインデックスとして扱われます。

```go
// 1人プレイ: 従来通り
in.BindKey(ActionJump, ebiten.KeySpace)
in.GetActionState(ActionJump)
```

## `state.Strength` について

`ActionState.Strength` はアクションの入力の強さを `0.0 〜 1.0` で表します。デバイスによって計算方法が異なります。

| デバイス                   | 計算方法                                                                |
| -------------------------- | ----------------------------------------------------------------------- |
| キーボード（単キー）       | 押していれば常に `1.0`                                                  |
| キーボード（軸）           | `√(dx² + dy²)` を `1.0` で clamp。斜め入力でも `1.0` が上限             |
| ゲームパッド（ボタン）     | 押していれば常に `1.0`                                                  |
| ゲームパッド（アナログ軸） | `√(x² + y²)` を `1.0` で clamp。アナログ入力の強さが実際に反映される    |
| バーチャルスティック       | 中心からの距離を半径で割った値。指を少し動かすだけなら `0.3` とかになる |

## ゲームパッドのデッドゾーン

スティックがニュートラルでも微妙な値が出るのはよくある話で、`BindGamepadAxisWithDeadzone` を使うとその辺を処理できます。

```go
// Strength が 0.2 以下の入力は無視する
in.BindGamepadAxisWithDeadzone(ActionMove, 0, 1, 0.2)
```

`BindGamepadAxis` はデッドゾーンなし（`0.0`）と同じです。

## 複数ゲームパッド

接続されているゲームパッドを全部巡回して、一番大きい入力を採用します。コントローラーを2本つないでいてもどちらかが動けば動く、くらいの感じです。

## ライセンス

MIT License
