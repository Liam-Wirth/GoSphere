package main

import "github.com/gdamore/tcell/v2"

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

type CubeFace struct {
	Char  rune
	Color tcell.Color
}

var defaultCubeFaces = []CubeFace{
	{'$', tcell.ColorRed},
	{'#', tcell.ColorGreen},
	{'@', tcell.ColorYellow},
	{'&', tcell.ColorBlue},
	{'%', tcell.ColorDarkCyan},
	{'|', tcell.ColorPurple},
}

func (cube *Cube) Generate(screen tcell.Screen) {
	// Precompute rotation angles
	rotations := Vec3{cube.A, cube.B, cube.C}

	// Cache repeated calculations
	halfWidth := cube.CubeWidth
	increment := cube.IncrementSpeed

	// Optimize loops by reducing nesting and precalculating values
	for i := 0; i < len(cube.CubeFaces); i++ {
		face := cube.CubeFaces[i]
		for cubeX := -halfWidth; cubeX < halfWidth; cubeX += increment {
			for cubeY := -halfWidth; cubeY < halfWidth; cubeY += increment {
				cubePoint := getCubePoint(i, cubeX, cubeY, cube)
				calculateForSurface(cubePoint, face.Char, face.Color, rotations, cube)
			}
		}
	}
}

func getCubePoint(faceIndex int, cubeX, cubeY float64, cube *Cube) Vec3 {
	switch faceIndex {
	case 0:
		return Vec3{cubeX, cubeY, -cube.CubeWidth}
	case 1:
		return Vec3{cube.CubeWidth, cubeY, cubeX}
	case 2:
		return Vec3{-cube.CubeWidth, cubeY, -cubeX}
	case 3:
		return Vec3{-cubeX, cubeY, cube.CubeWidth}
	case 4:
		return Vec3{cubeX, -cube.CubeWidth, -cubeY}
	case 5:
		return Vec3{cubeX, cube.CubeWidth, cubeY}
	}
	return Vec3{}
}

func calculateForSurface(cubePoint Vec3, ch rune, color tcell.Color, rotations Vec3, cube *Cube) {
	rotatedPoint := calcNewPoint(cubePoint, rotations)
	z := rotatedPoint.Z + float64(cube.DistanceFromCam)
	if z == 0 {
		return
	}
	ooz := 1 / z
	xp := int(float64(width)/2 + cube.HorizontalOffset + cube.K1*ooz*rotatedPoint.X*2)
	yp := int(float64(height)/2 + cube.K1*ooz*rotatedPoint.Y)
	if xp >= 0 && xp < width && yp >= 0 && yp < height {
		idx := xp + yp*width
		if ooz > zBuffer[idx] { //is this not supposed to be cullign?!?!?!?gg
			zBuffer[idx] = ooz
			buffer[idx] = ch
			colors[idx] = color
		}
	}
}

