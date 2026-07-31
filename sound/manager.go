package sound

import (
	"bytes"
	"fmt"
	"io"
	"sync"
	"sync/atomic"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

// Manager は BGM, SE, Voice などの各種サウンドを統括管理するメインマネージャーです。
type Manager struct {
	mu             sync.RWMutex
	audioContext   *audio.Context
	sampleRate     int
	masterVolume   float64
	isMasterMuted  bool
	typeVolumes    map[Type]float64
	typeMutes      map[Type]bool
	activeTracks   []*Track
	nextTrackID    atomic.Uint64
}

// NewManager は指定されたサンプルレート (例: 44100) で新しい Sound Manager を作成します。
func NewManager(sampleRate int) *Manager {
	ctx := audio.CurrentContext()
	if ctx == nil {
		ctx = audio.NewContext(sampleRate)
	}
	m := &Manager{
		audioContext:  ctx,
		sampleRate:    sampleRate,
		masterVolume:  1.0,
		typeVolumes:   make(map[Type]float64),
		typeMutes:     make(map[Type]bool),
		activeTracks:  make([]*Track, 0),
	}

	// デフォルトタイプの初期音量設定
	m.typeVolumes[TypeSE] = 1.0
	m.typeVolumes[TypeBGM] = 1.0
	m.typeVolumes[TypeVoice] = 1.0
	m.typeVolumes[TypeEnv] = 1.0

	return m
}

// SampleRate は現在のオーディオコンテキストのサンプルレートを取得します。
func (m *Manager) SampleRate() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.sampleRate
}

// SetSampleRate はサンプルレートを変更します。
func (m *Manager) SetSampleRate(sampleRate int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.sampleRate = sampleRate
	ctx := audio.CurrentContext()
	if ctx == nil {
		ctx = audio.NewContext(sampleRate)
	}
	m.audioContext = ctx
}


// --- 音量・ミュートコントロール ---

// SetMasterVolume はマスター全体の音量 (0.0 〜 1.0) を設定します。
func (m *Manager) SetMasterVolume(vol float64) {
	m.mu.Lock()
	m.masterVolume = vol
	tracks := append([]*Track(nil), m.activeTracks...)
	m.mu.Unlock()

	for _, t := range tracks {
		t.updatePlayerVolume()
	}
}

// MasterVolume は現在のマスター音量を取得します。
func (m *Manager) MasterVolume() float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.masterVolume
}

// SetVolume は指定したサウンドタイプ (TypeSE, TypeBGM 等) の音量を設定します。
func (m *Manager) SetVolume(soundType Type, vol float64) {
	m.mu.Lock()
	m.typeVolumes[soundType] = vol
	tracks := append([]*Track(nil), m.activeTracks...)
	m.mu.Unlock()

	for _, t := range tracks {
		if t.soundType == soundType {
			t.updatePlayerVolume()
		}
	}
}

// Volume は指定したサウンドタイプの現在の音量を取得します。
func (m *Manager) Volume(soundType Type) float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.volume(soundType)
}

func (m *Manager) volume(soundType Type) float64 {
	if vol, ok := m.typeVolumes[soundType]; ok {
		return vol
	}
	return 1.0
}

// SetMute は指定したサウンドタイプのミュート状態を設定します。
func (m *Manager) SetMute(soundType Type, mute bool) {
	m.mu.Lock()
	m.typeMutes[soundType] = mute
	tracks := append([]*Track(nil), m.activeTracks...)
	m.mu.Unlock()

	for _, t := range tracks {
		if t.soundType == soundType {
			t.updatePlayerVolume()
		}
	}
}

// IsMuted は指定したサウンドタイプがミュートされているかを取得します。
func (m *Manager) IsMuted(soundType Type) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.isMuted(soundType)
}

func (m *Manager) isMuted(soundType Type) bool {
	if mute, ok := m.typeMutes[soundType]; ok {
		return mute
	}
	return false
}

// --- 再生 (Play / PlayAt / CrossFade) ---

