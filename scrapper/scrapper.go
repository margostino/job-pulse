package scrapper

import (
	"github.com/go-rod/rod"
	"github.com/margostino/job-pulse/utils"
)

func (s Scrapper) GoPage(url string) *Scrapper {
	// TODO: only wait once
	s.page = s.browser.MustPage(url).MustWaitLoad()
	return &s
}

func (s Scrapper) GetLinkElements() rod.Elements {
	elements, err := s.page.Elements("li")
	utils.Check(err)
	return elements
}

func (s Scrapper) Click(selector string) {
	if s.exists(selector) {
		s.page.MustElement(selector).MustElement("*").MustClick()
	}
}

func (s Scrapper) Text(selector string) string {
	if s.exists(selector) {
		return s.page.MustElement(selector).MustText()
	}
	return ""
}

func (s Scrapper) Elements(selector string) rod.Elements {
	return s.page.MustElements(selector)
}

// TODO: improve the way to check existence
func (s Scrapper) exists(selector string) bool {
	elements, err := s.page.Elements(selector)
	return len(elements) > 0 || err != nil
}

func (s Scrapper) Close() {
	s.browser.MustClose()
}
