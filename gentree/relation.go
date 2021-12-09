package main

import (
	rand "crypto/rand"
	"fmt"
	"github.com/gin-contrib/location"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"math"
	"math/big"
	"net/http"
)

const (
	relFather  = "father"
	relMother  = "mother"
	relHusband = "husband"
)

type relationRecord struct {
	Id   int64  `json:"id" binding:"required"`
	Pid1 string `json:"pid1" binding:"required,alphanum|uuid"`
	Pid2 string `json:"pid2" binding:"required,alphanum|uuid"`
	Type string `json:"type" binding:"oneof=father mother husband"`
}

// Payload for the createPersonRelation request (POST /people/:pid/relations)
// The first person id is taken from the uri and the second is taken from the payload
type personRelationPayload struct {
	Pid  string `json:"pid" binding:"required,alphanum|uuid"`
	Type string `json:"type" binding:"oneof=father mother husband"`
}

func (r *personRelationPayload) toRelationRecord(targetPid string) relationRecord {
	return relationRecord{0, targetPid, r.Pid, r.Type}
}

// Payload for the createRelation request (POST /relations)
type relationPayload struct {
	Pid1 string `json:"pid1" binding:"required,alphanum|uuid"`
	Pid2 string `json:"pid2" binding:"required,alphanum|uuid"`
	Type string `json:"type" binding:"oneof=father mother husband"`
}

func (r *relationPayload) toRelationRecord() relationRecord {
	return relationRecord{0, r.Pid1, r.Pid2, r.Type}
}

var relations = map[int64]relationRecord{}

func getNextRelationId() (int64, error) {
	const maxAttempts = 5

	for i := 0; i < maxAttempts; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
		if err != nil {
			return 0, err
		}

		if _, found := relations[num.Int64()]; !found {
			return num.Int64(), nil
		}
	}

	msg := fmt.Sprintf("Failed to generate relation id %d attempts", maxAttempts)

	log.Warn(msg)

	return 0, AppError{errIdGenerationFailed, msg}
}

func getRelation(id int64) (relationRecord, bool, error) {
	log.Debugf("Retrieving relation record by id (%d)", id)

	relation, found := relations[id]

	if !found {
		log.Debugf("Relation record (%d) not found", id)

		return relation, false, nil
	}

	return relation, true, nil
}

/* Get all the relation records associated with the given person

   Return:
   * slice of relation records (empty if an error occurred)
   * error (if occurred and nil otherwise) */
func getRelations(pid string) ([]relationRecord, error) {
	log.Debugf("Retrieving all the relations of given person (%s)", pid)

	var result []relationRecord

	for _, r := range relations {
		if (r.Pid1 == pid) || (r.Pid2 == pid) {
			result = append(result, r)
		}
	}

	return result, nil
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

func createRelation(c *gin.Context) {
	log.Trace("Entry checkpoint")

	var payload relationPayload

	if err := c.ShouldBindJSON(&payload); err != nil {
		log.Infof("New relation data unmarshalling error: %s", err)

		c.JSON(http.StatusBadRequest, gin.H{"message": payloadErrorMsg})
		return
	}

	relation := payload.toRelationRecord()

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

	relation, found, err := getRelation(params.Rid)

	if !found {
		log.Infof("The relation with given id (%d) doesn't exist", params.Rid)
		c.JSON(http.StatusNotFound, gin.H{"message": "Unknown relation id"})
		return
	} else if err != nil {
		log.Errorf("An error occurred during the relation retrieval attempt (%s)", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": internalErrorMsg})
		return
	}

	c.JSON(http.StatusOK, relation)

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

	if _, found, err := getPerson(params.Pid); !found {
		log.Infof("The person with given id (%s) doesn't exist", params.Pid)
		c.JSON(http.StatusNotFound, gin.H{"message": "Unknown person id"})
		return
	} else if err != nil {
		log.Errorf("An error occurred during the person retrieval attempt (%s)", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": internalErrorMsg})
		return
	}

	relations, err := getRelations(params.Pid)

	if err != nil {
		log.Errorf("An error occurred during relations retrieval attempt (%s)", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": internalErrorMsg})
		return
	}

	c.JSON(http.StatusOK, relations)

	log.Infof("Found %d relations for the requested person (%s)", len(relations), params.Pid)
}
