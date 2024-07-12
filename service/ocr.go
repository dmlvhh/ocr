package service

import (
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"log"
	"os/exec"
	"regexp"
	"strconv"
)

// Point 定义一个结构体来表示坐标点
type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// ProductInfo 定义一个结构体来表示产品信息
type ProductInfo struct {
	Name        string  `json:"name"`
	Suitability float64 `json:"suitability"`
}

// DataItem 定义一个结构体来表示数据项
type DataItem struct {
	Coordinates []Point     `json:"coordinates"`
	ProductInfo ProductInfo `json:"product_info"`
}

func OCR(path string) (res []DataItem, err error) {
	// 模拟OCR输出的日志文本数据
	cmd := exec.Command("paddleocr", "--image_dir", path, "--use_angle_cls", "true", "--use_gpu", "false")
	fmt.Println(cmd)
	// 执行命令并获取输出
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Failed to execute command: %v", err)
	}

	decodedOutput, err := simplifiedchinese.GBK.NewDecoder().Bytes(output)
	if err != nil {
		fmt.Printf("编码 %+v 失败, err:%s\n", output, err)
		return
	}
	// 使用正则表达式匹配含有坐标点和产品信息的行
	pattern := regexp.MustCompile(`\[\[\[([\d.]+), ([\d.]+)\], \[([\d.]+), ([\d.]+)\], \[([\d.]+), ([\d.]+)\], \[([\d.]+), ([\d.]+)\]\], \('(.+)', ([\d.]+)\)\]`)

	// 匹配并解析数据
	matches := pattern.FindAllStringSubmatch(string(decodedOutput), -1)

	var dataItems []DataItem

	// 循环处理匹配的结果
	for _, match := range matches {
		var points []Point

		// 解析坐标点
		for i := 1; i <= 8; i += 2 {
			x, _ := strconv.ParseFloat(match[i], 64)
			y, _ := strconv.ParseFloat(match[i+1], 64)
			points = append(points, Point{X: x, Y: y})
		}

		// 解析产品信息
		productName := match[9]
		suitability, _ := strconv.ParseFloat(match[10], 64)

		// 构造DataItem并添加到数据项列表中
		dataItem := DataItem{
			Coordinates: points,
			ProductInfo: ProductInfo{
				Name:        productName,
				Suitability: suitability,
			},
		}
		dataItems = append(dataItems, dataItem)
	}
	res = dataItems
	return
}
