package phygo

import (
	"github.com/sofia-gros/ebiten/physics"
)

// phygoWorld implements physics.World for Phygo-2D
type phygoWorld struct {
	bodies []*phygoBodyWrapper
}

// NewWorld creates a new Phygo-2D physics world adapter.
func NewWorld() physics.World {
	// Initialize Phygo Space/World
	// e.g. space := phygo.NewSpace()
	wrapper := &phygoWorld{
		bodies: make([]*phygoBodyWrapper, 0),
	}
	return wrapper
}

func (w *phygoWorld) Step(dt float64) {
	// Call phygo step
	// w.space.Step(dt)
}

func (w *phygoWorld) CreateBody(options physics.BodyOptions) physics.Body {
	// Create the body in Phygo
	wrapperBody := &phygoBodyWrapper{
		world:   w,
		options: options,
	}

	// Iterate through Shapes and attach them to Phygo Body
	for _, shapeDef := range options.Shapes {
		switch s := shapeDef.Shape.(type) {
		case physics.BoxShape:
			_ = s
			// phygo.AddBox(...)
		case physics.CircleShape:
			_ = s
			// phygo.AddCircle(...)
		}
	}

	w.bodies = append(w.bodies, wrapperBody)
	return wrapperBody
}

func (w *phygoWorld) RemoveBody(body physics.Body) {
	bw, ok := body.(*phygoBodyWrapper)
	if !ok {
		return
	}
	// w.space.RemoveBody(bw.phygoBody)

	for i, b := range w.bodies {
		if b == bw {
			w.bodies = append(w.bodies[:i], w.bodies[i+1:]...)
			break
		}
	}
}

func (w *phygoWorld) Bodies() []physics.Body {
	res := make([]physics.Body, len(w.bodies))
	for i, b := range w.bodies {
		res[i] = b
	}
	return res
}

// phygoBodyWrapper implements physics.Body for Phygo-2D
type phygoBodyWrapper struct {
	world   *phygoWorld
	options physics.BodyOptions
	data    interface{}
	group   physics.Group
}

func (b *phygoBodyWrapper) Position() (float64, float64) {
	// return b.phygoBody.Position()
	return b.options.X, b.options.Y
}

func (b *phygoBodyWrapper) SetPosition(x, y float64) {
	b.options.X, b.options.Y = x, y
}

func (b *phygoBodyWrapper) Rotation() float64 {
	return 0
}

func (b *phygoBodyWrapper) SetRotation(angle float64) {
}

func (b *phygoBodyWrapper) Velocity() (float64, float64) {
	return 0, 0
}

func (b *phygoBodyWrapper) SetVelocity(vx, vy float64) {
}

func (b *phygoBodyWrapper) AngularVelocity() float64 {
	return 0
}

func (b *phygoBodyWrapper) SetAngularVelocity(w float64) {
}

func (b *phygoBodyWrapper) ApplyForce(fx, fy float64) {
}

func (b *phygoBodyWrapper) Type() physics.BodyType {
	return b.options.Type
}

func (b *phygoBodyWrapper) Shapes() []physics.ShapeDef {
	return b.options.Shapes
}

func (b *phygoBodyWrapper) SetGroup(group physics.Group) { b.group = group }
func (b *phygoBodyWrapper) Group() physics.Group         { return b.group }

func (b *phygoBodyWrapper) SetData(data interface{}) { b.data = data }
func (b *phygoBodyWrapper) Data() interface{}        { return b.data }
