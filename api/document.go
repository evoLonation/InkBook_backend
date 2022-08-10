package api

import (
	"backend/entity"
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"time"
)

type DocumentCreateRequest struct {
	Name       string `json:"name"`
	CreatorId  string `json:"creatorId"`
	TeamId     int    `json:"teamId"`
	ParentId   int    `json:"parentId"`
	TemplateId int    `json:"templateId"`
}

type DocumentProjectCreateRequest struct {
	Name       string `json:"name"`
	CreatorId  string `json:"creatorId"`
	ProjectId  int    `json:"projectId"`
	TemplateId int    `json:"templateId"`
}

type DocumentDeleteRequest struct {
	DocId     int    `json:"docId"`
	DeleterId string `json:"deleterId"`
}

type DocumentCompleteDeleteRequest struct {
	DocId int `json:"docId"`
}

type DocumentRenameRequest struct {
	DocId   int    `json:"docId"`
	NewName string `json:"newName"`
}

type DocumentSaveRequest struct {
	DocId   int    `json:"docId"`
	UserId  string `json:"userId"`
	Content string `json:"content"`
}

type DocumentExitRequest struct {
	DocId  int    `json:"docId"`
	UserId string `json:"userId"`
}

type DocumentApplyEditRequest struct {
	DocId  int    `json:"docId"`
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
	if document.DocId != 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "文档已存在",
		})
		return
	}

	var template entity.Template
	var content string
	if request.TemplateId != 0 {
		entity.Db.Find(&template, "template_id = ?", request.TemplateId)
		if template.TemplateId == 0 || template.Type != 1 {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"msg": "模板不存在",
			})
			return
		}
		content = template.Content
	} else {
		content = ""
	}

	document = entity.Document{
		Name:       request.Name,
		TeamId:     request.TeamId,
		ParentId:   request.ParentId,
		CreatorId:  request.CreatorId,
		CreateTime: time.Now(),
		ModifierId: request.CreatorId,
		ModifyTime: time.Now(),
		IsEditing:  false,
		IsDeleted:  false,
		DeleterId:  request.CreatorId,
		DeleteTime: time.Now(),
		Content:    content,
	}
	result := entity.Db.Create(&document)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":   "文档创建失败",
			"error": result.Error.Error(),
		})
		return
	}
	entity.Db.Where("name = ? and team_id = ?", request.Name, request.TeamId).First(&document)
	ctx.JSON(http.StatusOK, gin.H{
		"msg":   "文档创建成功",
		"docId": document.DocId,
	})
}

func DocumentProjectCreate(ctx *gin.Context) {
	var request DocumentProjectCreateRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var project entity.Project
	entity.Db.Find(&project, "project_id = ?", request.ProjectId)
	if project.ProjectId == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "项目不存在",
		})
		return
	}
	if project.IsDeleted {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "项目已删除",
		})
		return
	}

	var folder entity.Folder
	entity.Db.Find(&folder, "name = ? AND team_id = ?", project.Name+"的项目文档", project.TeamId)
	if folder.FolderId == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "项目文档文件夹不存在",
		})
		return
	}
	if folder.IsDeleted {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "项目文档文件夹已删除",
		})
		return
	}

	var document entity.Document
	entity.Db.Find(&document, "name = ? and team_id = ?", request.Name, project.TeamId)
	if document.DocId != 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "文档已存在",
		})
		return
	}

	var template entity.Template
	var content string
	if request.TemplateId != 0 {
		entity.Db.Find(&template, "template_id = ?", request.TemplateId)
		if template.TemplateId == 0 || template.Type != 1 {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"msg": "模板不存在",
			})
			return
		}
		content = template.Content
	} else {
		content = ""
	}

	document = entity.Document{
		Name:       request.Name,
		TeamId:     project.TeamId,
		ParentId:   folder.FolderId,
		CreatorId:  request.CreatorId,
		CreateTime: time.Now(),
		ModifierId: request.CreatorId,
		ModifyTime: time.Now(),
		IsEditing:  false,
		IsDeleted:  false,
		DeleterId:  request.CreatorId,
		DeleteTime: time.Now(),
		Content:    content,
	}
	result := entity.Db.Create(&document)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":   "文档创建失败",
			"error": result.Error.Error(),
		})
		return
	}
	entity.Db.Where("name = ? and team_id = ?", request.Name, project.TeamId).First(&document)
	ctx.JSON(http.StatusOK, gin.H{
		"msg":   "文档创建成功",
		"docId": document.DocId,
	})
}

