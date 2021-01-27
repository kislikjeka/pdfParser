package spider

import (
	"fmt"
	"github.com/cavaliercoder/grab"
	"github.com/tealeg/xlsx/v3"
	"os"
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

type Result struct {
	Name       string
	Class      int
	OlimpClass int
	Sum        float32
	Res        string
	School     string
}

func GetMunicip() {
	resDest := "./Files/Results"
	os.MkdirAll(resDest, os.ModePerm)
	//parsMunicip2019()
	parsMunicip2010(resDest)

	//var wg sync.WaitGroup
	//wg.Add(1)
	//maxChan := make(chan bool, 1)
	//excelFile := xlsx.NewFile()
	//defer excelFile.Save("Test.xlsx")
	//maxChan <-true
	//url := "https://reg.olimpiada.ru/district-olymp/public/winners/public.html?district=zeao&subject=geog&form=7&year=2016"
	//sheet, err := excelFile.AddSheet("Sheet")
	//if err != nil {
	//	fmt.Println(err)
	//}
	//rf := NewResultFile(sheet)
	//go processMun2010(url, maxChan, &rf, &wg)
	//wg.Wait()
}

func parsMunicip2019() {
	excelFile := xlsx.NewFile()
	defer excelFile.Save("Mun2019.xlsx")
	var wg sync.WaitGroup
	for _, discip := range disciplines {
		sheet, err := excelFile.AddSheet("Mun2019 - " + discip)
		if err != nil {
			fmt.Println(err)
		}
		rf := NewResultFile(sheet)
		wg.Add(1)
		go processMun2019(discip, &rf, &wg)
	}
	wg.Wait()
}

func parsMunicip2010(dest string) {
	var wg sync.WaitGroup
	maxChan := make(chan bool, maxFileDescriptors)
	for _, discip := range disciplines {
		for _, district := range districts {
			for year := 2010; year <= 2019; year++ {
				excelFile := xlsx.NewFile()
				filename := fmt.Sprintf("%s/%s_%s_%d.xlsx", dest, district, discip, year)
				defer excelFile.Save(filename)
				for class := 7; class <= 11; class++ {
					maxChan <- true
					wg.Add(1)
					sheet, err := excelFile.AddSheet(fmt.Sprintf("Class - %d", class))
					if err != nil {
						fmt.Println(err)
					}
					rf := NewResultFile(sheet)
					url := fmt.Sprintf("https://reg.olimpiada.ru/district-olymp/public/winners/public.html?district=%s&subject=%s&form=%d&year=%d", district, discip, class, year)
					fmt.Println("Process url:", url)
					go processMun2010(url, maxChan, &rf, &wg)
				}
			}
		}
	}
	wg.Wait()
}

//Download file from url to destination returning path to file
func DownloadFile(filepath string, url string) (string, error) {
	os.MkdirAll(filepath, os.ModePerm)
	resp, err := grab.Get(filepath, url)
	if err != nil {
		return "", err
	}

	return resp.Filename, nil
}

//Check if item is in array
func inArray(val interface{}, array interface{}) (exists bool, index int) {
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
