package api

import (
	"backend/entity"
	"github.com/gin-gonic/gin"
	"net/http"
	"sort"
	"strings"
	"time"
)

type FolderCreateRequest struct {
	Name      string `json:"name"`
	CreatorId string `json:"creatorId"`
	TeamId    int    `json:"teamId"`
	ParentId  int    `json:"parentId"`
}

type FolderDeleteRequest struct {
	FolderId  int    `json:"folderId"`
	DeleterId string `json:"deleterId"`
}

type FolderRenameRequest struct {
	FolderId int    `json:"folderId"`
	NewName  string `json:"newName"`
}

type FolderMoveRequest struct {
	DocId    int `json:"docId"`
	FolderId int `json:"folderId"`
}

func FolderCreate(ctx *gin.Context) {
	var request FolderCreateRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var folder entity.Folder
	entity.Db.Find(&folder, "name = ? and team_id = ?", request.Name, request.TeamId)
	if folder.FolderId != 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "文件夹已存在",
		})
		return
	}
	if strings.HasSuffix(request.Name, "的项目文档") {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "以“的项目文档”结尾是项目文档文件夹，自定义文件夹名不能以其结尾",
		})
		return
	}

	folder = entity.Folder{
		Name:       request.Name,
		TeamId:     request.TeamId,
		ParentId:   request.ParentId,
		CreatorId:  request.CreatorId,
		CreateTime: time.Now(),
		IsDeleted:  false,
		DeleterId:  request.CreatorId,
		DeleteTime: time.Now(),
	}
	result := entity.Db.Create(&folder)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":   "文件夹创建失败",
			"error": result.Error.Error(),
		})
		return
	}
	entity.Db.Where("name = ? and team_id = ?", request.Name, request.TeamId).First(&folder)
	ctx.JSON(http.StatusOK, gin.H{
		"msg":      "文件夹创建成功",
		"folderId": folder.FolderId,
	})
}

func FolderDelete(ctx *gin.Context) {
	var request FolderDeleteRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var folder entity.Folder
	entity.Db.Find(&folder, "folder_id = ?", request.FolderId)
	if folder.FolderId == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "文件夹不存在",
		})
		return
	}
	if folder.IsDeleted {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "文件夹已删除",
		})
		return
	}

	folder.IsDeleted = true
	folder.DeleterId = request.DeleterId
	folder.DeleteTime = time.Now()
	result := entity.Db.Model(&folder).Where("folder_id = ?", request.FolderId).Updates(&folder)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":   "文件夹删除失败",
			"error": result.Error.Error(),
		})
		return
	}

	var documents []entity.Document
	entity.Db.Where("parent_id = ?", request.FolderId).Find(&documents)
	for _, document := range documents {
		document.IsDeleted = true
		document.DeleterId = request.DeleterId
		document.DeleteTime = time.Now()
		result = entity.Db.Model(&document).Where("doc_id = ?", document.DocId).Updates(&document)
		if result.Error != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"msg":   "文件夹中文件删除失败",
				"error": result.Error.Error(),
			})
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg": "文件夹删除成功",
	})
}

func FolderCompleteDelete(ctx *gin.Context) {
	var request FolderDeleteRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var folder entity.Folder
	entity.Db.Find(&folder, "folder_id = ?", request.FolderId)
	if folder.FolderId == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "文件夹不存在",
		})
		return
	}
	if strings.HasSuffix(folder.Name, "的项目文档") {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "该文件夹为项目文档文件夹，不能被删除",
		})
		return
	}

	var documents []entity.Document
	entity.Db.Where("parent_id = ?", request.FolderId).Find(&documents)
	for _, document := range documents {
		document.IsDeleted = true
		document.DeleterId = request.DeleterId
		document.DeleteTime = time.Now()
		result := entity.Db.Model(&document).Where("doc_id = ?", document.DocId).Updates(&document)
		if result.Error != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"msg":   "文件夹中文件删除失败",
				"error": result.Error.Error(),
			})
			return
		}
		result = entity.Db.Model(&document).Where("doc_id = ?", document.DocId).Update("parent_id", 0)
	}

	result := entity.Db.Where("folder_id = ?", folder.FolderId).Delete(&folder)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":   "文件夹删除失败",
			"error": result.Error.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "文件夹删除成功",
	})
}

