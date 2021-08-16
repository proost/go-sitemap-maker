package main

import (
	"flag"
	"fmt"
	"github.com/proost/go-sitemap-maker/internal/filefmt/csv"
	"github.com/proost/go-sitemap-maker/internal/filefmt/xml"
	"github.com/proost/go-sitemap-maker/internal/scraper"
	"log"
	url2 "net/url"
	"os"
	"strings"
	"time"
)

var (
	term = flag.Int("duration", 3, "실행시간 - 기본 단위: Minute")
	csvFilePath = flag.String("csv", "", "URL 목록이 들어 있는 파일 경로")
	xmlDirPath = flag.String("dir", "", "사이트맵이 위치할 디렉토리 경로")
)

func assertNotEmpty(name string, value *string) {
	if value == nil || *value == "" {
		log.Fatalf("%s must be specified\n", name)
	}
}

func validateFileFormat(name *string, format string) {
	temp := strings.Split(*name, ".")
	if temp[len(temp) - 1] != format {
		log.Fatalf("Format of %s must be %s", *name, format)
	}
}

func validateDir(path *string) {
	fileInfo, err := os.Stat(*path)

	if err != nil {
		log.Fatalf("Can't not verify %s", *path)
	}
	if !fileInfo.IsDir(){
		log.Fatalf("%s is not directory", *path)
	}
}

func printMessage(message chan string) {
	for msg := range message {
		fmt.Println(msg)
	}
}

func init() {
	flag.Parse()

	// 파일 경로 체크
	assertNotEmpty("csv file", csvFilePath)
	assertNotEmpty("xml file", xmlDirPath)

	// 파일명 체크
	validateFileFormat(csvFilePath, "csv")

	// 디렉토리 체크
	validateDir(xmlDirPath)

	if _, err := os.Stat(*csvFilePath); os.IsNotExist(err) {
		log.Fatalf("%s not exist", *csvFilePath)
	}
}

func main() {
	var (
		err error
		urls []*url2.URL
		timer *time.Timer
		messageChan chan string
		jobs []*scraper.Job
		results map[string][]string
	)

	// Read File
	urls, err = csv.Read(*csvFilePath)
	if err != nil {
		log.Fatal(err)
	}

	// setup printing message
	messageChan = make(chan string)
	go printMessage(messageChan)

	// start jobs
	jobs = make([]*scraper.Job, 0)
	timer = time.NewTimer(time.Minute * time.Duration(*term))
	for _, url := range urls {

		origin := url
		// messageChan := messageChan

		job := scraper.NewJobFromUrl(origin)
		job.MessageChan = messageChan
		job.TimeoutChan = make(chan struct{})
		job.ResultChan = make(chan []string)

		go func() {
			job.Start()
		}()

		jobs = append(jobs, job)
	}

	<- timer.C
	messageChan <- "수집을 완료하였습니다! 파일 생성 중....."

	// end
	results = make(map[string][]string, 0)
	for _, job := range jobs {
		job.TimeoutChan <- struct{}{}
		result := <- job.ResultChan
		results[job.Origin().Host] = result
	}

	for domainName, links := range results {
		fileName := ""

		if domainName[len(domainName) - 1] == '/' {
			fileName = *xmlDirPath + domainName + ".xml"
		} else {
			fileName = *xmlDirPath + "/" + domainName + ".xml"
		}

		err = xml.Write(fileName, links)
		if err != nil {
			fmt.Printf("%v", err)
			messageChan <- fmt.Sprintf("%s의 사이트맵 파일을 만드는 데, 실패하였습니다.", domainName)
		}
	}

	messageChan <- "모든 작업을 완료하였습니다!"
}
