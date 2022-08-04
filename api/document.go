package api

import (
	"backend/entity"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"sort"
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

type DocumentRenameRequest struct {
	DocID   int    `json:"docId"`
	NewName string `json:"newName"`
}

type DocumentListRequest struct {
	ProjectID int `json:"projectId"`
}

type DocumentSaveRequest struct {
	DocID   int    `json:"docId"`
	UserId  string `json:"userId"`
	Content gin.H  `json:"content"`
}

type DocumentExitRequest struct {
	DocID  int    `json:"docId"`
	UserId string `json:"userId"`
}

func DocumentCreate(ctx *gin.Context) {
	var request DocumentCreateRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var document entity.Document
	entity.Db.Find(&document, "name = ?", request.Name)
	if document.DocID != 0 {
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
		Content:    "{}",
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
	if document.DocID == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "文档不存在",
		})
		return
	}
	if document.IsDeleted {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "文档已删除",
		})
		return
	}

	document.IsDeleted = true
	document.DeleterID = request.DeleterID
	document.DeleteTime = time.Now()
	result := entity.Db.Model(&document).Where("doc_id = ?", request.DocID).Updates(&document)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":   "文档删除失败",
			"error": result.Error.Error(),
		})
		return
	}
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
	if document.DocID == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "文档不存在",
		})
		return
	}

	result := entity.Db.Where("doc_id = ?", request.DocID).Delete(&document)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":   "文档删除失败",
			"error": result.Error.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "文档删除成功",
	})
}

func DocumentRename(ctx *gin.Context) {
	var request DocumentRenameRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var document entity.Document
	entity.Db.Find(&document, "doc_id = ?", request.DocID)
	if document.DocID == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "文档不存在",
		})
		return
	}

	var documents []entity.Document
	entity.Db.Where("name = ? and project_id = ?", request.NewName, document.ProjectID).Find(&documents)
	if len(documents) != 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "文档名称重复",
		})
		return
	}

	result := entity.Db.Model(&document).Where("doc_id = ?", request.DocID).Update("name", request.NewName)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":   "文档重命名失败",
			"error": result.Error.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "文档重命名成功",
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
	sort.SliceStable(documents, func(i, j int) bool {
		return documents[i].Name < documents[j].Name
	})
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
	sort.SliceStable(documents, func(i, j int) bool {
		return documents[i].Name < documents[j].Name
	})
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
	if document.DocID == 0 {
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

func DocumentSave(ctx *gin.Context) {
	var request DocumentSaveRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var document entity.Document
	entity.Db.Find(&document, "doc_id = ?", request.DocID)
	if document.DocID == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "文档不存在",
		})
		return
	}

	fmt.Println(request.Content)
	jsonContent, err := json.Marshal(request.Content)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"msg":   "JSON格式内容解析失败",
		})
		return
	}
	document.Content = string(jsonContent)
	document.ModifierID = request.UserId
	document.ModifyTime = time.Now()
	result := entity.Db.Model(&document).Where("doc_id = ?", request.DocID).Updates(&document)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "文档保存失败",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "文档保存成功",
	})
}

func DocumentExit(ctx *gin.Context) {
	var request DocumentExitRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var document entity.Document
	entity.Db.Find(&document, "doc_id = ?", request.DocID)
	if document.DocID == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "文档不存在",
		})
		return
	}
	if document.IsDeleted {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "文档在回收站中",
		})
		return
	}
	if !document.IsEditing {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "文档不在编辑状态",
		})
		return
	}

	result := entity.Db.Model(&document).Where("doc_id = ?", request.DocID).Update("editing_cnt", document.EditingCnt-1)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "文档退出编辑失败",
		})
		return
	}
	if document.EditingCnt == 0 {
		result = entity.Db.Model(&document).Where("doc_id = ?", request.DocID).Update("is_editing", false)
		if result.Error != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "文档退出编辑失败",
			})
			return
		}
	}

	document.ModifierID = request.UserId
	document.ModifyTime = time.Now()
	result = entity.Db.Model(&document).Where("doc_id = ?", request.DocID).Updates(&document)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "文档退出编辑失败",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":    "文档退出编辑成功",
		"remain": document.EditingCnt,
	})
}

func DocumentGet(ctx *gin.Context) {
	docId, ok := ctx.GetQuery("docId")
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "docId不能为空",
		})
		return
	}

	var document entity.Document
	entity.Db.Find(&document, "doc_id = ?", docId)
	if document.DocID == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "文档不存在",
		})
		return
	}
	if document.IsDeleted {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "文档在回收站中，无法编辑",
		})
		return
	}

	if document.IsEditing {
		result := entity.Db.Model(&document).Where("doc_id = ?", docId).Update("editing_cnt", document.EditingCnt+1)
		if result.Error != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "文档获取失败",
			})
			return
		}
		ctx.JSON(http.StatusConflict, gin.H{
			"msg": "文档正在编辑",
		})
		return
	}

	result := entity.Db.Model(&document).Where("doc_id = ?", docId).Update("is_editing", true)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "文档获取失败",
		})
		return
	}
	result = entity.Db.Model(&document).Where("doc_id = ?", docId).Update("editing_cnt", document.EditingCnt+1)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "文档获取失败",
		})
		return
	}

	var jsonContent gin.H
	if err := json.Unmarshal([]byte(document.Content), &jsonContent); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "JSON格式内容解析失败",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":     "文档获取成功",
		"content": jsonContent,
	})
}
