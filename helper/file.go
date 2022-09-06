package helper

import (
	"encoding/csv"
	"log"
	"os"

	"github.com/xuri/excelize/v2"
)

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

func ReadExcelFile(file string) ([][]string, error) {
	f, err := excelize.OpenFile(file)

	if err != nil {
		return nil, err
	}

	firstSheet := f.WorkBook.Sheets.Sheet[0].Name
	rows, err := f.GetRows(firstSheet)
	if err != nil {
		log.Fatalln(err)
	}

	return rows, nil
}

func DestroyFile(path string) error {
	err := os.Remove(path)
	return err
}
