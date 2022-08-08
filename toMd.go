package main

import (
	"fmt"
	md "github.com/JohannesKaufmann/html-to-markdown"
	"io"
	"log"
	"os"
)

func OopenFile(filename string) (*os.File, error) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		fmt.Println("文件不存在")
		return os.Create(filename) //创建文件
	}
	fmt.Println("文件存在")
	return os.OpenFile(filename, os.O_APPEND, 0666) //打开文件
}

func main() {
	converter := md.NewConverter("", true, nil)

	html := `<h1>建议收集</h1>`

	markdown, err := converter.ConvertString(html)
	if err != nil {
		log.Fatal(err)
	}
	file, err := os.Create("test.md")
	if err != nil { // 如果有错误，打印错误，同时返回
		fmt.Println("err = ", err)
		return
	}
	defer file.Close() // 在退出整个函数时，关闭文件
	f1, err1 := OopenFile("test.md")
	if err1 != nil {
		log.Fatal(err1.Error())
	}
	defer f1.Close()
	_, err2 := io.WriteString(f1, markdown) //写入文件(字符串)
	if err2 != nil {
		log.Fatal(err2.Error())
	}
	//fmt.Println("md ->", markdown)
}
