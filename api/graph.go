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

type GraphDeleteRequest struct {
	GraphID   int    `json:"graphId"`
	DeleterID string `json:"deleterId"`
}

type GraphCompleteDeleteRequest struct {
	GraphID int `json:"graphId"`
}

type GraphRenameRequest struct {
	GraphID int    `json:"graphId"`
	NewName string `json:"newName"`
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
	entity.Db.Where("name = ? and project_id = ?", request.Name, request.ProjectID).First(&graph)
	ctx.JSON(http.StatusOK, gin.H{
		"msg":     "UML图创建成功",
		"graphId": graph.GraphID,
	})
}

func GraphDelete(ctx *gin.Context) {
	var request GraphDeleteRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var graph entity.Graph
	entity.Db.Find(&graph, "graph_id = ?", request.GraphID)
	if graph == (entity.Graph{}) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "UML图不存在",
		})
		return
	}
	if graph.IsDeleted {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "UML图已删除",
		})
		return
	}

	graph.IsDeleted = true
	graph.DeleterID = request.DeleterID
	graph.DeleteTime = time.Now()
	result := entity.Db.Model(&graph).Where("graph_id = ?", request.GraphID).Updates(&graph)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":   "UML图删除失败",
			"error": result.Error.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "UML图删除成功",
	})
}

func GraphCompleteDelete(ctx *gin.Context) {

}

func GraphRename(ctx *gin.Context) {

}

func GraphList(ctx *gin.Context) {

}

func GraphRecycle(ctx *gin.Context) {

}

func GraphRecover(ctx *gin.Context) {

}

func GraphSave(ctx *gin.Context) {

}

func GraphExit(ctx *gin.Context) {

}

func GraphGet(ctx *gin.Context) {

}
