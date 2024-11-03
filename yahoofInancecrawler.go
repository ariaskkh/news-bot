package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

const repeatTimesOfScroll = 5

type NewsItem struct {
	Title string
	URL   string
}

type YahooFinanceCrawler struct {
	log      func(string)
	keywords *[]NewsKeyword
	message chan<- string
}

func CreateYahooFinanceCrawler(log func(string), message chan<- string, keywords *[]NewsKeyword) *YahooFinanceCrawler {
	return &YahooFinanceCrawler{
		log:      log,
		keywords: keywords,
		message: message,
	}
}

func (c *YahooFinanceCrawler) CrawlYahooNews() {
	if c.keywords == nil {
		log.Println("등록된 키워드가 없습니다.")
		return
	}
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()


	// Open Yahoo Finance news page
	var newsItems []NewsItem
	err := chromedp.Run(ctx,
		chromedp.Navigate("https://finance.yahoo.com/news/"),		
	)
	if err != nil {
		log.Fatal("Crawling error(navigate to yahoo): ", err)
	}

	// Scroll the page multiple times to load more content
	for i := 0; i < repeatTimesOfScroll; i++ {
		err = chromedp.Run(ctx,
			chromedp.ActionFunc(func(ctx context.Context) error {
				// Scroll to the bottom of the page
				return chromedp.Evaluate(`window.scrollTo(0, document.body.scrollHeight)`, nil).Do(ctx)
			}),
			chromedp.Sleep(2*time.Second), // Give some time for Lazy Loading
		)
	}
	if err != nil {
		log.Fatal("Crawling error(scrolling): ", err)
	}


	err = chromedp.Run(ctx,
		chromedp.Evaluate(`Array.from(document.querySelectorAll("a[href*='/news/'] h3")).map(item => ({ title: item.innerText, url: item.closest('a').href}))`, &newsItems),
	)
	if err != nil {
		log.Fatal("Crawling error(crawling, parsing): ", err)
	}

	for _, item := range newsItems {
		if c.ContainKeyword(item) {
			message := fmt.Sprintf("%s\n\n%s\n", item.Title, item.URL) // message form
			c.log(message)
			go func() { c.message <- message}()
		} else {
			continue
		}
	}

	// // test
	// item := newsItems[0]
	// if c.ContainKeyword(item) {
	// 	message := fmt.Sprintf("%s\n\n%s\n", item.Title, item.URL) // message form
	// 	c.log(message)
	// 	go func() { c.message <- message}()
	// }
}

func (c *YahooFinanceCrawler) ContainKeyword(item NewsItem) bool {
	if c.keywords == nil {
		log.Println("등록된 키워드가 없습니다.")
		return false
	}
	for _, keyword := range *c.keywords {
		if strings.Contains(strings.ToLower(item.Title), strings.ToLower(keyword.Text)) {
			return true
		}
	}
	return false
}