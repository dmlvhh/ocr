package utils

import (
	"github.com/gin-gonic/gin"
	"os"
	"path"
	"strconv"
	"time"
)

// UploadFile 文件上传
func UploadFile(c *gin.Context) (string, error) {
	// 读取file文件
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{
			"msg":  "上传失败",
			"code": 400,
			"data": nil,
		})
	}

	//文件大小校验，限制300KB大小
	size := file.Size
	if size/1024 > 300 {
		c.JSON(400, gin.H{
			"msg":  "超过允许上传文件大小，最大300KB",
			"code": 400,
			"data": nil,
		})
	}

	// 文件后缀校验
	extName := path.Ext(file.Filename)
	allowExtMap := map[string]bool{
		".jpg":  true,
		".png":  true,
		".gif":  true,
		".jpeg": true,
	}
	if _, ok := allowExtMap[extName]; !ok {
		c.JSON(400, gin.H{
			"msg":  "不允许上传此类文件",
			"code": 400,
			"data": nil,
		})
	}

	// 创建图片文件目录
	day := DateStr()
	dir := "./upload/" + day
	err = os.MkdirAll(dir, 0666)
	if err != nil {
		return "", err
	}

	//执行上传，毫秒级时间名称
	now := time.Now()
	unixNano := now.UnixNano() / int64(time.Millisecond)
	createFileName := strconv.FormatInt(unixNano, 10) + extName
	pathDir := path.Join(dir, createFileName)
	err = c.SaveUploadedFile(file, pathDir)
	if err != nil {
		return "", err
	}
	outerPath := "/" + pathDir
	return outerPath, nil
}

// DateStr 组装年月日字符串当做文件路径
func DateStr() string {
	template := "20060102"
	return time.Now().Format(template)
}
