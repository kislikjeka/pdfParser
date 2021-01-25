package pdf

import (
	"bytes"
	"fmt"
	"github.com/ledongthuc/pdf"
	"strings"
)

func GetKeyFromPdf(path string) string {
	pdf.DebugOn = true
	content, err := readPdf(path) // Read local pdf file
	if err != nil {
		panic(err)
	}
	items := strings.Split(content, "Ключ:")
	if len(items) < 2 {
		fmt.Println("Cant find key")
		return ""
	}
	key := strings.Split(items[1], "/")[0]

	return key
}

func readPdf(path string) (string, error) {
	f, r, err := pdf.Open(path)
	// remember close file
	defer f.Close()
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	b, err := r.GetPlainText()
	if err != nil {
		return "", err
	}
	buf.ReadFrom(b)

	return buf.String(), nil
}
