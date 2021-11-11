package main

import (
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
)

func TestPingRoute(t *testing.T) {
    r := gin.Default()
    customizeouter(r)

    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/ping", nil)
    r.ServeHTTP(w, req)

    assert.Equal(t, http.StatusOK, w.Code)
    assert.Equal(t, "{\"message\":\"pong\"}", w.Body.String())
}

func TestUpdateService(t *testing.T) {
    r := gin.Default()
    customizeouter(r)

    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/update?device_platform=unknown", nil)
    r.ServeHTTP(w, req)

    assert.Equal(t, http.StatusOK, w.Code)
    assert.Equal(t, "{\"download_url\":\"nil\",\"update_version_code\":\"\",\"md5\":\"\",\"title\":\"\",\"update_tips\":\"\"}", w.Body.String())
}


