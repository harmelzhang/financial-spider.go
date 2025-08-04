package xls

import (
	"bytes"

	"github.com/shakinm/xlsReader/xls"
)

// 读取 XLS 数据
func ReadXls(data []byte, sheetIndex int, skipRowNum int) ([][]string, error) {
	wb, err := xls.OpenReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	sheet, err := wb.GetSheet(sheetIndex)
	if err != nil {
		return nil, err
	}

	sheetData := make([][]string, 0)
	for i := 0; i < sheet.GetNumberRows(); i++ {
		// 跳过多少行
		if i <= skipRowNum-1 {
			continue
		}

		rowData := make([]string, 0)
		row, err := sheet.GetRow(i)
		if err != nil {
			return nil, err
		}

		for _, cellData := range row.GetCols() {
			rowData = append(rowData, cellData.GetString())
		}
		sheetData = append(sheetData, rowData)
	}

	return sheetData, nil
}
