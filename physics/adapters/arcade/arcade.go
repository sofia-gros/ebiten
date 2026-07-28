package arcade

import (
	"math"

	"github.com/sofia-gros/ebiten/physics"
)

// arcadeWorld implements physics.World
type arcadeWorld struct {
	bodies []*arcadeBody
	gx, gy float64
}

// NewWorld creates a new Arcade physics world.
func NewWorld() physics.World {
	return &arcadeWorld{
		bodies: make([]*arcadeBody, 0),
	}
}

func (w *arcadeWorld) SetGravity(gx, gy float64) {
	w.gx = gx
	w.gy = gy
}

func (w *arcadeWorld) Gravity() (float64, float64) {
	return w.gx, w.gy
}

func (w *arcadeWorld) Step(dt float64) {
	// 1. Integrate velocities
	for _, b := range w.bodies {
		if b.bType == physics.BodyTypeDynamic {
			b.vx += w.gx * dt
			b.vy += w.gy * dt
			b.x += b.vx * dt
			b.y += b.vy * dt
			b.rot += b.omega * dt
		}
	}

	// 2. Clear current frame overlaps
	for _, b := range w.bodies {
		b.currentOverlaps = make(map[*arcadeBody]bool)
	}

	// 3. Collision Detection and Resolution (naive O(N^2) for simplicity)
	for i := 0; i < len(w.bodies); i++ {
		b1 := w.bodies[i]
		for j := i + 1; j < len(w.bodies); j++ {
			b2 := w.bodies[j]

			// Skip if both are static
			if b1.bType == physics.BodyTypeStatic && b2.bType == physics.BodyTypeStatic {
				continue
			}

			if checkCollision(b1, b2) {
				// Mark as overlapping this frame
				b1.currentOverlaps[b2] = true
				b2.currentOverlaps[b1] = true

				// Callbacks
				fireCollisionCallbacks(b1, b2)

				// Resolution (only if neither is a sensor)
				if !b1.isSensor && !b2.isSensor {
					resolveCollision(b1, b2)
				}
			}
		}
	}

	// 4. Handle EndCollision
	for _, b1 := range w.bodies {
		for b2 := range b1.lastOverlaps {
			if !b1.currentOverlaps[b2] {
				// Collision ended
				if b1.options.OnCollisionEnd != nil {
					b1.options.OnCollisionEnd(b2)
				}
				if b2.options.OnCollisionEnd != nil && !b2.alreadyFiredEnd(b1) {
					b2.options.OnCollisionEnd(b1)
				}
			}
		}
		// Swap maps
		b1.lastOverlaps = b1.currentOverlaps
	}
}

func (w *arcadeWorld) CreateBody(options physics.BodyOptions) physics.Body {
	b := &arcadeBody{
		world:           w,
		options:         options,
		bType:           options.Type,
		x:               options.X,
		y:               options.Y,
		isSensor:        options.IsSensor,
		shapes:          options.Shapes,
		lastOverlaps:    make(map[*arcadeBody]bool),
		currentOverlaps: make(map[*arcadeBody]bool),
	}
	w.bodies = append(w.bodies, b)
	return b
}

func (w *arcadeWorld) RemoveBody(body physics.Body) {
	ab, ok := body.(*arcadeBody)
	if !ok {
		return
	}
	for i, b := range w.bodies {
		if b == ab {
			w.bodies = append(w.bodies[:i], w.bodies[i+1:]...)
			break
		}
	}
}

func (w *arcadeWorld) Bodies() []physics.Body {
	res := make([]physics.Body, len(w.bodies))
	for i, b := range w.bodies {
		res[i] = b
	}
	return res
}

// arcadeBody implements physics.Body
type arcadeBody struct {
	world   *arcadeWorld
	options physics.BodyOptions

	bType physics.BodyType
	x, y  float64
	vx, vy float64
	rot    float64
	omega  float64

	shapes []physics.ShapeDef

	isSensor bool
	userData interface{}
	group    physics.Group

	// Tracking for Begin/End events
	lastOverlaps    map[*arcadeBody]bool
	currentOverlaps map[*arcadeBody]bool
}

func (b *arcadeBody) Position() (float64, float64) { return b.x, b.y }
func (b *arcadeBody) SetPosition(x, y float64)     { b.x, b.y = x, y }
func (b *arcadeBody) Rotation() float64            { return b.rot }
func (b *arcadeBody) SetRotation(r float64)        { b.rot = r }
func (b *arcadeBody) Velocity() (float64, float64) { return b.vx, b.vy }
func (b *arcadeBody) SetVelocity(vx, vy float64)   { b.vx, b.vy = vx, vy }
func (b *arcadeBody) AngularVelocity() float64     { return b.omega }
func (b *arcadeBody) SetAngularVelocity(w float64) { b.omega = w }
func (b *arcadeBody) Type() physics.BodyType       { return b.bType }

func (b *arcadeBody) SetGroup(group physics.Group) { b.group = group }
func (b *arcadeBody) Group() physics.Group         { return b.group }

func (b *arcadeBody) SetData(data interface{}) { b.userData = data }
func (b *arcadeBody) Data() interface{}        { return b.userData }

func (b *arcadeBody) Shapes() []physics.ShapeDef   { return b.shapes }
func (b *arcadeBody) ApplyForce(fx, fy float64) {
	// Arcade engine ignores mass for simplicity, so force just acts as an instant velocity change.
	b.vx += fx
	b.vy += fy
}

