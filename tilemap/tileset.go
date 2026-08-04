package tilemap

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

// TileSource は Frame(index) でコマ画像を取得できる型を表すインターフェースです。
// (asset.SpriteSheet 等も自動的に適合します)
type TileSource interface {
	Frame(index int) *ebiten.Image
}

// Tileset は 1 枚のタイルセット画像または TileSource から切り出されたタイルコマ画像を管理します。
type Tileset struct {
	baseImage  *ebiten.Image
	tileWidth  int
	tileHeight int
	margin     int
	spacing    int
	tiles      map[int]*ebiten.Image
}

// NewTileset はベース画像、タイル幅・高さから Tileset を生成します。
func NewTileset(base *ebiten.Image, tileWidth, tileHeight int, opts ...int) *Tileset {
	margin := 0
	spacing := 0
	if len(opts) > 0 {
		margin = opts[0]
	}
	if len(opts) > 1 {
		spacing = opts[1]
	}

	ts := &Tileset{
		baseImage:  base,
		tileWidth:  tileWidth,
		tileHeight: tileHeight,
		margin:     margin,
		spacing:    spacing,
		tiles:      make(map[int]*ebiten.Image),
	}

	if base == nil || tileWidth <= 0 || tileHeight <= 0 {
		return ts
	}

	bounds := base.Bounds()
	imgW := bounds.Dx()
	imgH := bounds.Dy()

	tileID := 0
	for y := margin; y+tileHeight <= imgH-margin; y += tileHeight + spacing {
		for x := margin; x+tileWidth <= imgW-margin; x += tileWidth + spacing {
			subImg := base.SubImage(image.Rect(x, y, x+tileWidth, y+tileHeight)).(*ebiten.Image)
			ts.tiles[tileID] = subImg
			tileID++
		}
	}

	return ts
}

// NewTilesetFromSource は TileSource インターフェースから Tileset を生成します。
func NewTilesetFromSource(source TileSource, tileWidth, tileHeight int, maxTiles ...int) *Tileset {
	ts := &Tileset{
		tileWidth:  tileWidth,
		tileHeight: tileHeight,
		tiles:      make(map[int]*ebiten.Image),
	}

	if source == nil {
		return ts
	}

	limit := 256
	if len(maxTiles) > 0 && maxTiles[0] > 0 {
		limit = maxTiles[0]
	}

	for i := 0; i < limit; i++ {
		img := source.Frame(i)
		if img != nil {
			ts.tiles[i] = img
		} else {
			break
		}
	}

	return ts
}

// TileImage はタイルIDに対応する *ebiten.Image を取得します。
func (ts *Tileset) TileImage(id int) *ebiten.Image {
	if ts == nil {
		return nil
	}
	return ts.tiles[id]
}

// TileWidth は1タイルの幅を返します。
func (ts *Tileset) TileWidth() int {
	if ts == nil {
		return 0
	}
	return ts.tileWidth
}

// TileHeight は1タイルの高さを返します。
func (ts *Tileset) TileHeight() int {
	if ts == nil {
		return 0
	}
	return ts.tileHeight
}
