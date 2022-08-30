package slippy

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/faiface/pixel"
	"golang.org/x/sync/semaphore"
)

type Loader struct {
	sem     *semaphore.Weighted
	ctx     context.Context
	cancel  func()
	mu      sync.Mutex
	tiles   map[Tile]*pixel.PictureData
	pending map[Tile]bool
}

func NewLoader() *Loader {
	ctx, cancel := context.WithCancel(context.Background())
	return &Loader{
		sem:     semaphore.NewWeighted(3),
		ctx:     ctx,
		cancel:  cancel,
		tiles:   map[Tile]*pixel.PictureData{},
		pending: map[Tile]bool{},
	}
}

func (p *Loader) Cancel() {
	p.cancel()
}

func (p *Loader) Fetch(t Tile) (*pixel.PictureData, error) {
	url := fmt.Sprintf("https://tile.openstreetmap.org/%d/%d/%d.png", t.Z, t.X, t.Y)
	img, err := FetchImage(url)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}

func (l *Loader) Prefetch(t Tile) {

	// make sure we're not already fetching this tile
	l.mu.Lock()
	if l.pending[t] {
		l.mu.Unlock()
		return
	}
	l.pending[t] = true
	l.mu.Unlock()

	// fetch it
	err := l.sem.Acquire(l.ctx, 1)
	if err != nil {
		return
	}
	pic, err := l.Fetch(t)
	l.sem.Release(1)

	// add it to the map
	l.mu.Lock()
	switch err {
	case nil:
		l.tiles[t] = pic
	default:
		log.Println(err)
		delete(l.pending, t)
	}
	l.mu.Unlock()
}

func (l *Loader) Picture(t Tile) (*pixel.PictureData, bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if m, ok := l.tiles[t]; ok {
		return m, true
	}
	go l.Prefetch(t)
	return nil, false
}
