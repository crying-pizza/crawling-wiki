package crawling

import (
	"context"
	"fmt"
	"sync"

	// "net/http"
	"log"
	"strings"
	"net/url"

	"github.com/chromedp/chromedp"
	"github.com/PuerkitoBio/goquery"
)

type CrawlingResult struct {
	Topic string
	NewTopics []string
	Doc string
}

func CrawlWebste(mainCtx context.Context, wg *sync.WaitGroup, chTopic <-chan string, chResult chan<- CrawlingResult) {
	// fmt.Println("Crawling website:", url)
	// Add your crawling logic here
	defer wg.Done()

	ctx := mainCtx
	// options := append(chromedp.DefaultExecAllocatorOptions[:], chromedp.Flag("headless", false))
	// allocCtx, allocCancel := chromedp.NewExecAllocator(mainCtx, options...)
	allocCtx, allocCancel := chromedp.NewExecAllocator(mainCtx)
	defer allocCancel()
	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	for {
		select {
		case topic := <- chTopic:
			fmt.Printf("Crawling topic: %s\n", topic)
			newTopics, doc, err := getGetCrawlingResult(ctx, topic)
			if err != nil {
				log.Printf("Error getting crawling result: %v", err)
				continue
			}
			result := CrawlingResult{
				Topic: topic,
				NewTopics: newTopics,
				Doc: doc,
			}

			chResult <- result
		case <- mainCtx.Done():
			return
		}
	}
}

func getGetCrawlingResult(ctx context.Context, topic string) ([]string, string, error) {
	var title string
	var htmlContent string
	var resultTopic []string
	visitUrl := fmt.Sprintf("https://ko.wikipedia.org/wiki/%s", topic)

	err := chromedp.Run(ctx,
		chromedp.Navigate(visitUrl),
		chromedp.WaitVisible("body", chromedp.ByQuery),
		chromedp.Text("#firstHeading", &title, chromedp.ByQuery),
		chromedp.OuterHTML("#content", &htmlContent, chromedp.ByQuery),
	)
	defer ctx.Done()
	if err != nil {
		// log.Fatal(err)
		log.Println(err)
	}

	fmt.Println("Title:", title)
	// fmt.Println("HTML Content:", htmlContent)
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
    if err != nil {
        panic(err)
    }

	doc.Find("a[href]").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		text := s.Text()
		fmt.Printf("Link %d: %s (%s)\n", i, href, text)

		if strings.HasPrefix(href, "/wiki/") {
			trimmed := strings.TrimPrefix(href, "/wiki/")
			decoded, err := url.QueryUnescape(trimmed)
			if err != nil {
				panic(err)
			}
			resultTopic = append(resultTopic, decoded)
		}
	})

	return resultTopic, htmlContent, nil
}