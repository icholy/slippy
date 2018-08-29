package util

import (
	"errors"
	"image"
	"image/color"
	_ "image/png"
	"net/http"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"golang.org/x/image/colornames"
)

const UserAgent = "Slippy/Go-Test"

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
		return nil, errors.New(resp.Status)
	}
	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return nil, err
	}
	return img, nil
}

// DrawRec draws an outline of the specified rect
func DrawRect(t pixel.Target, r pixel.Rect, c color.Color) {
	m := imdraw.New(nil)
	m.Color = c
	m.Push(
		r.Min,
		pixel.V(r.Max.X, r.Min.Y),
		r.Max,
		pixel.V(r.Min.X, r.Max.Y),
	)
	m.Polygon(1)
	m.Draw(t)
}

// DrawVec draws a small circle at the specified vec
func DrawVec(t pixel.Target, v pixel.Vec) {
	m := imdraw.New(nil)
	m.Color = colornames.Blue
	m.Push(v)
	m.Circle(10, 5)
	m.Draw(t)
}
