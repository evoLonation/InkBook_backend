package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Start(address string) {
	router := gin.Default()

	// write routers here
	router.GET("/hello", hello)

	router.Run(address)
}

//该函数返回一个gin.H，gin.H是一个map，存储着键值对，将要返回给请求者
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func hello(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"msg": "hello, world!"})
}
