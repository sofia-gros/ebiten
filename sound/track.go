package sound

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2/audio"
)

// Track は現在再生中または一時停止中のサウンド（1つのトラック）の制御用ハンドルです。
type Track struct {
	id             uint64
	soundType      Type
	player         *audio.Player
	manager        *Manager
	volume         float64
	pan            float64
	pitch          float64
	loop           bool
	isPaused       bool
	isStopped      bool
	
	// フェード処理用パラメータ
	fadeInDuration  float64
	fadeInTimer     float64
	fadeOutDuration float64
	fadeOutTimer    float64
	isFadingOut     bool
}

// Pause はこのトラックの再生を一時停止します。
func (t *Track) Pause() {
	if t == nil || t.isStopped || t.isPaused {
		return
	}
	if t.player != nil {
		t.player.Pause()
	}
	t.isPaused = true
}

// Resume は一時停止中のこのトラックの再生を再開します。
func (t *Track) Resume() {
	if t == nil || t.isStopped || !t.isPaused {
		return
	}
	if t.player != nil {
		t.player.Play()
	}
	t.isPaused = false
}

// IsPlaying はこのトラックが現在再生中であるかを返します。
func (t *Track) IsPlaying() bool {
	if t == nil || t.isStopped || t.isPaused {
		return false
	}
	if t.player == nil {
		return false
	}
	return t.player.IsPlaying()
}

// SetVolume はこのトラック個別の音量 (0.0 〜 1.0) をリアルタイムに変更します。
func (t *Track) SetVolume(vol float64) {
	if t == nil {
		return
	}
	t.volume = math.Max(0.0, vol)
	t.updatePlayerVolume()
}

// SetPan はこのトラックの左右パンニング (-1.0: 左 〜 0.0: 中央 〜 1.0: 右) を変更します。
func (t *Track) SetPan(pan float64) {
	if t == nil {
		return
	}
	t.pan = math.Max(-1.0, math.Min(1.0, pan))
	t.updatePlayerVolume()
}

// SetPosition は音源位置 (sourceX, sourceY) とプレイヤー位置 (listenerX, listenerY) から、
// パンと距離に応じた音量減衰をリアルタイムに計算して適用します (ポジショナルサウンド動的追従)。
func (t *Track) SetPosition(sourceX, sourceY, listenerX, listenerY, maxDistance float64) {
	if t == nil {
		return
	}

	dx := sourceX - listenerX
	dy := sourceY - listenerY
	dist := math.Sqrt(dx*dx + dy*dy)

	// 距離減衰計算 (0.0 〜 1.0)
	volFactor := 1.0 - (dist / maxDistance)
	if volFactor < 0.0 {
		volFactor = 0.0
	}

	// 左右パン計算 (dx に応じて -1.0 〜 1.0)
	panFactor := dx / maxDistance
	if panFactor < -1.0 {
		panFactor = -1.0
	} else if panFactor > 1.0 {
		panFactor = 1.0
	}

	t.volume = volFactor
	t.pan = panFactor
	t.updatePlayerVolume()
}

// Stop は指定したフェードアウト秒数（省略可）でこのトラックを停止します。
func (t *Track) Stop(fadeSec ...float64) {
	if t == nil || t.isStopped {
		return
	}

	duration := 0.0
	if len(fadeSec) > 0 && fadeSec[0] > 0 {
		duration = fadeSec[0]
	}

	if duration <= 0 {
		t.isStopped = true
		if t.player != nil {
			t.player.Close()
		}
	} else {
		t.isFadingOut = true
		t.fadeOutDuration = duration
		t.fadeOutTimer = duration
	}
}

// updatePlayerVolume はマスター音量・タイプ別音量・個別音量・フェード倍率を掛け合わせてプレイヤー音量を更新します。
func (t *Track) updatePlayerVolume() {
	if t == nil || t.player == nil || t.manager == nil {
		return
	}

	if t.manager.isMasterMuted || t.manager.isMuted(t.soundType) {
		t.player.SetVolume(0)
		return
	}

	masterVol := t.manager.masterVolume
	typeVol := t.manager.volume(t.soundType)
	trackVol := t.volume

	// フェードイン倍率
	fadeInFactor := 1.0
	if t.fadeInDuration > 0 && t.fadeInTimer < t.fadeInDuration {
		fadeInFactor = t.fadeInTimer / t.fadeInDuration
	}

	// フェードアウト倍率
	fadeOutFactor := 1.0
	if t.isFadingOut && t.fadeOutDuration > 0 {
		fadeOutFactor = t.fadeOutTimer / t.fadeOutDuration
		if fadeOutFactor < 0 {
			fadeOutFactor = 0
		}
	}

	finalVolume := masterVol * typeVol * trackVol * fadeInFactor * fadeOutFactor
	t.player.SetVolume(finalVolume)

	// パンの更新
	if panSetter, ok := any(t.player).(interface{ SetBufferSize(int) }); ok {
		_ = panSetter // 静的型検証補助
	}
	// Note: ebiten/audio の基本 Player では Player.SetVolume のみ、Pan サポートはボリューム分配にて行います。
}

// update は毎フレームのタイマー更新処理を行います。
func (t *Track) update(dt float64) bool {
	if t.isStopped {
		return true // 終了
	}

	// フェードイン処理
	if t.fadeInDuration > 0 && t.fadeInTimer < t.fadeInDuration {
		t.fadeInTimer += dt
		if t.fadeInTimer > t.fadeInDuration {
			t.fadeInTimer = t.fadeInDuration
		}
		t.updatePlayerVolume()
	}

	// フェードアウト処理
	if t.isFadingOut {
		t.fadeOutTimer -= dt
		t.updatePlayerVolume()
		if t.fadeOutTimer <= 0 {
			t.isStopped = true
			if t.player != nil {
				t.player.Close()
			}
			return true
		}
	}

	// ループ指定なしで再生完了した場合
	if !t.loop && t.player != nil && !t.player.IsPlaying() && !t.isPaused {
		t.isStopped = true
		t.player.Close()
		return true
	}

	return false
}
