package scrapper

import "github.com/go-rod/rod"

type Scrapper struct {
	browser *rod.Browser
	page    *rod.Page
}
