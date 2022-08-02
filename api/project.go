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
	UserId string `json:"userId"`
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
	entity.Db.Create(&project)
	entity.Db.Where("name = ? AND team_id = ?", request.Name, request.TeamID).First(&project)
	ctx.JSON(http.StatusOK, gin.H{
		"msg":       "项目创建成功",
		"projectId": project.ProjectID,
	})
}

func ProjectDelete(ctx *gin.Context) {
	var request ProjectDeleteRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var project entity.Project
	entity.Db.Find(&project, "project_id = ?", request.ProjectID)
	if project == (entity.Project{}) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "项目不存在",
		})
		return
	}

	project.IsDeleted = true
	project.DeleteTime = time.Now()
	entity.Db.Model(&project).Where("project_id = ?", request.ProjectID).Updates(&project)
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "项目删除成功",
	})
}

func ProjectRename(ctx *gin.Context) {
	var request ProjectRenameRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var project entity.Project
	entity.Db.Find(&project, "project_id = ?", request.ProjectID)
	if project == (entity.Project{}) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "项目不存在",
		})
		return
	}

	project.Name = request.NewName
	entity.Db.Model(&project).Where("project_id = ?", request.ProjectID).Updates(&project)
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "项目重命名成功",
	})
}

func ProjectList(ctx *gin.Context) {
	teamId, ok := ctx.GetQuery("teamId")
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "teamId不能为空",
		})
		return
	}

	var projects []entity.Project
	entity.Db.Where("team_id = ? AND is_deleted = ?", teamId, false).Find(&projects)
	ctx.JSON(http.StatusOK, gin.H{

	}
}
