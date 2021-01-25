package parser

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/cavaliercoder/grab"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"pdfTimur/pdf"
	"pdfTimur/zip"
	"sync"
)

var disciplines = [...]string{
	"biol",
	//"law",
	//"phys",
	//"econ",
	//"astr",
	//"litr",
	//"bsvf",
	//"fren",
	//"ekol",
	//"engl",
	//"geog",
	//"iikt",
	//"amxk",
	//"span",
	//"hist",
	//"ital",
	//"chin",
	//"math",
	//"germ",
	//"soci",
	//"russ",
	//"techhome",
	//"techrobo",
	//"pcul",
	//"chem",
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
	for _, discip := range disciplines {
		var wg sync.WaitGroup

		wg.Add(1)
		parsMunicip2019(discip, &wg)
		wg.Wait()
	}
}

func parsMunicip2019(discp string, wg *sync.WaitGroup) {
	go func() {
		defer wg.Done()
		url := fmt.Sprintf("https://reg.olimpiada.ru/register/russia-olympiad-%s-2019-2/olympiad-protocol-static", discp)
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

		headerMap := make(map[int]string)
		table.Find("thead tr td").Each(func(i int, col *goquery.Selection) {
			headerMap[i] = col.Text()
		})

		table.Find("tbody tr").Each(func(i int, row *goquery.Selection) {

			line := make(map[int]string)

			row.Find("td").Each(func(i int, col *goquery.Selection) {
				if i == 0 {
					href, ok := col.Find("a").First().Attr("href")
					if ok {
						fileName, err := DownloadFile("./Files/Zip/", href)
						defer os.Remove(fileName)
						fileName = filepath.Join("./", fileName)
						pdfFiles, err := zip.UnzipPDF(fileName, "./Files/PDFs/")
						if err != nil {
							fmt.Println(err)
						}
						if len(pdfFiles) != 0 {
							key := pdf.GetKeyFromPdf(pdfFiles[0])
							line[i] = key
							defer os.Remove(pdfFiles[0])
						}
					} else {
						line[i] = col.Text()
					}
				} else {
					val := col.Text()
					line[i] = val
				}
			})
			for i, field := range line {
				fmt.Printf("%s - %s ", headerMap[i], field)
			}
			fmt.Println("")
		})
	}()

}

func DownloadFile(filepath string, url string) (string, error) {
	os.MkdirAll(filepath, os.ModePerm)
	resp, err := grab.Get(filepath, url)
	if err != nil {
		return "", err
	}

	return resp.Filename, nil
}
