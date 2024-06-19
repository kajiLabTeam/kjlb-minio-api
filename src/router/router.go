package router

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/NenfuAT/24AuthorizationServer/controller"
	"github.com/gin-gonic/gin"
)

func Init() {
	gin.DisableConsoleColor()
	// ログファイルを作成
	logFile, err := os.Create("log/server.log") // ファイルのパスを指定
	if err != nil {
		fmt.Println("ログファイルの作成に失敗しました:", err)
		return
	}

	// ログの出力先を設定
	gin.DefaultWriter = io.MultiWriter(logFile, os.Stdout) // ファイルとコンソールにログを出力

	r := gin.Default()
	r.GET("/hello", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello World!!")
	})

	r.GET("api/bucket/list", controller.GetBuckets)
	r.POST("api/bucket/create", controller.CreateBucket)

	r.POST("api/object/upload", controller.PostObject)
	r.POST("api/object/get", controller.GetObjectUrl)
	r.POST("api/object/list", controller.GetObjects)

	// サーバの起動とエラーハンドリング
	if err := r.Run("0.0.0.0:8000"); err != nil {
		fmt.Println("サーバの起動に失敗しました:", err)
	}
}
