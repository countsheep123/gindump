package main

import (
	"log"
	"net/http"

	"github.com/countsheep123/gindump"
	"github.com/gin-gonic/gin"
)

var logger = func(c *gin.Context) {
	c.Next()
	req, err := gindump.GetRequest(c)
	if err != nil {
		log.Println(err.Error())
	} else {
		log.Println(req)
	}
	res, err := gindump.GetResponse(c)
	if err != nil {
		log.Println(err.Error())
	} else {
		log.Println(res)
	}
}

func main() {
	r := gin.Default()
	r.Use(logger)
	r.Use(gindump.Dump())
	r.POST("/", handle)
	r.Run(":8080")
}

func handle(c *gin.Context) {
	var r map[string]interface{}
	if err := c.ShouldBindJSON(&r); err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"hello": "world",
	})
}
