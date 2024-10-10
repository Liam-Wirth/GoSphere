package main

import (
	"github.com/gdamore/tcell/v2"
)

// Adjust the color intensity based on the provided intensity value
func adjustColorIntensity(color tcell.Color, intensity float64) tcell.Color {
	r, g, b := color.RGB()
	ir := int32(float64(r) * intensity)
	ig := int32(float64(g) * intensity)
	ib := int32(float64(b) * intensity)
	return tcell.NewRGBColor(ir, ig, ib)
}

func intensityToColor(intensity float64) tcell.Color {
	if intensity < 0 {
		intensity = 0
	} else if intensity > 1 {
		intensity = 1
	}

	startColor := [3]uint8{0, 0, 255} // Blue
	endColor := [3]uint8{255, 0, 0}   // Red

	r := uint8(float64(startColor[0])*(1-intensity) + float64(endColor[0])*intensity)
	g := uint8(float64(startColor[1])*(1-intensity) + float64(endColor[1])*intensity)
	b := uint8(float64(startColor[2])*(1-intensity) + float64(endColor[2])*intensity)

	return tcell.NewRGBColor(int32(r), int32(g), int32(b))
}

func selectCharacterBasedOnSlope(dotProduct float64) rune {
    characters := []rune{ '▂', '▃', '▄', '▅', '▆', '▇', '█', '▉'}
    index := int(dotProduct * float64(len(characters)-1))
    return characters[index]
}


