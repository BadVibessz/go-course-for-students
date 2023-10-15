package app

import (
	"context"
	"homework6/internal/ads"
)

type App interface {
	CreateAd(ctx context.Context, title, text string, authorID int64) (int64, error)

	GetAdByID(ctx context.Context, id int64) (ads.Ad, error)
	GetAllAds(ctx context.Context) ([]ads.Ad, error)

	PublishAd(ctx context.Context, id int64) (ads.Ad, error)
	UnPublishAd(ctx context.Context, id int64) (ads.Ad, error)

	UpdateAd(ctx context.Context, id int64, newTitle, newText string, authorId int64) (ads.Ad, error)
}

type Repository interface {
	GetAdByID(ctx context.Context, id int64) (ads.Ad, error)
	GetAllAds(ctx context.Context) ([]ads.Ad, error)

	AddAd(ctx context.Context, ad ads.Ad) (int64, error)

	UpdateAd(ctx context.Context, id int64, newTitle, newText string, published bool) (ads.Ad, error)
}

type AdService struct {
	repo Repository
}

func (as *AdService) CreateAd(ctx context.Context, title, text string, authorID int64) (int64, error) {

	// todo: default params pattern?
	ad := ads.Ad{Title: title, Text: text, AuthorID: authorID}
	_, err := ad.Validate() // todo: code 400
	if err != nil {
		return -1, err
	}

	id, err := as.repo.AddAd(ctx, ad)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func (as *AdService) GetAdByID(ctx context.Context, id int64) (ads.Ad, error) {

	_, err := ads.ValidateID(id) // todo: code 400?
	if err != nil {
		return ads.Ad{}, err
	}

	ad, err := as.repo.GetAdByID(ctx, id)
	if err != nil {
		return ads.Ad{}, err
	}

	return ad, nil
}

func (as *AdService) GetAllAds(ctx context.Context) ([]ads.Ad, error) {

	ads, err := as.repo.GetAllAds(ctx)
	if err != nil {
		return nil, err
	}

	return ads, nil
}

func (as *AdService) PublishAd(ctx context.Context, id int64) (ads.Ad, error) {

	ad, err := as.GetAdByID(ctx, id)
	if err != nil {
		return ads.Ad{}, err
	}

	ad, err = as.repo.UpdateAd(ctx, id, ad.Title, ad.Text, true)
	if err != nil {
		return ads.Ad{}, err
	}

	return ad, nil
}

func (as *AdService) UnPublishAd(ctx context.Context, id int64) (ads.Ad, error) {

	ad, err := as.GetAdByID(ctx, id)
	if err != nil {
		return ads.Ad{}, err
	}

	// ad is not published yet
	if !ad.Published {
		return ad, nil
	}

	ad, err = as.repo.UpdateAd(ctx, id, ad.Title, ad.Text, false)
	if err != nil {
		return ads.Ad{}, err
	}

	return ad, nil
}

func (as *AdService) UpdateAd(ctx context.Context, id int64, newTitle, newText string, authorId int64) (ads.Ad, error) {

	ad, err := as.GetAdByID(ctx, id)
	if err != nil {
		return ads.Ad{}, err
	}

	if ad.AuthorID != authorId {
		return ads.Ad{}, err // todo: code 403
	}

	temp := ads.Ad{Title: newTitle, Text: newText}

	_, err = temp.Validate()
	if err != nil {
		return ads.Ad{}, err // todo: code 400
	}

	ad, err = as.repo.UpdateAd(ctx, id, newTitle, newText, ad.Published)
	if err != nil {
		return ads.Ad{}, err
	}

	return ad, nil
}

func NewApp(repo Repository) App {
	return &AdService{repo: repo}
}
