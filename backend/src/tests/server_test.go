package tests

import (
	"backend/src/api"
	"backend/src/dependencies"
	"backend/src/endpoints"
	"database/sql"
	"github.com/stretchr/testify/require"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type TestContext struct {
	db      *sql.DB
	handler http.Handler
	pgMock  sqlmock.Sqlmock
}

func Prepare(t *testing.T) TestContext {
	var context TestContext
	var err error
	context.db, context.pgMock, err = sqlmock.New()
	if err != nil {
		t.Errorf("failed to set up SQL mock")
	}

	handlers := endpoints.Handlers{
		Deps: dependencies.Collection{
			Db:     context.db,
			Config: nil,
		},
	}
	context.handler = api.Handler(api.NewStrictHandler(handlers, nil))
	return context
}

func Finalize(c TestContext) {
	c.db.Close()
}

func TestGetItems(t *testing.T) {
	t.Run("simple test api response", func(t *testing.T) {
		// set up mock
		c := Prepare(t)

		c.pgMock.ExpectPrepare("SELECT id, name, price FROM items").ExpectQuery().
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "name", "price"}).
					AddRow("id", "name", 1))

		request, _ := http.NewRequest(http.MethodGet, "/items", nil)
		response := httptest.NewRecorder()
		c.handler.ServeHTTP(response, request)

		expectedResponse := `{"items":[{"id":"id","name":"name","price":1}]}`

		require.EqualValues(t, 200, response.Code)
		require.JSONEq(t, expectedResponse, response.Body.String())
		Finalize(c)
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
		request, _ := http.NewRequest(http.MethodPost, "/items", strings.NewReader(requestBody))
		response := httptest.NewRecorder()
		c.handler.ServeHTTP(response, request)

		expectedResponse := `{"id":"id","name":"name","price":1}`

		require.EqualValues(t, 200, response.Code)
		require.JSONEq(t, expectedResponse, response.Body.String())
		Finalize(c)
	})
}
