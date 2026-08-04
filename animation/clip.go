package animation

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// SpriteProvider は Frame(index) メソッドでコマ画像を取得できる任意の型を表すインターフェースです。
// (asset.SpriteSheet なども自動的に適合します)
type SpriteProvider interface {
	Frame(index int) *ebiten.Image
}

// FrameCallback は特定コマへの到達やアニメーション完了時に呼ばれる関数型です。
type FrameCallback func()

// ClipOptions は Clip の初期化オプション構造体です。
type ClipOptions struct {
	Frames   []int     // コマインデックスの配列 (例: [0, 1, 2, 3])
	FPS      float64   // アニメーション再生速度 (1秒あたりのコマ数)
	Loop     bool      // ループ再生するか
	Reverse  bool      // 逆再生するか
	PingPong bool      // 往復再生するか (0->1->2->1->0)
	Durations []float64 // コマごとの個別の表示時間 (秒)。指定がある場合は FPS より優先されます
}

// Clip はアニメーションのコマシーケンス、速度、ループフラグ、イベントを保持するデータ定義です。
type Clip struct {
	name            string
	images          []*ebiten.Image
	durations       []float64
	fps             float64
	loop            bool
	reverse         bool
	pingPong        bool
	onComplete      FrameCallback
	onFrameHandlers map[int]FrameCallback
}

// NewClip は名前、画像リスト (または SpriteProvider)、オプションから Clip を生成します。
func NewClip(name string, source any, opts ...ClipOptions) *Clip {
	c := &Clip{
		name:            name,
		images:          make([]*ebiten.Image, 0),
		fps:             8.0,
		loop:            true,
		reverse:         false,
		pingPong:        false,
		onFrameHandlers: make(map[int]FrameCallback),
	}

	opt := ClipOptions{FPS: 8.0, Loop: true}
	if len(opts) > 0 {
		userOpt := opts[0]
		if userOpt.FPS > 0 {
			opt.FPS = userOpt.FPS
		}
		opt.Frames = userOpt.Frames
		opt.Loop = userOpt.Loop
		opt.Reverse = userOpt.Reverse
		opt.PingPong = userOpt.PingPong
		opt.Durations = userOpt.Durations
	}

	c.fps = opt.FPS
	c.loop = opt.Loop
	c.reverse = opt.Reverse
	c.pingPong = opt.PingPong

	// 画像ソースの抽出
	switch s := source.(type) {
	case []*ebiten.Image:
		if len(opt.Frames) > 0 {
			for _, idx := range opt.Frames {
				if idx >= 0 && idx < len(s) {
					c.images = append(c.images, s[idx])
				}
			}
		} else {
			c.images = append(c.images, s...)
		}
	case SpriteProvider:
		if len(opt.Frames) > 0 {
			for _, idx := range opt.Frames {
				if img := s.Frame(idx); img != nil {
					c.images = append(c.images, img)
				}
			}
		}
	case *ebiten.Image:
		c.images = append(c.images, s)
	}

	// コマ表示時間の計算
	numFrames := len(c.images)
	if numFrames > 0 {
		c.durations = make([]float64, numFrames)
		defaultDur := 1.0 / c.fps
		for i := 0; i < numFrames; i++ {
			if i < len(opt.Durations) && opt.Durations[i] > 0 {
				c.durations[i] = opt.Durations[i]
			} else {
				c.durations[i] = defaultDur
			}
		}
	}

	return c
}

// --- ゲッター・セッター & 動的フラグ操作 ---

func (c *Clip) Name() string                     { return c.name }
func (c *Clip) FrameCount() int                  { return len(c.images) }
func (c *Clip) FrameImage(idx int) *ebiten.Image {
	if idx < 0 || idx >= len(c.images) {
		return nil
	}
	return c.images[idx]
}

func (c *Clip) Duration(idx int) float64 {
	if idx < 0 || idx >= len(c.durations) {
		return 1.0 / c.fps
	}
	return c.durations[idx]
}

func (c *Clip) FPS() float64                     { return c.fps }
func (c *Clip) SetFPS(fps float64) {
	if fps <= 0 {
		return
	}
	c.fps = fps
	defaultDur := 1.0 / fps
	for i := range c.durations {
		c.durations[i] = defaultDur
	}
}

func (c *Clip) IsLoop() bool                     { return c.loop }
func (c *Clip) SetLoop(loop bool)                { c.loop = loop }

func (c *Clip) IsReverse() bool                  { return c.reverse }
func (c *Clip) SetReverse(reverse bool)          { c.reverse = reverse }

func (c *Clip) IsPingPong() bool                 { return c.pingPong }
func (c *Clip) SetPingPong(pingPong bool)        { c.pingPong = pingPong }

// OnComplete はアニメーションが全コマ再生完了した時に発火するコールバックを登録します。
func (c *Clip) OnComplete(fn FrameCallback) *Clip {
	c.onComplete = fn
	return c
}

// OnFrame は指定コマ index に到達した時に発火するコールバックを登録します。
func (c *Clip) OnFrame(frameIndex int, fn FrameCallback) *Clip {
	c.onFrameHandlers[frameIndex] = fn
	return c
}
