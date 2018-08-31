package slippy

import "github.com/faiface/pixel"

type Options struct {
	Center Coordinate
	Zoom   int
	Bounds pixel.Rect
}

type Map struct {
	opts     Options
	origin   pixel.Vec
	area     pixel.Rect
	tiles    []ImageTile
	provider *Provider
}

func New(opts Options) *Map {
	m := &Map{
		opts:     opts,
		provider: NewProvider(),
	}
	m.init()
	return m
}

func (m *Map) init() {
	var (
		opts   = m.opts
		bounds = opts.Bounds
		half   = bounds.Min.To(bounds.Max).Scaled(0.5)
	)
	m.origin = opts.Center.Vec(opts.Zoom).Sub(half)
	m.area = bounds.Moved(m.origin)
	m.tiles = fromRect(m.area, opts.Zoom)
}

func (m *Map) SetOptions(opts Options) {
	m.opts = opts
	m.init()
}

func (m *Map) Zoom() int {
	return m.opts.Zoom
}

func (m *Map) SetZoom(zoom int) {
	m.opts.Zoom = zoom
	m.init()
}

func (m *Map) Bounds() pixel.Rect {
	return m.opts.Bounds
}

func (m *Map) SetBounds(bounds pixel.Rect) {
	m.opts.Bounds = bounds
	m.init()
}

func (m *Map) Center() Coordinate {
	return m.opts.Center
}

func (m *Map) CenterVec() pixel.Vec {
	return m.opts.Center.Vec(m.opts.Zoom)
}

func (m *Map) SetCenter(center Coordinate) {
	m.opts.Center = center
	m.init()
}

func (m *Map) SetCenterVec(center pixel.Vec) {
	m.opts.Center = Coord(center, m.opts.Zoom)
	m.init()
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

// FetchAsync
func (m *Map) FetchAsync() {
	for i, t := range m.tiles {
		if !t.Loaded {
			pic, ok := m.provider.Picture(t.Tile)
			if ok {
				m.tiles[i].SetPicture(pic)
			}
		}
	}
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
