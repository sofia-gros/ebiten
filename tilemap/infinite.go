package tilemap

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// ChunkKey は チャンクのグリッド座標を表す構造体です。
type ChunkKey struct {
	X, Y int
}

// ChunkGeneratorFunc は動的チャンク生成用のコールバック関数型です。
type ChunkGeneratorFunc func(chunkX, chunkY int, chunkSize int) [][]int

// InfiniteTilemap はカメラの移動に合わせてチャンク単位で無限にマップを動的ロード・生成する構造体です。
type InfiniteTilemap struct {
	chunkSize  int // チャンクあたりのコマ数 (例: 16x16 コマ)
	tileWidth  int
	tileHeight int
	tileset    *Tileset
	chunks     map[ChunkKey]*StaticTilemap
	generator  ChunkGeneratorFunc
	solids     map[int]bool
}

// NewInfinite は無限タイルマップを生成します。
func NewInfinite(chunkSize, tileWidth, tileHeight int, tileset *Tileset, generator ChunkGeneratorFunc) *InfiniteTilemap {
	if chunkSize <= 0 {
		chunkSize = 16
	}
	return &InfiniteTilemap{
		chunkSize:  chunkSize,
		tileWidth:  tileWidth,
		tileHeight: tileHeight,
		tileset:    tileset,
		chunks:     make(map[ChunkKey]*StaticTilemap),
		generator:  generator,
		solids:     make(map[int]bool),
	}
}

// SetTileSolid は特定のタイルIDを通行不可 (Solid) に設定します。
func (inf *InfiniteTilemap) SetTileSolid(tileID int, solid bool) {
	inf.solids[tileID] = solid
	for _, chunk := range inf.chunks {
		chunk.SetTileSolid(tileID, solid)
	}
}

// GetChunk は指定したチャンク座標の StaticTilemap を取得 (無ければ自動生成) します。
func (inf *InfiniteTilemap) GetChunk(cx, cy int) *StaticTilemap {
	key := ChunkKey{X: cx, Y: cy}
	if chunk, ok := inf.chunks[key]; ok {
		return chunk
	}

	// チャンクの新規生成
	var chunkData [][]int
	if inf.generator != nil {
		chunkData = inf.generator(cx, cy, inf.chunkSize)
	}

	var chunk *StaticTilemap
	if len(chunkData) > 0 {
		chunk = NewStaticFromData(chunkData, inf.tileset, inf.tileWidth, inf.tileHeight)
	} else {
		chunk = NewStatic(inf.chunkSize, inf.chunkSize, inf.tileWidth, inf.tileHeight, inf.tileset)
	}

	for id, s := range inf.solids {
		chunk.SetTileSolid(id, s)
	}

	inf.chunks[key] = chunk
	return chunk
}

// GetTileAt はピクセル座標直下のタイルIDを取得します。
func (inf *InfiniteTilemap) GetTileAt(wx, wy float64) int {
	chunkX, chunkY, tx, ty := inf.WorldToChunkTile(wx, wy)
	chunk := inf.GetChunk(chunkX, chunkY)
	return chunk.GetTileAtTile(tx, ty)
}

// IsSolidAtPixel はピクセル座標が通行不可か判定します。
func (inf *InfiniteTilemap) IsSolidAtPixel(wx, wy float64) bool {
	chunkX, chunkY, tx, ty := inf.WorldToChunkTile(wx, wy)
	chunk := inf.GetChunk(chunkX, chunkY)
	return chunk.IsSolidAtTile(tx, ty)
}

// WorldToChunkTile はワールド座標を (チャンクX, チャンクY, チャンク内タイルのX, チャンク内タイルのY) へ変換します。
func (inf *InfiniteTilemap) WorldToChunkTile(wx, wy float64) (cx, cy, tx, ty int) {
	chunkPxWidth := float64(inf.chunkSize * inf.tileWidth)
	chunkPxHeight := float64(inf.chunkSize * inf.tileHeight)

	cx = int(floorDiv(wx, chunkPxWidth))
	cy = int(floorDiv(wy, chunkPxHeight))

	relX := wx - float64(cx)*chunkPxWidth
	relY := wy - float64(cy)*chunkPxHeight

	tx = int(relX / float64(inf.tileWidth))
	ty = int(relY / float64(inf.tileHeight))

	return cx, cy, tx, ty
}

func floorDiv(a, b float64) float64 {
	if b == 0 {
		return 0
	}
	res := a / b
	if res < 0 && float64(int(res)) != res {
		return res - 1.0
	}
	return res
}

// DrawRegion はカメラ等の表示領域内に含まれるチャンク群を自動計算してカリング描画します。
func (inf *InfiniteTilemap) DrawRegion(screen *ebiten.Image, viewX, viewY, viewW, viewH float64) {
	if screen == nil || inf.tileset == nil {
		return
	}

	chunkPxWidth := float64(inf.chunkSize * inf.tileWidth)
	chunkPxHeight := float64(inf.chunkSize * inf.tileHeight)

	startCX := int(floorDiv(viewX, chunkPxWidth))
	endCX := int(floorDiv(viewX+viewW, chunkPxWidth))
	startCY := int(floorDiv(viewY, chunkPxHeight))
	endCY := int(floorDiv(viewY+viewH, chunkPxHeight))

	for cy := startCY; cy <= endCY; cy++ {
		for cx := startCX; cx <= endCX; cx++ {
			chunk := inf.GetChunk(cx, cy)
			offsetWX := float64(cx) * chunkPxWidth
			offsetWY := float64(cy) * chunkPxHeight

			// チャンク領域内のカリング描画
			subViewX := viewX - offsetWX
			subViewY := viewY - offsetWY

			opts := &ebiten.DrawImageOptions{}
			opts.GeoM.Translate(offsetWX, offsetWY)

			// チャンク用一時キャンバス描画の代わりに、オフセット位置を計算して直接描画
			chunk.DrawRegion(screen, subViewX, subViewY, viewW, viewH)
		}
	}
}
