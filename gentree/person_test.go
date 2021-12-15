package main

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testPaginationJson struct {
	PrevUrl string `json:"prev_url"`
	NextUrl string `json:"next_url"`
}

type testPersonJson struct {
	Id      string `json:"id"`
	Given   string `json:"given_names"`
	Surname string `json:"surname"`
	Gender  string `json:"gender"`
}

type testErrorJson struct {
	Message string `json:"message"`
}

/* Test the person payload to the person record conversion function */
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

/* Test the person record to the person payload conversion function */
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

/* Test the person record list to the person payload list conversion function */
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

func TestCreatePersonRequestSuccess(t *testing.T) {
	router := setupRouter()

	people = map[string]personRecord{}

	person1 := testPersonJson{
		Id: "1",
		Given: "Dorota Justyna",
		Surname: "Zawadzka",
		Gender: gFemale}

	json1, err := json.Marshal(person1)
	require.Nil(t, err)

	res := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/people", bytes.NewBuffer(json1))
	router.ServeHTTP(res, req)

	assert.Equal(t, res.Code, http.StatusCreated)

	require.True(t, json.Valid(res.Body.Bytes()))
	responseData := testErrorJson{}
	err = json.Unmarshal(res.Body.Bytes(), &responseData)
	require.Nil(t, err)

	assert.Equal(t, "ok", responseData.Message)
	assert.Equal(t, "http://example.com/people/1", res.HeaderMap.Get("Location"))

	assert.Len(t, people, 1)
	assert.Equal(t, "1", people["1"].Id)
	assert.Equal(t, "Dorota Justyna", people["1"].Given)
	assert.Equal(t, "Zawadzka", people["1"].Surname)
	assert.Equal(t, gFemale, people["1"].Gender)
}

/* Check if payload errors are correctly handled */
func TestCreatePersonRequestPayload(t *testing.T) {
	router := setupRouter()

	people = map[string]personRecord{}

	// Invalid gender field value:

	person1 := testPersonJson{
		Id: "1",
		Given: "Dorota Justyna",
		Surname: "Zawadzka",
		Gender: "INVALID"}

	json1, err := json.Marshal(person1)
	require.Nil(t, err)

	res := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/people", bytes.NewBuffer(json1))
	router.ServeHTTP(res, req)

	assert.Equal(t, res.Code, http.StatusBadRequest)

	responseData := testErrorJson{}

	require.True(t, json.Valid(res.Body.Bytes()))
	err = json.Unmarshal(res.Body.Bytes(), &responseData)
	require.Nil(t, err)

	assert.Equal(t, payloadErrorMsg, responseData.Message)

	// Id field not specified:

	person2 := testPersonJson{
		Given: "Antoni",
		Surname: "Wiśniewski",
		Gender: gMale}

	json2, err := json.Marshal(person2)
	require.Nil(t, err)

	res = httptest.NewRecorder()
	req = httptest.NewRequest("POST", "/people", bytes.NewBuffer(json2))
	router.ServeHTTP(res, req)

	assert.Equal(t, res.Code, http.StatusBadRequest)

	responseData = testErrorJson{}

	require.True(t, json.Valid(res.Body.Bytes()))
	err = json.Unmarshal(res.Body.Bytes(), &responseData)
	require.Nil(t, err)

	assert.Equal(t, payloadErrorMsg, responseData.Message)
}

/* Test if the retrieve people endpoint correctly deals with empty database */
func TestRetrievePeopleRequestEmpty(t *testing.T) {
	router := setupRouter()

	people = map[string]personRecord{}

	res := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/people", nil)
	router.ServeHTTP(res, req)

	assert.Equal(t, res.Code, http.StatusOK)

	type responseJson struct {
		Pagination testPaginationJson `json:"pagination"`
		Records    []testPersonJson   `json:"records"`
	}

	var responseData responseJson

	require.True(t, json.Valid(res.Body.Bytes()))
	err := json.Unmarshal(res.Body.Bytes(), &responseData)
	require.Nil(t, err)

	assert.Len(t, responseData.Records, 0)
	assert.Empty(t, responseData.Pagination.NextUrl)
	assert.Empty(t, responseData.Pagination.PrevUrl)
}

