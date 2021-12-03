package main

import (
	"os"
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


func configLogger(args AppArgs) {
	log.SetLevel(args.LogLevel)
}


func main() {
	args, err := parseArgs()

	if err != nil {
		os.Exit(1)
	}

	configLogger(args)

	log.Trace("Entry checkpoint")
	r := gin.Default()
	r.POST("/person", createPerson)
	r.Run()
}
