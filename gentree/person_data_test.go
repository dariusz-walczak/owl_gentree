package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetPerson(t *testing.T) {
	people = map[string]personRecord{
		"P1": personRecord{"P1", "Jan", "Kowalski", gMale},
		"P2": personRecord{"P2", "Anna", "Nowak", gFemale},
		"P3": personRecord{"P3", "", "", gUnknown}}

	person, found, err := getPerson("P2")

	assert.Equal(t, person.Id, "P2")
	assert.Equal(t, person.Given, "Anna")
	assert.Equal(t, person.Surname, "Nowak")
	assert.Equal(t, person.Gender, gFemale)
	assert.True(t, found)
	assert.Nil(t, err)

	person, found, err = getPerson("P4")

	assert.Empty(t, person.Id)
	assert.Empty(t, person.Given)
	assert.Empty(t, person.Surname)
	assert.Empty(t, person.Gender)
	assert.False(t, found)
	assert.Nil(t, err)
}

/* Test the queryPeople function with a non-empty short people list
 *
 * 1. All the records are returned when the person filter is at the defaults
 * 2. Only the requested records are returned then the person ids filter is used */
func TestQueryPeople1Simple(t *testing.T) {
	people = map[string]personRecord{
		"P02": personRecord{"P02", "Anna", "Nowak", gFemale},
		"P04": personRecord{"P04", "Jagoda", "Szewczyk", gFemale},
		"P03": personRecord{"P03", "Antoni", "Michalak", gMale},
		"P05": personRecord{"P05", "Eustachy", "Sobczak", gMale},
		"P06": personRecord{"P06", "Blanka", "Baranowska", gFemale},
		"P01": personRecord{"P01", "Jan", "Kowalski", gMale},
	}

	// Case 1: All the records returned without filtering

	list, pagResult, err := queryPeople(paginationData{0, 10, 0, 10, 10}, personFilter{})

	assert.Len(t, list, 6)
	// The result table should be sorted by the person id field:
	assert.Equal(t, list[0].Id, "P01")
	assert.Equal(t, list[1].Id, "P02")
	assert.Equal(t, list[2].Id, "P03")
	assert.Equal(t, list[3].Id, "P04")
	assert.Equal(t, list[4].Id, "P05")
	assert.Equal(t, list[5].Id, "P06")
	assert.Equal(t, pagResult.PageIdx, 0)
	assert.Equal(t, pagResult.PageSize, 10)
	assert.Equal(t, pagResult.TotalCnt, 6)
	assert.Nil(t, err)

	// Case 2: Only the requested records returned with person ids filter

	list, pagResult, err = queryPeople(
		paginationData{
			PageIdx:     0,
			PageSize:    10,
			TotalCnt:    -1, // Should be ignored and overridden by the queryPeople function
			minPageSize: 10,
			maxPageSize: 10},
		personFilter{
			personIdsFilter{[]string{"P03", "P01", "P04"}, true}})

	assert.Len(t, list, 3)
	// The result table should be sorted by the person id field:
	assert.Equal(t, list[0].Id, "P01")
	assert.Equal(t, list[1].Id, "P03")
	assert.Equal(t, list[2].Id, "P04")
	assert.Equal(t, pagResult.PageIdx, 0)
	assert.Equal(t, pagResult.PageSize, 10)
	assert.Equal(t, pagResult.TotalCnt, 3)
	assert.Nil(t, err)
}

/* Test the queryPeople function when the people list is empty
 *
 * 1. Check the no filter scenario
 * 2. Check the person ids filter scenario */
func TestQueryPeopleEmpty(t *testing.T) {
	people = map[string]personRecord{}

	list, pagResult, err := queryPeople(
		paginationData{
			PageIdx:     0,
			PageSize:    10,
			TotalCnt:    -1, // Should be ignored and overridden by the queryPeople function
			minPageSize: 10,
			maxPageSize: 10},
		personFilter{})

	assert.Len(t, list, 0)
	assert.Equal(t, pagResult.PageIdx, 0)
	assert.Equal(t, pagResult.PageSize, 10)
	assert.Equal(t, pagResult.TotalCnt, 0)
	assert.Nil(t, err)

	list, pagResult, err = queryPeople(
		paginationData{
			PageIdx:     0,
			PageSize:    10,
			TotalCnt:    -1, // Should be ignored and overridden by the queryPeople function
			minPageSize: 10,
			maxPageSize: 10},
		personFilter{
			personIdsFilter{[]string{"missing1", "missing2", "IRRELEVANT630"}, true}})

	assert.Len(t, list, 0)
	assert.Equal(t, pagResult.PageIdx, 0)
	assert.Equal(t, pagResult.PageSize, 10)
	assert.Equal(t, pagResult.TotalCnt, 0)
	assert.Nil(t, err)
}

