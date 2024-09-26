package main

import (
	"math"
	"time"

	"github.com/gdamore/tcell/v2"
)

type Vec3 struct {
	X, Y, Z float64
}

// CubeFace holds the character and color for each face of the cube
type CubeFace struct {
	Char  rune
	Color tcell.Color
}

var (
	width, height       int           = 203, 74
	zBuffer             []float64     // Z-buffer to store depth values
	buffer              []rune        // Buffer to store ASCII
	colors              []tcell.Color // Buffer to store colors
	backgroundASCIICode rune          = ' '
	increment                         = 0.1
	autoSpin                          = true
)

// Cube represents a single cube with its attributes
type Cube struct {
	CubeWidth        float64
	DistanceFromCam  int
	HorizontalOffset float64
	K1               float64
	IncrementSpeed   float64
	A, B, C          float64
	Direction        bool
	CubeFaces        []CubeFace
}
type Sphere struct {
	Radius           float64
	DistanceFromCam  float64
	HorizontalOffset float64
	VerticalOffset   float64
	K1               float64
	RotationSpeed    float64
	A, B, C          float64 // Rotation angles
}

// Lighting to give shading effects (simplified lighting model)
var shading = []rune{' ', '.', '-', '=', '+', '*', '#', '%', '@'}
var grayScale = []tcell.Color{tcell.ColorBlack.TrueColor(), tcell.NewRGBColor(10, 10, 10).TrueColor(), tcell.NewRGBColor(20, 20, 20).TrueColor(), tcell.NewRGBColor(30, 30, 30).TrueColor(), tcell.NewRGBColor(40, 40, 40).TrueColor(), tcell.NewRGBColor(50, 50, 50).TrueColor(), tcell.NewRGBColor(60, 60, 60).TrueColor(), tcell.NewRGBColor(70, 70, 70).TrueColor(), tcell.NewRGBColor(80, 80, 80).TrueColor(), tcell.NewRGBColor(90, 90, 90).TrueColor(), tcell.NewRGBColor(100, 100, 100).TrueColor()}

// Define cube faces with associated colors and characters
var defaultCubeFaces = []CubeFace{
	{'$', tcell.ColorRed},      // Front face
	{'#', tcell.ColorGreen},    // Right face
	{'@', tcell.ColorYellow},   // Back face
	{'&', tcell.ColorBlue},     // Left face
	{'%', tcell.ColorDarkCyan}, // Bottom face
	{'|', tcell.ColorPurple},   // Top face
}

var miniCubeFaces = []CubeFace{
	{'*', tcell.ColorRed},       // Front face
	{'*', tcell.ColorGreen},     // Right face
	{'*', tcell.ColorYellow},    // Back face
	{'*', tcell.ColorBlue},      // Left face
	{'*', tcell.ColorLightBlue}, // Bottom face
	{'*', tcell.ColorPurple},    // Top face
}

func rotateX(p Vec3, theta float64) Vec3 {
	newY := (p.Y * math.Cos(theta)) - (p.Z * math.Sin(theta))
	newZ := (p.Y * math.Sin(theta)) + (p.Z * math.Cos(theta))
	return Vec3{p.X, newY, newZ}
}

func rotateY(p Vec3, theta float64) Vec3 {
	newX := (p.X * math.Cos(theta)) + (p.Z * math.Sin(theta))
	newZ := (p.Z * math.Cos(theta)) - (p.X * math.Sin(theta))
	return Vec3{newX, p.Y, newZ}
}

func rotateZ(p Vec3, theta float64) Vec3 {
	newX := (p.X * math.Cos(theta)) - (p.Y * math.Sin(theta))
	newY := (p.X * math.Sin(theta)) + (p.Y * math.Cos(theta))
	return Vec3{newX, newY, p.Z}
}

func calcNewPoint(p Vec3, rotations Vec3) Vec3 {
	newPoint := rotateZ(p, rotations.Z)
	newPoint = rotateY(newPoint, rotations.Y)
	newPoint = rotateX(newPoint, rotations.X)
	return newPoint
}

func (cube *Cube) Generate(screen tcell.Screen) {
	for cubeX := -cube.CubeWidth; cubeX < cube.CubeWidth; cubeX += cube.IncrementSpeed {
		for cubeY := -cube.CubeWidth; cubeY < cube.CubeWidth; cubeY += cube.IncrementSpeed {
			for i, face := range cube.CubeFaces {
				cubePoint := getCubePoint(i, cubeX, cubeY, cube)
				calculateForSurface(cubePoint, face.Char, face.Color, screen, cube)
			}
		}
	}
}

