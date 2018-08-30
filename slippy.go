package slippy

import (
	"fmt"
	"image"
	"math/rand"

	"github.com/faiface/pixel"

	"github.com/icholy/slippy/tiles"
	"github.com/icholy/slippy/util"
)

var (
	Placeholder = pixel.PictureDataFromImage(
		image.NewRGBA(image.Rect(0, 0, tiles.TileSize, tiles.TileSize)),
	)
)

type Options struct {
	Lat, Lon float64
	Zoom     int
	Bounds   pixel.Rect
}

type Map struct {
	opts   Options
	origin pixel.Vec
	area   pixel.Rect
	tiles  []ImageTile
}

func New(opts Options) *Map {
	var (
		origin = tiles.Vec(opts.Lat, opts.Lon, opts.Zoom)
		area   = opts.Bounds.Moved(origin)
	)
	return &Map{
		opts:   opts,
		origin: origin,
		area:   area,
		tiles:  RectTiles(area, opts.Zoom),
	}
}

func (m *Map) Fetch() error {
	for i := range m.tiles {
		if err := m.tiles[i].Fetch(); err != nil {
			return err
		}
	}
	return nil
}

func (m *Map) Draw(tg pixel.Target, mt pixel.Matrix) {
	reset := pixel.IM.Moved(pixel.ZV.Sub(m.origin)).Chained(mt)
	for _, t := range m.tiles {
		t.Draw(tg, reset)
	}
}

type ImageTile struct {
	tiles.Tile
	Sprite *pixel.Sprite
	Frame  pixel.Rect
}

func URL(t tiles.Tile) string {
	shards := []string{"a", "b", "c"}
	return fmt.Sprintf(
		"http://%[1]s.tile.openstreetmap.org/%[2]d/%[3]d/%[4]d.png",
		shards[rand.Intn(len(shards))], t.Z, t.X, t.Y,
	)
}

func (t *ImageTile) Fetch() error {
	url := URL(t.Tile)
	img, err := util.FetchImage(url)
	if err != nil {
		return err
	}
	pic := pixel.PictureDataFromImage(img)
	t.Sprite = pixel.NewSprite(pic, t.Frame)
	return nil
}

func (t ImageTile) DrawVec() pixel.Vec {
	offset := t.Frame.Min
	center := pixel.Vec{
		X: t.Frame.W() / 2,
		Y: t.Frame.H() / 2,
	}
	return t.Vec().Add(center).Add(offset)
}

func (t ImageTile) Draw(tg pixel.Target, m pixel.Matrix) {
	if t.Sprite == nil {
		return
	}
	v := t.DrawVec()
	t.Sprite.Draw(tg, pixel.IM.Moved(v).Chained(m))
}

func NewImageTile(t tiles.Tile, bounds pixel.Rect) ImageTile {
	var (
		pic     = Placeholder
		rect    = t.Rect()
		overlap = rect.Intersect(bounds)
		frame   = overlap.Moved(pixel.ZV.Sub(rect.Min))
	)
	return ImageTile{
		Tile:   t,
		Sprite: pixel.NewSprite(pic, frame),
		Frame:  frame,
	}
}

// RectTiles returns a slice of tiles requires to fully cover the rect
func RectTiles(bounds pixel.Rect, zoom int) []ImageTile {
	var (
		min = tiles.FromVec(bounds.Min, zoom)
		max = tiles.FromVec(bounds.Max, zoom)
		tt  []ImageTile
	)
	for x := min.X; x <= max.X; x++ {
		for y := max.Y; y <= min.Y; y++ {
			t := tiles.Tile{
				X: x,
				Y: y,
				Z: zoom,
			}
			tt = append(tt, NewImageTile(t, bounds))
		}
	}
	return tt
}
