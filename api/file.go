package api

import (
	"backend/entity"
	"github.com/gin-gonic/gin"
	"net/http"
	"sort"
	"strings"
)

func FileList(ctx *gin.Context) {
	teamId, ok := ctx.GetQuery("teamId")
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "teamId不能为空"})
		return
	}
	parentId, ok := ctx.GetQuery("parentId")
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "parentId不能为空"})
		return
	}

	var team entity.Team
	entity.Db.Find(&team, "team_id = ?", teamId)
	if team.TeamId == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "团队不存在",
		})
	}

	var folders []entity.Folder
	var folderList []gin.H
	if parentId == "-1" {
		entity.Db.Find(&folders, "team_id = ? AND parent_id = ?", teamId, 0)
	} else {
		entity.Db.Find(&folders, "team_id = ? AND parent_id = ?", teamId, parentId)
	}
	sort.SliceStable(folders, func(i, j int) bool {
		return folders[i].CreateTime.Unix() > folders[j].CreateTime.Unix()
	})
	for _, folder := range folders {
		if (parentId == "-1" && !strings.HasSuffix(folder.Name, "的项目文档")) ||
			(parentId != "-1" && strings.HasSuffix(folder.Name, "的项目文档")) {
			continue
		}
		var creator entity.User
		entity.Db.Where("user_id = ?", folder.CreatorId).Find(&creator)
		folderJson := gin.H{
			"folderId":   folder.FolderId,
			"name":       folder.Name,
			"creatorId":  folder.CreatorId,
			"createInfo": string(folder.CreateTime.Format("2006-01-02 15:04")) + " by " + creator.Nickname,
		}
		folderList = append(folderList, folderJson)
	}
	if parentId == "0" {
		folderList = append(folderList, gin.H{
			"folderId":   -1,
			"name":       "项目文档区",
			"creatorId":  0,
			"createInfo": "",
		})
	}
	if len(folderList) == 0 {
		folderList = make([]gin.H, 0)
	}

	var documents []entity.Document
	var docList []gin.H
	entity.Db.Find(&documents, "parent_id = ? AND team_id = ?", parentId, teamId)
	sort.SliceStable(documents, func(i, j int) bool {
		return documents[i].CreateTime.Unix() > documents[j].CreateTime.Unix()
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
		docList = make([]gin.H, 0)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":        "获取文件列表成功",
		"folderList": folderList,
		"docList":    docList,
	})
}
