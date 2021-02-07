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

var disciplines = [...]string{
	"biol",
	"law",
	"phys",
	"econ",
	"astr",
	"litr",
	"bsvf",
	"fren",
	"ekol",
	"engl",
	"geog",
	"iikt",
	"amxk",
	"span",
	"hist",
	"ital",
	"chin",
	"math",
	"germ",
	"soci",
	"russ",
	"techhome",
	"techrobo",
	"pcul",
	"chem",
}

var districts = [11]string{
	"zeao", "szao", "sao", "svao", "zao", "cao", "vao", "uzao", "uao", "uvao", "tinao",
}

var pdfDiscips = []string{
	"biol",
	"law",
	"phys",
	"econ",
}

const (
	// this is where you can specify how many maxFileDescriptors
	// you want to allow open
	maxFileDescriptors = 100
)

func GetRegional() {
	resDest := "./Files/Results/Reg"
	os.MkdirAll(resDest, os.ModePerm)
	//parsMunicip2019()
	parsReg2010_2017(resDest)
	parsReg2019(resDest)
}

func parsReg2010_2017(dest string) {
	var wg sync.WaitGroup
	maxChan := make(chan bool, maxFileDescriptors)
	for year := 2010; year <= 2017; year++ {
		excelFile := xlsx.NewFile()
		filename := fmt.Sprintf("%s/Reg_%d.xlsx", dest, year)
		defer excelFile.Save(filename)
		for _, discip := range disciplines {
			sheet, err := excelFile.AddSheet(fmt.Sprintf("%s", discip))
			if err != nil {
				fmt.Println(err)
			}
			rf := NewResultFile(sheet)
			for class := 7; class <= 11; class++ {
				maxChan <- true
				wg.Add(1)
				url := fmt.Sprintf("https://reg.olimpiada.ru/city-olymp/public/winners/public.html?subject=%s&form=%d&year=%d", discip, class, year)
				fmt.Println("Process url:", url)
				go processReg(url, maxChan, &rf, &wg)
			}
		}
	}

	wg.Wait()
}

func parsReg2019(dest string) {
	var wg sync.WaitGroup
	maxChan := make(chan bool, maxFileDescriptors)
	excelFile := xlsx.NewFile()
	filename := fmt.Sprintf("%s/Reg_%d.xlsx", dest, 2019)
	defer excelFile.Save(filename)
	for _, discip := range disciplines {
		sheet, err := excelFile.AddSheet(fmt.Sprintf("%s", discip))
		if err != nil {
			fmt.Println(err)
		}
		rf := NewResultFile(sheet)
		maxChan <- true
		wg.Add(1)
		url := fmt.Sprintf("https://reg.olimpiada.ru/register/russia-olympiad-%s-2019-3/participants-region", discip)
		fmt.Println("Process url:", url)
		go processReg(url, maxChan, &rf, &wg)
	}
	wg.Wait()

}

func processReg(url string, maxChan chan bool, rf *ResultFile, wg *sync.WaitGroup) {
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

	table.Find("tbody tr").Each(func(i int, row *goquery.Selection) {

		var line []string
		row.Find("td").Each(func(i int, col *goquery.Selection) {
			line = append(line, col.Text())
		})

		rf.WriteLine(line)
	})
}
