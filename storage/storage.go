package storage

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	e "tg-link-bot/lib/error"
)

type Storage interface {
	Save(page *Page) error
	PickRandom(userName string) (*Page, error)
	Remove(page *Page) error
	IsExists(page *Page) (bool, error)
}

type Page struct {
	URL      string
	UserName string
}

var ErrNoSavedPage = errors.New("no saved page")

func (p *Page) Hash() (string, error) {
	h := sha1.New()

	if _, err := io.WriteString(h, p.URL); err != nil {
		return "", e.Wrap("cant't calculate hash", err)
	}

	if _, err := io.WriteString(h, p.UserName); err != nil {
		return "", e.Wrap("cant't calculate hash", err)
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
