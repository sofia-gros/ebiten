package tilemap

import (
	"encoding/json"
	"fmt"
	"path"

	"github.com/hajimehoshi/ebiten/v2"
)

// Tiled JSON フォーマットのデータ型定義

type tiledProp struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value any    `json:"value"`
}

type tiledFrame struct {
	Duration int `json:"duration"` // ms
	TileID   int `json:"tileid"`
}

type tiledTileDef struct {
	ID        int          `json:"id"`
	Animation []tiledFrame `json:"animation"`
	Properties []tiledProp  `json:"properties"`
}

type tiledTilesetRef struct {
	FirstGID int            `json:"firstgid"`
	Source   string         `json:"source"`
	Image    string         `json:"image"`
	TileW    int            `json:"tilewidth"`
	TileH    int            `json:"tileheight"`
	Margin   int            `json:"margin"`
	Spacing  int            `json:"spacing"`
	Tiles    []tiledTileDef `json:"tiles"`
}

type tiledLayerRaw struct {
	Name    string `json:"name"`
	Type    string `json:"type"` // tilelayer
	Visible bool   `json:"visible"`
	Opacity float64 `json:"opacity"`
	Width   int    `json:"width"`
	Height  int    `json:"height"`
	Data    []int  `json:"data"`
}

type tiledMapRaw struct {
	Width      int               `json:"width"`
	Height     int               `json:"height"`
	TileWidth  int               `json:"tilewidth"`
	TileHeight int               `json:"tileheight"`
	Layers     []tiledLayerRaw   `json:"layers"`
	Tilesets   []tiledTilesetRef `json:"tilesets"`
}

// ImportTiledJSON は Tiled が出力した JSON バイト列と画像マップから MapGroup を全自動生成します。
func ImportTiledJSON(jsonBytes []byte, images map[string]*ebiten.Image) (*MapGroup, error) {
	if len(jsonBytes) == 0 {
		return nil, fmt.Errorf("empty Tiled JSON data")
	}

	var raw tiledMapRaw
	if err := json.Unmarshal(jsonBytes, &raw); err != nil {
		return nil, fmt.Errorf("failed to unmarshal Tiled JSON: %w", err)
	}

	mg := NewMapGroup()
	mg.tileWidth = raw.TileWidth
	mg.tileHeight = raw.TileHeight

	// 1. タイルセット画像の検索と作成
	var baseTileset *Tileset
	var firstGID = 1

	if len(raw.Tilesets) > 0 {
		tsRef := raw.Tilesets[0]
		firstGID = tsRef.FirstGID
		if firstGID <= 0 {
			firstGID = 1
		}

		imgKey := path.Base(tsRef.Image)
		var img *ebiten.Image
		if i, ok := images[imgKey]; ok {
			img = i
		} else if i, ok := images[tsRef.Image]; ok {
			img = i
		} else {
			for _, i := range images {
				img = i
				break
			}
		}

		if img != nil {
			tw := tsRef.TileW
			if tw <= 0 { tw = raw.TileWidth }
			th := tsRef.TileH
			if th <= 0 { th = raw.TileHeight }
			baseTileset = NewTileset(img, tw, th, tsRef.Margin, tsRef.Spacing)
		}
	}

	// 2. レイヤーの生成と格納
	for _, lRaw := range raw.Layers {
		if lRaw.Type != "tilelayer" || len(lRaw.Data) == 0 {
			continue
		}

		// アニメーションタイルの定義があるか判定
		hasAnimation := false
		if len(raw.Tilesets) > 0 {
			for _, tDef := range raw.Tilesets[0].Tiles {
				if len(tDef.Animation) > 0 {
					hasAnimation = true
					break
				}
			}
		}

		if hasAnimation {
			animMap := NewAnimated(lRaw.Width, lRaw.Height, raw.TileWidth, raw.TileHeight, baseTileset)
			animMap.SetVisible(lRaw.Visible)

			// データの流し込み
			for r := 0; r < lRaw.Height; r++ {
				for c := 0; c < lRaw.Width; c++ {
					gid := lRaw.Data[r*lRaw.Width+c]
					if gid >= firstGID {
						animMap.SetTile(c, r, gid-firstGID+1)
					}
				}
			}

			// アニメーションタイルの登録
			for _, tDef := range raw.Tilesets[0].Tiles {
				if len(tDef.Animation) > 0 {
					frames := make([]int, len(tDef.Animation))
					for i, f := range tDef.Animation {
						frames[i] = f.TileID + firstGID
					}
					fps := 4.0
					if len(tDef.Animation) > 0 && tDef.Animation[0].Duration > 0 {
						fps = 1000.0 / float64(tDef.Animation[0].Duration)
					}
					animMap.SetAnimatedTile(tDef.ID+firstGID, frames, fps)
				}
			}

			mg.AddLayer(animMap)
		} else {
			staticMap := NewStatic(lRaw.Width, lRaw.Height, raw.TileWidth, raw.TileHeight, baseTileset)
			staticMap.SetVisible(lRaw.Visible)

			for r := 0; r < lRaw.Height; r++ {
				for c := 0; c < lRaw.Width; c++ {
					gid := lRaw.Data[r*lRaw.Width+c]
					if gid >= firstGID {
						staticMap.SetTile(c, r, gid-firstGID+1)
					}
				}
			}
			mg.AddLayer(staticMap)
		}
	}

	// 3. solid: true プロパティの自動解析
	if len(raw.Tilesets) > 0 {
		for _, tDef := range raw.Tilesets[0].Tiles {
			for _, prop := range tDef.Properties {
				if prop.Name == "solid" {
					if b, ok := prop.Value.(bool); ok && b {
						mg.SetTileSolid(tDef.ID+firstGID, true)
					}
				}
			}
		}
	}

	return mg, nil
}
