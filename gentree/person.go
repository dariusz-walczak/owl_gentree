package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type personRecord struct {
	Id      string `json:"id" binding:"required,alphanum|uuid"`
	Given   string `json:"given_names"`
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

	if err := c.ShouldBindJSON(&person); err != nil {
		log.Infof("New person data unmarshalling error: %s", err)

		c.JSON(http.StatusBadRequest, gin.H{"message": payloadErrorMsg})
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

type specifyPersonUri struct {
	Pid string `uri:"id" binding:"required,alphanum|uuid"`
}

func replacePerson(c *gin.Context) {
	log.Trace("Entry checkpoint")

	var params specifyPersonUri

	if err := c.ShouldBindUri(&params); err != nil {
		log.Infof("Uri parameters unmarshalling error: %s", err)

		c.JSON(http.StatusBadRequest, gin.H{"message": uriErrorMsg})
		return
	}

	_, found, err := getPerson(params.Pid)

	if !found {
		log.Infof("The person with given id (%s) doesn't exist and can't be replaced", params.Pid)

		c.JSON(http.StatusNotFound, gin.H{"message": "Unknown person id"})
		return
	} else if err != nil {
		log.Errorf("An error occurred during the person retrieval attempt (%s)", err)

		c.JSON(http.StatusInternalServerError, gin.H{"message": internalErrorMsg})
		return
	}

	var person personRecord

	if err := c.ShouldBindJSON(&person); err != nil {
		log.Infof("Person data unmarshalling error: %s", err)

		c.JSON(http.StatusBadRequest, gin.H{"message": payloadErrorMsg})
		return
	}

	people[person.Id] = person

	c.JSON(http.StatusOK, gin.H{"message": "Person record replaced"})

	log.Infof("Replaced the person (%s) record", person.Id)
}

func retrievePerson(c *gin.Context) {
	log.Trace("Entry checkpoint")

	var params specifyPersonUri

	if err := c.ShouldBindUri(&params); err != nil {
		log.Infof("Uri parameters unmarshalling error: %s", err)

		c.JSON(http.StatusBadRequest, gin.H{"message": uriErrorMsg})
		return
	}

	person, found, err := getPerson(params.Pid)

	if !found {
		log.Infof("The person with given id (%s) doesn't exist", params.Pid)

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

	log.Infof("Found the requested person record (%s)", params.Pid)
}

func deletePerson(c *gin.Context) {
	log.Trace("Entry checkpoint")

	var params specifyPersonUri

	if err := c.ShouldBindUri(&params); err != nil {
		log.Infof("Uri parameters unmarshalling error: %s", err)

		c.JSON(http.StatusBadRequest, gin.H{"message": uriErrorMsg})
		return
	}

	_, found, err := getPerson(params.Pid)

	if !found {
		log.Infof("The person with given id (%s) doesn't exist", params.Pid)

		c.JSON(http.StatusNotFound, gin.H{"message": "Unknown person id"})
		return
	} else if err != nil {
		log.Errorf("An error occurred during the person retrieval attempt (%s)", err)

		c.JSON(http.StatusInternalServerError, gin.H{"message": internalErrorMsg})
		return
	}

	delete(people, params.Pid)

	log.Infof("Deleted the requested person record (%s)", params.Pid)
}
