package scene

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type sceneNode struct {
	scene           Scene
	isUpdateEnabled bool
	isDrawEnabled   bool
}

type command interface {
	execute(m *Manager)
}

// Manager は Ebitengine の Game インターフェースを実装する、シーンの統括管理者です。
type Manager struct {
	nodes        []*sceneNode
	commands     []command
	context      *Context
	screenWidth  int
	screenHeight int
}

// NewManager は新しいシーンマネージャーを作成します。
func NewManager(width, height int) *Manager {
	m := &Manager{
		nodes:        make([]*sceneNode, 0),
		commands:     make([]command, 0),
		screenWidth:  width,
		screenHeight: height,
	}
	m.context = &Context{manager: m}
	return m
}

// Context はこのマネージャーに関連付けられたコンテキストを取得します。
// アプリケーション起動時の初期シーンの設定などに使用します。
func (m *Manager) Context() *Context {
	return m.context
}

// Layout implements ebiten.Game.
func (m *Manager) Layout(outsideWidth, outsideHeight int) (int, int) {
	return m.screenWidth, m.screenHeight
}

// Update implements ebiten.Game.
func (m *Manager) Update() error {
	// シーン遷移や状態変更などの遅延コマンドを評価
	m.flushCommands()

	// Update可能なシーンを更新 (下から順に実行)
	for _, n := range m.nodes {
		if n.isUpdateEnabled {
			if err := n.scene.Update(m.context); err != nil {
				return err
			}
		}
	}
	return nil
}

// Draw implements ebiten.Game.
func (m *Manager) Draw(screen *ebiten.Image) {
	// Draw可能なシーンを下から順に描画（後から追加したシーンが上になる）
	for _, n := range m.nodes {
		if n.isDrawEnabled {
			n.scene.Draw(screen)
		}
	}
}

// flushCommands は溜まったコマンドを実行し、スライスをクリアします。
func (m *Manager) flushCommands() {
	for _, cmd := range m.commands {
		cmd.execute(m)
	}
	// 容量を維持したままスライスをリセット
	m.commands = m.commands[:0]
}
