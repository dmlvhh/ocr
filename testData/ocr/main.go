package main

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
	X float64
	Y float64
}

// ProductInfo 定义一个结构体来表示产品信息
type ProductInfo struct {
	Name        string
	Suitability float64
}

// DataItem 定义一个结构体来表示数据项
type DataItem struct {
	Coordinates []Point
	ProductInfo ProductInfo
}

func main() {
	// 模拟OCR输出的日志文本数据
	//	logData := `[2024/07/12 12:22:40] ppocr INFO: **********ppocr_img/imgs/11.jpg**********
	//[2024/07/12 12:22:40] ppocr DEBUG: dt_boxes num : 16, elapsed : 0.2732691764831543
	//[2024/07/12 12:22:40] ppocr DEBUG: cls num  : 16, elapsed : 0.10191226005554199
	//[2024/07/12 12:22:42] ppocr DEBUG: rec_res num  : 16, elapsed : 1.2383944988250732
	//[2024/07/12 12:22:42] ppocr INFO: [[[28.0, 37.0], [302.0, 39.0], [302.0, 72.0], [27.0, 70.0]], ('纯臻营养护发素', 0.9978386163711548)]
	//[2024/07/12 12:22:42] ppocr INFO: [[[26.0, 83.0], [173.0, 83.0], [173.0, 104.0], [26.0, 104.0]], ('产品信息/参数', 0.9898311495780945)]
	//[2024/07/12 12:22:42] ppocr INFO: [[[27.0, 112.0], [331.0, 112.0], [331.0, 135.0], [27.0, 135.0]], ('（45元/每公斤，100公斤起订）', 0.9659194946289062)]
	//[2024/07/12 12:22:42] ppocr INFO: [[[25.0, 143.0], [281.0, 143.0], [281.0, 165.0], [25.0, 165.0]], ('每瓶22元，1000瓶起订）', 0.9928649067878723)]
	//[2024/07/12 12:22:42] ppocr INFO: [[[26.0, 179.0], [300.0, 179.0], [300.0, 195.0], [26.0, 195.0]], ('【品牌】：代加工方式/OEMODM', 0.9843935966491699)]
	//[2024/07/12 12:22:42] ppocr INFO: [[[26.0, 210.0], [234.0, 210.0], [234.0, 227.0], [26.0, 227.0]], ('【品名】：纯臻营养护发素', 0.9963155388832092)]
	//[2024/07/12 12:22:42] ppocr INFO: [[[25.0, 239.0], [241.0, 239.0], [241.0, 259.0], [25.0, 259.0]], ('【产品编号】：YM-X-3011', 0.984801173210144)]
	//[2024/07/12 12:22:42] ppocr INFO: [[[413.0, 232.0], [430.0, 232.0], [430.0, 306.0], [413.0, 306.0]], ('ODMOEM', 0.9908038973808289)]
	//[2024/07/12 12:22:42] ppocr INFO: [[[24.0, 271.0], [180.0, 271.0], [180.0, 290.0], [24.0, 290.0]], ('【净含量】：220ml', 0.9892317056655884)]
	//[2024/07/12 12:22:42] ppocr INFO: [[[26.0, 303.0], [251.0, 303.0], [251.0, 319.0], [26.0, 319.0]], ('【适用人群】：适合所有肤质', 0.9909222722053528)]
	//[2024/07/12 12:22:42] ppocr INFO: [[[26.0, 335.0], [344.0, 335.0], [344.0, 352.0], [26.0, 352.0]], ('【主要成分】：鲸蜡硬脂醇、燕麦β-葡聚', 0.9828632473945618)]
	//[2024/07/12 12:22:42] ppocr INFO: [[[26.0, 364.0], [281.0, 364.0], [281.0, 384.0], [26.0, 384.0]], ('糖、椰油酰胺丙基甜菜碱、泛醌', 0.9505165219306946)]
	//[2024/07/12 12:22:42] ppocr INFO: [[[368.0, 368.0], [477.0, 368.0], [477.0, 389.0], [368.0, 389.0]], ('（成品包材）', 0.9920716285705566)]
	//[2024/07/12 12:22:42] ppocr INFO: [[[26.0, 397.0], [360.0, 397.0], [360.0, 414.0], [26.0, 414.0]], ('【主要功能】：可紧致头发磷层，从而达到', 0.9904318451881409)]
	//[2024/07/12 12:22:42] ppocr INFO: [[[28.0, 429.0], [370.0, 429.0], [370.0, 445.0], [28.0, 445.0]], ('即时持久改善头发光泽的效果，给干燥的头', 0.9874176383018494)]
	//[2024/07/12 12:22:42] ppocr INFO: [[[27.0, 458.0], [137.0, 458.0], [137.0, 479.0], [27.0, 479.0]], ('发足够的滋养', 0.9987382292747498)]`

	cmd := exec.Command("paddleocr", "--image_dir", "upload/20240712/1720766539649.jpg", "--use_angle_cls", "true", "--use_gpu", "false")
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

	// 打印提取的数据
	//fmt.Println("提取的数据项：")
	//for _, item := range dataItems {
	//	fmt.Printf("坐标点: %+v\n", item.Coordinates)
	//	fmt.Printf("产品名称: %s, 适合度: %f\n", item.ProductInfo.Name, item.ProductInfo.Suitability)
	//}
	fmt.Printf("%v", dataItems)
}
