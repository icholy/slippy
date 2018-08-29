package tiles

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"math/rand"
	"net/http"

	"github.com/buckhx/tiles"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
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
func Vec(lat, lon float64, zoom int) pixel.Vec {
	c := tiles.ClippedCoords(lat, lon)
	return PixelVec(c.ToPixel(zoom))
}

func Coordinate(v pixel.Vec, zoom int) (lat, lon float64) {
	c := VecPixel(v, zoom).ToCoords()
	return c.Lat, c.Lon
}

func VecPixel(v pixel.Vec, zoom int) tiles.Pixel {
	return tiles.Pixel{
		X: int(v.X),
		Y: -int(v.Y),
		Z: zoom,
	}
}

func VecTile(v pixel.Vec, zoom int) tiles.Tile {
	t, _ := VecPixel(v, zoom).ToTile()
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

func DrawRect(tg pixel.Target, r pixel.Rect, c color.Color) {
	m := imdraw.New(nil)
	m.Color = c
	m.Push(
		r.Min,
		pixel.V(r.Max.X, r.Min.Y),
		r.Max,
		pixel.V(r.Min.X, r.Max.Y),
	)
	m.Polygon(1)
	m.Draw(tg)
}

func Fill(r pixel.Rect, zoom int) []tiles.Tile {
	var (
		min = VecTile(r.Min, zoom)
		max = VecTile(r.Max, zoom)
		tt  []tiles.Tile
	)
	for x := min.X; x <= max.X; x++ {
		for y := max.Y; y <= min.Y; y++ {
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

func loadTiles(r pixel.Rect, zoom int) ([]Tile, error) {
	var tt []Tile
	for _, tile := range Fill(r, zoom) {
		t, err := LoadTile(tile)
		if err != nil {
			return nil, err
		}
		tt = append(tt, t)
	}
	return tt, nil
}
