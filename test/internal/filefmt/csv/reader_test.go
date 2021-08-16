package csv

import (
	"github.com/proost/go-sitemap-maker/internal/filefmt/csv"
	"testing"
)

const (
	correctTestFile = "./testdata/testReader.csv"
	invalidProtocolFile = "./testdata/testReaderInvalidProtocol.csv"
)

func TestRead(t *testing.T) {

	urls, err := csv.Read(correctTestFile)
	if err != nil {
		t.Error(err)
	}

	if len(urls) != 2 {
		t.Errorf("Length of urls in %s must be 2", correctTestFile)
	}
}

func TestReadInvalidProtocol(t *testing.T) {
	if _, err := csv.Read(invalidProtocolFile); err == nil {
		t.Error("Read Invalid Protocol Urls Must return error")
	}
}