func (b *arcadeBody) alreadyFiredEnd(other *arcadeBody) bool {
	// helper to prevent double-firing end callbacks
	return false
}

// ---- Collision Logic (Very basic AABB / Circle) ----

func fireCollisionCallbacks(b1, b2 *arcadeBody) {
	// b1 perspectives
	if !b1.lastOverlaps[b2] {
		if b1.options.OnCollisionBegin != nil {
			b1.options.OnCollisionBegin(b2)
		}
	}
	if b1.options.OnOverlap != nil {
		b1.options.OnOverlap(b2)
	}

	// b2 perspectives
	if !b2.lastOverlaps[b1] {
		if b2.options.OnCollisionBegin != nil {
			b2.options.OnCollisionBegin(b1)
		}
	}
	if b2.options.OnOverlap != nil {
		b2.options.OnOverlap(b1)
	}
}

func checkCollision(b1, b2 *arcadeBody) bool {
	// Compound shape checking
	for _, s1 := range b1.shapes {
		for _, s2 := range b2.shapes {
			if checkShapeCollision(b1, s1, b2, s2) {
				return true
			}
		}
	}
	return false
}

func checkShapeCollision(b1 *arcadeBody, s1 physics.ShapeDef, b2 *arcadeBody, s2 physics.ShapeDef) bool {
	// Calculate absolute positions of shapes (ignoring rotation for Arcade AABB for simplicity)
	x1 := b1.x + s1.OffsetX
	y1 := b1.y + s1.OffsetY
	x2 := b2.x + s2.OffsetX
	y2 := b2.y + s2.OffsetY

	// Box vs Box
	if box1, ok1 := s1.Shape.(physics.BoxShape); ok1 {
		if box2, ok2 := s2.Shape.(physics.BoxShape); ok2 {
			return math.Abs(x1-x2) < (box1.Width/2+box2.Width/2) && math.Abs(y1-y2) < (box1.Height/2+box2.Height/2)
		}
		if circ2, ok2 := s2.Shape.(physics.CircleShape); ok2 {
			return circleBoxCollide(x2, y2, circ2.Radius, x1, y1, box1.Width, box1.Height)
		}
	}
	// Circle vs Circle
	if circ1, ok1 := s1.Shape.(physics.CircleShape); ok1 {
		if circ2, ok2 := s2.Shape.(physics.CircleShape); ok2 {
			dx := x1 - x2
			dy := y1 - y2
			dist := math.Sqrt(dx*dx + dy*dy)
			return dist < (circ1.Radius + circ2.Radius)
		}
		if box2, ok2 := s2.Shape.(physics.BoxShape); ok2 {
			return circleBoxCollide(x1, y1, circ1.Radius, x2, y2, box2.Width, box2.Height)
		}
	}
	return false
}

func circleBoxCollide(cx, cy, cr, bx, by, bw, bh float64) bool {
	// Find closest point on box to circle center
	hw, hh := bw/2, bh/2
	clampX := math.Max(bx-hw, math.Min(cx, bx+hw))
	clampY := math.Max(by-hh, math.Min(cy, by+hh))

	dx := cx - clampX
	dy := cy - clampY
	return (dx*dx + dy*dy) < (cr * cr)
}

func resolveCollision(b1, b2 *arcadeBody) {
	if b1.bType != physics.BodyTypeDynamic && b2.bType != physics.BodyTypeDynamic {
		return
	}

	for _, s1 := range b1.shapes {
		for _, s2 := range b2.shapes {
			box1, ok1 := s1.Shape.(physics.BoxShape)
			box2, ok2 := s2.Shape.(physics.BoxShape)
			if ok1 && ok2 {
				x1, y1 := b1.x+s1.OffsetX, b1.y+s1.OffsetY
				x2, y2 := b2.x+s2.OffsetX, b2.y+s2.OffsetY

				dx := x2 - x1
				dy := y2 - y1

				halfWidths := box1.Width/2 + box2.Width/2
				halfHeights := box1.Height/2 + box2.Height/2

				overlapX := halfWidths - math.Abs(dx)
				overlapY := halfHeights - math.Abs(dy)

				if overlapX > 0 && overlapY > 0 {
					// Resolve on the axis of least penetration
					if overlapX < overlapY {
						if dx > 0 {
							if b1.bType == physics.BodyTypeDynamic { b1.x -= overlapX; b1.vx = 0 }
							if b2.bType == physics.BodyTypeDynamic { b2.x += overlapX; b2.vx = 0 }
						} else {
							if b1.bType == physics.BodyTypeDynamic { b1.x += overlapX; b1.vx = 0 }
							if b2.bType == physics.BodyTypeDynamic { b2.x -= overlapX; b2.vx = 0 }
						}
					} else {
						if dy > 0 {
							if b1.bType == physics.BodyTypeDynamic { b1.y -= overlapY; b1.vy = 0 }
							if b2.bType == physics.BodyTypeDynamic { b2.y += overlapY; b2.vy = 0 }
						} else {
							if b1.bType == physics.BodyTypeDynamic { b1.y += overlapY; b1.vy = 0 }
							if b2.bType == physics.BodyTypeDynamic { b2.y -= overlapY; b2.vy = 0 }
						}
					}
					return // Resolved
				}
			}
		}
	}
}
