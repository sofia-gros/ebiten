package animation

import (
	"fmt"
	"sort"
	"sync"
)

// Manager はゲーム内の複数の Animator を一括保持し、同期更新・一括ポーズ・グローバルタイムスケールを統括します。
type Manager struct {
	mu        sync.RWMutex
	animators map[string]*Animator
	list      []*Animator
	globalSpeed float64
	paused    bool
}

// NewManager は新しい Manager インスタンスを作成します。
func NewManager() *Manager {
	return &Manager{
		animators:   make(map[string]*Animator),
		list:        make([]*Animator, 0),
		globalSpeed: 1.0,
	}
}

// CreateAnimator は初期 Clip を指定して Animator を作成し、Manager の管理下に登録します。
func (m *Manager) CreateAnimator(defaultClip *Clip, opts ...AnimatorOptions) *Animator {
	animator := NewAnimator(defaultClip, opts...)
	m.Add(animator)
	return animator
}

// Add は既存の Animator を Manager に登録します。
func (m *Manager) Add(animator *Animator) {
	if animator == nil {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()

	key := fmt.Sprintf("%p", animator)
	m.animators[key] = animator
	m.rebuildList()
}

// Remove は Animator を Manager の一括管理対象から削除します。
func (m *Manager) Remove(animator *Animator) {
	if animator == nil {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()

	key := fmt.Sprintf("%p", animator)
	delete(m.animators, key)
	m.rebuildList()
}

// Update は登録された全 Animator を一括更新します。グローバルタイムスケールや一括ポーズも適用されます。
func (m *Manager) Update(dt float64) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.paused {
		return
	}

	effectiveDt := dt * m.globalSpeed
	for _, anim := range m.list {
		anim.Update(effectiveDt)
	}
}

// PauseAll は管理下の全 Animator を一括で一時停止します。
func (m *Manager) PauseAll() {
	m.mu.RLock()
	defer m.mu.RUnlock()
	m.paused = true
}

// ResumeAll は管理下の全 Animator の一時停止を一括解除します。
func (m *Manager) ResumeAll() {
	m.mu.RLock()
	defer m.mu.RUnlock()
	m.paused = false
}

// SetSpeed は全 Animator に適用されるグローバルタイムスケール (倍速/スローモーション) を設定します。
func (m *Manager) SetSpeed(speed float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if speed >= 0 {
		m.globalSpeed = speed
	}
}

func (m *Manager) Speed() float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.globalSpeed
}

func (m *Manager) rebuildList() {
	m.list = make([]*Animator, 0, len(m.animators))
	for _, a := range m.animators {
		m.list = append(m.list, a)
	}
	sort.SliceStable(m.list, func(i, j int) bool {
		return fmt.Sprintf("%p", m.list[i]) < fmt.Sprintf("%p", m.list[j])
	})
}
