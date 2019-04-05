package article

import (
	"errors"
	"strings"

	"github.com/maxmindlin/swarm/core"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

// Article organizes the different components of an article
type Article struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
}

// IsArticle determines whether or not a webpage is an article
func IsArticle(metaTags *goquery.Selection) bool {
	// Get the content of the "og:type" meta tag and check if its "article"
	return strings.EqualFold(core.SearchMetaTags(metaTags, "og:type"), "article")
}

/****
Getting the different components of an article will probably eventually be unique to the component.
Good way to incorporate that will probably be interfaces.
For now they are pretty duplicated though, because the process is naive.
****/

// GetArticleTitle gets the title of an article page
func GetArticleTitle(e *colly.HTMLElement, metaTags *goquery.Selection) string {
	title := core.SearchMetaTags(metaTags, "og:title")

	if title != "" {
		return title
	}
	// No meta tags gave title value, continue searching

	// Simply return the first title element
	return e.DOM.Find("title").First().Text()
}

// GetArticleDescription returns an article description string for a given page
func GetArticleDescription(e *colly.HTMLElement, metaTags *goquery.Selection) string {
	return core.SearchMetaTags(metaTags, "og:description")
}

// GetArticleURL returns a url linking back to the scraped article
func GetArticleURL(metaTags *goquery.Selection) string {
	return core.SearchMetaTags(metaTags, "og:url")
}

// IsKeywordRelevant returns a bool to whether or not an article is relevant
// to a given key word
func IsKeywordRelevant(article Article, word string) bool {
	word = strings.ToLower(word)

	// Is there an exact match for the word in the title
	inTitle := strings.Contains(strings.ToLower(article.Title), word)
	if inTitle {
		return true
	}

	// Is there an exact match in the description
	// This probably isnt good enough, should replace with a count and/or count % of total
	inDesc := strings.Contains(strings.ToLower(article.Description), word)
	if inDesc {
		return true
	}

	return false
}

// IsArticleRelevant determines if an article is relevant to any word
// in a set of keywords.
func IsArticleRelevant(article Article, keywords []string) bool {
	// The moment it is relevant to a keyword, return.
	// In the future should store keyword - article combos.
	for _, word := range keywords {
		if IsKeywordRelevant(article, word) {
			return true
		}
	}
	return false
}

// GatherStory determines if an HTML element is an Article,
// and if so, gathers relevant information into an Article object.
func GatherStory(e *colly.HTMLElement) (Article, error) {
	metaTags := core.GetMetaTags(e)
	if isArticle := IsArticle(metaTags); isArticle {
		// The HTML we were passed belongs to an article.
		// Begin gathering relevant information.
		temp := Article{}
		// Article content gathering should be concurrent.
		temp.Title = GetArticleTitle(e, metaTags)
		temp.Description = GetArticleDescription(e, metaTags)
		temp.URL = GetArticleURL(metaTags)

		return temp, nil
	}
	return Article{}, errors.New("HTML does not belong to an article")
}
