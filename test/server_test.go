package tests

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"

	"github.com/jolfzverb/pwstore/internal/dependencies"
	"github.com/jolfzverb/pwstore/internal/endpoints"
)

type TestContext struct {
	db      *sql.DB
	handler http.Handler
	pgMock  sqlmock.Sqlmock
	ctx     context.Context
}

func Prepare(t *testing.T) TestContext {
	t.Helper()
	var testContext TestContext
	var err error
	testContext.db, testContext.pgMock, err = sqlmock.New()
	if err != nil {
		t.Errorf("failed to set up SQL mock")
	}

	testContext.handler = endpoints.GetHandler(
		dependencies.Collection{
			DB:     testContext.db,
			Config: nil,
		})
	testContext.ctx = context.Background()
	return testContext
}

func Finalize(c TestContext) {
	c.db.Close()
}

func TestGetItems(t *testing.T) {
	t.Run("simple test api response", func(t *testing.T) {
		// set up mock
		c := Prepare(t)
		defer Finalize(c)

		c.pgMock.ExpectPrepare("SELECT id, name, price FROM items").ExpectQuery().
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "name", "price"}).
					AddRow("id", "name", 1))

		request, _ := http.NewRequestWithContext(c.ctx, http.MethodGet, "/items", nil)
		response := httptest.NewRecorder()
		c.handler.ServeHTTP(response, request)

		expectedResponse := `{"items":[{"id":"id","name":"name","price":1}]}`

		require.EqualValues(t, 200, response.Code)
		require.JSONEq(t, expectedResponse, response.Body.String())
	})
}

func TestPostItems(t *testing.T) {
	t.Run("simple test api response", func(t *testing.T) {
		c := Prepare(t)

		c.pgMock.ExpectPrepare("INSERT INTO items \\( name, price \\) VALUES \\( \\$1, \\$2 \\) RETURNING id, name, price").
			ExpectQuery().WithArgs("name", float64(1)).
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "name", "price"}).
					AddRow("id", "name", 1))

		requestBody := `{"name":"name","price":1}`
		request, _ := http.NewRequestWithContext(c.ctx, http.MethodPost, "/items", strings.NewReader(requestBody))
		response := httptest.NewRecorder()
		c.handler.ServeHTTP(response, request)

		expectedResponse := `{"id":"id","name":"name","price":1}`

		require.EqualValues(t, 200, response.Code)
		require.JSONEq(t, expectedResponse, response.Body.String())
		Finalize(c)
	})
}
