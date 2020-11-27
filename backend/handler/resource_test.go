package handler

import (
	"encoding/json"
	"fmt"
	"github.com/commonpool/backend/errors"
	"github.com/commonpool/backend/model"
	"github.com/commonpool/backend/resource"
	"github.com/commonpool/backend/router"
	"github.com/commonpool/backend/web"
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestCreateResource
// Should be able to create a resource
func TestSearchBySummaryAndType(t *testing.T) {
	tearDown()
	setup()

	// Creating the resource
	key := model.NewResourceKey(uuid.NewV4())
	r := resource.NewResource(
		key,
		resource.ResourceOffer,
		"author",
		"a superb summary",
		"Description",
		1,
		2,
	)
	rq := resource.NewCreateResourceQuery(&r)

	assert.NoError(t, rs.Create(rq).Error)

	_, _, rec, c := newRequest(echo.GET, "/api/resources?take=10&skip=0&query=superb&type=0", nil)
	err := h.SearchResources(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	fmt.Println(string(rec.Body.Bytes()))

	res := web.SearchResourcesResponse{}
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))
	assert.Equal(t, 1, res.TotalCount)
	assert.Equal(t, 0, res.Skip)
	assert.Equal(t, 10, res.Take)
	assert.Equal(t, 1, len(res.Resources))
	assert.Equal(t, key.String(), res.Resources[0].Id)
}

// TestCreateResource
// Should be able to create a resource
func TestCreateResource(t *testing.T) {
	tearDown()
	setup()

	js := `
	{
		"resource": {
			"summary":"summary",
			"description":"description",
			"type":0,
			"valueInHoursFrom":1,
			"valueInHoursTo":3
		}
	}`
	rec, c := newCreateResourceRequest(js)
	err := h.CreateResource(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)

	res := web.CreateResourceResponse{}
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))
	assert.Equal(t, "summary", res.Resource.Summary)
	assert.Equal(t, "description", res.Resource.Description)
	assert.Equal(t, resource.ResourceOffer, res.Resource.Type)
}

// TestCreateResourceInvalid400
// Should throw a 400 bad request when the request payload
// cannot be converted to a CreateResource object
func TestCreateResourceInvalid400(t *testing.T) {
	tearDown()
	setup()

	rec, c := newCreateResourceRequest(`
	{
		"resource": {
			"summary":123,
			"description":456,
			"type":0,
			"valueInHoursFrom":1,
			"valueInHoursTo":3
		}
	}`)
	err := h.CreateResource(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	res := errors.ErrorResponse{}
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))
	assert.Equal(t, "ErrCreateResourceBadRequest", res.Code)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

// TestCreateResourceEmptyName400
// Should return 400 if summary is empty
func TestCreateResourceEmptyName400(t *testing.T) {
	tearDown()
	setup()

	rec, c := newCreateResourceRequest(`
	{ 
		"resource": {
			"summary":"",
			"description":"description",
			"type":0,
			"valueInHoursFrom":1,
			"valueInHoursTo":3
		}
	}`)
	err := h.CreateResource(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// TestCreateResourceLongSummary400
// Should return 400 if summary is too long
func TestCreateResourceLongSummary400(t *testing.T) {
	tearDown()
	setup()

	var a = ""
	for i := 0; i <= 101; i++ {
		a = a + "A"
	}

	rec, c := newCreateResourceRequest(`
	{
		"resource": {
			"summary":"` + a + `",
			"description":"description",
			"type":0,
			"valueInHoursFrom":1,
			"valueInHoursTo":3
		}
	}`)
	err := h.CreateResource(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// TestGetResource
// Should be able to retrieve a resource
func TestGetResource(t *testing.T) {
	tearDown()
	setup()

	// Creating the resource
	key := model.NewResourceKey(uuid.NewV4())
	r := resource.NewResource(
		key,
		resource.ResourceOffer,
		user1.Subject,
		"Summary",
		"Description",
		1,
		2)
	rq := resource.NewCreateResourceQuery(&r)
	assert.NoError(t, rs.Create(rq).Error)

	// Getting the resource
	rec, c := newGetResourceRequest(key.String())
	err := h.GetResource(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	// Validating the payload
	res := web.GetResourceResponse{}
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))
	assert.Equal(t, "Summary", res.Resource.Summary)
	assert.Equal(t, "Description", res.Resource.Description)
	assert.Equal(t, resource.ResourceOffer, res.Resource.Type)
	assert.Equal(t, 1, res.Resource.ValueInHoursFrom)
	assert.Equal(t, 2, res.Resource.ValueInHoursTo)

}

// TestGetResourceBadId400
// Should throw 400 when id is of wrong format
func TestGetResourceBadId400(t *testing.T) {
	tearDown()
	setup()

	// Getting the resource
	rec, c := newGetResourceRequest("bla")
	err := h.GetResource(c)

	// Error assertions
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	res := errors.ErrorResponse{}
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))
	assert.Equal(t, "ErrInvalidResourceKey", res.Code)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

// TestGetUnknownResource404
// Should throw 404 when getting unknown resource
func TestGetUnknownResource404(t *testing.T) {
	tearDown()
	setup()

	// Setting up the request
	key := model.NewResourceKey(uuid.NewV4())
	rec, c := newGetResourceRequest(key.String())

	// Error assertions
	assert.NoError(t, h.GetResource(c))
	assert.Equal(t, http.StatusNotFound, rec.Code)

	res := errors.ErrorResponse{}
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))
	assert.Equal(t, "ErrResourceNotFound", res.Code)
	assert.Equal(t, http.StatusNotFound, res.StatusCode)

}

