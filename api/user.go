package api

import (
	"backend/entity"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

type IdentifyRequest struct {
	Email string `json:"email"`
}
type UserRegisterRequest struct {
	UserId   string `json:"userId"`
	Nickname string `json:"nickname"`
	Password string `json:"pwd"`
	Email    string `json:"email"`
	UserCode string `json:"userCode"`
	SendCode string `json:"sendCode"`
}
type UserLoginRequest struct {
	UserId string `form:"userId"`
	Email  string `form:"email"`
	Pwd    string `form:"pwd"`
}
type UserInfoRequest struct {
	UserId   string `json:"userId"`
	OldPwd   string `json:"oldPwd"`
	NewPwd   string `json:"newPwd"`
	UserCode string `json:"userCode"`
	SendCode string `json:"sendCode"`
	NewEmail string `json:"newEmail"`
	NewIntro string `json:"newIntro"`
	NewReal  string `json:"newReal"`
	NewNick  string `json:"newNick"`
}

func UserRegister(c *gin.Context) {
	var user entity.User
	var userRegisterRequest UserRegisterRequest
	err := c.ShouldBindJSON(&userRegisterRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "参数错误",
		})
		return
	}
	user.UserId = userRegisterRequest.UserId
	user.Nickname = userRegisterRequest.Nickname
	user.Email = userRegisterRequest.Email
	user.Password = userRegisterRequest.Password
	user.Url = "_defaultavatar.webp"
	var temp entity.User
	//同一邮箱只允许被注册一次
	entity.Db.Find(&temp, "email=?", user.Email)
	if temp != (entity.User{}) {
		c.JSON(http.StatusFound, gin.H{
			"msg":  "该邮箱已被注册",
			"user": user.Email,
		})
		return
	}
	sendCode := userRegisterRequest.SendCode
	userCode := userRegisterRequest.UserCode
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

