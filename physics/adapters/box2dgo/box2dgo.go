package box2dgo

import (
	"math"

	b2 "github.com/oliverbestmann/box2d-go"
	"github.com/sofia-gros/ebiten/physics"
)

type box2dgoWorld struct {
	world  b2.World
	bodies []*box2dgoBodyWrapper
	shapeMap map[b2.ShapeId]*box2dgoBodyWrapper
}

func NewWorld() physics.World {
	def := b2.DefaultWorldDef()
	def.Gravity = b2.Vec2{X: 0, Y: 100}

	world := b2.CreateWorld(def)

	return &box2dgoWorld{
		world:    world,
		bodies:   make([]*box2dgoBodyWrapper, 0),
		shapeMap: make(map[b2.ShapeId]*box2dgoBodyWrapper),
	}
}

func (w *box2dgoWorld) Step(dt float64) {
	// box2d-go prefers 4 sub-steps
	w.world.Step(float32(dt), 4)

	// Process events
	contactEvents := w.world.GetContactEvents()
	for _, e := range contactEvents.BeginEvents {
		wa, oka := w.shapeMap[e.ShapeIdA]
		wb, okb := w.shapeMap[e.ShapeIdB]
		if oka && okb {
			if wa.options.OnCollisionBegin != nil {
				wa.options.OnCollisionBegin(wb)
			}
			if wb.options.OnCollisionBegin != nil {
				wb.options.OnCollisionBegin(wa)
			}
		}
	}
	
	for _, e := range contactEvents.EndEvents {
		wa, oka := w.shapeMap[e.ShapeIdA]
		wb, okb := w.shapeMap[e.ShapeIdB]
		if oka && okb {
			if wa.options.OnCollisionEnd != nil {
				wa.options.OnCollisionEnd(wb)
			}
			if wb.options.OnCollisionEnd != nil {
				wb.options.OnCollisionEnd(wa)
			}
		}
	}
}

func (w *box2dgoWorld) SetGravity(gx, gy float64) {
	w.world.SetGravity(b2.Vec2{X: float32(gx), Y: float32(gy)})
}

func (w *box2dgoWorld) Gravity() (float64, float64) {
	g := w.world.GetGravity()
	return float64(g.X), float64(g.Y)
}

func (w *box2dgoWorld) CreateBody(options physics.BodyOptions) physics.Body {
	def := b2.DefaultBodyDef()
	def.Position = b2.Vec2{X: float32(options.X), Y: float32(options.Y)}
	
	if options.Type == physics.BodyTypeDynamic {
		def.Type1 = b2.DynamicBody
	} else {
		def.Type1 = b2.StaticBody
	}

	body := w.world.CreateBody(def)

	wrapper := &box2dgoBodyWrapper{
		world:   w,
		body:    body,
		options: options,
		shapes:  make([]b2.Shape, 0),
	}

	for _, shapeDef := range options.Shapes {
		sd := b2.DefaultShapeDef()
		sd.Density = float32(options.Density)
		if sd.Density == 0 {
			sd.Density = 1.0
		}
		
		sd.IsSensor = 0
		if options.IsSensor {
			sd.IsSensor = 1
		}
		// Enable contact events for this shape
		sd.EnableContactEvents = 1

		
		center := b2.Vec2{X: float32(shapeDef.OffsetX), Y: float32(shapeDef.OffsetY)}
		rot := b2.Rot{C: float32(math.Cos(shapeDef.Rotation)), S: float32(math.Sin(shapeDef.Rotation))}
		
		var shape b2.Shape
		switch s := shapeDef.Shape.(type) {
		case physics.BoxShape:
			poly := b2.MakeOffsetBox(float32(s.Width/2), float32(s.Height/2), center, rot)
			shape = body.CreatePolygonShape(sd, poly)
		case physics.CircleShape:
			circ := b2.Circle{Radius: float32(s.Radius), Center: center}
			shape = body.CreateCircleShape(sd, circ)
		}
		
		shape.SetFriction(float32(options.Friction))
		shape.SetRestitution(float32(options.Restitution))
		
		wrapper.shapes = append(wrapper.shapes, shape)
		w.shapeMap[shape.Id] = wrapper
	}

	w.bodies = append(w.bodies, wrapper)
	return wrapper
}

func (w *box2dgoWorld) RemoveBody(body physics.Body) {
	bw, ok := body.(*box2dgoBodyWrapper)
	if !ok {
		return
	}
	
	for _, shape := range bw.shapes {
		delete(w.shapeMap, shape.Id)
	}
	bw.body.DestroyBody()

	for i, b := range w.bodies {
		if b == bw {
			w.bodies = append(w.bodies[:i], w.bodies[i+1:]...)
			break
		}
	}
}

func (w *box2dgoWorld) Bodies() []physics.Body {
	res := make([]physics.Body, len(w.bodies))
	for i, b := range w.bodies {
		res[i] = b
	}
	return res
}

type box2dgoBodyWrapper struct {
	world   *box2dgoWorld
	body    b2.Body
	shapes  []b2.Shape
	options physics.BodyOptions
	data    interface{}
	group   physics.Group
}

func (b *box2dgoBodyWrapper) Position() (float64, float64) {
	pos := b.body.GetPosition()
	return float64(pos.X), float64(pos.Y)
}

func (b *box2dgoBodyWrapper) SetPosition(x, y float64) {
	rot := b.body.GetRotation()
	b.body.SetTransform(b2.Vec2{X: float32(x), Y: float32(y)}, rot)
}

func (b *box2dgoBodyWrapper) Rotation() float64 {
	rot := b.body.GetRotation()
	return float64(math.Atan2(float64(rot.S), float64(rot.C)))
}

func (b *box2dgoBodyWrapper) SetRotation(angle float64) {
	pos := b.body.GetPosition()
	rot := b2.Rot{C: float32(math.Cos(angle)), S: float32(math.Sin(angle))}
	b.body.SetTransform(pos, rot)
}

func (b *box2dgoBodyWrapper) Velocity() (float64, float64) {
	v := b.body.GetLinearVelocity()
	return float64(v.X), float64(v.Y)
}

func (b *box2dgoBodyWrapper) SetVelocity(vx, vy float64) {
	b.body.SetLinearVelocity(b2.Vec2{X: float32(vx), Y: float32(vy)})
}

func (b *box2dgoBodyWrapper) AngularVelocity() float64 {
	return float64(b.body.GetAngularVelocity())
}

func (b *box2dgoBodyWrapper) SetAngularVelocity(w float64) {
	b.body.SetAngularVelocity(float32(w))
}

func (b *box2dgoBodyWrapper) ApplyForce(fx, fy float64) {
	center := b.body.GetWorldCenterOfMass()
	b.body.ApplyForce(b2.Vec2{X: float32(fx), Y: float32(fy)}, center, 1)
}

func (b *box2dgoBodyWrapper) Type() physics.BodyType {
	return b.options.Type
}

func (b *box2dgoBodyWrapper) Shapes() []physics.ShapeDef {
	return b.options.Shapes
}

func (b *box2dgoBodyWrapper) SetGroup(group physics.Group) { b.group = group }
func (b *box2dgoBodyWrapper) Group() physics.Group         { return b.group }

func (b *box2dgoBodyWrapper) SetData(data interface{}) { b.data = data }
func (b *box2dgoBodyWrapper) Data() interface{}        { return b.data }
