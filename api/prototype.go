package api

import (
	"backend/entity"
	"github.com/gin-gonic/gin"
	"net/http"
	"sort"
	"time"
)

type PrototypeCreateRequest struct {
	Name      string `json:"name"`
	CreatorID string `json:"creatorId"`
	ProjectID int    `json:"projectId"`
}

type PrototypeDeleteRequest struct {
	ProtoID   int    `json:"protoId"`
	DeleterID string `json:"deleterId"`
}

type PrototypeCompleteDeleteRequest struct {
	ProtoID int `json:"protoId"`
}

type ProtoRenameRequest struct {
	ProtoID int    `json:"protoId"`
	NewName string `json:"newName"`
}

type ProtoSaveRequest struct {
	ProtoID int    `json:"protoId"`
	UserId  string `json:"userId"`
	Content string `json:"content"`
}

type ProtoExitRequest struct {
	ProtoID int    `json:"protoId"`
	UserId  string `json:"userId"`
}

type ProtoApplyEditRequest struct {
	ProtoID int    `json:"protoId"`
	UserId  string `json:"userId"`
}

var protoEditorMap = make(map[int][]string)
var protoEditTimeMap = make(map[int]time.Time)
var protoUserTimeMap = make(map[string]time.Time)

func PrototypeCreate(ctx *gin.Context) {
	var request PrototypeCreateRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var prototype entity.Prototype
	entity.Db.Find(&prototype, "name = ?", request.Name)
	if prototype.ProtoID != 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "原型已存在",
		})
		return
	}

	prototype = entity.Prototype{
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
		Content:    "",
	}
	result := entity.Db.Create(&prototype)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "原型创建失败",
			"err": result.Error.Error(),
		})
		return
	}
	entity.Db.Where("name = ? AND project_id = ?", request.Name, request.ProjectID).First(&prototype)
	ctx.JSON(http.StatusOK, gin.H{
		"msg":     "原型创建成功",
		"protoId": prototype.ProtoID,
	})
}

func PrototypeDelete(ctx *gin.Context) {
	var request PrototypeDeleteRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var prototype entity.Prototype
	entity.Db.Find(&prototype, "proto_id = ?", request.ProtoID)
	if prototype.ProtoID == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "原型不存在",
		})
		return
	}
	if prototype.IsDeleted {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "原型已删除",
		})
		return
	}

	prototype.IsDeleted = true
	prototype.DeleterID = request.DeleterID
	prototype.DeleteTime = time.Now()
	result := entity.Db.Model(&prototype).Where("proto_id = ?", request.ProtoID).Updates(&prototype)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "原型删除失败",
			"err": result.Error.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "原型删除成功",
	})
}

func PrototypeCompleteDelete(ctx *gin.Context) {
	var request PrototypeCompleteDeleteRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var prototype entity.Prototype
	entity.Db.Find(&prototype, "proto_id = ?", request.ProtoID)
	if prototype.ProtoID == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "原型不存在",
		})
		return
	}

	result := entity.Db.Where("proto_id = ?", request.ProtoID).Delete(&prototype)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "原型删除失败",
			"err": result.Error.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "原型删除成功",
	})
}

func PrototypeRename(ctx *gin.Context) {
	var request ProtoRenameRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var prototype entity.Prototype
	entity.Db.Find(&prototype, "proto_id = ?", request.ProtoID)
	if prototype.ProtoID == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "原型不存在",
		})
		return
	}

	var prototypes []entity.Prototype
	entity.Db.Where("name = ? AND project_id = ?", request.NewName, prototype.ProjectID).Find(&prototypes)
	if len(prototypes) != 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "原型名称重复",
		})
		return
	}

	result := entity.Db.Model(&prototype).Where("proto_id = ?", request.ProtoID).Update("name", request.NewName)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "原型重命名失败",
			"err": result.Error.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "原型重命名成功",
	})
}

func PrototypeList(ctx *gin.Context) {
	projectId, ok := ctx.GetQuery("projectId")
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "projectId不能为空",
		})
		return
	}

	var prototypes []entity.Prototype
	var prototypeList []gin.H
	entity.Db.Where("project_id = ?", projectId).Find(&prototypes)
	sort.SliceStable(prototypes, func(i, j int) bool {
		return prototypes[i].Name < prototypes[j].Name
	})
	for _, prototype := range prototypes {
		if prototype.IsDeleted {
			continue
		}
		var creator, modifier entity.User
		entity.Db.Where("user_id = ?", prototype.CreatorID).First(&creator)
		entity.Db.Where("user_id = ?", prototype.ModifierID).First(&modifier)
		prototypeJson := gin.H{
			"protoId":    prototype.ProtoID,
			"protoName":  prototype.Name,
			"creatorId":  prototype.CreatorID,
			"createInfo": string(prototype.CreateTime.Format("2006-01-02 15:04")) + " by " + creator.Nickname,
			"modifierId": prototype.ModifierID,
			"modifyInfo": string(prototype.ModifyTime.Format("2006-01-02 15:04")) + " by " + modifier.Nickname,
		}
		prototypeList = append(prototypeList, prototypeJson)
	}
	if len(prototypeList) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"msg":       "当前项目没有原型",
			"protoList": make([]entity.Prototype, 0),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":       "获取原型列表成功",
		"protoList": prototypeList,
	})
}

