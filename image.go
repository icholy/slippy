package slippy

import (
	"fmt"
	"image"
	"math/rand"

	"github.com/faiface/pixel"
)

var (
	Placeholder = pixel.PictureDataFromImage(
		image.NewRGBA(image.Rect(0, 0, TileSize, TileSize)),
	)
)

type ImageTile struct {
	Tile
	Sprite *pixel.Sprite
	Frame  pixel.Rect
	Loaded bool
}

func (t *ImageTile) Fetch() error {
	shards := []string{"a", "b", "c"}
	url := fmt.Sprintf(
		"http://%[1]s.tile.openstreetmap.org/%[2]d/%[3]d/%[4]d.png",
		shards[rand.Intn(len(shards))], t.Z, t.X, t.Y,
	)
	img, err := FetchImage(url)
	if err != nil {
		return err
	}
	pic := pixel.PictureDataFromImage(img)
	t.Sprite = pixel.NewSprite(pic, t.Frame)
	t.Loaded = true
	return nil
}

func (t ImageTile) Offset() pixel.Vec {
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
	offset := t.Offset()
	t.Sprite.Draw(tg, pixel.IM.Moved(offset).Chained(m))
}

func NewImageTile(t Tile, bounds pixel.Rect) ImageTile {
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

// fromRect returns a slice of tiles requires to fully cover the rect
func fromRect(bounds pixel.Rect, zoom int) []ImageTile {
	var (
		min   = fromVec(bounds.Min, zoom)
		max   = fromVec(bounds.Max, zoom)
		tiles []ImageTile
	)
	for x := min.X; x <= max.X; x++ {
		for y := max.Y; y <= min.Y; y++ {
			t := Tile{
				X: x,
				Y: y,
				Z: zoom,
			}
			tiles = append(tiles, NewImageTile(t, bounds))
		}
	}
	return tiles
}
