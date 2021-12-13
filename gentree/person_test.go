package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPersonPayloadToRecord(t *testing.T) {
	p := personPayload{
		Id:      "123",
		Given:   "Józefa Alina",
		Surname: "Jankowska",
		Gender:  gFemale}

	r := p.toRecord()

	assert.Equal(t, r.Id, "123")
	assert.Equal(t, r.Given, "Józefa Alina")
	assert.Equal(t, r.Surname, "Jankowska")
	assert.Equal(t, r.Gender, gFemale)

	p = personPayload{Id: "XYZ"}

	r = p.toRecord()

	assert.Equal(t, r.Id, "XYZ")
	assert.Empty(t, r.Given)
	assert.Empty(t, r.Surname)
	assert.Equal(t, r.Gender, gUnknown)
}

func TestPersonRecordToPayload(t *testing.T) {
	r := personRecord{
		Id:      "P0001",
		Given:   "Edward Krzysztof",
		Surname: "Kamiński",
		Gender:  gMale}

	p := r.toPayload()

	assert.Equal(t, p.Id, "P0001")
	assert.Equal(t, p.Given, "Edward Krzysztof")
	assert.Equal(t, p.Surname, "Kamiński")
	assert.Equal(t, p.Gender, gMale)
}

func TestPersonListToPayload(t *testing.T) {
	l := personList{
		personRecord{Id: "A1", Given: "Aleksander", Surname: "Cieślak", Gender: gMale},
		personRecord{Id: "A2", Given: "Czesław", Surname: "Gajewski", Gender: gMale},
		personRecord{Id: "A3", Given: "Zofia", Surname: "Krajewska", Gender: gFemale}}

	p := l.toPayload()

	assert.Len(t, p, 3)
	assert.Equal(t, p[0].Id, "A1")
	assert.Equal(t, p[0].Given, "Aleksander")
	assert.Equal(t, p[0].Surname, "Cieślak")
	assert.Equal(t, p[0].Gender, gMale)

	assert.Equal(t, p[1].Id, "A2")
	assert.Equal(t, p[1].Given, "Czesław")
	assert.Equal(t, p[1].Surname, "Gajewski")
	assert.Equal(t, p[1].Gender, gMale)

	assert.Equal(t, p[2].Id, "A3")
	assert.Equal(t, p[2].Given, "Zofia")
	assert.Equal(t, p[2].Surname, "Krajewska")
	assert.Equal(t, p[2].Gender, gFemale)
}
