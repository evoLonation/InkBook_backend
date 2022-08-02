package api

import (
	"backend/entity"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type DocumentCreateRequest struct {
	Name      string `json:"name"`
	CreatorID int    `json:"creatorId"`
	ProjectID int    `json:"projectId"`
}

type DocumentDeleteRequest struct {
	DocID     int `json:"docId"`
	DeleterID int `json:"deleterId"`
}

type DocumentListRequest struct {
	ProjectID int `json:"projectId"`
}

func DocumentCreate(ctx *gin.Context) {
	var request DocumentCreateRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var document entity.Document
	entity.Db.Find(&document, "name = ?", request.Name)
	if document != (entity.Document{}) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "文档已存在",
		})
		return
	}

	document = entity.Document{
		Name:       request.Name,
		ProjectID:  request.ProjectID,
		CreatorID:  request.CreatorID,
		CreateTime: time.Now(),
		ModifierID: request.CreatorID,
		ModifyTime: time.Now(),
		IsEditing:  false,
		IsDeleted:  false,
		DeleterID:  request.ProjectID,
	}
	entity.Db.Create(&document)
	entity.Db.Where("name = ? AND project_id = ?", request.Name, request.ProjectID).First(&document)
	ctx.JSON(http.StatusOK, gin.H{
		"msg":   "文档创建成功",
		"docId": document.DocID,
	})
}

func DocumentDelete(ctx *gin.Context) {
	var request DocumentDeleteRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var document entity.Document
	entity.Db.Find(&document, "doc_id = ?", request.DocID)
	if document == (entity.Document{}) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "文档不存在",
		})
		return
	}

	document.IsDeleted = true
	document.DeleterID = request.DeleterID
	document.DeleteTime = time.Now()
	entity.Db.Model(&document).Where("doc_id = ?", request.DocID).Updates(&document)
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "文档删除成功",
	})
}

func DocumentList(ctx *gin.Context) {
	var request DocumentListRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var documents []entity.Document
	var docList []gin.H
	entity.Db.Where("project_id = ?", request.ProjectID).Find(&documents)
	for _, document := range documents {
		if document.IsDeleted {
			continue
		}
		var creator entity.User
		entity.Db.Where("user_id = ?", document.CreatorID).Find(&creator)
		documentJson := gin.H{
			"docId":      document.DocID,
			"docName":    document.Name,
			"creatorId":  document.CreatorID,
			"createInfo": document.CreateTime.Format("2022-01-01 00:00") + " by " + creator.Nickname,
		}
		docList = append(docList, documentJson)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"docList": docList,
	})
}
