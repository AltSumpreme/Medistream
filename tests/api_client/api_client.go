package apiclient

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

type TestClient struct {
	Router *gin.Engine
}

func NewTestClient(router *gin.Engine) *TestClient {
	return &TestClient{
		Router: router,
	}
}

func (c *TestClient) PerformRequest(method string, path string, body interface{}, headers map[string]string) *httptest.ResponseRecorder {
	var reqBody io.Reader

	if body != nil {
		jsonBody, _ := json.Marshal(body)
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req := httptest.NewRequest(method, path, reqBody)
	req.Header.Set("Content-Type", "application/json")

	// Add custom headers if provided
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	w := httptest.NewRecorder()
	c.Router.ServeHTTP(w, req)
	return w
}

// Convenience Wrapper for http requests

func (c *TestClient) Get(path string, headers map[string]string) *httptest.ResponseRecorder {
	return c.PerformRequest("GET", path, nil, headers)
}

func (c *TestClient) Post(path string, body interface{}, headers map[string]string) *httptest.ResponseRecorder {
	return c.PerformRequest("POST", path, body, headers)
}

func (c *TestClient) Put(path string, body interface{}, headers map[string]string) *httptest.ResponseRecorder {
	return c.PerformRequest("PUT", path, body, headers)
}

func (c *TestClient) Delete(path string, headers map[string]string) *httptest.ResponseRecorder {
	return c.PerformRequest("DELETE", path, nil, headers)
}
