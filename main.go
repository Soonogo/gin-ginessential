package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type User struct {
	gorm.Model
	Name      string `gorm:"type:varchar(20);not null"`
	Telephone string `gorm:"type:varchar(11);not null;unique"`
	Password  string `gorm:"size:255;not null"`
}

func main() {
	db := InitDB()
	close1, _ := db.DB()
	defer close1.Close()
	r := gin.Default()
	r.POST("/api/auth/register", func(c *gin.Context) {
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
			name = RandomName(10)
		}

		if isTelephoneExit(db, telephone) {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "User is Already"})
			return
		}
		newUser := &User{
			Name:      name,
			Telephone: telephone,
			Password:  password,
		}

		db.Create(&newUser)
		log.Println(name, telephone, password)
		c.JSON(http.StatusOK, gin.H{
			"message": "Sign In Successful",
		})
	})
	panic(r.Run()) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func isTelephoneExit(db *gorm.DB, telephone string) bool {
	var user User
	db.Where("telephone = ?", telephone).First(&user)
	if user.ID != 0 {
		return true
	}
	return false
}

func RandomName(n int) string {
	var letters = []byte("asdfghjklqwertyuiopzxcvbnmASDFGHJKLQWERTYUIOPZXCVBNM")
	result := make([]byte, n)

	rand.Seed(time.Now().Unix())
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}

func InitDB() *gorm.DB {
	host := "localhost"
	port := "9999"
	database := "ginessential"
	username := "root"
	password := "root"
	charset := "utf8"
	args := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true",
		username,
		password,
		host,
		port,
		database,
		charset)
	db, err := gorm.Open(mysql.Open(args), &gorm.Config{})
	if err != nil {
		panic("failed to connect database,err:" + err.Error())
	}
	db.AutoMigrate(&User{})
	return db
}
