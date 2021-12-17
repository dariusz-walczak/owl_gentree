package main

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testIitRelationJson struct {
	Pid1 string `json:"pid1"`
	Pid2 string `json:"pid2"`
	Type string `json:"type"`
}

type testItRelationJson struct {
	Pid  string `json:"pid"`
	Type string `json:"type"`
}

type testRelationIdJson struct {
	Message    string `json:"message"`
	RelationId int64  `json:"relation_id"`
}

func testRelationIdRes(t *testing.T, res *httptest.ResponseRecorder) testRelationIdJson {
	payload := testRelationIdJson{}
	testJsonRes(t, res, &payload)
	return payload
}

/* Test if both the create relation endpoints correctly record new, valid relations

   Test two variants of the create relation action:
   * general (/relations)
   * person specific (/people/:pid/relations)

   1. Test the general relation endpoint with the father relation
   2. Test the general relation endpoint with the mother relation
   3. Test the general relation endpoint with the husband relation
   4. Test the person specific relation endpoint with the father relation
   5. Test the person specific relation endpoint with the mother relation
   6. Test the person specific relation endpoint with the husband relation */
func TestCreateRelationRequestSuccess(t *testing.T) {
	router := setupRouter()

	people = map[string]personRecord{
		"F1P1": personRecord{
			Id:      "F1P1",
			Given:   "Ignacy",
			Surname: "Marciniak",
			Gender:  gMale},
		"F1P2": personRecord{
			Id:      "F1P2",
			Given:   "Sylwia",
			Surname: "Rutkowska",
			Gender:  gFemale},
		"F1P3": personRecord{
			Id:      "F1P3",
			Given:   "Luiza",
			Surname: "Marciniak",
			Gender:  gFemale},
		"F2P1": personRecord{
			Id:      "F2P1",
			Given:   "Cezary",
			Surname: "Cieślak",
			Gender:  gMale},
		"F2P2": personRecord{
			Id:      "F2P2",
			Given:   "Alicja",
			Surname: "Baranowska",
			Gender:  gFemale},
		"F2P3": personRecord{
			Id:      "F2P3",
			Given:   "Denis",
			Surname: "Cieślak",
			Gender:  gMale}}

	// Case 1: General father relation

	iitRelation := testIitRelationJson{
		Pid1: "F1P1",
		Pid2: "F1P3",
		Type: "father"}

	res := testMakeRequest(router, "POST", "/relations", testJsonBody(t, iitRelation))

	assert.Equal(t, http.StatusCreated, res.Code)

	resData1 := testRelationIdRes(t, res)

	assert.Equal(t, "Relation created", resData1.Message)

	// Case 2: General mother relation

	iitRelation = testIitRelationJson{
		Pid1: "F1P2",
		Pid2: "F1P3",
		Type: "mother"}

	res = testMakeRequest(router, "POST", "/relations", testJsonBody(t, iitRelation))

	assert.Equal(t, http.StatusCreated, res.Code)

	resData2 := testRelationIdRes(t, res)

	assert.Equal(t, "Relation created", resData2.Message)

	// Case 3: General husband relation

	iitRelation = testIitRelationJson{
		Pid1: "F1P1",
		Pid2: "F1P2",
		Type: "husband"}

	res = testMakeRequest(router, "POST", "/relations", testJsonBody(t, iitRelation))

	assert.Equal(t, http.StatusCreated, res.Code)

	resData3 := testRelationIdRes(t, res)

	assert.Equal(t, "Relation created", resData3.Message)

	// Case 4: Person specific father relation

	itRelation := testItRelationJson{
		Pid:  "F2P3",
		Type: "father"}

	res = testMakeRequest(router, "POST", "/people/F2P1/relations", testJsonBody(t, itRelation))

	assert.Equal(t, http.StatusCreated, res.Code)

	resData4 := testRelationIdRes(t, res)

	assert.Equal(t, "Relation created", resData4.Message)

	// Case 5: Person specific mother relation

	itRelation = testItRelationJson{
		Pid:  "F2P3",
		Type: "mother"}

	res = testMakeRequest(router, "POST", "/people/F2P2/relations", testJsonBody(t, itRelation))

	assert.Equal(t, http.StatusCreated, res.Code)

	resData5 := testRelationIdRes(t, res)

	assert.Equal(t, "Relation created", resData5.Message)

	// Case 6: Person specific husband relation

	itRelation = testItRelationJson{
		Pid:  "F2P2",
		Type: "husband"}

	res = testMakeRequest(router, "POST", "/people/F2P1/relations", testJsonBody(t, itRelation))

	assert.Equal(t, http.StatusCreated, res.Code)

	resData6 := testRelationIdRes(t, res)

	assert.Equal(t, "Relation created", resData6.Message)

	// Check the final state of the relation table:

	assert.Len(t, relations, 6)

	assert.Equal(t, resData1.RelationId, relations[resData1.RelationId].Id)
	assert.Equal(t, "F1P1", relations[resData1.RelationId].Pid1)
	assert.Equal(t, "F1P3", relations[resData1.RelationId].Pid2)
	assert.Equal(t, relFather, relations[resData1.RelationId].Type)

	assert.Equal(t, resData2.RelationId, relations[resData2.RelationId].Id)
	assert.Equal(t, "F1P2", relations[resData2.RelationId].Pid1)
	assert.Equal(t, "F1P3", relations[resData2.RelationId].Pid2)
	assert.Equal(t, relMother, relations[resData2.RelationId].Type)

	assert.Equal(t, resData3.RelationId, relations[resData3.RelationId].Id)
	assert.Equal(t, "F1P1", relations[resData3.RelationId].Pid1)
	assert.Equal(t, "F1P2", relations[resData3.RelationId].Pid2)
	assert.Equal(t, relHusband, relations[resData3.RelationId].Type)

	assert.Equal(t, resData4.RelationId, relations[resData4.RelationId].Id)
	assert.Equal(t, "F2P1", relations[resData4.RelationId].Pid1)
	assert.Equal(t, "F2P3", relations[resData4.RelationId].Pid2)
	assert.Equal(t, relFather, relations[resData4.RelationId].Type)

	assert.Equal(t, resData5.RelationId, relations[resData5.RelationId].Id)
	assert.Equal(t, "F2P2", relations[resData5.RelationId].Pid1)
	assert.Equal(t, "F2P3", relations[resData5.RelationId].Pid2)
	assert.Equal(t, relMother, relations[resData5.RelationId].Type)

	assert.Equal(t, resData6.RelationId, relations[resData6.RelationId].Id)
	assert.Equal(t, "F2P1", relations[resData6.RelationId].Pid1)
	assert.Equal(t, "F2P2", relations[resData6.RelationId].Pid2)
	assert.Equal(t, relHusband, relations[resData6.RelationId].Type)
}

