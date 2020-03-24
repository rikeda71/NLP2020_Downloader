package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/joho/godotenv"
)

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func getNameAndPass() (string, string) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
		log.Fatal("error with loading `.env`")
		os.Exit(-1)
	}
	username := os.Getenv("USERNAME")
	password := os.Getenv("PASSWORD")
	return username, password
}

func requestPageWithAuth(url string) *http.Response {
	client := &http.Client{Timeout: time.Duration(30) * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
		os.Exit(-1)
	}
	username, password := getNameAndPass()
	req.SetBasicAuth(username, password)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		os.Exit(-1)
	}
	return resp
}

func downloadProcPdf(pdfurl string, title string, dirStr string, pidStr string) {
	pdfres := requestPageWithAuth(pdfurl)
	defer pdfres.Body.Close()

	fname := fmt.Sprintf("%s/%s_%s.pdf", dirStr, pidStr, strings.Replace(title, "/", "／", -1))
	// ファイルの存在チェック 途中からダウンロードを始めるために
	if fileExists(fname) {
		fmt.Println(fname + " exists. So, skip download")
	} else {
		file, err := os.Create(fname)
		if err != nil {
			log.Fatal(err)
			os.Exit(-1)
		}
		defer file.Close()
		fmt.Println("download pdf from `" + pdfurl + "`")
		io.Copy(file, pdfres.Body)

		time.Sleep(time.Second * 1)
	}

}

func main() {
	// get request to nlp2020 proceedings page
	baseurl := "https://www.anlp.jp/nlp2020/program_online/"
	procurl := baseurl + "index.html"
	dirStr := "proceedings"
	fmt.Println(baseurl)
	var resp *http.Response = requestPageWithAuth(procurl)
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		log.Fatal(err)
		os.Exit(-1)
	}
	for i := 16; i < 64; i++ {
		selector := fmt.Sprintf("body > div > div > div.span9 > div:nth-child(%d) > table > tbody", i)
		doc.Find(selector).Each(func(_ int, s1 *goquery.Selection) {
			s1.Find("tr").Each(func(_ int, s2 *goquery.Selection) {
				pid := s2.Find(".pid")
				pidStr, _ := pid.Find("span").Attr("id")
				pdfurl, _ := pid.Find("a").Attr("href")
				pdfurl = baseurl + pdfurl
				title := s2.Find("span.title").Text()
				// 発表しない論文は `title_no` class がついている
				if len(title) < 1 {
					title = s2.Find("span.title_no").Text()
				}
				downloadProcPdf(pdfurl, title, dirStr, pidStr)

				// ポスターがダウンロードできなくなっていたので必要ない
				// s2.Find("td > a").Each(func(_ int, s3 *goquery.Selection) {
				// 	fmt.Println(s3.Text())
				// 	fmt.Println(s3.Html())
				// 	fmt.Println(s3.Attr("href"))
				// })
			})
		})
	}
}
