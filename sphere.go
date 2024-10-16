// TODO: generalize a lot of this code/apply it to lighting calculations for other shapes
package main

import (
	// "fmt"
	"github.com/gdamore/tcell/v2"
	"math"
	"runtime"
	"sync"
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
	cameraPos        Vec3
	lightPos         Vec3
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
	sphere.cameraPos = Vec3{0, 0, -sphere.DistanceFromCam}

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
	// TODO: decouple goroutines from sphere rendering and move it to general rendering function, right now cause cube is single threaded, rendering gets blocked if I do 2 shapes at a time
	runtime.GOMAXPROCS(runtime.NumCPU()) // Ensure all CPUs are used
	sphere.lightPos = Vec3{X: 1, Y: 1, Z: -2}

	aspectRatio := sphere.K2
	K2 := sphere.K1 / aspectRatio

	cosA, sinA := math.Cos(sphere.A), math.Sin(sphere.A)
	cosB, sinB := math.Cos(sphere.B), math.Sin(sphere.B)
	cosC, sinC := math.Cos(sphere.C), math.Sin(sphere.C)

	numWorkers := runtime.NumCPU()
	pointChunks := chunkPoints(sphere.Points, numWorkers)
	var wg sync.WaitGroup

	for workerID, chunk := range pointChunks {
		wg.Add(1)
		go func(points []SpherePoint, workerID int) {
			defer wg.Done()
			for _, point := range points {
				x := point.X
				y := point.Y
				z := point.Z

				rotatedPoint := rotatePoint(x, y, z, cosA, sinA, cosB, sinB, cosC, sinC)

				nx := x / sphere.Radius
				ny := y / sphere.Radius
				nz := z / sphere.Radius
				normal := Vec3{nx, ny, nz}

				rotatedNormal := rotateNormal(normal, cosA, sinA, cosB, sinB, cosC, sinC)

				//  Dynamic lighting, I could theoretically map the x and y values to cursor position, and have it like that
				// kinda cheating by zooming out lmao
				// TODO: Uncouple light dir from camera pos in culling
				lightDir := Vec3{
					X: sphere.lightPos.X - rotatedPoint.X,
					Y: sphere.lightPos.Y - rotatedPoint.Y,
					Z: sphere.lightPos.Z - rotatedPoint.Z,
				}
				lightDir = normalize(lightDir)

				// this is because the dot product is technically the light direction, they are coupled currently because the code is written such that the light is "supposed" to be coming from the camera
				// cause cpu i'm minimizing calculations, so if I move the light dir away from the camera pos too much, it fucks with culling/rendering
				// hence why it looks all weird rn
				// dotProduct := rotatedNormal.X*lightDir.X + rotatedNormal.Y*lightDir.Y + rotatedNormal.Z*lightDir.Z
				dotProduct := rotatedNormal.X*lightDir.X + rotatedNormal.Y*lightDir.Y + rotatedNormal.Z*lightDir.Z
				if dotProduct < 0 {
					continue
				}

				zProj := rotatedPoint.Z + sphere.DistanceFromCam
				if zProj == 0 {
					continue
				}
				ooz := 1 / zProj
				xp := int(float64(width)/2 + sphere.HorizontalOffset + sphere.K1*ooz*rotatedPoint.X)
				yp := int(float64(height)/2 + sphere.VerticalOffset - K2*ooz*rotatedPoint.Y)

				if xp >= 0 && xp < width && yp >= 0 && yp < height {
					idx := xp + yp*width
					if ooz > zBuffer[idx] {
						zBuffer[idx] = ooz

						intensity := dotProduct
						color := adjustColorIntensity(point.Color, intensity)
						buffer[idx] = '@'
						// buffer[idx] = '█'
						colors[idx] = color
					}
				}
			}
		}(chunk, workerID)
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
		return tcell.ColorFuchsia
	} else {
		return tcell.ColorBlack
	}

}
func stripeColorFunction(phi, theta float64) tcell.Color {
	if int(theta/(math.Pi/4))%2 == 0 {
		return tcell.ColorRed
	} else {
		return tcell.ColorBlue
	}
}
func gradientColorFunction(phi, theta float64) tcell.Color {
	// Map phi from [0, π] to [0, 1]
	t := phi / math.Pi
	r := int32(255 * t)
	g := int32(255 * (1 - t))
	b := int32(128)
	return tcell.NewRGBColor(r, g, b)
}
