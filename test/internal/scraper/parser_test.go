package scraper

import (
	"github.com/proost/go-sitemap-maker/internal/scraper"
	"html"
	"io/ioutil"
	url2 "net/url"
	"testing"
)

const testFile = "testdata/test.html"

func TestParse(t *testing.T) {
	url, _ := url2.Parse("https://www.google.com")

	content, err := ioutil.ReadFile(testFile)
	if err != nil {
		t.Errorf("Can't open %s", testFile)
	}

	result, err := scraper.Parse(url, html.UnescapeString(string(content)))
	if err != nil {
		t.Error(err)
	}

	if len(result) != 2 {
		t.Error("Fail to correct parsing")
	}

	expected := []string{"https://www.google.com", "https://www.google.com/ncr"}
	for i, ex := range expected {
		if result[i].String() != ex {
			t.Errorf("Fail to correct parsing %s", result[i].String())
		}
	}
}

func TestParseDifferentHost(t *testing.T) {
	url, _ := url2.Parse("https://www.google.co.kr")

	content, err := ioutil.ReadFile(testFile)
	if err != nil {
		t.Errorf("Can't open %s", testFile)
	}

	result, err := scraper.Parse(url, html.UnescapeString(string(content)))
	if err != nil {
		t.Error(err)
	}

	if len(result) != 0 {
		t.Error("Fail to collect right links")
	}
}
