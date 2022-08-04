package api

import (
	"backend/entity"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

type TeamInfoRequest struct {
	UserId     string `json:"userId"`
	TeamId     int    `json:"teamId"`
	OperatorId string `json:"operatorId"`
	MemberId   string `json:"memberId"`
	NewIntro   string `json:"newIntro"`
	NewName    string `json:"newName"`
}
type TeamMember struct {
	UserId   string `json:"userId"`
	Identity int    `json:"identity"`
	Name     string `json:"name"`
	Intro    string `json:"intro"`
}

func TeamCreate(c *gin.Context) {
	var team entity.Team
	err := c.ShouldBindJSON(&team)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "参数错误",
		})
		return
	} else {
		userId := team.CaptainId
		var user entity.User
		entity.Db.Find(&user, "user_id=?", userId)
		if user == (entity.User{}) {
			c.JSON(http.StatusNotFound, gin.H{
				"msg": "用户不存在",
			})
			return
		}
		team.Url = "_defaultavatar.webp"
		entity.Db.Create(&team)
		teamMember := entity.TeamMember{
			TeamId:   team.TeamId,
			MemberId: team.CaptainId,
			Identity: 0,
		}
		entity.Db.Create(&teamMember)
		c.JSON(200, gin.H{
			"msg":    "创建成功",
			"teamId": team.TeamId,
		})
		return
	}
}
func TeamDismiss(c *gin.Context) {
	var teamInfoRequest TeamInfoRequest
	err := c.ShouldBindJSON(&teamInfoRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "参数错误",
		})
		return
	}
	teamId := teamInfoRequest.TeamId
	userId := teamInfoRequest.UserId
	var team entity.Team
	entity.Db.Find(&team, "team_id=?", teamId)
	if team == (entity.Team{}) {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "团队不存在",
		})
		return
	}
	if team.CaptainId != userId {
		c.JSON(http.StatusConflict, gin.H{
			"msg": "该用户非该团队队长",
		})
		return
	}
	entity.Db.Where("team_id = ?", teamId).Delete(&entity.Team{})
	c.JSON(http.StatusOK, gin.H{
		"msg": "解散成功",
	})
	return
}
func TeamModifyName(c *gin.Context) {
	var teamInfoRequest TeamInfoRequest
	err := c.ShouldBindJSON(&teamInfoRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "参数错误",
		})
		return
	}
	teamId := teamInfoRequest.TeamId
	newName := teamInfoRequest.NewName
	var team entity.Team
	entity.Db.Find(&team, "team_id=?", teamId)
	if team == (entity.Team{}) {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "团队不存在",
		})
		return
	}
	entity.Db.Model(&team).Where("team_id=?", teamId).Update("name", newName)
	c.JSON(http.StatusOK, gin.H{
		"msg": "更新成功",
	})
	return
}
func TeamModifyIntro(c *gin.Context) {
	var teamInfoRequest TeamInfoRequest
	err := c.ShouldBindJSON(&teamInfoRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "参数错误",
		})
		return
	}
	teamId := teamInfoRequest.TeamId
	newIntro := teamInfoRequest.NewIntro
	var team entity.Team
	entity.Db.Find(&team, "team_id=?", teamId)
	if team == (entity.Team{}) {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "团队不存在",
		})
		return
	}
	entity.Db.Model(&team).Where("team_id=?", teamId).Update("intro", newIntro)
	c.JSON(http.StatusOK, gin.H{
		"msg": "更新成功",
	})
	return
}
func TeamInformation(c *gin.Context) {
	teamId, ok := c.GetQuery("teamId")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "参数错误",
		})
		return
	}
	var team entity.Team
	entity.Db.Find(&team, "team_id=?", teamId)
	if team == (entity.Team{}) {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "团队不存在",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":   "查找成功",
		"name":  team.Name,
		"intro": team.Intro,
	})
	return
}
func GetMember(c *gin.Context) {
	teamId, ok := c.GetQuery("teamId")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "参数错误",
		})
		return
	}
	var members []TeamMember
	id, _ := strconv.Atoi(teamId)
	var team entity.Team
	entity.Db.Find(&team, "team_id=?", id)
	if team == (entity.Team{}) {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "团队不存在",
		})
		return
	}
	//selectErr := entity.Db.Table("users").Select("users.user_id as user_id,users.nickname as name,users.intro as intro,team_members.identity as identity").Joins("join team_members on team_members.member_id = users.user_id where team_members.team_id <> ? ", id).Scan(&members).Error
	entity.Db.Table("users").Select("users.user_id as user_id,users.nickname as name,users.intro as intro,team_members.identity as identity").Joins("join team_members on team_members.member_id = users.user_id ").Where("team_id = ?", id).Scan(&members)
	var memberList []gin.H
	for _, member := range members {
		projectJson := gin.H{
			"name":     member.Name,
			"intro":    member.Intro,
			"userId":   member.UserId,
			"identity": member.Identity,
		}
		memberList = append(memberList, projectJson)
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":     "查找成功",
		"members": memberList,
	})
}
func Remove(c *gin.Context) {
	var teamInfoRequest TeamInfoRequest
	err := c.ShouldBindJSON(&teamInfoRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "参数错误",
		})
		return
	}
	teamId := teamInfoRequest.TeamId
	operatorId := teamInfoRequest.OperatorId
	memberId := teamInfoRequest.MemberId
	var team entity.Team
	entity.Db.Find(&team, "team_id=?", teamId)
	if team == (entity.Team{}) {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "团队不存在",
		})
		return
	}
	var operator entity.TeamMember
	entity.Db.Find(&operator, "team_id=? and member_id=?", teamId, operatorId)
	if operator == (entity.TeamMember{}) {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "操作者不在团队中",
		})
		return
	}
	var member entity.TeamMember
	entity.Db.Find(&member, "team_id=? and member_id=?", teamId, memberId)
	if member == (entity.TeamMember{}) {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "用户不在团队中",
		})
		return
	}
	if operator.Identity >= member.Identity {
		c.JSON(http.StatusConflict, gin.H{
			"msg": "无移除该用户权限",
		})
		return
	}
	entity.Db.Where("team_id = ? and member_id = ?", teamId, memberId).Delete(&entity.TeamMember{})
	c.JSON(http.StatusOK, gin.H{
		"msg": "移除成功",
	})
	return
}
func Transfer(c *gin.Context) {
	var teamInfoRequest TeamInfoRequest
	err := c.ShouldBindJSON(&teamInfoRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "参数错误",
		})
		return
	}
	teamId := teamInfoRequest.TeamId
	operatorId := teamInfoRequest.OperatorId
	memberId := teamInfoRequest.MemberId
	var team entity.Team
	entity.Db.Find(&team, "team_id=?", teamId)
	if team == (entity.Team{}) {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "团队不存在",
		})
		return
	}
	var operator entity.TeamMember
	entity.Db.Find(&operator, "team_id=? and member_id=?", teamId, operatorId)
	if operator == (entity.TeamMember{}) {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "操作者不在团队中",
		})
		return
	}
	var member entity.TeamMember
	entity.Db.Find(&member, "team_id=? and member_id=?", teamId, memberId)
	if member == (entity.TeamMember{}) {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "用户不在团队中",
		})
		return
	}
	if operator.Identity >= member.Identity || operator.Identity == 0 {
		c.JSON(http.StatusConflict, gin.H{
			"msg": "无转移权限",
		})
		return
	}
	entity.Db.Model(&operator).Where("team_id=? and member_id=?", teamId, operatorId).Update("identity", 2)
	entity.Db.Model(&member).Where("team_id=? and member_id=?", teamId, memberId).Update("identity", 1)
	c.JSON(http.StatusOK, gin.H{
		"msg": "转移成功",
	})
	return
}
func SetAdmin(c *gin.Context) {
	var teamInfoRequest TeamInfoRequest
	err := c.ShouldBindJSON(&teamInfoRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "参数错误",
		})
		return
	}
	teamId := teamInfoRequest.TeamId
	operatorId := teamInfoRequest.OperatorId
	memberId := teamInfoRequest.MemberId
	var team entity.Team
	entity.Db.Find(&team, "team_id=?", teamId)
	if team == (entity.Team{}) {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "团队不存在",
		})
		return
	}
	var operator entity.TeamMember
	entity.Db.Find(&operator, "team_id=? and member_id=?", teamId, operatorId)
	if operator == (entity.TeamMember{}) {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "操作者不在团队中",
		})
		return
	}
	var member entity.TeamMember
	entity.Db.Find(&member, "team_id=? and member_id=?", teamId, memberId)
	if member == (entity.TeamMember{}) {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "用户不在团队中",
		})
		return
	}
	if operator.Identity != 0 {
		c.JSON(http.StatusConflict, gin.H{
			"msg": "无设置管理员权限",
		})
		return
	}
	if member.Identity == 1 {
		c.JSON(http.StatusConflict, gin.H{
			"msg": "该用户已是管理员",
		})
		return
	}
	entity.Db.Model(&member).Where("team_id=? and member_id=?", teamId, memberId).Update("identity", 1)
	c.JSON(http.StatusOK, gin.H{
		"msg": "设置成功",
	})
	return
}
func RemoveAdmin(c *gin.Context) {
	var teamInfoRequest TeamInfoRequest
	err := c.ShouldBindJSON(&teamInfoRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "参数错误",
		})
		return
	}
	teamId := teamInfoRequest.TeamId
	operatorId := teamInfoRequest.OperatorId
	memberId := teamInfoRequest.MemberId
	var team entity.Team
	entity.Db.Find(&team, "team_id=?", teamId)
	if team == (entity.Team{}) {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "团队不存在",
		})
		return
	}
	var operator entity.TeamMember
	entity.Db.Find(&operator, "team_id=? and member_id=?", teamId, operatorId)
	if operator == (entity.TeamMember{}) {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "操作者不在团队中",
		})
		return
	}
	var member entity.TeamMember
	entity.Db.Find(&member, "team_id=? and member_id=?", teamId, memberId)
	if member == (entity.TeamMember{}) {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "用户不在团队中",
		})
		return
	}
	if operator.Identity != 0 {
		c.JSON(http.StatusConflict, gin.H{
			"msg": "无设置管理员权限",
		})
		return
	}
	if member.Identity == 0 {
		c.JSON(http.StatusConflict, gin.H{
			"msg": "该用户不是管理员",
		})
		return
	}
	entity.Db.Model(&member).Where("team_id=? and member_id=?", teamId, memberId).Update("identity", 2)
	c.JSON(http.StatusOK, gin.H{
		"msg": "撤销成功",
	})
	return
}
func Leave(c *gin.Context) {
	var teamInfoRequest TeamInfoRequest
	err := c.ShouldBindJSON(&teamInfoRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "参数错误",
		})
		return
	}
	teamId := teamInfoRequest.TeamId
	userId := teamInfoRequest.UserId
	var team entity.Team
	entity.Db.Find(&team, "team_id=?", teamId)
	if team == (entity.Team{}) {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "团队不存在",
		})
		return
	}
	var operator entity.TeamMember
	entity.Db.Find(&operator, "team_id=? and member_id=?", teamId, userId)
	if operator == (entity.TeamMember{}) {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "用户不在团队中",
		})
		return
	}
	if operator.Identity == 0 {
		c.JSON(http.StatusConflict, gin.H{
			"msg": "创建者无法离开团队",
		})
		return
	}
	entity.Db.Where("team_id = ? and member_id=?", teamId, userId).Delete(&entity.TeamMember{})
	c.JSON(http.StatusOK, gin.H{
		"msg": "离开成功",
	})
	return
}
func TeamModifyAvatar(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	filename := header.Filename
	out, err := os.Create("./localFile/team/" + filename)
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
	teamId, ok := c.GetPostForm("team_id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "参数错误",
		})
		return
	}
	var team entity.Team
	entity.Db.Find(&team, "teamId=?", teamId)
	if team == (entity.Team{}) {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "团队不存在",
		})
		return
	}
	entity.Db.Model(&team).Where("team_id=?", teamId).Update("url", filename)
	c.JSON(http.StatusOK, gin.H{
		"msg": "更新成功",
	})
}
func TeamGetAvatar(c *gin.Context) {
	teamId, ok := c.GetQuery("teamId")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "参数错误",
		})
		return
	}
	var team entity.Team
	entity.Db.Find(&team, "team_id=?", teamId)
	if team == (entity.Team{}) {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "团队不存在",
		})
		return
	}
	c.File("./localFile/team/" + team.Url)
	c.JSON(http.StatusOK, gin.H{
		"msg": "成功",
	})
}
func Confirm(c *gin.Context) {
	var teamInfoRequest TeamInfoRequest
	err := c.ShouldBindJSON(&teamInfoRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "参数错误",
		})
		return
	}
	teamId := teamInfoRequest.TeamId
	userId := teamInfoRequest.UserId
	var team entity.Team
	entity.Db.Find(&team, "team_id=?", teamId)
	if team == (entity.Team{}) {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "团队不存在",
		})
		return
	}
	teamMember := entity.TeamMember{
		TeamId:   teamId,
		MemberId: userId,
		Identity: 2,
	}
	entity.Db.Create(&teamMember)
	c.JSON(http.StatusOK, gin.H{
		"msg": "加入成功",
	})
	return
}
func Apply(c *gin.Context) {
	var teamInfoRequest TeamInfoRequest
	err := c.ShouldBindJSON(&teamInfoRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "参数错误",
		})
		return
	}
	teamId := teamInfoRequest.TeamId
	userId := teamInfoRequest.UserId
	var team entity.Team
	entity.Db.Find(&team, "team_id=?", teamId)
	if team == (entity.Team{}) {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "团队不存在",
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
	var member entity.TeamMember
	entity.Db.First(&member, "member_id=? and team_id =?", userId, teamId)
	if member != (entity.TeamMember{}) {
		c.JSON(http.StatusOK, gin.H{
			"msg":  "用户已在团队中",
			"code": 1,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "请求成功",
		"code": 0,
	})
	return
}
func getAdminNum(c *gin.Context) {
	teamId, ok := c.GetQuery("teamId")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "参数错误",
		})
		return
	}
	var team entity.Team
	entity.Db.Find(&team, "team_id=?", teamId)
	if team == (entity.Team{}) {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "团队不存在",
		})
		return
	}
	var members []entity.TeamMember
	entity.Db.Where("team_id = ?", teamId).Find(&members)
	var cnt = 0
	for _, member := range members {
		if member.Identity == 1 {
			cnt++
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "查找成功",
		"num": cnt,
	})
}
func getIdentity(c *gin.Context) {
	teamId, ok := c.GetQuery("teamId")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "参数错误",
		})
		return
	}
	userId, ok := c.GetQuery("userId")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "参数错误",
		})
		return
	}
	var team entity.Team
	entity.Db.Find(&team, "team_id=?", teamId)
	if team == (entity.Team{}) {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "团队不存在",
		})
		return
	}
	var member entity.TeamMember
	entity.Db.Find(&member, "member_id=? and team_id =?", userId, teamId)
	if member == (entity.TeamMember{}) {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "用户不在团队中",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":      "查找成功",
		"identity": member.Identity,
	})
}
