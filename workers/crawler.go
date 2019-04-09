package workers

import (
	"fmt"

	"github.com/gocolly/colly"
	"github.com/maxmindlin/swarm/core/article"
	"github.com/maxmindlin/swarm/model"
)

// Crawl initiates a crawling session given a starting URL
func Crawl(start string) []model.Article {
	stories := []model.Article{}
	keywords := []string{
		"Trump",
		"music",
	}

	c := colly.NewCollector(
		// Disallow popular bottomless non-article sites
		colly.DisallowedDomains(
			// Theres gotta be a way to regexp these
			"www.youtube.com",
			"www.instagram.com",
			"www.twitter.com",
			"www.facebook.com",
			"instagram.com",
			"twitter.com",
			"youtube.com",
			"facebook.com",
		),

		// Cache responses to prevent multiple download of pages
		// even if the collector is restarted
		colly.CacheDir("./news_cache"),
		colly.MaxDepth(2),
	)

	c.OnHTML("body", func(e *colly.HTMLElement) {
		// Once the body loads, analyze the contents of the page
		story, err := article.GatherStory(e, keywords)
		if err != nil {
			return
		}

		if len(story.Keywords) == 0 {
			stories = append(stories, story)
		}
	})

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		// Find all links on the page
		link := e.Attr("href")
		e.Request.Visit(link)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.Visit(start)

	return stories
}
