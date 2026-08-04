package tilemap

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

func TestStaticTilemapAndData(t *testing.T) {
	tilesetImg := ebiten.NewImage(32, 32)
	tileset := NewTileset(tilesetImg, 16, 16)

	mapData := [][]int{
		{1, 1, 1},
		{1, 0, 1},
		{1, 1, 1},
	}

	m := NewStaticFromData(mapData, tileset)

	if m.Width() != 3 || m.Height() != 3 {
		t.Errorf("expected 3x3 map, got %dx%d", m.Width(), m.Height())
	}
	if m.GetTileAtTile(1, 1) != 0 {
		t.Errorf("expected tile at (1,1) to be 0, got %d", m.GetTileAtTile(1, 1))
	}
	if m.GetTileAtTile(0, 0) != 1 {
		t.Errorf("expected tile at (0,0) to be 1, got %d", m.GetTileAtTile(0, 0))
	}

	// 座標変換テスト
	tx, ty := m.WorldToTile(20.0, 20.0) // 16px tile -> (1, 1)
	if tx != 1 || ty != 1 {
		t.Errorf("expected WorldToTile (20,20) -> (1,1), got (%d,%d)", tx, ty)
	}

	wx, wy := m.TileToWorld(1, 1)
	if wx != 16.0 || wy != 16.0 {
		t.Errorf("expected TileToWorld (1,1) -> (16,16), got (%.1f,%.1f)", wx, wy)
	}
}

func TestSolidAndAreaQuery(t *testing.T) {
	tilesetImg := ebiten.NewImage(32, 32)
	tileset := NewTileset(tilesetImg, 16, 16)

	mapData := [][]int{
		{1, 1, 1, 1},
		{1, 0, 9, 1},
		{1, 1, 1, 1},
	}

	m := NewStaticFromData(mapData, tileset)
	m.SetTileSolid(1, true)

	// 通行判定テスト
	if !m.IsSolidAtTile(0, 0) {
		t.Errorf("expected (0,0) to be solid")
	}
	if m.IsSolidAtTile(1, 1) {
		t.Errorf("expected (1,1) not to be solid")
	}

	// Area クエリテスト
	area := m.GetArea(0, 0, 64, 48)
	foundChests := area.FindTiles(9)
	if len(foundChests) != 1 || foundChests[0].X != 2 || foundChests[0].Y != 1 {
		t.Errorf("expected chest at (2,1), got %v", foundChests)
	}

	// エリア置換テスト
	replacedCount := area.ReplaceTile(9, 0)
	if replacedCount != 1 {
		t.Errorf("expected 1 tile replaced, got %d", replacedCount)
	}
	if m.GetTileAtTile(2, 1) != 0 {
		t.Errorf("expected tile at (2,1) to be replaced to 0, got %d", m.GetTileAtTile(2, 1))
	}
}

func TestCollisionBoxCalculation(t *testing.T) {
	tilesetImg := ebiten.NewImage(32, 32)
	tileset := NewTileset(tilesetImg, 16, 16)

	mapData := [][]int{
		{1, 1, 1, 1},
		{1, 0, 0, 1},
		{1, 1, 1, 1},
	}

	mg := NewMapGroup()
	st := NewStaticFromData(mapData, tileset)
	mg.AddLayer(st)
	mg.SetTileSolid(1, true)

	boxes := mg.CreateCollisionBoxes()
	if len(boxes) == 0 {
		t.Fatalf("expected collision boxes, got empty slice")
	}

	// 壁の全周が正しく AABB 結合矩形に変換されているか検証
	totalArea := 0.0
	for _, b := range boxes {
		totalArea += b.Width * b.Height
	}

	// 壁タイル数: 10コマ * (16x16) = 2560 px^2
	expectedArea := 10.0 * 16.0 * 16.0
	if totalArea != expectedArea {
		t.Errorf("expected total collision box area %.1f, got %.1f", expectedArea, totalArea)
	}
}

func TestTiledJSONImport(t *testing.T) {
	jsonContent := []byte(`{
		"width": 3,
		"height": 3,
		"tilewidth": 16,
		"tileheight": 16,
		"layers": [
			{
				"name": "Ground",
				"type": "tilelayer",
				"visible": true,
				"opacity": 1.0,
				"width": 3,
				"height": 3,
				"data": [1, 1, 1, 1, 2, 1, 1, 1, 1]
			}
		],
		"tilesets": [
			{
				"firstgid": 1,
				"image": "tileset.png",
				"tilewidth": 16,
				"tileheight": 16,
				"tiles": [
					{
						"id": 0,
						"properties": [
							{ "name": "solid", "type": "bool", "value": true }
						]
					}
				]
			}
		]
	}`)

	img := ebiten.NewImage(32, 32)
	images := map[string]*ebiten.Image{"tileset.png": img}

	mg, err := ImportTiledJSON(jsonContent, images)
	if err != nil {
		t.Fatalf("failed to import Tiled JSON: %v", err)
	}

	if !mg.IsSolidAtPixel(8.0, 8.0) {
		t.Errorf("expected pixel (8,8) to be solid by Tiled property")
	}
}
