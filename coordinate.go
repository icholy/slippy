package slippy

import (
	"fmt"
	"math"

	"github.com/faiface/pixel"
)

// Coordinate is a simple struct for hold WGS-84 Lat Lon coordinates in degrees
type Coordinate struct {
	Lat, Lon float64
}

// Coord converts a vec into a coordinate
func Coord(v pixel.Vec, zoom int) Coordinate {
	size := float64(mapDimensions(zoom))
	x := (clip(v.X, 0, size-1) / size) - 0.5
	y := 0.5 - (clip(-v.Y, 0, size-1) / size)
	lat := 90 - 360*math.Atan(math.Exp(-y*2*math.Pi))/math.Pi
	lon := 360.0 * x
	return ClippedCoords(lat, lon)
}

// Vec converts a lat lon into a vec
func Vec(lat, lon float64, zoom int) pixel.Vec {
	return ClippedCoords(lat, lon).Vec(zoom)
}

// Vec gets the vec of the coord at the zoom level
func (c Coordinate) Vec(zoom int) pixel.Vec {
	x := (c.Lon + 180) / 360.0
	sinLat := math.Sin(c.Lat * math.Pi / 180.0)
	y := 0.5 - math.Log((1+sinLat)/(1-sinLat))/(4*math.Pi)
	size := float64(mapDimensions(zoom))
	return pixel.V(
		clip(x*size+0.5, 0, size-1),
		-clip(y*size+0.5, 0, size-1),
	)
}

// Tile returns the tile which contains the coordinate
func (c Coordinate) Tile(zoom int) Tile {
	return fromVec(c.Vec(zoom), zoom)
}

// String returns a string representation of the coordinate
func (c Coordinate) String() string {
	return fmt.Sprintf("Coordinate(%v, %v)", c.Lat, c.Lon)
}

// ClippedCoords that have been clipped to Max/Min Lat/Lon
// This can be used as a constructor to assert bad values will be clipped
func ClippedCoords(lat, lon float64) Coordinate {
	return Coordinate{
		Lat: clip(lat, MinLat, MaxLat),
		Lon: clip(lon, MinLon, MaxLon),
	}
}
