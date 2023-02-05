package main

import (
	"aaa/common"
	"github.com/gin-gonic/gin"
)

func main() {
	print("2")
	db := common.InitDB()
	println(db, "22")

	r := gin.Default()
	r = CollectRoute(r)
	panic(r.Run()) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
