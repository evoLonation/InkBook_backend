package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"os"
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
	userGroup.GET("/team", UserTeam)
	userGroup.GET("/modify-avatar", func(ctx *gin.Context) {
		file, header, err := ctx.Request.FormFile("file")
		filename := header.Filename
		fmt.Println(header.Filename)
		out, err := os.Create("./localFile/" + filename)

		if err != nil {
			log.Fatal(err)
		}
		defer out.Close()
		_, err = io.Copy(out, file)
		if err != nil {
			log.Fatal(err)
		}
	})
	userGroup.POST("/get-avatar", func(ctx *gin.Context) {
		ctx.File("./localFile/_defaultavatar.webp")
	})

	//team
	teamGroup := router.Group("/api/team")
	teamGroup.POST("/register", TeamRegister)

	//project
	projectGroup := router.Group("/api/project")
	{
		projectGroup.POST("/create", ProjectCreate)
		projectGroup.POST("/delete", ProjectDelete)
		projectGroup.POST("/rename", ProjectRename)
		projectGroup.GET("/list", ProjectList)
	}

	//document
	documentGroup := router.Group("/api/document")
	{
		documentGroup.POST("/create", DocumentCreate)
		documentGroup.POST("/delete", DocumentDelete)
		documentGroup.GET("/list", DocumentList)
	}

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
