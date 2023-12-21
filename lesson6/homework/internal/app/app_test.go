package app

import (
	"context"
	"github.com/stretchr/testify/assert"
	"homework6/internal/ads"
	"testing"
)

//type Repository interface {
//	Add(ad *ads.Ad) error
//	GetAd(id int64) *ads.Ad
//	ChangeAdStatus(id int64, published bool) (*ads.Ad, error)
//	UpdateAd(id int64, title string, text string) (*ads.Ad, error)
//}

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
func (r *repo) UpdateAd(id int64, title string, text string) (ad *ads.Ad, err error) {
	ad = r.GetAd(id)

	if title != "" {
		ad, err = ad.TitleUpdate(title)
	}
	if text != "" {
		ad, err = ad.TextUpdate(text)
	}

	return
}

func TestApp(t *testing.T) {
	ctx := context.Background()

	app := NewApp(&repo{ads: []*ads.Ad{}})

	t.Run("add_ad", func(t *testing.T) {
		ad, err := app.CreateAd(ctx, "title", "text", 123)

		assert.Equal(t, err, nil)
		assert.Equal(t, ad.ID, int64(0))
	})
}
