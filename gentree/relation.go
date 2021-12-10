package main

import (
	"fmt"
	"github.com/gin-contrib/location"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type relationPayload struct {
	Id   int64  `json:"id" binding:"required"`
	Pid1 string `json:"pid1" binding:"required,alphanum|uuid"`
	Pid2 string `json:"pid2" binding:"required,alphanum|uuid"`
	Type string `json:"type" binding:"oneof=father mother husband"`
}

/* Create a relation record from a payload struct

   Used by request handlers when communicating with the storage backend.

   Returns:
   * relation record (used for interaction with the storage backend) */
func (p *relationPayload) toRelationRecord() relationRecord {
	return relationRecord{p.Id, p.Pid1, p.Pid2, p.Type}
}

/* Convert a relation record to payload data

   Used by request handlers when responding with data provided by the storage backend.

   Returns:
   * relation payload */
func (r *relationRecord) toPayload() relationPayload {
	return relationPayload{r.Id, r.Pid1, r.Pid2, r.Type}
}

/* Convert a list of relation records to payload data

   Used by request handlers when responding with data provided by the storage backend.

   Returns:
   * slice of relation payload structures */
func (relations relationList) toPayload() []relationPayload {
	payload := make([]relationPayload, 0, len(relations))

	for _, r := range relations {
		payload = append(payload, r.toPayload())
	}

	return payload
}

/* Relation payload accepted by the createPersonRelation handler

   The 'it' prefix stands for "id and type" */
type itRelationPayload struct {
	// Target person identifier
	Pid  string `json:"pid" binding:"required,alphanum|uuid"`
	Type string `json:"type" binding:"oneof=father mother husband"`
}

/* Create a relation record from a payload struct

   Used by request handlers when communicating with the storage backend.

   Params:
   * sourcePid - id of the relation source person (not included in the payload) */
func (p *itRelationPayload) toRelationRecord(sourcePid string) relationRecord {
	return relationRecord{0, sourcePid, p.Pid, p.Type}
}

/* Relation payload accepted by the createRelation handler

   The 'iit' prefix stands for "id, id, and type" */
type iitRelationPayload struct {
	Pid1 string `json:"pid1" binding:"required,alphanum|uuid"`
	Pid2 string `json:"pid2" binding:"required,alphanum|uuid"`
	Type string `json:"type" binding:"oneof=father mother husband"`
}

/* Create a relation record from a payload struct

   Used by request handlers when communicating with the storage backend. */
func (p *iitRelationPayload) toRelationRecord() relationRecord {
	return relationRecord{0, p.Pid1, p.Pid2, p.Type}
}

/* Structure used to extract relation id from an URI */
type specifyRelationUri struct {
	Rid int64 `uri:"rid" binding:"required"`
}

/* Lower level, shared implementation of the create relation handlers

   The upper level handlers are adapters taking relation record parameters from different
   sources and passing them to this function */
func doCreateRelation(c *gin.Context, relation relationRecord) {
	log.Trace("Entry checkpoint")

	if _, found, err := queryRelationByData(relation.Pid1, relation.Type, relation.Pid2); found {
		log.Infof(
			"A relation (%d) matching given attributes (%s, %s, %s) already exists",
			relation.Id, relation.Pid1, relation.Type, relation.Pid2)
		c.JSON(
			http.StatusBadRequest,
			gin.H{"message": fmt.Sprintf("Relation (%s, %s, %s) already exists",
				relation.Pid1, relation.Type, relation.Pid2)})
		return
	} else if err != nil {
		log.Infof("An error occurred during the relation retrieval attempt (%s)", err)

		c.JSON(http.StatusInternalServerError, gin.H{"message": internalErrorMsg})
		return
	}

	valid, err := validateRelation(relation)

	if !valid {
		log.Infof(
			"The relation (%s, %s, %s) is not valid",
			relation.Pid1, relation.Type, relation.Pid2)

		c.JSON(http.StatusBadRequest,
			gin.H{"message": fmt.Sprintf("Relation (%s, %s, %s) is invalid",
				relation.Pid1, relation.Type, relation.Pid2)})
		return
	} else if err != nil {
		log.Errorf("An error occurred during the relation validation (%s)", err)

		c.JSON(http.StatusInternalServerError, gin.H{"message": internalErrorMsg})
		return
	}

	id, err := getNextRelationId()

	if err != nil {
		log.Infof("An error occurred during the relation id generation (%s)", err)

		c.JSON(http.StatusInternalServerError, gin.H{"message": internalErrorMsg})
		return
	}

	relation.Id = id
	relations[id] = relation

	url := location.Get(c)
	url.Path = fmt.Sprintf("/relations/%d", id)

	c.Header("Location", url.String())
	c.JSON(http.StatusCreated, gin.H{"message": "ok"})

	log.Infof("Created a new relation (%d) record", relation.Id)
}

/* Handle a create relation request

   All the required data is expected in the request payload (iitRelationPayload) */
func createRelation(c *gin.Context) {
	log.Trace("Entry checkpoint")

	var payload iitRelationPayload

	if err := c.ShouldBindJSON(&payload); err != nil {
		log.Infof("New relation data unmarshalling error: %s", err)

		c.JSON(http.StatusBadRequest, gin.H{"message": payloadErrorMsg})
		return
	}

	doCreateRelation(c, payload.toRelationRecord())
}

/* Handle a create relation request

   The source person id is expected to be a part of the URI (specifyPersonUri). The rest of the
   data is expected in the request payload (itRelationPayload) */
func createPersonRelation(c *gin.Context) {
	log.Trace("Entry checkpoint")

	var params specifyPersonUri

	if err := c.ShouldBindUri(&params); err != nil {
		log.Infof("Uri parameters unmarshalling error: %s", err)

		c.JSON(http.StatusBadRequest, gin.H{"message": uriErrorMsg})
		return
	}

	var payload itRelationPayload

	if err := c.ShouldBindJSON(&payload); err != nil {
		log.Infof("New relation data unmarshalling error: %s", err)

		c.JSON(http.StatusBadRequest, gin.H{"message": payloadErrorMsg})
		return
	}

	doCreateRelation(c, payload.toRelationRecord(params.Pid))
}

/* Retrieve a relation

   The relation id is expected to be a part of the URI (specifyRelationUri) */
func retrieveRelation(c *gin.Context) {
	log.Trace("Entry checkpoint")

	var params specifyRelationUri

	if err := c.ShouldBindUri(&params); err != nil {
		log.Infof("Uri parameters unmarshalling error: %s", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": uriErrorMsg})
		return
	}

	relation, found, err := queryRelationById(params.Rid)

	if !found {
		log.Infof("The relation with given id (%d) doesn't exist", params.Rid)
		c.JSON(http.StatusNotFound, gin.H{"message": "Unknown relation id"})
		return
	} else if err != nil {
		log.Errorf("An error occurred during the relation retrieval attempt (%s)", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": internalErrorMsg})
		return
	}

	c.JSON(http.StatusOK, relation.toPayload())

	log.Infof("Found the requested relation record (%d)", params.Rid)
}

/* Retrieve all the relations of given person

   The person id is expected to be a part of the URI (specifyPersonUri). */
func retrievePersonRelations(c *gin.Context) {
	log.Trace("Entry checkpoint")

	var params specifyPersonUri

	if err := c.ShouldBindUri(&params); err != nil {
		log.Infof("Uri parameters unmarshalling error: %s", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": uriErrorMsg})
		return
	}

	var pagination pagePaginationQuery

	if err := c.ShouldBindQuery(&pagination); err != nil {
		log.Infof("Query parameters unmarshalling error: %s", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": queryErrorMsg})
		return
	}

	pagination.applyDefaults()

	if _, found, err := getPerson(params.Pid); !found {
		log.Infof("The person with given id (%s) doesn't exist", params.Pid)
		c.JSON(http.StatusNotFound, gin.H{"message": "Unknown person id"})
		return
	} else if err != nil {
		log.Errorf("An error occurred during the person retrieval attempt (%s)", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": internalErrorMsg})
		return
	}

	relations, err := queryRelationsByPerson(params.Pid, pagination.Page, pagination.Limit)

	if err != nil {
		log.Errorf("An error occurred during relations retrieval attempt (%s)", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": internalErrorMsg})
		return
	}

	c.JSON(http.StatusOK, relations.toPayload())

	log.Infof("Found %d relations for the requested person (%s)", len(relations), params.Pid)
}
