package tween

// Group は特定のカテゴリ (UIアニメーション群、エネミー演出群等) ごとに Tween をまとめて一括管理・制御するグループ構造体です。
type Group struct {
	manager *Manager
}

// NewGroup は新しい Tween Group を作成します。
func NewGroup() *Group {
	return &Group{
		manager: NewManager(),
	}
}

// New はこの Group に属する待機状態の Tween を作成します。
func (g *Group) New(opts ...Option) *Tween {
	return g.manager.New(opts...)
}

// FromTo はこの Group に属する Tween を作成します。
func (g *Group) FromTo(start, end, duration float64) *Tween {
	return g.manager.FromTo(start, end, duration)
}

// Add はこの Group に Tween を追加登録します。
func (g *Group) Add(tw *Tween) {
	g.manager.Add(tw)
}

// Remove はこの Group から Tween を除外します。
func (g *Group) Remove(tw *Tween) {
	g.manager.Remove(tw)
}

// Update はこの Group 内のすべての Tween の時間を一括進行させます。
func (g *Group) Update(dt float64) {
	g.manager.Update(dt)
}

// PauseAll はこの Group 内のすべての Tween を一括一時停止します。
func (g *Group) PauseAll() {
	g.manager.PauseAll()
}

// ResumeAll はこの Group 内のすべての Tween の一時停止を一括解除します。
func (g *Group) ResumeAll() {
	g.manager.ResumeAll()
}

// RestartAll はこの Group 内のすべての Tween を最初から一括リスタートします。
func (g *Group) RestartAll() {
	g.manager.RestartAll()
}

// Clear はこの Group 内のすべての Tween を消去・強制停止します。
func (g *Group) Clear() {
	g.manager.Clear()
}

// Count はこの Group で管理されている active な Tween の総数を返します。
func (g *Group) Count() int {
	return g.manager.Count()
}
