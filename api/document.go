package api

import (
	"backend/entity"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type DocumentCreateRequest struct {
	Name      string `json:"name"`
	CreatorID string `json:"creatorId"`
	ProjectID int    `json:"projectId"`
}

type DocumentDeleteRequest struct {
	DocID     int    `json:"docId"`
	DeleterID string `json:"deleterId"`
}

type DocumentCompleteDeleteRequest struct {
	DocID int `json:"docId"`
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
		DeleterID:  request.CreatorID,
		DeleteTime: time.Now(),
	}
	result := entity.Db.Create(&document)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":   "文档创建失败",
			"error": result.Error.Error(),
		})
		return
	}
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

func DocumentCompleteDelete(ctx *gin.Context) {
	var request DocumentCompleteDeleteRequest
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

	entity.Db.Where("doc_id = ?", request.DocID).Delete(&document)
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "文档删除成功",
	})
}

func DocumentList(ctx *gin.Context) {
	projectId, ok := ctx.GetQuery("projectId")
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "projectId不能为空",
		})
		return
	}

	var documents []entity.Document
	var docList []gin.H
	entity.Db.Where("project_id = ?", projectId).Find(&documents)
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
			"createInfo": string(document.CreateTime.Format("2006-01-02 15:04")) + " by " + creator.Nickname,
		}
		docList = append(docList, documentJson)
	}
	if len(docList) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"msg": "当前项目没有文档",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"docList": docList,
	})
}

func DocumentRecycle(ctx *gin.Context) {
	projectId, ok := ctx.GetQuery("projectId")
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "projectId不能为空",
		})
		return
	}

	var documents []entity.Document
	var docList []gin.H
	entity.Db.Where("project_id = ?", projectId).Find(&documents)
	for _, document := range documents {
		if !document.IsDeleted {
			continue
		}
		var creator entity.User
		entity.Db.Where("user_id = ?", document.CreatorID).Find(&creator)
		documentJson := gin.H{
			"docId":      document.DocID,
			"docName":    document.Name,
			"creatorId":  document.CreatorID,
			"createInfo": string(document.CreateTime.Format("2006-01-02 15:04")) + " by " + creator.Nickname,
		}
		docList = append(docList, documentJson)
	}
	if len(docList) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"msg": "当前回收站没有文档",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"docList": docList,
	})
}

func DocumentRecover(ctx *gin.Context) {
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
	if !document.IsDeleted {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "文档不在回收站中",
		})
		return
	}

	entity.Db.Model(&document).Where("doc_id = ?", request.DocID).Update("is_deleted", false)
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "文档恢复成功",
	})
}
