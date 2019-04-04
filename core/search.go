package core

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

// SearchMetaTags will search all the given meta tags for a target property
// then return the content of that tag
// If you can get what you want from the meta tags, its usually a slam dunk.
func SearchMetaTags(tags *goquery.Selection, propTarget string) string {
	var output string
	tags.EachWithBreak(func(_ int, s *goquery.Selection) bool {
		property, _ := s.Attr("property")
		if strings.EqualFold(property, propTarget) {
			content, _ := s.Attr("content")
			output = content
			return false
		}
		return true
	})
	return output
}

// GetMetaTags returns the meta tags on a page
func GetMetaTags(e *colly.HTMLElement) *goquery.Selection {
	return e.DOM.ParentsUntil("~").Find("meta")
}
