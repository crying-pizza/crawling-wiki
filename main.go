package main

import (
	"context"
	// "fmt"
	"sync"
	// "time"

	// "log"

	// "github.com/chromedp/chromedp"
	"crawling_wiki/crawling"
)

func main() {
	crawlerCount := 2

	notCrawledTopics := make(map[string]struct{})
	crawlingTopics := make(map[string]struct{})
	crawledTopics := make(map[string]struct{})

	// docs := make(map[string]string)
	chTopic := make(chan string)
	chCrawlingResult := make(chan crawling.CrawlingResult)
	// chNewTopic := make(chan []string)
	// chDoc := make(chan string)

	firstTopic := "대한민국"
	notCrawledTopics[firstTopic] = struct{}{}

	// var mu sync.Mutex
	var wg sync.WaitGroup
	ctx := context.Background()
	for i := 0; i < crawlerCount; i++ {
		wg.Add(1)
		go crawling.CrawlWebste(ctx, &wg, chTopic, chCrawlingResult)
	}

	for (len(notCrawledTopics) > 0 || len(crawlingTopics) > 0) && len(crawledTopics) < 10 {
		// if len(notCrawledTopics) == 0 {
		// 	fmt.Println("All Topics crawled")
		// 	break
		// }

		var notCrawledTopic string
		for notCrawledTopic = range notCrawledTopics {
			break
		}

		if notCrawledTopic == "" {
			crawlingResult := <-chCrawlingResult
			updateMaps(notCrawledTopics, crawlingTopics, crawledTopics, &crawlingResult)

		} else {
			select {
			case crawlingResult := <-chCrawlingResult:
				updateMaps(notCrawledTopics, crawlingTopics, crawledTopics, &crawlingResult)


				// if len(notCrawledTopics) == 0 {
				// 	break
				// }

			case chTopic <- notCrawledTopic:
				delete(notCrawledTopics, notCrawledTopic)
				crawlingTopics[notCrawledTopic] = struct{}{}
				// fmt.Println(crawlingResult.NewTopic)
			}
		}
	}

	ctx.Done()

	close(chTopic)
	close(chCrawlingResult)

	wg.Wait()
}

func updateMaps(notCrawledTopics, crawlingTopics, crawledTopics map[string]struct{}, crawlingResult *crawling.CrawlingResult) {
	// fmt.Println(crawlingResult.Topic)
	for _, newTopic := range crawlingResult.NewTopics {
		_, ok1 := notCrawledTopics[newTopic]
		_, ok2 := crawledTopics[newTopic]
		_, ok3 := crawlingTopics[newTopic]
		if !ok1 && !ok2 && !ok3{
			notCrawledTopics[newTopic] = struct{}{}
		}
	}

	delete(crawlingTopics, crawlingResult.Topic)
	crawledTopics[crawlingResult.Topic] = struct{}{}
}
