package api

import (
	"backend/entity"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
	teamId, ok := c.GetPostForm("teamId")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "参数错误",
		})
		return
	}
	userId, ok := c.GetPostForm("userId")
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
	teamId, ok := c.GetPostForm("teamId")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "参数错误",
		})
		return
	}
	newName, ok := c.GetPostForm("newName")
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
	entity.Db.Model(&team).Update("name", newName)
	c.JSON(http.StatusOK, gin.H{
		"msg": "更新成功",
	})
	return
}
func TeamModifyIntro(c *gin.Context) {
	teamId, ok := c.GetPostForm("teamId")
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
		"msg":   "查找成功",
		"teams": members,
	})

}
func Remove(c *gin.Context) {
	teamId, ok := c.GetQuery("teamId")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "参数错误",
		})
		return
	}
	operatorId, ok := c.GetQuery("operatorId")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "参数错误",
		})
		return
	}
	memberId, ok := c.GetQuery("memberId")
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
	teamId, ok := c.GetQuery("teamId")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "参数错误",
		})
		return
	}
	operatorId, ok := c.GetQuery("operatorId")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "参数错误",
		})
		return
	}
	memberId, ok := c.GetQuery("memberId")
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
	if operator.Identity >= member.Identity && operator.Identity != 0 {
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
func Leave(c *gin.Context) {
	teamId, ok := c.GetQuery("teamId")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "参数错误",
		})
		return
	}
	operatorId, ok := c.GetQuery("operatorId")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "参数错误",
		})
		return
	}
	memberId, ok := c.GetQuery("memberId")
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
	if operator.Identity >= member.Identity && operator.Identity != 0 {
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
