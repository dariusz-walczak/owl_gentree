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

func TestQueryPeople1(t *testing.T) {
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