func Identifying(c *gin.Context) {
	var identifyRequest IdentifyRequest
	err := c.ShouldBindJSON(&identifyRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "参数错误",
		})
		return
	}
	email := []string{identifyRequest.Email}
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
	err := c.ShouldBindQuery(&userLoginRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "参数错误",
		})
		return
	}
	username := userLoginRequest.UserId
	email := userLoginRequest.Email
	password := userLoginRequest.Pwd
	if username == "" {
		flag = 1
	}
	var loginUser entity.User
	if flag == 0 {
		entity.Db.Find(&loginUser, "user_id=?", username)
	} else {
		entity.Db.Find(&loginUser, "email=?", email)
	}
	if loginUser == (entity.User{}) {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "用户不存在",
		})
		return
	}
	if loginUser.Password != password {
		c.JSON(http.StatusConflict, gin.H{
			"msg": "用户名或密码错误",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":      "登陆成功",
		"userId":   loginUser.UserId,
		"nickName": loginUser.Nickname,
	})
}
func UserInformation(c *gin.Context) {
	var userLoginRequest UserLoginRequest
	err := c.ShouldBindQuery(&userLoginRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "参数错误",
		})
		return
	}
	userId := userLoginRequest.UserId
	var user entity.User
	entity.Db.Find(&user, "user_id=?", userId)
	if user == (entity.User{}) {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "用户不存在",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":      "返回成功",
		"nickname": user.Nickname,
		"realname": user.Realname,
		"email":    user.Email,
		"intro":    user.Intro,
	})
}
func UserModifyPassword(c *gin.Context) {
	var userInfoRequest UserInfoRequest
	err := c.ShouldBindJSON(&userInfoRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "参数错误",
		})
		return
	}
	userId := userInfoRequest.UserId
	oldPwd := userInfoRequest.OldPwd
	newPwd := userInfoRequest.NewPwd
	var user entity.User
	entity.Db.Find(&user, "user_id=?", userId)
	if user == (entity.User{}) {
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
	entity.Db.Model(&user).Where("user_id = ?", userId).Update("password", newPwd)
	c.JSON(http.StatusOK, gin.H{
		"msg": "更新成功",
	})
}
func UserModifyEmail(c *gin.Context) {
	var userInfoRequest UserInfoRequest
	err := c.ShouldBindJSON(&userInfoRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "参数错误",
		})
		return
	}
	userId := userInfoRequest.UserId
	newEmail := userInfoRequest.NewEmail
	sendCode := userInfoRequest.SendCode
	userCode := userInfoRequest.UserCode
	if sendCode != userCode {
		c.JSON(http.StatusConflict, gin.H{
			"msg": "验证码错误",
		})
		return
	}
	var user entity.User
	entity.Db.Find(&user, "user_id=?", userId)
	if user == (entity.User{}) {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "用户不存在",
		})
		return
	}
	entity.Db.Model(&user).Where("user_id = ?", userId).Update("email", newEmail)
	c.JSON(http.StatusOK, gin.H{
		"msg": "更新成功",
	})
}
func UserModifyIntroduction(c *gin.Context) {
	var userInfoRequest UserInfoRequest
	err := c.ShouldBindJSON(&userInfoRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "参数错误",
		})
		return
	}
	userId := userInfoRequest.UserId
	newIntro := userInfoRequest.NewIntro
	var user entity.User
	entity.Db.Find(&user, "user_id=?", userId)
	if user == (entity.User{}) {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "用户不存在",
		})
		return
	}
	entity.Db.Model(&user).Where("user_id = ?", userId).Update("intro", newIntro)
	c.JSON(http.StatusOK, gin.H{
		"msg": "更新成功",
	})
}
func UserModifyNickname(c *gin.Context) {
	var userInfoRequest UserInfoRequest
	err := c.ShouldBindJSON(&userInfoRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "参数错误",
		})
		return
	}
	userId := userInfoRequest.UserId
	newNick := userInfoRequest.NewNick
	var user entity.User
	entity.Db.Find(&user, "user_id=?", userId)
	if user == (entity.User{}) {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "用户不存在",
		})
		return
	}
	entity.Db.Model(&user).Where("user_id = ?", userId).Update("nickname", newNick)
	c.JSON(http.StatusOK, gin.H{
		"msg": "更新成功",
	})
}
func UserModifyRealname(c *gin.Context) {
	var userInfoRequest UserInfoRequest
	err := c.ShouldBindJSON(&userInfoRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "参数错误",
		})
		return
	}
	userId := userInfoRequest.UserId
	newReal := userInfoRequest.NewReal
	var user entity.User
	entity.Db.Find(&user, "user_id=?", userId)
	if user == (entity.User{}) {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "用户不存在",
		})
		return
	}
	entity.Db.Model(&user).Where("user_id = ?", userId).Update("realname", newReal)
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
	var user entity.User
	entity.Db.Find(&user, "user_id=?", userId)
	if user == (entity.User{}) {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "用户不存在",
		})
		return
	}
	var teams []entity.Team
	entity.Db.Table("teams").Select("teams.team_id as team_id,teams.name as name,teams.intro as intro").Joins("join team_members on team_members.team_id = teams.team_id", userId).Where("member_id=?", userId).Scan(&teams)
	var teamList []gin.H
	for _, team := range teams {
		projectJson := gin.H{
			"name":   team.Name,
			"intro":  team.Intro,
			"teamId": team.TeamId,
		}
		teamList = append(teamList, projectJson)
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":   "查找成功",
		"teams": teamList,
	})
}
func UserModifyAvatar(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	filename := header.Filename
	out, err := os.Create("./localFile/user/" + filename)
	if err != nil {
		log.Fatal(err)
	}
	defer func(out *os.File) {
		err := out.Close()
		if err != nil {

		}
	}(out)
	_, err = io.Copy(out, file)
	if err != nil {
		log.Fatal(err)
	}
	userId, ok := c.GetPostForm("userId")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "参数错误",
		})
		return
	}
	var user entity.User
	entity.Db.Find(&user, "user_id=?", userId)
	if user == (entity.User{}) {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "用户不存在",
		})
		return
	}
	entity.Db.Model(&user).Where("user_id = ?", userId).Update("url", filename)
	c.JSON(http.StatusOK, gin.H{
		"msg": "更新成功",
	})
}
func UserGetAvatar(c *gin.Context) {
	userId, ok := c.GetQuery("userId")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "参数错误",
		})
		return
	}
	var user entity.User
	entity.Db.Find(&user, "user_id=?", userId)
	if user == (entity.User{}) {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "用户不存在",
		})
		return
	}
	c.File("./localFile/user/" + user.Url)
}
