package physics

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Renderable pairs a physical Body with an ebiten.Image for automatic drawing.
type Renderable struct {
	Body  Body
	Image *ebiten.Image
}

// Manager orchestrates the physics simulation and optional rendering.
type Manager struct {
	world       World
	renderables []Renderable
	debugColor  color.Color
}

// NewManager creates a new physics Manager.
// Note: You must call SetWorld() before starting the simulation.
func NewManager() *Manager {
	return &Manager{
		renderables: make([]Renderable, 0),
		debugColor:  color.RGBA{R: 255, G: 0, B: 0, A: 255}, // Default debug color is red
	}
}

// SetWorld injects a specific physics engine adapter into the manager.
func (m *Manager) SetWorld(w World) {
	m.world = w
}

// SetGravity sets the global gravity for the current physics world.
func (m *Manager) SetGravity(gx, gy float64) {
	if m.world != nil {
		m.world.SetGravity(gx, gy)
	}
}

// Gravity returns the current global gravity.
func (m *Manager) Gravity() (float64, float64) {
	if m.world != nil {
		return m.world.Gravity()
	}
	return 0, 0
}

// World returns the current physics world.
func (m *Manager) World() World {
	return m.world
}

// CreateBody creates a new physical body in the active world.
func (m *Manager) CreateBody(options BodyOptions) Body {
	if m.world == nil {
		panic("physics: CreateBody called before SetWorld")
	}

	// Normalize single Shape into Shapes slice for convenience
	if len(options.Shapes) == 0 && options.Shape != nil {
		options.Shapes = []ShapeDef{
			{Shape: options.Shape},
		}
	}

	return m.world.CreateBody(options)
}


// RemoveBody removes a body from the active world and also removes it from renderables.
func (m *Manager) RemoveBody(body Body) {
	if m.world != nil {
		m.world.RemoveBody(body)
	}

	// Remove from renderables if it exists
	for i := len(m.renderables) - 1; i >= 0; i-- {
		if m.renderables[i].Body == body {
			m.renderables = append(m.renderables[:i], m.renderables[i+1:]...)
		}
	}
}

// AddRenderable registers a body and its image for automatic batch rendering when Draw() is called.
func (m *Manager) AddRenderable(body Body, img *ebiten.Image) {
	m.renderables = append(m.renderables, Renderable{
		Body:  body,
		Image: img,
	})
}

// Update advances the physics simulation by dt (delta time in seconds).
func (m *Manager) Update(dt float64) {
	if m.world != nil {
		m.world.Step(dt)
	}
}

// Draw automatically draws all registered renderables.
// This is strictly opt-in; you can completely ignore this method and draw bodies manually.
func (m *Manager) Draw(screen *ebiten.Image) {
	for _, r := range m.renderables {
		x, y := r.Body.Position()
		rot := r.Body.Rotation()
		op := &ebiten.DrawImageOptions{}

		// Assuming the image center should match the body center
		w, h := r.Image.Bounds().Dx(), r.Image.Bounds().Dy()
		op.GeoM.Translate(-float64(w)/2, -float64(h)/2) // Center the image

		if rot != 0 {
			op.GeoM.Rotate(rot)
		}
		op.GeoM.Translate(x, y) // Move to world position
		screen.DrawImage(r.Image, op)
	}
}

// DrawDebug renders the wireframes of all bodies in the physics world for debugging purposes.
func (m *Manager) DrawDebug(screen *ebiten.Image) {
	if m.world == nil {
		return
	}

	for _, b := range m.world.Bodies() {
		bx, by := b.Position()
		bRot := b.Rotation()

		// Get the shape definitions from the body options (if available)
		// We'll assume the body holds its creation options in its UserData for simplicity,
		// or the adapter exposes its shapes. 
		// Since our Body interface doesn't strictly expose shapes yet, 
		// we'll use a hack to read it if the data is a BodyOptions struct.
		// A proper implementation might require a Shapes() method on the Body interface.
		
		// For now, let's add a hypothetical Shapes() method if we can, or just draw a point.
		// Wait, the user didn't request a Shapes() method on Body, but we need it to draw debug.
		// Let's add a DrawDebug proxy if we want, or define how to get shapes.
		
		// Let's assume the underlying Body exposes its shapes, or we just draw a point for now if we can't.
		// Actually, let's check if the Body implements an interface that provides shapes.
		if shapeProvider, ok := b.(interface{ Shapes() []ShapeDef }); ok {
			for _, sdef := range shapeProvider.Shapes() {
				// Calculate world position of the shape
				sin, cos := math.Sincos(bRot)
				// Rotate offset
				ox := sdef.OffsetX*cos - sdef.OffsetY*sin
				oy := sdef.OffsetX*sin + sdef.OffsetY*cos
				sx, sy := bx+ox, by+oy
				sRot := bRot + sdef.Rotation

				switch s := sdef.Shape.(type) {
				case BoxShape:
					drawRotatedRect(screen, sx, sy, s.Width, s.Height, sRot, m.debugColor)
				case CircleShape:
					drawCircle(screen, sx, sy, s.Radius, m.debugColor)
				}
			}
		} else {
			// Fallback: Just draw a point at the body's center
			drawCircle(screen, bx, by, 2, m.debugColor)
		}
	}
}

func drawRotatedRect(screen *ebiten.Image, cx, cy, w, h, rot float64, clr color.Color) {
	sin, cos := math.Sincos(rot)
	hw, hh := float32(w/2), float32(h/2)
	
	// Corners relative to center
	corners := []struct{ x, y float32 }{
		{-hw, -hh}, {hw, -hh}, {hw, hh}, {-hw, hh},
	}

	// Rotate and translate
	var pts [4]struct{ x, y float32 }
	for i, c := range corners {
		rx := float64(c.x)*cos - float64(c.y)*sin
		ry := float64(c.x)*sin + float64(c.y)*cos
		pts[i].x = float32(cx + rx)
		pts[i].y = float32(cy + ry)
	}

	vector.StrokeLine(screen, pts[0].x, pts[0].y, pts[1].x, pts[1].y, 1, clr, true)
	vector.StrokeLine(screen, pts[1].x, pts[1].y, pts[2].x, pts[2].y, 1, clr, true)
	vector.StrokeLine(screen, pts[2].x, pts[2].y, pts[3].x, pts[3].y, 1, clr, true)
	vector.StrokeLine(screen, pts[3].x, pts[3].y, pts[0].x, pts[0].y, 1, clr, true)
}

func drawCircle(screen *ebiten.Image, cx, cy, r float64, clr color.Color) {
	// A simple circle using a polygon approximation
	segments := 16
	var prevX, prevY float32
	
	for i := 0; i <= segments; i++ {
		angle := float64(i) * 2 * math.Pi / float64(segments)
		sin, cos := math.Sincos(angle)
		px := float32(cx + r*cos)
		py := float32(cy + r*sin)
		
		if i > 0 {
			vector.StrokeLine(screen, prevX, prevY, px, py, 1, clr, true)
		}
		prevX, prevY = px, py
	}
}
