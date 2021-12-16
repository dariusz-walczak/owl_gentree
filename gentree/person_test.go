package main

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
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

func testJsonBody(t *testing.T, payload interface{}) io.Reader {
	strData, err := json.Marshal(payload)
	require.Nil(t, err)
	return bytes.NewBuffer(strData)
}

func testMakeRequest(router *gin.Engine, method string, url string, body io.Reader) *httptest.ResponseRecorder {
	res := httptest.NewRecorder()
	req := httptest.NewRequest(method, url, body)
	router.ServeHTTP(res, req)

	return res
}

func testJsonRes(t *testing.T, res *httptest.ResponseRecorder, payload interface{}) {
	require.True(t, json.Valid(res.Body.Bytes()))
	err := json.Unmarshal(res.Body.Bytes(), &payload)
	require.Nil(t, err)
}

func testPersonRes(t *testing.T, res *httptest.ResponseRecorder) testPersonJson {
	payload := testPersonJson{}
	testJsonRes(t, res, &payload)
	return payload
}

type testErrorJson struct {
	Message string `json:"message"`
}

func testErrorRes(t *testing.T, res *httptest.ResponseRecorder) testErrorJson {
	payload := testErrorJson{}
	testJsonRes(t, res, &payload)
	return payload
}

type testLocationJson struct {
	Message  string `json:"message"`
	Location string `json:"location"`
}

func testLocationRes(t *testing.T, res *httptest.ResponseRecorder) testLocationJson {
	payload := testLocationJson{}
	testJsonRes(t, res, &payload)
	return payload
}

type testPersonListJson struct {
	Pagination testPaginationJson `json:"pagination"`
	Records    []testPersonJson   `json:"records"`
}

func testPersonListRes(t *testing.T, res *httptest.ResponseRecorder) testPersonListJson {
	payload := testPersonListJson{}
	testJsonRes(t, res, &payload)
	return payload
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

	person := testPersonJson{
		Id:      "1",
		Given:   "Dorota Justyna",
		Surname: "Zawadzka",
		Gender:  gFemale}

	res := testMakeRequest(router, "POST", "/people", testJsonBody(t, person))

	assert.Equal(t, res.Code, http.StatusCreated)

	resData := testErrorRes(t, res)

	assert.Equal(t, "ok", resData.Message)
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

	person := testPersonJson{
		Id:      "1",
		Given:   "Eliza",
		Surname: "Wojciechowska",
		Gender:  "INVALID"}

	res := testMakeRequest(router, "POST", "/people", testJsonBody(t, person))

	assert.Equal(t, res.Code, http.StatusBadRequest)

	resData := testErrorRes(t, res)

	assert.Equal(t, payloadErrorMsg, resData.Message)

	// Id field not specified:

	person = testPersonJson{
		Given:   "Antoni",
		Surname: "Wiśniewski",
		Gender:  gMale}

	res = testMakeRequest(router, "POST", "/people", testJsonBody(t, person))

	assert.Equal(t, res.Code, http.StatusBadRequest)

	resData = testErrorRes(t, res)

	assert.Equal(t, payloadErrorMsg, resData.Message)
}

/* Check if the create person handler correctly handles the case of already existing person.
   Confirm that the returned location string works as expected */
func TestCreatePersonRequestExists(t *testing.T) {
	router := setupRouter()

	people = map[string]personRecord{
		"X99": personRecord{
			Id:      "X99",
			Given:   "Marian",
			Surname: "Zakrzewski",
			Gender:  gMale}}

	person := testPersonJson{
		Id:      "X99",
		Given:   "Maria",
		Surname: "Zakrzewska",
		Gender:  gFemale}

	res := testMakeRequest(router, "POST", "/people", testJsonBody(t, person))

	assert.Equal(t, res.Code, http.StatusBadRequest)

	resData1 := testLocationRes(t, res)

	assert.Equal(t, "Person (X99) already exists", resData1.Message)
	assert.Equal(t, "http://example.com/people/X99", resData1.Location)

	// Retrieve the record to confirm that the Id field not specified:
	res = testMakeRequest(router, "GET", resData1.Location, nil)

	assert.Equal(t, res.Code, http.StatusOK)

	resData2 := testPersonRes(t, res)

	assert.Equal(t, "X99", resData2.Id)
	assert.Equal(t, "Marian", resData2.Given)
	assert.Equal(t, "Zakrzewski", resData2.Surname)
	assert.Equal(t, gMale, resData2.Gender)
}

