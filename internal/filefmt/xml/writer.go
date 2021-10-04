package xml

import (
	"encoding/xml"
	"fmt"
	"os"
	"time"
)

type Urlset struct {
	XMLName        xml.Name `xml:"urlset"`
	Text           string   `xml:",chardata"`
	SchemaLocation string   `xml:"xsi:schemaLocation,attr"`
	Xsi            string   `xml:"xmlns:xsi,attr"`
	Xmlns          string   `xml:"xmlns,attr"`
	URL            []URL    `xml:"url"`
}

type URL struct {
	Text    string `xml:",chardata"`
	Loc     string `xml:"loc"`
	Lastmod string `xml:"lastmod"`
}

func Write(fileName string, results []string) error {
	var (
		err error
		xmlFile *os.File
		encoder *xml.Encoder
		urls []URL
		urlset *Urlset
	)

	// create File
	xmlFile, err = os.Create(fileName)
	if err != nil {
		return err
	}
	defer xmlFile.Close()

	// write header
	_, err = xmlFile.WriteString(xml.Header)
	if err != nil {
		return err
	}

	encoder = xml.NewEncoder(xmlFile)
	encoder.Indent("", "\t")

	// make URL list
	urls = makeUrlElements(results)

	urlset = newUrlsetWithAttributes()
	urlset.URL = urls

	err = encoder.Encode(urlset)
	if err != nil {
		return err
	}

	return nil
}

func makeUrlElements(results []string) []URL {
	urls := make([]URL, 0)
	now := getCurrentTimeString()
	for _, link := range results {
		url := URL{
			Loc: link,
			Lastmod: now,
		}
		urls = append(urls, url)
	}

	return urls
}

func getCurrentTimeString() string {
	now := time.Now()
	return fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d",
		now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
}

func newUrlsetWithAttributes() *Urlset {
	return &Urlset{
		Xmlns: "https://www.sitemaps.org/schemas/sitemap/0.9",
		Xsi: "http://www.w3.org/2001/XMLSchema-instance",
		SchemaLocation: "http://www.sitemaps.org/schemas/sitemap/0.9 http://www.sitemaps.org/schemas/sitemap/0.9/sitemap.xsd",
	}
}
