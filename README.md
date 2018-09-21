# Slippy [![GoDoc](https://godoc.org/github.com/icholy/slippy?status.svg)](https://godoc.org/github.com/icholy/slippy)

> Draw [OpenStreetMap](https://www.openstreetmap.org) tiles in [pixel](https://github.com/faiface/pixel)

``` go
m := slippy.New(slippy.Options{
	Zoom:   10,
	Center: slippy.C(43.174366, -79.231511),
	Bounds: win.Bounds(),
})

for !win.Closed() {
	m.FetchAsync()
	m.Draw(win, pixel.IM)
	win.Update()
}
```
