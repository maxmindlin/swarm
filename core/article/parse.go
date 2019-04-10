package article

import (
	"errors"
	"strings"
	"sync"

	"github.com/maxmindlin/swarm/core"
	"github.com/maxmindlin/swarm/model"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

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
// Meant to be run concurrently as part of a sync wait group
func GetArticleTitle(e *colly.HTMLElement, metaTags *goquery.Selection, article *model.Article, wg *sync.WaitGroup) {
	title := core.SearchMetaTags(metaTags, "og:title")

	if title != "" {
		// return title
		article.Title = title
		wg.Done()
		return
	}
	// No meta tags gave title value, continue searching

	// Simply return the first title element
	article.Title = e.DOM.Find("title").First().Text()
	wg.Done()
}

// GetArticleDescription returns an article description string for a given page
func GetArticleDescription(e *colly.HTMLElement, metaTags *goquery.Selection, article *model.Article, wg *sync.WaitGroup) {
	article.Description = core.SearchMetaTags(metaTags, "og:description")
	wg.Done()
}

// GetArticleURL returns a url linking back to the scraped article
func GetArticleURL(metaTags *goquery.Selection, article *model.Article, wg *sync.WaitGroup) {
	article.URL = core.SearchMetaTags(metaTags, "og:url")
	wg.Done()
}

// IsKeywordRelevant returns a bool to whether or not an article is relevant
// to a given key word
func IsKeywordRelevant(article model.Article, word string) bool {
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

// GatherArticleKeywords returns a slice of keywords relevant to article
func GatherArticleKeywords(article model.Article, keywords []string) []string {
	var matched []string
	for _, word := range keywords {
		if IsKeywordRelevant(article, word) {
			matched = append(matched, word)
		}
	}

	return matched
}

// GatherStory determines if an HTML element is an Article,
// and if so, gathers relevant information into an Article object.
func GatherStory(e *colly.HTMLElement, keywords []string) (model.Article, error) {
	metaTags := core.GetMetaTags(e)
	if isArticle := IsArticle(metaTags); isArticle {
		// The HTML we were passed belongs to an article.
		// Begin gathering relevant information.
		temp := model.Article{}
		var wg sync.WaitGroup
		wg.Add(3)
		go GetArticleTitle(e, metaTags, &temp, &wg)
		go GetArticleDescription(e, metaTags, &temp, &wg)
		go GetArticleURL(metaTags, &temp, &wg)
		wg.Wait()
		temp.Keywords = GatherArticleKeywords(temp, keywords)

		return temp, nil
	}
	return model.Article{}, errors.New("HTML does not belong to an article")
}
