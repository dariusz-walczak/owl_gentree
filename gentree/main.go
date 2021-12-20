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

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(location.Default())

	r.DELETE("/people/:pid", deletePerson)
	r.GET("/people", retrievePeople)
	r.GET("/people/:pid", retrievePerson)
	r.POST("/people", createPerson)
	r.PUT("/people/:pid", replacePerson)

	r.DELETE("/relations/:rid", deleteRelation)
	r.GET("/relations", retrieveRelations)
	r.GET("/relations/:rid", retrieveRelation)
	r.POST("/relations", createRelation)
	r.PUT("/relations/:rid", replaceRelation)

	r.GET("/people/:pid/relations", retrievePersonRelations)
	r.POST("/people/:pid/relations", createPersonRelation)

	return r
}

func main() {
	args, err := parseArgs()

	if err != nil {
		os.Exit(1)
	}

	configLogger(args)

	log.Trace("Entry checkpoint")

	router := setupRouter()

	if err := router.Run(); err != nil {
		log.Fatalf("An error occurred during the gin server run attempt (%s)", err)
	}
}
