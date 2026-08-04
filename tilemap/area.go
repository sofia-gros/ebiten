package tilemap

// TilePos はタイルグリッド上の位置 (X, Y) を表します。
type TilePos struct {
	X int
	Y int
}

// TileArea は GetArea / GetTileArea によって指定された部分領域に対するクエリ操作構造体です。
type TileArea struct {
	tilemap TilemapInterface
	startX  int
	startY  int
	width   int
	height  int
}

// TilemapInterface は TileArea がクエリを実行するために使用する共通インターフェースです。
type TilemapInterface interface {
	GetTileAtTile(tx, ty int) int
	SetTile(tx, ty, tileID int)
	Width() int
	Height() int
}

func newTileArea(m TilemapInterface, startX, startY, width, height int) *TileArea {
	return &TileArea{
		tilemap: m,
		startX:  startX,
		startY:  startY,
		width:   width,
		height:  height,
	}
}

// FindTiles は指定したエリア内から特定の tileID を持つ全座標スライスを検索して返します。
func (ta *TileArea) FindTiles(tileID int) []TilePos {
	result := make([]TilePos, 0)
	if ta == nil || ta.tilemap == nil {
		return result
	}

	for r := ta.startY; r < ta.startY+ta.height; r++ {
		for c := ta.startX; c < ta.startX+ta.width; c++ {
			if ta.tilemap.GetTileAtTile(c, r) == tileID {
				result = append(result, TilePos{X: c, Y: r})
			}
		}
	}
	return result
}

// ReplaceTile は指定したエリア内にある targetTileID を新しい newTileID に置換します。
func (ta *TileArea) ReplaceTile(targetTileID, newTileID int) int {
	if ta == nil || ta.tilemap == nil {
		return 0
	}

	count := 0
	for r := ta.startY; r < ta.startY+ta.height; r++ {
		for c := ta.startX; c < ta.startX+ta.width; c++ {
			if ta.tilemap.GetTileAtTile(c, r) == targetTileID {
				ta.tilemap.SetTile(c, r, newTileID)
				count++
			}
		}
	}
	return count
}

// Tiles は指定したエリア内の 2 次元タイルIDスライス ([][]int) を取得します。
func (ta *TileArea) Tiles() [][]int {
	if ta == nil || ta.tilemap == nil || ta.width <= 0 || ta.height <= 0 {
		return nil
	}

	res := make([][]int, ta.height)
	for r := 0; r < ta.height; r++ {
		res[r] = make([]int, ta.width)
		for c := 0; c < ta.width; c++ {
			res[r][c] = ta.tilemap.GetTileAtTile(ta.startX+c, ta.startY+r)
		}
	}
	return res
}