// Play は指定した音声データ (バイトスライス、`io.Reader` 等) を再生し、操作用の *Track ハンドルを返します。
func (m *Manager) Play(src any, opts ...Option) *Track {
	opt := defaultOption(opts)

	r, err := m.toReader(src)
	if err != nil {
		return nil
	}

	var player *audio.Player
	if opt.Loop {
		// ループプレイヤーの生成
		type lengthGetter interface {
			Length() int64
		}
		type sizeGetter interface {
			Size() int64
		}

		var length int64
		if lg, ok := r.(lengthGetter); ok {
			length = lg.Length()
		} else if sg, ok := r.(sizeGetter); ok {
			length = sg.Size()
		}

		if seeker, ok := r.(io.ReadSeeker); ok && length > 0 {
			loopStream := audio.NewInfiniteLoop(seeker, length)
			player, err = m.audioContext.NewPlayer(loopStream)
		} else {
			player, err = m.audioContext.NewPlayer(r)
		}
	} else {
		player, err = m.audioContext.NewPlayer(r)
	}


	if err != nil || player == nil {
		return nil
	}

	trackID := m.nextTrackID.Add(1)
	track := &Track{
		id:             trackID,
		soundType:      opt.Type,
		player:         player,
		manager:        m,
		volume:         opt.Volume,
		pan:            opt.Pan,
		pitch:          opt.Pitch,
		loop:           opt.Loop,
		fadeInDuration: opt.FadeInDuration,
	}

	track.updatePlayerVolume()
	player.Play()

	m.mu.Lock()
	m.activeTracks = append(m.activeTracks, track)
	m.mu.Unlock()

	return track
}

// PlayAt は位置情報を伴ってサウンドを再生し、距離減衰と左右パンニングを自動計算して適用します。
func (m *Manager) PlayAt(src any, sourceX, sourceY, listenerX, listenerY, maxDistance float64, opts ...Option) *Track {
	track := m.Play(src, opts...)
	if track != nil {
		track.SetPosition(sourceX, sourceY, listenerX, listenerY, maxDistance)
	}
	return track
}

// CrossFade は現在再生されている指定タイプ（デフォルトは TypeBGM）のサウンドをフェードアウトさせながら、新しい曲をフェードイン再生します。
func (m *Manager) CrossFade(src any, fadeSec float64, opts ...Option) *Track {
	opt := defaultOption(opts)

	// 既存の同タイプサウンドをフェードアウト停止
	m.StopType(opt.Type, fadeSec)

	// 新しいサウンドをフェードイン再生
	opt.FadeInDuration = fadeSec
	return m.Play(src, opt)
}

// --- 停止 (Stop / StopType / StopAll) ---

// Stop は指定したトラックをフェードアウト秒数（省略時は即時）で停止します。
func (m *Manager) Stop(track *Track, fadeSec ...float64) {
	if track != nil {
		track.Stop(fadeSec...)
	}
}

// StopType は指定したサウンドタイプ (TypeSE, TypeBGM 等) の再生中トラックを一括して停止します。
func (m *Manager) StopType(soundType Type, fadeSec ...float64) {
	m.mu.RLock()
	tracks := append([]*Track(nil), m.activeTracks...)
	m.mu.RUnlock()

	for _, t := range tracks {
		if t.soundType == soundType {
			t.Stop(fadeSec...)
		}
	}
}

// StopAll はすべての再生中トラックを一括して停止します。
func (m *Manager) StopAll(fadeSec ...float64) {
	m.mu.RLock()
	tracks := append([]*Track(nil), m.activeTracks...)
	m.mu.RUnlock()

	for _, t := range tracks {
		t.Stop(fadeSec...)
	}
}

// --- 一時停止 / 再開 (PauseAll / ResumeAll) ---

// PauseAll はすべての再生中トラックを一時停止します (ポーズ画面用)。
func (m *Manager) PauseAll() {
	m.mu.RLock()
	tracks := append([]*Track(nil), m.activeTracks...)
	m.mu.RUnlock()

	for _, t := range tracks {
		t.Pause()
	}
}

// ResumeAll は一時停止中のすべてのトラックの再生を再開します。
func (m *Manager) ResumeAll() {
	m.mu.RLock()
	tracks := append([]*Track(nil), m.activeTracks...)
	m.mu.RUnlock()

	for _, t := range tracks {
		t.Resume()
	}
}

// Update は毎フレームのフェードイン/アウト処理および終了済みトラックの自動クリーンアップを行います。
func (m *Manager) Update(dt float64) {
	m.mu.Lock()
	defer m.mu.Unlock()

	n := 0
	for _, t := range m.activeTracks {
		finished := t.update(dt)
		if !finished {
			m.activeTracks[n] = t
			n++
		}
	}
	m.activeTracks = m.activeTracks[:n]
}

// Clear はすべての再生を即時停止し、トラックプールを全クリアします。
func (m *Manager) Clear() {
	m.StopAll()
	m.mu.Lock()
	defer m.mu.Unlock()
	m.activeTracks = m.activeTracks[:0]
}

// toReader は入力ソース ([]byte, io.Reader) を WAV/PCM ストリームへ安全に変換します。
func (m *Manager) toReader(src any) (io.Reader, error) {
	switch v := src.(type) {
	case []byte:
		stream, err := wav.DecodeWithSampleRate(m.sampleRate, bytes.NewReader(v))
		if err == nil {
			return stream, nil
		}
		return bytes.NewReader(v), nil
	case io.Reader:
		return v, nil
	default:
		return nil, fmt.Errorf("unsupported sound source type: %T", src)
	}
}