/* Test if the replace person endpoint overrides the existing record completely

   The person identifier is taken from the request uri. The identifier specified in the payload
   should be ignored if provided */
func TestReplacePersonRequestSuccess(t *testing.T) {
	router := setupRouter()

	people = map[string]personRecord{
		"5rjk": personRecord{
			Id:      "5rjk",
			Given:   "Honorata",
			Surname: "Czarnecka",
			Gender:  gFemale}}

	person := testPersonJson{
		Id:      "ignored",
		Given:   "Gniewomir",
		Surname: "Baranek",
		Gender:  gMale}

	res := testMakeRequest(router, "PUT", "/people/5rjk", testJsonBody(t, person))

	assert.Equal(t, res.Code, http.StatusOK)

	resData := testErrorRes(t, res)

	assert.Equal(t, "Person record replaced", resData.Message)

	assert.Len(t, people, 1)
	assert.Equal(t, "5rjk", people["5rjk"].Id)
	assert.Equal(t, "Gniewomir", people["5rjk"].Given)
	assert.Equal(t, "Baranek", people["5rjk"].Surname)
	assert.Equal(t, gMale, people["5rjk"].Gender)
}

/* Test if the retrieve people endpoint correctly deals with empty database */
func TestRetrievePeopleRequestEmpty(t *testing.T) {
	router := setupRouter()

	people = map[string]personRecord{}

	res := testMakeRequest(router, "GET", "/people", nil)

	assert.Equal(t, res.Code, http.StatusOK)

	resData := testPersonListRes(t, res)

	assert.Len(t, resData.Records, 0)
	assert.Empty(t, resData.Pagination.NextUrl)
	assert.Empty(t, resData.Pagination.PrevUrl)
}

