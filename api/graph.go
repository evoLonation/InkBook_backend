package api

import (
	"backend/entity"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"sort"
	"time"
)

type GraphCreateRequest struct {
	Name      string `json:"name"`
	CreatorID string `json:"creatorId"`
	ProjectID int    `json:"projectId"`
}

type GraphDeleteRequest struct {
	GraphID   int    `json:"graphId"`
	DeleterID string `json:"deleterId"`
}

type GraphCompleteDeleteRequest struct {
	GraphID int `json:"graphId"`
}

type GraphRenameRequest struct {
	GraphID int    `json:"graphId"`
	NewName string `json:"newName"`
}

type GraphListRequest struct {
	ProjectID int `json:"projectId"`
}

type GraphSaveRequest struct {
	GraphID int    `json:"graphId"`
	UserId  string `json:"userId"`
	Content string `json:"content"`
}

type GraphExitRequest struct {
	GraphID int    `json:"graphId"`
	UserId  string `json:"userId"`
}

type GraphApplyEditRequest struct {
	GraphID int    `json:"graphId"`
	UserId  string `json:"userId"`
}

var graphEditorMap = make(map[int][]string)
var graphEditTimeMap = make(map[int]time.Time)
var graphUserTimeMap = make(map[string]time.Time)

func GraphCreate(ctx *gin.Context) {
	var request GraphCreateRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var graph entity.Graph
	entity.Db.Find(&graph, "name = ?", request.Name)
	if graph != (entity.Graph{}) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "UML图已存在",
		})
		return
	}

	graph = entity.Graph{
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
		EditingCnt: 0,
	}
	result := entity.Db.Create(&graph)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":   "UML图创建失败",
			"error": result.Error.Error(),
		})
		return
	}
	entity.Db.Where("name = ? and project_id = ?", request.Name, request.ProjectID).First(&graph)
	ctx.JSON(http.StatusOK, gin.H{
		"msg":     "UML图创建成功",
		"graphId": graph.GraphID,
	})
}

func GraphDelete(ctx *gin.Context) {
	var request GraphDeleteRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var graph entity.Graph
	entity.Db.Find(&graph, "graph_id = ?", request.GraphID)
	if graph == (entity.Graph{}) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "UML图不存在",
		})
		return
	}
	if graph.IsDeleted {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "UML图已删除",
		})
		return
	}

	graph.IsDeleted = true
	graph.DeleterID = request.DeleterID
	graph.DeleteTime = time.Now()
	result := entity.Db.Model(&graph).Where("graph_id = ?", request.GraphID).Updates(&graph)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":   "UML图删除失败",
			"error": result.Error.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "UML图删除成功",
	})
}

func GraphCompleteDelete(ctx *gin.Context) {
	var request GraphCompleteDeleteRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var graph entity.Graph
	entity.Db.Find(&graph, "graph_id = ?", request.GraphID)
	if graph == (entity.Graph{}) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "UML图不存在",
		})
		return
	}

	result := entity.Db.Where("graph_id = ?", request.GraphID).Delete(&graph)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":   "UML图删除失败",
			"error": result.Error.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "UML图删除成功",
	})
}

func GraphRename(ctx *gin.Context) {
	var request GraphRenameRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var graph entity.Graph
	entity.Db.Find(&graph, "graph_id = ?", request.GraphID)
	if graph == (entity.Graph{}) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "UML图不存在",
		})
		return
	}

	var graphs []entity.Graph
	entity.Db.Where("name = ? and project_id = ?", request.NewName, graph.ProjectID).Find(&graphs)
	if len(graphs) != 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "UML图名称重复",
		})
		return
	}

	result := entity.Db.Model(&graph).Where("graph_id = ?", request.GraphID).Update("name", request.NewName)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":   "UML图重命名失败",
			"error": result.Error.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "UML图重命名成功",
	})
}

func GraphList(ctx *gin.Context) {
	projectId, ok := ctx.GetQuery("projectId")
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "projectId不能为空",
		})
		return
	}

	var graphs []entity.Graph
	var graphList []gin.H
	entity.Db.Where("project_id = ?", projectId).Find(&graphs)
	sort.SliceStable(graphs, func(i, j int) bool {
		return graphs[i].Name < graphs[j].Name
	})
	for _, graph := range graphs {
		if graph.IsDeleted {
			continue
		}
		var creator entity.User
		entity.Db.Where("user_id = ?", graph.CreatorID).Find(&creator)
		graphJson := gin.H{
			"graphId":    graph.GraphID,
			"name":       graph.Name,
			"creatorId":  graph.CreatorID,
			"CreateInfo": string(graph.CreateTime.Format("2006-01-02 15:04:05")) + " by " + creator.Nickname,
			"ModifierID": graph.ModifierID,
			"ModifyInfo": string(graph.ModifyTime.Format("2006-01-02 15:04:05")) + " by " + creator.Nickname,
		}
		graphList = append(graphList, graphJson)
	}
	if len(graphList) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"msg":       "当前项目没有UML图",
			"graphList": make([]entity.Graph, 0),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":       "UML图列表获取成功",
		"graphList": graphList,
	})
}

