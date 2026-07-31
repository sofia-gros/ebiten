# ebitensound

[English](./README_en.md)

ebitensound は Ebitengine 向けの直感的なサウンド管理ライブラリです。
BGM や SE、ボイスの再生・音量制御、クロスフェード、リアルタイム 2D ポジショナルサウンド（音源追従）、フェードアウト停止を一元管理できます。

---

## 特徴

- **サウンドタイプ拡張 (`Type`)**:
  - `sound.TypeSE`, `sound.TypeBGM`, `sound.TypeVoice`, `sound.TypeEnv` などの標準タイプに加え、`sound.TypeCustom + iota` で独自のサウンドタイプを追加可能。
- **タイプ別音量制御 (`SetVolume`)**:
  - 設定画面などで「SE音量」「BGM音量」をタイプごとに独立して括り調整可能。
- **シンプルなフェード停止 & クロスフェード (`Stop`, `CrossFade`)**:
  - 指定秒数かけて滑らかにフェードアウト停止、または別BGMへクロスフェード移行。
- **2D ポジショナルサウンド (`PlayAt` / `SetPosition`)**:
  - プレイヤーと音源（敵や爆発など）の距離に応じた左右パン (`Pan`) と音量減衰をリアルタイム追従更新。

---

## インストール

```bash
go get github.com/sofia-gros/ebiten/sound
```

---

## 使い方

### クイックスタート

効果音 (SE) を即座に再生したり、BGM をフェードアウト停止させる基本的なコード例です。


```go
package main

import (
	"github.com/sofia-gros/ebiten/sound"
)

func main() {
	// 1. サウンドマネージャーの作成 (サンプルレート: 44100Hz)
	m := sound.NewManager(44100)

	// 2. SE の即時再生 (シンプルアプローチ)
	seTrack := m.Play(seBytes)
	_ = seTrack

	// 3. オプション付き BGM の再生 (ループ再生 + 音量 0.8)
	bgmTrack := m.Play(bgmBytes, sound.Option{
		Type:   sound.TypeBGM,
		Volume: 0.8,
		Loop:   true,
	})

	// 4. 1.5秒かけてフェードアウト停止
	m.Stop(bgmTrack, 1.5)
}
```

---

### 全機能の使い方

ゲーム設定画面でのタイプ別音量調整、音源追従 (2D Positional Audio)、およびシーン切り替え時の BGM クロスフェードを行う全機能の使い方です。


```go
package main

import (
	"github.com/sofia-gros/ebiten/sound"
)

// 独自サウンドタイプの定義
const (
	TypeBattle sound.Type = sound.TypeCustom + iota // カスタムタイプ 4
)

type Game struct {
	m          *sound.Manager
	enemyTrack *sound.Track
}

func (g *Game) Init() {
	g.m = sound.NewManager(44100)

	// タイプごとのマスター音量調整 (ゲーム設定画面などで使用)
	g.m.SetVolume(sound.TypeSE, 0.7)  // SE 全体を 70% に設定
	g.m.SetVolume(sound.TypeBGM, 0.5) // BGM 全体を 50% に設定

	// 2D ポジショナルサウンドの再生開始
	// (ターゲット位置 x:200, y:100 / 聴き手・プレイヤー位置 x:100, y:100)
	g.enemyTrack = g.m.PlayAt(engineSEBytes, 200, 100, 100, 100, sound.Option{
		Type: sound.TypeSE,
		Loop: true,
	})
}

func (g *Game) Update() error {
	dt := 1.0 / 60.0

	// 敵が右へ移動したため、音源の 2D 位置をリアルタイム更新！ (左右パンと音量減衰が自動変更される)
	enemyX, enemyY := 350.0, 100.0
	playerX, playerY := 100.0, 100.0
	g.enemyTrack.SetPosition(enemyX, enemyY, playerX, playerY)

	return nil
}

func (g *Game) ChangeToBossScene(nextBgmBytes []byte) {
	// 次の BGM へ 2.0秒かけてスムーズにクロスフェード移行
	g.m.CrossFade(nextBgmBytes, 2.0, sound.Option{
		Type: sound.TypeBGM,
		Loop: true,
	})
}
```

---

## 主要 API リファレンス

### `sound.Manager`
- `NewManager(sampleRate int)`: サウンドマネージャーを作成。
- `Play(audioBytes, opts...)`: サウンドを即時再生。戻り値 `*Track` を受け取る。
- `PlayAt(audioBytes, targetX, targetY, listenerX, listenerY, opts...)`: 2D ポジショナルサウンドとして再生。
- `CrossFade(nextAudioBytes, fadeSec, opts...)`: 現在再生中の BGM から新しい BGM へ指定秒数でクロスフェード移行。
- `Stop(trackOrType, ...fadeSec)`: 特定の `*Track` または `sound.Type` 全体を指定秒数でフェードアウト停止。
- `SetVolume(soundType, volume)`: 指定したサウンドタイプ全体のマスター音量 (`0.0 〜 1.0`) を設定。
- `PauseAll()` / `ResumeAll()`: 全トラックの一時停止 / 再開。

### `sound.Track`
- `SetPosition(targetX, targetY, listenerX, listenerY)`: 音源と聴き手の位置を更新。
- `SetVolume(vol)`: このトラック個別の音量を変更。
- `Stop(...fadeSec)`: このトラックを停止。

---

## ライセンス

MIT License
