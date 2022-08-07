package api

import (
	"backend/entity"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type FolderCreateRequest struct {
	Name      string `json:"name"`
	CreatorId string `json:"creatorId"`
	TeamId    int    `json:"teamId"`
	ParentId  int    `json:"parentId"`
}

func FolderCreate(ctx *gin.Context) {
	var request FolderCreateRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var folder entity.Folder
	entity.Db.Find(&folder, "name = ? and team_id = ?", request.Name, request.TeamId)
	if folder.FolderId != 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "文件夹已存在",
		})
		return
	}

	folder = entity.Folder{
		Name:       request.Name,
		TeamId:     request.TeamId,
		ParentId:   request.ParentId,
		CreatorId:  request.CreatorId,
		CreateTime: time.Now(),
		IsDeleted:  false,
		DeleterId:  request.CreatorId,
		DeleteTime: time.Now(),
	}
	result := entity.Db.Create(&folder)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":   "文件夹创建失败",
			"error": result.Error.Error(),
		})
		return
	}
	entity.Db.Where("name = ? and team_id = ?", request.Name, request.TeamId).First(&folder)
	ctx.JSON(http.StatusOK, gin.H{
		"msg":      "文件夹创建成功",
		"folderId": folder.FolderId,
	})
}
