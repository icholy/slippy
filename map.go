package slippy

import "github.com/faiface/pixel"

// Options contains configuration options for a map
type Options struct {
	Center Coord
	Zoom   int
	Bounds pixel.Rect
}

// Map draws a slippy map onto a pixel target
type Map struct {
	opts   Options
	origin pixel.Vec
	area   pixel.Rect
	tiles  []ImageTile
	loader *Loader
}

// New creates a new map instance
func New(opts Options) *Map {
	m := &Map{
		opts:   opts,
		loader: NewLoader(),
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

// SetOptions updates the map options
func (m *Map) SetOptions(opts Options) {
	m.opts = opts
	m.init()
}

// Zoom returns the zoom level
func (m *Map) Zoom() int {
	return m.opts.Zoom
}

// SetZoom sets the zoom level
func (m *Map) SetZoom(zoom int) {
	if 1 <= zoom && zoom <= MaxZ {
		m.opts.Zoom = zoom
		m.init()
	}
}

// Bounds returns the map view bounds
func (m *Map) Bounds() pixel.Rect {
	return m.opts.Bounds
}

// SetBounds sets the map view bounds
func (m *Map) SetBounds(bounds pixel.Rect) {
	m.opts.Bounds = bounds
	m.init()
}

// Center returns the coordinate of the center of the map
func (m *Map) Center() Coord {
	return m.opts.Center
}

// Visible checks if a coordinate is contained withing
// the currently visible map view
func (m *Map) Visible(c Coord) bool {
	return m.Bounds().Contains(m.Vec(c))
}

// CenterVec returns the vec that corresponds to the coordinate
// of the center of the map
func (m *Map) CenterVec() pixel.Vec {
	return m.Center().Vec(m.Zoom())
}

// SetCenter sets the center of the map to the provided coordinate
func (m *Map) SetCenter(center Coord) {
	m.opts.Center = center
	m.init()
}

// SetCenterVec sets the center of the map view to
// the coordinate which corresponds to the provided vec
func (m *Map) SetCenterVec(center pixel.Vec) {
	m.opts.Center = FromVec(center, m.Zoom())
	m.init()
}

// Coord returns the pixels coordinate assuming the camera is at 0,0
// and the map was drawn with the identify matrix
func (m *Map) Coord(v pixel.Vec) Coord {
	return FromVec(v.Sub(m.Bounds().Min).Add(m.origin), m.Zoom())
}

// Vec returns the coordinate's pixel assuming the camera is at 0,0
// and the map was drawn with the identify matrix
func (m *Map) Vec(c Coord) pixel.Vec {
	return c.Vec(m.Zoom()).Add(m.Bounds().Min).Sub(m.origin)
}

// FetchAsync
func (m *Map) FetchAsync() {
	for i, t := range m.tiles {
		if !t.Loaded {
			pic, ok := m.loader.Picture(t.Tile)
			if ok {
				m.tiles[i].SetPicture(pic)
			}
		}
	}
}

// FetchSync the tile imagery
func (m *Map) FetchSync() error {
	for i, t := range m.tiles {
		if !t.Loaded {
			if err := m.tiles[i].Fetch(); err != nil {
				return err
			}
		}
	}
	return nil
}

func (m *Map) Draw(tg pixel.Target, mt pixel.Matrix) {
	reset := pixel.IM.Moved(pixel.ZV.Sub(m.origin)).Chained(mt)
	for _, t := range m.tiles {
		if t.Loaded {
			t.Draw(tg, reset)
		}
	}
}

type Pusher interface {
	Push(pts ...pixel.Vec)
}

// Push is a convinience method for pushing coordinates
// into an imdraw.IMDraw
func (m *Map) Push(p Pusher, coords ...Coord) {
	for _, c := range coords {
		p.Push(m.Vec(c))
	}
}
