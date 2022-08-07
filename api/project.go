package api

import (
	"backend/entity"
	"github.com/gin-gonic/gin"
	"net/http"
	"sort"
	"time"
)

type ProjectCreateRequest struct {
	Name   string `json:"name"`
	TeamId int    `json:"teamId"`
	UserId string `json:"userId"`
	Detail string `json:"detail"`
	ImgURL string `json:"imgUrl"`
}

type ProjectDeleteRequest struct {
	ProjectId int `json:"projectId"`
}

type ProjectRenameRequest struct {
	ProjectId int    `json:"projectId"`
	NewName   string `json:"newName"`
}

type ProjectModifyIntroRequest struct {
	ProjectId int    `json:"projectId"`
	NewIntro  string `json:"newIntro"`
}

type ProjectModifyImgRequest struct {
	ProjectId int    `json:"projectId"`
	NewImgURL string `json:"newImgUrl"`
}

func ProjectCreate(ctx *gin.Context) {
	var request ProjectCreateRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var project entity.Project
	entity.Db.Find(&project, "name = ?", request.Name)
	if project != (entity.Project{}) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "项目已存在",
		})
		return
	}

	project = entity.Project{
		Name:       request.Name,
		TeamId:     request.TeamId,
		CreatorId:  request.UserId,
		CreateTime: time.Now(),
		IsDeleted:  false,
		DeleteTime: time.Now(),
		Intro:      request.Detail,
		ImgURL:     request.ImgURL,
	}
	if project.Intro == "" {
		project.Intro = "暂无项目简介"
	}
	result := entity.Db.Create(&project)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":   "项目创建失败",
			"error": result.Error.Error(),
		})
		return
	}

	folder := entity.Folder{
		Name:       request.Name,
		TeamId:     request.TeamId,
		CreatorId:  request.UserId,
		CreateTime: time.Now(),
		IsDeleted:  false,
		DeleterId:  request.UserId,
		DeleteTime: time.Now(),
	}
	result = entity.Db.Create(&folder)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":   "项目文件夹创建失败",
			"error": result.Error.Error(),
		})
		return
	}

	entity.Db.Where("name = ? AND team_id = ?", request.Name, request.TeamId).First(&project)
	entity.Db.Where("name = ? AND team_id = ?", request.Name, request.TeamId).First(&folder)
	ctx.JSON(http.StatusOK, gin.H{
		"msg":       "项目创建成功",
		"projectId": project.ProjectId,
		"folderId":  folder.FolderId,
	})
}

func ProjectDelete(ctx *gin.Context) {
	var request ProjectDeleteRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var project entity.Project
	entity.Db.Find(&project, "project_id = ?", request.ProjectId)
	if project == (entity.Project{}) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "项目不存在",
		})
		return
	}
	if project.IsDeleted {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "项目已删除",
		})
		return
	}

	project.IsDeleted = true
	project.DeleteTime = time.Now()
	result := entity.Db.Model(&project).Where("project_id = ?", request.ProjectId).Updates(&project)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": result.Error.Error(),
			"msg":   "项目删除失败",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "项目删除成功",
	})
}

func ProjectCompleteDelete(ctx *gin.Context) {
	var request ProjectDeleteRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var project entity.Project
	entity.Db.Find(&project, "project_id = ?", request.ProjectId)
	if project == (entity.Project{}) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "项目不存在",
		})
		return
	}

	result := entity.Db.Where("project_id = ?", request.ProjectId).Delete(&project)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": result.Error.Error(),
			"msg":   "项目删除失败",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "项目删除成功",
	})
}

func ProjectRename(ctx *gin.Context) {
	var request ProjectRenameRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var project entity.Project
	entity.Db.Find(&project, "project_id = ?", request.ProjectId)
	if project == (entity.Project{}) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "项目不存在",
		})
		return
	}

	var projects []entity.Project
	entity.Db.Where("name = ? AND team_id = ?", request.NewName, project.TeamId).Find(&projects)
	if len(projects) != 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "项目名称重复",
		})
		return
	}

	result := entity.Db.Model(&project).Where("project_id = ?", request.ProjectId).Update("name", request.NewName)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": result.Error.Error(),
			"msg":   "项目重命名失败",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "项目重命名成功",
	})
}

func ProjectModifyIntro(ctx *gin.Context) {
	var request ProjectModifyIntroRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var project entity.Project
	entity.Db.Find(&project, "project_id = ?", request.ProjectId)
	if project == (entity.Project{}) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "项目不存在",
		})
		return
	}

	if request.NewIntro == "" {
		request.NewIntro = "暂无项目简介"
	}
	entity.Db.Model(&project).Where("project_id = ?", request.ProjectId).Update("intro", request.NewIntro)
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "项目简介修改成功",
	})
}

func ProjectModifyImg(ctx *gin.Context) {
	var request ProjectModifyImgRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var project entity.Project
	entity.Db.Find(&project, "project_id = ?", request.ProjectId)
	if project == (entity.Project{}) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "项目不存在",
		})
		return
	}

	project.ImgURL = request.NewImgURL
	entity.Db.Model(&project).Where("project_id = ?", request.ProjectId).Update("img_url", request.NewImgURL)
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "项目图片修改成功",
	})
}