/* Test if the retrieve people endpoint correctly handles result data pagination */
func TestRetrievePeopleRequestPagination(t *testing.T) {
	type responseJson struct {
		Pagination testPaginationJson `json:"pagination"`
		Records    []testPersonJson   `json:"records"`
	}

	var responseData responseJson

	people = map[string]personRecord{
		"P01": personRecord{"P01", "Lidia", "Błaszczyk", gFemale},
		"P02": personRecord{"P02", "Lara", "Szymańska", gFemale},
		"P03": personRecord{"P03", "Radosław", "Kołodziej", gMale},
		"P04": personRecord{"P04", "Antonina", "Kozłowska", gFemale},
		"P05": personRecord{"P05", "Marcela", "Szymczak", gFemale},
		"P06": personRecord{"P06", "Bruno", "Maciejewski", gMale},
		"P07": personRecord{"P07", "Mirosława", "Czarnecka", gFemale},
		"P08": personRecord{"P08", "Elena", "Szewczyk", gFemale},
		"P09": personRecord{"P09", "Ariel", "Zalewski", gMale},
		"P10": personRecord{"P10", "Florian", "Jankowski", gMale},
		"P11": personRecord{"P11", "Borys", "Kalinowski", gMale},
		"P12": personRecord{"P12", "Oliwia", "Cieślak", gFemale},
		"P13": personRecord{"P13", "Natalia", "Ziółkowska", gFemale},
	}

	router := setupRouter()

	// Request the first page:

	res := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/people?limit=10&page=0", nil)
	router.ServeHTTP(res, req)

	assert.Equal(t, res.Code, http.StatusOK)

	require.True(t, json.Valid(res.Body.Bytes()))
	err := json.Unmarshal(res.Body.Bytes(), &responseData)
	require.Nil(t, err)

	assert.Len(t, responseData.Records, 10)
	assert.Equal(t, "P01", responseData.Records[0].Id)
	assert.Equal(t, "Lidia", responseData.Records[0].Given)
	assert.Equal(t, "Błaszczyk", responseData.Records[0].Surname)
	assert.Equal(t, gFemale, responseData.Records[0].Gender)

	assert.Equal(t, "P02", responseData.Records[1].Id)
	assert.Equal(t, "Lara", responseData.Records[1].Given)
	assert.Equal(t, "Szymańska", responseData.Records[1].Surname)
	assert.Equal(t, gFemale, responseData.Records[1].Gender)

	assert.Equal(t, "P03", responseData.Records[2].Id)
	assert.Equal(t, "Radosław", responseData.Records[2].Given)
	assert.Equal(t, "Kołodziej", responseData.Records[2].Surname)
	assert.Equal(t, gMale, responseData.Records[2].Gender)

	assert.Equal(t, "P04", responseData.Records[3].Id)
	assert.Equal(t, "Antonina", responseData.Records[3].Given)
	assert.Equal(t, "Kozłowska", responseData.Records[3].Surname)
	assert.Equal(t, gFemale, responseData.Records[3].Gender)

	assert.Equal(t, "P05", responseData.Records[4].Id)
	assert.Equal(t, "Marcela", responseData.Records[4].Given)
	assert.Equal(t, "Szymczak", responseData.Records[4].Surname)
	assert.Equal(t, gFemale, responseData.Records[4].Gender)

	assert.Equal(t, "P06", responseData.Records[5].Id)
	assert.Equal(t, "Bruno", responseData.Records[5].Given)
	assert.Equal(t, "Maciejewski", responseData.Records[5].Surname)
	assert.Equal(t, gMale, responseData.Records[5].Gender)

	assert.Equal(t, "P07", responseData.Records[6].Id)
	assert.Equal(t, "Mirosława", responseData.Records[6].Given)
	assert.Equal(t, "Czarnecka", responseData.Records[6].Surname)
	assert.Equal(t, gFemale, responseData.Records[6].Gender)

	assert.Equal(t, "P08", responseData.Records[7].Id)
	assert.Equal(t, "Elena", responseData.Records[7].Given)
	assert.Equal(t, "Szewczyk", responseData.Records[7].Surname)
	assert.Equal(t, gFemale, responseData.Records[7].Gender)

	assert.Equal(t, "P09", responseData.Records[8].Id)
	assert.Equal(t, "Ariel", responseData.Records[8].Given)
	assert.Equal(t, "Zalewski", responseData.Records[8].Surname)
	assert.Equal(t, gMale, responseData.Records[8].Gender)

	assert.Equal(t, "P10", responseData.Records[9].Id)
	assert.Equal(t, "Florian", responseData.Records[9].Given)
	assert.Equal(t, "Jankowski", responseData.Records[9].Surname)
	assert.Equal(t, gMale, responseData.Records[9].Gender)

	assert.Equal(t, "http://example.com/people?limit=10&page=1", responseData.Pagination.NextUrl)
	assert.Empty(t, responseData.Pagination.PrevUrl)

	// Request the second page:

	res = httptest.NewRecorder()
	req = httptest.NewRequest("GET", responseData.Pagination.NextUrl, nil)
	router.ServeHTTP(res, req)

	assert.Equal(t, res.Code, http.StatusOK)

	require.True(t, json.Valid(res.Body.Bytes()))
	responseData = responseJson{}
	err = json.Unmarshal(res.Body.Bytes(), &responseData)
	require.Nil(t, err)

	assert.Len(t, responseData.Records, 3)
	assert.Equal(t, "P11", responseData.Records[0].Id)
	assert.Equal(t, "Borys", responseData.Records[0].Given)
	assert.Equal(t, "Kalinowski", responseData.Records[0].Surname)
	assert.Equal(t, gMale, responseData.Records[0].Gender)

	assert.Equal(t, "P12", responseData.Records[1].Id)
	assert.Equal(t, "Oliwia", responseData.Records[1].Given)
	assert.Equal(t, "Cieślak", responseData.Records[1].Surname)
	assert.Equal(t, gFemale, responseData.Records[1].Gender)

	assert.Equal(t, "P13", responseData.Records[2].Id)
	assert.Equal(t, "Natalia", responseData.Records[2].Given)
	assert.Equal(t, "Ziółkowska", responseData.Records[2].Surname)
	assert.Equal(t, gFemale, responseData.Records[2].Gender)

	assert.Empty(t, responseData.Pagination.NextUrl)
	assert.Equal(t, "http://example.com/people?limit=10&page=0", responseData.Pagination.PrevUrl)
}

