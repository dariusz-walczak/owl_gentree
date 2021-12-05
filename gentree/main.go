package main

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"os"
)

func configLogger(args AppArgs) {
	log.SetLevel(args.LogLevel)
	log.SetReportCaller(true)
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
	r.GET("/person/:id", retrievePerson)
	r.DELETE("/person/:id", deletePerson)
	r.PUT("/person/:id", replacePerson)
	r.Run()
}
