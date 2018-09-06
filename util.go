package slippy

import (
	"errors"
	"image"
	"image/color"
	_ "image/png"
	"math"
	"net/http"
)

// Earth Parameters
const (
	MaxZ         = 23
	MinLat       = -85.05112878
	MaxLat       = 85.05112878
	MinLon       = -180
	MaxLon       = 180
	EarthRadiusM = 6378137
	TileSize     = 256
)

// clamp the value
func clamp(val, min, max float64) float64 {
	if min > max {
		panic("clamp: min cannot be greater than max")
	}
	return math.Min(math.Max(val, min), max)
}

// Gets the size of the x, y dimensions in pixels at the given zoom level
// x == y since the map is a square
func mapDimensions(zoom int) int {
	return TileSize << uint(zoom)
}

var UserAgent = "Slippy/Go-Test"

func FetchImage(url string) (image.Image, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", UserAgent)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return BlankTile(), nil
		}
		return nil, errors.New(resp.Status)
	}
	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func BlankTile() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, TileSize, TileSize))
	for x := 0; x < TileSize; x++ {
		for y := 0; y < TileSize; y++ {
			img.Set(x, y, color.Black)
		}
	}
	return img
}
