package api

import (
	"backend/entity"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type TemplateCreateRequest struct {
	Name      string `json:"name"`
	Type      int    `json:"type"`
	CreatorId string `json:"creatorId"`
	Intro     string `json:"intro"`
	Content   string `json:"content"`
}

type TemplateDeleteRequest struct {
	TemplateId int `json:"templateId"`
}

type TemplateRenameRequest struct {
	TemplateId int    `json:"templateId"`
	NewName    string `json:"newName"`
}

type TemplateModifyIntroRequest struct {
	TemplateId int    `json:"templateId"`
	Intro      string `json:"intro"`
}

type TemplateModifyContentRequest struct {
	TemplateId int    `json:"templateId"`
	Content    string `json:"content"`
}

func TemplateCreate(ctx *gin.Context) {
	var request TemplateCreateRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var template entity.Template
	entity.Db.Find(&template, "name = ? and type = ?", request.Name, request.Type)
	if template.TemplateId != 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "模板已存在",
		})
		return
	}
	if request.Type != 1 && request.Type != 2 && request.Type != 3 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "模板类型错误",
		})
		return
	}

	template = entity.Template{
		Name:       request.Name,
		Type:       request.Type,
		CreatorId:  request.CreatorId,
		CreateTime: time.Now(),
		Intro:      request.Intro,
		Content:    request.Content,
	}
	if template.Intro == "" {
		template.Intro = "暂无模板简介"
	}
	result := entity.Db.Create(&template)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":   "创建模板失败",
			"error": result.Error.Error(),
		})
		return
	}
	entity.Db.Where("name = ? and type = ?", request.Name, request.Type).First(&template)
	ctx.JSON(http.StatusOK, gin.H{
		"msg":        "创建模板成功",
		"templateId": template.TemplateId,
	})
}

func TemplateCompleteDelete(ctx *gin.Context) {
	var request TemplateDeleteRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var template entity.Template
	entity.Db.Find(&template, "template_id = ?", request.TemplateId)
	if template.TemplateId == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "模板不存在",
		})
		return
	}

	result := entity.Db.Where("template_id = ?", request.TemplateId).Delete(&template)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":   "删除模板失败",
			"error": result.Error.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "删除模板成功",
	})
}

func TemplateRename(ctx *gin.Context) {
	var request TemplateRenameRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var template entity.Template
	entity.Db.Find(&template, "template_id = ?", request.TemplateId)
	if template.TemplateId == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "模板不存在",
		})
		return
	}

	result := entity.Db.Model(&template).Where("template_id = ?", request.TemplateId).Update("name", request.NewName)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":   "模板重命名失败",
			"error": result.Error.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "模板重命名成功",
	})
}

func TemplateModifyIntro(ctx *gin.Context) {
	var request TemplateModifyIntroRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var template entity.Template
	entity.Db.Find(&template, "template_id = ?", request.TemplateId)
	if template.TemplateId == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "模板不存在",
		})
		return
	}

	result := entity.Db.Model(&template).Where("template_id = ?", request.TemplateId).Update("intro", request.Intro)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":   "模板修改简介失败",
			"error": result.Error.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "模板修改简介成功",
	})
}

func TemplateModifyContent(ctx *gin.Context) {
	var request TemplateModifyContentRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var template entity.Template
	entity.Db.Find(&template, "template_id = ?", request.TemplateId)
	if template.TemplateId == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "模板不存在",
		})
		return
	}

	result := entity.Db.Model(&template).Where("template_id = ?", request.TemplateId).Update("content", request.Content)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":   "模板修改内容失败",
			"error": result.Error.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "模板修改内容成功",
	})
}

func TemplateList(ctx *gin.Context) {
	templateType, ok := ctx.GetQuery("type")
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "type参数不能为空",
		})
		return
	}

	var templates []entity.Template
	var templateList []gin.H
	entity.Db.Where("type = ?", templateType).Find(&templates)
	for _, template := range templates {
		var creator entity.User
		entity.Db.Where("user_id = ?", template.CreatorId).Find(&creator)
		templateJson := gin.H{
			"templateId": template.TemplateId,
			"name":       template.Name,
			"type":       template.Type,
			"creatorId":  template.CreatorId,
			"createInfo": string(template.CreateTime.Format("2006-01-02 15:04")) + " by " + creator.Nickname,
			"intro":      template.Intro,
			"content":    template.Content,
		}
		templateList = append(templateList, templateJson)
	}
	if len(templateList) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"msg":          "当前类型没有模板",
			"templateList": make([]gin.H, 0),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":          "获取模板列表成功",
		"templateList": templateList,
	})
}

func TemplateGet(ctx *gin.Context) {
	templateId, ok := ctx.GetQuery("templateId")
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "templateId参数不能为空",
		})
		return
	}

	var template entity.Template
	entity.Db.Find(&template, "template_id = ?", templateId)
	if template.TemplateId == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "模板不存在",
		})
		return
	}
	var creator entity.User
	entity.Db.Where("user_id = ?", template.CreatorId).Find(&creator)
	ctx.JSON(http.StatusOK, gin.H{
		"msg":      "获取模板成功",
		"template": template.Content,
	})
}