// TestUpdateResource
// Should be able to update a resource summary, description
func TestUpdateResource(t *testing.T) {
	tearDown()
	setup()

	mockLoggedInAs(user1)

	// Creating the resource
	key := model.NewResourceKey(uuid.NewV4())
	r := resource.NewResource(
		key,
		resource.ResourceOffer,
		user1.Subject,
		"Summary",
		"Description",
		1,
		2)
	rq := resource.NewCreateResourceQuery(&r)
	assert.NoError(t, rs.Create(rq).Error)

	// Setting up the request
	js := `
	{
		"resource":{
			"summary":"new summary",
			"description":"new description",
			"valueInHoursFrom":3,
			"valueInHoursTo":4
		}
	}`
	rec, c := newUpdateResourceRequest(key, js)

	// Updating the resource
	assert.NoError(t, h.UpdateResource(c))
	assert.Equal(t, http.StatusOK, rec.Code)

	// Validate the update
	res := web.GetResourceResponse{}
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))
	assert.Equal(t, "new summary", res.Resource.Summary)
	assert.Equal(t, "new description", res.Resource.Description)
	assert.Equal(t, 3, res.Resource.ValueInHoursFrom)
	assert.Equal(t, 4, res.Resource.ValueInHoursTo)
}

func newGetResourceRequest(key string) (*httptest.ResponseRecorder, echo.Context) {
	path := "/api/resources/:id"
	_, _, rec, c := newRequest(echo.GET, path, nil)
	c.SetPath(path)
	c.SetParamNames("id")
	c.SetParamValues(key)
	return rec, c
}

func newUpdateResourceRequest(key model.ResourceKey, js string) (*httptest.ResponseRecorder, echo.Context) {
	path := "/api/resources/:id"
	_, _, rec, c := newRequest(echo.PUT, path, &js)
	c.SetPath(path)
	c.SetParamNames("id")
	c.SetParamValues(key.String())
	return rec, c
}

func newRequest(method string, target string, reqJson *string) (*echo.Echo, *http.Request, *httptest.ResponseRecorder, echo.Context) {
	e := router.NewRouter()
	var req *http.Request
	if reqJson != nil {
		req = httptest.NewRequest(method, target, strings.NewReader(*reqJson))
	} else {
		req = httptest.NewRequest(method, target, nil)
	}

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	return e, req, rec, c
}

func newCreateResourceRequest(js string) (*httptest.ResponseRecorder, echo.Context) {
	_, _, rec, c := newRequest(echo.POST, "/api/resources", &js)
	return rec, c
}

func createResource(t *testing.T, summary string, description string, resType resource.ResourceType) web.CreateResourceResponse {
	payload := web.CreateResourceRequest{
		Resource: web.CreateResourcePayload{
			Summary:          summary,
			Description:      description,
			Type:             resType,
			ValueInHoursFrom: 1,
			ValueInHoursTo:   3,
			SharedWith:       []web.InputResourceSharing{},
		},
	}
	js, err := json.Marshal(payload)
	assert.NoError(t, err)

	rec, c := newCreateResourceRequest(string(js))
	err = h.CreateResource(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
	resource := web.CreateResourceResponse{}
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resource))
	return resource
}
