package app

import (
	"fmt"
	"github.com/valyala/fasthttp"
	"homework6/internal/ads"
)

type App interface {
	CreateAd(ctx *fasthttp.RequestCtx, title string, text string, authorId int64) (*ads.Ad, error)
	ChangeAdStatus(ctx *fasthttp.RequestCtx, id int64, userId int64, published bool) (*ads.Ad, error)
	UpdateAd(ctx *fasthttp.RequestCtx, id int64, userId int64, title string, text string) (*ads.Ad, error)
}

type app struct {
	repo Repository
}

func (a *app) CreateAd(ctx *fasthttp.RequestCtx, title string, text string, authorId int64) (*ads.Ad, error) {
	ad, err := ads.New(title, text, authorId)
	fmt.Println("app CreateAd", ad, err)
	if err != nil {
		ctx.Error("Bad request", 400)
	}
	err = a.repo.Add(&ad)
	if err != nil {
		return nil, err
	}
	return &ad, nil
}

func (a *app) ChangeAdStatus(ctx *fasthttp.RequestCtx, id int64, userId int64, published bool) (*ads.Ad, error) {
	adAuthor := a.repo.GetAd(id).AuthorID
	if adAuthor != userId {
		ctx.Error("Unauthorized", 403)
	}
	ad, err := a.repo.ChangeAdStatus(id, published)
	if err != nil {
		ctx.Error("Bad request", 400)
	}
	return ad, nil
}

func (a *app) UpdateAd(ctx *fasthttp.RequestCtx, id int64, userId int64, title string, text string) (*ads.Ad, error) {
	adAuthor := a.repo.GetAd(id).AuthorID
	fmt.Println("app UpdateAd adAuthor:", adAuthor, "userId:", userId, adAuthor == userId)
	if adAuthor != userId {
		ctx.Error("Unauthorized", 403)
	}
	ad, err := a.repo.UpdateAd(id, title, text)
	fmt.Println("app UpdateAd result", ad, err)
	if err != nil {
		ctx.Error("Bad request", 400)
		//return nil, err
	}
	return ad, nil
}

type Repository interface {
	Add(ad *ads.Ad) error
	GetAd(id int64) *ads.Ad
	ChangeAdStatus(id int64, published bool) (*ads.Ad, error)
	UpdateAd(id int64, title string, text string) (*ads.Ad, error)
}

func NewApp(repo Repository) App {
	return &app{repo: repo}
}
