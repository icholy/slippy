package tiles

import (
	"testing"
)

func TestPixel(t *testing.T) {
	pixelTests := []struct {
		coords Coordinate
		pixel  Pixel
	}{
		{ClippedCoords(40.0, -105.0), Pixel{6827, 12405, 7}},
	}
	errf := "%+v coords.ToPixel(%v) %+v -> %+v"
	for _, test := range pixelTests {
		pixel := test.coords.pixel(test.pixel.Z)
		if pixel != test.pixel {
			t.Errorf(errf, test.coords, test.pixel, pixel)
		}
	}
}
