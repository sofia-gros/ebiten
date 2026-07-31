module github.com/sofia-gros/ebiten/ui

go 1.25.1

require (
	github.com/hajimehoshi/ebiten/v2 v2.9.9
	github.com/sofia-gros/ebiten/camera v0.0.0
	github.com/sofia-gros/ebiten/pad v0.0.0
)

replace (
	github.com/sofia-gros/ebiten/camera => ../camera
	github.com/sofia-gros/ebiten/pad => ../pad
)
