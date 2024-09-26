package main

import (
	"math"
	"time"

	"github.com/gdamore/tcell/v2"
)

type Vec3 struct {
	X, Y, Z float64
}

var (
	A, B, C             float64
	cubeWidth           float64 = 20
	width, height       int     = 203, 74
	zBuffer             []float64
	buffer              []rune
	colors              []tcell.Color
	backgroundASCIICode rune = ' '
	distanceFromCam     int  = 71
	horizontalOffset    float64
	K1                  float64 = 40
	incrementSpeed      float64 = 0.6
	direction           bool    = true
)

// CubeFace holds the character and color for each face of the cube
type CubeFace struct {
	Char  rune
	Color tcell.Color
}

// Define cube faces with associated colors and characters
var cubeFaces = []CubeFace{
	{'$', tcell.ColorRed},    // Front face
	{'#', tcell.ColorGreen},  // Right face
	{'@', tcell.ColorYellow}, // Back face
	{'&', tcell.ColorBlue},   // Left face
	{'G', tcell.ColorLightBlue},   // Bottom face
	{'?', tcell.ColorPurple},// Top face
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

func calculateForSurface(cubePoint Vec3, ch rune, color tcell.Color, screen tcell.Screen) {
	rotatedPoint := calcNewPoint(cubePoint, Vec3{A, B, C})
	z := rotatedPoint.Z + float64(distanceFromCam)
	ooz := 1 / z
	xp := int(float64(width)/2 + horizontalOffset + K1*ooz*rotatedPoint.X*2)
	yp := int(float64(height)/2 + K1*ooz*rotatedPoint.Y)
	idx := xp + yp*width
	if idx >= 0 && idx < width*height {
		if ooz > zBuffer[idx] {
			zBuffer[idx] = ooz
			buffer[idx] = ch
			colors[idx] = color
		}
	}
}

func generateCube(screen tcell.Screen) {
	horizontalOffset = -2 * cubeWidth
	for cubeX := -cubeWidth; cubeX < cubeWidth; cubeX += incrementSpeed {
		for cubeY := -cubeWidth; cubeY < cubeWidth; cubeY += incrementSpeed {
			// Loop through each face and calculate its projection with associated character and color
			for i, face := range cubeFaces {
				cubePoint := getCubePoint(i, cubeX, cubeY)
				calculateForSurface(cubePoint, face.Char, face.Color, screen)
			}
		}
	}
}

// Helper function to determine which point to project for each face
func getCubePoint(faceIndex int, cubeX, cubeY float64) Vec3 {
	switch faceIndex {
	case 0:
		return Vec3{cubeX, cubeY, -cubeWidth}
	case 1:
		return Vec3{cubeWidth, cubeY, cubeX}
	case 2:
		return Vec3{-cubeWidth, cubeY, -cubeX}
	case 3:
		return Vec3{-cubeX, cubeY, cubeWidth}
	case 4:
		return Vec3{cubeX, -cubeWidth, -cubeY}
	case 5:
		return Vec3{cubeX, cubeWidth, cubeY}
	}
	return Vec3{}
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

	go func() {
		for {
			ev := s.PollEvent()
			switch ev := ev.(type) {
			case *tcell.EventKey:
				if ev.Key() == tcell.KeyCtrlC || ev.Rune() == 'q' {
					close(quit)
					return
				}
			}
		}
	}()

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

			generateCube(s)

			// Render the cube with colors
			for y := 0; y < height; y++ {
				for x := 0; x < width; x++ {
					idx := y*width + x
					ch := buffer[idx]
					color := colors[idx]
					if ch != backgroundASCIICode {
						s.SetContent(x, y, ch, nil, tcell.StyleDefault.Foreground(color))
					} else {
						s.SetContent(x, y, ' ', nil, tcell.StyleDefault.Background(tcell.ColorDefault))
					}
				}
			}
			s.Show()

			A += 0.05
			B += 0.05
			C += 0.01

			if distanceFromCam <= 50 {
				direction = false
			} else if distanceFromCam >= 200 {
				direction = true
			}
			if direction {
				distanceFromCam -= 0
			} else {
				distanceFromCam += 0
			}

			time.Sleep(time.Millisecond * 16)
		}
	}
}

