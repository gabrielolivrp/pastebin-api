package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/gabrielolivrp/pastebin-api/internal/module/paste"
	"github.com/gabrielolivrp/pastebin-api/pkg/config"
	"github.com/gabrielolivrp/pastebin-api/pkg/logging"
	"github.com/gin-gonic/gin"
)

func pasteTestSetup(s *TestSuite) *gin.Engine {
	logger, err := logging.NewLogger(config.Test)
	s.Require().NoError(err)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	paste.RegisterRoutes(router.Group(""), logger, s.dbClient, s.cacheClient)

	return router
}

func (s *TestSuite) TestPaste() {
	s.Run("Create Paste", func() {
		router := pasteTestSetup(s)

		body := map[string]interface{}{
			"title":   "Title Test",
			"content": "# Hello World",
			"lang":    "md",
		}
		bodyBytes, err := json.Marshal(body)
		s.Require().NoError(err)

		req := httptest.NewRequest(http.MethodPost, "/pastes", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		s.Equal(http.StatusCreated, w.Code)

		var resp map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &resp)
		s.Require().NoError(err)

		data := resp["data"].(map[string]interface{})
		s.NotEmpty(data["id"])
		s.NotEmpty(data["created_at"])
		s.NotEmpty(data["expires_at"])
		s.Equal(body["title"], data["title"])
		s.Equal(body["content"], data["content"])
		s.Equal(body["lang"], data["lang"])
	})

	s.Run("Create Paste Invalid Body", func() {
		router := pasteTestSetup(s)

		body := map[string]interface{}{
			"title":   "Title Test",
			"content": "# Hello World",
		}
		bodyBytes, err := json.Marshal(body)
		s.Require().NoError(err)

		req := httptest.NewRequest(http.MethodPost, "/pastes", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)
		s.Equal(http.StatusBadRequest, w.Code)
		var resp map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &resp)
		s.Require().NoError(err)

		apiError := resp["error"].(map[string]interface{})
		s.Equal(apiError["title"], "Validation failed")
		s.Contains(apiError["errors"],
			map[string]interface{}{
				"field":   "Lang",
				"code":    "required",
				"message": "Key: 'CreatePasteRequest.Lang' Error:Field validation for 'Lang' failed on the 'required' tag",
			},
		)
	})

	s.Run("Get Paste By ID", func() {
		router := pasteTestSetup(s)

		body := map[string]interface{}{
			"title":   "Title Test",
			"content": "# Hello World",
			"lang":    "md",
		}
		bodyBytes, err := json.Marshal(body)
		s.Require().NoError(err)

		req := httptest.NewRequest(http.MethodPost, "/pastes", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		var createResp map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &createResp)
		s.Require().NoError(err)

		pasteID := createResp["data"].(map[string]interface{})["id"]

		req = httptest.NewRequest(http.MethodGet, "/pastes/"+pasteID.(string), nil)
		w = httptest.NewRecorder()

		router.ServeHTTP(w, req)

		s.Equal(http.StatusOK, w.Code)

		var getResp map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &getResp)
		s.Require().NoError(err)

		data := getResp["data"].(map[string]interface{})
		s.Equal(pasteID, data["id"])
		s.Equal(body["title"], data["title"])
		s.Equal(body["content"], data["content"])
		s.Equal(body["lang"], data["lang"])
	})
}
