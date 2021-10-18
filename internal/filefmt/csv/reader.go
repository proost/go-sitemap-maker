package csv

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	url2 "net/url"
	"os"
)

func Read(filePath string) ([]url2.URL, error){
	var (
		err	error
		file *os.File
		rows [][]string
		urls []url2.URL
	)

	// Open File
	file, err = os.Open(filePath)
	if err != nil {
		log.Printf("Can't Open %s", filePath)
		return nil, err
	}
	defer file.Close()

	// Read File
	rows, err = csv.NewReader(bufio.NewReader(file)).ReadAll()
	if err != nil {
		log.Printf("Can't Open %s", filePath)
		return nil, err
	}

	// Get Urls By Row
	urls = make([]url2.URL, 0)
	for _, row := range rows {
		url, err := url2.Parse(row[0])
		if err != nil {
			return nil, err
		}
		if !url.IsAbs() {
			return nil, fmt.Errorf("%s is not correct format", url.String())
		}
		urls = append(urls, *url)
	}

	return urls, nil
}