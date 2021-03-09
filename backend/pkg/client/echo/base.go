package echo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/commonpool/backend/pkg/client"
	"github.com/commonpool/backend/pkg/exceptions"
	"github.com/labstack/echo/v4"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
)

type Client struct {
	echo           *echo.Echo
	authentication client.Authentication
}

func (e *Client) newRequest(method, target string, body io.Reader) *http.Request {
	return httptest.NewRequest(method, target, body)
}

func (e *Client) read(value interface{}) (io.Reader, error) {
	if value == nil {
		return nil, nil
	}
	bts, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(bts), nil
}

func (e *Client) assertResponseCode(statusCode int, resp *http.Response, body []byte) error {
	if resp.StatusCode != statusCode {
		var errResponse exceptions.WebServiceException
		err := json.Unmarshal(body, &errResponse)
		if err != nil {
			return fmt.Errorf("could not unmarshal error response")
		}
		if errResponse.Code == "" {
			return fmt.Errorf("error occured on server side, but could not decode error: %s", string(body))
		}
		return errResponse
	}
	return nil
}

func (e *Client) createRequest(ctx context.Context, method, target string, postBody interface{}) (*http.Request, error) {
	postBodyReader, err := e.read(postBody)
	if err != nil {
		return nil, err
	}
	httpRequest := e.newRequest(method, target, postBodyReader)
	httpRequest = httpRequest.WithContext(ctx)
	client.JsonContentType(httpRequest)
	if err := e.authentication.Apply(httpRequest); err != nil {
		return nil, err
	}
	return httpRequest, nil
}

func (e *Client) readResponse(response *http.Response, assertResponseCode int, out interface{}) error {
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	if err := e.assertResponseCode(assertResponseCode, response, body); err != nil {
		return err
	}
	if out == nil {
		return nil
	}
	return json.Unmarshal(body, out)
}

func (e *Client) doRequest(request *http.Request) *http.Response {
	recorder := httptest.NewRecorder()
	e.echo.ServeHTTP(recorder, request)
	resp := recorder.Result()
	return resp
}

func (e *Client) do(ctx context.Context, method, target string, assertResponseCode int, payload interface{}, output interface{}) error {
	httpRequest, err := e.createRequest(ctx, method, target, payload)
	if err != nil {
		return err
	}
	resp := e.doRequest(httpRequest)
	if err := e.readResponse(resp, assertResponseCode, output); err != nil {
		return err
	}
	return nil
}

func NewEchoClient(echo *echo.Echo, authentication client.Authentication) *Client {
	return &Client{
		echo:           echo,
		authentication: authentication,
	}
}
