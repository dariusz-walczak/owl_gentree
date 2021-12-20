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

/* Convert a relation record to payload data

   This function is used by request handlers when responding with data provided by the storage
   backend.

   Returns:
   * relation payload */
func (r *relationRecord) toPayload() relationPayload {
	return relationPayload{r.Id, r.Pid1, r.Pid2, r.Type}
}

/* Convert a list of relation records to payload data

   This function is used by request handlers when responding with data provided by the storage
   backend.

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

   This function is used by request handlers when communicating with the storage backend.

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

   This function is used by request handlers when communicating with the storage backend. */
func (p *iitRelationPayload) toRelationRecord() relationRecord {
	return relationRecord{0, p.Pid1, p.Pid2, p.Type}
}

/* The structure used to extract relation id from a URI */
type specifyRelationUri struct {
	Rid int64 `uri:"rid" binding:"required"`
}

type paginationQuery struct {
	Page  int `form:"page" binding:"min=0"`
	Limit int `form:"limit" binding:"isdefault|min=10,max=100"`
}

func (p *paginationQuery) toPaginationData() paginationData {
	// Apply defaults:
	pageSize := p.Limit

	if pageSize == 0 {
		pageSize = 20
	}

	return paginationData{p.Page, pageSize, 0, minPageSize, maxPageSize}
}

/* Compose an URL allowing retrieval of the given relation

   Params:
   * c - gin context
   * rid - the relation identifier

   Return:
   * URL string */
func makeRetrieveRelationUrl(c *gin.Context, rid int64) string {
	u := location.Get(c)
	u.Path = fmt.Sprintf("/relations/%d", rid)
	return u.String()
}

/* Lower level, shared implementation of the create relation handlers

   The upper-level handlers are adapters taking relation record parameters from different
   sources and passing them to this function */
func doCreateRelation(c *gin.Context, relation relationRecord) {
	log.Trace("Entry checkpoint")

	if existing, found, err := queryRelationByData(relation.Pid1, relation.Type, relation.Pid2); found {
		log.Infof(
			"A relation (%d) matching given attributes (%s, %s, %s) already exists",
			existing.Id, existing.Pid1, existing.Type, existing.Pid2)
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"message": fmt.Sprintf("Relation (%s, %s, %s) already exists",
					existing.Pid1, existing.Type, existing.Pid2),
				"location": makeRetrieveRelationUrl(c, existing.Id),
			})

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

	c.Header("Location", makeRetrieveRelationUrl(c, id))
	c.JSON(http.StatusCreated, gin.H{"message": "Relation created", "relation_id": id})

	log.Infof("Created a new relation (%d) record", relation.Id)
}

/* Handle a create relation request

   The function will retrieve all the input data from the request payload (iitRelationPayload) */
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

   The function will retrieve the source person id from the request URI (specifyPersonUri), and the
   rest of the data from the request payload (itRelationPayload) */
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

/* Delete a relation

   The function will extract the relation id from the request URI (specifyRelationUri) */
func deleteRelation(c *gin.Context) {
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

	delete(relations, params.Rid)
	c.JSON(http.StatusOK, gin.H{"message": "Relation deleted"})

	log.Infof(
		"Deleted the requested relation (%d) record:  %s, %s, %s",
		relation.Id, relation.Pid1, relation.Type, relation.Pid2)
}

/* Retrieve a relation

   The function will extract the relation id from the request URI (specifyRelationUri) */
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

/* Retrieve all the existing relations */
func retrieveRelations(c *gin.Context) {
	log.Trace("Entry checkpoint")

	var pagQuery paginationQuery

	if err := c.ShouldBindQuery(&pagQuery); err != nil {
		log.Infof("Query parameters unmarshalling error: %s", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": queryErrorMsg})
		return
	}

	relations, pagData, err := queryRelationsByPerson("", pagQuery.toPaginationData())

	if err != nil {
		log.Errorf("An error occurred during relations retrieval attempt (%s)", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": internalErrorMsg})
		return
	}

	reqUrl := location.Get(c)
	reqUrl.Path = "/relations"

	c.JSON(http.StatusOK, gin.H{
		"pagination": pagData.getJson(*reqUrl),
		"records":    relations.toPayload(),
	})

	log.Infof("Found %d relations", len(relations))
}

/* Retrieve all the relations of the given person

   The function will extract the person id from the request URI (specifyPersonUri) */
func retrievePersonRelations(c *gin.Context) {
	log.Trace("Entry checkpoint")

	var params specifyPersonUri

	if err := c.ShouldBindUri(&params); err != nil {
		log.Infof("Uri parameters unmarshalling error: %s", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": uriErrorMsg})
		return
	}

	var pagQuery paginationQuery

	if err := c.ShouldBindQuery(&pagQuery); err != nil {
		log.Infof("Query parameters unmarshalling error: %s", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": queryErrorMsg})
		return
	}

	if _, found, err := getPerson(params.Pid); !found {
		log.Infof("The person with given id (%s) doesn't exist", params.Pid)
		c.JSON(http.StatusNotFound, gin.H{"message": "Unknown person id"})
		return
	} else if err != nil {
		log.Errorf("An error occurred during the person retrieval attempt (%s)", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": internalErrorMsg})
		return
	}

	relations, pagData, err := queryRelationsByPerson(params.Pid, pagQuery.toPaginationData())

	if err != nil {
		log.Errorf("An error occurred during relations retrieval attempt (%s)", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": internalErrorMsg})
		return
	}

	reqUrl := location.Get(c)
	reqUrl.Path = fmt.Sprintf("/people/%s/relations", params.Pid)

	c.JSON(http.StatusOK, gin.H{
		"pagination": pagData.getJson(*reqUrl),
		"records":    relations.toPayload(),
	})

	log.Infof("Found %d relations for the requested person (%s)", len(relations), params.Pid)
}
