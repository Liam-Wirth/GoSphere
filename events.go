package main

import (
	"github.com/gdamore/tcell/v2"
)

func handleEvents(s tcell.Screen, quit chan struct{}, cube1 *Cube, sphere *Sphere) {
	for {
		ev := s.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyCtrlC:
				fallthrough
			case tcell.KeyESC:
				fallthrough
			case tcell.KeyCtrlQ:
				close(quit)
				return
			case tcell.KeyUp:
				cube1.A += increment
			case tcell.KeyDown:
				cube1.A -= increment
			case tcell.KeyLeft:
				cube1.B += increment
			case tcell.KeyRight:
				cube1.B -= increment
			}

			switch ev.Rune() {
			case 'q':
				close(quit)
				return
			case '[':
				if increment < 1 {
					increment += 0.1
				}
			case ']':
				if increment > 0.1 {
					increment -= 0.1
				}
			case ' ':
				autoSpin = !autoSpin
                                sphere.BuildSurface() //redraw
			case '1':
				sphere.Resolution -= 0.01
			case '2':
				sphere.Resolution += 0.01
			case 'w':
				sphere.DistanceFromCam += 0.001
			case 's':
				sphere.DistanceFromCam -= 0.001

			case 'm':
				// Open the menu
				// showMenu(s, cube1, sphere)
			}
		}
	}
}
