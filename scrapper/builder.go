package scrapper

import "github.com/go-rod/rod"

func New() *Scrapper {
	return &Scrapper{
		browser: rod.New().MustConnect(),
	}
}
