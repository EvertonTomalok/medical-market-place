package mocks

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/EvertonTomalok/marketplace-health/internal/infra/postgres/database"
)

func NewMockEqualMatcher(t *testing.T) (*database.Postgres, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return &database.Postgres{DB: db}, mock
}

func NewMockRegexMatcher(t *testing.T) (*database.Postgres, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return &database.Postgres{DB: db}, mock
}

func NewMockDefault(t *testing.T) (*database.Postgres, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	return &database.Postgres{DB: db}, mock
}
