package main

import (
	"fmt"
	"math"
	// "sort"
	"time"
)

type Vec3 struct {
	X, Y, Z float64
}

var (
	// these are angles in radians
	A, B, C float64

	cubeWidth           float64   = 20
	width, height       int       = 203, 74
	zBuffer             []float64 = make([]float64, width*height)
	buffer              []byte    = make([]byte, width*height)
	backgroundASCIICode byte      = ' '
	distanceFromCam     int       = 100
	horizontalOffset    float64
	K1                  float64 = 40

	incrementSpeed float64 = 0.6

	x, y, z         float64
	point, newPoint Vec3

	ooz float64

	xp, yp int

	idx int
	//true for forward, false for backward
	direction bool = true
)

// Implemented based on rotation matrices I learned about and figured out on scratch paper.
func rotateX(p Vec3, theta float64) Vec3 {
	newY := (p.Y * math.Cos(theta)) - (p.Z * math.Sin(theta))
	newZ := (p.Y * math.Sin(theta)) + (p.Z * math.Cos(theta))

	return Vec3{
		p.X, newY, newZ,
	}
}

func rotateY(p Vec3, theta float64) Vec3 {
	newX := (p.X * math.Cos(theta)) + (p.Z * math.Sin(theta))
	newZ := (p.Z * math.Cos(theta)) - (p.X * math.Sin(theta))

	return Vec3{
		newX, p.Y, newZ,
	}
}

func rotateZ(p Vec3, theta float64) Vec3 {
	newX := (p.X * math.Cos(theta)) - (p.Y * math.Sin(theta))
	newY := (p.X * math.Sin(theta)) + (p.Y * math.Cos(theta))

	return Vec3{
		newX, newY, p.Z,
	}
}

func calcNewPoint(p Vec3, rotations Vec3) Vec3 {

	newPoint := rotateZ(p, rotations.Z)
	newPoint = rotateY(newPoint, rotations.Y)
	newPoint = rotateX(newPoint, rotations.X)
	return newPoint
}

func calculateForSurface(cubePoint Vec3, ch byte) {
	rotatedPoint := calcNewPoint(cubePoint, Vec3{A, B, C})
	z := rotatedPoint.Z + float64(distanceFromCam)
	ooz = 1 / z

	xp = int(float64(width)/2 + horizontalOffset + K1*ooz*rotatedPoint.X*2)
	yp = int(float64(height)/2 + K1*ooz*rotatedPoint.Y)

	idx = xp + yp*width
	if idx >= 0 && idx < width*height {
		if ooz > zBuffer[idx] {
			zBuffer[idx] = ooz
			buffer[idx] = ch
		}
	}
}

func generateCube() {
	horizontalOffset = -2 * cubeWidth
	for cubeX := -cubeWidth; cubeX < cubeWidth; cubeX += incrementSpeed {
		for cubeY := -cubeWidth; cubeY < cubeWidth; cubeY += incrementSpeed {
			cubePoint := Vec3{cubeX, cubeY, -cubeWidth}
			calculateForSurface(cubePoint, '@')

			cubePoint = Vec3{cubeWidth, cubeY, cubeX}
			calculateForSurface(cubePoint, '$')

			cubePoint = Vec3{-cubeWidth, cubeY, -cubeX}
			calculateForSurface(cubePoint, '~')

			cubePoint = Vec3{-cubeX, cubeY, cubeWidth}
			calculateForSurface(cubePoint, '#')

			cubePoint = Vec3{cubeX, -cubeWidth, -cubeY}
			calculateForSurface(cubePoint, ';')

			cubePoint = Vec3{cubeX, cubeWidth, cubeY}
			calculateForSurface(cubePoint, '+')
		}
	}
}

func main() {
	zBuffer = make([]float64, width*height)
	buffer = make([]byte, width*height)

	// Clear the console
	fmt.Print("\033[H\033[2J")

	for {
		// Clear buffers
		for i := 0; i < width*height; i++ {
			buffer[i] = backgroundASCIICode
			zBuffer[i] = 0
		}

		// Generate the cube
		generateCube()
		// Print the buffer
		fmt.Print("\x1b[H")
		for k := 0; k < width*height; k++ {
			if k%width == 0 {
				fmt.Println()
			} else {
				fmt.Print(string(buffer[k]))
			}
		}

		A += 0.05
		B += 0.05
		C += 0.01
		// Adjust distanceFromCam based on direction
		if distanceFromCam <= 50 {
			direction = false // Start moving farther away
		} else if distanceFromCam >= 200 {
			direction = true // Start moving closer
		}
		if direction {
			distanceFromCam -= 3
		} else {
			distanceFromCam += 3
		}
		time.Sleep(time.Millisecond * 10)
	}
}
