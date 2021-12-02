package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

func createPerson(c *gin.Context) {
	log.Trace("Entry checkpoint")

	c.JSON(201, gin.H{
		"message": "hello",
	})

	log.Info("Created new person resource")

	log.Trace("Exit checkpoint")
}

func main() {
	log.Trace("Entry checkpoint")
	r := gin.Default()
	r.POST("/person", createPerson)
	r.Run()
	log.Trace("Exit checkpoint")
}
