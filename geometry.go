package main

import (
	"math"

	"github.com/gdamore/tcell/v2"
)

type Vec3 struct {
	X, Y, Z float64
}

func rotateX(p Vec3, cosTheta, sinTheta float64) Vec3 {
	newY := p.Y*cosTheta - p.Z*sinTheta
	newZ := p.Y*sinTheta + p.Z*cosTheta
	return Vec3{p.X, newY, newZ}
}

func rotateY(p Vec3, cosTheta, sinTheta float64) Vec3 {
	newX := p.X*cosTheta + p.Z*sinTheta
	newZ := p.Z*cosTheta - p.X*sinTheta
	return Vec3{newX, p.Y, newZ}
}

func rotateZ(p Vec3, cosTheta, sinTheta float64) Vec3 {
	newX := p.X*cosTheta - p.Y*sinTheta
	newY := p.X*sinTheta + p.Y*cosTheta
	return Vec3{newX, newY, p.Z}
}

func calcNewPoint(p Vec3, rotations Vec3) Vec3 {
	cosA, sinA := math.Cos(rotations.X), math.Sin(rotations.X)
	cosB, sinB := math.Cos(rotations.Y), math.Sin(rotations.Y)
	cosC, sinC := math.Cos(rotations.Z), math.Sin(rotations.Z)

	newPoint := rotateZ(p, cosC, sinC)
	newPoint = rotateY(newPoint, cosB, sinB)
	newPoint = rotateX(newPoint, cosA, sinA)
	return newPoint
}

func normalize(v Vec3) Vec3 {
	length := math.Hypot(math.Hypot(v.X, v.Y), v.Z)
	if length == 0 {
		return Vec3{0, 0, 0}
	}
	return Vec3{v.X / length, v.Y / length, v.Z / length}
}

type Model interface {
	Generate(screen tcell.Screen)
	BuildSurface()
}

