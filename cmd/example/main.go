package main

import (
	"fmt"
	"log"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"

	"github.com/icholy/slippy"
)

func run() error {
	cfg := pixelgl.WindowConfig{
		Title:  "Pixel Rocks!",
		Bounds: pixel.R(0, 0, 1024, 768),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		return err
	}

	padding := pixel.V(10, 10)
	m := slippy.New(slippy.Options{
		Lat:  43.174366,
		Lon:  -79.231511,
		Zoom: 10,
		Bounds: pixel.Rect{
			Min: win.Bounds().Min.Add(padding),
			Max: win.Bounds().Max.Sub(padding),
		},
	})

	if err := m.Fetch(); err != nil {
		return err
	}

	m.Draw(win, pixel.IM)

	for !win.Closed() {

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
