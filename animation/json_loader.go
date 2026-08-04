package animation

import (
	"encoding/json"
	"fmt"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

// Aseprite / TexturePacker 汎用 JSON の構造体定義 (Hash / Array 形式の両対応)

type asepriteRect struct {
	X int `json:"x"`
	Y int `json:"y"`
	W int `json:"w"`
	H int `json:"h"`
}

type asepriteFrameItem struct {
	Filename string       `json:"filename"`
	Frame    asepriteRect `json:"frame"`
	Duration int          `json:"duration"` // ms
}

type asepriteTag struct {
	Name      string `json:"name"`
	From      int    `json:"from"`
	To        int    `json:"to"`
	Direction string `json:"direction"` // forward, reverse, pingpong
}

type asepriteMeta struct {
	FrameTags []asepriteTag `json:"frameTags"`
}

type asepriteDataRaw struct {
	Frames json.RawMessage `json:"frames"`
	Meta   asepriteMeta    `json:"meta"`
}

// LoadClipsFromJSON はベース画像 (*ebiten.Image) と Aseprite/TexturePacker 等の標準 JSON バイト列から
// 含まれるアニメーション Clip 群を全自動解析・解析作成します。
func LoadClipsFromJSON(baseImg *ebiten.Image, jsonBytes []byte) ([]*Clip, error) {
	if baseImg == nil || len(jsonBytes) == 0 {
		return nil, fmt.Errorf("invalid base image or empty JSON data")
	}

	var raw asepriteDataRaw
	if err := json.Unmarshal(jsonBytes, &raw); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	// 1. コマフレーム情報の抽出 (Hash 形式または Array 形式)
	var framesList []asepriteFrameItem

	// Array 形式の解析試行
	var arrFrames []asepriteFrameItem
	if err := json.Unmarshal(raw.Frames, &arrFrames); err == nil && len(arrFrames) > 0 {
		framesList = arrFrames
	} else {
		// Hash 形式の解析試行
		var mapFrames map[string]asepriteFrameItem
		if err := json.Unmarshal(raw.Frames, &mapFrames); err == nil && len(mapFrames) > 0 {
			for name, item := range mapFrames {
				item.Filename = name
				framesList = append(framesList, item)
			}
		}
	}

	if len(framesList) == 0 {
		return nil, fmt.Errorf("no frames found in JSON data")
	}

	// 2. コマ画像 (*ebiten.Image) および表示時間の切り出し
	frameImages := make([]*ebiten.Image, len(framesList))
	frameDurations := make([]float64, len(framesList))

	for i, f := range framesList {
		rect := image.Rect(f.Frame.X, f.Frame.Y, f.Frame.X+f.Frame.W, f.Frame.Y+f.Frame.H)
		frameImages[i] = baseImg.SubImage(rect).(*ebiten.Image)
		if f.Duration > 0 {
			frameDurations[i] = float64(f.Duration) / 1000.0 // ms -> seconds
		} else {
			frameDurations[i] = 1.0 / 8.0
		}
	}

	clips := make([]*Clip, 0)

	// 3. meta.frameTags から Clip を自動構築
	if len(raw.Meta.FrameTags) > 0 {
		for _, tag := range raw.Meta.FrameTags {
			from := tag.From
			to := tag.To
			if from < 0 {
				from = 0
			}
			if to >= len(frameImages) {
				to = len(frameImages) - 1
			}

			tagImages := make([]*ebiten.Image, 0)
			tagDurs := make([]float64, 0)

			for idx := from; idx <= to; idx++ {
				tagImages = append(tagImages, frameImages[idx])
				tagDurs = append(tagDurs, frameDurations[idx])
			}

			isRev := (tag.Direction == "reverse")
			isPingPong := (tag.Direction == "pingpong")

			clip := NewClip(tag.Name, tagImages, ClipOptions{
				Loop:      true,
				Reverse:   isRev,
				PingPong:  isPingPong,
				Durations: tagDurs,
			})
			clips = append(clips, clip)
		}
	} else {
		// タグが無い場合は全コマを "default" クリップとして作成
		clip := NewClip("default", frameImages, ClipOptions{
			Loop:      true,
			Durations: frameDurations,
		})
		clips = append(clips, clip)
	}

	return clips, nil
}

// CreateAnimatorFromJSON はベース画像と JSON データから Animator を自動生成し、Manager に登録します。
func (m *Manager) CreateAnimatorFromJSON(baseImg *ebiten.Image, jsonBytes []byte, opts ...AnimatorOptions) (*Animator, error) {
	clips, err := LoadClipsFromJSON(baseImg, jsonBytes)
	if err != nil {
		return nil, err
	}
	if len(clips) == 0 {
		return nil, fmt.Errorf("no clips generated from JSON")
	}

	animator := NewAnimator(clips[0], opts...)
	for i := 1; i < len(clips); i++ {
		animator.AddClip(clips[i])
	}

	m.Add(animator)
	return animator, nil
}
