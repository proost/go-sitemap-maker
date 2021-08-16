package xml

import (
	"encoding/xml"
	url2 "net/url"
	"os"
	"time"
)

var now = time.Now().Format("2021-08-15T22:04:05")

type Urlset struct {
	XMLName        xml.Name `xml:"urlset"`
	Text           string   `xml:",chardata"`
	SchemaLocation string   `xml:"schemaLocation,attr"`
	Xsi            string   `xml:"xsi,attr"`
	Xmlns          string   `xml:"xmlns,attr"`
	URL            []URL    `xml:"url"`
}

type URL struct {
	Text    string `xml:",chardata"`
	Loc     string `xml:"loc"`
	Lastmod string `xml:"lastmod"`
}

func Write(fileName string, results []*url2.URL) error {
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

func makeUrlElements(results []*url2.URL) []URL {
	urls := make([]URL, 0)
	for _, link := range results {
		url := URL{
			Loc: link.String(),
			Lastmod: now,
		}
		urls = append(urls, url)
	}

	return urls
}

func newUrlsetWithAttributes() *Urlset {
	return &Urlset{
		SchemaLocation: "http://www.sitemaps.org/schemas/sitemap/0.9 http://www.sitemaps.org/schemas/sitemap/0.9/sitemap.xsd",
		Xsi: "http://www.w3.org/2001/XMLSchema-instance",
		Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9",
	}
}
