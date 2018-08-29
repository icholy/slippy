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
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

func TilePictureData(t tiles.Tile) (*pixel.PictureData, error) {
	shards := []string{"a", "b", "c"}
	url := fmt.Sprintf(
		"http://%[1]s.tile.openstreetmap.org/%[2]d/%[3]d/%[4]d.png",
		shards[rand.Intn(len(shards))], t.Z, t.X, t.Y,
	)
	fmt.Println("URL", url)
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
func CoordinateVec(c tiles.Coordinate, zoom int) pixel.Vec {
	return PixelVec(c.ToPixel(zoom))
}

func VecTile(v pixel.Vec, zoom int) tiles.Tile {
	p := tiles.Pixel{
		X: int(v.X),
		Y: -int(v.Y),
		Z: zoom,
	}
	t, _ := p.ToTile()
	return t
}

func PixelVec(p tiles.Pixel) pixel.Vec {
	return pixel.V(
		float64(p.X),
		-float64(p.Y),
	)
}

// TileVec returns a vector for the bottom left corner
// of the tile
func TileVec(t tiles.Tile) pixel.Vec {
	p := t.ToPixel()
	return PixelVec(tiles.Pixel{
		X: p.X,
		Y: p.Y + tiles.TileSize,
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

func DrawTile(tg pixel.Target, t tiles.Tile, s *pixel.Sprite) {
	m := float64(tiles.TileSize) / 2
	v := TileVec(t).Add(pixel.V(m, m))
	s.Draw(tg, pixel.IM.Moved(v))
}

func DrawRect(tg pixel.Target, r pixel.Rect) {
	m := imdraw.New(nil)
	m.Color = colornames.Black
	m.Push(
		r.Min,
		pixel.V(r.Max.X, r.Min.Y),
		r.Max,
		pixel.V(r.Min.X, r.Max.Y),
	)
	m.Polygon(5)
	m.Draw(tg)
}

func Fill(r pixel.Rect, zoom int) []tiles.Tile {
	var (
		min = VecTile(r.Min, zoom)
		max = VecTile(r.Max, zoom)
		tt  []tiles.Tile
	)
	for x := min.X; x <= max.X; x++ {
		for y := min.Y; y <= max.Y; y++ {
			tt = append(tt, tiles.Tile{
				X: x,
				Y: y,
				Z: zoom,
			})
		}
	}
	return tt
}

func DrawVec(tg pixel.Target, v pixel.Vec) {
	m := imdraw.New(nil)
	m.Color = colornames.Blue
	m.Push(v)
	m.Circle(10, 5)
	m.Draw(tg)
}

type Tile struct {
	Tile   tiles.Tile
	Sprite *pixel.Sprite
}

func LoadTile(t tiles.Tile) (Tile, error) {
	pic, err := TilePictureData(t)
	if err != nil {
		return Tile{}, err
	}
	return Tile{
		Tile:   t,
		Sprite: pixel.NewSprite(pic, pic.Bounds()),
	}, nil
}

func (t Tile) Rect() pixel.Rect {
	return TileRect(t.Tile)
}

func (t Tile) Vec() pixel.Vec {
	return TileVec(t.Tile)
}

func (t Tile) Draw(tg pixel.Target) {
	DrawTile(tg, t.Tile, t.Sprite)
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
	coord := tiles.ClippedCoords(43.174366, -79.231511)
	origin := CoordinateVec(coord, zoom)

	t, err := LoadTile(VecTile(origin, zoom))
	if err != nil {
		return err
	}

	frame := pixel.R(0, 0, 400, 400).Moved(origin)

	camera := pixel.ZV.Sub(origin).Add(pixel.V(100, 100))
	win.SetMatrix(pixel.IM.Moved(camera))

	win.Clear(colornames.Skyblue)

	t.Draw(win)

	DrawRect(win, frame)
	DrawVec(win, origin)
	DrawVec(win, t.Vec())
	DrawRect(win, pixel.R(0, 0, 256, 256).Moved(t.Vec()))

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
