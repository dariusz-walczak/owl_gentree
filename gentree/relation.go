package main

import (
	rand "crypto/rand"
	"fmt"
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
	Id   int64  `json:"id" binding:"required,hexadecimal"`
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
   * Relation record structure (uninitialized if not found or when an error occurred)
   * Success flag (true if one and only one record was found, and false otherwise)
   * Error (if occurred) */
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

func createRelation(c *gin.Context) {
	log.Trace("Entry checkpoint")

	var payload relationPayload

	if err := c.ShouldBindJSON(&payload); err != nil {
		log.Infof("New relation data unmarshalling error: %s", err)

		c.JSON(http.StatusBadRequest, gin.H{"message": payloadErrorMsg})
		return
	}

	relation := payload.toRelationRecord()

	//validateRelation : gender vs type

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
		log.Infof("An error occurred during the relations retrieval attempt (%s)", err)

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

	c.JSON(http.StatusCreated, gin.H{"message": "ok"})

	log.Infof("Created a new relation (%d) record", relation.Id)
}
