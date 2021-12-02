package main

import (
	"log"
	"github.com/gin-gonic/gin"
)

func hello(c *gin.Context) {
	log.Println("Hello: Entry checkpoint")
	c.JSON(200, gin.H{
		"message": "hello",
	})

	log.Println("Hello: Exit checkpoint")
}

func main() {
	log.Println("Entry checkpoint")
	r := gin.Default()
	r.GET("/hello", hello)
	r.Run()
	log.Println("Exit checkpoint")
}
