package utils

import (
	"fmt"
	"log"

	"github.com/xuri/excelize/v2"
)

// ExcelHandler 处理Excel文件的结构体
type ExcelHandler struct {
	InputFile  *excelize.File
	OutputFile *excelize.File
}

// NewExcelHandler 创建新的Excel处理器
func NewExcelHandler(inputPath string) (*ExcelHandler, error) {
	inputFile, err := excelize.OpenFile(inputPath)
	if err != nil {
		return nil, fmt.Errorf("打开Excel文件失败: %v", err)
	}

	return &ExcelHandler{
		InputFile:  inputFile,
		OutputFile: excelize.NewFile(),
	}, nil
}

// GetRows 获取工作表中的所有行
func (h *ExcelHandler) GetRows() ([][]string, error) {
	return h.InputFile.GetRows("Sheet1")
}

// SaveOutput 保存处理结果到输出文件
func (h *ExcelHandler) SaveOutput(outputPath string) error {
	return h.OutputFile.SaveAs(outputPath)
}

// WriteResult 写入处理结果到输出文件
func (h *ExcelHandler) WriteResult(rowIndex int, input, output string) {
	h.OutputFile.SetCellValue("Sheet1", fmt.Sprintf("A%d", rowIndex+1), input)
	h.OutputFile.SetCellValue("Sheet1", fmt.Sprintf("B%d", rowIndex+1), output)
}

// Close 关闭Excel文件
func (h *ExcelHandler) Close() {
	if err := h.InputFile.Close(); err != nil {
		log.Printf("关闭输入文件失败: %v", err)
	}
}
