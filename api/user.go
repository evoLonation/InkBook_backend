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

type IdentifyRequest struct {
	email string
}
type UserRegisterRequest struct {
	userCode string
	sendCode string
}
type UserLoginRequest struct {
	userId string
	email  string
	pwd    string
}
type UserInfoRequest struct {
	userId   string
	oldPwd   string
	newPwd   string
	userCode string
	sendCode string
	newEmail string
	newIntro string
	newReal  string
	newNick  string
}

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
		var userRegisterRequest UserRegisterRequest
		err := c.ShouldBind(&userRegisterRequest)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "参数错误",
			})
			return
		}
		sendCode := userRegisterRequest.sendCode
		userCode := userRegisterRequest.userCode
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
	var identifyRequest IdentifyRequest
	err := c.ShouldBind(&identifyRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "参数错误",
		})
		return
	}
	email := []string{identifyRequest.email}
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	code := fmt.Sprintf("%06v", rnd.Int31n(1000000))
	err = SendEmail(email, code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "参数错误",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "发送成功",
		"code": code,
	})
}

func UserLogin(c *gin.Context) {
	var flag = 0
	var userLoginRequest UserLoginRequest
	err := c.ShouldBind(&userLoginRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "参数错误",
		})
		return
	}
	username := userLoginRequest.userId
	email := userLoginRequest.email
	password := userLoginRequest.pwd
	if username == "" {
		flag = 1
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
	var userInfoRequest UserInfoRequest
	err := c.ShouldBind(&userInfoRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "参数错误",
		})
		return
	}
	userId := userInfoRequest.userId
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
	var userInfoRequest UserInfoRequest
	err := c.ShouldBind(&userInfoRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "参数错误",
		})
		return
	}
	userId := userInfoRequest.userId
	oldPwd := userInfoRequest.oldPwd
	newPwd := userInfoRequest.newPwd
	var user entity.User
	selectErr := entity.Db.Find(&user, "userId=?", userId).Error
	errors.Is(selectErr, gorm.ErrRecordNotFound)
	if selectErr != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "用户不存在",
		})
		return
	}
	if user.Password != oldPwd {
		c.JSON(http.StatusConflict, gin.H{
			"msg": "密码错误",
		})
		return
	}
	entity.Db.Model(&user).Update("password", newPwd)
	c.JSON(http.StatusOK, gin.H{
		"msg": "更新成功",
	})
}
func UserModifyEmail(c *gin.Context) {
	var userInfoRequest UserInfoRequest
	err := c.ShouldBind(&userInfoRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "参数错误",
		})
		return
	}
	userId := userInfoRequest.userId
	newEmail := userInfoRequest.newEmail
	sendCode := userInfoRequest.sendCode
	userCode := userInfoRequest.userCode
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
	entity.Db.Model(&user).Update("email", newEmail)
	c.JSON(http.StatusOK, gin.H{
		"msg": "更新成功",
	})
}
func UserModifyIntroduction(c *gin.Context) {
	var userInfoRequest UserInfoRequest
	err := c.ShouldBind(&userInfoRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "参数错误",
		})
		return
	}
	userId := userInfoRequest.userId
	newIntro := userInfoRequest.newIntro
	var user entity.User
	selectErr := entity.Db.Find(&user, "userId=?", userId).Error
	errors.Is(selectErr, gorm.ErrRecordNotFound)
	if selectErr != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "用户不存在",
		})
		return
	}
	entity.Db.Model(&user).Update("intro", newIntro)
	c.JSON(http.StatusOK, gin.H{
		"msg": "更新成功",
	})
}
func UserModifyNickname(c *gin.Context) {
	var userInfoRequest UserInfoRequest
	err := c.ShouldBind(&userInfoRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "参数错误",
		})
		return
	}
	userId := userInfoRequest.userId
	newNick := userInfoRequest.newNick
	var user entity.User
	selectErr := entity.Db.Find(&user, "userId=?", userId).Error
	errors.Is(selectErr, gorm.ErrRecordNotFound)
	if selectErr != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "用户不存在",
		})
		return
	}
	entity.Db.Model(&user).Update("nickname", newNick)
	c.JSON(http.StatusOK, gin.H{
		"msg": "更新成功",
	})
}
func UserModifyRealname(c *gin.Context) {
	var userInfoRequest UserInfoRequest
	err := c.ShouldBind(&userInfoRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "参数错误",
		})
		return
	}
	userId := userInfoRequest.userId
	newReal := userInfoRequest.newReal
	var user entity.User
	selectErr := entity.Db.Find(&user, "userId=?", userId).Error
	errors.Is(selectErr, gorm.ErrRecordNotFound)
	if selectErr != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "用户不存在",
		})
		return
	}
	entity.Db.Model(&user).Update("realname", newReal)
	c.JSON(http.StatusOK, gin.H{
		"msg": "更新成功",
	})
}
func UserTeam(c *gin.Context) {
	var userInfoRequest UserInfoRequest
	err := c.ShouldBind(&userInfoRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "参数错误",
		})
		return
	}
	userId := userInfoRequest.userId
	var teams []entity.Team
	selectErr := entity.Db.Table("teams").Select("teams.team_id as team_id,teams.name as name,teams.intro as intro").Joins("left join team_member on team_member.team_id = teams.team_id where team_member.user_id <> ? ", userId).Scan(&teams).Error
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
