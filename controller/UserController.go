package controller

import (
	"aaa/common"
	"aaa/model"
	"aaa/util"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
)

func Register(c *gin.Context) {
	db := common.GetDB()
	name := c.PostForm("name")
	telephone := c.PostForm("telephone")
	password := c.PostForm("password")
	if len(telephone) != 11 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "no 11 number"})
		return
	}
	if len(password) < 6 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "password too short"})
		return
	}
	if len(name) == 0 {
		name = util.RandomName(10)
	}

	if isTelephoneExit(db, telephone) {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "User is Already"})
		return
	}
	newUser := &model.User{
		Name:      name,
		Telephone: telephone,
		Password:  password,
	}

	db.Create(&newUser)
	log.Println(name, telephone, password)
	c.JSON(http.StatusOK, gin.H{
		"message": "Sign In Successful",
	})
}

func isTelephoneExit(db *gorm.DB, telephone string) bool {
	var user *model.User
	db.Where("telephone = ?", telephone).First(&user)
	if user.ID != 0 {
		return true
	}
	return false
}
