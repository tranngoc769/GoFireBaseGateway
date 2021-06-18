package util

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestGetHeaderOk(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ginCtx, _ := gin.CreateTestContext(w)
	ginCtx.Request, _ = http.NewRequest("GET", "/", nil)
	ginCtx.Request.Header.Set("id", "foo")
	_, value := GetHeader(ginCtx, "id")
	require.Equal(t, true, value)
}

func TestGetHeaderFail(t *testing.T) {
	testCases := []struct {
		name          string
		headers       map[string]string
		checkResponse func(value bool)
	}{
		{
			name: "Fail - Missing",
			headers: map[string]string{
				"xid": "foo",
			},
			checkResponse: func(value bool) {
				require.Equal(t, false, value)
			},
		},
		{
			name: "Fail - Empty",
			headers: map[string]string{
				"id": "",
			},
			checkResponse: func(value bool) {
				require.Equal(t, false, value)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			ginCtx, _ := gin.CreateTestContext(w)
			ginCtx.Request, _ = http.NewRequest("GET", "/", nil)
			for key, value := range tc.headers {
				ginCtx.Request.Header.Set(key, value)
			}
			_, value := GetHeader(ginCtx, "id")
			tc.checkResponse(value)
		})
	}
}
