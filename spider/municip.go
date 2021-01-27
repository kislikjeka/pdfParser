package spider

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"pdfTimur/pdf"
	"pdfTimur/zip"
	"sync"
)

func processMun2019(discp string, rf *ResultFile, wg *sync.WaitGroup) {
	defer wg.Done()
	url := fmt.Sprintf("https://reg.olimpiada.ru/register/russia-olympiad-%s-2019-2/olympiad-protocol-static", discp)
	fmt.Println("Process url:", url)
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	table := doc.Find(".beauty_table").First()

	var headerMap []string
	table.Find("thead tr td").Each(func(i int, col *goquery.Selection) {
		headerMap = append(headerMap, col.Text())
	})

	rf.WriteLine(headerMap)

	table.Find("tbody tr").Each(func(i int, row *goquery.Selection) {

		var line []string

		row.Find("td").Each(func(i int, col *goquery.Selection) {
			if i == 0 {
				href, ok := col.Find("a").First().Attr("href")
				if exists, _ := inArray(discp, pdfDiscips); ok && exists {
					fileName, err := DownloadFile("./Files/Zip/", href)
					defer os.Remove(fileName)
					fileName = filepath.Join("./", fileName)
					pdfFiles, err := zip.UnzipPDF(fileName, "./Files/PDFs/")
					if err != nil {
						fmt.Println(err)
					}
					if len(pdfFiles) != 0 {
						key := pdf.GetKeyFromPdf(pdfFiles[0])
						line = append(line, key)
						defer os.Remove(pdfFiles[0])
					}
				} else {
					line = append(line, "")
				}
			} else {
				line = append(line, col.Text())
			}
		})
		rf.WriteLine(line)
	})
}

func processMun2010(url string, maxChan chan bool, rf *ResultFile, wg *sync.WaitGroup) {
	defer wg.Done()
	defer func(maxChan chan bool) { <-maxChan }(maxChan)
	res, err := http.Get(url)
	if err != nil {
		log.Println(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Printf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Println(err)
	}

	table := doc.Find("table").First()

	table.Find("tbody tr").Each(func(i int, row *goquery.Selection) {

		var line []string

		row.Find("td").Each(func(i int, col *goquery.Selection) {
			line = append(line, col.Text())
		})

		rf.WriteLine(line)
	})
}
