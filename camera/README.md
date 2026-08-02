# ebitencamera

[English](./README_en.md)

ebitencamera は Ebitengine 向けの 2D カメラライブラリです。
特定のアセットや物理エンジン構造体に一切依存せず、純粋な座標制御 (`SetPos`)、ズーム、回転、画面振動 (Shake)、可視領域境界制限 (Bounds)、ビューポート切り抜き、Custom Shader 自動切替、複数カメラの ZIndex 自動ソート描画を行えます。

---

## 特徴

- **純粋な座標制御 (`SetPos` / `Move` / `MoveTo`)**: 特定のオブジェクト構造体に固定されず、プレイヤーや敵の座標を直接指定してカメラを操作可能。
- **柔軟なクロージャ描画 (`cam.Render`)**:
  - コールバッククロージャ内で描画処理を実行。`pad.Draw` や物理エンジンのデバッグ描画など、独自の描画メソッドを持つライブラリ・コンポーネントをそのままカメラ空間内で描画。
  - 中間バッファ生成が不要な高速ダイレクト描画。
- **カメラ状態としてのシェーダー自動切替 (`SetShader` / `ClearShader`)**:
  - `cam.SetShader(shader, opts)` を呼ぶと、描画コードを変更することなく自動的にシェーダー描画へ切り替わります。
- **マルチカメラ & ZIndex 自動ソート (`Group`)**:
  - 右上のミニマップ表示や画面分割など複数のカメラを `Group` に追加すると、`ZIndex` 優先度順に自動ソートして一括描画できます。
- **相互座標変換 & カリング機能**:
  - `ScreenToWorld(x, y)` / `WorldToScreen(x, y)` による双方向座標変換。
  - `VisibleBounds()` で現在画面内に映っているワールド領域 `(minX, minY, maxX, maxY)` を取得し、画面外オブジェクトの描画スキップ（カリング）が容易。

---

## インストール

```bash
go get github.com/sofia-gros/ebiten/camera
```

---

## 使い方

### クイックスタート

単一のカメラを使ってプレイヤーに追従し、クロージャ内でマップやプレイヤーを描画する基本的な記述方法です。


```go
package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sofia-gros/ebiten/camera"
)

type Game struct {
	cam     *camera.Camera
	playerX float64
	playerY float64
}

func (g *Game) Init() {
	// 幅 640, 高さ 480 のカメラを作成
	g.cam = camera.New(640, 480)
}

func (g *Game) Update() error {
	dt := 1.0 / 60.0

	// プレイヤーの座標へ滑らかにカメラを追従移動 (Lerp 0.1)
	g.cam.MoveTo(g.playerX, g.playerY, 0.1)

	// 画面振動の更新など
	g.cam.Update(dt)

	// マウス位置のスクリーン座標をワールド座標へ変換
	mx, my := ebiten.CursorPosition()
	worldX, worldY := g.cam.ScreenToWorld(float64(mx), float64(my))
	_, _ = worldX, worldY

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Render クロージャ内で描画を行うことで、カメラ空間が自動適用される
	g.cam.Render(screen, func(target *ebiten.Image) {
		// マップやスプライト、pad など独自の Render メソッドを直接描画
		g.drawMap(target)
		g.drawPlayer(target)
	})
}
```

---

### 全機能の使い方

メイン画面と右上のミニマップ用カメラを `Group` で管理し、状態に応じた Custom Shader の自動切替を行う全機能の使い方です。


```go
package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sofia-gros/ebiten/camera"
)

type Game struct {
	group     *camera.Group
	mainCam   *camera.Camera
	miniCam   *camera.Camera
	darkShader *ebiten.Shader
}

func (g *Game) Init() {
	// メインカメラ (ZIndex: 0)
	g.mainCam = camera.New("main", 640, 480)
	g.mainCam.SetZIndex(0)

	// ミニマップ用カメラ (ZIndex: 10 で最前面表示、画面右上 160x120 枠)
	g.miniCam = camera.New("mini", 640, 480)
	g.miniCam.SetViewport(460, 20, 160, 120)
	g.miniCam.SetZoom(0.2) // 広域表示
	g.miniCam.SetZIndex(10)

	// カメラグループの作成
	g.group = camera.NewGroup(g.mainCam, g.miniCam)
}

func (g *Game) Update() error {
	dt := 1.0 / 60.0

	// 両方のカメラ位置をプレイヤーに設定
	g.mainCam.SetPos(playerX, playerY)
	g.miniCam.SetPos(playerX, playerY)

	// 暗闇エリアに入ったらメインカメラにシェーダーを適用
	if inDarkArea {
		g.mainCam.SetShader(g.darkShader)
	} else {
		g.mainCam.ClearShader() // 通常表示に戻す
	}

	g.group.Update(dt)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// ZIndex 順 (メイン ➔ ミニマップ) に自動ソートして一括描画
	g.group.Render(screen, func(cam *camera.Camera, target *ebiten.Image) {
		g.drawWorld(target)
	})
}
```

---

## 主要 API リファレンス

### `camera.Camera`
- `New(width, height)`: 新しいカメラを作成。
- `SetPos(x, y)` / `Pos() (x, y)`: カメラ中心位置の設定・取得。
- `Move(dx, dy)` / `MoveTo(targetX, targetY, speed)`: カメラの移動・滑らかな補間移動。
- `SetZoom(zoom)` / `Zoom()`: ズーム倍率の設定・取得。
- `SetRotation(angle)` / `Rotation()`: 回転角度の設定・取得。
- `Forward() (dirX, dirY)` / `Right() (rtX, rtY)`: カメラの向き（正面・右方向）の単位ベクトルを取得（回転時の相対 WASD 移動などに利用）。
- `Direction() (dirX, dirY)`: カメラの向きベクトルを取得。
- `SetBounds(minX, minY, maxX, maxY)`: ワールド移動境界制限の設定。
- `SetViewport(x, y, w, h)`: 画面出力領域の設定 (ミニマップ用)。
- `Shake(strength, durationSec)`: 画面振動を発生させる。
- `SetShader(shader, opts...)` / `ClearShader()`: カスタムシェーダーの適用・解除。
- `ScreenToWorld(screenX, screenY)` / `WorldToScreen(worldX, worldY)`: 座標相互変換。
- `VisibleBounds()`: 現在カメラ内に見えているワールド領域 `(minX, minY, maxX, maxY)` を取得。
- `Render(screen, drawFunc)`: カメラのトランスフォーム・ビューポート・シェーダーを適用して描画クロージャを実行。


### `camera.Group`
- `NewGroup(cameras...)`: 複数カメラを管理するグループを作成。
- `Add(cam)` / `Remove(cam)`: カメラの追加・除外。
- `Update(dt)`: グループ内の全カメラを更新。
- `Render(screen, drawFunc)`: ZIndex 順に自動ソートして全カメラを一括描画。

---

## ライセンス

MIT License
