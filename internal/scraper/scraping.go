package scraper

import (
	"fmt"
	url2 "net/url"
)

type Job struct {
	origin *url2.URL
	result jobResult
	TimeoutChan chan struct{}
	ResultChan chan []*url2.URL
	MessageChan chan string
}

type jobResult map[*url2.URL]struct{}

func (jr jobResult) resultToArray() []*url2.URL {
	keys := make([]*url2.URL, 0)
	for k, _ := range jr {
		keys = append(keys, k)
	}
	return keys
}

func (jr jobResult) updateResult(urls []*url2.URL) {
	for _, url := range urls {
		url := url

		jr[url] = struct{}{}
	}
}

func NewJobFromUrl(url *url2.URL) *Job {
	w := Job{
		origin: url,
		result: make(jobResult),
	}

	return &w
}

func (w *Job) updateResult(foundLinks []*url2.URL) {
	w.result.updateResult(foundLinks)
}

func (w *Job) Origin() *url2.URL {
	return w.origin
}

func (w *Job) Start() {
	// init
	count := 0
	links := make([]*url2.URL, 0)
	links = append(links, w.origin)

	for {
		if len(links) > 0 {
			next := links[0]

			// Get Content from url
			content, err := GetUrlContent(next.String())
			if err != nil {
				continue
			}

			// find links
			foundLinks, err := Parse(next, content)
			if err != nil {
				continue
			}

			// update links
			w.updateResult(foundLinks)
			links = append(links, foundLinks...)
			links = links[1:]
		}

		select {
		case <- w.TimeoutChan:
			w.ResultChan <- w.result.resultToArray()
			return
		default:
			if len(w.result) - count >= 50000 {
				count = len(w.result)
				w.MessageChan <- fmt.Sprintf("%s에서 %d개의 url이 수집 중 입니다.....", w.origin.String(), count)
			}
		}
	}
}
