package controller

import (
	"aaa/common"
	"aaa/dto"
	"aaa/model"
	"aaa/response"
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
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "no 11 number")

		return
	}
	if len(password) < 6 {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "password too short")
		return
	}
	if len(name) == 0 {
		name = util.RandomName(10)
	}

	if isTelephoneExit(db, telephone) {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "User is Already")
		return
	}

	hasedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "encryption error")
		return
	}
	newUser := &model.User{
		Name:      name,
		Telephone: telephone,
		Password:  string(hasedPassword),
	}

	db.Create(&newUser)
	log.Println(name, telephone, password)
	response.Success(c, nil, "Sign In Successful")

}

func Login(c *gin.Context) {
	db := common.GetDB()
	telephone := c.PostForm("telephone")
	password := c.PostForm("password")
	if len(telephone) != 11 {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "no 11 number")
		return
	}
	if len(password) < 6 {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "password too short")
		return
	}
	var user model.User
	db.Where("telephone = ?", telephone).First(&user)
	if user.Id == 0 {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "User not register")
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "password error")
		return
	}
	token, err := common.ReleaseToken(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "500", "msg": "system error"})
		log.Println("token generate error : %v", err)
		return
	}

	response.Success(c, gin.H{"token": token}, "Login In Successful")
}

func Info(ctx *gin.Context) {
	user, _ := ctx.Get("user")
	ctx.JSON(http.StatusOK, gin.H{"code": "200", "data": gin.H{"user": dto.ToUserParse(user.(model.User))}})
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
