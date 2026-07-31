package ui

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

// NineSlice は画像を 9 つの領域に分割し、角を歪ませずに綺麗に拡大縮小描画する描画構造体です。
type NineSlice struct {
	image        *ebiten.Image
	borderTop    int
	borderRight  int
	borderBottom int
	borderLeft   int
}

// NewNineSlice はテクスチャ画像と 4 方向の境界幅 (px) を指定して NineSlice を作成します。
func NewNineSlice(img *ebiten.Image, top, right, bottom, left int) *NineSlice {
	return &NineSlice{
		image:        img,
		borderTop:    top,
		borderRight:  right,
		borderBottom: bottom,
		borderLeft:   left,
	}
}

// Draw は指定された描画対象 (target) の領域 (x, y, width, height) に 9スライス伸長描画を行います。
func (n *NineSlice) Draw(target *ebiten.Image, destX, destY, destW, destH float64) {
	if n == nil || n.image == nil || target == nil || destW <= 0 || destH <= 0 {
		return
	}

	srcW, srcH := n.image.Bounds().Dx(), n.image.Bounds().Dy()

	// スライス用 9 パッチのソース座標・ディスティネーション座標を計算
	srcX := []int{0, n.borderLeft, srcW - n.borderRight, srcW}
	srcY := []int{0, n.borderTop, srcH - n.borderBottom, srcH}

	dstX := []float64{destX, destX + float64(n.borderLeft), destX + destW - float64(n.borderRight), destX + destW}
	dstY := []float64{destY, destY + float64(n.borderTop), destY + destH - float64(n.borderBottom), destY + destH}

	for row := 0; row < 3; row++ {
		for col := 0; col < 3; col++ {
			sRect := image.Rect(srcX[col], srcY[row], srcX[col+1], srcY[row+1])
			if sRect.Dx() <= 0 || sRect.Dy() <= 0 {
				continue
			}

			dw := dstX[col+1] - dstX[col]
			dh := dstY[row+1] - dstY[row]
			if dw <= 0 || dh <= 0 {
				continue
			}

			subImg := n.image.SubImage(sRect).(*ebiten.Image)
			opts := &ebiten.DrawImageOptions{}

			scaleX := dw / float64(sRect.Dx())
			scaleY := dh / float64(sRect.Dy())

			opts.GeoM.Scale(scaleX, scaleY)
			opts.GeoM.Translate(dstX[col], dstY[row])

			target.DrawImage(subImg, opts)
		}
	}
}
