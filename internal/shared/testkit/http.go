package testkit

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

func NewGinEngine(routes func(r *gin.Engine)) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	routes(r)
	return r
}

func NewRequest(method, path string, body any) *http.Request {
	bBody, _ := json.Marshal(body)
	req := httptest.NewRequest(method, path, bytes.NewReader(bBody))
	req.Header.Set("Content-Type", "application/json")
	return req
}

func ExecRequest(r *gin.Engine, req *http.Request) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
