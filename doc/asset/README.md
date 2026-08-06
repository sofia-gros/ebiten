# ebitenasset

[English](./README_en.md)

ebitenasset は Ebitengine 向けの直感的なアセットロード・キャッシュ管理ライブラリです。
画像、コマ割りスプライトシート、音声、Tilemap マップデータ、アニメーション定義、JSON/YAML データを Phaser 風のシンプルな操作感で読み込み・取得・解放できます。

---

## 特徴

- **シンプルなアセット管理**: `Loader` に画像・音声・データを登録して一括ロード。
- **2通りの画像ロード (`Image` vs `Sprite`)**:
  - `Image`: 背景や一枚絵用の単純なテクスチャ。
  - `Sprite`: アニメーションやコマ割り（FrameWidth, FrameHeight）を持つスプライトシート。
- **構造体データへの即時変換 (`DataAs[T]`)**:
  - JSON や YAML データファイルをロード後、Go の構造体へ型安全にアンマーシャル。
- **メモリ解放制御 (`Unload` / `Clear`)**:
  - 不要になった特定アセットやシーン終了時のキャッシュ一括削除。

---

## インストール

```bash
go get github.com/sofia-gros/ebiten/asset
```

---

## 使い方

### クイックスタート

一枚絵の背景画像や JSON データの読み込み、取得、描画を行う基本的なコード例です。


```go
package main

import (
	"embed"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sofia-gros/ebiten/asset"
)

//go:embed assets/*
var assetsFS embed.FS

type GameConfig struct {
	Title   string  `json:"title"`
	Volume  float64 `json:"volume"`
}

type Game struct {
	m *asset.Manager
}

func (g *Game) Init() {
	// Manager の作成 (FS を指定)
	g.m = asset.NewManager(assetsFS)

	// アセットの事前登録
	g.m.Image("bg", "assets/background.png")
	g.m.Data("config", "assets/config.json")

	// ロードの実行
	if err := g.m.Load(); err != nil {
		panic(err)
	}

	// JSON データを型安全に構造体へ展開
	var cfg GameConfig
	if err := asset.DataAs(g.m, "config", &cfg); err == nil {
		println("Game Title:", cfg.Title)
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	// キャッシュから画像を取得して描画
	if bgImg := g.m.GetImage("bg"); bgImg != nil {
		screen.DrawImage(bgImg, nil)
	}
}
```

---

### 全機能の使い方

コマ割りスプライトシート、アニメーション定義、BGM/SE 音声データ、Tilemap マップデータを含む全機能の使い方です。


```go
package main

import (
	"embed"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sofia-gros/ebiten/asset"
)

//go:embed assets/*
var assetsFS embed.FS

type Game struct {
	m *asset.Manager
}

func (g *Game) Init() {
	g.m = asset.NewManager(assetsFS)

	// 1. コマ割りスプライトシート (32x32px のコマ)
	g.m.Sprite("player", "assets/player_sheet.png", asset.SpriteOptions{
		FrameWidth:  32,
		FrameHeight: 32,
	})

	// 2. アニメーションデータ・Tilemap
	g.m.Animation("player_anim", "assets/player_anim.json")
	g.m.Tilemap("level1", "assets/level1.json")

	// 3. 音声ファイル (BGM / SE)
	g.m.Audio("bgm_main", "assets/bgm.mp3")
	g.m.Audio("se_jump", "assets/jump.wav")

	// 一括ロード実行
	if err := g.m.Load(); err != nil {
		panic(err)
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	// コマ割りスプライトから特定フレーム (コマ番号 2) を切り出して描画
	if sheet := g.m.GetSprite("player"); sheet != nil {
		frameImg := sheet.GetFrame(2)
		screen.DrawImage(frameImg, nil)
	}
}

func (g *Game) OnSceneChange() {
	// シーン切り替え時に不要アセットをメモリ解放
	g.m.Unload("level1")
}
```

---

## 主要 API リファレンス

### `asset.Manager`
- `NewManager(fs embed.FS)`: アセットマネージャーを作成。
- `Image(key, path)`: 単一画像を登録。
- `Sprite(key, path, SpriteOptions)`: コマ割りスプライトシートを登録。
- `Audio(key, path)`: BGM / SE 音声データを登録。
- `Data(key, path)` / `Tilemap(key, path)` / `Animation(key, path)`: JSON/YAML データを登録。
- `Load() error`: 登録された全アセットを一括ロード。
- `GetImage(key)`: ロード済み画像を `*ebiten.Image` として取得。
- `GetSprite(key)`: コマ割りスプライトシート (`*asset.SpriteSheet`) を取得。
- `GetAudio(key)`: 音声データバイト列 (`[]byte`) を取得。
- `Unload(key)` / `Clear()`: 特定アセットまたは全キャッシュのメモリ解放。

---

## ライセンス

MIT License
