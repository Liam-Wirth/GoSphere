// TODO: Split into smaller files, and start to flesh this out instead into a console based 3d renderer
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

type SpherePoint struct {
	X, Y, Z float64     // Original coordinates
	Color   tcell.Color // Assigned color
	// symb    rune
}

type Sphere struct {
	Radius           float64
	DistanceFromCam  float64
	HorizontalOffset float64
	VerticalOffset   float64
	K1               float64
	K2               float64
	RotationSpeed    float64
	A, B, C          float64 // Rotation angles
	//below might be the worst way possible to do this, but fuck it, I need to bake some sorta "texture"/pattern onto the sphere, and have that pattern stay constant (as in on the same point in the sphere) even when the sphere rotates
	// i.e if I have pointA, and it's given colorA, and it's rotated such that it goes in the place where pointB was, instead of assuming the color of pointB, it keeps it's old color

	Points []SpherePoint // Surface points

	Resolution    float64                              // Sampling resolution
	ColorFunction func(phi, theta float64) tcell.Color // Function to assign colors
}

// Lighting to give shading effects (simplified lighting model)
var shading = []rune{'.', '-', '/', '=', '+', '*', '#', '%', '@'}
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

func normalize(v Vec3) Vec3 {
	length := math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
	return Vec3{v.X / length, v.Y / length, v.Z / length}
}

func checkerboardColorFunction(phi, theta float64) tcell.Color {
	u := int(phi/(math.Pi/8))
	v := int(theta/(math.Pi/8))

	if (u+v)%2 == 0 {
		return tcell.ColorPurple.TrueColor()
	} else {
		return tcell.ColorWhite.TrueColor()
	}
}

// FIX: The logic for adjusting intensity/brightness of the color is not working properly.
func adjustColorIntensity(color tcell.Color, intensity float64) tcell.Color {
	// Decompose the color into RGB components
	r, g, b := color.RGB()
	//intensity is clamped from 0 -> 1
	intensity = 1 - intensity
	intensity *= 0.1
	// ir := int32(math.Floor(float64(r) * intensity))
	// ig := int32(math.Floor(float64(g) * intensity))
	// ib := int32(math.Floor(float64(b) * intensity))
	ir := r
	ig := g
	ib := b
	return tcell.NewRGBColor(ir, ig, ib)
}

func (sphere *Sphere) BuildSurface() {
	sphere.Points = []SpherePoint{}

	for theta := 0.0; theta < 2*math.Pi; theta += sphere.Resolution {
		for phi := 0.0; phi < math.Pi; phi += sphere.Resolution {
			// Calculate the 3D coordinates
			x := sphere.Radius * math.Sin(phi) * math.Cos(theta)
			y := sphere.Radius * math.Sin(phi) * math.Sin(theta)
			z := sphere.Radius * math.Cos(phi)

			// Assign color based on phi and theta
			color := sphere.ColorFunction(phi, theta)

			// Create the SpherePoint
			point := SpherePoint{
				X:     x,
				Y:     y,
				Z:     z,
				Color: color,
			}

			sphere.Points = append(sphere.Points, point)
		}
	}
}

