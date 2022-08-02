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
			"code": -1,
		})
		return
	} else {
		var temp entity.User
		entity.Db.Find(&temp, "name=?", team.Name)
		if temp.ID != 0 {
			c.JSON(200, gin.H{
				"code": 1,
			})
			return
		}
		entity.Db.Create(&team)
		c.JSON(200, gin.H{
			"code": 0,
		})
		return
	}
}
