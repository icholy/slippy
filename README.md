# Slippy

> A slippy map package for pixel

``` go
m := slippy.New(slippy.Options{
	Zoom:   10,
	Center: slippy.ClippedCoords(43.174366, -79.231511),
	Bounds: win.Bounds(),
})

for !win.Closed() {
	m.FetchAsync()
	m.Draw(win, pixel.IM)
	win.Update()
}
```
