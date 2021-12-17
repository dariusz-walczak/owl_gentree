package main

import (
	"github.com/stretchr/testify/assert"
	"net/http"
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
func TestCreatePersonRequestSuccess1(t *testing.T) {
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

	resData := testErrorRes(t, res)

	assert.Equal(t, "Relation created", resData.Message)
	assert.Len(t, relations, 1)

	// Case 2: General mother relation

	iitRelation = testIitRelationJson{
		Pid1: "F1P2",
		Pid2: "F1P3",
		Type: "mother"}

	res = testMakeRequest(router, "POST", "/relations", testJsonBody(t, iitRelation))

	assert.Equal(t, http.StatusCreated, res.Code)

	resData = testErrorRes(t, res)

	assert.Equal(t, "Relation created", resData.Message)
	assert.Len(t, relations, 2)

	// Case 3: General husband relation

	iitRelation = testIitRelationJson{
		Pid1: "F1P1",
		Pid2: "F1P2",
		Type: "husband"}

	res = testMakeRequest(router, "POST", "/relations", testJsonBody(t, iitRelation))

	assert.Equal(t, http.StatusCreated, res.Code)

	resData = testErrorRes(t, res)

	assert.Equal(t, "Relation created", resData.Message)
	assert.Len(t, relations, 3)

	// Case 4: Person specific father relation

	itRelation := testItRelationJson{
		Pid:  "F2P3",
		Type: "father"}

	res = testMakeRequest(router, "POST", "/people/F2P1/relations", testJsonBody(t, itRelation))

	assert.Equal(t, http.StatusCreated, res.Code)

	resData = testErrorRes(t, res)

	assert.Equal(t, "Relation created", resData.Message)
	assert.Len(t, relations, 4)

	// Case 5: Person specific mother relation

	itRelation = testItRelationJson{
		Pid:  "F2P3",
		Type: "mother"}

	res = testMakeRequest(router, "POST", "/people/F2P2/relations", testJsonBody(t, itRelation))

	assert.Equal(t, http.StatusCreated, res.Code)

	resData = testErrorRes(t, res)

	assert.Equal(t, "Relation created", resData.Message)
	assert.Len(t, relations, 5)

	// Case 6: Person specific husband relation

	itRelation = testItRelationJson{
		Pid:  "F2P2",
		Type: "husband"}

	res = testMakeRequest(router, "POST", "/people/F2P1/relations", testJsonBody(t, itRelation))

	assert.Equal(t, http.StatusCreated, res.Code)

	resData = testErrorRes(t, res)

	assert.Equal(t, "Relation created", resData.Message)
	assert.Len(t, relations, 6)
}
