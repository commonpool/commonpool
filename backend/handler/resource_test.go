package handler

import (
	"encoding/json"
	"fmt"
	"github.com/commonpool/backend/errors"
	"github.com/commonpool/backend/model"
	"github.com/commonpool/backend/router"
	"github.com/commonpool/backend/web"
	"github.com/labstack/echo/v4"
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
	key := model.NewResourceKey()
	resource := model.NewResource(
		key,
		model.Offer,
		"author",
		"a superb summary",
		"Description",
		1,
		2,
	)
	assert.NoError(t, rs.Create(&resource))

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
			"type":0
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
	assert.Equal(t, model.Offer, res.Resource.Type)
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
			"type":0
		}
	}`)
	err := h.CreateResource(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	res := errors.ErrorResponse{}
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))
	assert.Equal(t, errors.ErrCreateResourceCannotBind, res.Message)
	assert.Equal(t, errors.ErrCreateResourceCannotBindCode, res.Code)
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
			"type":0
		}
	}`)
	err := h.CreateResource(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	res := errors.ErrorResponse{}
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))
	assert.Equal(t, errors.ErrSummaryEmptyOrNull, res.Message)
	assert.Equal(t, errors.ErrSummaryEmptyOrNullCode, res.Code)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	fmt.Println(string(rec.Body.Bytes()))
}

// TestCreateResourceEmptyName400
// Should return 400 if summary is empty
func TestCreateResourceLongSummary400(t *testing.T) {
	tearDown()
	setup()

	var a = ""
	for i := 0; i < 101; i++ {
		a = a + "A"
	}

	rec, c := newCreateResourceRequest(`
	{
		"resource": {
			"summary":"` + a + `",
			"description":"description",
			"type":0
		}
	}`)
	err := h.CreateResource(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	res := errors.ErrorResponse{}
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))
	assert.Equal(t, errors.ErrSummaryTooLong, res.Message)
	assert.Equal(t, errors.ErrSummaryTooLongCode, res.Code)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	fmt.Println(string(rec.Body.Bytes()))
}

// TestGetResource
// Should be able to retrieve a resource
func TestGetResource(t *testing.T) {
	tearDown()
	setup()

	// Creating the resource
	key := model.NewResourceKey()
	resource := model.NewResource(
		key,
		model.Offer,
		"author",
		"Summary",
		"Description",
		1,
		2)
	assert.NoError(t, rs.Create(&resource))

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
	assert.Equal(t, model.Offer, res.Resource.Type)
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
	assert.Equal(t, errors.ErrUuidParseError, res.Message)
	assert.Equal(t, errors.ErrUuidParseErrorCode, res.Code)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

// TestGetUnknownResource404
// Should throw 404 when getting unknown resource
func TestGetUnknownResource404(t *testing.T) {
	tearDown()
	setup()

	// Setting up the request
	key := model.NewResourceKey()
	rec, c := newGetResourceRequest(key.String())

	// Error assertions
	assert.NoError(t, h.GetResource(c))
	assert.Equal(t, http.StatusNotFound, rec.Code)

	res := errors.ErrorResponse{}
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))
	assert.Equal(t, fmt.Sprintf(errors.ErrResourceNotFoundMsg, key.String()), res.Message)
	assert.Equal(t, errors.ErrResourceNotFoundCode, res.Code)
	assert.Equal(t, http.StatusNotFound, res.StatusCode)

}

// TestUpdateResource
// Should be able to update a resource summary, description and type
func TestUpdateResource(t *testing.T) {
	tearDown()
	setup()

	// Creating the resource
	key := model.NewResourceKey()
	resource := model.NewResource(
		key,
		model.Offer,
		"author",
		"Summary",
		"Description",
		1,
		2)
	assert.NoError(t, rs.Create(&resource))

	// Setting up the request
	js := `
	{
		"resource":{
			"summary":"new summary",
			"description":"new description",
			"type":1,
			"timeSensitivity":20,
			"exchangeValue":30,
			"necessityLevel":40,
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
	assert.Equal(t, model.Request, res.Resource.Type)
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

func newCreateResourceRequest(js string) (*httptest.ResponseRecorder, echo.Context) {
	_, _, rec, c := newRequest(echo.POST, "/api/resources", &js)
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