// Sphere Generate projects points on the surface of the sphere
// FIX: don't work :(
func (sphere *Sphere) Generate(screen tcell.Screen) {
	for theta := 0.0; theta < math.Pi; theta += 0.07 { // Polar angle
		for phi := 0.0; phi < 2*math.Pi; phi += 0.07 { // Azimuthal angle
			// Parametric equations for the sphere
			x := sphere.Radius * math.Sin(theta) * math.Cos(phi)
			y := sphere.Radius * math.Sin(theta) * math.Sin(phi)
			z := sphere.Radius * math.Cos(theta)

			// Rotate the sphere
			point := rotateX(Vec3{x, y, z}, sphere.A)
			point = rotateY(point, sphere.B)
			point = rotateZ(point, sphere.C)

			// Project 3D point to 2D
			zProj := point.Z + sphere.DistanceFromCam
			ooz := 1 / zProj
			xp := int(float64(width)/2 + sphere.HorizontalOffset + sphere.K1*ooz*point.X*2)
			yp := int(float64(height)/2 + sphere.VerticalOffset + sphere.K1*ooz*point.Y)

			idx := xp + yp*width
			if idx >= 0 && idx < width*height {
				if ooz > zBuffer[idx] {
					zBuffer[idx] = ooz
					// Lighting calculation (simplified using z)
					lightIdx := int((ooz - 0.5) * 10)
					if lightIdx < 0 {
						lightIdx = 0
					} else if lightIdx >= len(shading) {
						lightIdx = len(shading) - 1
					}

					buffer[idx] = shading[lightIdx]
					if idx >= 0 && idx < width*height {
						if ooz > zBuffer[idx] {
							zBuffer[idx] = ooz
							buffer[idx] = shading[lightIdx]
							colors[idx] = tcell.ColorPurple
						}
					}

				}
			}
		}
	}
}

func calculateForSurface(cubePoint Vec3, ch rune, color tcell.Color, screen tcell.Screen, cube *Cube) {
	rotatedPoint := calcNewPoint(cubePoint, Vec3{cube.A, cube.B, cube.C})
	z := rotatedPoint.Z + float64(cube.DistanceFromCam)
	ooz := 1 / z
	xp := int(float64(width)/2 + cube.HorizontalOffset + cube.K1*ooz*rotatedPoint.X*2)
	yp := int(float64(height)/2 + cube.K1*ooz*rotatedPoint.Y)
	idx := xp + yp*width
	if idx >= 0 && idx < width*height {
		if ooz > zBuffer[idx] {
			zBuffer[idx] = ooz
			buffer[idx] = ch
			colors[idx] = color
		}
	}
}

// Helper function to determine which point to project for each face
// Helper function to determine which point to project for each face
func getCubePoint(faceIndex int, cubeX, cubeY float64, cube *Cube) Vec3 {
	switch faceIndex {
	case 0:
		return Vec3{cubeX, cubeY, -cube.CubeWidth} // Use cube.CubeWidth
	case 1:
		return Vec3{cube.CubeWidth, cubeY, cubeX} // Use cube.CubeWidth
	case 2:
		return Vec3{-cube.CubeWidth, cubeY, -cubeX} // Use cube.CubeWidth
	case 3:
		return Vec3{-cubeX, cubeY, cube.CubeWidth} // Use cube.CubeWidth
	case 4:
		return Vec3{cubeX, -cube.CubeWidth, -cubeY} // Use cube.CubeWidth
	case 5:
		return Vec3{cubeX, cube.CubeWidth, cubeY} // Use cube.CubeWidth
	}
	return Vec3{}
}

func drawbuffers(screen tcell.Screen) {
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			idx := x + y*width
			screen.SetContent(x, y, buffer[idx], nil, tcell.StyleDefault.Foreground(colors[idx]))
		}
	}

}

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

	cube1 := &Cube{
		CubeWidth:        20,
		DistanceFromCam:  100,
		HorizontalOffset: -40, // Left side of the screen
		K1:               40,
		IncrementSpeed:   0.3,
		CubeFaces:        defaultCubeFaces,
	}
	go func() {
		for {
			ev := s.PollEvent()
			switch ev := ev.(type) {
			case *tcell.EventKey:
				if ev.Key() == tcell.KeyCtrlC || ev.Rune() == 'q' {
					close(quit)
					return
				}

				if ev.Rune() == '[' {
					if increment < 1 {
						increment += 0.1
					}
				}
				if ev.Rune() == ']' {
					if increment > 0.1 {
						increment -= 0.1
					}
				}
				if ev.Key() == tcell.KeyUp {
					cube1.A += increment
				}
				if ev.Key() == tcell.KeyDown {
					cube1.A -= increment
				}
				if ev.Key() == tcell.KeyLeft {
					cube1.B += increment
				}
				if ev.Key() == tcell.KeyRight {
					cube1.B -= increment
				}
				if ev.Rune() == ' ' {
					autoSpin = !autoSpin
				}
			}
		}
	}()
	sphere := &Sphere{
		Radius:           46,
		DistanceFromCam:  89,
		HorizontalOffset: 0,
		VerticalOffset:   0,
		K1:               20,
	}

	for {
		select {
		case <-quit:
			return
		default:
			width, height = s.Size()
			zBuffer = make([]float64, width*height)
			buffer = make([]rune, width*height)
			colors = make([]tcell.Color, width*height)

			// Clear buffers
			for i := 0; i < width*height; i++ {
				buffer[i] = backgroundASCIICode
				zBuffer[i] = 0
			}

			// sphere.Generate(s)
			cube1.Generate(s) // Call Generate on each cube
			drawbuffers(s)

			s.Show()
			if autoSpin {
				cube1.A += 0.04
				cube1.B += 0.03
				cube1.C += 0.02
			}

			sphere.A += 0.04
			sphere.B += 0.03
			sphere.C += 0.02

			// Update distanceFromCam and direction for each cube individually
			if cube1.DistanceFromCam <= 50 {
				cube1.Direction = false
			} else if cube1.DistanceFromCam >= 200 {
				cube1.Direction = true
			}
			if cube1.Direction {
				cube1.DistanceFromCam -= 0 // Adjust as needed
			} else {
				cube1.DistanceFromCam += 0 // Adjust as needed
			}

			time.Sleep(time.Millisecond * 16)
		}
	}
}

