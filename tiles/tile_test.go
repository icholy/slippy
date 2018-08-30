package tiles_test

import (
	"fmt"
	"testing"

	"github.com/icholy/slippy/tiles"
)

func TestTilePixel(t *testing.T) {
	tileTests := []struct {
		tile  tiles.Tile
		pixel tiles.Pixel
	}{
		{tiles.Tile{X: 26, Y: 48, Z: 7}, tiles.Pixel{X: 6656, Y: 12288, Z: 7}},
		//{Tile{26, 48, 7}, Pixel{6827, 12405, 7}},
	}
	errf := "Tile%+v: %+v -> %+v"
	for _, test := range tileTests {
		pixel := test.tile.Pixel()
		if pixel != test.pixel {
			t.Errorf(errf, test.tile, test.pixel, pixel)
		}
	}
}

var (
	// These are globals to make sure that the compiler doesn't skip benchmarks
	bT tiles.Tile
)

func BenchmarkTileFromCoordinate(b *testing.B) {
	var t tiles.Tile
	z := 18
	lat := 40.7484
	lon := -73.9857
	for i := 0; i < b.N; i++ {
		t = tiles.FromCoordinate(lat, lon, z)
	}
	bT = t
}

func ExampleFromCoordinate() {
	esbLat := 40.7484
	esbLon := -73.9857
	tile := tiles.FromCoordinate(esbLat, esbLon, 18)
	fmt.Println(tile)
}
