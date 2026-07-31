package physics

// BodyType represents the type of a physics body (static, dynamic).
type BodyType int

const (
	BodyTypeStatic BodyType = iota
	BodyTypeDynamic
)

// Group represents a collision group or layer (e.g. 1 for Player, 2 for Enemy).
type Group uint32

// Shape represents a collision shape.
// It uses a string type identifier so adapters can type-assert to specific shapes.
type Shape interface {
	ShapeType() string
}

// BoxShape represents an Axis-Aligned Bounding Box (AABB) or rotatable box shape.
type BoxShape struct {
	Width  float64
	Height float64
}

func (s BoxShape) ShapeType() string { return "box" }

// CircleShape represents a circular collision shape.
type CircleShape struct {
	Radius float64
}

func (s CircleShape) ShapeType() string { return "circle" }

// ShapeDef represents a shape relative to a body's origin, allowing for compound shapes.
type ShapeDef struct {
	Shape    Shape
	OffsetX  float64
	OffsetY  float64
	Rotation float64
}

// CollisionCallback is called for collision events.
// other is the Body that this body collided with.
type CollisionCallback func(other Body)

// BodyOptions contains all possible parameters for creating a new body.
// Some physics engine adapters may not support all properties.
type BodyOptions struct {
	Type BodyType
	X, Y float64

	// Shape is a convenience field for single-shape bodies (e.g., Shape: physics.BoxShape{Width: 32, Height: 32}).
	Shape Shape

	// Shapes allows creating a compound body by connecting multiple shapes together with offsets/rotations.
	Shapes []ShapeDef

	Friction    float64 // Friction coefficient

	Restitution float64 // Bounciness/Restitution coefficient
	Density     float64 // Density (used for mass calculation)
	IsSensor    bool    // If true, the body detects collisions but has no physical response

	// Event callbacks
	OnCollisionBegin CollisionCallback // Fired once when a collision starts
	OnCollisionEnd   CollisionCallback // Fired once when a collision ends
	OnOverlap        CollisionCallback // Fired every frame while bodies overlap
}

// Body represents a physical object in the world.
type Body interface {
	// ---- State Getters ----
	Position() (x, y float64)
	Rotation() float64
	Velocity() (vx, vy float64)
	AngularVelocity() float64
	Type() BodyType

	// ---- State Setters ----
	SetPosition(x, y float64)
	SetRotation(angle float64)
	SetVelocity(vx, vy float64)
	SetAngularVelocity(omega float64)
	ApplyForce(fx, fy float64)

	// ---- Shapes (for Debug Draw) ----
	Shapes() []ShapeDef

	// ---- Collision Group / Tags ----
	SetGroup(group Group)
	Group() Group

	// ---- User Data ----
	SetData(data interface{})
	Data() interface{}
}

// World represents the physics simulation environment.
type World interface {
	Step(dt float64)
	CreateBody(options BodyOptions) Body
	RemoveBody(body Body)
	Bodies() []Body // Required for Debug Draw
	SetGravity(gx, gy float64) // Set global gravity
	Gravity() (float64, float64) // Get global gravity
}
