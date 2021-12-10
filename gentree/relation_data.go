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

   The id mustn't exist in the relations table. The generation is performed using cryptographic
   random function. The function fails with an error if it wasn't able to generate a new unique
   number in multiple attempts.

   Returns:
   * new relation record identifier (unique in the scope of the relations map)
   * error (if occurred or when generation failed and nil otherwise) */
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

   Params:
   * pid - the person identifier
   * pageIdx - the zero based index of the results page to be returned (must be greater or equal to
     zero)
   * pageSize - the maximum size of the page to be returned (must be between minPageSize and
     maxPageSize)

   Return:
   * slice of relation records (empty if an error occurred)
   * error (if occurred and nil otherwise) */
func queryRelationsByPerson(pid string, pageIdx int, pageSize int) (relationList, error) {
	log.Debugf("Retrieving all the relations of given person (%s)", pid)

	if err := checkPaginationParams(pageIdx, pageSize); err != nil {
		return []relationRecord{}, err
	}

	// Extract slice of all the values (relation records) of the relations map
	sorted := make(relationList, 0, len(relations))

	for _, r := range relations {
		sorted = append(sorted, r)
	}

	sort.Slice(sorted, func(i, j int) bool { return sorted[i].Id < sorted[j].Id })

	first := minInt(pageIdx*pageSize, len(sorted))
	last := minInt((pageIdx+1)*pageSize, len(sorted))

	return sorted[first:last], nil
}
