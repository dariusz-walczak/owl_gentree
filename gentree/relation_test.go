package main

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testFullRelationJson struct {
	Id   int64  `json:"id"`
	Pid1 string `json:"pid1"`
	Pid2 string `json:"pid2"`
	Type string `json:"type"`
}

func testFullRelationRes(t *testing.T, res *httptest.ResponseRecorder) testFullRelationJson {
	payload := testFullRelationJson{}
	testJsonRes(t, res, &payload)
	return payload
}

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

type testRelationListJson struct {
	Pagination testPaginationJson     `json:"pagination"`
	Records    []testFullRelationJson `json:"records"`
}

func testRelationListRes(t *testing.T, res *httptest.ResponseRecorder) testRelationListJson {
	payload := testRelationListJson{}
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

/* Test if both the create relation endpoints handle relation validation errors correctly

   1. The general handler should indicate an error when the relation is invalid
   2. The person specific handler should indicate an error when the relation is invalid */
func TestCreateRelationRequestValidationError(t *testing.T) {
	router := setupRouter()

	people = map[string]personRecord{
		"9596": personRecord{
			Id:      "9596",
			Given:   "Dorian",
			Surname: "Piotrowski",
			Gender:  gMale}}

	relations = map[int64]relationRecord{}

	// Case 1: General handler

	iitRelation := testIitRelationJson{
		Pid1: "9596",
		Pid2: "3141",
		Type: "father"}

	res := testMakeRequest(router, "POST", "/relations", testJsonBody(t, iitRelation))

	assert.Equal(t, http.StatusBadRequest, res.Code)

	resData := testErrorRes(t, res)

	assert.Equal(t, "Relation (9596, father, 3141) is invalid", resData.Message)

	// Case 2: Person specific handler
	itRelation := testItRelationJson{
		Pid:  "3141",
		Type: "father"}

	res = testMakeRequest(
		router, "POST", "/people/9596/relations", testJsonBody(t, itRelation))

	assert.Equal(t, http.StatusBadRequest, res.Code)

	resData = testErrorRes(t, res)

	assert.Equal(t, "Relation (9596, father, 3141) is invalid", resData.Message)
}

/* Test the delete relation endpoint

   1. Check the success scenario (existing record correctly deleted)
   2. Check the case of invalid relation id format (part of the url)
   3. Check the case of missing relation */
func TestDeleteRelationRequest(t *testing.T) {
	router := setupRouter()

	people = map[string]personRecord{
		"cf": personRecord{
			Id:      "cf",
			Given:   "Bartłomiej",
			Surname: "Sokołowski",
			Gender:  gMale},
		"9c": personRecord{
			Id:      "9c",
			Given:   "Celina",
			Surname: "Kamińska",
			Gender:  gFemale},
		"b4": personRecord{
			Id:      "b4",
			Given:   "Lidia",
			Surname: "Sokołowska",
			Gender:  gFemale}}

	relations = map[int64]relationRecord{
		88128: relationRecord{Id: 88128, Pid1: "9c", Pid2: "b4", Type: relMother},
		86917: relationRecord{Id: 86917, Pid1: "cf", Pid2: "b4", Type: relFather},
		51235: relationRecord{Id: 51235, Pid1: "cf", Pid2: "9c", Type: relHusband}}

	// Case 1: Successful deletion

	res := testMakeRequest(router, "DELETE", "/relations/51235", nil)

	assert.Equal(t, http.StatusOK, res.Code)

	resData := testErrorRes(t, res)

	assert.Equal(t, "Relation deleted", resData.Message)

	// TODO: people map

	// Case 2: Invalid relation id

	res = testMakeRequest(router, "DELETE", "/relations/e90252c6", nil)

	assert.Equal(t, http.StatusBadRequest, res.Code)

	resData = testErrorRes(t, res)

	assert.Equal(t, uriErrorMsg, resData.Message)

	// Case 3: Missing relation

	res = testMakeRequest(router, "DELETE", "/relations/28944", nil)

	assert.Equal(t, http.StatusNotFound, res.Code)

	resData = testErrorRes(t, res)

	assert.Equal(t, "Unknown relation id", resData.Message)
}

/* Test the replace relation endpoint

   1. Test the success scenario (existing record correctly replaced)
   2. Test the case of invalid relation id format (part of url)
   3. Test the case of missing relation
   4. Test the case of invalid person id (payload) */
func TestReplaceRelationRequest(t *testing.T) {
	router := setupRouter()

	people = map[string]personRecord{
		"A": personRecord{
			Id:      "A",
			Given:   "Magda",
			Surname: "Krawczyk",
			Gender:  gFemale},
		"B": personRecord{
			Id:      "B",
			Given:   "Edyta",
			Surname: "Andrzejewska",
			Gender:  gFemale},
		"C": personRecord{
			Id:      "C",
			Given:   "Konstanty",
			Surname: "Andrzejewski",
			Gender:  gMale},
		"D": personRecord{
			Id:      "D",
			Given:   "Dorian",
			Surname: "Czarnecki",
			Gender:  gMale},
	}

	relations = map[int64]relationRecord{
		1: relationRecord{Id: 1, Pid1: "A", Pid2: "B", Type: relMother},
		2: relationRecord{Id: 2, Pid1: "C", Pid2: "B", Type: relFather}}

	// Case 1: Replacement success

	iitRelation := testIitRelationJson{
		Pid1: "D",
		Pid2: "B",
		Type: relFather}

	res := testMakeRequest(router, "PUT", "/relations/2", testJsonBody(t, iitRelation))

	assert.Equal(t, http.StatusOK, res.Code)

	resData := testErrorRes(t, res)

	assert.Equal(t, "Relation record replaced", resData.Message)

	assert.Len(t, relations, 2)

	assert.Equal(t, int64(1), relations[1].Id)
	assert.Equal(t, "A", relations[1].Pid1)
	assert.Equal(t, "B", relations[1].Pid2)
	assert.Equal(t, relMother, relations[1].Type)

	assert.Equal(t, int64(2), relations[2].Id)
	assert.Equal(t, "D", relations[2].Pid1)
	assert.Equal(t, "B", relations[2].Pid2)
	assert.Equal(t, relFather, relations[2].Type)

	// Case 2: Invalid relation id

	iitRelation = testIitRelationJson{
		Pid1: "C",
		Pid2: "B",
		Type: relFather}

	res = testMakeRequest(router, "PUT", "/relations/R001", testJsonBody(t, iitRelation))

	assert.Equal(t, http.StatusBadRequest, res.Code)

	resData = testErrorRes(t, res)

	assert.Equal(t, uriErrorMsg, resData.Message)

	assert.Len(t, relations, 2)

	assert.Equal(t, int64(1), relations[1].Id)
	assert.Equal(t, "A", relations[1].Pid1)
	assert.Equal(t, "B", relations[1].Pid2)
	assert.Equal(t, relMother, relations[1].Type)

	assert.Equal(t, int64(2), relations[2].Id)
	assert.Equal(t, "D", relations[2].Pid1)
	assert.Equal(t, "B", relations[2].Pid2)
	assert.Equal(t, relFather, relations[2].Type)

	// Case 3: Missing relation

	// iitRelation assignment skipped intentionally (the one defined for the 2nd case is ok)

	res = testMakeRequest(router, "PUT", "/relations/3", testJsonBody(t, iitRelation))

	assert.Equal(t, http.StatusNotFound, res.Code)

	resData = testErrorRes(t, res)

	assert.Equal(t, "Unknown relation id", resData.Message)

	assert.Len(t, relations, 2)

	assert.Equal(t, int64(1), relations[1].Id)
	assert.Equal(t, "A", relations[1].Pid1)
	assert.Equal(t, "B", relations[1].Pid2)
	assert.Equal(t, relMother, relations[1].Type)

	assert.Equal(t, int64(2), relations[2].Id)
	assert.Equal(t, "D", relations[2].Pid1)
	assert.Equal(t, "B", relations[2].Pid2)
	assert.Equal(t, relFather, relations[2].Type)

	// Case 4: Invalid person id

	iitRelation = testIitRelationJson{
		Pid1: "",
		Pid2: "B",
		Type: relFather}

	res = testMakeRequest(router, "PUT", "/relations/2", testJsonBody(t, iitRelation))

	assert.Equal(t, http.StatusBadRequest, res.Code)

	resData = testErrorRes(t, res)

	assert.Equal(t, payloadErrorMsg, resData.Message)

	assert.Len(t, relations, 2)

	assert.Equal(t, int64(1), relations[1].Id)
	assert.Equal(t, "A", relations[1].Pid1)
	assert.Equal(t, "B", relations[1].Pid2)
	assert.Equal(t, relMother, relations[1].Type)

	assert.Equal(t, int64(2), relations[2].Id)
	assert.Equal(t, "D", relations[2].Pid1)
	assert.Equal(t, "B", relations[2].Pid2)
	assert.Equal(t, relFather, relations[2].Type)
}

/* Test the retrieve relation endpoint

   1. Test the success scenario (existing record correctly returned)
   2. Test the case of invalid relation id format (part of url)
   3. Test the case of missing relation */
func TestRetrieveRelationRequest(t *testing.T) {
	router := setupRouter()

	people = map[string]personRecord{
		"f6b6": personRecord{
			Id:      "f6b6",
			Given:   "Florian",
			Surname: "Krajewski",
			Gender:  gMale},
		"b0dc": personRecord{
			Id:      "b0dc",
			Given:   "Julia",
			Surname: "Kozłowska",
			Gender:  gFemale},
		"f870": personRecord{
			Id:      "f870",
			Given:   "Krzysztof",
			Surname: "Krajewski",
			Gender:  gMale}}

	relations = map[int64]relationRecord{
		20547: relationRecord{Id: 20547, Pid1: "b0dc", Pid2: "f870", Type: relMother},
		11646: relationRecord{Id: 11646, Pid1: "f6b6", Pid2: "f870", Type: relFather}}

	// Case 1: Successful retrieval

	res := testMakeRequest(router, "GET", "/relations/20547", nil)

	assert.Equal(t, http.StatusOK, res.Code)

	resData1 := testFullRelationRes(t, res)

	assert.Equal(t, int64(20547), resData1.Id)
	assert.Equal(t, "b0dc", resData1.Pid1)
	assert.Equal(t, "f870", resData1.Pid2)
	assert.Equal(t, "mother", resData1.Type)

	// Case 2: Invalid relation id

	res = testMakeRequest(router, "GET", "/relations/c909debd", nil)

	assert.Equal(t, http.StatusBadRequest, res.Code)

	resData2 := testErrorRes(t, res)

	assert.Equal(t, uriErrorMsg, resData2.Message)

	// Case 3: Missing relation

	res = testMakeRequest(router, "GET", "/relations/55752", nil)

	assert.Equal(t, http.StatusNotFound, res.Code)

	resData3 := testErrorRes(t, res)

	assert.Equal(t, "Unknown relation id", resData3.Message)
}

/* Test if the retrieve relations endpoint works as expected

   1. Test the successful retrieval of the first page (existing records sorted, divided into parts
      and the page returned)
   2. Test the successful retrieval of the second page
   3. Test the handling of the negative pagination page index (to test pagination binding error) */
func TestRetrieveRelationsRequest(t *testing.T) {
	router := setupRouter()

	people = map[string]personRecord{
		"P01": personRecord{Id: "P01", Given: "Marian", Surname: "Zawadzki", Gender: gMale},
		"P02": personRecord{Id: "P02", Given: "Marlena", Surname: "Pawlak", Gender: gFemale},
		"P03": personRecord{Id: "P03", Given: "Urszula", Surname: "Zawadzka", Gender: gFemale},
		"P04": personRecord{Id: "P04", Given: "Mikołaj", Surname: "Zawadzki", Gender: gMale},
		"P05": personRecord{Id: "P05", Given: "Radosław", Surname: "Malinowski", Gender: gMale},
		"P06": personRecord{Id: "P06", Given: "Weronika", Surname: "Krajewska", Gender: gFemale},
		"P07": personRecord{Id: "P07", Given: "Dorota", Surname: "Malinowska", Gender: gFemale},
		"P08": personRecord{Id: "P08", Given: "Magdalena", Surname: "Malinowska", Gender: gFemale},
		"P09": personRecord{Id: "P09", Given: "Anna", Surname: "Malinowska", Gender: gFemale},
		"P10": personRecord{Id: "P10", Given: "Emanuel", Surname: "Witkowski", Gender: gMale}}

	relations = map[int64]relationRecord{
		10: relationRecord{Id: 10, Pid1: "P01", Pid2: "P02", Type: relHusband},
		11: relationRecord{Id: 11, Pid1: "P01", Pid2: "P03", Type: relFather},
		12: relationRecord{Id: 12, Pid1: "P02", Pid2: "P03", Type: relMother},
		13: relationRecord{Id: 13, Pid1: "P01", Pid2: "P04", Type: relFather},
		14: relationRecord{Id: 14, Pid1: "P02", Pid2: "P04", Type: relMother},
		15: relationRecord{Id: 15, Pid1: "P05", Pid2: "P06", Type: relHusband},
		16: relationRecord{Id: 16, Pid1: "P05", Pid2: "P07", Type: relFather},
		17: relationRecord{Id: 17, Pid1: "P06", Pid2: "P07", Type: relMother},
		18: relationRecord{Id: 18, Pid1: "P05", Pid2: "P08", Type: relFather},
		19: relationRecord{Id: 19, Pid1: "P06", Pid2: "P08", Type: relMother},
		20: relationRecord{Id: 20, Pid1: "P05", Pid2: "P09", Type: relFather},
		21: relationRecord{Id: 21, Pid1: "P06", Pid2: "P09", Type: relMother},
		22: relationRecord{Id: 22, Pid1: "P04", Pid2: "P07", Type: relHusband},
		23: relationRecord{Id: 23, Pid1: "P10", Pid2: "P09", Type: relHusband}}

	// Case 1: Successful retrieval of the first page

	res := testMakeRequest(router, "GET", "/relations?limit=10&page=0", nil)

	assert.Equal(t, http.StatusOK, res.Code)

	resData1 := testRelationListRes(t, res)

	assert.Len(t, resData1.Records, 10)
	assert.Equal(t, int64(10), resData1.Records[0].Id)
	assert.Equal(t, "P01", resData1.Records[0].Pid1)
	assert.Equal(t, "P02", resData1.Records[0].Pid2)
	assert.Equal(t, relHusband, resData1.Records[0].Type)

	assert.Equal(t, int64(11), resData1.Records[1].Id)
	assert.Equal(t, "P01", resData1.Records[1].Pid1)
	assert.Equal(t, "P03", resData1.Records[1].Pid2)
	assert.Equal(t, relFather, resData1.Records[1].Type)

	assert.Equal(t, int64(12), resData1.Records[2].Id)
	assert.Equal(t, "P02", resData1.Records[2].Pid1)
	assert.Equal(t, "P03", resData1.Records[2].Pid2)
	assert.Equal(t, relMother, resData1.Records[2].Type)

	assert.Equal(t, int64(13), resData1.Records[3].Id)
	assert.Equal(t, "P01", resData1.Records[3].Pid1)
	assert.Equal(t, "P04", resData1.Records[3].Pid2)
	assert.Equal(t, relFather, resData1.Records[3].Type)

	assert.Equal(t, int64(14), resData1.Records[4].Id)
	assert.Equal(t, "P02", resData1.Records[4].Pid1)
	assert.Equal(t, "P04", resData1.Records[4].Pid2)
	assert.Equal(t, relMother, resData1.Records[4].Type)

	assert.Equal(t, int64(15), resData1.Records[5].Id)
	assert.Equal(t, "P05", resData1.Records[5].Pid1)
	assert.Equal(t, "P06", resData1.Records[5].Pid2)
	assert.Equal(t, relHusband, resData1.Records[5].Type)

	assert.Equal(t, int64(16), resData1.Records[6].Id)
	assert.Equal(t, "P05", resData1.Records[6].Pid1)
	assert.Equal(t, "P07", resData1.Records[6].Pid2)
	assert.Equal(t, relFather, resData1.Records[6].Type)

	assert.Equal(t, int64(17), resData1.Records[7].Id)
	assert.Equal(t, "P06", resData1.Records[7].Pid1)
	assert.Equal(t, "P07", resData1.Records[7].Pid2)
	assert.Equal(t, relMother, resData1.Records[7].Type)

	assert.Equal(t, int64(18), resData1.Records[8].Id)
	assert.Equal(t, "P05", resData1.Records[8].Pid1)
	assert.Equal(t, "P08", resData1.Records[8].Pid2)
	assert.Equal(t, relFather, resData1.Records[8].Type)

	assert.Equal(t, int64(19), resData1.Records[9].Id)
	assert.Equal(t, "P06", resData1.Records[9].Pid1)
	assert.Equal(t, "P08", resData1.Records[9].Pid2)
	assert.Equal(t, relMother, resData1.Records[9].Type)

	assert.Equal(t, "http://example.com/relations?limit=10&page=1", resData1.Pagination.NextUrl)
	assert.Empty(t, resData1.Pagination.PrevUrl)

	// Case 2: Successful retrieval of the second page

	res = testMakeRequest(router, "GET", resData1.Pagination.NextUrl, nil)

	assert.Equal(t, http.StatusOK, res.Code)

	resData2 := testRelationListRes(t, res)

	assert.Len(t, resData2.Records, 4)

	assert.Equal(t, int64(20), resData2.Records[0].Id)
	assert.Equal(t, "P05", resData2.Records[0].Pid1)
	assert.Equal(t, "P09", resData2.Records[0].Pid2)
	assert.Equal(t, relFather, resData2.Records[0].Type)

	assert.Equal(t, int64(21), resData2.Records[1].Id)
	assert.Equal(t, "P06", resData2.Records[1].Pid1)
	assert.Equal(t, "P09", resData2.Records[1].Pid2)
	assert.Equal(t, relMother, resData2.Records[1].Type)

	assert.Equal(t, int64(22), resData2.Records[2].Id)
	assert.Equal(t, "P04", resData2.Records[2].Pid1)
	assert.Equal(t, "P07", resData2.Records[2].Pid2)
	assert.Equal(t, relHusband, resData2.Records[2].Type)

	assert.Equal(t, int64(23), resData2.Records[3].Id)
	assert.Equal(t, "P10", resData2.Records[3].Pid1)
	assert.Equal(t, "P09", resData2.Records[3].Pid2)
	assert.Equal(t, relHusband, resData2.Records[3].Type)

	assert.Empty(t, resData2.Pagination.NextUrl)
	assert.Equal(t, "http://example.com/relations?limit=10&page=0", resData2.Pagination.PrevUrl)

	// Case 3: Negative pagination page index

	res = testMakeRequest(router, "GET", "/relations?page=-1", nil)

	assert.Equal(t, http.StatusBadRequest, res.Code)

	resData3 := testErrorRes(t, res)

	assert.Equal(t, queryErrorMsg, resData3.Message)
}

/* Test if the retrieve person relations endpoint works as expected

   1. Test the successful retrieval of the first page (matching records sorted, divided into parts
      and the page returned)
   2. Test the successful retrieval of the second page
   3. Test the handling of the invalid source person id format (part of the url)
   4. Test the handling of the invalid pagination query variables (too small page size)
   5. Test the handling of not existing person */
func TestRetrievePersonRelationsRequest(t *testing.T) {
	router := setupRouter()

	people = map[string]personRecord{
		"P01": personRecord{Id: "P01", Given: "Albert", Surname: "Michalski", Gender: gMale},
		"P02": personRecord{Id: "P02", Given: "Wioletta", Surname: "Piotrowska", Gender: gFemale},
		"P03": personRecord{Id: "P03", Given: "Adela", Surname: "Michalska", Gender: gFemale},
		"P04": personRecord{Id: "P04", Given: "Joachim", Surname: "Maciejewski", Gender: gMale},
		"P05": personRecord{Id: "P05", Given: "Damian", Surname: "Michalski", Gender: gMale},
		"P06": personRecord{Id: "P06", Given: "Karolina", Surname: "Kołodziej", Gender: gFemale},
		"P07": personRecord{Id: "P07", Given: "Laura", Surname: "Michalska", Gender: gFemale},
		"P08": personRecord{Id: "P08", Given: "Czesław", Surname: "Maciejewski", Gender: gMale},
		"P09": personRecord{Id: "P09", Given: "Gniewomir", Surname: "Michalski", Gender: gFemale},
		"P10": personRecord{Id: "P10", Given: "Ewa", Surname: "Michalska", Gender: gFemale},
		"P11": personRecord{Id: "P11", Given: "Patrycja", Surname: "Michalska", Gender: gFemale},
		"P12": personRecord{Id: "P12", Given: "Monika", Surname: "Maciejewska", Gender: gFemale},
		"P13": personRecord{Id: "P13", Given: "Marcelina", Surname: "Michalska", Gender: gFemale},
		"P14": personRecord{Id: "P14", Given: "Ludwik", Surname: "Michalski", Gender: gMale},
		"P15": personRecord{Id: "P15", Given: "Gustaw", Surname: "Michalski", Gender: gMale},
		"P16": personRecord{Id: "P16", Given: "Edyta", Surname: "Michalska", Gender: gFemale},
		"P17": personRecord{Id: "P17", Given: "Barbara", Surname: "Michalska", Gender: gFemale}}

	relations = map[int64]relationRecord{
		10: relationRecord{Id: 10, Pid1: "P01", Pid2: "P02", Type: relHusband},
		11: relationRecord{Id: 11, Pid1: "P01", Pid2: "P03", Type: relFather},
		12: relationRecord{Id: 12, Pid1: "P02", Pid2: "P03", Type: relMother},
		13: relationRecord{Id: 13, Pid1: "P04", Pid2: "P03", Type: relHusband},
		14: relationRecord{Id: 14, Pid1: "P01", Pid2: "P05", Type: relFather},
		15: relationRecord{Id: 15, Pid1: "P02", Pid2: "P05", Type: relMother},
		16: relationRecord{Id: 16, Pid1: "P05", Pid2: "P06", Type: relHusband},
		17: relationRecord{Id: 17, Pid1: "P05", Pid2: "P07", Type: relFather},
		18: relationRecord{Id: 18, Pid1: "P06", Pid2: "P07", Type: relMother},
		19: relationRecord{Id: 19, Pid1: "P04", Pid2: "P08", Type: relFather},
		20: relationRecord{Id: 20, Pid1: "P03", Pid2: "P08", Type: relMother},
		21: relationRecord{Id: 21, Pid1: "P05", Pid2: "P09", Type: relFather},
		22: relationRecord{Id: 22, Pid1: "P06", Pid2: "P09", Type: relMother},
		23: relationRecord{Id: 23, Pid1: "P05", Pid2: "P10", Type: relFather},
		24: relationRecord{Id: 24, Pid1: "P06", Pid2: "P10", Type: relMother},
		25: relationRecord{Id: 25, Pid1: "P05", Pid2: "P11", Type: relFather},
		26: relationRecord{Id: 26, Pid1: "P06", Pid2: "P11", Type: relMother},
		27: relationRecord{Id: 27, Pid1: "P04", Pid2: "P12", Type: relFather},
		28: relationRecord{Id: 28, Pid1: "P03", Pid2: "P12", Type: relMother},
		29: relationRecord{Id: 29, Pid1: "P05", Pid2: "P13", Type: relFather},
		30: relationRecord{Id: 30, Pid1: "P06", Pid2: "P13", Type: relMother},
		31: relationRecord{Id: 31, Pid1: "P05", Pid2: "P14", Type: relFather},
		32: relationRecord{Id: 32, Pid1: "P06", Pid2: "P14", Type: relMother},
		33: relationRecord{Id: 33, Pid1: "P05", Pid2: "P15", Type: relFather},
		34: relationRecord{Id: 34, Pid1: "P06", Pid2: "P15", Type: relMother},
		35: relationRecord{Id: 35, Pid1: "P05", Pid2: "P16", Type: relFather},
		36: relationRecord{Id: 36, Pid1: "P06", Pid2: "P16", Type: relMother},
		37: relationRecord{Id: 37, Pid1: "P05", Pid2: "P17", Type: relFather},
		38: relationRecord{Id: 38, Pid1: "P06", Pid2: "P17", Type: relMother}}

	// Case 1: Successful retrieval of the first page

	res := testMakeRequest(router, "GET", "/people/P05/relations?limit=10&page=0", nil)

	assert.Equal(t, http.StatusOK, res.Code)

	resData1 := testRelationListRes(t, res)

	assert.Len(t, resData1.Records, 10)

	assert.Equal(t, int64(14), resData1.Records[0].Id)
	assert.Equal(t, "P01", resData1.Records[0].Pid1)
	assert.Equal(t, "P05", resData1.Records[0].Pid2)
	assert.Equal(t, relFather, resData1.Records[0].Type)

	assert.Equal(t, int64(15), resData1.Records[1].Id)
	assert.Equal(t, "P02", resData1.Records[1].Pid1)
	assert.Equal(t, "P05", resData1.Records[1].Pid2)
	assert.Equal(t, relMother, resData1.Records[1].Type)

	assert.Equal(t, int64(16), resData1.Records[2].Id)
	assert.Equal(t, "P05", resData1.Records[2].Pid1)
	assert.Equal(t, "P06", resData1.Records[2].Pid2)
	assert.Equal(t, relHusband, resData1.Records[2].Type)

	assert.Equal(t, int64(17), resData1.Records[3].Id)
	assert.Equal(t, "P05", resData1.Records[3].Pid1)
	assert.Equal(t, "P07", resData1.Records[3].Pid2)
	assert.Equal(t, relFather, resData1.Records[3].Type)

	assert.Equal(t, int64(21), resData1.Records[4].Id)
	assert.Equal(t, "P05", resData1.Records[4].Pid1)
	assert.Equal(t, "P09", resData1.Records[4].Pid2)
	assert.Equal(t, relFather, resData1.Records[4].Type)

	assert.Equal(t, int64(23), resData1.Records[5].Id)
	assert.Equal(t, "P05", resData1.Records[5].Pid1)
	assert.Equal(t, "P10", resData1.Records[5].Pid2)
	assert.Equal(t, relFather, resData1.Records[5].Type)

	assert.Equal(t, int64(25), resData1.Records[6].Id)
	assert.Equal(t, "P05", resData1.Records[6].Pid1)
	assert.Equal(t, "P11", resData1.Records[6].Pid2)
	assert.Equal(t, relFather, resData1.Records[6].Type)

	assert.Equal(t, int64(29), resData1.Records[7].Id)
	assert.Equal(t, "P05", resData1.Records[7].Pid1)
	assert.Equal(t, "P13", resData1.Records[7].Pid2)
	assert.Equal(t, relFather, resData1.Records[7].Type)

	assert.Equal(t, int64(31), resData1.Records[8].Id)
	assert.Equal(t, "P05", resData1.Records[8].Pid1)
	assert.Equal(t, "P14", resData1.Records[8].Pid2)
	assert.Equal(t, relFather, resData1.Records[8].Type)

	assert.Equal(t, int64(33), resData1.Records[9].Id)
	assert.Equal(t, "P05", resData1.Records[9].Pid1)
	assert.Equal(t, "P15", resData1.Records[9].Pid2)
	assert.Equal(t, relFather, resData1.Records[9].Type)

	assert.Equal(
		t, "http://example.com/people/P05/relations?limit=10&page=1", resData1.Pagination.NextUrl)
	assert.Empty(t, resData1.Pagination.PrevUrl)

	// Case 2: Successful retrieval of the second page

	res = testMakeRequest(router, "GET", resData1.Pagination.NextUrl, nil)

	assert.Equal(t, http.StatusOK, res.Code)

	resData2 := testRelationListRes(t, res)

	assert.Len(t, resData2.Records, 2)

	assert.Equal(t, int64(35), resData2.Records[0].Id)
	assert.Equal(t, "P05", resData2.Records[0].Pid1)
	assert.Equal(t, "P16", resData2.Records[0].Pid2)
	assert.Equal(t, relFather, resData2.Records[0].Type)

	assert.Equal(t, int64(37), resData2.Records[1].Id)
	assert.Equal(t, "P05", resData2.Records[1].Pid1)
	assert.Equal(t, "P17", resData2.Records[1].Pid2)
	assert.Equal(t, relFather, resData2.Records[1].Type)

	assert.Empty(t, resData2.Pagination.NextUrl)
	assert.Equal(
		t, "http://example.com/people/P05/relations?limit=10&page=0", resData2.Pagination.PrevUrl)

	// Case 3: Invalid source person id format

	res = testMakeRequest(router, "GET", "/people/%231/relations", nil)

	assert.Equal(t, http.StatusBadRequest, res.Code)

	resData3 := testErrorRes(t, res)

	assert.Equal(t, uriErrorMsg, resData3.Message)

	// Case 4: Too small page size

	res = testMakeRequest(router, "GET", "/people/1/relations?limit=5", nil)

	assert.Equal(t, http.StatusBadRequest, res.Code)

	resData4 := testErrorRes(t, res)

	assert.Equal(t, queryErrorMsg, resData4.Message)

	// Case 5: Non-existing person

	res = testMakeRequest(router, "GET", "/people/P00/relations", nil)

	assert.Equal(t, http.StatusNotFound, res.Code)

	resData5 := testErrorRes(t, res)

	assert.Equal(t, "Unknown person id", resData5.Message)
}
