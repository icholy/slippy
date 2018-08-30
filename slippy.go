package slippy

import (
	"fmt"
	"image"
	_ "image/png"
	"log"
	"math/rand"

	"github.com/faiface/pixel"

	"github.com/icholy/slippy/tiles"
	"github.com/icholy/slippy/util"
)

func TilePictureData(t tiles.Tile) *pixel.PictureData {
	url := URL(t)
	img, err := util.FetchImage(url)
	if err != nil {
		log.Println(err)
		img = image.NewRGBA(image.Rect(0, 0, tiles.TileSize, tiles.TileSize))
	}
	return pixel.PictureDataFromImage(img)
}

type Tile struct {
	t tiles.Tile
	s *pixel.Sprite
}

func URL(t tiles.Tile) string {
	shards := []string{"a", "b", "c"}
	return fmt.Sprintf(
		"http://%[1]s.tile.openstreetmap.org/%[2]d/%[3]d/%[4]d.png",
		shards[rand.Intn(len(shards))], t.Z, t.X, t.Y,
	)
}

func LoadTile(t tiles.Tile) (Tile, error) {
	pic := TilePictureData(t)
	return Tile{
		t: t,
		s: pixel.NewSprite(pic, pic.Bounds()),
	}, nil
}

func (t Tile) Rect() pixel.Rect {
	return t.t.Rect()
}

func (t Tile) Vec() pixel.Vec {
	return t.t.Vec()
}

func (t Tile) Draw(tg pixel.Target) {
	m := float64(tiles.TileSize) / 2
	v := t.t.Vec().Add(pixel.V(m, m))
	t.s.Draw(tg, pixel.IM.Moved(v))
}

func loadTiles(r pixel.Rect, zoom int) ([]Tile, error) {
	var tt []Tile
	for _, tile := range RectTiles(r, zoom) {
		t, err := LoadTile(tile.Tile)
		if err != nil {
			return nil, err
		}
		tt = append(tt, t)
	}
	return tt, nil
}
