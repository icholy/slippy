package tiles

import (
	"fmt"
	"testing"
)

func TestTilePixel(t *testing.T) {
	tileTests := []struct {
		tile  Tile
		pixel Pixel
	}{
		{Tile{X: 26, Y: 48, Z: 7}, Pixel{X: 6656, Y: 12288, Z: 7}},
		//{Tile{26, 48, 7}, Pixel{6827, 12405, 7}},
	}
	errf := "Tile%+v: %+v -> %+v"
	for _, test := range tileTests {
		pixel := test.tile.pixel()
		if pixel != test.pixel {
			t.Errorf(errf, test.tile, test.pixel, pixel)
		}
	}
}

var (
	// These are globals to make sure that the compiler doesn't skip benchmarks
	bT Tile
)

func BenchmarkTileFromCoordinate(b *testing.B) {
	var t Tile
	z := 18
	lat := 40.7484
	lon := -73.9857
	for i := 0; i < b.N; i++ {
		t = FromCoordinate(lat, lon, z)
	}
	bT = t
}
