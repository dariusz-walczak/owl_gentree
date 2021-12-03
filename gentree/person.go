package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"net/http"
)


type personRecord struct {
	Id string `json:"id"`
	Given string `json:"given_names"`
	Surname string `json:"surname"`
}


func createPerson(c *gin.Context) {
	log.Trace("Entry checkpoint")

	var person personRecord

    if err := c.BindJSON(&person); err != nil {
		log.Infof("New person data unmarshalling error: %s", err)

		c.JSON(http.StatusBadRequest, gin.H{"message": "Payload format error"})
		return
    }

	log.Infof("id: %s", person.Id)
	log.Infof("given name(s): %s", person.Given)
	log.Infof("surname: %s", person.Surname)

	c.JSON(http.StatusCreated, gin.H{"message": "ok"})

	log.Info("Created new person resource")
}
