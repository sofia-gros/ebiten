package tilemap

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

// AnimatedTileDef は特定のタイルIDに対するアニメーションコマ定義です。
type AnimatedTileDef struct {
	Frames   []int   // コマタイルIDのリスト
	FPS      float64 // アニメーション速度
	timer    float64
	currIdx  int
}

// AnimatedTilemap は水面や風に揺れる草花などアニメーションタイルを保持するレイヤーです。
type AnimatedTilemap struct {
	StaticTilemap
	animDefs map[int]*AnimatedTileDef
}

// NewAnimated は AnimatedTilemap を作成します。
func NewAnimated(width, height, tileWidth, tileHeight int, tileset *Tileset) *AnimatedTilemap {
	return &AnimatedTilemap{
		StaticTilemap: *NewStatic(width, height, tileWidth, tileHeight, tileset),
		animDefs:      make(map[int]*AnimatedTileDef),
	}
}

// SetAnimatedTile は特定のタイルIDにアニメーションコマ定義を登録します。
func (m *AnimatedTilemap) SetAnimatedTile(baseTileID int, frames []int, fps float64) {
	if len(frames) == 0 {
		return
	}
	if fps <= 0 {
		fps = 4.0
	}
	m.animDefs[baseTileID] = &AnimatedTileDef{
		Frames:  frames,
		FPS:     fps,
		timer:   0,
		currIdx: 0,
	}
}

// Update はアニメーションタイルのコマ送りタイマーを更新します。
func (m *AnimatedTilemap) Update(dt float64) {
	for _, def := range m.animDefs {
		if len(def.Frames) <= 1 {
			continue
		}
		def.timer += dt * def.FPS
		frameDur := 1.0
		for def.timer >= frameDur {
			def.timer -= frameDur
			def.currIdx = (def.currIdx + 1) % len(def.Frames)
		}
	}
}

// DrawRegion はアニメーションコマの切り替えを適用して描画します。
func (m *AnimatedTilemap) DrawRegion(screen *ebiten.Image, viewX, viewY, viewW, viewH float64) {
	if !m.visible || screen == nil || m.tileset == nil || m.width <= 0 || m.height <= 0 {
		return
	}

	startCol := int(math.Floor(viewX / float64(m.tileWidth)))
	endCol := int(math.Ceil((viewX + viewW) / float64(m.tileWidth)))
	startRow := int(math.Floor(viewY / float64(m.tileHeight)))
	endRow := int(math.Ceil((viewY + viewH) / float64(m.tileHeight)))

	if startCol < 0 { startCol = 0 }
	if endCol > m.width { endCol = m.width }
	if startRow < 0 { startRow = 0 }
	if endRow > m.height { endRow = m.height }

	opts := &ebiten.DrawImageOptions{}

	for r := startRow; r < endRow; r++ {
		for c := startCol; c < endCol; c++ {
			tileID := m.data[r*m.width+c]
			if tileID <= 0 {
				continue
			}

			// アニメーションコマの適用
			renderTileID := tileID
			if def, ok := m.animDefs[tileID]; ok && len(def.Frames) > 0 {
				renderTileID = def.Frames[def.currIdx]
			}

			img := m.tileset.TileImage(renderTileID)
			if img != nil {
				opts.GeoM.Reset()
				opts.GeoM.Translate(float64(c*m.tileWidth), float64(r*m.tileHeight))
				screen.DrawImage(img, opts)
			}
		}
	}
}
