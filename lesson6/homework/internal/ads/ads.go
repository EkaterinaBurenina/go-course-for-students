package ads

import (
	"errors"
)

type Ad struct {
	ID        int64
	Title     string
	Text      string
	AuthorID  int64
	Published bool
}

func New(title string, text string, authorId int64) (Ad, error) {
	if title == "" {
		return Ad{}, errors.New("empty title")
	}
	if len(title) > 100 {
		return Ad{}, errors.New("too long title")
	}
	if text == "" {
		return Ad{}, errors.New("empty text")
	}
	if len(text) > 500 {
		return Ad{}, errors.New("too long text")
	}

	return Ad{
		ID:       NewAdId(),
		Title:    title,
		Text:     text,
		AuthorID: authorId,
	}, nil
}

func (a *Ad) TitleUpdate(t string) (*Ad, error) {
	if t == "" {
		return a, errors.New("empty title")
	}
	if len(t) > 100 {
		return a, errors.New("too long title")
	}
	a.Title = t
	return a, nil
}

func (a *Ad) TextUpdate(t string) (*Ad, error) {
	if t == "" {
		return a, errors.New("empty text")
	}
	if len(t) > 500 {
		return a, errors.New("too long text")
	}
	a.Text = t
	return a, nil
}

func (a *Ad) PublishUpdate(f bool) (*Ad, error) {
	a.Published = f
	return a, nil
}
