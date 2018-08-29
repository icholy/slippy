package main

import (
	"errors"
	"fmt"
	"image"
	_ "image/png"
	"log"
	"math/rand"
	"net/http"

	"github.com/buckhx/tiles"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

func TilePictureData(t tiles.Tile) (*pixel.PictureData, error) {
	shards := []string{"a", "b", "c"}
	url := fmt.Sprintf(
		"http://%[1]s.tile.openstreetmap.org/%[2]d/%[3]d/%[4]d.png",
		shards[rand.Intn(len(shards))], t.Z, t.X, t.Y,
	)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Go-Tile-Pixel-Experiment")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}
	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}

// CoordinateVec return the vector for the specified WGS84 coordiantes
func CoordinateVec(latitude, longitude float64, zoom int) pixel.Vec {
	c := tiles.FromCoordinate(latitude, longitude, zoom)
	return PixelVec(c.ToPixel())
}

func PixelVec(p tiles.Pixel) pixel.Vec {
	return pixel.V(
		float64(p.X),
		float64(p.Y),
	)
}

// TileVec returns a vector for the bottom left corner
// of the tile
func TileVec(t tiles.Tile) pixel.Vec {
	p := t.ToPixel()
	return PixelVec(tiles.Pixel{
		X: p.X,
		Y: p.Y - tiles.TileSize,
	})
}

// TileRect returns a rectangle of the tile
func TileRect(t tiles.Tile) pixel.Rect {
	v := TileVec(t)
	return pixel.R(
		v.X,
		v.Y,
		v.X+float64(tiles.TileSize),
		v.Y+float64(tiles.TileSize),
	)
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

	win.Clear(colornames.Skyblue)

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
