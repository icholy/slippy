package main

import (
	"fmt"
	"log"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"

	"github.com/icholy/slippy"
)

func run() error {
	// Create the window
	cfg := pixelgl.WindowConfig{
		Title:  "OpenStreetMaps",
		Bounds: pixel.R(0, 0, 1024, 768),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		return err
	}

	// create the map
	padding := pixel.V(10, 10)
	m := slippy.New(slippy.Options{
		Zoom:   10,
		Center: slippy.ClippedCoords(43.174366, -79.231511),
		Bounds: pixel.Rect{
			Min: win.Bounds().Min.Add(padding),
			Max: win.Bounds().Max.Sub(padding),
		},
	})

	// fetch the tiles
	if err := m.Fetch(); err != nil {
		return err
	}

	// draw the map
	m.Draw(win, pixel.IM)

	for !win.Closed() {
		// print clicked location coordinates
		if win.JustPressed(pixelgl.MouseButtonLeft) {
			fmt.Println("Clicked", m.Coord(win.MousePosition()))
		}
		win.Update()
	}
	return nil
}

func main() {
	pixelgl.Run(func() {
		if err := run(); err != nil {
			log.Fatal(err)
		}
	})

}
