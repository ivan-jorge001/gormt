package tools

import (
	"log"
	"os"
	"path"
	"strings"
)

// BuildDir 创建目录
func BuildDir(absDir string) error {
	return os.MkdirAll(path.Dir(absDir), os.ModePerm) //生成多级目录
}

// GetModelPath 获取程序运行目录
func GetModelPath() string {
	dir, _ := os.Getwd()
	return strings.Replace(dir, "\\", "/", -1)
}

// WriteFile 写入文件
func WriteFile(filename string, src []string, isClear bool) bool {
	err2 := BuildDir(filename)
	if err2 != nil {
		log.Printf("处理目录失败:%s", err2.Error())
		return false
	}

	flag := os.O_CREATE | os.O_WRONLY | os.O_TRUNC
	if !isClear {
		flag = os.O_CREATE | os.O_RDWR | os.O_APPEND
	}

	f, err := os.OpenFile(filename, flag, 0666)
	if err != nil {
		log.Printf("%s", err)
		return false
	}

	defer f.Close()

	for _, v := range src {
		_, err3 := f.WriteString(v)
		if err3 != nil {
			log.Printf("生成失败:%s", err3.Error())
			continue
		}
	}

	return true
}
