package tilemap

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

// StaticTilemap は動かない静的なタイルマップレイヤーを管理します。
type StaticTilemap struct {
	width      int // 横コマ数
	height     int // 縦コマ数
	tileWidth  int // 1タイルの幅px
	tileHeight int // 1タイルの高さpx
	data       []int
	tileset    *Tileset
	solids     map[int]bool
	visible    bool
	opacity    float64
}

// NewStatic はサイズとタイルセットを指定して空の StaticTilemap を生成します。
func NewStatic(width, height, tileWidth, tileHeight int, tileset *Tileset) *StaticTilemap {
	return &StaticTilemap{
		width:      width,
		height:     height,
		tileWidth:  tileWidth,
		tileHeight: tileHeight,
		data:       make([]int, width*height),
		tileset:    tileset,
		solids:     make(map[int]bool),
		visible:    true,
		opacity:    1.0,
	}
}

// NewStaticFromData は 2 次元配列 ([][]int) から直接 StaticTilemap を構築します。
func NewStaticFromData(data [][]int, tileset *Tileset, tileDimensions ...int) *StaticTilemap {
	if len(data) == 0 || len(data[0]) == 0 {
		return NewStatic(0, 0, 16, 16, tileset)
	}

	rows := len(data)
	cols := len(data[0])

	tw := 16
	th := 16
	if tileset != nil {
		if tileset.TileWidth() > 0 {
			tw = tileset.TileWidth()
		}
		if tileset.TileHeight() > 0 {
			th = tileset.TileHeight()
		}
	}
	if len(tileDimensions) > 0 {
		tw = tileDimensions[0]
	}
	if len(tileDimensions) > 1 {
		th = tileDimensions[1]
	}

	m := NewStatic(cols, rows, tw, th, tileset)
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			if c < len(data[r]) {
				m.SetTile(c, r, data[r][c])
			}
		}
	}
	return m
}

// Update は静的レイヤーのダミー更新処理です (MapLayer インターフェースを満たします)。
func (m *StaticTilemap) Update(dt float64) {}

// --- 基本情報 ---

func (m *StaticTilemap) Width() int       { return m.width }

func (m *StaticTilemap) Height() int      { return m.height }
func (m *StaticTilemap) TileWidth() int   { return m.tileWidth }
func (m *StaticTilemap) TileHeight() int  { return m.tileHeight }
func (m *StaticTilemap) Visible() bool    { return m.visible }
func (m *StaticTilemap) SetVisible(v bool){ m.visible = v }

// --- タイル取得・設置 ---

func (m *StaticTilemap) SetTile(tx, ty, tileID int) {
	if tx < 0 || tx >= m.width || ty < 0 || ty >= m.height {
		return
	}
	m.data[ty*m.width+tx] = tileID
}

func (m *StaticTilemap) GetTileAtTile(tx, ty int) int {
	if tx < 0 || tx >= m.width || ty < 0 || ty >= m.height {
		return -1
	}
	return m.data[ty*m.width+tx]
}

func (m *StaticTilemap) GetTileAt(worldX, worldY float64) int {
	tx, ty := m.WorldToTile(worldX, worldY)
	return m.GetTileAtTile(tx, ty)
}

func (m *StaticTilemap) Fill(tileID int) {
	for i := range m.data {
		m.data[i] = tileID
	}
}

// --- 座標相互変換 ---

func (m *StaticTilemap) WorldToTile(wx, wy float64) (int, int) {
	if m.tileWidth <= 0 || m.tileHeight <= 0 {
		return 0, 0
	}
	tx := int(math.Floor(wx / float64(m.tileWidth)))
	ty := int(math.Floor(wy / float64(m.tileHeight)))
	return tx, ty
}

func (m *StaticTilemap) TileToWorld(tx, ty int) (float64, float64) {
	wx := float64(tx * m.tileWidth)
	wy := float64(ty * m.tileHeight)
	return wx, wy
}

// --- 通行判定 (Solid) ---

func (m *StaticTilemap) SetTileSolid(tileID int, solid bool) {
	m.solids[tileID] = solid
}

func (m *StaticTilemap) IsSolidAtTile(tx, ty int) bool {
	tileID := m.GetTileAtTile(tx, ty)
	if tileID < 0 {
		return true // マップ外は原則通行不可
	}
	return m.solids[tileID]
}

func (m *StaticTilemap) IsSolidAtPixel(wx, wy float64) bool {
	tx, ty := m.WorldToTile(wx, wy)
	return m.IsSolidAtTile(tx, ty)
}

// --- 描画 (全画面 / 視域カリング DrawRegion) ---

func (m *StaticTilemap) Draw(screen *ebiten.Image) {
	m.DrawRegion(screen, 0, 0, float64(m.width*m.tileWidth), float64(m.height*m.tileHeight))
}

func (m *StaticTilemap) DrawRegion(screen *ebiten.Image, viewX, viewY, viewW, viewH float64) {
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
			if tileID <= 0 { // 0 以下は透明・空タイル
				continue
			}

			img := m.tileset.TileImage(tileID)
			if img != nil {
				opts.GeoM.Reset()
				opts.GeoM.Translate(float64(c*m.tileWidth), float64(r*m.tileHeight))
				screen.DrawImage(img, opts)
			}
		}
	}
}

// GetArea はピクセル座標での範囲指定クエリ構造体 (TileArea) を返します。
func (m *StaticTilemap) GetArea(wx, wy, ww, wh float64) *TileArea {
	tx, ty := m.WorldToTile(wx, wy)
	tw := int(math.Ceil(ww / float64(m.tileWidth)))
	th := int(math.Ceil(wh / float64(m.tileHeight)))
	return m.GetTileArea(tx, ty, tw, th)
}

// GetTileArea はタイル座標での範囲指定クエリ構造体 (TileArea) を返します。
func (m *StaticTilemap) GetTileArea(tx, ty, tw, th int) *TileArea {
	return newTileArea(m, tx, ty, tw, th)
}