func PrototypeRecycle(ctx *gin.Context) {
	projectId, ok := ctx.GetQuery("projectId")
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "projectId不能为空",
		})
		return
	}

	var prototypes []entity.Prototype
	var prototypeList []gin.H
	entity.Db.Where("project_id = ?", projectId).Find(&prototypes)
	sort.SliceStable(prototypes, func(i, j int) bool {
		return prototypes[i].Name < prototypes[j].Name
	})
	for _, prototype := range prototypes {
		if !prototype.IsDeleted {
			continue
		}
		var creator, modifier entity.User
		entity.Db.Where("user_id = ?", prototype.CreatorID).First(&creator)
		entity.Db.Where("user_id = ?", prototype.ModifierID).First(&modifier)
		prototypeJson := gin.H{
			"protoId":    prototype.ProtoID,
			"protoName":  prototype.Name,
			"creatorId":  prototype.CreatorID,
			"createInfo": string(prototype.CreateTime.Format("2006-01-02 15:04")) + " by " + creator.Nickname,
			"modifierId": prototype.ModifierID,
			"modifyInfo": string(prototype.ModifyTime.Format("2006-01-02 15:04")) + " by " + modifier.Nickname,
		}
		prototypeList = append(prototypeList, prototypeJson)
	}
	if len(prototypeList) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"msg":       "当前项目回收站没有原型",
			"protoList": make([]entity.Prototype, 0),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":       "回收站原型列表获取成功",
		"protoList": prototypeList,
	})
}

func PrototypeRecover(ctx *gin.Context) {
	var request PrototypeDeleteRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var prototype entity.Prototype
	entity.Db.Find(&prototype, "proto_id = ?", request.ProtoID)
	if prototype.ProtoID == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "原型不存在",
		})
		return
	}
	if !prototype.IsDeleted {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "原型不在回收站中",
		})
		return
	}

	result := entity.Db.Model(&prototype).Where("proto_id = ?", request.ProtoID).Update("is_deleted", false)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "原型恢复失败",
			"err": result.Error.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "原型恢复成功",
	})
}

func PrototypeSave(ctx *gin.Context) {
	var request ProtoSaveRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var prototype entity.Prototype
	entity.Db.Find(&prototype, "proto_id = ?", request.ProtoID)
	if prototype.ProtoID == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "原型不存在",
		})
		return
	}

	prototype.Content = request.Content
	prototype.ModifierID = request.UserId
	prototype.ModifyTime = time.Now()
	result := entity.Db.Model(&prototype).Where("proto_id = ?", request.ProtoID).Updates(&prototype)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "原型保存失败",
			"err": result.Error.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "原型保存成功",
	})
}

func PrototypeExit(ctx *gin.Context) {
	var request ProtoExitRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var prototype entity.Prototype
	entity.Db.Find(&prototype, "proto_id = ?", request.ProtoID)
	if prototype.ProtoID == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "原型不存在",
		})
		return
	}
	if prototype.IsDeleted {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "原型在回收站中",
		})
		return
	}

	if time.Now().Sub(protoEditTimeMap[request.ProtoID]) > time.Second*3 {
		protoEditorMap[request.ProtoID] = []string{}
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "原型不在编辑状态",
		})
		return
	}

	editors := protoEditorMap[request.ProtoID]
	var nowEditors []string
	for _, editor := range editors {
		if editor == request.UserId || time.Now().Sub(protoUserTimeMap[editor]) > time.Second*3 {
			continue
		}
		nowEditors = append(nowEditors, editor)
	}

	prototype.ModifierID = request.UserId
	prototype.ModifyTime = time.Now()
	result := entity.Db.Model(&prototype).Where("proto_id = ?", request.ProtoID).Updates(&prototype)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "原型退出编辑失败",
			"err": result.Error.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":    "原型退出编辑成功",
		"remain": len(nowEditors),
	})
}

func PrototypeGet(ctx *gin.Context) {
	protoId, ok := ctx.GetQuery("protoId")
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "protoId不能为空",
		})
		return
	}

	var prototype entity.Prototype
	entity.Db.Find(&prototype, "proto_id = ?", protoId)
	if prototype.ProtoID == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "原型不存在",
		})
		return
	}
	if prototype.IsDeleted {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "原型在回收站中, 无法编辑",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":     "获取原型成功",
		"content": prototype.Content,
	})
}

func PrototypeApplyEdit(ctx *gin.Context) {
	var request ProtoApplyEditRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var prototype entity.Prototype
	entity.Db.Find(&prototype, "proto_id = ?", request.ProtoID)
	if prototype.ProtoID == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "原型不存在",
		})
		return
	}
	if prototype.IsDeleted {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "原型在回收站中, 无法编辑",
		})
		return
	}

	editors := protoEditorMap[request.ProtoID]
	protoUserTimeMap[request.UserId] = time.Now()
	var nowEditors []string
	for _, editor := range editors {
		if editor == request.UserId || time.Now().Sub(protoUserTimeMap[editor]) > time.Second*3 {
			continue
		}
		nowEditors = append(nowEditors, editor)
	}

	nowEditors = append(nowEditors, request.UserId)
	protoEditorMap[request.ProtoID] = nowEditors
	protoEditTimeMap[request.ProtoID] = time.Now()
	ctx.JSON(http.StatusOK, gin.H{
		"msg":          "原型申请编辑成功",
		"nowEditorNum": len(nowEditors),
	})
}
