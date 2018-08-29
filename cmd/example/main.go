package main

import (
	_ "image/png"
	"log"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"

	"github.com/icholy/tiles"
)

func loadTiles(r pixel.Rect, zoom int) ([]tiles.Tile, error) {
	var tt []tiles.Tile
	for _, tile := range tiles.Fill(r, zoom) {
		t, err := tiles.LoadTile(tile)
		if err != nil {
			return nil, err
		}
		tt = append(tt, t)
	}
	return tt, nil
}

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

	tt, err := loadTiles(frame, zoom)
	if err != nil {
		return err
	}

	camera := pixel.ZV.Sub(origin).Add(pixel.V(200, 200))
	win.SetMatrix(pixel.IM.Moved(camera))

	win.Clear(colornames.Skyblue)

	for _, t := range tt {
		t.Draw(win)
		tiles.DrawRect(win, t.Rect(), colornames.Black)
	}

	tiles.DrawRect(win, frame, colornames.Blue)
	tiles.DrawVec(win, origin)

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
