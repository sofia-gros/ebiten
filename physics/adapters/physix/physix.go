package physix

import (
	"github.com/sofia-gros/ebiten/physics"
)

// physixWorld implements physics.World for Physix-go
type physixWorld struct {
	bodies []*physixBodyWrapper
}

// NewWorld creates a new Physix-go physics world adapter.
func NewWorld() physics.World {
	// Initialize Physix Environment
	wrapper := &physixWorld{
		bodies: make([]*physixBodyWrapper, 0),
	}
	return wrapper
}

func (w *physixWorld) Step(dt float64) {
	// Call physix step
}

func (w *physixWorld) CreateBody(options physics.BodyOptions) physics.Body {
	// Create the body in Physix
	wrapperBody := &physixBodyWrapper{
		world:   w,
		options: options,
	}

	w.bodies = append(w.bodies, wrapperBody)
	return wrapperBody
}

func (w *physixWorld) RemoveBody(body physics.Body) {
	bw, ok := body.(*physixBodyWrapper)
	if !ok {
		return
	}

	for i, b := range w.bodies {
		if b == bw {
			w.bodies = append(w.bodies[:i], w.bodies[i+1:]...)
			break
		}
	}
}

func (w *physixWorld) Bodies() []physics.Body {
	res := make([]physics.Body, len(w.bodies))
	for i, b := range w.bodies {
		res[i] = b
	}
	return res
}

// physixBodyWrapper implements physics.Body for Physix-go
type physixBodyWrapper struct {
	world   *physixWorld
	options physics.BodyOptions
	data    interface{}
}

func (b *physixBodyWrapper) Position() (float64, float64) {
	return b.options.X, b.options.Y
}

func (b *physixBodyWrapper) SetPosition(x, y float64) {
	b.options.X, b.options.Y = x, y
}

func (b *physixBodyWrapper) Rotation() float64 {
	return 0
}

func (b *physixBodyWrapper) SetRotation(angle float64) {
}

func (b *physixBodyWrapper) Velocity() (float64, float64) {
	return 0, 0
}

func (b *physixBodyWrapper) SetVelocity(vx, vy float64) {
}

func (b *physixBodyWrapper) AngularVelocity() float64 {
	return 0
}

func (b *physixBodyWrapper) SetAngularVelocity(w float64) {
}

func (b *physixBodyWrapper) ApplyForce(fx, fy float64) {
}

func (b *physixBodyWrapper) Type() physics.BodyType {
	return b.options.Type
}

func (b *physixBodyWrapper) Shapes() []physics.ShapeDef {
	return b.options.Shapes
}

func (b *physixBodyWrapper) SetData(data interface{}) {
	b.data = data
}

func (b *physixBodyWrapper) Data() interface{} {
	return b.data
}
