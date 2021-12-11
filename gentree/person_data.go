package main

import (
	log "github.com/sirupsen/logrus"
	"sort"
)

// Possible gender values
const (
	gMale    = "male"
	gFemale  = "female"
	gUnknown = "unknown"
)

type personRecord struct {
	Id      string `json:"id" binding:"required,alphanum|uuid"`
	Given   string `json:"given_names"`
	Surname string `json:"surname"`
	Gender  string `json:"gender" binding:"isdefault|oneof=male female unknown"`
}

type personList []personRecord

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

/* Query all the person records

   Params:
   * pag - pagination data specifying the range of records to be returned

   Return:
   * slice of person records (empty if an error occurred)
   * updated pagination data (empty if an error occurred; copy of the pag parameter with the total
     record count field updated otherwise)
   * error (if occurred and nil otherwise) */
func queryPeople(pag paginationData) (personList, paginationData, error) {
	log.Debugf("Retrieving all the people")

	if err := pag.validate(); err != nil {
		return []personRecord{}, paginationData{}, err
	}

	// Extract slice of all the values (person records) of the person map
	sorted := make(personList, 0, len(people))

	for _, r := range people {
		sorted = append(sorted, r)
	}

	sort.Slice(sorted, func(i, j int) bool { return sorted[i].Id < sorted[j].Id })

	first := minInt(pag.PageIdx*pag.PageSize, len(sorted))
	last := minInt((pag.PageIdx+1)*pag.PageSize, len(sorted))

	pag.TotalCnt = len(people)

	return sorted[first:last], pag, nil
}
