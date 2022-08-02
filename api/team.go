package api

import (
	"backend/entity"
	"github.com/gin-gonic/gin"
	"net/http"
)

func TeamRegister(c *gin.Context) {
	var team entity.Team
	err := c.ShouldBind(&team)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "参数错误",
		})
		return
	} else {
		entity.Db.Create(&team)
		//entity.Db
		c.JSON(200, gin.H{
			"code": 0,
		})
		return
	}
}
