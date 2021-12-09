package main

import (
	"github.com/gin-contrib/location"
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
	r.Use(location.Default())

	r.POST("/people", createPerson)
	r.GET("/people/:pid", retrievePerson)
	r.DELETE("/people/:pid", deletePerson)
	r.PUT("/people/:pid", replacePerson)

	r.POST("/relations", createRelation)
	r.GET("/relations/:rid", retrieveRelation)
	r.POST("/people/:pid/relations", createPersonRelation)
	r.GET("/people/:pid/relations", retrievePersonRelations)

	if err := r.Run(); err != nil {
		log.Fatalf("An error occurred during the gin server run attempt (%s)", err)
	}
}
