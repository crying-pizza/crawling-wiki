package main

import (
	"context"
	// "fmt"
	"sync"
	"flag"

	// "log"

	// "github.com/chromedp/chromedp"
	"crawling_wiki/crawling"
)

func main() {
	crawlerCount := flag.Int("crawlerCount", 2, "Number of crawlers")
	startTopic := flag.String("startTopic", "대한민국", "Topic that a crawler visit")
	numMaxTopic := flag.Int("numMaxTopic", 10, "Number of max topics")

	flag.Parse()

	notCrawledTopics := make(map[string]struct{})
	crawlingTopics := make(map[string]struct{})
	crawledTopics := make(map[string]struct{})

	// docs := make(map[string]string)
	chTopic := make(chan string)
	chCrawlingResult := make(chan crawling.CrawlingResult)
	// chNewTopic := make(chan []string)
	// chDoc := make(chan string)

	notCrawledTopics[*startTopic] = struct{}{}

	// var mu sync.Mutex
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	for i := 0; i < *crawlerCount; i++ {
		wg.Add(1)
		go crawling.CrawlWebste(ctx, &wg, chTopic, chCrawlingResult)
	}

	for (len(notCrawledTopics) > 0 || len(crawlingTopics) > 0) && len(crawledTopics) < *numMaxTopic {
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

			case chTopic <- notCrawledTopic:
				delete(notCrawledTopics, notCrawledTopic)
				crawlingTopics[notCrawledTopic] = struct{}{}
			}
		}
	}

	cancel()

	// Drain the channel
	go func() {
		for range chCrawlingResult {
			
		}
	}()

	wg.Wait()

	close(chTopic)
	close(chCrawlingResult)

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
