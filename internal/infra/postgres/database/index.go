package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/newrelic/go-agent/v3/integrations/nrpq"
	log "github.com/sirupsen/logrus"
)

var Started bool

type Postgres struct {
	*sql.DB
}

func NewPostgresDatabase(ctx context.Context, host string) (*Postgres, error) {
	db, err := sql.Open("postgres", host)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	Started = true

	return &Postgres{DB: db}, nil
}

func MakeMigration(ctx context.Context, database *Postgres, migrationName string) {
	driver, err := postgres.WithInstance(database.DB, &postgres.Config{
		MigrationsTable: fmt.Sprintf("%s-database-schema-migrations", migrationName),
	})
	if err != nil {
		log.Panic(err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://./db/migrations/postgres",
		migrationName,
		driver,
	)
	if err != nil {
		log.Panicf("Error connecting migrator %+v", err)
	}
	if err := m.Up(); err != nil {
		if string(err.Error()) != "no change" {
			log.Panicf("Error making the migration -> %+v", err)
		}
	}
}

func (p *Postgres) Close() {
	_ = p.DB.Close()
}
