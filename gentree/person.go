package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"net/http"
)


type personRecord struct {
	Id string `json:"id" binding:"required,min=3"`
	Given string `json:"given_names"`
	Surname string `json:"surname"`
}


var people = map[string]personRecord{}


/* Retrieve a person record by id
 * Returns:
 * * Person record structure (uninitialized if not found)
 * * Success flag (true if the record was found and false otherwise)
 * * Error (if occurred) */
func getPerson(pid string) (personRecord, bool, error) {
	log.Debugf("Retrieving person record by id (%s)", pid)

	person, found := people[pid]

	if !found {
		log.Debugf("Person record (%s) not found", pid)

		return person, false, nil
	}

	return person, true, nil
}


func createPerson(c *gin.Context) {
	log.Trace("Entry checkpoint")

	var person personRecord

    if err := c.BindJSON(&person); err != nil {
		log.Infof("New person data unmarshalling error: %s", err)

		c.JSON(http.StatusBadRequest, gin.H{"message": "Payload format error"})
		return
    }

	if _, found, err := getPerson(person.Id); found {
		log.Infof("A person with given id (%s) already exists", person.Id)

		c.JSON(
			http.StatusBadRequest,
			gin.H{"message": fmt.Sprintf("Person (%s) already exists", person.Id)})
		return
	} else if err != nil {
		log.Errorf("An error occurred during the person retrieval attempt (%s)", err)

		c.JSON(http.StatusInternalServerError, gin.H{"message": internalErrorMsg})
		return
	}

	people[person.Id] = person

	c.JSON(http.StatusCreated, gin.H{"message": "ok"})

	log.Infof("Created a new person (%s) record", person.Id)
}


func retrievePerson(c *gin.Context) {
	log.Trace("Entry checkpoint")

	pid := c.Param("id")

	person, found, err := getPerson(pid)

	if !found {
		log.Infof("The person with given id (%s) doesn't exist", pid)

		c.JSON(http.StatusNotFound, gin.H{"message": "Unknown person id"})
		return
	} else if err != nil {
		log.Errorf("An error occurred during the person retrieval attempt (%s)", err)

		c.JSON(http.StatusInternalServerError, gin.H{"message": internalErrorMsg})
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{"id": person.Id, "given_names": person.Given, "surname": person.Surname})

	log.Infof("Found the requested person record (%s)", pid)
}
