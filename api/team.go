package api

import (
	"backend/entity"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"io"
	"log"
	"net/http"
	"os"
)

type TeamInfoRequest struct {
	UserId     string `json:"userId"`
	TeamId     int    `json:"teamId"`
	OperatorId string `json:"operatorId"`
	MemberId   string `json:"memberId"`
	NewIntro   string `json:"newIntro"`
	NewName    string `json:"newName"`
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
		entity.Db.Create(&team)
		teamMember := entity.TeamMember{
			TeamId:   team.ID,
			MemberId: team.CaptainID,
			Identity: 0,
		}
		entity.Db.Create(&teamMember)
		c.JSON(200, gin.H{
			"msg":    "创建成功",
			"teamId": team.ID,
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
	selectErr := entity.Db.Find(&team, "team_id=?", teamId).Error
	errors.Is(selectErr, gorm.ErrRecordNotFound)
	if selectErr != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "团队不存在",
		})
		return
	}
	if team.CaptainID != userId {
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
	selectErr := entity.Db.Find(&team, "team_id=?", teamId).Error
	errors.Is(selectErr, gorm.ErrRecordNotFound)
	if selectErr != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "团队不存在",
		})
		return
	}
	entity.Db.Model(&team).Update("name", newName)
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
	selectErr := entity.Db.Find(&team, "team_id=?", teamId).Error
	errors.Is(selectErr, gorm.ErrRecordNotFound)
	if selectErr != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "团队不存在",
		})
		return
	}
	entity.Db.Model(&team).Update("intro", newIntro)
	c.JSON(http.StatusOK, gin.H{
		"msg": "更新成功",
	})
	return
}
func TeamInformation(c *gin.Context) {
	teamId, ok := c.GetPostForm("teamId")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "参数错误",
		})
		return
	}
	var team entity.Team
	selectErr := entity.Db.Find(&team, "team_id=?", teamId).Error
	errors.Is(selectErr, gorm.ErrRecordNotFound)
	if selectErr != nil {
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
	var members []entity.TeamMember
	selectErr := entity.Db.Table("users").Select("user.user_id as user_id,user.nickname as name,user.intro as intro,team_member.identity as identity").Joins("left join team_member on team_member.user_id = users.user_id where team_member.team_id <> ? ", teamId).Scan(&members).Error
	errors.Is(selectErr, gorm.ErrRecordNotFound)
	if selectErr != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "团队不存在",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":     "查找成功",
		"members": members,
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
	selectErr := entity.Db.Find(&team, "team_id=?", teamId).Error
	errors.Is(selectErr, gorm.ErrRecordNotFound)
	if selectErr != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "团队不存在",
		})
		return
	}
	var operator entity.TeamMember
	selectErr = entity.Db.Find(&operator, "team_id=? and member_id=?", teamId, operatorId).Error
	errors.Is(selectErr, gorm.ErrRecordNotFound)
	if selectErr != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "操作者不在团队中",
		})
		return
	}
	var member entity.TeamMember
	selectErr = entity.Db.Find(&member, "team_id=? and member_id=?", teamId, memberId).Error
	errors.Is(selectErr, gorm.ErrRecordNotFound)
	if selectErr != nil {
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
	entity.Db.Where("team_id = ? and member-id = ?", teamId, memberId).Delete(&entity.TeamMember{})
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
	selectErr := entity.Db.Find(&team, "team_id=?", teamId).Error
	errors.Is(selectErr, gorm.ErrRecordNotFound)
	if selectErr != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "团队不存在",
		})
		return
	}
	var operator entity.TeamMember
	selectErr = entity.Db.Find(&operator, "team_id=? and member_id=?", teamId, operatorId).Error
	errors.Is(selectErr, gorm.ErrRecordNotFound)
	if selectErr != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "操作者不在团队中",
		})
		return
	}
	var member entity.TeamMember
	selectErr = entity.Db.Find(&member, "team_id=? and member_id=?", teamId, memberId).Error
	errors.Is(selectErr, gorm.ErrRecordNotFound)
	if selectErr != nil {
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
	temp := operator.Identity
	entity.Db.Model(&operator).Update("identity", member.Identity)
	entity.Db.Model(&member).Update("identity", temp)
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
	selectErr := entity.Db.Find(&team, "team_id=?", teamId).Error
	errors.Is(selectErr, gorm.ErrRecordNotFound)
	if selectErr != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "团队不存在",
		})
		return
	}
	var operator entity.TeamMember
	selectErr = entity.Db.Find(&operator, "team_id=? and member_id=?", teamId, operatorId).Error
	errors.Is(selectErr, gorm.ErrRecordNotFound)
	if selectErr != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "操作者不在团队中",
		})
		return
	}
	var member entity.TeamMember
	selectErr = entity.Db.Find(&member, "team_id=? and member_id=?", teamId, memberId).Error
	errors.Is(selectErr, gorm.ErrRecordNotFound)
	if selectErr != nil {
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
	entity.Db.Model(&member).Update("identity", 1)
	c.JSON(http.StatusOK, gin.H{
		"msg": "设置成功",
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
	selectErr := entity.Db.Find(&team, "team_id=?", teamId).Error
	errors.Is(selectErr, gorm.ErrRecordNotFound)
	if selectErr != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "团队不存在",
		})
		return
	}
	var operator entity.TeamMember
	selectErr = entity.Db.Find(&operator, "team_id=? and member_id=?", teamId, userId).Error
	errors.Is(selectErr, gorm.ErrRecordNotFound)
	if selectErr != nil {
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
	entity.Db.Where("team_id = ?and memberId=?", teamId, userId).Delete(&entity.TeamMember{})
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
	teamId, ok := c.GetPostForm("teamId")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "参数错误",
		})
		return
	}
	var team entity.Team
	selectErr := entity.Db.Find(&team, "teamId=?", teamId).Error
	errors.Is(selectErr, gorm.ErrRecordNotFound)
	if selectErr != nil {
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
	selectErr := entity.Db.Find(&team, "teamId=?", teamId).Error
	errors.Is(selectErr, gorm.ErrRecordNotFound)
	if selectErr != nil {
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
	selectErr := entity.Db.Find(&team, "teamId=?", teamId).Error
	errors.Is(selectErr, gorm.ErrRecordNotFound)
	if selectErr != nil {
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
	selectErr := entity.Db.Find(&team, "teamId=?", teamId).Error
	errors.Is(selectErr, gorm.ErrRecordNotFound)
	if selectErr != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "团队不存在",
		})
		return
	}
	var user entity.User
	selectErr = entity.Db.Find(&user, "userId=?", userId).Error
	errors.Is(selectErr, gorm.ErrRecordNotFound)
	if selectErr != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "用户不存在",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "请求成功",
	})
	return
}
