package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

// InitLogger 初始化日志系统
func InitLogger() (*os.File, error) {

	err := os.MkdirAll("logs", 0755)
	if err != nil {
		return nil, fmt.Errorf("创建日志目录失败: %v", err)
	}

	// 生成日志文件名，格式为：logs/程序运行时间.log
	logFileName := filepath.Join("logs", fmt.Sprintf("%s.log", time.Now().Format("2006-01-02_15-04-05")))

	// 创建或打开日志文件
	logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("创建日志文件失败: %v", err)
	}

	// 设置日志输出到文件和终端
	log.SetOutput(os.Stdout)
	// 设置日志格式：时间 文件:行号 日志级别 内容
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// 创建多重写入器
	multiWriter := NewMultiWriter(os.Stdout, logFile)
	log.SetOutput(multiWriter)

	log.Printf("日志系统初始化完成，日志文件：%s", logFileName)
	return logFile, nil
}

// MultiWriter 实现多重写入器
type MultiWriter struct {
	writers []*os.File
}

// NewMultiWriter 创建新的多重写入器
func NewMultiWriter(stdout *os.File, files ...*os.File) *MultiWriter {
	writers := append([]*os.File{stdout}, files...)
	return &MultiWriter{writers: writers}
}

// Write 实现io.Writer接口
func (t *MultiWriter) Write(p []byte) (n int, err error) {
	for _, w := range t.writers {
		n, err = w.Write(p)
		if err != nil {
			return
		}
		if n != len(p) {
			err = fmt.Errorf("写入不完整")
			return
		}
	}
	return len(p), nil
}
