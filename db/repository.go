package db

import (
	"Quote/config"
	"Quote/db/entity/entities/apertureEntity"
	"Quote/db/entity/entities/clientEntity"
	"Quote/db/entity/entities/colorEntity"
	"Quote/db/entity/entities/fileEntity"
	"Quote/db/entity/entities/glassEntity"
	"Quote/db/entity/entities/profileEntity"
	"Quote/db/entity/entities/quoteEntity"
	"Quote/db/entity/entities/supplierEntity"
	"Quote/db/entity/entities/userEntity"
	quoteLogger "Quote/logger"
	"fmt"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

const connectionString string = "host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=Europe/Rome"

var repoInstance *Repo

type Repo struct {
	DB *gorm.DB
}

func (c *Repo) dbSetup() error {
	err := c.DB.AutoMigrate(
		clientEntity.New(),
		colorEntity.New(),
		glassEntity.New(),
		profileEntity.New(),
		apertureEntity.New(),
		supplierEntity.New(),
		quoteEntity.New(),
		userEntity.New(),
		fileEntity.New(),
	)
	if err != nil {
		return err
	}

	return nil
}

func InitRepo() error {
	repoInstance = new(Repo)

	cfg := config.GetInstance()
	log := quoteLogger.GetInstance()

	connStr := fmt.Sprintf(connectionString, cfg.PostgresHost, cfg.PostgresPort, cfg.PostgresUser,
		cfg.PostgresPassword, cfg.PostgresDatabase, cfg.PostgresSsl)

	var gormConfig *gorm.Config
	if !log.Desugar().Core().Enabled(zap.DebugLevel) {
		gormConfig = &gorm.Config{}
	} else {
		gormConfig = &gorm.Config{Logger: gormLogger.Default.LogMode(gormLogger.Silent)}
	}
	gormConfig.DisableForeignKeyConstraintWhenMigrating = true

	db, err := gorm.Open(postgres.Open(connStr), gormConfig)
	if err != nil {
		return err
	}

	repoInstance.DB = db
	log.Infof("connected to database : %s", cfg.PostgresDatabase)

	if cfg.DBSetup {
		log.Info("migration starts...")

		err = repoInstance.dbSetup()
		if err != nil {
			log.Errorf("failed to auto migrate db: %s", err.Error())
		}
	}

	return nil
}

func GetDB() *gorm.DB {
	return repoInstance.DB
}
