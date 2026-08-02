module ebiten/examples/psvita_build

go 1.26.5

require (
	github.com/hajimehoshi/ebiten/v2 v2.8.6
	github.com/sofia-gros/ebiten/camera v0.0.0
	github.com/sofia-gros/ebiten/physics v0.0.0
	github.com/sofia-gros/ebiten/physics/adapters/arcade v0.0.0
	ebiten/psvita v0.0.0
)

replace (
	github.com/sofia-gros/ebiten/camera => "../../camera"
	github.com/sofia-gros/ebiten/physics => "../../physics"
	github.com/sofia-gros/ebiten/physics/adapters/arcade => "../../physics/adapters/arcade"
	ebiten/psvita => "../../psvita"
)
