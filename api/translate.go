package api

import (
	"fmt"
	md "github.com/JohannesKaufmann/html-to-markdown"
	pdf "github.com/adrg/go-wkhtmltopdf"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"os"
)

type TranslateRequest struct {
	Data string `json:"data"`
}

// OpenFile 判断文件是否存在  存在则OpenFile 不存在则Create
func OpenFile(filename string) (*os.File, error) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		fmt.Println("文件不存在")
		return os.Create(filename) //创建文件
	}
	fmt.Println("文件存在")
	return os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, os.ModeAppend) //打开文件
}

func Translate2Pdf(c *gin.Context) {
	var translateRequest TranslateRequest
	err := c.ShouldBindJSON(&translateRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "参数错误",
		})
		return
	}
	data := translateRequest.Data
	file, err := os.Create("sample.html")
	if err != nil { // 如果有错误，打印错误，同时返回
		fmt.Println("err = ", err)
		return
	}
	defer file.Close() // 在退出整个函数时，关闭文件
	f1, err1 := OpenFile("sample.html")
	if err1 != nil {
		log.Fatal(err1.Error())
	}
	defer f1.Close()
	_, err2 := io.WriteString(f1, data) //写入文件(字符串)
	if err2 != nil {
		log.Fatal(err2.Error())
	}
	err = pdf.Init()
	if err != nil {
		return
	}
	defer pdf.Destroy()

	object, err := pdf.NewObject("sample.html")
	if err != nil {
		log.Fatal(err)
	}
	object.Header.ContentCenter = "[title]"
	object.Header.DisplaySeparator = true

	object.Footer.ContentLeft = "[date]"
	object.Footer.ContentCenter = "Sample footer information"
	object.Footer.ContentRight = "[page]"
	object.Footer.DisplaySeparator = true

	// Create converter.
	converter, err := pdf.NewConverter()
	if err != nil {
		log.Fatal(err)
	}
	defer converter.Destroy()

	// Add created objects to the converter.
	converter.Add(object)
	// Set converter options.
	converter.Title = "Output Pdf"
	converter.PaperSize = pdf.A4
	converter.Orientation = pdf.Landscape
	converter.MarginTop = "1cm"
	converter.MarginBottom = "1cm"
	converter.MarginLeft = "10mm"
	converter.MarginRight = "10mm"

	// Convert objects and save the output PDF document.
	outFile, err := os.Create("out.pdf")
	if err != nil {
		log.Fatal(err)
	}
	defer func(outFile *os.File) {
		err := outFile.Close()
		if err != nil {

		}
	}(outFile)
	if err := converter.Run(outFile); err != nil {
		log.Fatal(err)
	}
	//c.File("out.pdf")
}
func Translate2Md(c *gin.Context) {
	var translateRequest TranslateRequest
	err := c.ShouldBindJSON(&translateRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "参数错误",
		})
		return
	}
	data := translateRequest.Data
	converter := md.NewConverter("", true, nil)

	markdown, err := converter.ConvertString(data)
	if err != nil {
		log.Fatal(err)
	}
	file, err := os.Create("out.md")
	if err != nil { // 如果有错误，打印错误，同时返回
		fmt.Println("err = ", err)
		return
	}
	defer file.Close() // 在退出整个函数时，关闭文件
	f1, err1 := OpenFile("out.md")
	if err1 != nil {
		log.Fatal(err1.Error())
	}
	defer f1.Close()
	_, err2 := io.WriteString(f1, markdown) //写入文件(字符串)
	if err2 != nil {
		log.Fatal(err2.Error())
	}
	//c.File("out.md")
}
func DownloadPdf(c *gin.Context) {
	c.File("out.pdf")
}
func DownloadMd(c *gin.Context) {
	c.File("out.md")
}
