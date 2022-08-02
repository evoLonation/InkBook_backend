package api

import (
	"backend/entity"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type ProjectCreateRequest struct {
	Name   string `json:"name"`
	TeamID int    `json:"teamId"`
	UserId int    `json:"userId"`
	Detail string `json:"detail"`
	ImgURL string `json:"imgUrl"`
}

type ProjectDeleteRequest struct {
	ProjectID int `json:"projectId"`
}

type ProjectRenameRequest struct {
	ProjectID int    `json:"projectId"`
	NewName   string `json:"newName"`
}

func ProjectCreate(ctx *gin.Context) {
	var request ProjectCreateRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var project entity.Project
	entity.Db.Find(&project, "name = ?", request.Name)
	if project != (entity.Project{}) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "项目已存在",
		})
		return
	}

	project = entity.Project{
		Name:       request.Name,
		TeamID:     request.TeamID,
		CreatorID:  request.UserId,
		CreateTime: time.Now(),
		IsDeleted:  false,
		Intro:      request.Detail,
		ImgURL:     request.ImgURL,
	}
}

func ProjectDelete(ctx *gin.Context) {
	// TODO
}

func ProjectRename(ctx *gin.Context) {
	// TODO
}

func ProjectList(ctx *gin.Context) {
	// TODO
}
