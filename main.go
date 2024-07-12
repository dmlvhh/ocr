package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"ocr/service"
	"ocr/utils"
)

var Log log.Logger

func main() {
	r := gin.Default()
	r.POST("/ocr", OCR)
	r.POST("/uploadFile", UploadFile)
	r.Run(":9999")
}

// OCR 识别api
func OCR(c *gin.Context) {
	path, err := utils.UploadFile(c)
	if err != nil {
		c.JSON(400, gin.H{
			"msg":  "上传失败",
			"code": 400,
			"data": nil,
		})
		Log.Fatal(err)
	}
	data, err := service.OCR(path[1:])
	if err != nil {
		c.JSON(400, gin.H{
			"msg":  "识别失败",
			"code": 400,
			"data": nil,
		})
		Log.Fatal(err)
	}
	c.JSON(200, gin.H{
		"msg":  "识别成功",
		"code": 200,
		"data": data,
	})
}

// UploadFile 文件上传api
func UploadFile(c *gin.Context) {
	path, err := utils.UploadFile(c)
	if err != nil {
		c.JSON(400, gin.H{
			"msg":  "上传失败",
			"code": 400,
			"data": nil,
		})
		return
	}
	c.JSON(200, gin.H{
		"msg":  "上传成功",
		"code": 200,
		"data": map[string]string{
			"path": path,
		},
	})
}
