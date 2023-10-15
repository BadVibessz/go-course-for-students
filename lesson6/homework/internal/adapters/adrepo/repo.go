package adrepo

import (
	"context"
	"homework6/internal/ads"
	"homework6/internal/app"
	"sync"
)

type AdRepository struct {
	hashMap map[string]int
	mutex   sync.Mutex
}

func New() app.Repository {
	return &AdRepository{hashMap: make(map[string]int), mutex: sync.Mutex{}}
}

func (ar *AdRepository) AddAd(ctx context.Context, ad ads.Ad) (int64, error) {

	select {
	case <-ctx.Done():
		return -1, ctx.Err()

	// else
	default:
		id := len(ar.hashMap) // todo:

	}

}

func (ar *AdRepository) GetAdByID(ctx context.Context, id int64) (ads.Ad, error) {

}

func (ar *AdRepository) GetAllAds(ctx context.Context) ([]ads.Ad, error) {
	//TODO implement me
	panic("implement me")
}

func (ar *AdRepository) UpdateAd(ctx context.Context, id int64, newTitle, newText string, published bool) (ads.Ad, error) {
	//TODO implement me
	panic("implement me")
}

// todo: implement methods
