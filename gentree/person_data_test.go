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

func TestQueryPeople1Simple(t *testing.T) {
	people = map[string]personRecord{
		"P02": personRecord{"P02", "Anna", "Nowak", gFemale},
		"P04": personRecord{"P04", "Jagoda", "Szewczyk", gFemale},
		"P03": personRecord{"P03", "Antoni", "Michalak", gMale},
		"P05": personRecord{"P05", "Eustachy", "Sobczak", gMale},
		"P06": personRecord{"P06", "Blanka", "Baranowska", gFemale},
		"P01": personRecord{"P01", "Jan", "Kowalski", gMale},
	}

	list, pagResult, err := queryPeople(paginationData{0, 10, 0, 10, 10})

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
}

func TestQueryPeopleEmpty(t *testing.T) {
	people = map[string]personRecord{}

	list, pagResult, err := queryPeople(
		paginationData{
			PageIdx:     0,
			PageSize:    10,
			TotalCnt:    -1, // Should be ignored and overridden by the queryPeople function
			minPageSize: 10,
			maxPageSize: 10})

	assert.Len(t, list, 0)
	assert.Equal(t, pagResult.PageIdx, 0)
	assert.Equal(t, pagResult.PageSize, 10)
	assert.Equal(t, pagResult.TotalCnt, 0)
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
			maxPageSize: 2})

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
			maxPageSize: 2})

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
			maxPageSize: 2})

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
			maxPageSize: 2})

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
			maxPageSize: 2})

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
			maxPageSize: 2})

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
			maxPageSize: 2})

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
			maxPageSize: 2})

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
			maxPageSize: 2})

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
			maxPageSize: 2})

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
			maxPageSize: 2})

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
			maxPageSize: 100})

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
			maxPageSize: 90})

	assert.Len(t, list, 0)
	assert.Empty(t, pagResult.PageIdx)
	assert.Empty(t, pagResult.PageSize)
	assert.Empty(t, pagResult.TotalCnt)
	assert.ErrorIs(
		t, err, AppError{errInvalidArgument, "The page size (200) is out of bounds ([50, 90])"})
}