func DocumentDelete(ctx *gin.Context) {
	var request DocumentDeleteRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var document entity.Document
	entity.Db.Find(&document, "doc_id = ?", request.DocId)
	if document.DocId == 0 {
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
	document.DeleterId = request.DeleterId
	document.DeleteTime = time.Now()
	result := entity.Db.Model(&document).Where("doc_id = ?", request.DocId).Updates(&document)
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
	entity.Db.Find(&document, "doc_id = ?", request.DocId)
	if document.DocId == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "文档不存在",
		})
		return
	}

	result := entity.Db.Where("doc_id = ?", request.DocId).Delete(&document)
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
	entity.Db.Find(&document, "doc_id = ?", request.DocId)
	if document.DocId == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "文档不存在",
		})
		return
	}

	var documents []entity.Document
	entity.Db.Where("name = ? and parent_id = ?", request.NewName, document.ParentId).Find(&documents)
	if len(documents) != 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "文档名称重复",
		})
		return
	}

	result := entity.Db.Model(&document).Where("doc_id = ?", request.DocId).Update("name", request.NewName)
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

	var project entity.Project
	entity.Db.Find(&project, "project_id = ?", projectId)
	if project.ProjectId == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "项目不存在",
		})
		return
	}

	var folder entity.Folder
	entity.Db.Find(&folder, "team_id = ? and name = ?", project.TeamId, project.Name+"的项目文档")
	if folder.FolderId == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "项目文件夹不存在",
		})
		return
	}

	var documents []entity.Document
	var docList []gin.H
	entity.Db.Where("parent_id = ?", folder.FolderId).Find(&documents)
	sort.SliceStable(documents, func(i, j int) bool {
		return documents[i].CreateTime.Unix() > documents[i].CreateTime.Unix()
	})
	for _, document := range documents {
		if document.IsDeleted {
			continue
		}
		var creator, modifier entity.User
		entity.Db.Where("user_id = ?", document.CreatorId).Find(&creator)
		entity.Db.Where("user_id = ?", document.ModifierId).Find(&modifier)
		documentJson := gin.H{
			"docId":      document.DocId,
			"docName":    document.Name,
			"creatorId":  document.CreatorId,
			"createInfo": string(document.CreateTime.Format("2006-01-02 15:04")) + " by " + creator.Nickname,
			"modifierId": document.ModifierId,
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
	teamId, ok := ctx.GetQuery("teamId")
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "teamId不能为空",
		})
		return
	}

	var documents []entity.Document
	var docList []gin.H
	entity.Db.Where("team_id = ?", teamId).Find(&documents)
	sort.SliceStable(documents, func(i, j int) bool {
		return documents[i].CreateTime.Unix() > documents[i].CreateTime.Unix()
	})
	for _, document := range documents {
		if !document.IsDeleted {
			continue
		}
		var creator, modifier entity.User
		entity.Db.Where("user_id = ?", document.CreatorId).Find(&creator)
		entity.Db.Where("user_id = ?", document.ModifierId).Find(&modifier)
		documentJson := gin.H{
			"docId":      document.DocId,
			"docName":    document.Name,
			"creatorId":  document.CreatorId,
			"createInfo": string(document.CreateTime.Format("2006-01-02 15:04")) + " by " + creator.Nickname,
			"modifierId": document.ModifierId,
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
	entity.Db.Find(&document, "doc_id = ?", request.DocId)
	if document.DocId == 0 {
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

	result := entity.Db.Model(&document).Where("doc_id = ?", request.DocId).Update("is_deleted", false)
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
	entity.Db.Find(&document, "doc_id = ?", request.DocId)
	if document.DocId == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "文档不存在",
		})
		return
	}

	document.Content = request.Content
	document.ModifierId = request.UserId
	document.ModifyTime = time.Now()
	result := entity.Db.Model(&document).Where("doc_id = ?", request.DocId).Updates(&document)
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
	entity.Db.Find(&document, "doc_id = ?", request.DocId)
	if document.DocId == 0 {
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

	if time.Now().Sub(docEditTimeMap[request.DocId]) > time.Second*3 {
		docEditorMap[request.DocId] = []string{}
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "文档不在编辑状态",
		})
		return
	}

	editors := docEditorMap[request.DocId]
	var nowEditors []string
	for _, editor := range editors {
		if editor == request.UserId || time.Now().Sub(docUserTimeMap[editor]) > time.Second*3 {
			continue
		}
		nowEditors = append(nowEditors, editor)
	}

	document.ModifierId = request.UserId
	document.ModifyTime = time.Now()
	result := entity.Db.Model(&document).Where("doc_id = ?", request.DocId).Updates(&document)
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
	if document.DocId == 0 {
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

	ctx.JSON(http.StatusOK, gin.H{
		"msg":     "文档获取成功",
		"content": document.Content,
	})
}

func DocumentApplyEdit(ctx *gin.Context) {
	var request DocumentApplyEditRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var document entity.Document
	entity.Db.Find(&document, "doc_id = ?", request.DocId)
	if document.DocId == 0 {
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

	editors := docEditorMap[request.DocId]
	docUserTimeMap[request.UserId] = time.Now()
	var nowEditors []string
	for _, editor := range editors {
		if editor == request.UserId || time.Now().Sub(docUserTimeMap[editor]) > time.Second*3 {
			continue
		}
		nowEditors = append(nowEditors, editor)
	}

	nowEditors = append(nowEditors, request.UserId)
	docEditorMap[request.DocId] = nowEditors
	docEditTimeMap[request.DocId] = time.Now()
	ctx.JSON(http.StatusOK, gin.H{
		"msg":          "申请编辑状态成功",
		"nowEditorNum": len(nowEditors),
		"editorList":   nowEditors,
	})
}
func getMd5(file []byte) string {
	h := md5.New()
	h.Write(file)
	bytesBuffer := bytes.NewBuffer(h.Sum(nil))
	var x int64
	binary.Read(bytesBuffer, binary.BigEndian, &x)
	if x < 0 {
		x *= -1
	}
	return strconv.FormatInt(x, 16)
}
func DocumentImg(c *gin.Context) {
	file, header, _ := c.Request.FormFile("file")
	filename := header.Filename
	fileContent, _ := header.Open()
	byteContainer, err := ioutil.ReadAll(fileContent)
	p := getMd5(byteContainer)
	fmt.Println(p)
	out, err := os.Create("./localFile/document/" + p + filename)
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
	c.JSON(http.StatusOK, gin.H{"url": "http://43.138.71.108/api/url/localFile/document/" + p + filename})
}

func Url(c *gin.Context) {
	url := c.Param("url")
	//c.JSON(http.StatusOK, gin.H{"url": url})
	c.File("./localFile/document/" + url)
}
