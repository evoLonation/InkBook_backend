package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Start(address string) {
	router := gin.Default()

	// write routers here
	router.GET("/hello", hello)

	//user
	userGroup := router.Group("/api/user")
	userGroup.POST("/register", UserRegister)
	userGroup.GET("/login", UserLogin)
	userGroup.GET("/information", UserInformation)
	userGroup.POST("/modify/password", UserModifyPassword)
	userGroup.POST("/modify/email", UserModifyEmail)
	userGroup.POST("/modify/intro", UserModifyIntroduction)
	userGroup.POST("/modify/nickname", UserModifyNickname)
	userGroup.POST("/modify/realname", UserModifyRealname)
	userGroup.POST("/send-identifying", Identifying)
	userGroup.GET("/team")

	//team
	teamGroup := router.Group("/api/team")
	teamGroup.POST("/register", TeamRegister)

	err := router.Run(address)
	if err != nil {
		return
	}
}

// 该函数返回一个gin.H，gin.H是一个map，存储着键值对，将要返回给请求者
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func hello(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"msg": "hello, world!"})
}