func ProjectListTeam(ctx *gin.Context) {
	teamId, ok := ctx.GetQuery("teamId")
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "teamId不能为空",
		})
		return
	}

	var projects []entity.Project
	var projectList []gin.H
	entity.Db.Where("team_id = ?", teamId).Find(&projects)
	sort.SliceStable(projects, func(i, j int) bool {
		return projects[i].CreateTime.Unix() > projects[j].CreateTime.Unix()
	})
	for _, project := range projects {
		if project.IsDeleted {
			continue
		}
		var creator entity.User
		entity.Db.Find(&creator, "user_id = ?", project.CreatorId)
		projectJson := gin.H{
			"id":         project.ProjectId,
			"name":       project.Name,
			"detail":     project.Intro,
			"imgUrl":     project.ImgURL,
			"creatorId":  project.CreatorId,
			"createInfo": string(project.CreateTime.Format("2006-01-02 15:04:05")) + " by " + creator.Nickname,
		}
		projectList = append(projectList, projectJson)
	}
	if len(projectList) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"msg":      "当前团队没有项目",
			"projects": make([]entity.Project, 0),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"projects": projectList,
	})
}

func ProjectListUser(ctx *gin.Context) {
	userId, ok := ctx.GetQuery("userId")
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "userId不能为空",
		})
		return
	}

	var teams []entity.TeamMember
	var projectList []gin.H
	entity.Db.Where("member_id = ?", userId).Find(&teams)
	sort.SliceStable(teams, func(i, j int) bool {
		return teams[i].TeamId < teams[j].TeamId
	})
	for _, team := range teams {
		var projects []entity.Project
		entity.Db.Where("team_id = ?", team.TeamId).Find(&projects)
		sort.SliceStable(projects, func(i, j int) bool {
			return projects[i].CreateTime.Unix() > projects[j].CreateTime.Unix()
		})
		for _, project := range projects {
			if project.IsDeleted {
				continue
			}
			var creator entity.User
			entity.Db.Find(&creator, "user_id = ?", project.CreatorId)
			projectJson := gin.H{
				"id":         project.ProjectId,
				"teamId":     team.TeamId,
				"name":       project.Name,
				"detail":     project.Intro,
				"imgUrl":     project.ImgURL,
				"creatorId":  project.CreatorId,
				"createInfo": string(project.CreateTime.Format("2006-01-02 15:04:05")) + " by " + creator.Nickname,
			}
			projectList = append(projectList, projectJson)
		}
	}
	if len(projectList) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"msg":      "当前用户没有项目",
			"projects": make([]entity.Project, 0),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"projects": projectList,
	})
}

func ProjectRecycle(ctx *gin.Context) {
	teamId, ok := ctx.GetQuery("teamId")
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "teamId不能为空",
		})
		return
	}

	var projects []entity.Project
	var projectList []gin.H
	entity.Db.Where("team_id = ?", teamId).Find(&projects)
	sort.SliceStable(projects, func(i, j int) bool {
		return projects[i].CreateTime.Unix() > projects[j].CreateTime.Unix()
	})
	for _, project := range projects {
		if !project.IsDeleted {
			continue
		}
		var creator entity.User
		entity.Db.Find(&creator, "user_id = ?", project.CreatorId)
		projectJson := gin.H{
			"id":         project.ProjectId,
			"name":       project.Name,
			"detail":     project.Intro,
			"imgUrl":     project.ImgURL,
			"creatorId":  project.CreatorId,
			"createInfo": string(project.CreateTime.Format("2006-01-02 15:04:05")) + " by " + creator.Nickname,
		}
		projectList = append(projectList, projectJson)
	}
	if len(projectList) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"msg":      "当前回收站中没有项目",
			"projects": make([]entity.Project, 0),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"projects": projectList,
	})
}

func ProjectRecover(ctx *gin.Context) {
	var request ProjectDeleteRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var project entity.Project
	entity.Db.Find(&project, "project_id = ?", request.ProjectId)
	if project == (entity.Project{}) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "项目不存在",
		})
		return
	}
	if !project.IsDeleted {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "项目不在回收站中",
		})
		return
	}

	entity.Db.Model(&project).Where("project_id = ?", request.ProjectId).Update("is_deleted", false)
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "项目恢复成功",
	})
}

func ProjectSearch(ctx *gin.Context) {
	keyword, ok := ctx.GetQuery("keyword")
	teamId, ok := ctx.GetQuery("teamId")
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "keyword不能为空",
		})
		return
	}

	var projects []entity.Project
	var projectList []gin.H
	entity.Db.Where("name LIKE ? AND team_id = ?", keyword, teamId).Find(&projects)
	sort.SliceStable(projects, func(i, j int) bool {
		return projects[i].CreateTime.Unix() > projects[j].CreateTime.Unix()
	})
	for _, project := range projects {
		if project.IsDeleted {
			continue
		}
		var creator entity.User
		entity.Db.Find(&creator, "user_id = ?", project.CreatorId)
		projectJson := gin.H{
			"id":         project.ProjectId,
			"name":       project.Name,
			"detail":     project.Intro,
			"imgUrl":     project.ImgURL,
			"creatorId":  project.CreatorId,
			"createInfo": string(project.CreateTime.Format("2006-01-02 15:04:05")) + " by " + creator.Nickname,
		}
		projectList = append(projectList, projectJson)
	}
	if len(projectList) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"msg":      "没有搜索到相关项目",
			"projects": make([]entity.Project, 0),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":      "项目搜索成功",
		"projects": projectList,
	})
}
