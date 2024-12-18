//go:build integration

package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"testing"

	_ "github.com/jackc/pgx/stdlib" // postgresql provider
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const (
	dbDefault           = "host=localhost port=10905 user=user password=qwerty dbname=otus_calendar sslmode=disable"
	dbDriverDefault     = "sql"
	restDefault         = "http://127.0.0.1:8090"
	amqpURIDefault      = "amqp://guest:guest@127.0.0.1:5672/local"
	exchangeTypeDefault = "direct"
	exchangeNameDefault = "event-exchange"
	routingKeyDefault   = "event-key"
)

type SuiteIntegrationTest struct {
	suite.Suite
	db           *sqlx.DB
	dbDriver     string
	rest         string
	amqpURI      string
	exchangeType string
	exchangeName string
	routingKey   string
}

func (s *SuiteIntegrationTest) needCheckSql() bool {
	return s.dbDriver == "sql"
}

func TestSuit(t *testing.T) {
	suite.Run(t, &SuiteIntegrationTest{})
}

func (s *SuiteIntegrationTest) SetupSuite() {
	dbDSN, ok := os.LookupEnv("TEST_DB_DSN")
	if !ok {
		dbDSN = dbDefault
	}
	db, err := sqlx.Open("pgx", dbDSN)
	require.NoError(s.T(), err)
	err = db.Ping()
	require.NoError(s.T(), err)
	s.db = db

	rest, ok := os.LookupEnv("TEST_REST")
	if !ok {
		rest = restDefault
	}
	s.rest = rest

	dbDriver, ok := os.LookupEnv("TEST_DB_DRIVER")
	if !ok {
		dbDriver = dbDriverDefault
	}
	s.dbDriver = dbDriver

	amqpURI, ok := os.LookupEnv("TEST_DB_AMQP_URI")
	if !ok {
		amqpURI = amqpURIDefault
	}
	s.amqpURI = amqpURI

	exchangeType, ok := os.LookupEnv("TEST_DB_EXCHANGE_TYPE")
	if !ok {
		exchangeType = exchangeTypeDefault
	}
	s.exchangeType = exchangeType

	exchangeName, ok := os.LookupEnv("TEST_DB_EXCHANGE_NAME")
	if !ok {
		exchangeName = exchangeNameDefault
	}
	s.exchangeName = exchangeName

	routingKey, ok := os.LookupEnv("TEST_DB_ROUTING_KEY")
	if !ok {
		routingKey = routingKeyDefault
	}
	s.routingKey = routingKey
}

func (s *SuiteIntegrationTest) TearDownTest() {
	_, err := s.db.Exec("truncate event cascade")
	require.NoError(s.T(), err)
}

func (s *SuiteIntegrationTest) httpRequest(method string, url string, body any) *http.Response {
	bodyJson, err := json.Marshal(body)
	require.NoError(s.T(), err)

	request, err := http.NewRequest(method, s.rest+url, bytes.NewBuffer(bodyJson))
	require.NoError(s.T(), err)

	request.Header.Add("Content-Type", "application/json; charset=utf-8")
	request.Header.Add("Accept", "application/json")

	response, err := http.DefaultClient.Do(request)
	require.NoError(s.T(), err)

	return response
}
