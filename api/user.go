package api

import (
	"backend/entity"
	"github.com/gin-gonic/gin"
)

func UserRegister(c *gin.Context) {
	var user entity.User
	err := c.ShouldBind(&user)
	if err != nil {
		c.JSON(200, gin.H{
			"code": -1,
		})
		return
	} else {
		var temp entity.User
		entity.Db.Find(&temp, "email=?", user.Email)
		if temp.ID != 0 {
			c.JSON(200, gin.H{
				"code": 1,
			})
			return
		}
		entity.Db.Create(&user)
		c.JSON(200, gin.H{
			"code": 0,
		})
		return
	}
}

func UserLogin(c *gin.Context) {
	username, ok := c.GetPostForm("username")
	if !ok {
		c.JSON(200, gin.H{
			"code": -1,
		})
		return
	}
	password, ok := c.GetPostForm("password")
	if !ok {
		c.JSON(200, gin.H{
			"code": -1,
		})
		return
	}
	var loginUser entity.User
	entity.Db.Find(&loginUser, "nickname=?", username)
	if loginUser.ID == 0 {
		c.JSON(200, gin.H{
			"code": 1,
		})
		return
	}
	if loginUser.Password != password {
		c.JSON(200, gin.H{
			"code": 2,
		})
		return
	}
	c.JSON(200, gin.H{
		"code":     0,
		"id":       loginUser.ID,
		"username": loginUser.Nickname,
	})
}