func GraphRecycle(ctx *gin.Context) {
	projectId, ok := ctx.GetQuery("projectId")
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "projectId不能为空",
		})
		return
	}

	var graphs []entity.Graph
	var graphList []gin.H
	entity.Db.Where("project_id = ?", projectId).Find(&graphs)
	sort.SliceStable(graphs, func(i, j int) bool {
		return graphs[i].Name < graphs[j].Name
	})
	for _, graph := range graphs {
		if !graph.IsDeleted {
			continue
		}
		var creator, modifier entity.User
		entity.Db.Where("user_id = ?", graph.CreatorID).Find(&creator)
		entity.Db.Where("user_id = ?", graph.ModifierID).Find(&modifier)
		graphJson := gin.H{
			"graphId":    graph.GraphID,
			"name":       graph.Name,
			"creatorId":  graph.CreatorID,
			"CreateInfo": string(graph.CreateTime.Format("2006-01-02 15:04:05")) + " by " + creator.Nickname,
			"ModifierID": graph.ModifierID,
			"ModifyInfo": string(graph.ModifyTime.Format("2006-01-02 15:04:05")) + " by " + modifier.Nickname,
		}
		graphList = append(graphList, graphJson)
	}
	if len(graphList) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"msg":       "当前回收站没有UML图",
			"graphList": make([]entity.Graph, 0),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":       "回收站UML图列表获取成功",
		"graphList": graphList,
	})
}

func GraphRecover(ctx *gin.Context) {
	var request GraphDeleteRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var graph entity.Graph
	entity.Db.Find(&graph, "graph_id = ?", request.GraphID)
	if graph == (entity.Graph{}) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "UML图不存在",
		})
		return
	}
	if !graph.IsDeleted {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "UML图不在回收站中",
		})
		return
	}

	result := entity.Db.Model(&graph).Where("graph_id = ?", request.GraphID).Update("is_deleted", false)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":   "UML图恢复失败",
			"error": result.Error.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "UML图恢复成功",
	})
}

func GraphSave(ctx *gin.Context) {
	var request GraphSaveRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var graph entity.Graph
	entity.Db.Find(&graph, "graph_id = ?", request.GraphID)
	if graph == (entity.Graph{}) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "UML图不存在",
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
	graph.Content = "{\"content\": " + request.Content + "}"
	graph.ModifierID = request.UserId
	graph.ModifyTime = time.Now()
	result := entity.Db.Where("graph_id = ?", request.GraphID).Updates(&graph)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "UML图保存失败",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "UML图保存成功",
	})
}

func GraphExit(ctx *gin.Context) {
	var request GraphExitRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var graph entity.Graph
	entity.Db.Find(&graph, "graph_id = ?", request.GraphID)
	if graph == (entity.Graph{}) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "UML图不存在",
		})
		return
	}
	if graph.IsDeleted {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "UML图在回收站中",
		})
		return
	}

	if time.Now().Sub(graphEditTimeMap[request.GraphID]) > time.Second*3 {
		graphEditTimeMap[request.GraphID] = time.Now()
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "UML图不在编辑状态",
		})
	}

	editors := graphEditorMap[request.GraphID]
	var nowEditors []string
	for _, editor := range editors {
		if editor == request.UserId || time.Now().Sub(graphUserTimeMap[editor]) > time.Second*3 {
			continue
		}
		nowEditors = append(nowEditors, editor)
	}

	graph.ModifierID = request.UserId
	graph.ModifyTime = time.Now()
	result := entity.Db.Model(&graph).Where("graph_id = ?", request.GraphID).Updates(&graph)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "UML图退出编辑失败",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":    "UML图退出编辑成功",
		"remain": len(nowEditors),
	})
}

func GraphGet(ctx *gin.Context) {
	graphId, ok := ctx.GetQuery("graphId")
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "graphId不能为空",
		})
		return
	}

	var graph entity.Graph
	entity.Db.Find(&graph, "graph_id = ?", graphId)
	if graph == (entity.Graph{}) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "UML图不存在",
		})
		return
	}
	if graph.IsDeleted {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "UML图在回收站中, 无法编辑",
		})
		return
	}

	var jsonContent gin.H
	if err := json.Unmarshal([]byte(graph.Content), &jsonContent); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":   "JSON格式内容解析失败",
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":   "UML图获取成功",
		"graph": jsonContent["content"],
	})
}

func GraphApplyEdit(ctx *gin.Context) {
	var request GraphApplyEditRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var graph entity.Graph
	entity.Db.Find(&graph, "graph_id = ?", request.GraphID)
	if graph == (entity.Graph{}) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "UML图不存在",
		})
		return
	}
	if graph.IsDeleted {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "UML图已被删除",
		})
		return
	}

	editors := graphEditorMap[request.GraphID]
	graphUserTimeMap[request.UserId] = time.Now()
	var nowEditors []string
	for _, editor := range editors {
		if editor == request.UserId || time.Now().Sub(graphUserTimeMap[editor]) > time.Second*3 {
			continue
		}
		nowEditors = append(nowEditors, editor)
	}

	nowEditors = append(nowEditors, request.UserId)
	graphEditorMap[request.GraphID] = nowEditors
	graphEditTimeMap[request.GraphID] = time.Now()
	ctx.JSON(http.StatusOK, gin.H{
		"msg":          "UML图申请编辑成功",
		"nowEditorNum": len(nowEditors),
	})
}