/* Test the queryPeople function with a non-empty people list and person ids filter border cases
 *
 * 1. Check the empty ids set (with the enabled flag set to true)
 * 2. Check the unknown id case
 * 3. Check a single, known id case
 * 4. Check the multiple, some unknown, ids case
 * 5. Check the all existing ids case
 * 6. Check the all existing ids case with some extra unknown ids */
func TestQueryPeoplePidsFilter(t *testing.T) {
	people = map[string]personRecord{
		"y 002": personRecord{"y 002", "Zuzanna", "Dąbrowska", gFemale},
		"y 001": personRecord{"y 001", "Bogumiła", "Bąk", gFemale},
		"y 003": personRecord{"y 003", "Edward", "Szymczak", gMale},
		"y 004": personRecord{"y 004", "Jerzy", "Sokołowski", gMale},
		"y 005": personRecord{"y 005", "Lila", "Gajewska", gFemale},
	}

	// Case 1: Empty ids set

	list, pagResult, err := queryPeople(
		paginationData{
			PageIdx:     0,
			PageSize:    10,
			TotalCnt:    -1, // Should be ignored and overridden by the queryPeople function
			minPageSize: 10,
			maxPageSize: 10},
		personFilter{
			personIdsFilter{[]string{}, true}})

	assert.Len(t, list, 0)
	assert.Equal(t, pagResult.PageIdx, 0)
	assert.Equal(t, pagResult.PageSize, 10)
	assert.Equal(t, pagResult.TotalCnt, 0)
	assert.Nil(t, err)

	// Case 2: Unknown id

	list, pagResult, err = queryPeople(
		paginationData{
			PageIdx:     0,
			PageSize:    10,
			TotalCnt:    -1, // Should be ignored and overridden by the queryPeople function
			minPageSize: 10,
			maxPageSize: 10},
		personFilter{
			personIdsFilter{[]string{"unknown", "XYZ"}, true}})

	assert.Len(t, list, 0)
	assert.Equal(t, pagResult.PageIdx, 0)
	assert.Equal(t, pagResult.PageSize, 10)
	assert.Equal(t, pagResult.TotalCnt, 0)
	assert.Nil(t, err)

	// Case 3: Known id

	list, pagResult, err = queryPeople(
		paginationData{
			PageIdx:     0,
			PageSize:    10,
			TotalCnt:    -1, // Should be ignored and overridden by the queryPeople function
			minPageSize: 10,
			maxPageSize: 10},
		personFilter{
			personIdsFilter{[]string{"y 003"}, true}})

	assert.Len(t, list, 1)
	assert.Equal(t, list[0].Id, "y 003")
	assert.Equal(t, pagResult.PageIdx, 0)
	assert.Equal(t, pagResult.PageSize, 10)
	assert.Equal(t, pagResult.TotalCnt, 1)
	assert.Nil(t, err)

	// Case 4: Multiple ids, some unknown

	list, pagResult, err = queryPeople(
		paginationData{
			PageIdx:     0,
			PageSize:    10,
			TotalCnt:    -1, // Should be ignored and overridden by the queryPeople function
			minPageSize: 10,
			maxPageSize: 10},
		personFilter{
			personIdsFilter{[]string{"y 002", "P01", "y 001", "y 005"}, true}})

	assert.Len(t, list, 3)
	assert.Equal(t, list[0].Id, "y 001")
	assert.Equal(t, list[1].Id, "y 002")
	assert.Equal(t, list[2].Id, "y 005")
	assert.Equal(t, pagResult.PageIdx, 0)
	assert.Equal(t, pagResult.PageSize, 10)
	assert.Equal(t, pagResult.TotalCnt, 3)
	assert.Nil(t, err)

	// Case 5: All existing ids are in the filter

	list, pagResult, err = queryPeople(
		paginationData{
			PageIdx:     0,
			PageSize:    10,
			TotalCnt:    -1, // Should be ignored and overridden by the queryPeople function
			minPageSize: 10,
			maxPageSize: 10},
		personFilter{
			personIdsFilter{[]string{"y 002", "y 004", "y 001", "y 003", "y 005"}, true}})

	assert.Len(t, list, 5)
	assert.Equal(t, list[0].Id, "y 001")
	assert.Equal(t, list[1].Id, "y 002")
	assert.Equal(t, list[2].Id, "y 003")
	assert.Equal(t, list[3].Id, "y 004")
	assert.Equal(t, list[4].Id, "y 005")
	assert.Equal(t, pagResult.PageIdx, 0)
	assert.Equal(t, pagResult.PageSize, 10)
	assert.Equal(t, pagResult.TotalCnt, 5)
	assert.Nil(t, err)

	// Case 6: All existing ids and some unknown

	list, pagResult, err = queryPeople(
		paginationData{
			PageIdx:     0,
			PageSize:    10,
			TotalCnt:    -1, // Should be ignored and overridden by the queryPeople function
			minPageSize: 10,
			maxPageSize: 10},
		personFilter{
			personIdsFilter{
				[]string{"P001", "y 004", "y 001", "y 002", "P999", "y 003", "y 005"}, true}})

	assert.Len(t, list, 5)
	assert.Equal(t, list[0].Id, "y 001")
	assert.Equal(t, list[1].Id, "y 002")
	assert.Equal(t, list[2].Id, "y 003")
	assert.Equal(t, list[3].Id, "y 004")
	assert.Equal(t, list[4].Id, "y 005")
	assert.Equal(t, pagResult.PageIdx, 0)
	assert.Equal(t, pagResult.PageSize, 10)
	assert.Equal(t, pagResult.TotalCnt, 5)
	assert.Nil(t, err)
}

