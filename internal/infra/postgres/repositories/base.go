package repositories

import (
	"github.com/EvertonTomalok/marketplace-health/internal/infra/postgres/database"
	"github.com/sirupsen/logrus"
)

type BaseRepository struct {
	db     *database.Postgres //lint:ignore U1000 Ignore unused function base struct
	logger *logrus.Logger     //lint:ignore U1000 Ignore unused function base struct
	path   string             //lint:ignore U1000 Ignore unused function base struct
}

func NewBaseRepository(db *database.Postgres, logger *logrus.Logger, path string) BaseRepository {
	return BaseRepository{
		db: db, logger: logger, path: path,
	}
}
