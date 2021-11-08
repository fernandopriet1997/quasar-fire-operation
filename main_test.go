package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPostTopSecret(t *testing.T) {
	router := setupRouter()

	jsonParam := `
	{
		"satellites": [
			{
				"name": "kenobi", 
				"distance": 100, 
				"message": ["Este", "", "", "mensaje", ""]
			},
			{
				"name": "skywalker", 
				"distance": 115.5, 
				"message": ["", "es", "", "", "secreto"]
			},
			{
				"name": "sato",
				"distance": 142.7, 
				"message": ["Este", "", "un", "", ""]
			}
		]
	}`

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("POST", "/topsecret", strings.NewReader(string(jsonParam)))
	req.Header.Add("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}
func TestPostTopSecretSplit(t *testing.T) {
	router := setupRouter()

	jsonParam := `
	{
		"distance": 100, 
		"message": ["Este", "", "", "mensaje", ""]
	}`

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("POST", "/topsecret_split/kenobi", strings.NewReader(string(jsonParam)))
	req.Header.Add("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}
func TetsPostSetPosition(t *testing.T) {
	router := setupRouter()

	jsonParam := `
	{
		"x": 100, 
		"y": -100
	}`

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("POST", "/config/kenobi", strings.NewReader(string(jsonParam)))
	req.Header.Add("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}
func TetsResetDefault(t *testing.T) {

	router := setupRouter()

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("DELETE", "/topsecret_reset", nil)
	req.Header.Add("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}
