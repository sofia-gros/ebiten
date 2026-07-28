package phygo

import (
	"github.com/ab-dek/Phygo-2D"
	"github.com/sofia-gros/ebiten/physics"
)

// phygoWorld implements physics.World for Phygo-2D
type phygoWorld struct {
	bodies []*phygoBodyWrapper
	gx, gy float64
}

// NewWorld creates a new Phygo-2D physics world adapter.
func NewWorld() physics.World {
	wrapper := &phygoWorld{
		bodies: make([]*phygoBodyWrapper, 0),
	}
	// Phygo uses global state, but we track wrappers here.
	wrapper.SetGravity(0, 100.0)
	return wrapper
}

func (w *phygoWorld) Step(dt float64) {
	phygo.UpdatePhysics(float32(dt))

	// Manually check overlaps since Phygo doesn't expose ContactListener
	for i := 0; i < len(w.bodies); i++ {
		b1 := w.bodies[i]
		for j := i + 1; j < len(w.bodies); j++ {
			b2 := w.bodies[j]

			// skip static vs static
			if b1.options.Type == physics.BodyTypeStatic && b2.options.Type == physics.BodyTypeStatic {
				continue
			}

			hit, _, _ := phygo.CheckCollision(b1.phygoBody, b2.phygoBody)

			if hit {
				// overlap started
				if !b1.lastOverlaps[b2] {
					if b1.options.OnCollisionBegin != nil {
						b1.options.OnCollisionBegin(b2)
					}
					if b2.options.OnCollisionBegin != nil {
						b2.options.OnCollisionBegin(b1)
					}
				}
				// overlapping
				if b1.options.OnOverlap != nil {
					b1.options.OnOverlap(b2)
				}
				if b2.options.OnOverlap != nil {
					b2.options.OnOverlap(b1)
				}
				b1.currentOverlaps[b2] = true
				b2.currentOverlaps[b1] = true
			}
		}
	}

	// overlap ended
	for _, b := range w.bodies {
		for other := range b.lastOverlaps {
			if !b.currentOverlaps[other] {
				if b.options.OnCollisionEnd != nil {
					b.options.OnCollisionEnd(other)
				}
			}
		}
		b.lastOverlaps = b.currentOverlaps
		b.currentOverlaps = make(map[*phygoBodyWrapper]bool)
	}
}

func (w *phygoWorld) SetGravity(gx, gy float64) {
	w.gx = gx
	w.gy = gy
	phygo.SetGravity(float32(gx/phygoVelScale), float32(gy/phygoVelScale))
}

func (w *phygoWorld) Gravity() (float64, float64) {
	return w.gx, w.gy
}

func (w *phygoWorld) CreateBody(options physics.BodyOptions) physics.Body {
	isStatic := options.Type == physics.BodyTypeStatic
	pos := phygo.NewVector(float32(options.X), float32(options.Y))
	density := float32(options.Density)
	if density == 0 {
		density = 1.0
	}

	var pb *phygo.Body
	
	// Phygo only supports 1 shape per body, we take the first one
	if len(options.Shapes) > 0 {
		switch s := options.Shapes[0].Shape.(type) {
		case physics.BoxShape:
			pb = phygo.CreateBodyRectangle(pos, float32(s.Width), float32(s.Height), density, isStatic)
		case physics.CircleShape:
			pb = phygo.CreateBodyCircle(pos, float32(s.Radius), density, isStatic)
		}
	}
	
	// Fallback to a tiny box if no shape
	if pb == nil {
		pb = phygo.CreateBodyRectangle(pos, 1, 1, density, isStatic)
	}

	pb.UseGravity = true
	pb.SetRestitution(float32(options.Restitution))
	pb.SetStaticFriction(float32(options.Friction))
	pb.SetDynamicFriction(float32(options.Friction))

	wrapperBody := &phygoBodyWrapper{
		world:           w,
		options:         options,
		phygoBody:       pb,
		lastOverlaps:    make(map[*phygoBodyWrapper]bool),
		currentOverlaps: make(map[*phygoBodyWrapper]bool),
	}

	w.bodies = append(w.bodies, wrapperBody)
	return wrapperBody
}

func (w *phygoWorld) RemoveBody(body physics.Body) {
	bw, ok := body.(*phygoBodyWrapper)
	if !ok {
		return
	}
	phygo.RemoveBody(bw.phygoBody)

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

// phygoScale is 50.0 (ppu in Phygo-2D).
// Because Phygo-2D does `position += Velocity * ppu * time` and `GetPos()` does `position * ppu`,
// effective pixel movement is `Velocity * ppu^2 * time`.
// Thus we must divide all velocities and forces by 2500.0 to get pixel-perfect values.
const phygoVelScale = 2500.0

// phygoBodyWrapper implements physics.Body for Phygo-2D
type phygoBodyWrapper struct {
	world           *phygoWorld
	options         physics.BodyOptions
	data            interface{}
	group           physics.Group
	phygoBody       *phygo.Body
	lastOverlaps    map[*phygoBodyWrapper]bool
	currentOverlaps map[*phygoBodyWrapper]bool
}

func (b *phygoBodyWrapper) Position() (float64, float64) {
	pos := b.phygoBody.GetPos()
	return float64(pos.X), float64(pos.Y)
}

func (b *phygoBodyWrapper) SetPosition(x, y float64) {
	b.phygoBody.MoveTo(phygo.NewVector(float32(x), float32(y)))
}

func (b *phygoBodyWrapper) Rotation() float64 {
	return float64(b.phygoBody.Rotation)
}

func (b *phygoBodyWrapper) SetRotation(angle float64) {
	b.phygoBody.RotateTo(float32(angle))
}

func (b *phygoBodyWrapper) Velocity() (float64, float64) {
	return float64(b.phygoBody.Velocity.X) * phygoVelScale, float64(b.phygoBody.Velocity.Y) * phygoVelScale
}

func (b *phygoBodyWrapper) SetVelocity(vx, vy float64) {
	b.phygoBody.Velocity = phygo.NewVector(float32(vx/phygoVelScale), float32(vy/phygoVelScale))
}

func (b *phygoBodyWrapper) AngularVelocity() float64 {
	return float64(b.phygoBody.AngularVelocity) * 50.0
}

func (b *phygoBodyWrapper) SetAngularVelocity(w float64) {
	b.phygoBody.AngularVelocity = float32(w / 50.0)
}

func (b *phygoBodyWrapper) ApplyForce(fx, fy float64) {
	b.phygoBody.ApplyForce(phygo.NewVector(float32(fx/phygoVelScale), float32(fy/phygoVelScale)))
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
