package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"

	"github.com/icholy/slippy"
)

func run() error {

	// create the window
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

	for !win.Closed() {

		m.FetchAsync()

		// draw the map
		m.Draw(win, pixel.IM)

		switch {
		case win.Pressed(pixelgl.KeyLeft):
			m.SetCenterVec(m.CenterVec().Sub(pixel.V(10, 0)))
		case win.Pressed(pixelgl.KeyRight):
			m.SetCenterVec(m.CenterVec().Add(pixel.V(10, 0)))
		case win.Pressed(pixelgl.KeyUp):
			m.SetCenterVec(m.CenterVec().Add(pixel.V(0, 10)))
		case win.Pressed(pixelgl.KeyDown):
			m.SetCenterVec(m.CenterVec().Sub(pixel.V(0, 10)))
		case win.JustPressed(pixelgl.KeyEqual):
			m.SetZoom(m.Zoom() + 1)
		case win.JustPressed(pixelgl.KeyMinus):
			m.SetZoom(m.Zoom() - 1)
		}

		// print clicked location coordinates
		if win.Pressed(pixelgl.MouseButtonLeft) {
			coord := m.Coord(win.MousePosition())
			fmt.Println("Clicked", coord)
			drawVec(win, m.Vec(coord))
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

func drawVec(t pixel.Target, v pixel.Vec) {
	imd := imdraw.New(nil)
	imd.Color = randColor()
	imd.Push(v)
	imd.Circle(5, 2)
	imd.Draw(t)
}

func randColor() color.Color {
	names := colornames.Names
	name := names[rand.Intn(len(names))]
	return colornames.Map[name]
}
