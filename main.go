package main

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/corona10/goimagehash"
	"github.com/xuri/excelize/v2"
)

type FileInfo struct {
	FileName string
	Hash     uint64 // pHash值是64位整数
}

func main() {
	fmt.Println("开始处理图片文件...")

	var dir1Path, dir2Path string
	fmt.Print("请输入第一个文件夹的绝对路径: ")
	fmt.Scan(&dir1Path)
	fmt.Print("请输入第二个文件夹的绝对路径: ")
	fmt.Scan(&dir2Path)

	for _, path := range []string{dir1Path, dir2Path} {
		if info, err := os.Stat(path); os.IsNotExist(err) || !info.IsDir() {
			log.Fatalf("文件夹不存在或不是目录: %s", path)
		}
	}

	// 获取两个文件夹中的图片信息
	fmt.Println("正在计算文件哈希值...")
	files1 := getImageFilesInfo(dir1Path)
	files2 := getImageFilesInfo(dir2Path)

	// 创建所有文件的比对表格
	createComparisonExcel(files1, files2)

	// 创建重复文件的比对表格
	createDuplicateExcel(files1, files2)

	fmt.Println("处理完成！")
}

// 获取目录下所有图片文件的信息
func getImageFilesInfo(dirPath string) []FileInfo {
	var files []FileInfo
	filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && isImageFile(strings.ToLower(filepath.Ext(path))) {
			hash, err := calculateImageHash(path)
			if err != nil {
				log.Printf("计算文件哈希值失败 %s: %v", path, err)
				return nil
			}
			files = append(files, FileInfo{
				FileName: filepath.Base(path),
				Hash:     hash,
			})
		}
		return nil
	})
	return files
}

// 创建包含所有文件哈希值比对的Excel表格
func createComparisonExcel(files1, files2 []FileInfo) {
	f := excelize.NewFile()
	defer f.Close()

	sheet := "哈希值比对"
	f.SetSheetName("Sheet1", sheet)

	f.SetCellValue(sheet, "A1", "文件1名称")
	f.SetCellValue(sheet, "B1", "文件1哈希值")
	f.SetCellValue(sheet, "C1", "文件2名称")
	f.SetCellValue(sheet, "D1", "文件2哈希值")

	maxRows := len(files1)
	if len(files2) > maxRows {
		maxRows = len(files2)
	}

	for i := 0; i < maxRows; i++ {
		row := i + 2
		if i < len(files1) {
			f.SetCellValue(sheet, fmt.Sprintf("A%d", row), files1[i].FileName)
			f.SetCellValue(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("%016x", files1[i].Hash))
		}
		if i < len(files2) {
			f.SetCellValue(sheet, fmt.Sprintf("C%d", row), files2[i].FileName)
			f.SetCellValue(sheet, fmt.Sprintf("D%d", row), fmt.Sprintf("%016x", files2[i].Hash))
		}
	}

	timestamp := time.Now().Format("2006-01-02_15-04-05")
	filename := fmt.Sprintf("文件哈希值比对_%s.xlsx", timestamp)
	if err := f.SaveAs(filename); err != nil {
		log.Fatal("保存Excel文件失败:", err)
	}
}

// 创建重复文件比对的Excel表格
func createDuplicateExcel(files1, files2 []FileInfo) {
	f := excelize.NewFile()
	defer f.Close()

	sheet := "重复文件比对"
	f.SetSheetName("Sheet1", sheet)

	f.SetCellValue(sheet, "A1", "文件1名称")
	f.SetCellValue(sheet, "B1", "文件1哈希值")
	f.SetCellValue(sheet, "C1", "文件2名称")
	f.SetCellValue(sheet, "D1", "文件2哈希值")

	row := 2
	for _, f1 := range files1 {
		for _, f2 := range files2 {
			if f1.Hash == f2.Hash {
				f.SetCellValue(sheet, fmt.Sprintf("A%d", row), f1.FileName)
				f.SetCellValue(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("%016x", f1.Hash))
				f.SetCellValue(sheet, fmt.Sprintf("C%d", row), f2.FileName)
				f.SetCellValue(sheet, fmt.Sprintf("D%d", row), fmt.Sprintf("%016x", f2.Hash))
				row++
			}
		}
	}

	timestamp := time.Now().Format("2006-01-02_15-04-05")
	filename := fmt.Sprintf("重复文件比对_%s.xlsx", timestamp)
	if err := f.SaveAs(filename); err != nil {
		log.Fatal("保存Excel文件失败:", err)
	}
}

// 判断是否为图片文件
func isImageFile(ext string) bool {
	imageExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".bmp":  true,
		".webp": true,
	}
	return imageExts[ext]
}

// 计算图片的pHash值
func calculateImageHash(filePath string) (uint64, error) {

	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return 0, err
	}

	hash, err := goimagehash.PerceptionHash(img)
	if err != nil {
		return 0, err
	}

	return hash.GetHash(), nil
}
