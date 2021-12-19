package main

/* This file defines storage data structures and functions interacting with the data storage
   backend */

import (
	rand "crypto/rand"
	"fmt"
	log "github.com/sirupsen/logrus"
	"math"
	"math/big"
	"sort"
)

const (
	relFather  = "father"
	relMother  = "mother"
	relHusband = "husband"
)

type relationRecord struct {
	Id   int64
	Pid1 string
	Pid2 string
	Type string
}

type relationList []relationRecord

var relations = map[int64]relationRecord{}

/* Generate a new, unique relation id

   The function generates a random, unique identifier using a cryptographic function. It always
   uses the entire identifier pool, so the chance for a conflict increases with the increasing
   number of existing records. If multiple attempts of unique identifier generation fail, the
   function will fail, too.

   Returns:
   * new relation record identifier (unique in the scope of the relations table)
   * error (if occurred or when the generation failed and nil otherwise) */
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

/* Query a relation record by relation id

   Returns:
   * relation record (uninitialized if not found or when an error occurred)
   * success flag (true if the relation was found and false otherwise)
   * error (if occurred and nil otherwise) */
func queryRelationById(id int64) (relationRecord, bool, error) {
	log.Debugf("Retrieving relation record by id (%d)", id)

	relation, found := relations[id]

	if !found {
		log.Debugf("Relation record (%d) not found", id)

		return relation, false, nil
	}

	return relation, true, nil
}

/* Query all the relation records associated with the given person

   All the existing records will be returned if the person identifier is an empty string

   Params:
   * pid - the person identifier (ignored if it is an empty string)
   * pag - pagination data specifying the range of records to be returned

   Return:
   * slice of relation records (empty if an error occurred)
   * updated pagination data (empty if an error occurred; copy of the pag parameter with the total
     record count field updated otherwise)
   * error (if occurred and nil otherwise) */
func queryRelationsByPerson(pid string, pag paginationData) (relationList, paginationData, error) {
	log.Debugf("Retrieving all the relations of given person (%s)", pid)

	if err := pag.validate(); err != nil {
		return []relationRecord{}, paginationData{}, err
	}

	// Extract slice of all the values (relation records) of the relations map
	sorted := make(relationList, 0, len(relations))

	for _, r := range relations {
		if (r.Pid1 == pid) || (r.Pid2 == pid) || (pid == "") {
			sorted = append(sorted, r)
		}
	}

	sort.Slice(sorted, func(i, j int) bool { return sorted[i].Id < sorted[j].Id })

	first := minInt(pag.PageIdx*pag.PageSize, len(sorted))
	last := minInt((pag.PageIdx+1)*pag.PageSize, len(sorted))

	pag.TotalCnt = len(sorted)

	return sorted[first:last], pag, nil
}

func queryRelationsByData(pid1 string, typ string, pid2 string) ([]relationRecord, error) {
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
func queryRelationByData(pid1 string, typ string, pid2 string) (relationRecord, bool, error) {
	relations, err := queryRelationsByData(pid1, typ, pid2)

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
   *** The first person in the husband relation must be a male, and the second one must be a female
   ** Relations of the same type already exist for some of the people (case of multiple fathers or
      mothers)
   *** The second person in the father relation mustn't have any other father relation in which
       they are the target (the second) person
   *** The second person in the mother relation mustn't have any other mother relation in which
       they are the target (the second) person
   * The current approach to the relation types is minimal on purpose:
   ** The first implementation is more straightforward in terms of data errors detection with such
      a simple model
   ** The need to add less common relations (same-sex partnerships, child adoption, etc.) is
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
		other, found, err := queryRelationByData("", r.Type, r.Pid2)

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
