package main

// TODO: Update this to take anything that can qualify as like a type of "Model"
import (
	"math"
	// "time"

	"github.com/gdamore/tcell/v2"
)

func renderLoop(s tcell.Screen, quit chan struct{}, cube1 *Cube, sphere *Sphere) {
	for {
		select {
		case <-quit:
			return
		default:
			w, h := s.Size()
			if w != width || h != height {
				width = w
				height = h
				// Optionally handle resizing logic if necessary
			}

			// Clear buffers without reallocating
			clearBuffers()

			// Generate sphere
			sphere.Generate(s)

			// Optionally generate cube
			// cube1.Generate(s)

			// Draw buffers to the screen
			drawBuffers(s)

			// Show the updated screen
			s.Show()

			// Update sphere rotation angles
			sphere.B += sphere.RotationSpeed
                        sphere.A -= sphere.RotationSpeed

                        // cube1.A += 0.05
                        // cube1.B += 0.05
                        // cube1.C += 0.01
			// time.Sleep(time.Millisecond * 16)
		}
	}
}

func clearBuffers() {
	// Only clear the portion of the buffers that will be used
	size := width * height
	for i := 0; i < size; i++ {
		buffer[i] = backgroundASCIICode
		zBuffer[i] = -math.MaxFloat64
		colors[i] = tcell.ColorBlack
	}
}

func drawBuffers(screen tcell.Screen) {
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			idx := x + y*width
			screen.SetContent(x, y, buffer[idx], nil, tcell.StyleDefault.Foreground(colors[idx]))
		}
	}
}

