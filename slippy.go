package slippy

import (
	"fmt"
	_ "image/png"
	"math/rand"

	"github.com/faiface/pixel"

	"github.com/icholy/slippy/tiles"
	"github.com/icholy/slippy/util"
)

func TilePictureData(t tiles.Tile) (*pixel.PictureData, error) {
	url := URL(t)
	img, err := util.FetchImage(url)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}

type ImageTile struct {
	tiles.Tile
	Sprite *pixel.Sprite
}

func URL(t tiles.Tile) string {
	shards := []string{"a", "b", "c"}
	return fmt.Sprintf(
		"http://%[1]s.tile.openstreetmap.org/%[2]d/%[3]d/%[4]d.png",
		shards[rand.Intn(len(shards))], t.Z, t.X, t.Y,
	)
}

func LoadTile(t tiles.Tile) (ImageTile, error) {
	pic, err := TilePictureData(t)
	if err != nil {
		return ImageTile{}, err
	}
	return ImageTile{
		Tile:   t,
		Sprite: pixel.NewSprite(pic, pic.Bounds()),
	}, nil
}

func (t ImageTile) Draw(tg pixel.Target) {
	m := float64(tiles.TileSize) / 2
	v := t.Vec().Add(pixel.V(m, m))
	t.Sprite.Draw(tg, pixel.IM.Moved(v))
}
