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
	userGroup.GET("/team", UserTeam)
	userGroup.POST("/modify/avatar", UserModifyAvatar)
	userGroup.GET("/get-avatar", UserGetAvatar)

	//team
	teamGroup := router.Group("/api/team")
	teamGroup.POST("/create", TeamCreate)
	teamGroup.POST("/dismiss", TeamDismiss)
	teamGroup.POST("/modify/name", TeamModifyName)
	teamGroup.POST("/modify/intro", TeamModifyIntro)
	teamGroup.POST("/information", TeamInformation)
	teamGroup.POST("/member", GetMember)
	teamGroup.POST("/remove", Remove)
	teamGroup.POST("/transfer", Transfer)
	teamGroup.POST("/setAdmin", SetAdmin)
	teamGroup.POST("/leave", Leave)
	teamGroup.POST("/modify-avatar", TeamModifyAvatar)
	teamGroup.GET("/get-avatar", TeamGetAvatar)
	teamGroup.POST("/confirm", Confirm)
	teamGroup.POST("/apply", Apply)

	//project
	projectGroup := router.Group("/api/project")
	{
		projectGroup.POST("/create", ProjectCreate)
		projectGroup.POST("/delete", ProjectDelete)
		projectGroup.POST("/complete-delete", ProjectCompleteDelete)
		projectGroup.POST("/rename", ProjectRename)
		projectGroup.GET("/list", ProjectList)
		projectGroup.POST("/modify/intro", ProjectModifyIntro)
		projectGroup.POST("/modify/img", ProjectModifyImg)
	}

	//document
	documentGroup := router.Group("/api/document")
	{
		documentGroup.POST("/create", DocumentCreate)
		documentGroup.POST("/delete", DocumentDelete)
		documentGroup.POST("/complete-delete", DocumentCompleteDelete)
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
