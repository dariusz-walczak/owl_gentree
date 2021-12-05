package main

import (
	"os"
	log "github.com/sirupsen/logrus"
	"github.com/gin-gonic/gin"
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
	r.Run()
}
