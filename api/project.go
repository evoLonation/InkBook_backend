package api

import "github.com/gin-gonic/gin"

type ProjectCreateRequest struct {
	Name   string `json:"name"`
	TeamID int    `json:"teamId"`
	UserId int    `json:"userId"`
	Detail string `json:"detail"`
}

func ProjectCreate(ctx *gin.Context) {
	// TODO
}

func ProjectDelete(ctx *gin.Context) {
	// TODO
}

func ProjectRename(ctx *gin.Context) {
	// TODO
}

func ProjectList(ctx *gin.Context) {
	// TODO
}
