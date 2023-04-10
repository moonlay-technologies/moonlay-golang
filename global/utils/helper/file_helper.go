package helper

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"order-service/app/models/constants"

	"github.com/360EntSecGroup-Skylar/excelize"
)

func GenerateExportBufferFile(data [][]string, fileType string) (*bytes.Buffer, error) {
	var b *bytes.Buffer
	var err error = nil
	if fileType == constants.FILE_CSV_TYPE {
		b, err = GenerateCsv(data)
	} else {
		b, err = GenerateXlsx(data)
	}
	if err != nil {
		return nil, err
	}
	return b, nil
}

func GenerateXlsx(data [][]string) (*bytes.Buffer, error) {
	fmt.Println("xlsx")
	x := excelize.NewFile()
	for k, v := range data {
		for i, m := range v {
			x.SetCellValue(x.GetSheetName(1), fmt.Sprintf("%s%d", constants.ExcelCollumnMapper()[i], k+1), m)
		}
	}
	b, err := x.WriteToBuffer()
	if err != nil {
		fmt.Println("err convert = ", err)
		return nil, err
	}
	return b, nil
}

func GenerateCsv(data [][]string) (*bytes.Buffer, error) {
	fmt.Println("csv")
	b := new(bytes.Buffer)
	writer := csv.NewWriter(b)
	defer writer.Flush()
	for _, v := range data {
		if err := writer.Write(v); err != nil {
			fmt.Println("error fill", err)
			return nil, err
		}
	}
	return b, nil
}
