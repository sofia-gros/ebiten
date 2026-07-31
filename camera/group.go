package camera

import (
	"fmt"
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
)

// Group は複数の Camera を一括して管理し、ZIndex (優先度) 順に自動ソートしてレンダリングするマネージャー構造体です。
type Group struct {
	cameras map[string]*Camera
	list    []*Camera
}

// NewGroup は新しい Camera Group を作成します。
func NewGroup(cameras ...*Camera) *Group {
	g := &Group{
		cameras: make(map[string]*Camera),
		list:    make([]*Camera, 0),
	}
	for _, c := range cameras {
		g.Add(c)
	}
	return g
}

// Add はグループに Camera を追加します。
func (g *Group) Add(cam *Camera) {
	if cam == nil {
		return
	}
	key := cam.Name()
	if key == "" {
		key = fmt.Sprintf("%p", cam)
	}
	g.cameras[key] = cam
	g.rebuildList()
}

// Get は名前でグループ内の Camera を取得します。
func (g *Group) Get(name string) *Camera {
	return g.cameras[name]
}

// Remove はグループから Camera を除外します。
func (g *Group) Remove(cam *Camera) {
	if cam == nil {
		return
	}
	key := cam.Name()
	if key == "" {
		key = fmt.Sprintf("%p", cam)
	}
	delete(g.cameras, key)
	g.rebuildList()
}

// Update はグループ内のすべての Camera の Update (Shake 減衰など) を一括更新します。
func (g *Group) Update(dt float64) {
	for _, c := range g.list {
		c.Update(dt)
	}
}

// Render はグループ内の全 Camera を ZIndex 順 (昇順: 値が小さいカメラが背面、大きいカメラが前面) に
// 自動ソートして一括レンダリングします。
func (g *Group) Render(screen *ebiten.Image, drawFunc func(cam *Camera, target *ebiten.Image)) {
	if screen == nil || drawFunc == nil {
		return
	}

	g.sortList()

	for _, cam := range g.list {
		cam.Render(screen, func(target *ebiten.Image) {
			drawFunc(cam, target)
		})
	}
}

func (g *Group) rebuildList() {
	g.list = make([]*Camera, 0, len(g.cameras))
	for _, c := range g.cameras {
		g.list = append(g.list, c)
	}
	g.sortList()
}

func (g *Group) sortList() {
	sort.SliceStable(g.list, func(i, j int) bool {
		return g.list[i].zIndex < g.list[j].zIndex
	})
}