func TestQueryPeoplePaging(t *testing.T) {
	people = map[string]personRecord{
		"P01": personRecord{"P01", "Anna", "Kowalska", gFemale},
	}

	list, pagResult, err := queryPeople(
		paginationData{
			PageIdx:     0,
			PageSize:    2,
			TotalCnt:    -1, // Should be ignored and overridden by the queryPeople function
			minPageSize: 2,
			maxPageSize: 2},
		personFilter{})

	assert.Len(t, list, 1)
	assert.Equal(t, list[0].Id, "P01")
	assert.Equal(t, pagResult.PageIdx, 0)
	assert.Equal(t, pagResult.PageSize, 2)
	assert.Equal(t, pagResult.TotalCnt, 1)
	assert.Nil(t, err)

	list, pagResult, err = queryPeople(
		paginationData{
			PageIdx:     1,
			PageSize:    2,
			TotalCnt:    -1, // Should be ignored and overridden by the queryPeople function
			minPageSize: 2,
			maxPageSize: 2},
		personFilter{})

	assert.Len(t, list, 0)
	assert.Equal(t, pagResult.PageIdx, 1)
	assert.Equal(t, pagResult.PageSize, 2)
	assert.Equal(t, pagResult.TotalCnt, 1)
	assert.Nil(t, err)

	people["P03"] = personRecord{"P03", "Żaneta", "Rutkowska", gFemale}

	list, pagResult, err = queryPeople(
		paginationData{
			PageIdx:     0,
			PageSize:    2,
			TotalCnt:    -1, // Should be ignored and overridden by the queryPeople function
			minPageSize: 2,
			maxPageSize: 2},
		personFilter{})

	assert.Len(t, list, 2)
	assert.Equal(t, list[0].Id, "P01")
	assert.Equal(t, list[1].Id, "P03")
	assert.Equal(t, pagResult.PageIdx, 0)
	assert.Equal(t, pagResult.PageSize, 2)
	assert.Equal(t, pagResult.TotalCnt, 2)
	assert.Nil(t, err)

	list, pagResult, err = queryPeople(
		paginationData{
			PageIdx:     1,
			PageSize:    2,
			TotalCnt:    -1, // Should be ignored and overridden by the queryPeople function
			minPageSize: 2,
			maxPageSize: 2},
		personFilter{})

	assert.Len(t, list, 0)
	assert.Equal(t, pagResult.PageIdx, 1)
	assert.Equal(t, pagResult.PageSize, 2)
	assert.Equal(t, pagResult.TotalCnt, 2)
	assert.Nil(t, err)

	people["P04"] = personRecord{"P04", "Anatol", "Chmielewski", gMale}

	list, pagResult, err = queryPeople(
		paginationData{
			PageIdx:     0,
			PageSize:    2,
			TotalCnt:    -1, // Should be ignored and overridden by the queryPeople function
			minPageSize: 2,
			maxPageSize: 2},
		personFilter{})

	assert.Len(t, list, 2)
	assert.Equal(t, list[0].Id, "P01")
	assert.Equal(t, list[1].Id, "P03")
	assert.Equal(t, pagResult.PageIdx, 0)
	assert.Equal(t, pagResult.PageSize, 2)
	assert.Equal(t, pagResult.TotalCnt, 3)
	assert.Nil(t, err)

	list, pagResult, err = queryPeople(
		paginationData{
			PageIdx:     1,
			PageSize:    2,
			TotalCnt:    -1, // Should be ignored and overridden by the queryPeople function
			minPageSize: 2,
			maxPageSize: 2},
		personFilter{})

	assert.Len(t, list, 1)
	assert.Equal(t, list[0].Id, "P04")
	assert.Equal(t, pagResult.PageIdx, 1)
	assert.Equal(t, pagResult.PageSize, 2)
	assert.Equal(t, pagResult.TotalCnt, 3)
	assert.Nil(t, err)

	list, pagResult, err = queryPeople(
		paginationData{
			PageIdx:     2,
			PageSize:    2,
			TotalCnt:    -1, // Should be ignored and overridden by the queryPeople function
			minPageSize: 2,
			maxPageSize: 2},
		personFilter{})

	assert.Len(t, list, 0)
	assert.Equal(t, pagResult.PageIdx, 2)
	assert.Equal(t, pagResult.PageSize, 2)
	assert.Equal(t, pagResult.TotalCnt, 3)
	assert.Nil(t, err)

	// Note that the 'P02' identifier puts this record on the first page
	people["P02"] = personRecord{"P02", "Michał", "Jasiński", gMale}

	list, pagResult, err = queryPeople(
		paginationData{
			PageIdx:     0,
			PageSize:    2,
			TotalCnt:    -1, // Should be ignored and overridden by the queryPeople function
			minPageSize: 2,
			maxPageSize: 2},
		personFilter{})

	assert.Len(t, list, 2)
	assert.Equal(t, list[0].Id, "P01")
	assert.Equal(t, list[1].Id, "P02")
	assert.Equal(t, pagResult.PageIdx, 0)
	assert.Equal(t, pagResult.PageSize, 2)
	assert.Equal(t, pagResult.TotalCnt, 4)
	assert.Nil(t, err)

	list, pagResult, err = queryPeople(
		paginationData{
			PageIdx:     1,
			PageSize:    2,
			TotalCnt:    -1, // Should be ignored and overridden by the queryPeople function
			minPageSize: 2,
			maxPageSize: 2},
		personFilter{})

	assert.Len(t, list, 2)
	assert.Equal(t, list[0].Id, "P03")
	assert.Equal(t, list[1].Id, "P04")
	assert.Equal(t, pagResult.PageIdx, 1)
	assert.Equal(t, pagResult.PageSize, 2)
	assert.Equal(t, pagResult.TotalCnt, 4)
	assert.Nil(t, err)

	list, pagResult, err = queryPeople(
		paginationData{
			PageIdx:     2,
			PageSize:    2,
			TotalCnt:    -1, // Should be ignored and overridden by the queryPeople function
			minPageSize: 2,
			maxPageSize: 2},
		personFilter{})

	assert.Len(t, list, 0)
	assert.Equal(t, pagResult.PageIdx, 2)
	assert.Equal(t, pagResult.PageSize, 2)
	assert.Equal(t, pagResult.TotalCnt, 4)
	assert.Nil(t, err)
}

