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

type FolderDeleteRequest struct {
	FolderId  int    `json:"folderId"`
	DeleterId string `json:"deleterId"`
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

func FolderDelete(ctx *gin.Context) {
	var request FolderDeleteRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var folder entity.Folder
	entity.Db.Find(&folder, "folder_id = ?", request.FolderId)
	if folder.FolderId == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "文件夹不存在",
		})
		return
	}
	if folder.IsDeleted {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "文件夹已删除",
		})
		return
	}

	folder.IsDeleted = true
	folder.DeleterId = request.DeleterId
	folder.DeleteTime = time.Now()
	result := entity.Db.Model(&folder).Where("folder_id = ?", request.FolderId).Updates(&folder)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":   "文件夹删除失败",
			"error": result.Error.Error(),
		})
		return
	}

	var documents []entity.Document
	entity.Db.Where("parent_id = ?", request.FolderId).Find(&documents)
	for _, document := range documents {
		document.IsDeleted = true
		document.DeleterId = request.DeleterId
		document.DeleteTime = time.Now()
		result = entity.Db.Model(&document).Where("document_id = ?", document.DocId).Updates(&document)
		if result.Error != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"msg":   "文件夹中文件删除失败",
				"error": result.Error.Error(),
			})
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg": "文件夹删除成功",
	})
}