/* Test if the retrieve people endpoint correctly handles result data pagination */
func TestRetrievePeopleRequestPagination(t *testing.T) {
	router := setupRouter()

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

	// Request the first page:

	res := testMakeRequest(router, "GET", "/people?limit=10&page=0", nil)

	assert.Equal(t, res.Code, http.StatusOK)

	resData := testPersonListRes(t, res)

	assert.Len(t, resData.Records, 10)
	assert.Equal(t, "P01", resData.Records[0].Id)
	assert.Equal(t, "Lidia", resData.Records[0].Given)
	assert.Equal(t, "Błaszczyk", resData.Records[0].Surname)
	assert.Equal(t, gFemale, resData.Records[0].Gender)

	assert.Equal(t, "P02", resData.Records[1].Id)
	assert.Equal(t, "Lara", resData.Records[1].Given)
	assert.Equal(t, "Szymańska", resData.Records[1].Surname)
	assert.Equal(t, gFemale, resData.Records[1].Gender)

	assert.Equal(t, "P03", resData.Records[2].Id)
	assert.Equal(t, "Radosław", resData.Records[2].Given)
	assert.Equal(t, "Kołodziej", resData.Records[2].Surname)
	assert.Equal(t, gMale, resData.Records[2].Gender)

	assert.Equal(t, "P04", resData.Records[3].Id)
	assert.Equal(t, "Antonina", resData.Records[3].Given)
	assert.Equal(t, "Kozłowska", resData.Records[3].Surname)
	assert.Equal(t, gFemale, resData.Records[3].Gender)

	assert.Equal(t, "P05", resData.Records[4].Id)
	assert.Equal(t, "Marcela", resData.Records[4].Given)
	assert.Equal(t, "Szymczak", resData.Records[4].Surname)
	assert.Equal(t, gFemale, resData.Records[4].Gender)

	assert.Equal(t, "P06", resData.Records[5].Id)
	assert.Equal(t, "Bruno", resData.Records[5].Given)
	assert.Equal(t, "Maciejewski", resData.Records[5].Surname)
	assert.Equal(t, gMale, resData.Records[5].Gender)

	assert.Equal(t, "P07", resData.Records[6].Id)
	assert.Equal(t, "Mirosława", resData.Records[6].Given)
	assert.Equal(t, "Czarnecka", resData.Records[6].Surname)
	assert.Equal(t, gFemale, resData.Records[6].Gender)

	assert.Equal(t, "P08", resData.Records[7].Id)
	assert.Equal(t, "Elena", resData.Records[7].Given)
	assert.Equal(t, "Szewczyk", resData.Records[7].Surname)
	assert.Equal(t, gFemale, resData.Records[7].Gender)

	assert.Equal(t, "P09", resData.Records[8].Id)
	assert.Equal(t, "Ariel", resData.Records[8].Given)
	assert.Equal(t, "Zalewski", resData.Records[8].Surname)
	assert.Equal(t, gMale, resData.Records[8].Gender)

	assert.Equal(t, "P10", resData.Records[9].Id)
	assert.Equal(t, "Florian", resData.Records[9].Given)
	assert.Equal(t, "Jankowski", resData.Records[9].Surname)
	assert.Equal(t, gMale, resData.Records[9].Gender)

	assert.Equal(t, "http://example.com/people?limit=10&page=1", resData.Pagination.NextUrl)
	assert.Empty(t, resData.Pagination.PrevUrl)

	// Request the second page:

	res = testMakeRequest(router, "GET", resData.Pagination.NextUrl, nil)

	assert.Equal(t, res.Code, http.StatusOK)

	resData = testPersonListRes(t, res)

	assert.Len(t, resData.Records, 3)
	assert.Equal(t, "P11", resData.Records[0].Id)
	assert.Equal(t, "Borys", resData.Records[0].Given)
	assert.Equal(t, "Kalinowski", resData.Records[0].Surname)
	assert.Equal(t, gMale, resData.Records[0].Gender)

	assert.Equal(t, "P12", resData.Records[1].Id)
	assert.Equal(t, "Oliwia", resData.Records[1].Given)
	assert.Equal(t, "Cieślak", resData.Records[1].Surname)
	assert.Equal(t, gFemale, resData.Records[1].Gender)

	assert.Equal(t, "P13", resData.Records[2].Id)
	assert.Equal(t, "Natalia", resData.Records[2].Given)
	assert.Equal(t, "Ziółkowska", resData.Records[2].Surname)
	assert.Equal(t, gFemale, resData.Records[2].Gender)

	assert.Empty(t, resData.Pagination.NextUrl)
	assert.Equal(t, "http://example.com/people?limit=10&page=0", resData.Pagination.PrevUrl)
}

/* Test if the retrieve people endpoint correctly handles invalid pagination parameters */
func TestRetrievePeopleRequestPaginationParams(t *testing.T) {
	router := setupRouter()

	// Request negative page:

	res := testMakeRequest(router, "GET", "/people?page=-1", nil)

	assert.Equal(t, res.Code, http.StatusBadRequest)

	resData := testErrorRes(t, res)

	assert.Equal(t, queryErrorMsg, resData.Message)

	// Request too small page size:

	res = testMakeRequest(router, "GET", "/people?limit=5", nil)

	assert.Equal(t, res.Code, http.StatusBadRequest)

	resData = testErrorRes(t, res)

	assert.Equal(t, queryErrorMsg, resData.Message)

	// Request too big page size:

	res = testMakeRequest(router, "GET", "/people?limit=1000", nil)

	assert.Equal(t, res.Code, http.StatusBadRequest)

	resData = testErrorRes(t, res)

	assert.Equal(t, queryErrorMsg, resData.Message)
}
