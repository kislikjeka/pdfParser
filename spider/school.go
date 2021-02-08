package spider

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/tealeg/xlsx/v3"
	"log"
	"net/http"
	"os"
	"sync"
)

func GetSchool() {
	resDest := "./Files/Results/School"
	os.MkdirAll(resDest, os.ModePerm)
	//parsMunicip2019()
	//parsMunicip2010(resDest)
	parsSchool2015_2018(resDest)
}

func parsSchool2015_2018(dest string) {
	var wg sync.WaitGroup
	maxChan := make(chan bool, maxFileDescriptors)
	for year := 2015; year <= 2020; year++ {
		excelFile := xlsx.NewFile()
		filename := fmt.Sprintf("%s/School_%d.xlsx", dest, year)
		defer excelFile.Save(filename)
		for _, discip := range disciplines {
			sheet, err := excelFile.AddSheet(fmt.Sprintf("%s", discip))
			if err != nil {
				fmt.Println(err)
			}
			rf := NewResultFile(sheet)
			for class := 5; class <= 11; class++ {
				maxChan <- true
				wg.Add(1)
				if year == 2020 {
					url := fmt.Sprintf("https://reg.olimpiada.ru/register/russia-olympiad-%s-%d-1-%d/olympiad-protocol-static", discip, year, class)
					fmt.Println("Process url:", url)
					go processSchool2020(url, maxChan, &rf, &wg)
				} else if year == 2019 {
					url := fmt.Sprintf("https://reg.olimpiada.ru/register/russia-olympiad-%s-%d-1/public-list/prizewinners?form-number=%d", discip, year, class)
					fmt.Println("Process url:", url)
					go processSchool2019(url, maxChan, &rf, &wg)
				} else {
					url := fmt.Sprintf("https://reg.olimpiada.ru/register/russia-olympiad-%s-%d-1/public-list/prizewinners?form-number=%d", discip, year, class)
					fmt.Println("Process url:", url)
					go processSchool(url, maxChan, &rf, &wg)
				}
			}
		}
	}

	wg.Wait()
}

func processSchool(url string, maxChan chan bool, rf *ResultFile, wg *sync.WaitGroup) {
	defer wg.Done()
	defer func(maxChan chan bool) { <-maxChan }(maxChan)
	res, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Printf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Println(err)
	}

	doc.Find("tbody tr").Each(func(i int, row *goquery.Selection) {

		var line []string
		row.Find("td").Each(func(i int, col *goquery.Selection) {
			line = append(line, col.Text())
		})

		rf.WriteLine(line)
	})
}

func processSchool2019(url string, maxChan chan bool, rf *ResultFile, wg *sync.WaitGroup) {
	defer wg.Done()
	defer func(maxChan chan bool) { <-maxChan }(maxChan)
	res, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Printf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Println(err)
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
			line = append(line, col.Text())
		})

		rf.WriteLine(line)
	})
}

func processSchool2020(url string, maxChan chan bool, rf *ResultFile, wg *sync.WaitGroup) {
	processSchool2019(url, maxChan, rf, wg)
}
