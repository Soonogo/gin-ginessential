package controller

import (
	"aaa/common"
	"aaa/model"
	"aaa/util"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
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

	hasedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"code": 500, "msg": "encryption error"})
		return
	}
	newUser := &model.User{
		Name:      name,
		Telephone: telephone,
		Password:  string(hasedPassword),
	}

	db.Create(&newUser)
	log.Println(name, telephone, password)
	c.JSON(http.StatusOK, gin.H{
		"code":    "200",
		"message": "Sign In Successful",
	})
}

func Login(c *gin.Context) {
	db := common.GetDB()
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
	var user model.User
	db.Where("telephone = ?", telephone).First(&user)
	if user.Id == 0 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "User not register"})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "400", "msg": "password error"})
		return
	}
	token, err := common.ReleaseToken(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "500", "msg": "system error"})
		log.Println("token generate error : %v", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    "200",
		"data":    gin.H{"token": token},
		"message": "Login In Successful",
	})
}

func Info(ctx *gin.Context) {
	user, _ := ctx.Get("user")
	ctx.JSON(http.StatusOK, gin.H{"code": "200", "data": gin.H{"user": user}})
	return
}

func isTelephoneExit(db *gorm.DB, telephone string) bool {
	var user model.User
	db.Where("telephone = ?", telephone).First(&user)
	println(user.Id, telephone, "userid")
	log.Println(user)
	if user.Id != 0 {
		return true
	}
	return false
}
