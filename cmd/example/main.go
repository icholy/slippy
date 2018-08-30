package main

import (
	"log"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"

	"github.com/icholy/slippy"
	"github.com/icholy/slippy/tiles"
	"github.com/icholy/slippy/util"
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

	zoom := 10

	origin := tiles.Vec(43.174366, -79.231511, zoom)
	frame := pixel.R(0, 0, 450, 500).Moved(origin)

	tt := slippy.RectTiles(frame, zoom)
	for i := range tt {
		if err := tt[i].Fetch(); err != nil {
			return err
		}
	}

	camera := pixel.ZV.Sub(origin).Add(pixel.V(200, 200))
	win.SetMatrix(pixel.IM.Moved(camera))

	win.Clear(colornames.Skyblue)

	for _, t := range tt {
		t.Draw(win)
		util.DrawRect(win, t.Rect(), colornames.Black)
	}

	util.DrawRect(win, frame, colornames.Blue)
	util.DrawVec(win, origin)

	for !win.Closed() {
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