/* Test if both the create relation endpoints correctly handle an attempt to create an already
   existing relation

   Test two variants of the create relation action:
   1. general (/relations)
   2. person specific (/people/:pid/relations) */
func TestCreateRelationRequestExists(t *testing.T) {
	router := setupRouter()

	people = map[string]personRecord{
		"A": personRecord{
			Id:      "A",
			Given:   "Mirosław",
			Surname: "Woźniak",
			Gender:  gMale},
		"B": personRecord{
			Id:      "B",
			Given:   "Oksana",
			Surname: "Włodarczyk",
			Gender:  gFemale},
		"C": personRecord{
			Id:      "C",
			Given:   "Olimpia",
			Surname: "Woźniak",
			Gender:  gFemale}}

	relations = map[int64]relationRecord{
		1: relationRecord{Id: 1, Pid1: "A", Pid2: "B", Type: relHusband},
		2: relationRecord{Id: 2, Pid1: "B", Pid2: "C", Type: relMother},
		3: relationRecord{Id: 3, Pid1: "A", Pid2: "C", Type: relFather}}

	// Case 1: General relation

	iitRelation := testIitRelationJson{
		Pid1: "A",
		Pid2: "B",
		Type: "husband"}

	res := testMakeRequest(router, "POST", "/relations", testJsonBody(t, iitRelation))

	assert.Equal(t, http.StatusBadRequest, res.Code)

	resData := testLocationRes(t, res)

	assert.Equal(t, "Relation (A, husband, B) already exists", resData.Message)
	assert.Equal(t, "http://example.com/relations/1", resData.Location)

	// Case 2: Person specific relation

	itRelation := testItRelationJson{
		Pid:  "C",
		Type: "mother"}

	res = testMakeRequest(router, "POST", "/people/B/relations", testJsonBody(t, itRelation))

	assert.Equal(t, http.StatusBadRequest, res.Code)

	resData = testLocationRes(t, res)

	assert.Equal(t, "Relation (B, mother, C) already exists", resData.Message)
	assert.Equal(t, "http://example.com/relations/2", resData.Location)
}

/* Test if the create person specific relation endpoint handles invalid person id format
   correctly */
func TestCreateRelationRequestPidError(t *testing.T) {
	router := setupRouter()

	people = map[string]personRecord{
		"A": personRecord{
			Id:      "A",
			Given:   "Patrycja",
			Surname: "Ziółkowska",
			Gender:  gFemale},
		"B": personRecord{
			Id:      "B",
			Given:   "Ireneusz",
			Surname: "Szymczak",
			Gender:  gMale}}

	relations = map[int64]relationRecord{}

	relation := testItRelationJson{
		Pid:  "A",
		Type: "husband"}

	res := testMakeRequest(router, "POST", "/people/B!/relations", testJsonBody(t, relation))

	assert.Equal(t, http.StatusBadRequest, res.Code)

	resData := testErrorRes(t, res)

	assert.Equal(t, uriErrorMsg, resData.Message)
}

/* Test if both the create relation endpoints handle payload data format errors correctly

   1. The general handler should indicate an error when the relation type is invalid
   2. The person specific handler should indicate an error when the relation type is
      invalid */
func TestCreateRelationRequestPayloadError(t *testing.T) {
	router := setupRouter()

	people = map[string]personRecord{
		"526839f0": personRecord{
			Id:      "526839f0",
			Given:   "Dominik",
			Surname: "Wiśniewski",
			Gender:  gMale},
		"38e205fa": personRecord{
			Id:      "38e205fa",
			Given:   "Jacek",
			Surname: "Baran",
			Gender:  gMale}}

	relations = map[int64]relationRecord{}

	// Case 1: General handler

	iitRelation := testIitRelationJson{
		Pid1: "526839f0",
		Pid2: "38e205fa",
		Type: "invalid"}

	res := testMakeRequest(router, "POST", "/relations", testJsonBody(t, iitRelation))

	assert.Equal(t, http.StatusBadRequest, res.Code)

	resData := testErrorRes(t, res)

	assert.Equal(t, payloadErrorMsg, resData.Message)

	// Case 2: Person specific handler
	itRelation := testItRelationJson{
		Pid:  "38e205fa",
		Type: "invalid"}

	res = testMakeRequest(
		router, "POST", "/people/526839f0/relations", testJsonBody(t, itRelation))

	assert.Equal(t, http.StatusBadRequest, res.Code)

	resData = testErrorRes(t, res)

	assert.Equal(t, payloadErrorMsg, resData.Message)
}
