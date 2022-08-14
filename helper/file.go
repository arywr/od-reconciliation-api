package helper

import (
	"encoding/csv"
	"os"
)

// import (
// 	"encoding/csv"
// 	"fmt"
// 	"io"
// 	"net/http"
// 	"os"
// 	"path/filepath"
// 	"time"
// )

func ReadCSVFile(file string) (*csv.Reader, *os.File, error) {
	osFile, err := os.Open(file)
	if err != nil {
		return nil, nil, err
	}

	reader := csv.NewReader(osFile)
	reader.FieldsPerRecord = -1
	reader.Comma = ';'
	reader.LazyQuotes = true

	return reader, osFile, nil
}

// func ReadCsvFile(request *http.Request, file string) (*csv.Reader, *os.File, error) {
// 	_, _, day := time.Now().Date()

// 	var readFile string

// 	if file != "" {
// 		readFile = file
// 	} else {
// 		readFile = fmt.Sprintf("temp/jalin %d0722.csv", int(day)-request.Day)
// 	}

// 	osFile, err := os.Open(readFile)
// 	PanicIfError(err)

// 	reader := csv.NewReader(osFile)
// 	reader.FieldsPerRecord = -1
// 	reader.Comma = ';'
// 	reader.LazyQuotes = true

// 	return reader, osFile, nil
// }

// func UploadFile(file interface{}, name string) (string, error) {
// 	var fileName string

// 	current := time.Now().Unix()

// 	file, handler, err := request.FormFile("file")
// 	defer file.Close()

// 	ext := filepath.Ext(handler.Filename)

// 	fileName = fmt.Sprintf("./temp/%d_%s%s", current, request.FormValue("platformId"), ext)

// 	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0666)
// 	PanicIfError(err)

// 	io.Copy(f, file)

// 	return fileName, nil
// }

// func DestroyFile(path string) {
// 	err := os.Remove(path)
// 	PanicIfError(err)
// }
