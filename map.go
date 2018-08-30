package slippy

import "github.com/faiface/pixel"

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
		origin = Vec(opts.Lat, opts.Lon, opts.Zoom)
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