// Sphere Generate projects points on the surface of the sphere
// FIX: don't work :(
func (sphere *Sphere) Generate(screen tcell.Screen) {
	aspectRatio := sphere.K2
	K2 := sphere.K1 / aspectRatio

	for _, point := range sphere.Points {
		// Original coordinates
		x := point.X
		y := point.Y
		z := point.Z

		// Rotate the point
		rotatedPoint := rotateX(Vec3{x, y, z}, sphere.A)
		rotatedPoint = rotateY(rotatedPoint, sphere.B)
		rotatedPoint = rotateZ(rotatedPoint, sphere.C)

		// Compute normal vector before rotation (for lighting)
		nx := x / sphere.Radius
		ny := y / sphere.Radius
		nz := z / sphere.Radius
		normal := Vec3{nx, ny, nz}

		// Rotate the normal vector
		rotatedNormal := rotateX(normal, sphere.A)
		rotatedNormal = rotateY(rotatedNormal, sphere.B)
		rotatedNormal = rotateZ(rotatedNormal, sphere.C)

		// Light direction from camera to point
		cameraPos := Vec3{0, 0, -sphere.DistanceFromCam}
		lightDir := Vec3{
			X: cameraPos.X - rotatedPoint.X,
			Y: cameraPos.Y - rotatedPoint.Y,
			Z: cameraPos.Z - rotatedPoint.Z,
		}
		lightDir = normalize(lightDir)

		// Compute the dot product
		dotProduct := rotatedNormal.X*lightDir.X + rotatedNormal.Y*lightDir.Y + rotatedNormal.Z*lightDir.Z

		// Back-face culling
		if dotProduct < 0 {
			continue
		}

		// Clamp dotProduct to [0, 1]
		if dotProduct > 1 {
			dotProduct = 1
		}

		// Project 3D point to 2D
		zProj := rotatedPoint.Z + sphere.DistanceFromCam
		if zProj == 0 {
			continue // Avoid division by zero
		}
		ooz := 1 / zProj
		xp := int(math.Round(float64(width)/2 + sphere.HorizontalOffset + sphere.K1*ooz*rotatedPoint.X))
		yp := int(math.Round(float64(height)/2 + sphere.VerticalOffset - K2*ooz*rotatedPoint.Y))

		if xp >= 0 && xp < width && yp >= 0 && yp < height {
			idx := xp + yp*width
			if ooz > zBuffer[idx] {
				zBuffer[idx] = ooz

				// Clamp intensity to [0, 1]
				intensity := dotProduct
				if intensity < 0 {
					intensity = 0
				} else if intensity > 1 {
					intensity = 1
				}

				// Adjust the point's color with the lighting intensity
				color := adjustColorIntensity(point.Color, intensity)

				buffer[idx] = '█'
				colors[idx] = color
			}
		}
	}
}

func intensityToColor(intensity float64) tcell.Color {
	// Ensure intensity is clamped between 0 and 1
	if intensity < 0 {
		intensity = 0
	} else if intensity > 1 {
		intensity = 1
	}

	// Define start and end colors for the gradient
	startColor := [3]uint8{0, 0, 255} // Blue
	endColor := [3]uint8{255, 0, 0}   // Red

	// Interpolate between the start and end colors
	r := uint8(float64(startColor[0])*(1-intensity) + float64(endColor[0])*intensity)
	g := uint8(float64(startColor[1])*(1-intensity) + float64(endColor[1])*intensity)
	b := uint8(float64(startColor[2])*(1-intensity) + float64(endColor[2])*intensity)

	return tcell.NewRGBColor(int32(r), int32(g), int32(b))
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

	// Initialize the cube (if needed)
	cube1 := &Cube{
		CubeWidth:        20,
		DistanceFromCam:  100,
		HorizontalOffset: -40, // Left side of the screen
		K1:               40,
		IncrementSpeed:   0.3,
		CubeFaces:        defaultCubeFaces,
	}

	// Initialize the sphere
	sphere := &Sphere{
		Radius:           36,
		DistanceFromCam:  54,
		HorizontalOffset: 0,
		VerticalOffset:   0,
		K1:               20,
		K2:               1.5,  // Aspect ratio
		Resolution:       0.02, // Adjust for desired quality
		ColorFunction:    checkerboardColorFunction,
	}
	sphere.BuildSurface()

	// Event handling goroutine
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

	// Main rendering loop
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
				zBuffer[i] = -math.MaxFloat64 // Initialize to very small value
				colors[i] = tcell.ColorBlack  // Initialize background color
			}

			// Generate sphere
			sphere.Generate(s)

			// Optionally generate cube
			// cube1.Generate(s)

			// Draw buffers to the screen
			drawbuffers(s)

			// Show the updated screen
			s.Show()

			// Apply automatic rotation if enabled

			// Update sphere rotation angles
			sphere.B += 0.03

			// Update cube distance (if needed)
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
