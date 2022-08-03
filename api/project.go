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

type ProjectModifyIntroRequest struct {
	ProjectID int    `json:"projectId"`
	NewIntro  string `json:"newIntro"`
}

type ProjectModifyImgRequest struct {
	ProjectID int    `json:"projectId"`
	NewImgURL string `json:"newImgUrl"`
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
	if project.Intro == "" {
		project.Intro = "暂无项目简介"
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
	if project.IsDeleted {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "项目已删除",
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

func ProjectCompleteDelete(ctx *gin.Context) {
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

	entity.Db.Where("project_id = ?", request.ProjectID).Delete(&project)
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

	var projects []entity.Project
	entity.Db.Where("name = ? AND team_id = ?", request.NewName, project.TeamID).Find(&projects)
	if len(projects) != 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "项目名称重复",
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
	var projectList []gin.H
	entity.Db.Where("team_id = ?", teamId).Find(&projects)
	for _, project := range projects {
		if project.IsDeleted {
			continue
		}
		projectJson := gin.H{
			"id":     project.ProjectID,
			"name":   project.Name,
			"detail": project.Intro,
			"imgUrl": project.ImgURL,
		}
		projectList = append(projectList, projectJson)
	}
	ctx.JSON(http.StatusOK, gin.H{
		"projects": projectList,
	})
}

func ProjectModifyIntro(ctx *gin.Context) {
	var request ProjectModifyIntroRequest
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

	project.Intro = request.NewIntro
	entity.Db.Model(&project).Where("project_id = ?", request.ProjectID).Updates(&project)
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "项目简介修改成功",
	})
}

func ProjectModifyImg(ctx *gin.Context) {
	var request ProjectModifyImgRequest
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

	project.ImgURL = request.NewImgURL
	entity.Db.Model(&project).Where("project_id = ?", request.ProjectID).Updates(&project)
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "项目图片修改成功",
	})
}