func FolderRename(ctx *gin.Context) {
	var request FolderRenameRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var folder entity.Folder
	entity.Db.Find(&folder, "folder_id = ?", request.FolderId)
	if folder.FolderId == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "文件夹不存在",
		})
		return
	}
	if folder.IsDeleted {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "文件夹已删除",
		})
		return
	}

	if strings.HasSuffix(folder.Name, "的项目文档") {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "该文件夹为项目文档文件夹，名称不能修改",
		})
		return
	}
	if strings.HasSuffix(request.NewName, "的项目文档") {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "新名称不能为项目文档文件夹",
		})
		return
	}

	result := entity.Db.Model(&folder).Where("folder_id = ?", request.FolderId).Update("name", request.NewName)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":   "文件夹重命名失败",
			"error": result.Error.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "文件夹重命名成功",
	})
}

func FolderList(ctx *gin.Context) {
	teamId, ok := ctx.GetQuery("teamId")
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "teamId不能为空",
		})
		return
	}

	var folders []entity.Folder
	var folderList []gin.H
	var projectFolderList []gin.H
	entity.Db.Where("team_id = ?", teamId).Find(&folders)
	sort.SliceStable(folders, func(i, j int) bool {
		return folders[i].CreateTime.Unix() > folders[j].CreateTime.Unix()
	})
	for _, folder := range folders {
		if folder.IsDeleted {
			continue
		}
		var creator entity.User
		entity.Db.Where("user_id = ?", folder.CreatorId).Find(&creator)
		if strings.HasSuffix(folder.Name, "的项目文档") {
			var project entity.Project
			entity.Db.Where("name = ?", strings.TrimSuffix(folder.Name, "的项目文档")).Find(&project)
			folderJson := gin.H{
				"folderId":   folder.FolderId,
				"name":       folder.Name,
				"creatorId":  folder.CreatorId,
				"createInfo": string(folder.CreateTime.Format("2006-01-02 15:04")) + " by " + creator.Nickname,
				"projectId":  project.ProjectId,
			}
			projectFolderList = append(projectFolderList, folderJson)
		} else {
			folderJson := gin.H{
				"folderId":   folder.FolderId,
				"name":       folder.Name,
				"creatorId":  folder.CreatorId,
				"createInfo": string(folder.CreateTime.Format("2006-01-02 15:04")) + " by " + creator.Nickname,
			}
			folderList = append(folderList, folderJson)
		}
	}
	if len(projectFolderList) == 0 {
		folderList = append(folderList, gin.H{
			"项目文档区": make([]gin.H, 0),
		})
	} else {
		folderList = append(folderList, gin.H{
			"name":              "项目文档区",
			"projectFolderList": projectFolderList,
		})
	}
	if len(folderList) == 1 && len(projectFolderList) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"msg":        "当前团队没有文件夹",
			"folderList": folderList,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":        "文件夹列表获取成功",
		"folderList": folderList,
	})
}

func FolderMove(ctx *gin.Context) {
	var request FolderMoveRequest
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

	var folder entity.Folder
	if request.FolderId != 0 {
		entity.Db.Find(&folder, "folder_id = ?", request.FolderId)
		if folder.FolderId == 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"msg": "文件夹不存在",
			})
			return
		}
		if folder.IsDeleted {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"msg": "文件夹已删除",
			})
			return
		}
	}

	result := entity.Db.Model(&document).Where("doc_id = ?", request.DocId).Update("parent_id", request.FolderId)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":   "文档移动失败",
			"error": result.Error.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "文档移动成功",
	})
}
