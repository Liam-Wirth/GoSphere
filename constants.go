package main

import (
	"math"

	"github.com/gdamore/tcell/v2"
)

var (
	width, height       int
	zBuffer             []float64
	buffer              []rune
	colors              []tcell.Color
	backgroundASCIICode rune    = ' '
	increment           float64 = 0.1
	autoSpin            bool    = true
	pi_16               float64 = math.Pi / 16
	maxWidth            int     = 500 // Adjust as per maximum expected width
	maxHeight           int     = 200 // Adjust as per maximum expected height
)

var shading = []rune{'.', '-', '/', '=', '+', '*', '#', '%', '@'}

var grayScale = []tcell.Color{
	tcell.ColorBlack,
	tcell.NewRGBColor(10, 10, 10),
	tcell.NewRGBColor(20, 20, 20),
	tcell.NewRGBColor(30, 30, 30),
	tcell.NewRGBColor(40, 40, 40),
	tcell.NewRGBColor(50, 50, 50),
	tcell.NewRGBColor(60, 60, 60),
	tcell.NewRGBColor(70, 70, 70),
	tcell.NewRGBColor(80, 80, 80),
	tcell.NewRGBColor(90, 90, 90),
	tcell.NewRGBColor(100, 100, 100),
}

