package helper

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"order-service/app/models/constants"
	"strconv"

	"github.com/360EntSecGroup-Skylar/excelize"
)

func GenerateExportBufferFile(data [][]interface{}, fileType string) (*bytes.Buffer, error) {
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

func GenerateXlsx(data [][]interface{}) (*bytes.Buffer, error) {
	fmt.Println("xlsx")
	x := excelize.NewFile()
	for k, v := range data {
		for i, m := range v {
			if m != nil {
				x.SetCellValue(x.GetSheetName(1), fmt.Sprintf("%s%d", constants.ExcelCollumnMapper()[i], k+1), m)
			} else {
				x.SetCellValue(x.GetSheetName(1), fmt.Sprintf("%s%d", constants.ExcelCollumnMapper()[i], k+1), "NULL")
			}
		}
	}
	b, err := x.WriteToBuffer()
	if err != nil {
		fmt.Println("err convert = ", err)
		return nil, err
	}
	return b, nil
}

func GenerateCsv(data [][]interface{}) (*bytes.Buffer, error) {
	fmt.Println("csv")
	b := new(bytes.Buffer)
	writer := csv.NewWriter(b)
	defer writer.Flush()
	for _, v := range data {
		for _, value := range v {
			var row = []string{}
			switch t := value.(type) {
			case float32:
				row = append(row, strconv.FormatFloat(float64(value.(float32)), 'f', -1, 32))
			case float64:
				row = append(row, strconv.FormatFloat(float64(value.(float64)), 'f', -1, 64))
			case string:
				if t != "" {
					row = append(row, t)
				} else {
					row = append(row, "NULL")
				}
			case []byte:
				if string(t) != "" {
					row = append(row, string(t))
				} else {
					row = append(row, "NULL")
				}
			case nil:
				row = append(row, "NULL")
			case bool:
				row = append(row, strconv.FormatBool(value.(bool)))
			case int:
				row = append(row, strconv.Itoa(int(value.(int))))
			default:
				if string(t.(string)) != "" {
					row = append(row, string(string(t.(string))))
				} else {
					row = append(row, "NULL")
				}
			}
			if err := writer.Write(row); err != nil {
				fmt.Println("error fill", err)
				return nil, err
			}
		}
	}
	return b, nil
}
