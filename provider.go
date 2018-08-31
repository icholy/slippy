package slippy

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"sync"

	"github.com/faiface/pixel"
	"golang.org/x/sync/semaphore"
)

type Provider struct {
	sem     *semaphore.Weighted
	ctx     context.Context
	cancel  func()
	mu      sync.Mutex
	tiles   map[Tile]*pixel.PictureData
	pending map[Tile]bool
}

func NewProvider() *Provider {
	ctx, cancel := context.WithCancel(context.Background())
	return &Provider{
		sem:     semaphore.NewWeighted(3),
		ctx:     ctx,
		cancel:  cancel,
		tiles:   map[Tile]*pixel.PictureData{},
		pending: map[Tile]bool{},
	}
}

func (p Provider) Cancel() {
	p.cancel()
}

func (p Provider) Fetch(t Tile) (*pixel.PictureData, error) {
	shards := []string{"a", "b", "c"}
	url := fmt.Sprintf(
		"http://%[1]s.tile.openstreetmap.org/%[2]d/%[3]d/%[4]d.png",
		shards[rand.Intn(len(shards))], t.Z, t.X, t.Y,
	)
	img, err := FetchImage(url)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}

func (p *Provider) Prefetch(t Tile) {

	// make sure we're not already fetching this tile
	p.mu.Lock()
	if p.pending[t] {
		p.mu.Unlock()
		return
	}
	p.pending[t] = true
	p.mu.Unlock()

	// fetch it
	err := p.sem.Acquire(p.ctx, 1)
	if err != nil {
		return
	}
	pic, err := p.Fetch(t)
	p.sem.Release(1)

	// add it to the map
	p.mu.Lock()
	switch err {
	case nil:
		p.tiles[t] = pic
	default:
		log.Println(err)
		delete(p.pending, t)
	}
	p.mu.Unlock()
}

func (p *Provider) Picture(t Tile) (*pixel.PictureData, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if m, ok := p.tiles[t]; ok {
		return m, true
	}
	go p.Prefetch(t)
	return nil, false
}
