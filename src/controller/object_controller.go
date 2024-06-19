package controller

import (
	"net/http"

	"github.com/NenfuAT/24AuthorizationServer/helper"
	"github.com/NenfuAT/24AuthorizationServer/service"
	"github.com/gin-gonic/gin"
)

func PostObject(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
		return
	}

	err := helper.AuthBasic(authHeader)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	bucketName := c.Request.FormValue("bucket")
	if bucketName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Bucket name is required",
		})
		return
	}

	path := c.Request.FormValue("path")

	// フォームからファイルを取得
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Unable to get file",
		})
		return
	}
	defer file.Close()
	fileName := header.Filename

	err = service.PostObject(bucketName, path, file, fileName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"bucket": bucketName,
		"path":   path,
		"file":   fileName,
	})
}

func GetObjects(c *gin.Context) {
	var objectNames []string

	type Response struct {
		Objects []string `json:"objects"`
	}
	type Request struct {
		Bucket string `json:"bucket"`
		Prefix string `json:"prefix"`
	}

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
		return
	}

	err := helper.AuthBasic(authHeader)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	var req Request
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	objecs, err := service.GetObjects(req.Bucket, req.Prefix)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	for _, object := range objecs.Contents {
		objectNames = append(objectNames, *object.Key)
	}
	if len(objectNames) == 0 {
		objectNames = []string{} // 空の配列を返す
	}

	responseData := Response{
		Objects: objectNames,
	}

	c.JSON(http.StatusOK, responseData)
}

func GetObjectUrl(c *gin.Context) {
	type Request struct {
		Bucket string `json:"bucket"`
		Key    string `json:"key"`
	}

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
		return
	}

	err := helper.AuthBasic(authHeader)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	var req Request
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var objectNames []string
	objects, err := service.GetObjects(req.Bucket, req.Key)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for _, object := range objects.Contents {
		objectNames = append(objectNames, *object.Key)
	}
	for _, name := range objectNames {
		if name != req.Key {
			c.JSON(http.StatusBadRequest, gin.H{"error": "file not found"})
			return
		}

	}
	if len(objectNames) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file not found"})
		return
	}

	url, err := service.GetObjectUrl(req.Bucket, req.Key)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"url": url,
	})
}
