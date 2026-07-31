package camera

import (
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

// Camera は 2D ゲームにおける注視点、ズーム、回転、画面振動、ビューポート切り抜き、
// カスタムシェーダー切り替えを統合管理するカメラオブジェクトです。
type Camera struct {
	name       string
	width      float64
	height     float64
	x          float64
	y          float64
	zoom       float64
	rotation   float64
	zIndex     int
	viewport   *Viewport
	bounds     *Bounds
	hasBounds  bool
	
	// 画面振動 (Shake) 用パラメータ
	shakeStrength float64
	shakeDuration float64
	shakeTimer    float64
	shakeOffsetX  float64
	shakeOffsetY  float64

	// カスタムシェーダー状態
	shader       *ebiten.Shader
	shaderOption ShaderOption
	hasShader    bool

	// レンダリング用一時オフスクリーンバッファ (シェーダー/ビューポート用)
	offscreen *ebiten.Image
}

// New は画面の幅と高さ (px) を指定して新しい Camera を作成します。
func New(width, height float64) *Camera {
	return NewWithName("", width, height)
}

// NewWithName は名前、幅、高さを指定して Camera を作成します。
func NewWithName(name string, width, height float64) *Camera {
	return &Camera{
		name:   name,
		width:  width,
		height: height,
		zoom:   1.0,
	}
}

// --- ゲッターメソッド群 ---

// Name はカメラの名前を取得します。
func (c *Camera) Name() string { return c.name }

// Pos は現在のカメラの中心ワールド座標 (x, y) を返します。
func (c *Camera) Pos() (float64, float64) { return c.x, c.y }

// X はカメラの中心 X 座標を返します (sound.PlayAt などのポジショナル位置指定に便利)。
func (c *Camera) X() float64 { return c.x }

// Y はカメラの中心 Y 座標を返します。
func (c *Camera) Y() float64 { return c.y }

// Zoom は現在のズーム倍率を返します (デフォルト 1.0)。
func (c *Camera) Zoom() float64 { return c.zoom }

// Rotation は現在の回転角度 (ラジアン) を返します。
func (c *Camera) Rotation() float64 { return c.rotation }

// ZIndex はカメラの描画優先度を返します。
func (c *Camera) ZIndex() int { return c.zIndex }

// HasShader は現在シェーダーがセットされているかを返します。
func (c *Camera) HasShader() bool { return c.hasShader }

// IsShaking は現在カメラが画面振動中であるかを返します。
func (c *Camera) IsShaking() bool { return c.shakeTimer > 0 }

// Viewport は設定されているビューポート情報 (x, y, w, h) を返します。
// 設定されていない場合は全画面 (0, 0, width, height) を返します。
func (c *Camera) Viewport() (float64, float64, float64, float64) {
	if c.viewport != nil {
		return c.viewport.X, c.viewport.Y, c.viewport.Width, c.viewport.Height
	}
	return 0, 0, c.width, c.height
}

// VisibleBounds は現在カメラ画面内に見えているワールド空間の包囲領域 Bounds (minX, minY, maxX, maxY) を返します。
// 画面外のオブジェクトの描画や計算をスキップするカリング（Culling）処理に最適です。
func (c *Camera) VisibleBounds() (float64, float64, float64, float64) {
	vw, vh := c.width, c.height
	if c.viewport != nil {
		vw, vh = c.viewport.Width, c.viewport.Height
	}

	halfW := (vw / 2.0) / c.zoom
	halfH := (vh / 2.0) / c.zoom

	// 回転時の包含矩形対角長を考慮
	if c.rotation != 0 {
		diag := math.Sqrt(halfW*halfW + halfH*halfH)
		halfW, halfH = diag, diag
	}

	return c.x - halfW, c.y - halfH, c.x + halfW, c.y + halfH
}

// --- セッター・操作メソッド群 ---

// SetPos はカメラの中心ワールド座標 (x, y) を設定します。
func (c *Camera) SetPos(x, y float64) {
	c.x = x
	c.y = y
	c.clampToBounds()
}

// Move は現在のカメラ座標から指定オフセット (dx, dy) 移動させます。
func (c *Camera) Move(dx, dy float64) {
	c.x += dx
	c.y += dy
	c.clampToBounds()
}

// MoveTo は目標座標 (targetX, targetY) に向かって指定スピード / Lerp率 (0.0 〜 1.0) で滑らかに直線補間移動させます。
func (c *Camera) MoveTo(targetX, targetY, speed float64) {
	if speed >= 1.0 {
		c.SetPos(targetX, targetY)
		return
	}
	c.x += (targetX - c.x) * speed
	c.y += (targetY - c.y) * speed
	c.clampToBounds()
}

// SetZoom はズーム倍率を設定します (1.0 が等倍、2.0 で2倍拡大、0.5 で2倍広域)。
func (c *Camera) SetZoom(zoom float64) {
	if zoom <= 0.001 {
		zoom = 0.001
	}
	c.zoom = zoom
}

// SetRotation はカメラの回転角度 (ラジアン) を設定します。
func (c *Camera) SetRotation(angle float64) {
	c.rotation = angle
}

// SetZIndex はカメラの描画優先度 (重ね合わせ順序) を設定します。
func (c *Camera) SetZIndex(zIndex int) {
	c.zIndex = zIndex
}

// SetViewport は画面内でのカメラ出力領域 (x, y, width, height) を設定します (ミニマップや画面分割用)。
func (c *Camera) SetViewport(x, y, width, height float64) {
	c.viewport = &Viewport{
		X:      x,
		Y:      y,
		Width:  width,
		Height: height,
	}
}

// ClearViewport はビューポート設定を解除し、全画面出力に戻します。
func (c *Camera) ClearViewport() {
	c.viewport = nil
}

// SetBounds はカメラが移動できるワールド境界領域 (minX, minY, maxX, maxY) を制限します。
func (c *Camera) SetBounds(minX, minY, maxX, maxY float64) {
	c.bounds = &Bounds{MinX: minX, MinY: minY, MaxX: maxX, MaxY: maxY}
	c.hasBounds = true
	c.clampToBounds()
}

// ClearBounds はカメラの移動境界制限を解除します。
func (c *Camera) ClearBounds() {
	c.bounds = nil
	c.hasBounds = false
}

// Shake は画面を指定された強さ (strength: px) と時間 (durationSec: 秒) でランダムに振動させます。
func (c *Camera) Shake(strength, durationSec float64) {
	c.shakeStrength = strength
	c.shakeDuration = durationSec
	c.shakeTimer = durationSec
}

// SetShader はカメラに Custom Shader を適用します。適用中、カメラは自動的にシェーダー描画モードに切り替わります。
func (c *Camera) SetShader(shader *ebiten.Shader, opts ...ShaderOption) {
	c.shader = shader
	c.hasShader = (shader != nil)
	if len(opts) > 0 {
		c.shaderOption = opts[0]
	} else {
		c.shaderOption = ShaderOption{}
	}
}

// ClearShader はシェーダーの設定を解除し、元のダイレクト描画モードに戻します。
func (c *Camera) ClearShader() {
	c.shader = nil
	c.hasShader = false
}

// Update は画面振動 (Shake) やタイマー進行を毎フレーム更新します。
func (c *Camera) Update(dt float64) {
	if c.shakeTimer > 0 {
		c.shakeTimer -= dt
		if c.shakeTimer <= 0 {
			c.shakeTimer = 0
			c.shakeOffsetX = 0
			c.shakeOffsetY = 0
		} else {
			// 指数減衰しながらランダム振動
			factor := c.shakeTimer / c.shakeDuration
			currStrength := c.shakeStrength * factor
			c.shakeOffsetX = (rand.Float64()*2.0 - 1.0) * currStrength
			c.shakeOffsetY = (rand.Float64()*2.0 - 1.0) * currStrength
		}
	}
}

// clampToBounds はカメラ座標が bounds の範囲を超えないよう自動修正します。
func (c *Camera) clampToBounds() {
	if !c.hasBounds || c.bounds == nil {
		return
	}

	vw, vh := c.width, c.height
	if c.viewport != nil {
		vw, vh = c.viewport.Width, c.viewport.Height
	}

	halfW := (vw / 2.0) / c.zoom
	halfH := (vh / 2.0) / c.zoom

	if c.x-halfW < c.bounds.MinX {
		c.x = c.bounds.MinX + halfW
	}
	if c.x+halfW > c.bounds.MaxX {
		c.x = c.bounds.MaxX - halfW
	}
	if c.y-halfH < c.bounds.MinY {
		c.y = c.bounds.MinY + halfH
	}
	if c.y+halfH > c.bounds.MaxY {
		c.y = c.bounds.MaxY - halfH
	}
}

// --- 座標相互変換 (ScreenToWorld / WorldToScreen) ---

// ScreenToWorld は画面上のローカル座標 (screenX, screenY) を、カメラ座標系を反映したワールド空間座標 (worldX, worldY) に変換します。
func (c *Camera) ScreenToWorld(screenX, screenY float64) (float64, float64) {
	vw, vh := c.width, c.height
	vx, vy := 0.0, 0.0
	if c.viewport != nil {
		vx, vy = c.viewport.X, c.viewport.Y
		vw, vh = c.viewport.Width, c.viewport.Height
	}

	// 1. ビューポート原点からの相対位置
	relX := screenX - vx - (vw / 2.0)
	relY := screenY - vy - (vh / 2.0)

	// 2. ズーム解除
	relX /= c.zoom
	relY /= c.zoom

	// 3. 回転解除
	if c.rotation != 0 {
		cosA := math.Cos(-c.rotation)
		sinA := math.Sin(-c.rotation)
		rx := relX*cosA - relY*sinA
		ry := relX*sinA + relY*cosA
		relX, relY = rx, ry
	}

	// 4. カメラ位置とシェイクオフセットを加算
	camX := c.x + c.shakeOffsetX
	camY := c.y + c.shakeOffsetY

	return relX + camX, relY + camY
}

// WorldToScreen はワールド空間座標 (worldX, worldY) を画面上の表示位置座標 (screenX, screenY) に変換します。
func (c *Camera) WorldToScreen(worldX, worldY float64) (float64, float64) {
	vw, vh := c.width, c.height
	vx, vy := 0.0, 0.0
	if c.viewport != nil {
		vx, vy = c.viewport.X, c.viewport.Y
		vw, vh = c.viewport.Width, c.viewport.Height
	}

	camX := c.x + c.shakeOffsetX
	camY := c.y + c.shakeOffsetY

	// 1. カメラ座標からの相対オフセット
	relX := worldX - camX
	relY := worldY - camY

	// 2. 回転適用
	if c.rotation != 0 {
		cosA := math.Cos(c.rotation)
		sinA := math.Sin(c.rotation)
		rx := relX*cosA - relY*sinA
		ry := relX*sinA + relY*cosA
		relX, relY = rx, ry
	}

	// 3. ズーム適用
	relX *= c.zoom
	relY *= c.zoom

	// 4. ビューポートスクリーン座標への移動
	return relX + vx + (vw / 2.0), relY + vy + (vh / 2.0)
}

// Apply は指定された ebiten.DrawImageOptions の GeoM に対して、このカメラのトランスフォーム行列を適用します。
func (c *Camera) Apply(opts *ebiten.DrawImageOptions) {
	if opts == nil {
		return
	}

	vw, vh := c.width, c.height
	vx, vy := 0.0, 0.0
	if c.viewport != nil {
		vx, vy = c.viewport.X, c.viewport.Y
		vw, vh = c.viewport.Width, c.viewport.Height
	}

	camX := c.x + c.shakeOffsetX
	camY := c.y + c.shakeOffsetY

	// カメラ平行移動・回転・ズームを行列合成
	opts.GeoM.Translate(-camX, -camY)
	if c.rotation != 0 {
		opts.GeoM.Rotate(c.rotation)
	}
	opts.GeoM.Scale(c.zoom, c.zoom)
	opts.GeoM.Translate(vx+(vw/2.0), vy+(vh/2.0))
}

// --- レンダリング (Render) ---

// Render は描画クロージャ (drawFunc) を実行し、カメラの座標・ビューポート・シェーダー・シェイクを適用して screen へ描画します。
func (c *Camera) Render(screen *ebiten.Image, drawFunc func(target *ebiten.Image)) {
	if screen == nil || drawFunc == nil {
		return
	}

	// 出力画面および表示可能領域の決定
	vw, vh := int(c.width), int(c.height)
	vx, vy := 0, 0

	if c.viewport != nil {
		vx, vy = int(c.viewport.X), int(c.viewport.Y)
		vw, vh = int(c.viewport.Width), int(c.viewport.Height)
	}

	// ビューポート範囲への安全なクリッピング
	targetScreen := screen
	if c.viewport != nil {
		targetScreen = screen.SubImage(rectToRectangle(vx, vy, vw, vh)).(*ebiten.Image)
	}

	// シェーダーが有効な場合
	if c.hasShader && c.shader != nil {
		if c.offscreen == nil || c.offscreen.Bounds().Dx() != vw || c.offscreen.Bounds().Dy() != vh {
			c.offscreen = ebiten.NewImage(vw, vh)
		}
		c.offscreen.Clear()

		// オフスクリーンに描画
		drawFunc(c.offscreen)

		// シェーダーを描画適用
		sOpts := &ebiten.DrawRectShaderOptions{}
		sOpts.GeoM.Translate(float64(vx), float64(vy))
		if c.shaderOption.Uniforms != nil {
			sOpts.Uniforms = c.shaderOption.Uniforms
		}
		sOpts.Images = c.shaderOption.Images
		sOpts.Images[0] = c.offscreen

		targetScreen.DrawRectShader(vw, vh, c.shader, sOpts)
		return
	}

	// ダイレクト描画モード (高速・中間テクスチャ生成ゼロ)
	drawFunc(targetScreen)
}
