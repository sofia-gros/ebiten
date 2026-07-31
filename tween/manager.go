package tween

import (
	"sync"
)

// Manager は複数の Tween の時間更新・ライフサイクルを一括管理するマネージャー構造体です。
type Manager struct {
	mu     sync.RWMutex
	tweens []*Tween
}

// DefaultManager は全自動登録用のデフォルトグローバルマネージャーです。
var DefaultManager = NewManager()

// NewManager は新しいマネージャーを作成します。
func NewManager() *Manager {
	return &Manager{
		tweens: make([]*Tween, 0),
	}
}

// New はこのマネージャーに紐づく待機状態の Tween を作成します。
func (m *Manager) New(opts ...Option) *Tween {
	tw := New(opts...)
	tw.manager = m
	return tw
}

// FromTo はこのマネージャーに紐づく Tween を作成します。
func (m *Manager) FromTo(start, end, duration float64) *Tween {
	tw := FromTo(start, end, duration)
	tw.manager = m
	return tw
}

// Add はマネージャーに Tween を追加登録します。
func (m *Manager) Add(tw *Tween) {
	if tw == nil {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()

	tw.manager = m
	for _, existing := range m.tweens {
		if existing == tw {
			return
		}
	}
	m.tweens = append(m.tweens, tw)
}

// Remove はマネージャーから Tween を除外します。
func (m *Manager) Remove(tw *Tween) {
	if tw == nil {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()

	newTweens := make([]*Tween, 0, len(m.tweens))
	for _, existing := range m.tweens {
		if existing != tw {
			newTweens = append(newTweens, existing)
		}
	}
	m.tweens = newTweens
}

// Update は登録されているすべての Tween の時間を一括進行させ、完了したものを自動消去します。
func (m *Manager) Update(dt float64) {
	m.mu.Lock()
	activeTweens := make([]*Tween, len(m.tweens))
	copy(activeTweens, m.tweens)
	m.mu.Unlock()

	var remaining []*Tween
	for _, tw := range activeTweens {
		finished := tw.updateInternal(dt)
		if !finished {
			remaining = append(remaining, tw)
		}
	}

	m.mu.Lock()
	m.tweens = remaining
	m.mu.Unlock()
}

// PauseAll はマネージャー内のすべての Tween を一括一時停止します。
func (m *Manager) PauseAll() {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, tw := range m.tweens {
		tw.Pause()
	}
}

// ResumeAll はマネージャー内のすべての Tween の一括一時停止を解除します。
func (m *Manager) ResumeAll() {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, tw := range m.tweens {
		tw.Resume()
	}
}

// RestartAll はマネージャー内のすべての Tween を最初から一括リスタートします。
func (m *Manager) RestartAll() {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, tw := range m.tweens {
		tw.Restart()
	}
}

// Clear はマネージャー内のすべての Tween を消去・強制停止します。
func (m *Manager) Clear() {
	m.mu.Lock()
	currentTweens := make([]*Tween, len(m.tweens))
	copy(currentTweens, m.tweens)
	m.tweens = make([]*Tween, 0)
	m.mu.Unlock()

	for _, tw := range currentTweens {
		tw.stopInternal()
	}
}


// Count は現在マネージャーで管理されている Tween の総数を返します。
func (m *Manager) Count() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.tweens)
}

// --- パッケージレベルの全自動ヘルパー関数 ---

// Update は DefaultManager に登録されている全 Tween の時間を一括更新します。
func Update(dt float64) {
	DefaultManager.Update(dt)
}

// PauseAll は DefaultManager の全 Tween を一括一時停止します。
func PauseAll() {
	DefaultManager.PauseAll()
}

// ResumeAll は DefaultManager の全 Tween の一括一時停止を解除します。
func ResumeAll() {
	DefaultManager.ResumeAll()
}

// RestartAll は DefaultManager の全 Tween を一括リスタートします。
func RestartAll() {
	DefaultManager.RestartAll()
}

// Clear は DefaultManager の全 Tween をクリアします。
func Clear() {
	DefaultManager.Clear()
}
