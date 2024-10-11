package main

import (
	"math"

	"github.com/gdamore/tcell/v2"
)

func main() {
	// Initialize tcell screen
	s, err := tcell.NewScreen()
	if err != nil {
		panic(err)
	}
	if err := s.Init(); err != nil {
		panic(err)
	}
	defer s.Fini()

	s.SetStyle(tcell.StyleDefault.Foreground(tcell.ColorWhite))
	s.Clear()

	quit := make(chan struct{})

	// Initialize buffers once
	initBuffers()

	// Initialize the cube (if needed)
	cube1 := &Cube{
		CubeWidth:        4,
		DistanceFromCam:  13,
		HorizontalOffset: -46, // Left side of the screen
		K1:               40,
		IncrementSpeed:   0.048,
		CubeFaces:        defaultCubeFaces,
	}

	// Initialize the sphere
	sphere := &Sphere{
		Radius:           1.05,
		DistanceFromCam:  1.1,
		HorizontalOffset: 0,
		VerticalOffset:   0,
		K1:               20,
		K2:               1.7,   // Aspect ratio
		Resolution:       0.002, // Adjust for desired quality
		ColorFunction:    checkerboardColorFunction,
                RotationSpeed: 0.004,
                C: math.Pi/3,
	}
	sphere.BuildSurface()

	// Event handling goroutine
	go handleEvents(s, quit, cube1, sphere)

	// Main rendering loop
	renderLoop(s, quit, cube1, sphere)
}

func initBuffers() {
	zBuffer = make([]float64, maxWidth*maxHeight)
	buffer = make([]rune, maxWidth*maxHeight)
	colors = make([]tcell.Color, maxWidth*maxHeight)
}

