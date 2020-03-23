package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/joho/godotenv"
)

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

func requestProceedingsPage() *http.Response {
	url := "https://www.anlp.jp/nlp2020/program_online/index.html"
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

func main() {
	// get request to nlp2020 proceedings page
	var resp *http.Response = requestProceedingsPage()
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
				fmt.Println(s2.Text())
			})
		})
	}
}
