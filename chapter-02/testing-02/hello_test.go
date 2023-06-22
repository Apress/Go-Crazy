package main

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAPIRouteNew(t *testing.T) {
	router := router()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/new/janus", nil)
	router.ServeHTTP(w, req)

	var img ImageAnswer
	json.Unmarshal(w.Body.Bytes(), &img)
	assert.NotEmpty(t, img)

}

func TestAPIRouteLoad(t *testing.T) {
	router := router()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/new/janus", nil)
	router.ServeHTTP(w, req)

	var img ImageAnswer
	json.Unmarshal(w.Body.Bytes(), &img)

	req, _ = http.NewRequest("GET", "/new/load/"+img.Id, nil)
	router.ServeHTTP(w, req)

	fmt.Printf("%+v\n", img)
	assert.Equal(t, "image/jpg", w.Header().Get("Content-Type"))

	time.Sleep(3 * time.Second)

	req, _ = http.NewRequest("GET", "/new/load/"+img.Id, nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, "image/png", w.Header().Get("Content-Type"))
}
