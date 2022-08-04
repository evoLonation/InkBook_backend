package api

import (
	"backend/entity"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type GraphCreateRequest struct {
	Name      string `json:"name"`
	CreatorID string `json:"creatorId"`
	ProjectID int    `json:"projectId"`
}

func GraphCreate(ctx *gin.Context) {
	var request GraphCreateRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var graph entity.Graph
	entity.Db.Find(&graph, "name = ?", request.Name)
	if graph != (entity.Graph{}) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "UML图已存在",
		})
		return
	}

	graph = entity.Graph{
		Name:       request.Name,
		ProjectID:  request.ProjectID,
		CreatorID:  request.CreatorID,
		CreateTime: time.Now(),
		ModifierID: request.CreatorID,
		ModifyTime: time.Now(),
		IsEditing:  false,
		IsDeleted:  false,
		DeleterID:  request.CreatorID,
		DeleteTime: time.Now(),
		Content:    "{}",
		EditingCnt: 0,
	}
	result := entity.Db.Create(&graph)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":   "UML图创建失败",
			"error": result.Error.Error(),
		})
		return
	}
	entity.Db.Where("name = ?, and project_id = ?", request.Name, request.ProjectID).First(&graph)
	ctx.JSON(http.StatusOK, gin.H{
		"msg":     "UML图创建成功",
		"graphId": graph.GraphID,
	})
}
