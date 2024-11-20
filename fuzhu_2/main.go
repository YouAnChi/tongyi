package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"fuzhu_2/api"
	"fuzhu_2/types"
	"fuzhu_2/utils"
)

func main() {
	// 初始化日志系统
	logFile, err := utils.InitLogger()
	if err != nil {
		log.Fatalf("初始化日志系统失败: %v", err)
	}
	defer logFile.Close()

	startTime := time.Now()
	log.Printf("程序开始执行，正在打开输入文件 'input.xlsx'...")

	// 初始化Excel处理器
	excelHandler, err := utils.NewExcelHandler("input.xlsx")
	if err != nil {
		log.Fatalf("❌ %v", err)
	}
	defer excelHandler.Close()

	// 读取所有行
	rows, err := excelHandler.GetRows()
	if err != nil {
		log.Fatalf("❌ 读取工作表失败: %v", err)
	}
	log.Printf("✅ 成功读取输入文件，共有 %d 行数据需要处理", len(rows))

	// 初始化API客户端
	apiClient := api.NewAPIClient("sk-ad297c6e95034aa896725120f452bac1")

	// 配置并发处理参数
	maxWorkers := 16
	resultChan := make(chan types.Result, len(rows))
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, maxWorkers)

	// 并发处理数据
	log.Printf("开始并发处理数据，并发数: %d", maxWorkers)
	for i, row := range rows {
		if len(row) == 0 {
			log.Printf("⚠️ 跳过第 %d 行：空行", i+1)
			continue
		}

		wg.Add(1)
		go func(rowIndex int, input string) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			output := apiClient.ProcessText(input)
			resultChan <- types.Result{
				RowIndex: rowIndex,
				Input:    input,
				Output:   output,
			}
			log.Printf("已处理第 %d 行", rowIndex+1)
		}(i, row[0])
	}

	// 等待所有处理完成
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// 收集并保存结果
	for result := range resultChan {
		excelHandler.WriteResult(result.RowIndex, result.Input, result.Output)
	}

	// 生成带时间戳的输出文件名
	outputFileName := fmt.Sprintf("output_%s.xlsx", time.Now().Format("2006-01-02_15-04-05"))

	// 保存输出文件
	if err := excelHandler.SaveOutput(outputFileName); err != nil {
		log.Printf("保存文件失败: %v", err)
	}

	// 输出统计信息
	log.Printf("✅ 处理完成！")
	log.Printf("总行数: %d", len(rows))
	log.Printf("总耗时: %v", time.Since(startTime))
	log.Printf("结果已保存到 %s", outputFileName)
}

//types/types.go: 定义所有数据结构
//utils/excel.go: 处理Excel文件相关的功能
//api/client.go: 处理API调用相关的功能
//main.go: 程序的主要流程控制

// 使用goroutine并发处理数据：
// 最多同时运行16个工作协程
// 使用channel控制并发数量
// 使用WaitGroup等待所有处理完成
