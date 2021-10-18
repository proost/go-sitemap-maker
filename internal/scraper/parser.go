package scraper

import (
	url2 "net/url"
	"regexp"
)

var pattern = regexp.MustCompile("<a.*?href=\"(.*?)\"")

func Parse(url url2.URL, content string) ([]url2.URL, error) {
	var (
		err error
		links	[]url2.URL
		matches [][]string
	)

	matches = pattern.FindAllStringSubmatch(content, -1)

	for _, val := range matches {
		var link *url2.URL

		if link, err = url.Parse(val[1]); err != nil {
			return links, err
		}

		if !link.IsAbs() {
			link.Scheme = url.Scheme
			link.Host = url.Host
		}

		if link.Host == url.Host {
			links = append(links, *link)
		}
	}

	return links, nil
}
