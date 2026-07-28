package box2d

import (
	"fmt"
	"github.com/ByteArena/box2d"
	"github.com/sofia-gros/ebiten/physics"
)

// b2WorldWrapper implements physics.World for Box2D
type b2WorldWrapper struct {
	world    *box2d.B2World
	listener *contactListener
	bodies   []*b2BodyWrapper
}

// NewWorld creates a new Box2D physics world adapter.
func NewWorld() physics.World {
	// Initialize world with gravity so things fall by default (e.g. 100 pixels/sec^2)
	gravity := box2d.MakeB2Vec2(0.0, 100.0)
	w := box2d.MakeB2World(gravity)

	wrapper := &b2WorldWrapper{
		world:  &w,
		bodies: make([]*b2BodyWrapper, 0),
	}

	// Setup contact listener
	wrapper.listener = &contactListener{}
	wrapper.world.SetContactListener(wrapper.listener)

	return wrapper
}

func (w *b2WorldWrapper) Step(dt float64) {
	// Typically 8 velocity iterations and 3 position iterations are recommended.
	w.world.Step(dt, 8, 3)
}

func (w *b2WorldWrapper) SetGravity(gx, gy float64) {
	w.world.SetGravity(box2d.MakeB2Vec2(gx, gy))
}

func (w *b2WorldWrapper) Gravity() (float64, float64) {
	g := w.world.GetGravity()
	return g.X, g.Y
}

func (w *b2WorldWrapper) CreateBody(options physics.BodyOptions) physics.Body {
	def := box2d.MakeB2BodyDef()
	def.Position.Set(options.X, options.Y)
	
	if options.Type == physics.BodyTypeDynamic {
		def.Type = box2d.B2BodyType.B2_dynamicBody
	} else {
		def.Type = box2d.B2BodyType.B2_staticBody
	}

	b2Body := w.world.CreateBody(&def)

	wrapperBody := &b2BodyWrapper{
		world:   w,
		b2Body:  b2Body,
		options: options,
	}
	b2Body.SetUserData(wrapperBody)

	// Create fixtures for each shape
	for _, shapeDef := range options.Shapes {
		fd := box2d.MakeB2FixtureDef()
		fd.Friction = options.Friction
		fd.Restitution = options.Restitution
		fd.Density = options.Density
		fd.IsSensor = options.IsSensor

		switch s := shapeDef.Shape.(type) {
		case physics.BoxShape:
			poly := box2d.MakeB2PolygonShape()
			center := box2d.MakeB2Vec2(shapeDef.OffsetX, shapeDef.OffsetY)
			poly.SetAsBoxFromCenterAndAngle(s.Width/2, s.Height/2, center, shapeDef.Rotation)
			fd.Shape = &poly
		case physics.CircleShape:
			circ := box2d.MakeB2CircleShape()
			circ.M_radius = s.Radius
			circ.M_p = box2d.MakeB2Vec2(shapeDef.OffsetX, shapeDef.OffsetY)
			fd.Shape = &circ
		}

		b2Body.CreateFixtureFromDef(&fd)
	}

	// We store the wrapper inside the Box2D body's UserData so the ContactListener can find it.
	b2Body.SetUserData(wrapperBody)

	w.bodies = append(w.bodies, wrapperBody)
	return wrapperBody
}

func (w *b2WorldWrapper) RemoveBody(body physics.Body) {
	bw, ok := body.(*b2BodyWrapper)
	if !ok {
		return
	}
	w.world.DestroyBody(bw.b2Body)

	for i, b := range w.bodies {
		if b == bw {
			w.bodies = append(w.bodies[:i], w.bodies[i+1:]...)
			break
		}
	}
}

func (w *b2WorldWrapper) Bodies() []physics.Body {
	res := make([]physics.Body, len(w.bodies))
	for i, b := range w.bodies {
		res[i] = b
	}
	return res
}

// b2BodyWrapper implements physics.Body for Box2D
type b2BodyWrapper struct {
	world   *b2WorldWrapper
	b2Body  *box2d.B2Body
	options physics.BodyOptions
	data    interface{}
	group   physics.Group
}

