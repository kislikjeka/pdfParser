package spider

import (
	"github.com/cavaliercoder/grab"
	"os"
	"reflect"
)

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
