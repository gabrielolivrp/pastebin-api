package test

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/gabrielolivrp/pastebin-api/pkg/cache"
	"github.com/gabrielolivrp/pastebin-api/pkg/config"
	"github.com/gabrielolivrp/pastebin-api/pkg/database"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/modules/redis"
	"github.com/testcontainers/testcontainers-go/wait"
)

type TestSuite struct {
	suite.Suite
	ctx            context.Context
	dbContainer    *postgres.PostgresContainer
	cacheContainer *redis.RedisContainer
	dbClient       database.Client
	cacheClient    cache.Client
	config         *config.Config
}

func (s *TestSuite) SetupSuite() {
	s.ctx = context.Background()
	config, err := config.LoadConfig("../.env.test")
	if err != nil {
		panic(err)
	}
	s.config = config
	s.databaseContainerSetup()
	s.cacheContainerSetup()
}

func (s *TestSuite) TearDownSuite() {
	if s.dbContainer != nil {
		err := s.dbContainer.Terminate(s.ctx)
		s.Require().NoError(err)
	}
}

func (s *TestSuite) SetupTest() {
	s.clearDatabase()
}

func (s *TestSuite) TearDownTest() {
	s.clearDatabase()
}

func (s *TestSuite) clearDatabase() {
	var tables []string
	err := s.dbClient.DB().
		Raw("SELECT tablename FROM pg_tables WHERE schemaname = 'public'").
		Scan(&tables).Error
	s.Require().NoError(err)

	if len(tables) == 0 {
		return
	}

	for _, table := range tables {
		err := s.dbClient.DB().Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table)).Error
		s.Require().NoError(err)
	}
}

func (s *TestSuite) databaseContainerSetup() {
	pgContainer, err := postgres.Run(s.ctx,
		"postgres:17.6",
		postgres.WithDatabase(s.config.DB.Database),
		postgres.WithUsername(s.config.DB.Username),
		postgres.WithPassword(s.config.DB.Password),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(1).
				WithStartupTimeout(60*time.Second),
		),
	)
	s.Require().NoError(err)
	s.dbContainer = pgContainer

	host, err := pgContainer.Host(s.ctx)
	s.Require().NoError(err)
	if host == "localhost" || host == "::1" {
		host = "127.0.0.1"
	}

	port, err := pgContainer.MappedPort(s.ctx, "5432/tcp")
	s.Require().NoError(err)

	time.Sleep(2 * time.Second)

	dsn, err := pgContainer.ConnectionString(s.ctx, "sslmode=disable")
	s.Require().NoError(err)

	migrationsDir := filepath.Join("..", "migrations")
	m, err := migrate.New("file://"+migrationsDir, dsn)
	s.Require().NoError(err)

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		s.Require().NoError(err)
	}

	client, err := database.NewClient(database.ClientConfig{
		Host:     host,
		Port:     port.Port(),
		Username: s.config.DB.Username,
		Password: s.config.DB.Password,
		Database: s.config.DB.Database,
		SSLMode:  "disable",
	})
	s.Require().NoError(err)
	s.dbClient = client
}

func (s *TestSuite) cacheContainerSetup() {
	rdContainer, err := redis.Run(s.ctx,
		"redis:8.2.1",
		redis.WithSnapshotting(10, 1),
		redis.WithLogLevel(redis.LogLevelVerbose),
		testcontainers.WithWaitStrategy(
			wait.ForLog("Ready to accept connections").
				WithOccurrence(1).
				WithStartupTimeout(60*time.Second),
		),
	)

	s.Require().NoError(err)
	s.cacheContainer = rdContainer

	host, err := rdContainer.Host(s.ctx)
	s.Require().NoError(err)
	if host == "localhost" || host == "::1" {
		host = "127.0.0.1"
	}

	port, err := rdContainer.MappedPort(s.ctx, "6379/tcp")
	s.Require().NoError(err)

	client, err := cache.NewClient(cache.ClientConfig{
		Host:     host,
		Port:     port.Port(),
		Password: s.config.Cache.Password,
	})
	s.Require().NoError(err)
	s.cacheClient = client
}

func Test(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
