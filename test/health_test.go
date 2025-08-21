package test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/gabrielolivrp/pastebin-api/internal/module/health"
	"github.com/gabrielolivrp/pastebin-api/pkg/cache"
	"github.com/gabrielolivrp/pastebin-api/pkg/config"
	"github.com/gabrielolivrp/pastebin-api/pkg/database"
	"github.com/gabrielolivrp/pastebin-api/pkg/logging"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewMockFailingDB() database.Client {
	return &mockFailingDB{}
}

type mockFailingDB struct{}

func (d *mockFailingDB) DB() *gorm.DB {
	return &gorm.DB{}
}

func (d *mockFailingDB) Ping() error {
	return errors.New("Database connection failed")
}

func NewMockFailingCache() cache.Client {
	return &mockFailingCache{}
}

type mockFailingCache struct{}

func (d *mockFailingCache) Ping() error {
	return errors.New("Database connection failed")
}

func (d *mockFailingCache) Get(ctx context.Context, key string) (string, error) {
	return "", errors.New("Cache connection failed")
}

func (d *mockFailingCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return errors.New("Cache connection failed")
}

func (d *mockFailingCache) Has(ctx context.Context, key string) (bool, error) {
	return false, errors.New("Cache connection failed")
}

func (s *TestSuite) TestHealth() {
	s.Run("Health Check OK", func() {
		logger, err := logging.NewLogger(config.Test)
		s.Require().NoError(err)

		gin.SetMode(gin.TestMode)
		router := gin.Default()
		health.RegisterRoutes(router.Group(""), logger, s.dbClient, s.cacheClient)

		req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		s.Equal(http.StatusOK, w.Code)

		var resp map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &resp)
		s.Require().NoError(err)
		s.Equal("ok", resp["data"].(map[string]interface{})["status"])
	})

	s.Run("Health Check Fail", func() {
		logger, err := logging.NewLogger(config.Test)
		s.Require().NoError(err)

		gin.SetMode(gin.TestMode)
		router := gin.Default()
		health.RegisterRoutes(router.Group(""), logger, NewMockFailingDB(), s.cacheClient)

		req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		s.Equal(http.StatusServiceUnavailable, w.Code)

		var resp map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &resp)
		s.Require().NoError(err)

		apiError := resp["error"].(map[string]interface{})
		s.Equal("Service Unavailable", apiError["title"])
		s.Contains(apiError["detail"], "Some resources are unavailable")
	})
}
