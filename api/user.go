package api

import (
	"backend/entity"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"math/rand"
	"net/http"
	"time"
)

func UserRegister(c *gin.Context) {
	var user entity.User
	err := c.ShouldBind(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "参数错误",
		})
		return
	} else {
		var temp entity.User
		//同一邮箱只允许被注册一次
		entity.Db.Find(&temp, "email=?", user.Email)
		if temp.UserId != "" {
			c.JSON(http.StatusFound, gin.H{
				"msg":  "该邮箱已被注册",
				"user": user.Email,
			})
			return
		}
		userCode, ok := c.GetPostForm("userCode")
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "参数错误",
			})
			return
		}
		sendCode, ok := c.GetPostForm("sendCode")
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "参数错误",
			})
			return
		}
		if sendCode != userCode {
			c.JSON(http.StatusConflict, gin.H{
				"msg": "验证码错误",
			})
			return
		}
		entity.Db.Create(&user)
		c.JSON(http.StatusOK, gin.H{
			"msg": "注册成功",
		})
		return
	}
}

func Identifying(c *gin.Context) {
	email := []string{c.PostForm("email")}
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	code := fmt.Sprintf("%06v", rnd.Int31n(1000000))
	err := SendEmail(email, code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "参数错误",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "发送成功",
	})
}

func UserLogin(c *gin.Context) {
	var flag = 0
	username, ok := c.GetQuery("username")
	var email string
	if !ok {
		email, ok = c.GetQuery("email")
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "参数错误",
			})
			return
		}
		flag = 1
	}
	password, ok := c.GetQuery("password")
	if !ok {
		c.JSON(200, gin.H{
			"msg": "参数错误",
		})
		return
	}
	var loginUser entity.User
	var selectErr error
	if flag == 0 {
		selectErr = entity.Db.Find(&loginUser, "nickname=?", username).Error
	} else {
		selectErr = entity.Db.Find(&loginUser, "email=?", email).Error
	}
	errors.Is(selectErr, gorm.ErrRecordNotFound)
	if selectErr != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "用户不存在",
		})
		return
	}
	if loginUser.Password != password {
		c.JSON(http.StatusConflict, gin.H{
			"msg": "密码错误",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":      "登陆成功",
		"id":       loginUser.UserId,
		"username": loginUser.Nickname,
	})
}
func UserInformation(c *gin.Context) {
	userId, ok := c.GetQuery("userId")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "参数错误",
		})
		return
	}
	var user entity.User
	selectErr := entity.Db.Find(&user, "userId=?", userId).Error
	errors.Is(selectErr, gorm.ErrRecordNotFound)
	if selectErr != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "用户不存在",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":          "返回成功",
		"nickname":     user.Nickname,
		"email":        user.Email,
		"introduction": user.Intro,
	})
}
func UserModifyPassword(c *gin.Context) {
	userId, ok := c.GetPostForm("userId")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "参数错误",
		})
		return
	}
	newPwd, ok := c.GetPostForm("newPwd")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "参数错误",
		})
		return
	}
	oldpwd, ok := c.GetPostForm("OldPwd")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "参数错误",
		})
		return
	}
	var user entity.User
	selectErr := entity.Db.Find(&user, "userId=?", userId).Error
	errors.Is(selectErr, gorm.ErrRecordNotFound)
	if selectErr != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "用户不存在",
		})
		return
	}
	if user.Password != oldpwd {
		c.JSON(http.StatusConflict, gin.H{
			"msg": "密码错误",
		})
		return
	}
	entity.Db.Update("password", newPwd)
	c.JSON(http.StatusOK, gin.H{
		"msg": "更新成功",
	})
}
func UserModifyEmail(c *gin.Context) {
	userId, ok := c.GetPostForm("userId")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "参数错误",
		})
		return
	}
	newEmail, ok := c.GetPostForm("newEmail")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "参数错误",
		})
		return
	}
	userCode, ok := c.GetPostForm("userCode")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "参数错误",
		})
		return
	}
	sendCode, ok := c.GetPostForm("sendCode")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "参数错误",
		})
		return
	}
	if sendCode != userCode {
		c.JSON(http.StatusConflict, gin.H{
			"msg": "验证码错误",
		})
		return
	}
	var user entity.User
	selectErr := entity.Db.Find(&user, "userId=?", userId).Error
	errors.Is(selectErr, gorm.ErrRecordNotFound)
	if selectErr != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "用户不存在",
		})
		return
	}
	entity.Db.Update("email", newEmail)
	c.JSON(http.StatusOK, gin.H{
		"msg": "更新成功",
	})
}
func UserModifyIntroduction(c *gin.Context) {
	userId, ok := c.GetPostForm("userId")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "参数错误",
		})
		return
	}
	newIntro, ok := c.GetPostForm("newIntro")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "参数错误",
		})
		return
	}
	var user entity.User
	selectErr := entity.Db.Find(&user, "userId=?", userId).Error
	errors.Is(selectErr, gorm.ErrRecordNotFound)
	if selectErr != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "用户不存在",
		})
		return
	}
	entity.Db.Update("intro", newIntro)
	c.JSON(http.StatusOK, gin.H{
		"msg": "更新成功",
	})
}
func UserModifyNickname(c *gin.Context) {
	userId, ok := c.GetPostForm("userId")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "参数错误",
		})
		return
	}
	newNick, ok := c.GetPostForm("newNick")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "参数错误",
		})
		return
	}
	var user entity.User
	selectErr := entity.Db.Find(&user, "userId=?", userId).Error
	errors.Is(selectErr, gorm.ErrRecordNotFound)
	if selectErr != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "用户不存在",
		})
		return
	}
	entity.Db.Update("nickname", newNick)
	c.JSON(http.StatusOK, gin.H{
		"msg": "更新成功",
	})
}
func UserModifyRealname(c *gin.Context) {
	userId, ok := c.GetPostForm("userId")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "参数错误",
		})
		return
	}
	newReal, ok := c.GetPostForm("newReal")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "参数错误",
		})
		return
	}
	var user entity.User
	selectErr := entity.Db.Find(&user, "userId=?", userId).Error
	errors.Is(selectErr, gorm.ErrRecordNotFound)
	if selectErr != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "用户不存在",
		})
		return
	}
	entity.Db.Update("newReal", newReal)
	c.JSON(http.StatusOK, gin.H{
		"msg": "更新成功",
	})
}
func UserTeam(c *gin.Context) {
	userId, ok := c.GetQuery("userId")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "参数错误",
		})
		return
	}
	var teams []entity.Team
	selectErr := entity.Db.Table("teams").Select("teams.team_id as team_id,teams.name as name,teams.intro as intro").Joins("left join user_team on user_team.team_id = teams.team_id where user_team.user_id <> ? ", userId).Scan(&teams).Error
	errors.Is(selectErr, gorm.ErrRecordNotFound)
	if selectErr != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "用户不存在",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":   "查找成功",
		"teams": teams,
	})
}
