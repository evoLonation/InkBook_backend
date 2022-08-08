/*
package main

import (

	"fmt"
	"github.com/adrg/go-wkhtmltopdf"
	"log"

)

func main() {

	pdf.Init()

	defer pdf.Destroy()

	// Create object from url

	object1, err := pdf.NewObject("https://www.baidu.com/")

	if err != nil {

		log.Fatal(err)

	}

	object1.SetOption("footer.right", "[page]")

	// Create converter

	converter, _ := pdf.NewConverter()

	defer converter.Destroy()

	// Add created objects to the converter

	converter.AddObject(object1)

	// Add converter options

	converter.SetOption("documentTitle", "Sample document")

	converter.SetOption("margin.left", "10mm")

	converter.SetOption("margin.right", "10mm")

	converter.SetOption("margin.top", "10mm")

	converter.SetOption("margin.bottom", "10mm")

	// Convert the objects and get the output PDF document

	output, err := converter.Convert()

	if err != nil {

		log.Fatal(err)

	}

	fmt.Println(string(output))

}
*/
package main

import (
	"fmt"
	pdf "github.com/adrg/go-wkhtmltopdf"
	"io"
	"log"
	"os"
)

// OpenFile 判断文件是否存在  存在则OpenFile 不存在则Create
func OpenFile(filename string) (*os.File, error) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		fmt.Println("文件不存在")
		return os.Create(filename) //创建文件
	}
	fmt.Println("文件存在")
	return os.OpenFile(filename, os.O_APPEND, 0666) //打开文件
}

func test() {
	file, err := os.Create("sample3.html")
	if err != nil { // 如果有错误，打印错误，同时返回
		fmt.Println("err = ", err)
		return
	}
	defer file.Close() // 在退出整个函数时，关闭文件
	f1, err1 := OpenFile("sample3.html")
	if err1 != nil {
		log.Fatal(err1.Error())
	}
	defer f1.Close()
	_, err2 := io.WriteString(f1, "测试文件1") //写入文件(字符串)
	if err2 != nil {
		log.Fatal(err2.Error())
	}
	err = pdf.Init()
	if err != nil {
		return
	}
	defer pdf.Destroy()

	// Create object from file.
	object, err := pdf.NewObject("sample.html")
	if err != nil {
		log.Fatal(err)
	}
	object.Header.ContentCenter = "[title]"
	object.Header.DisplaySeparator = true

	// Create object from URL.
	/*object2, err := pdf.NewObject("https://google.com")
	if err != nil {
		log.Fatal(err)
	}*/
	object.Footer.ContentLeft = "[date]"
	object.Footer.ContentCenter = "Sample footer information"
	object.Footer.ContentRight = "[page]"
	object.Footer.DisplaySeparator = true

	/*// Create object from reader.
	inFile, err := os.Open("sample2.html")
	if err != nil {
		log.Fatal(err)
	}
	defer inFile.Close()

	object3, err := pdf.NewObjectFromReader(inFile)
	if err != nil {
		log.Fatal(err)
	}
	object3.Zoom = 1.5
	object3.TOC.Title = "Table of Contents"*/

	// Create converter.
	converter, err := pdf.NewConverter()
	if err != nil {
		log.Fatal(err)
	}
	defer converter.Destroy()

	// Add created objects to the converter.
	converter.Add(object)
	//converter.Add(object2)
	//converter.Add(object3)

	// Set converter options.
	converter.Title = "Sample document"
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
}
