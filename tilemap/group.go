package tilemap

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// CollisionBox は physics エンジン等へ流し込むための物理衝突矩形構造体です。
type CollisionBox struct {
	X      float64 // 矩形中心 X (px)
	Y      float64 // 矩形中心 Y (px)
	Width  float64 // 矩形幅 (px)
	Height float64 // 矩形高さ (px)
}

// MapLayer は MapGroup 内で管理される個別のレイヤーインターフェースです。
type MapLayer interface {
	Update(dt float64)
	DrawRegion(screen *ebiten.Image, viewX, viewY, viewW, viewH float64)
	GetTileAt(worldX, worldY float64) int
	IsSolidAtPixel(worldX, worldY float64) bool
	Width() int
	Height() int
	TileWidth() int
	TileHeight() int
}

// MapGroup は複数の StaticTilemap / AnimatedTilemap レイヤーを保持し、
// 一括描画・一括更新・一括判定・Tiled JSON の受け皿となる統合マネージャーです。
type MapGroup struct {
	layers     []MapLayer
	solids     map[int]bool
	tileWidth  int
	tileHeight int
}

// NewMapGroup は MapGroup を作成します。
func NewMapGroup() *MapGroup {
	return &MapGroup{
		layers: make([]MapLayer, 0),
		solids: make(map[int]bool),
	}
}

// AddLayer はレイヤーを追加します。
func (mg *MapGroup) AddLayer(layer MapLayer) {
	if layer != nil {
		mg.layers = append(mg.layers, layer)
		if mg.tileWidth <= 0 {
			mg.tileWidth = layer.TileWidth()
			mg.tileHeight = layer.TileHeight()
		}
	}
}

// Update は全レイヤーのアニメーション等を更新します。
func (mg *MapGroup) Update(dt float64) {
	for _, l := range mg.layers {
		l.Update(dt)
	}
}

// DrawRegion は全レイヤーを表示領域カリングで重ね合わせて描画します。
func (mg *MapGroup) DrawRegion(screen *ebiten.Image, viewX, viewY, viewW, viewH float64) {
	for _, l := range mg.layers {
		l.DrawRegion(screen, viewX, viewY, viewW, viewH)
	}
}

// Draw は全画面描画を行います。
func (mg *MapGroup) Draw(screen *ebiten.Image) {
	mg.DrawRegion(screen, 0, 0, 10000, 10000)
}

// SetTileSolid は指定したタイルIDを通行不可 (Solid) に設定します。
func (mg *MapGroup) SetTileSolid(tileID int, solid bool) {
	mg.solids[tileID] = solid
	for _, l := range mg.layers {
		if st, ok := l.(*StaticTilemap); ok {
			st.SetTileSolid(tileID, solid)
		} else if anim, ok := l.(*AnimatedTilemap); ok {
			anim.SetTileSolid(tileID, solid)
		}
	}
}

// IsSolidAtPixel は指定したピクセル座標が通行不可か判定します。
func (mg *MapGroup) IsSolidAtPixel(wx, wy float64) bool {
	for _, l := range mg.layers {
		if l.IsSolidAtPixel(wx, wy) {
			return true
		}
	}
	return false
}

// GetTileAt は最前面のレイヤーのタイルIDを取得します。
func (mg *MapGroup) GetTileAt(wx, wy float64) int {
	for i := len(mg.layers) - 1; i >= 0; i-- {
		id := mg.layers[i].GetTileAt(wx, wy)
		if id > 0 {
			return id
		}
	}
	return 0
}

func (mg *MapGroup) WorldToTile(wx, wy float64) (int, int) {
	if len(mg.layers) > 0 {
		if st, ok := mg.layers[0].(*StaticTilemap); ok {
			return st.WorldToTile(wx, wy)
		}
	}
	return int(wx / 16), int(wy / 16)
}

func (mg *MapGroup) TileToWorld(tx, ty int) (float64, float64) {
	if len(mg.layers) > 0 {
		if st, ok := mg.layers[0].(*StaticTilemap); ok {
			return st.TileToWorld(tx, ty)
		}
	}
	return float64(tx * 16), float64(ty * 16)
}

// GetArea はピクセル座標での範囲指定クエリ構造体 (TileArea) を返します。
func (mg *MapGroup) GetArea(wx, wy, ww, wh float64) *TileArea {
	if len(mg.layers) > 0 {
		if st, ok := mg.layers[0].(*StaticTilemap); ok {
			return st.GetArea(wx, wy, ww, wh)
		}
	}
	return nil
}

// CreateCollisionBoxes は solid: true の隣接壁タイル群を自動結合した物理 AABB 矩形リスト (CollisionBox) を算出します。
func (mg *MapGroup) CreateCollisionBoxes() []CollisionBox {
	boxes := make([]CollisionBox, 0)
	if len(mg.layers) == 0 {
		return boxes
	}

	layer := mg.layers[0]
	w := layer.Width()
	h := layer.Height()
	tw := float64(layer.TileWidth())
	th := float64(layer.TileHeight())

	visited := make([][]bool, h)
	for r := 0; r < h; r++ {
		visited[r] = make([]bool, w)
	}

	for r := 0; r < h; r++ {
		for c := 0; c < w; c++ {
			wx := float64(c)*tw + tw/2.0
			wy := float64(r)*th + th/2.0

			if mg.IsSolidAtPixel(wx, wy) && !visited[r][c] {
				// 結合可能な水平幅を計算
				boxW := 1
				for c+boxW < w {
					nextWX := float64(c+boxW)*tw + tw/2.0
					if mg.IsSolidAtPixel(nextWX, wy) && !visited[r][c+boxW] {
						boxW++
					} else {
						break
					}
				}

				// 結合可能な垂直高さを計算
				boxH := 1
				for r+boxH < h {
					canExpand := true
					for subC := c; subC < c+boxW; subC++ {
						subWX := float64(subC)*tw + tw/2.0
						subWY := float64(r+boxH)*th + th/2.0
						if !mg.IsSolidAtPixel(subWX, subWY) || visited[r+boxH][subC] {
							canExpand = false
							break
						}
					}
					if canExpand {
						boxH++
					} else {
						break
					}
				}

				// 訪問フラグの更新
				for subR := r; subR < r+boxH; subR++ {
					for subC := c; subC < c+boxW; subC++ {
						visited[subR][subC] = true
					}
				}

				// AABB 物理結合矩形の追加 (中心座標 X, Y, 幅, 高さ)
				centerX := float64(c)*tw + (float64(boxW)*tw)/2.0
				centerY := float64(r)*th + (float64(boxH)*th)/2.0
				boxes = append(boxes, CollisionBox{
					X:      centerX,
					Y:      centerY,
					Width:  float64(boxW) * tw,
					Height: float64(boxH) * th,
				})
			}
		}
	}

	return boxes
}
