package asset

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

// SpriteSheet は元の画像から指定サイズでコマ割り（グリッド切り出し）されたスプライト情報を保持します。
type SpriteSheet struct {
	baseImage   *ebiten.Image
	frameWidth  int
	frameHeight int
	margin      int
	spacing     int
	frames      []*ebiten.Image
}

// NewSpriteSheet はベース画像から指定のサイズ・間隔でコマ割りした SpriteSheet を生成します。
func NewSpriteSheet(base *ebiten.Image, frameWidth, frameHeight, margin, spacing int) *SpriteSheet {
	ss := &SpriteSheet{
		baseImage:   base,
		frameWidth:  frameWidth,
		frameHeight: frameHeight,
		margin:      margin,
		spacing:     spacing,
		frames:      make([]*ebiten.Image, 0),
	}

	if base == nil || frameWidth <= 0 || frameHeight <= 0 {
		if base != nil {
			ss.frames = append(ss.frames, base)
		}
		return ss
	}

	bounds := base.Bounds()
	imgWidth := bounds.Dx()
	imgHeight := bounds.Dy()

	// マージンとスペーシングを考慮して縦横のコマ数を計算
	for y := margin; y+frameHeight <= imgHeight-margin; y += frameHeight + spacing {
		for x := margin; x+frameWidth <= imgWidth-margin; x += frameWidth + spacing {
			subImg := base.SubImage(image.Rect(x, y, x+frameWidth, y+frameHeight)).(*ebiten.Image)
			ss.frames = append(ss.frames, subImg)
		}
	}

	// 1コマも切り出せなかった場合はベース画像を全体の1コマとして扱う
	if len(ss.frames) == 0 {
		ss.frames = append(ss.frames, base)
	}

	return ss
}

// Frame は指定したインデックス (0始まり) のコマ画像 (*ebiten.Image) を返します。
// 範囲外のインデックスが指定された場合は、最後のコマまたはダミーコマを安全に返します。
func (ss *SpriteSheet) Frame(index int) *ebiten.Image {
	if ss == nil || len(ss.frames) == 0 {
		return nil
	}
	if index < 0 {
		return ss.frames[0]
	}
	if index >= len(ss.frames) {
		return ss.frames[len(ss.frames)-1]
	}
	return ss.frames[index]
}

// Frames は切り出されたすべてのコマ画像 (*ebiten.Image) のスライスを返します。
func (ss *SpriteSheet) Frames() []*ebiten.Image {
	if ss == nil {
		return nil
	}
	return ss.frames
}

// Count は切り出されたコマの総数を返します。
func (ss *SpriteSheet) Count() int {
	if ss == nil {
		return 0
	}
	return len(ss.frames)
}

