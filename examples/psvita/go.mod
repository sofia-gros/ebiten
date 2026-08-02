module ebiten/examples/psvita

go 1.26.5

require (
	github.com/hajimehoshi/ebiten/v2 v2.8.6
	ebiten/camera v0.0.0
	ebiten/physics v0.0.0
	ebiten/psvita v0.0.0
)

replace (
	ebiten/camera => "../../camera"
	ebiten/physics => "../../physics"
	ebiten/psvita => "../../psvita"
)