func TestQueryPeopleValidation(t *testing.T) {
	people = map[string]personRecord{
		"P01": personRecord{"P01", "Anna", "Kowalska", gFemale},
		"P02": personRecord{"P02", "Błażej", "Czerwiński", gMale},
		"P03": personRecord{"P03", "Bianka", "Wysocka", gFemale},
	}

	// Page index smaller than 0:
	list, pagResult, err := queryPeople(
		paginationData{
			PageIdx:     -1,
			PageSize:    2,
			TotalCnt:    -1,
			minPageSize: 2,
			maxPageSize: 2},
		personFilter{})

	assert.Len(t, list, 0)
	assert.Empty(t, pagResult.PageIdx)
	assert.Empty(t, pagResult.PageSize)
	assert.Empty(t, pagResult.TotalCnt)
	assert.ErrorIs(t, err, AppError{errInvalidArgument, "The page index is negative (-1)"})

	// Page size smaller than the minimum
	list, pagResult, err = queryPeople(
		paginationData{
			PageIdx:     0,
			PageSize:    2,
			TotalCnt:    -1,
			minPageSize: 10,
			maxPageSize: 100},
		personFilter{})

	assert.Len(t, list, 0)
	assert.Empty(t, pagResult.PageIdx)
	assert.Empty(t, pagResult.PageSize)
	assert.Empty(t, pagResult.TotalCnt)
	assert.ErrorIs(
		t, err, AppError{errInvalidArgument, "The page size (2) is out of bounds ([10, 100])"})

	// Page size greater than the maximum
	list, pagResult, err = queryPeople(
		paginationData{
			PageIdx:     0,
			PageSize:    200,
			TotalCnt:    -1,
			minPageSize: 50,
			maxPageSize: 90},
		personFilter{})

	assert.Len(t, list, 0)
	assert.Empty(t, pagResult.PageIdx)
	assert.Empty(t, pagResult.PageSize)
	assert.Empty(t, pagResult.TotalCnt)
	assert.ErrorIs(
		t, err, AppError{errInvalidArgument, "The page size (200) is out of bounds ([50, 90])"})
}
