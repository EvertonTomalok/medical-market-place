package app

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"

	"github.com/EvertonTomalok/marketplace-health/internal/app/usecases"
	"github.com/EvertonTomalok/marketplace-health/internal/infra/postgres/database"
	postgresDB "github.com/EvertonTomalok/marketplace-health/internal/infra/postgres/database"
	"github.com/EvertonTomalok/marketplace-health/internal/infra/postgres/repositories"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

const LocalHost string = "0.0.0.0"
const DefaultPort string = "8000"
const DefaultDatabasePort string = "5432"
const DefaultPostgresHostTemplate string = "postgres://postgres:postgres@%s:%s/%s?sslmode=disable"

var DefaultDatabase string = fmt.Sprintf(
	DefaultPostgresHostTemplate,
	LocalHost,
	DefaultDatabasePort,
	"postgres",
)

var connections = make([]*postgresDB.Postgres, 0)

type Config struct {
	App struct {
		Port     string
		Host     string
		LogLevel string
		Database struct {
			Host               string
			Port               string
			Name               string
			ConnMaxLifetime    int
			MaxOpenConnections int
			MaxIdleConnections int
			MakeMigration      bool
		}
	}
}

func InjectRepositories(ctx context.Context, cfg Config, logger *logrus.Logger) {
	usecases.DocumentWorkerRepository = repositories.NewDocumentWorkerRepository(nil, logger)
	usecases.FacilityRequirementsRepository = repositories.NewFacilityRequirementRepository(nil, logger)
	usecases.FacilityRepository = repositories.NewFacilityRepository(nil, logger)
	usecases.ShiftRepository = repositories.NewShiftRepository(nil, logger)
	usecases.WorkerRepository = repositories.NewWorkerRepository(nil, logger)
}

func InitDB(ctx context.Context, config *Config, logger *logrus.Logger) {
	conn, err := database.NewPostgresDatabase(ctx, config.App.Database.Host)
	if err != nil {
		logger.Panicf("Error initializing DB. Err: %+v", err)
	}
	repositories.Conn = conn

	InjectRepositories(ctx, *config, logger)

	maxDelimiter := 12
	if len(config.App.Database.Host) < 12 {
		maxDelimiter = 5
	}
	logger.Infof(
		"Database connection is ready at %s***:%s/%s",
		config.App.Database.Host[0:maxDelimiter],
		config.App.Database.Port,
		config.App.Database.Name,
	)
}

func Configure(ctx context.Context) Config {
	_ = godotenv.Load()

	viper.SetDefault("App.Host", LocalHost)
	viper.SetDefault("App.Port", DefaultPort)
	viper.SetDefault("App.LogLevel", "DEBUG")
	viper.SetDefault("App.Database.Host", DefaultDatabase)
	viper.SetDefault("App.Database.Port", "5432")
	viper.SetDefault("App.Database.Name", "postgres")
	viper.SetDefault("App.Database.ConnMaxLifetime", 15*time.Minute)
	viper.SetDefault("App.Database.MaxOpenConnections", 100)
	viper.SetDefault("App.Database.MaxIdleConnections", 50)
	viper.SetDefault("App.Database.MakeMigration", false)
	viper.AutomaticEnv()

	var cfg Config

	if err := viper.Unmarshal(&cfg); err != nil {
		log.Panicf("unmarshaling config: %+v", err)
	}
	log.Print("configuration loaded")

	return cfg
}

func CloseConnections(ctx context.Context) {
	repositories.Conn.Close()

	for _, db := range connections {
		db.Close()
	}
	connections = make([]*postgresDB.Postgres, 0)
}