/* Test if the retrieve people endpoint correctly handles invalid pagination parameters */
func TestRetrievePeopleRequestPaginationParams(t *testing.T) {
	router := setupRouter()

	// Request negative page:

	res := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/people?page=-1", nil)
	router.ServeHTTP(res, req)

	assert.Equal(t, res.Code, http.StatusBadRequest)

	responseData := testErrorJson{}

	require.True(t, json.Valid(res.Body.Bytes()))
	err := json.Unmarshal(res.Body.Bytes(), &responseData)
	require.Nil(t, err)

	assert.Equal(t, queryErrorMsg, responseData.Message)

	// Request too small page size:

	res = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "/people?limit=5", nil)
	router.ServeHTTP(res, req)

	assert.Equal(t, res.Code, http.StatusBadRequest)

	responseData = testErrorJson{}

	require.True(t, json.Valid(res.Body.Bytes()))
	err = json.Unmarshal(res.Body.Bytes(), &responseData)
	require.Nil(t, err)

	assert.Equal(t, queryErrorMsg, responseData.Message)

	// Request too big page size:

	res = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "/people?limit=1000", nil)
	router.ServeHTTP(res, req)

	assert.Equal(t, res.Code, http.StatusBadRequest)

	responseData = testErrorJson{}

	require.True(t, json.Valid(res.Body.Bytes()))
	err = json.Unmarshal(res.Body.Bytes(), &responseData)
	require.Nil(t, err)

	assert.Equal(t, queryErrorMsg, responseData.Message)
}
