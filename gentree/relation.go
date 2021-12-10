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

func findRelations(pid1 string, typ string, pid2 string) ([]relationRecord, error) {
	log.Debugf("Looking for matching relations (%s, %s, %s)", pid1, typ, pid2)

	var result []relationRecord

	for _, r := range relations {
		if (pid1 == "" || pid1 == r.Pid1) && (typ == "" || typ == r.Type) &&
			(pid2 == "" || pid2 == r.Pid2) {
			result = append(result, r)
		}
	}

	log.Debugf("Found %d matching relations", len(result))

	return result, nil
}

/* Find the relation record matching given attributes
   In the case of multiple matching relations, return a duplicate error.
   Returns:
   * relation record structure (uninitialized if not found or when an error occurred)
   * success flag (true if one and only one record was found, and false otherwise)
   * error (if occurred and nil otherwise) */
func findRelation(pid1 string, typ string, pid2 string) (relationRecord, bool, error) {
	relations, err := findRelations(pid1, typ, pid2)

	if err != nil {
		return relationRecord{}, false, err
	} else if len(relations) > 1 {
		msg := fmt.Sprintf(
			"%d duplicated relation records found: %s, %s, %s",
			len(relations), pid1, typ, pid2)
		return relationRecord{}, false, AppError{errDuplicateFound, msg}
	} else if len(relations) > 0 {
		return relations[0], true, nil
	}

	return relationRecord{}, false, nil
}

/* Check if the relation record is valid considering people records and other existing relation
   records.

   Return:
   * Outcome flag (true if the relation is valid and false otherwise)
   * Error (if occurred and nil otherwise)

   Design Assumptions:
   * The relation is considered invalid when:
   ** At least one of the related people doesn't exist
   ** The people gender is inconsistent with the gender implied by the relation type:
   *** The first person in the father relation must be a male
   *** The first person in the mother relation must be a female
   *** The first person in the husband relation must be a male and the second one must be a female
   ** Relations of the same type already exist for some of the people (multiple fathers/mothers case)
   *** The second person in the father relation mustn't have any other father relation in which they
       are the target (the second) person
   *** The second person in the mother relation mustn't have any other mother relation in which they
       are the target (the second) person
   * The current approach to the relation types is absolutely minimal on purpose:
   ** The first implementation is easier in terms of data errors detection with such a simple model
   ** The need to add less common relations (same sex partnerships, child adoption, etc.) is
      recognized but planned as an extension when the basic functionality works. */
func validateRelation(r relationRecord) (bool, error) {
	p1, found, err := getPerson(r.Pid1)

	if !found {
		log.Infof(
			"The person (%s) referenced by the relation (%s, %s, %s) doesn't exist",
			r.Pid1, r.Pid1, r.Type, r.Pid2)

		return false, nil
	} else if err != nil {
		log.Tracef("An error occurred during the person retrieval attempt (%s)", err)

		return false, err
	}

	if (p1.Gender != gMale) && (r.Type == relFather || r.Type == relHusband) {
		log.Infof(
			"Unexpected person (%s) gender (%s): for the '%s' relation, '%s' is expected",
			p1.Id, p1.Gender, r.Type, gMale)

		return false, nil
	} else if (p1.Gender != gFemale) && (r.Type == relMother) {
		log.Infof(
			"Unexpected person (%s) gender (%s): for the '%s' relation, '%s' is expected",
			p1.Id, p1.Gender, r.Type, gFemale)

		return false, nil
	}

	p2, found, err := getPerson(r.Pid2)

	if !found {
		log.Debugf(
			"The person (%s) referenced by the relation (%s, %s, %s) doesn't exist",
			r.Pid2, r.Pid1, r.Type, r.Pid2)

		return false, nil
	} else if err != nil {
		log.Tracef("An error occurred during the person retrieval attempt (%s)", err)

		return false, err
	}

	if (p2.Gender != gFemale) && (r.Type == relHusband) {
		log.Infof(
			"Unexpected person (%s) gender (%s): for the '%s' relation, '%s' is expected",
			p2.Id, p2.Gender, r.Type, gFemale)

		return false, nil
	}

	// Check the multiple fathers/mothers case:

	if (r.Type == relFather) || (r.Type == relMother) {
		other, found, err := findRelation("", r.Type, r.Pid2)

		if found {
			log.Infof(
				"Found another (%d) %s relation for the target person (%s)",
				other.Id, r.Type, r.Pid2)

			return false, nil
		} else if err != nil {
			log.Tracef("An error occurred during the relation retrieval attempt (%s)", err)

			return false, err
		}
	}

	return true, nil
}

/* The lower level implementation of the create relation request.

   The upper level handlers are adapters taking relation record parameters from different
   sources */
func doCreateRelation(c *gin.Context, relation relationRecord) {
	log.Trace("Entry checkpoint")

	if _, found, err := findRelation(relation.Pid1, relation.Type, relation.Pid2); found {
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

type specifyRelationUri struct {
	Rid int64 `uri:"rid" binding:"required"`
}

// Retrieve the relation specified through the relation id (provided in the uri)
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

// Retrieve all the relations of the given person
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
