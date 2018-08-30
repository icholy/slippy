package slippy

import "github.com/faiface/pixel"

type Options struct {
	Center Coordinate
	Zoom   int
	Bounds pixel.Rect
}

type Map struct {
	opts   Options
	origin pixel.Vec
	area   pixel.Rect
	tiles  []ImageTile
}

func New(opts Options) *Map {
	var (
		bounds = opts.Bounds
		half   = bounds.Min.To(bounds.Max).Scaled(0.5)
		origin = opts.Center.Vec(opts.Zoom).Sub(half)
		area   = bounds.Moved(origin)
	)
	return &Map{
		opts:   opts,
		origin: origin,
		area:   area,
		tiles:  fromRect(area, opts.Zoom),
	}
}

// Coord returns the pixels coordinate assuming the camera is at 0,0
// and the map was drawn with the identify matrix
func (m *Map) Coord(v pixel.Vec) Coordinate {
	return Coord(v.Sub(m.opts.Bounds.Min).Add(m.origin), m.opts.Zoom)
}

// Vec returns the coordinate's pixel assuming the camera is at 0,0
// and the map was drawn with the identify matrix
func (m *Map) Vec(c Coordinate) pixel.Vec {
	return c.Vec(m.opts.Zoom).Add(m.opts.Bounds.Min).Sub(m.origin)
}

// Fetch the tile imagery
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
