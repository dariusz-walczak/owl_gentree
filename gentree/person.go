package main

import (
	"fmt"
	"github.com/gin-contrib/location"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

/* Intermediate structure used to bind person payload and respond with person data */
type personPayload struct {
	Id      string `json:"id" binding:"required,alphanum|uuid"`
	Given   string `json:"given_names"`
	Surname string `json:"surname"`
	Gender  string `json:"gender" binding:"isdefault|oneof=male female unknown"`
}

/* Create a person record from a person payload

   This function is used by request handlers when communicating with the storage backend

   Returns:
   * person record */
func (p *personPayload) toRecord() personRecord {
	gender := p.Gender

	if gender == "" {
		gender = gUnknown
	}

	return personRecord{p.Id, p.Given, p.Surname, gender}
}

/* Convert a person record to person payload

   This function is used by request handlers when responding with data provided by the storage
   backend.

   Returns:
   * relation payload */
func (r *personRecord) toPayload() personPayload {
	return personPayload{r.Id, r.Given, r.Surname, r.Gender}
}

/* Convert a list of person records to payload

   This function is used by request handlers when responding with data provided by the storage
   backend.

   Returns:
   * slice of person payload structures */
func (list personList) toPayload() []personPayload {
	payload := make([]personPayload, 0, len(list))

	for _, r := range list {
		payload = append(payload, r.toPayload())
	}

	return payload
}

/* A structure used to extract person identifier from a URI */
type specifyPersonUri struct {
	Pid string `uri:"pid" binding:"required,alphanum|uuid"`
}

/* Handle a create person request

   The function will retrieve all the input data from the request payload (personPayload) */
func createPerson(c *gin.Context) {
	log.Trace("Entry checkpoint")

	var person personPayload

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

	people[person.Id] = person.toRecord()

	c.JSON(http.StatusCreated, gin.H{"message": "ok"})

	log.Infof("Created a new person (%s) record", person.Id)
}

/* Handle a replace person request

   The function will extract the person id from the request URI (specifyPersonUri), and the rest of
   the data from the request payload (personPayload) */
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

	var person personPayload

	if err := c.ShouldBindJSON(&person); err != nil {
		log.Infof("Person data unmarshalling error: %s", err)

		c.JSON(http.StatusBadRequest, gin.H{"message": payloadErrorMsg})
		return
	}

	people[person.Id] = person.toRecord()

	c.JSON(http.StatusOK, gin.H{"message": "Person record replaced"})

	log.Infof("Replaced the person (%s) record", person.Id)
}

/* Handle a retrieva person request

   The function will extract the person id from the request URI (specifyPersonUri) */
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

	c.JSON(http.StatusOK, person.toPayload())

	log.Infof("Found the requested person record (%s)", params.Pid)
}

/* Handle a retrieve all people request */
func retrievePeople(c *gin.Context) {
	log.Trace("Retrieving all the person records")

	var pagQuery paginationQuery

	if err := c.ShouldBindQuery(&pagQuery); err != nil {
		log.Infof("Query parameters unmarshalling error: %s", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": queryErrorMsg})
		return
	}

	people, pagData, err := queryPeople(pagQuery.toPaginationData())

	if err != nil {
		log.Errorf("An error occurred during people retrieval attempt (%s)", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": internalErrorMsg})
		return
	}

	reqUrl := location.Get(c)
	reqUrl.Path = "/people"

	c.JSON(http.StatusOK, gin.H{
		"pagination": pagData.getJson(*reqUrl),
		"records":    people.toPayload(),
	})

	log.Infof("Found %d persons", len(relations))
}

/* Handle a delete person request

   The function will extract the person id from the request URI (specifyPersonUri) */
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
