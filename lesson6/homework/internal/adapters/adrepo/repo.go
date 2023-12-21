package adrepo

import (
	"homework6/internal/ads"
	"homework6/internal/app"
)

type repo struct {
	ads []*ads.Ad
}

func (r *repo) Add(ad *ads.Ad) error {
	r.ads = append(r.ads, ad)
	return nil
}

func (r *repo) GetAd(id int64) *ads.Ad {
	return r.ads[id]
}

func (r *repo) ChangeAdStatus(id int64, published bool) (*ads.Ad, error) {
	ad, err := r.GetAd(id).PublishUpdate(published)
	return ad, err
}

func (r *repo) UpdateAd(id int64, title string, text string) (*ads.Ad, error) {
	ad := r.GetAd(id)

	ad, err := ad.TitleUpdate(title)
	if err != nil {
		return ad, err
	}
	ad, err = ad.TextUpdate(text)
	if err != nil {
		return ad, err
	}

	return ad, nil
}

func New() app.Repository {
	return &repo{ads: []*ads.Ad{}}
}
