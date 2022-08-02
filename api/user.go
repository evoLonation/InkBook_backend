package api

import (
	"backend/entity"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

func UserRegister(c *gin.Context) {
	var user entity.User
	err := c.ShouldBind(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
		})
		return
	} else {
		var temp entity.User
		entity.Db.Find(&temp, "email=?", user.Email)
		if temp.ID != 0 {
			c.JSON(http.StatusFound, gin.H{
				"code": 1,
			})
			return
		}
		entity.Db.Create(&user)
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
		})
		return
	}
}

func UserLogin(c *gin.Context) {
	username, ok := c.GetQuery("username")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
		})
		return
	}
	password, ok := c.GetPostForm("password")
	if !ok {
		c.JSON(200, gin.H{
			"code": ok,
		})
		return
	}
	var loginUser entity.User
	selectErr := entity.Db.Find(&loginUser, "nickname=?", username).Error
	errors.Is(selectErr, gorm.ErrRecordNotFound)
	if selectErr != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 1,
		})
		return
	}
	if loginUser.Password != password {
		c.JSON(http.StatusConflict, gin.H{
			"code": 2,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":     0,
		"id":       loginUser.ID,
		"username": loginUser.Nickname,
	})
}
