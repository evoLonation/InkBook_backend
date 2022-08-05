package api

import (
	"backend/entity"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"os"
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

type DocumentSaveRequest struct {
	DocID   int    `json:"docId"`
	UserId  string `json:"userId"`
	Content string `json:"content"`
}

type DocumentExitRequest struct {
	DocID  int    `json:"docId"`
	UserId string `json:"userId"`
}

type DocumentApplyEditRequest struct {
	DocID  int    `json:"docId"`
	UserId string `json:"userId"`
}

var docEditorMap = make(map[int][]string)
var docEditTimeMap = make(map[int]time.Time)
var docUserTimeMap = make(map[string]time.Time)

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
			"msg": "文档已存在",
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
		Content:    "{\"content\": {}}",
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
			"msg": "文档不存在",
		})
		return
	}
	if document.IsDeleted {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "文档已删除",
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
			"msg": "文档不存在",
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
			"msg": "文档不存在",
		})
		return
	}

	var documents []entity.Document
	entity.Db.Where("name = ? and project_id = ?", request.NewName, document.ProjectID).Find(&documents)
	if len(documents) != 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "文档名称重复",
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
			"msg": "projectId不能为空",
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
		var creator, modifier entity.User
		entity.Db.Where("user_id = ?", document.CreatorID).Find(&creator)
		entity.Db.Where("user_id = ?", document.ModifierID).Find(&modifier)
		documentJson := gin.H{
			"docId":      document.DocID,
			"docName":    document.Name,
			"creatorId":  document.CreatorID,
			"createInfo": string(document.CreateTime.Format("2006-01-02 15:04")) + " by " + creator.Nickname,
			"modifierId": document.ModifierID,
			"modifyInfo": string(document.ModifyTime.Format("2006-01-02 15:04")) + " by " + modifier.Nickname,
		}
		docList = append(docList, documentJson)
	}
	if len(docList) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"msg":     "当前项目没有文档",
			"docList": make([]entity.Document, 0),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":     "文档列表获取成功",
		"docList": docList,
	})
}

func DocumentRecycle(ctx *gin.Context) {
	projectId, ok := ctx.GetQuery("projectId")
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "projectId不能为空",
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
		var creator, modifier entity.User
		entity.Db.Where("user_id = ?", document.CreatorID).Find(&creator)
		entity.Db.Where("user_id = ?", document.ModifierID).Find(&modifier)
		documentJson := gin.H{
			"docId":      document.DocID,
			"docName":    document.Name,
			"creatorId":  document.CreatorID,
			"createInfo": string(document.CreateTime.Format("2006-01-02 15:04")) + " by " + creator.Nickname,
			"modifierId": document.ModifierID,
			"modifyInfo": string(document.ModifyTime.Format("2006-01-02 15:04")) + " by " + modifier.Nickname,
		}
		docList = append(docList, documentJson)
	}
	if len(docList) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"msg":     "当前回收站没有文档",
			"docList": make([]entity.Document, 0),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":     "回收站文档列表获取成功",
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
			"msg": "文档不存在",
		})
		return
	}
	if !document.IsDeleted {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "文档不在回收站中",
		})
		return
	}

	result := entity.Db.Model(&document).Where("doc_id = ?", request.DocID).Update("is_deleted", false)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":   "文档恢复失败",
			"error": result.Error.Error(),
		})
		return
	}
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
			"msg": "文档不存在",
		})
		return
	}

	//jsonContent, err := json.Marshal(request.Content)
	//if err != nil {
	//	ctx.JSON(http.StatusBadRequest, gin.H{
	//		"error": err.Error(),
	//		"msg":   "JSON格式内容解析失败",
	//	})
	//	return
	//}
	document.Content = "{\"content\":" + string(request.Content) + "}"
	document.ModifierID = request.UserId
	document.ModifyTime = time.Now()
	result := entity.Db.Model(&document).Where("doc_id = ?", request.DocID).Updates(&document)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "文档保存失败",
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
			"msg": "文档不存在",
		})
		return
	}
	if document.IsDeleted {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "文档在回收站中",
		})
		return
	}

	if time.Now().Sub(docEditTimeMap[request.DocID]) > time.Second*3 {
		docEditorMap[request.DocID] = []string{}
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "文档不在编辑状态",
		})
		return
	}

	editors := docEditorMap[request.DocID]
	var nowEditors []string
	for _, editor := range editors {
		if editor == request.UserId || time.Now().Sub(docUserTimeMap[editor]) > time.Second*3 {
			continue
		}
		nowEditors = append(nowEditors, editor)
	}

	document.ModifierID = request.UserId
	document.ModifyTime = time.Now()
	result := entity.Db.Model(&document).Where("doc_id = ?", request.DocID).Updates(&document)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "文档退出编辑失败",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":    "文档退出编辑成功",
		"remain": len(nowEditors),
	})
}

func DocumentGet(ctx *gin.Context) {
	docId, ok := ctx.GetQuery("docId")
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "docId不能为空",
		})
		return
	}

	var document entity.Document
	entity.Db.Find(&document, "doc_id = ?", docId)
	if document.DocID == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "文档不存在",
		})
		return
	}
	if document.IsDeleted {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "文档在回收站中，无法编辑",
		})
		return
	}

	var jsonContent gin.H
	if err := json.Unmarshal([]byte(document.Content), &jsonContent); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "JSON格式内容解析失败",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":     "文档获取成功",
		"content": jsonContent["content"],
	})
}

func DocumentApplyEdit(ctx *gin.Context) {
	var request DocumentApplyEditRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var document entity.Document
	entity.Db.Find(&document, "doc_id = ?", request.DocID)
	if document.DocID == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "文档不存在",
		})
		return
	}
	if document.IsDeleted {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "文档已被删除",
		})
		return
	}

	editors := docEditorMap[request.DocID]
	docUserTimeMap[request.UserId] = time.Now()
	var nowEditors []string
	for _, editor := range editors {
		if editor == request.UserId || time.Now().Sub(docUserTimeMap[editor]) > time.Second*3 {
			continue
		}
		nowEditors = append(nowEditors, editor)
	}

	nowEditors = append(nowEditors, request.UserId)
	docEditorMap[request.DocID] = nowEditors
	docEditTimeMap[request.DocID] = time.Now()
	ctx.JSON(http.StatusOK, gin.H{
		"msg":          "申请编辑状态成功",
		"nowEditorNum": len(nowEditors),
	})
}

func DocumentImg(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	filename := header.Filename
	out, err := os.Create("./localFile/document/" + filename)
	if err != nil {
		log.Fatal(err)
	}
	defer func(out *os.File) {
		err := out.Close()
		if err != nil {

		}
	}(out)
	_, err = io.Copy(out, file)
	if err != nil {
		log.Fatal(err)
	}
	c.JSON(http.StatusOK, gin.H{"url": "http://43.138.71.108/api/url/" + "./localFile/document/" + filename})
}

func Url(c *gin.Context) {
	url := c.Param("url")
	//c.JSON(http.StatusOK, gin.H{"url": url})
	c.File("./localFile/document/" + url)
}
