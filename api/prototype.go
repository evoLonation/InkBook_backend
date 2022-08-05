package api

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

//var protoEditorMap = make(map[int][]string)
//var protoEditTimeMap = make(map[int]time.Time)
//var protoUserTimeMap = make(map[string]time.Time)
//
//func PrototypeCreate(ctx *gin.Context) {
//	var request PrototypeCreateRequest
//	if err := ctx.ShouldBindJSON(&request); err != nil {
//		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//
//	var prototype entity.Prototype
//	entity.Db.Find(&prototype, "name = ?", request.Name)
//	if prototype.Name != "" {
//		ctx.JSON(http.StatusBadRequest, gin.H{
//			"msg": "原型名称已存在",
//		})
//		return
//	}
//
//	prototype = entity.Prototype{
//		Name:       request.Name,
//		ProjectID:  request.ProjectID,
//		CreatorID:  request.CreatorID,
//		CreateTime: time.Now(),
//		ModifierID: request.CreatorID,
//		ModifyTime: time.Now(),
//		IsEditing:  false,
//		IsDeleted:  false,
//	}
//}
//
//func PrototypeDelete(ctx *gin.Context) {
//	var request DocumentDeleteRequest
//	if err := ctx.ShouldBindJSON(&request); err != nil {
//		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//
//	var document entity.Document
//	entity.Db.Find(&document, "doc_id = ?", request.DocID)
//	if document.DocID == 0 {
//		ctx.JSON(http.StatusBadRequest, gin.H{
//			"msg": "文档不存在",
//		})
//		return
//	}
//	if document.IsDeleted {
//		ctx.JSON(http.StatusBadRequest, gin.H{
//			"msg": "文档已删除",
//		})
//		return
//	}
//
//	document.IsDeleted = true
//	document.DeleterID = request.DeleterID
//	document.DeleteTime = time.Now()
//	result := entity.Db.Model(&document).Where("doc_id = ?", request.DocID).Updates(&document)
//	if result.Error != nil {
//		ctx.JSON(http.StatusBadRequest, gin.H{
//			"msg":   "文档删除失败",
//			"error": result.Error.Error(),
//		})
//		return
//	}
//	ctx.JSON(http.StatusOK, gin.H{
//		"msg": "文档删除成功",
//	})
//}
//
//func PrototypeCompleteDelete(ctx *gin.Context) {
//	var request DocumentCompleteDeleteRequest
//	if err := ctx.ShouldBindJSON(&request); err != nil {
//		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//
//	var document entity.Document
//	entity.Db.Find(&document, "doc_id = ?", request.DocID)
//	if document.DocID == 0 {
//		ctx.JSON(http.StatusBadRequest, gin.H{
//			"msg": "文档不存在",
//		})
//		return
//	}
//
//	result := entity.Db.Where("doc_id = ?", request.DocID).Delete(&document)
//	if result.Error != nil {
//		ctx.JSON(http.StatusBadRequest, gin.H{
//			"msg":   "文档删除失败",
//			"error": result.Error.Error(),
//		})
//		return
//	}
//	ctx.JSON(http.StatusOK, gin.H{
//		"msg": "文档删除成功",
//	})
//}
//
//func PrototypeRename(ctx *gin.Context) {
//	var request DocumentRenameRequest
//	if err := ctx.ShouldBindJSON(&request); err != nil {
//		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//
//	var document entity.Document
//	entity.Db.Find(&document, "doc_id = ?", request.DocID)
//	if document.DocID == 0 {
//		ctx.JSON(http.StatusBadRequest, gin.H{
//			"msg": "文档不存在",
//		})
//		return
//	}
//
//	var documents []entity.Document
//	entity.Db.Where("name = ? and project_id = ?", request.NewName, document.ProjectID).Find(&documents)
//	if len(documents) != 0 {
//		ctx.JSON(http.StatusBadRequest, gin.H{
//			"msg": "文档名称重复",
//		})
//		return
//	}
//
//	result := entity.Db.Model(&document).Where("doc_id = ?", request.DocID).Update("name", request.NewName)
//	if result.Error != nil {
//		ctx.JSON(http.StatusBadRequest, gin.H{
//			"msg":   "文档重命名失败",
//			"error": result.Error.Error(),
//		})
//		return
//	}
//	ctx.JSON(http.StatusOK, gin.H{
//		"msg": "文档重命名成功",
//	})
//}
//
//func PrototypeList(ctx *gin.Context) {
//	projectId, ok := ctx.GetQuery("projectId")
//	if !ok {
//		ctx.JSON(http.StatusBadRequest, gin.H{
//			"msg": "projectId不能为空",
//		})
//		return
//	}
//
//	var documents []entity.Document
//	var docList []gin.H
//	entity.Db.Where("project_id = ?", projectId).Find(&documents)
//	sort.SliceStable(documents, func(i, j int) bool {
//		return documents[i].Name < documents[j].Name
//	})
//	for _, document := range documents {
//		if document.IsDeleted {
//			continue
//		}
//		var creator, modifier entity.User
//		entity.Db.Where("user_id = ?", document.CreatorID).Find(&creator)
//		entity.Db.Where("user_id = ?", document.ModifierID).Find(&modifier)
//		documentJson := gin.H{
//			"docId":      document.DocID,
//			"docName":    document.Name,
//			"creatorId":  document.CreatorID,
//			"createInfo": string(document.CreateTime.Format("2006-01-02 15:04")) + " by " + creator.Nickname,
//			"modifierId": document.ModifierID,
//			"modifyInfo": string(document.ModifyTime.Format("2006-01-02 15:04")) + " by " + modifier.Nickname,
//		}
//		docList = append(docList, documentJson)
//	}
//	if len(docList) == 0 {
//		ctx.JSON(http.StatusOK, gin.H{
//			"msg":     "当前项目没有文档",
//			"docList": make([]entity.Document, 0),
//		})
//		return
//	}
//	ctx.JSON(http.StatusOK, gin.H{
//		"msg":     "文档列表获取成功",
//		"docList": docList,
//	})
//}
//
//func PrototypeRecycle(ctx *gin.Context) {
//	projectId, ok := ctx.GetQuery("projectId")
//	if !ok {
//		ctx.JSON(http.StatusBadRequest, gin.H{
//			"msg": "projectId不能为空",
//		})
//		return
//	}
//
//	var documents []entity.Document
//	var docList []gin.H
//	entity.Db.Where("project_id = ?", projectId).Find(&documents)
//	sort.SliceStable(documents, func(i, j int) bool {
//		return documents[i].Name < documents[j].Name
//	})
//	for _, document := range documents {
//		if !document.IsDeleted {
//			continue
//		}
//		var creator, modifier entity.User
//		entity.Db.Where("user_id = ?", document.CreatorID).Find(&creator)
//		entity.Db.Where("user_id = ?", document.ModifierID).Find(&modifier)
//		documentJson := gin.H{
//			"docId":      document.DocID,
//			"docName":    document.Name,
//			"creatorId":  document.CreatorID,
//			"createInfo": string(document.CreateTime.Format("2006-01-02 15:04")) + " by " + creator.Nickname,
//			"modifierId": document.ModifierID,
//			"modifyInfo": string(document.ModifyTime.Format("2006-01-02 15:04")) + " by " + modifier.Nickname,
//		}
//		docList = append(docList, documentJson)
//	}
//	if len(docList) == 0 {
//		ctx.JSON(http.StatusOK, gin.H{
//			"msg":     "当前回收站没有文档",
//			"docList": make([]entity.Document, 0),
//		})
//		return
//	}
//	ctx.JSON(http.StatusOK, gin.H{
//		"msg":     "回收站文档列表获取成功",
//		"docList": docList,
//	})
//}
//
//func PrototypeRecover(ctx *gin.Context) {
//	var request DocumentDeleteRequest
//	if err := ctx.ShouldBindJSON(&request); err != nil {
//		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//
//	var document entity.Document
//	entity.Db.Find(&document, "doc_id = ?", request.DocID)
//	if document.DocID == 0 {
//		ctx.JSON(http.StatusBadRequest, gin.H{
//			"msg": "文档不存在",
//		})
//		return
//	}
//	if !document.IsDeleted {
//		ctx.JSON(http.StatusBadRequest, gin.H{
//			"msg": "文档不在回收站中",
//		})
//		return
//	}
//
//	result := entity.Db.Model(&document).Where("doc_id = ?", request.DocID).Update("is_deleted", false)
//	if result.Error != nil {
//		ctx.JSON(http.StatusBadRequest, gin.H{
//			"msg":   "文档恢复失败",
//			"error": result.Error.Error(),
//		})
//		return
//	}
//	ctx.JSON(http.StatusOK, gin.H{
//		"msg": "文档恢复成功",
//	})
//}
//
//func PrototypeSave(ctx *gin.Context) {
//	var request DocumentSaveRequest
//	if err := ctx.ShouldBindJSON(&request); err != nil {
//		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//
//	var document entity.Document
//	entity.Db.Find(&document, "doc_id = ?", request.DocID)
//	if document.DocID == 0 {
//		ctx.JSON(http.StatusBadRequest, gin.H{
//			"msg": "文档不存在",
//		})
//		return
//	}
//
//	//jsonContent, err := json.Marshal(request.Content)
//	//if err != nil {
//	//	ctx.JSON(http.StatusBadRequest, gin.H{
//	//		"error": err.Error(),
//	//		"msg":   "JSON格式内容解析失败",
//	//	})
//	//	return
//	//}
//	document.Content = "{\"content\":\"" + string(request.Content) + "\"}"
//	fmt.Println(document.Content)
//	document.ModifierID = request.UserId
//	document.ModifyTime = time.Now()
//	result := entity.Db.Model(&document).Where("doc_id = ?", request.DocID).Updates(&document)
//	if result.Error != nil {
//		ctx.JSON(http.StatusBadRequest, gin.H{
//			"msg": "文档保存失败",
//		})
//		return
//	}
//	ctx.JSON(http.StatusOK, gin.H{
//		"msg": "文档保存成功",
//	})
//}
//
//func PrototypeExit(ctx *gin.Context) {
//	var request DocumentExitRequest
//	if err := ctx.ShouldBindJSON(&request); err != nil {
//		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//
//	var document entity.Document
//	entity.Db.Find(&document, "doc_id = ?", request.DocID)
//	if document.DocID == 0 {
//		ctx.JSON(http.StatusBadRequest, gin.H{
//			"msg": "文档不存在",
//		})
//		return
//	}
//	if document.IsDeleted {
//		ctx.JSON(http.StatusBadRequest, gin.H{
//			"msg": "文档在回收站中",
//		})
//		return
//	}
//
//	if time.Now().Sub(docEditTimeMap[request.DocID]) > time.Second*3 {
//		docEditorMap[request.DocID] = []string{}
//		ctx.JSON(http.StatusBadRequest, gin.H{
//			"msg": "文档不在编辑状态",
//		})
//		return
//	}
//
//	editors := docEditorMap[request.DocID]
//	var nowEditors []string
//	for _, editor := range editors {
//		if editor == request.UserId || time.Now().Sub(docUserTimeMap[editor]) > time.Second*3 {
//			continue
//		}
//		nowEditors = append(nowEditors, editor)
//	}
//
//	document.ModifierID = request.UserId
//	document.ModifyTime = time.Now()
//	result := entity.Db.Model(&document).Where("doc_id = ?", request.DocID).Updates(&document)
//	if result.Error != nil {
//		ctx.JSON(http.StatusBadRequest, gin.H{
//			"msg": "文档退出编辑失败",
//		})
//		return
//	}
//	ctx.JSON(http.StatusOK, gin.H{
//		"msg":    "文档退出编辑成功",
//		"remain": len(nowEditors),
//	})
//}
//
//func PrototypeGet(ctx *gin.Context) {
//	docId, ok := ctx.GetQuery("docId")
//	if !ok {
//		ctx.JSON(http.StatusBadRequest, gin.H{
//			"msg": "docId不能为空",
//		})
//		return
//	}
//
//	var document entity.Document
//	entity.Db.Find(&document, "doc_id = ?", docId)
//	if document.DocID == 0 {
//		ctx.JSON(http.StatusBadRequest, gin.H{
//			"msg": "文档不存在",
//		})
//		return
//	}
//	if document.IsDeleted {
//		ctx.JSON(http.StatusBadRequest, gin.H{
//			"msg": "文档在回收站中，无法编辑",
//		})
//		return
//	}
//
//	var jsonContent gin.H
//	if err := json.Unmarshal([]byte(document.Content), &jsonContent); err != nil {
//		ctx.JSON(http.StatusBadRequest, gin.H{
//			"msg": "JSON格式内容解析失败",
//		})
//		return
//	}
//	ctx.JSON(http.StatusOK, gin.H{
//		"msg":     "文档获取成功",
//		"content": jsonContent["content"],
//	})
//}
//
//func PrototypeApplyEdit(ctx *gin.Context) {
//	var request DocumentApplyEditRequest
//	if err := ctx.ShouldBindJSON(&request); err != nil {
//		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//
//	var document entity.Document
//	entity.Db.Find(&document, "doc_id = ?", request.DocID)
//	if document.DocID == 0 {
//		ctx.JSON(http.StatusBadRequest, gin.H{
//			"msg": "文档不存在",
//		})
//		return
//	}
//	if document.IsDeleted {
//		ctx.JSON(http.StatusBadRequest, gin.H{
//			"msg": "文档已被删除",
//		})
//		return
//	}
//
//	editors := docEditorMap[request.DocID]
//	docUserTimeMap[request.UserId] = time.Now()
//	var nowEditors []string
//	for _, editor := range editors {
//		if editor == request.UserId || time.Now().Sub(docUserTimeMap[editor]) > time.Second*3 {
//			continue
//		}
//		nowEditors = append(nowEditors, editor)
//	}
//
//	nowEditors = append(nowEditors, request.UserId)
//	docEditorMap[request.DocID] = nowEditors
//	docEditTimeMap[request.DocID] = time.Now()
//	ctx.JSON(http.StatusOK, gin.H{
//		"msg":          "申请编辑状态成功",
//		"nowEditorNum": len(nowEditors),
//	})
//}