func (b *b2BodyWrapper) Position() (float64, float64) {
	pos := b.b2Body.GetPosition()
	return pos.X, pos.Y
}

func (b *b2BodyWrapper) SetPosition(x, y float64) {
	b.b2Body.SetTransform(box2d.MakeB2Vec2(x, y), b.b2Body.GetAngle())
}

func (b *b2BodyWrapper) Rotation() float64 {
	return b.b2Body.GetAngle()
}

func (b *b2BodyWrapper) SetRotation(angle float64) {
	b.b2Body.SetTransform(b.b2Body.GetPosition(), angle)
}

func (b *b2BodyWrapper) Velocity() (float64, float64) {
	v := b.b2Body.GetLinearVelocity()
	return v.X, v.Y
}

func (b *b2BodyWrapper) SetVelocity(vx, vy float64) {
	b.b2Body.SetLinearVelocity(box2d.MakeB2Vec2(vx, vy))
}

func (b *b2BodyWrapper) AngularVelocity() float64 {
	return b.b2Body.GetAngularVelocity()
}

func (b *b2BodyWrapper) SetAngularVelocity(w float64) {
	b.b2Body.SetAngularVelocity(w)
}

func (b *b2BodyWrapper) ApplyForce(fx, fy float64) {
	center := b.b2Body.GetWorldCenter()
	b.b2Body.ApplyForce(box2d.MakeB2Vec2(fx, fy), center, true)
}

func (b *b2BodyWrapper) Type() physics.BodyType {
	return b.options.Type
}

func (b *b2BodyWrapper) Shapes() []physics.ShapeDef {
	return b.options.Shapes
}

func (b *b2BodyWrapper) SetGroup(group physics.Group) { b.group = group }
func (b *b2BodyWrapper) Group() physics.Group         { return b.group }

func (b *b2BodyWrapper) SetData(data interface{}) { b.data = data }
func (b *b2BodyWrapper) Data() interface{}        { return b.data }

// contactListener implements box2d.B2ContactListenerInterface
type contactListener struct {
	wrapper *b2WorldWrapper
}

func (c *contactListener) getBodies(contact box2d.B2ContactInterface) (*b2BodyWrapper, *b2BodyWrapper) {
	fixA := contact.GetFixtureA()
	fixB := contact.GetFixtureB()
	if fixA == nil || fixB == nil {
		return nil, nil
	}
	bodyA := fixA.GetBody().GetUserData()
	bodyB := fixB.GetBody().GetUserData()

	wa, _ := bodyA.(*b2BodyWrapper)
	wb, _ := bodyB.(*b2BodyWrapper)
	return wa, wb
}

func (c *contactListener) BeginContact(contact box2d.B2ContactInterface) {
	fmt.Println("BeginContact CALLED!")
	wa, wb := c.getBodies(contact)
	if wa == nil || wb == nil {
		fmt.Println("[Box2D Debug] getBodies returned nil. wa:", wa, "wb:", wb)
		return
	}
	if wa.options.OnCollisionBegin != nil {
		wa.options.OnCollisionBegin(wb)
	}
	if wb.options.OnCollisionBegin != nil {
		wb.options.OnCollisionBegin(wa)
	}
}

func (c *contactListener) EndContact(contact box2d.B2ContactInterface) {
	wa, wb := c.getBodies(contact)
	if wa == nil || wb == nil {
		return
	}
	if wa.options.OnCollisionEnd != nil {
		wa.options.OnCollisionEnd(wb)
	}
	if wb.options.OnCollisionEnd != nil {
		wb.options.OnCollisionEnd(wa)
	}
}

func (c *contactListener) PreSolve(contact box2d.B2ContactInterface, oldManifold box2d.B2Manifold) {
	// PreSolve fires every frame while touching
	wa, wb := c.getBodies(contact)
	if wa == nil || wb == nil {
		return
	}
	if wa.options.OnOverlap != nil {
		wa.options.OnOverlap(wb)
	}
	if wb.options.OnOverlap != nil {
		wb.options.OnOverlap(wa)
	}
}

func (c *contactListener) PostSolve(contact box2d.B2ContactInterface, impulse *box2d.B2ContactImpulse) {
}
