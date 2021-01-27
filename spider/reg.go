package spider

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/cavaliercoder/grab"
	"github.com/tealeg/xlsx/v3"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"pdfTimur/pdf"
	"pdfTimur/zip"
	"reflect"
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

var pdfDiscips = []string{
	"biol",
	"law",
	"phys",
	"econ",
}

type Result struct {
	Name       string
	Class      int
	OlimpClass int
	Sum        float32
	Res        string
	School     string
}

func GetMunicip() {
	excelFile := xlsx.NewFile()
	defer excelFile.Save("Result.xlsx")
	var wg sync.WaitGroup
	for _, discip := range disciplines {

		sheet, err := excelFile.AddSheet("Mun2019 - " + discip)
		if err != nil {
			fmt.Println(err)
		}
		rf := NewResultFile(sheet)
		wg.Add(1)
		go parsMunicip2019(discip, &wg, &rf)
	}
	wg.Wait()
}

func parsMunicip2019(discp string, wg *sync.WaitGroup, rf *ResultFile) {
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
				if exists, _ := in_array(discp, pdfDiscips); ok && exists {
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

func DownloadFile(filepath string, url string) (string, error) {
	os.MkdirAll(filepath, os.ModePerm)
	resp, err := grab.Get(filepath, url)
	if err != nil {
		return "", err
	}

	return resp.Filename, nil
}

func in_array(val interface{}, array interface{}) (exists bool, index int) {
	exists = false
	index = -1

	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) == true {
				index = i
				exists = true
				return
			}
		}
	}

	return
}
