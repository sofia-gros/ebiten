# ebitentilemap

[English](./README_en.md)

ebitentilemap は Ebitengine 向けの軽量かつ高機能な 2D タイルマップライブラリです。
Tiled Map Editor (`.json` / `.tmj`) の全自動インポートに対応しているほか、2次元配列データからのコード直接生成、3種類のタイルマップ構成、視域カリング描画（`DrawRegion`）、直感的なエリア指定クエリ、物理エンジン連携用 AABB 結合矩形算出機能を提供します。

`camera`, `asset`, `physics` などの自作ライブラリには一切依存せず、`ebiten/v2` のみに依存する完全独立モジュールです。

---

## 特徴

- **完全独立設計**: `ebiten/v2` のみに直接依存。他ライブラリなしで単体利用可能。
- **3種類のタイルマップ構成**:
  - `StaticTilemap`: 固定地面・壁用。超高速レンダリング。
  - `AnimatedTilemap`: 水面や揺れる草花のアニメーションタイル。
  - `InfiniteTilemap`: 広大なオープンワールド用チャンク型自動生成。
- **Tiled Map Editor 全自動インポート (`ImportTiledJSON`)**:
  - Tiled 出力 JSON から静的レイヤー・アニメーション・`solid: true` 通行判定を一括自動解析。
- **コード上での直接生成 (`NewStaticFromData`)**:
  - 2次元配列 (`[][]int`) を渡すだけで 1 行で即座にタイルマップを定義。
- **視域カリング描画 (`DrawRegion`)**:
  - カメラ等の表示領域 (`viewX, viewY, viewW, viewH`) 内のタイルのみを動的カリング描画。
- **直感的なエリア指定クエリ (`GetArea`)**:
  - `GetArea(x, y, w, h)` で部分領域を取得し、`area.FindTiles(id)` (範囲内検索) や `area.ReplaceTile(oldID, newID)` (範囲置換・草刈り演出) をチェーン記述。
- **物理エンジン連携 (`CreateCollisionBoxes`)**:
  - `solid: true` の壁タイル群を自動結合した物理 AABB 矩形リスト (`CollisionBox`) を算出。

---

## 使い方

```go
package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sofia-gros/ebiten/tilemap"
)

type Game struct {
	mapGroup *tilemap.MapGroup
	codeMap  *tilemap.StaticTilemap
	playerX  float64
	playerY  float64
}

func (g *Game) Init(tilesetImg *ebiten.Image, tiledJSONBytes []byte) {
	// 1. Tiled JSON からの一括全自動インポート
	g.mapGroup, _ = tilemap.ImportTiledJSON(tiledJSONBytes, map[string]*ebiten.Image{
		"tileset.png": tilesetImg,
	})

	// 2. 2次元配列からコード上で直接マップ生成
	tileset := tilemap.NewTileset(tilesetImg, 16, 16)
	mapData := [][]int{
		{1, 1, 1, 1},
		{1, 0, 9, 1}, // 9: 宝箱
		{1, 1, 1, 1},
	}
	g.codeMap = tilemap.NewStaticFromData(mapData, tileset)
	g.codeMap.SetTileSolid(1, true)
}

func (g *Game) Update() error {
	dt := 1.0 / 60.0

	// 直感的なピクセル指定通行判定
	if !g.mapGroup.IsSolidAtPixel(g.playerX+2, g.playerY) {
		g.playerX += 2
	}

	// エリア指定クエリ (草刈り演出など)
	if ebiten.IsKeyJustPressed(ebiten.KeySpace) {
		area := g.mapGroup.GetArea(g.playerX-32, g.playerY-32, 64, 64)
		area.ReplaceTile(5, 6) // 周辺の草(5)を刈った草(6)に置換
	}

	g.mapGroup.Update(dt)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// 表示領域のカリング描画 (画面サイズ 640x480)
	g.mapGroup.DrawRegion(screen, g.playerX-320, g.playerY-240, 640, 480)
}
```
