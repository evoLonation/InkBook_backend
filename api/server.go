package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Start(address string) {
	router := gin.Default()

	// write routers here
	router.GET("/hello", hello)
	router.GET("/api/url/./localFile/document/:url", Url)
	router.POST("/api/translate", Translate)

	//user
	userGroup := router.Group("/api/user")
	userGroup.POST("/register", UserRegister)
	userGroup.GET("/login", UserLogin)
	userGroup.GET("/information", UserInformation)
	userGroup.POST("/modify/password", UserModifyPassword)
	userGroup.POST("/modify/email", UserModifyEmail)
	userGroup.POST("/modify/introduction", UserModifyIntroduction)
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
	teamGroup.POST("/modify/introduction", TeamModifyIntro)
	teamGroup.GET("/information", TeamInformation)
	teamGroup.GET("/member", GetMember)
	teamGroup.POST("/remove", Remove)
	teamGroup.POST("/transfer", Transfer)
	teamGroup.POST("/setAdmin", SetAdmin)
	teamGroup.POST("/removeAdmin", RemoveAdmin)
	teamGroup.POST("/leave", Leave)
	teamGroup.POST("/modify-avatar", TeamModifyAvatar)
	teamGroup.GET("/get-avatar", TeamGetAvatar)
	teamGroup.POST("/confirm", Confirm)
	teamGroup.POST("/apply", Apply)
	teamGroup.GET("/getAdminNum", getAdminNum)
	teamGroup.GET("/getIdentity", getIdentity)
	teamGroup.GET("/search", SearchTeam)

	//project
	projectGroup := router.Group("/api/project")
	{
		projectGroup.POST("/create", ProjectCreate)
		projectGroup.POST("/delete", ProjectDelete)
		projectGroup.POST("/complete-delete", ProjectCompleteDelete)
		projectGroup.POST("/rename", ProjectRename)
		projectGroup.POST("/modify/intro", ProjectModifyIntro)
		projectGroup.POST("/modify/img", ProjectModifyImg)
		projectGroup.GET("/list-team", ProjectListTeam)
		projectGroup.GET("/list-user", ProjectListUser)
		projectGroup.GET("/recycle", ProjectRecycle)
		projectGroup.POST("/recover", ProjectRecover)
		projectGroup.GET("/search", ProjectSearch)
	}

	//folder
	folderGroup := router.Group("/api/folder")
	{
		folderGroup.POST("/create", FolderCreate)
		folderGroup.POST("/complete-delete", FolderCompleteDelete)
		folderGroup.POST("/rename", FolderRename)
		folderGroup.GET("/list", FolderList)
		folderGroup.POST("/move", FolderMove)
	}

	//document
	documentGroup := router.Group("/api/document")
	{
		documentGroup.POST("/create", DocumentCreate)
		documentGroup.POST("/delete", DocumentDelete)
		documentGroup.POST("/complete-delete", DocumentCompleteDelete)
		documentGroup.POST("/rename", DocumentRename)
		documentGroup.GET("/list", DocumentList)
		documentGroup.GET("/project", DocumentProject)
		documentGroup.GET("/recycle", DocumentRecycle)
		documentGroup.POST("/recover", DocumentRecover)
		documentGroup.POST("/save", DocumentSave)
		documentGroup.POST("/exit", DocumentExit)
		documentGroup.GET("/get", DocumentGet)
		documentGroup.POST("/apply-edit", DocumentApplyEdit)
		documentGroup.POST("/img", DocumentImg)
	}

	//graph
	graphGroup := router.Group("/api/graph")
	{
		graphGroup.POST("/create", GraphCreate)
		graphGroup.POST("/delete", GraphDelete)
		graphGroup.POST("/complete-delete", GraphCompleteDelete)
		graphGroup.POST("/rename", GraphRename)
		graphGroup.GET("/list", GraphList)
		graphGroup.GET("/recycle", GraphRecycle)
		graphGroup.POST("/recover", GraphRecover)
		graphGroup.POST("/save", GraphSave)
		graphGroup.POST("/exit", GraphExit)
		graphGroup.GET("/get", GraphGet)
		graphGroup.POST("/apply-edit", GraphApplyEdit)
	}

	//prototype
	prototypeGroup := router.Group("/api/prototype")
	{
		prototypeGroup.POST("/create", PrototypeCreate)
		prototypeGroup.POST("/delete", PrototypeDelete)
		prototypeGroup.POST("/complete-delete", PrototypeCompleteDelete)
		prototypeGroup.POST("/rename", PrototypeRename)
		prototypeGroup.GET("/list", PrototypeList)
		prototypeGroup.GET("/recycle", PrototypeRecycle)
		prototypeGroup.POST("/recover", PrototypeRecover)
		prototypeGroup.POST("/save", PrototypeSave)
		prototypeGroup.POST("/exit", PrototypeExit)
		prototypeGroup.GET("/get", PrototypeGet)
		prototypeGroup.POST("/apply-edit", PrototypeApplyEdit)
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
