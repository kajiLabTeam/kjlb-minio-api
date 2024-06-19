package controller

import (
	"net/http"

	"github.com/NenfuAT/24AuthorizationServer/helper"
	"github.com/NenfuAT/24AuthorizationServer/service"
	"github.com/gin-gonic/gin"
)

func GetBuckets(c *gin.Context) {
	type Response struct {
		Buckets []string `json:"buckets"`
	}
	var bucketNames []string

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
	// サービスからバケットのリストを取得
	buckets, err := service.GetBuckets()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Unable to create session",
		})
		return
	}

	// バケットの名前を抽出

	for _, bucket := range buckets.Buckets {
		bucketNames = append(bucketNames, *bucket.Name)
	}
	responseData := Response{
		Buckets: bucketNames,
	}

	// JSONレスポンスを返す
	c.JSON(http.StatusOK, responseData)
}

func CreateBucket(c *gin.Context) {
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
	bucketName := c.Query("name")

	err = service.CreateBucket(bucketName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"bucket": bucketName,
	})
}
