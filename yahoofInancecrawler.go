package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/chromedp/chromedp"
)

type NewsItem struct {
	Title string
	URL   string
}

type YahooFinanceCrawler struct {
	log      func(string)
	keywords []NewsKeyword
	message chan<- string
}

func CreateYahooFinanceCrawler(log func(string), message chan<- string) *YahooFinanceCrawler {
	return &YahooFinanceCrawler{
		log:      log,
		keywords: []NewsKeyword{},
		message: message,
	}
}

func (c *YahooFinanceCrawler) CrawlYahooNews() {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()


	// Open Yahoo Finance news page
	var newsItems []NewsItem
	err := chromedp.Run(ctx,
		chromedp.Navigate("https://finance.yahoo.com/news/"),
		chromedp.Sleep(2*time.Second), // Give some time for Lazy Loading
		chromedp.Evaluate(`Array.from(document.querySelectorAll("a[href*='/news/'] h3")).map(item => ({ title: item.innerText, url: item.closest('a').href}))`, &newsItems),
	)
	if err != nil {
		log.Fatal("Crawling error: ", err)
	}

	for _, item := range newsItems {
		message := fmt.Sprintf("%s\n\n%s\n", item.Title, item.URL) // message form
		c.log(message)
		go func() { c.message <- message}()
	}
}

func (c *YahooFinanceCrawler) AddKeyword(newKeyword string) bool {
	for _, keyword := range c.keywords {
		if keyword.Text == newKeyword {
			c.log(fmt.Sprintf("Keyword '%s' already exists in the list", newKeyword))
			return false
		}
	}
	c.keywords = append(c.keywords, NewsKeyword{Text: newKeyword})
	c.log(fmt.Sprintf("Keyword '%s' added successfully", newKeyword))
	return true
}

func (c *YahooFinanceCrawler) RemoveKeyword(keywordToRemove string) bool {
	for i, keyword := range c.keywords {
		if keyword.Text == keywordToRemove {
			c.keywords = append(c.keywords[:i], c.keywords[i+1:]...)
			c.log(fmt.Sprintf("Keyword '%s' removed successfully", keywordToRemove))
			return true
		}
	}
	c.log(fmt.Sprintf("Keyword '%s' does not exist in the list", keywordToRemove))
	return false
}
