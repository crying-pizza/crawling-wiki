# Crawling Wiki

This is my first application written in **Go**.  
It crawls Korean Wikipedia pages asynchronously.

## 🚀 How to Build & Run

```bash
go build .
./crawling_wiki
```

## 🧭 Usage

```bash
./crawling_wiki -h
```

## 📝 Example
```bash
./crawling_wiki -startTopic="대한민국" -crawlerCount=3 -numMaxTopic=20
```

## ⚙️ Flags
- -startTopic: The starting topic to crawl (default: "대한민국")
- -crawlerCount: Number of concurrent crawler goroutines (default: 2)
- -numMaxTopic: Maximum number of topics to crawl (default: 10)