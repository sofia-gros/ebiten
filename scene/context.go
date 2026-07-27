package scene

import "reflect"

// Context はシーンからマネージャーを操作するためのAPIを提供します。
type Context struct {
	manager *Manager
}

// Start は現在のすべてのシーンを破棄し、新しいシーンを起動します。
func (c *Context) Start(s Scene) {
	c.manager.commands = append(c.manager.commands, &cmdStart{scene: s})
}

// Overlay は現在のシーンの上に新しいシーンを重ねて起動します。
func (c *Context) Overlay(s Scene) {
	c.manager.commands = append(c.manager.commands, &cmdOverlay{scene: s})
}

// Hide は指定したシーンの Update と Draw を停止します。
func (c *Context) Hide(s Scene) {
	c.manager.commands = append(c.manager.commands, &cmdChangeState{target: s, action: "hide"})
}

// Show は指定したシーンの Update と Draw を再開し、最前面に移動します。
func (c *Context) Show(s Scene) {
	c.manager.commands = append(c.manager.commands, &cmdChangeState{target: s, action: "show"})
}

// Stop は指定したシーンの Update を停止します（Draw は継続されます）。
func (c *Context) Stop(s Scene) {
	c.manager.commands = append(c.manager.commands, &cmdChangeState{target: s, action: "stop"})
}

// Run は指定したシーンの Update を再開します。
func (c *Context) Run(s Scene) {
	c.manager.commands = append(c.manager.commands, &cmdChangeState{target: s, action: "run"})
}

// Destroy は指定したシーンを破棄（スタックから削除）します。
func (c *Context) Destroy(s Scene) {
	c.manager.commands = append(c.manager.commands, &cmdDestroy{target: s})
}

// --- 以下、内部コマンドの実装 ---

func findNodeByType(nodes []*sceneNode, target Scene) int {
	if target == nil {
		return -1
	}
	targetType := reflect.TypeOf(target)
	// 末尾（最上位）から検索する
	for i := len(nodes) - 1; i >= 0; i-- {
		if reflect.TypeOf(nodes[i].scene) == targetType {
			return i
		}
	}
	return -1
}

type cmdStart struct {
	scene Scene
}

func (c *cmdStart) execute(m *Manager) {
	m.nodes = []*sceneNode{{
		scene:           c.scene,
		isUpdateEnabled: true,
		isDrawEnabled:   true,
	}}
}

type cmdOverlay struct {
	scene Scene
}

func (c *cmdOverlay) execute(m *Manager) {
	m.nodes = append(m.nodes, &sceneNode{
		scene:           c.scene,
		isUpdateEnabled: true,
		isDrawEnabled:   true,
	})
}

type cmdChangeState struct {
	target Scene
	action string // "hide", "show", "stop", "run"
}

func (c *cmdChangeState) execute(m *Manager) {
	idx := findNodeByType(m.nodes, c.target)
	if idx == -1 {
		return
	}
	node := m.nodes[idx]

	switch c.action {
	case "hide":
		node.isUpdateEnabled = false
		node.isDrawEnabled = false
	case "show":
		node.isUpdateEnabled = true
		node.isDrawEnabled = true
		// 最前面でない場合は最前面（末尾）に移動する
		if idx != len(m.nodes)-1 {
			m.nodes = append(m.nodes[:idx], m.nodes[idx+1:]...)
			m.nodes = append(m.nodes, node)
		}
	case "stop":
		node.isUpdateEnabled = false
	case "run":
		node.isUpdateEnabled = true
	}
}

type cmdDestroy struct {
	target Scene
}

func (c *cmdDestroy) execute(m *Manager) {
	idx := findNodeByType(m.nodes, c.target)
	if idx == -1 {
		return
	}
	// 指定インデックスの要素を削除
	m.nodes = append(m.nodes[:idx], m.nodes[idx+1:]...)
}
