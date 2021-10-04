## 실행 방법
1. 빌드
    > go build cmd/sitemap-maker/sitemapmaker.go
    * 크로스 컴파일이 필요하면: https://golang.org/doc/install/source#environment
2. 실행
   > 예시:  ./sitemapmaker -duration=1 -csv=./test.csv -dir=./
   * 옵션
      * duration: 제작기가 도는 시간. 기본은 3분. 단위는 Minute
      * csv: 읽어들일 csv 파일의 위치
      * dir: 결과값 xml파일들이 놓일 디렉토리 경로
   * csv 파일 구성
      * 수집할 root url들을 line by line으로 적어 놓는다.
      * 프로토콜과 도메인이 모두 적혀 있어야 한다.
      * 예시
         ```markdown
         <!-- test.xml -->
         https://google.com
         https://naver.com
         ```