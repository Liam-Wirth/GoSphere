package main

import (
	"math"
	"runtime"
	"sync"

	"github.com/gdamore/tcell/v2"
)

type Sphere struct {
	Radius           float64
	DistanceFromCam  float64
	HorizontalOffset float64
	VerticalOffset   float64
	K1               float64
	K2               float64
	RotationSpeed    float64
	A, B, C          float64
	Points           []SpherePoint
	Resolution       float64
	ColorFunction    func(phi, theta float64) tcell.Color
}

type SpherePoint struct {
	X, Y, Z float64
	Color   tcell.Color
}

func (sphere *Sphere) BuildSurface() {
	numTheta := int(2*math.Pi/sphere.Resolution) + 1
	numPhi := int(math.Pi/sphere.Resolution) + 1
	totalPoints := numTheta * numPhi
	sphere.Points = make([]SpherePoint, 0, totalPoints)

	for theta := 0.0; theta <= 2*math.Pi; theta += sphere.Resolution {
		sinTheta := math.Sin(theta)
		cosTheta := math.Cos(theta)
		for phi := 0.0; phi <= math.Pi; phi += sphere.Resolution {
			sinPhi := math.Sin(phi)
			cosPhi := math.Cos(phi)

			x := sphere.Radius * sinPhi * cosTheta
			y := sphere.Radius * sinPhi * sinTheta
			z := sphere.Radius * cosPhi

			color := sphere.ColorFunction(phi, theta)

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

func (sphere *Sphere) Generate(screen tcell.Screen) {
	aspectRatio := sphere.K2
	K2 := sphere.K1 / aspectRatio

	// Precompute rotation cosines and sines
	cosA, sinA := math.Cos(sphere.A), math.Sin(sphere.A)
	cosB, sinB := math.Cos(sphere.B), math.Sin(sphere.B)
	cosC, sinC := math.Cos(sphere.C), math.Sin(sphere.C)

	// Use concurrency to process points
	numWorkers := runtime.NumCPU()
	pointChunks := chunkPoints(sphere.Points, numWorkers)
	var wg sync.WaitGroup

	for _, chunk := range pointChunks {
		wg.Add(1)
		go func(points []SpherePoint) {
			defer wg.Done()
			for _, point := range points {
				x := point.X
				y := point.Y
				z := point.Z

				// Rotate the point
				rotatedPoint := rotatePoint(x, y, z, cosA, sinA, cosB, sinB, cosC, sinC)

				// Compute normal vector before rotation (for lighting)
				nx := x / sphere.Radius
				ny := y / sphere.Radius
				nz := z / sphere.Radius
				normal := Vec3{nx, ny, nz}

				// Rotate the normal vector
				rotatedNormal := rotateNormal(normal, cosA, sinA, cosB, sinB, cosC, sinC)

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

				zProj := rotatedPoint.Z + sphere.DistanceFromCam
				if zProj == 0 {
					continue // Avoid division by zero
				}
				ooz := 1 / zProj
				xp := int(float64(width)/2 + sphere.HorizontalOffset + sphere.K1*ooz*rotatedPoint.X)
				yp := int(float64(height)/2 + sphere.VerticalOffset - K2*ooz*rotatedPoint.Y)

				if xp >= 0 && xp < width && yp >= 0 && yp < height {
					idx := xp + yp*width
					if ooz > zBuffer[idx] {
						zBuffer[idx] = ooz

						// Clamp intensity to [0, 1]
						intensity := dotProduct
						if intensity > 1 {
							intensity = 1
						}

						// Adjust the point's color with the lighting intensity
						color := adjustColorIntensity(point.Color, intensity)

						buffer[idx] = '█'
						colors[idx] = color
					}
				}
			}
		}(chunk)
	}
	wg.Wait()
}

func chunkPoints(points []SpherePoint, numChunks int) [][]SpherePoint {
	chunkSize := (len(points) + numChunks - 1) / numChunks
	var chunks [][]SpherePoint
	for i := 0; i < len(points); i += chunkSize {
		end := i + chunkSize
		if end > len(points) {
			end = len(points)
		}
		chunks = append(chunks, points[i:end])
	}
	return chunks
}

func rotatePoint(x, y, z, cosA, sinA, cosB, sinB, cosC, sinC float64) Vec3 {
	// Apply Z rotation
	newX := x*cosC - y*sinC
	newY := x*sinC + y*cosC
	newZ := z

	// Apply Y rotation
	x = newX*cosB + newZ*sinB
	z = newZ*cosB - newX*sinB

	// Apply X rotation
	y = newY*cosA - z*sinA
	z = newY*sinA + z*cosA

	return Vec3{x, y, z}
}

func rotateNormal(v Vec3, cosA, sinA, cosB, sinB, cosC, sinC float64) Vec3 {
	// Apply Z rotation
	newX := v.X*cosC - v.Y*sinC
	newY := v.X*sinC + v.Y*cosC
	newZ := v.Z

	// Apply Y rotation
	x := newX*cosB + newZ*sinB
	z := newZ*cosB - newX*sinB

	// Apply X rotation
	y := newY*cosA - z*sinA
	z = newY*sinA + z*cosA

	return Vec3{x, y, z}
}
func checkerboardColorFunction(phi, theta float64) tcell.Color {
	u := int(phi / pi_16)
	v := int(theta / pi_16)

	if (u+v)%2 == 0 {
		return tcell.ColorPurple
	} else {
		return tcell.ColorWhite
	}
}
