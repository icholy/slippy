package slippy

import (
	"fmt"
	"math"

	"github.com/faiface/pixel"
)

// Coord is a simple struct for hold WGS-84 Lat Lon coordinates in degrees
type Coord [2]float64

// Lon returns the longitude
func (c Coord) Lon() float64 { return c[0] }

// Lat returns the latitude
func (c Coord) Lat() float64 { return c[1] }

// FromVec converts a vec into a coordinate
func FromVec(v pixel.Vec, zoom int) Coord {
	size := float64(mapDimensions(zoom))
	x := (clamp(v.X, 0, size-1) / size) - 0.5
	y := 0.5 - (clamp(-v.Y, 0, size-1) / size)
	lat := 90 - 360*math.Atan(math.Exp(-y*2*math.Pi))/math.Pi
	lon := 360.0 * x
	return C(lat, lon)
}

// V converts a lat lon into a vec
func V(lat, lon float64, zoom int) pixel.Vec {
	return C(lat, lon).Vec(zoom)
}

// Vec gets the vec of the coord at the zoom level
func (c Coord) Vec(zoom int) pixel.Vec {
	x := (c.Lon() + 180) / 360.0
	sinLat := math.Sin(c.Lat() * math.Pi / 180.0)
	y := 0.5 - math.Log((1+sinLat)/(1-sinLat))/(4*math.Pi)
	size := float64(mapDimensions(zoom))
	return pixel.V(
		clamp(x*size+0.5, 0, size-1),
		-clamp(y*size+0.5, 0, size-1),
	)
}

// Tile returns the tile which contains the coordinate
func (c Coord) Tile(zoom int) Tile {
	return fromVec(c.Vec(zoom), zoom)
}

// String returns a string representation of the coordinate
func (c Coord) String() string {
	return fmt.Sprintf("Coordinate(%v, %v)", c.Lat(), c.Lon())
}

// C that have been clipped to Max/Min Lat/Lon
// This can be used as a constructor to assert bad values will be clipped
func C(lat, lon float64) Coord {
	return Coord{
		clamp(lon, MinLon, MaxLon),
		clamp(lat, MinLat, MaxLat),
	}
}
