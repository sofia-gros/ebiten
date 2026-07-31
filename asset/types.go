package asset

// Tilemap は Tiled Map Editor 等で作成されたマップデータを表す組み込み型です。
type Tilemap struct {
	Width      int         `json:"width" yaml:"width"`           // マップの横タイル数
	Height     int         `json:"height" yaml:"height"`         // マップの縦タイル数
	TileWidth  int         `json:"tilewidth" yaml:"tilewidth"`   // タイル1枚の幅(px)
	TileHeight int         `json:"tileheight" yaml:"tileheight"` // タイル1枚の高さ(px)
	Layers     []TileLayer `json:"layers" yaml:"layers"`         // マップレイヤー一覧
	Tilesets   []Tileset   `json:"tilesets" yaml:"tilesets"`     // タイルセット情報一覧
}

// TileLayer はマップ内の個々のレイヤー情報を表します。
type TileLayer struct {
	Name    string    `json:"name" yaml:"name"`       // レイヤー名
	Width   int       `json:"width" yaml:"width"`     // レイヤーの横タイル数
	Height  int       `json:"height" yaml:"height"`   // レイヤーの縦タイル数
	Data    []int     `json:"data" yaml:"data"`       // タイルIDの配置配列
	Visible bool      `json:"visible" yaml:"visible"` // 表示フラグ
	Opacity float64   `json:"opacity" yaml:"opacity"` // 透明度 (0.0 〜 1.0)
}

// Tileset はマップで使用されるタイルセットの情報を表します。
type Tileset struct {
	FirstGID   int    `json:"firstgid" yaml:"firstgid"`     // 開始GID
	Name       string `json:"name" yaml:"name"`             // タイルセット名
	Image      string `json:"image" yaml:"image"`           // 画像ファイルパス
	ImageWidth int    `json:"imagewidth" yaml:"imagewidth"` // 画像全体の横幅(px)
	ImageHeight int   `json:"imageheight" yaml:"imageheight"`// 画像全体の縦幅(px)
	TileWidth  int    `json:"tilewidth" yaml:"tilewidth"`   // 1タイルの横幅(px)
	TileHeight int    `json:"tileheight" yaml:"tileheight"` // 1タイルの縦幅(px)
}

// AnimationData はスプライトアニメーションの定義情報を表す組み込み型です。
type AnimationData struct {
	FrameRate float64            `json:"framerate" yaml:"framerate"` // デフォルトフレームレート (FPS)
	Loop      bool               `json:"loop" yaml:"loop"`           // デフォルトループ再生フラグ
	Tags      map[string]AnimTag `json:"tags" yaml:"tags"`           // タグ名ごとのアニメーション定義 ("walk", "attack" 等)
}

// AnimTag は特定のアニメーション動作（タグ）のフレーム構成や時間を表します。
type AnimTag struct {
	Name      string    `json:"name" yaml:"name"`           // タグ名
	Frames    []int     `json:"frames" yaml:"frames"`       // 使用するコマ(フレーム)インデックスの配列 [0, 1, 2, 1]
	Durations []float64 `json:"durations" yaml:"durations"` // 各コマの表示時間(秒) (省略時は FrameRate から自動算出)
	Loop      bool      `json:"loop" yaml:"loop"`           // このタグ独自のループフラグ
}
