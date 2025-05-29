// Package messagepersistence provides the postgres message service.
package messagepersistence

import (
	"context"
	"os"
	"testing"
	"time"

	"emperror.dev/errors"
	"github.com/stretchr/testify/suite"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"gorm.io/gorm"

	uuid "github.com/satori/go.uuid"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/config"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/config/environment"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/persistmessage"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	defaultLogger "github.com/raphaeldiscky/go-food-micro/internal/pkg/logger/defaultlogger"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger/external/fxlog"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger/zap"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/mapper"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/postgresgorm"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/postgresgorm/helpers/gormextensions"
)

// postgresMessageServiceTest is the test suite for the postgres message service.
type postgresMessageServiceTest struct {
	suite.Suite
	DB                  *gorm.DB
	logger              logger.Logger
	messagingRepository persistmessage.MessagePersistenceService
	dbContext           *PostgresMessagePersistenceDBContext
	storeMessages       []*persistmessage.StoreMessage
	ctx                 context.Context
	dbFilePath          string
	app                 *fxtest.App
}

// TestPostgresMessageService tests the postgres message service.
func TestPostgresMessageService(t *testing.T) {
	suite.Run(
		t,
		&postgresMessageServiceTest{logger: defaultLogger.GetLogger()},
	)
}

// func (c *postgresMessageServiceTest) SetupSuite() {
//	opts, err := gorm2.NewGormTestContainers(defaultLogger.GetLogger()).
//		PopulateContainerOptions(context.Background(), c.T())
//	c.Require().NoError(err)
//
//	gormDB, err := postgresgorm.NewGorm(opts)
//	c.Require().NoError(err)
//	c.DB = gormDB
//
//	err = migrationDatabase(gormDB)
//	c.Require().NoError(err)
//
//	c.dbContext = NewPostgresMessagePersistenceDBContext(gormDB)
//	c.messagingRepository = NewPostgresMessageService(
//		c.dbContext,
//		defaultLogger.GetLogger(),
//	)
//}

// SetupTest sets up the test.
func (c *postgresMessageServiceTest) SetupTest() {
	var gormDBContext *PostgresMessagePersistenceDBContext
	var gormOptions *postgresgorm.GormOptions

	app := fxtest.New(
		c.T(),
		config.ModuleFunc(environment.Test),
		zap.Module,
		fxlog.FxLogger,
		core.Module,
		postgresgorm.Module,
		fx.Decorate(
			func(cfg *postgresgorm.GormOptions) (*postgresgorm.GormOptions, error) {
				// using sql-lite with a database file
				cfg.UseSQLLite = true

				return cfg, nil
			},
		),
		fx.Provide(NewPostgresMessagePersistenceDBContext),
		fx.Populate(&gormDBContext),
		fx.Populate(&gormOptions),
	).RequireStart()

	c.dbContext = gormDBContext
	c.dbFilePath = gormOptions.DND()
	c.app = app

	c.initDB()
}

// TearDownTest tears down the test.
func (c *postgresMessageServiceTest) TearDownTest() {
	err := c.cleanupDB()
	c.Require().NoError(err)

	mapper.ClearMappings()

	c.app.RequireStop()
}

// func (c *postgresMessageServiceTest) SetupTest() {
//	ctx := context.Background()
//	c.ctx = ctx
//	p, err := seedData(context.Background(), c.DB)
//	c.Require().NoError(err)
//	c.storeMessages = p
//}
//
// func (c *postgresMessageServiceTest) TearDownTest() {
//	err := c.cleanupPostgresData()
//	c.Require().NoError(err)
//}

// BeginTx begins the transaction.
func (c *postgresMessageServiceTest) BeginTx() {
	c.logger.Info("starting transaction")
	tx := c.dbContext.DB().Begin()
	gormContext := gormextensions.SetTxToContext(c.ctx, tx)
	c.ctx = gormContext
}

// CommitTx commits the transaction.
func (c *postgresMessageServiceTest) CommitTx() {
	tx := gormextensions.GetTxFromContextIfExists(c.ctx)
	if tx != nil {
		c.logger.Info("committing transaction")
		tx.Commit()
	}
}

// TestAdd tests the add.
func (c *postgresMessageServiceTest) TestAdd() {
	message := &persistmessage.StoreMessage{
		ID:            uuid.NewV4(),
		MessageStatus: persistmessage.Processed,
		Data:          "test data 3",
		DataType:      "string",
		CreatedAt:     time.Now(),
		DeliveryType:  persistmessage.Outbox,
	}

	c.BeginTx()
	err := c.messagingRepository.Add(c.ctx, message)
	c.CommitTx()

	c.Require().NoError(err)

	m, err := c.messagingRepository.GetByID(c.ctx, message.ID)
	if err != nil {
		return
	}

	c.Assert().NotNil(m)
	c.Assert().Equal(message.ID, m.ID)
}

// initDB initializes the database.
func (c *postgresMessageServiceTest) initDB() {
	err := migrateGorm(c.dbContext.DB())
	c.Require().NoError(err)

	storeMessages, err := seedData(c.dbContext.DB())
	c.Require().NoError(err)

	c.storeMessages = storeMessages
}

// cleanupDB cleans up the database.
func (c *postgresMessageServiceTest) cleanupDB() error {
	sqldb, err := c.dbContext.DB().DB()
	if err != nil {
		return err
	}
	err = sqldb.Close()
	if err != nil {
		return err
	}

	// removing sql-lite file
	err = os.Remove(c.dbFilePath)

	return err
}

// migrateGorm migrates the gorm.
func migrateGorm(db *gorm.DB) error {
	err := db.AutoMigrate(&persistmessage.StoreMessage{})
	if err != nil {
		return err
	}

	return nil
}

// seedData seeds the data.
func seedData(
	db *gorm.DB,
) ([]*persistmessage.StoreMessage, error) {
	messages := []*persistmessage.StoreMessage{
		{
			ID:            uuid.NewV4(),
			MessageStatus: persistmessage.Processed,
			Data:          "test data",
			DataType:      "string",
			CreatedAt:     time.Now(),
			DeliveryType:  persistmessage.Outbox,
		},
		{
			ID:            uuid.NewV4(),
			MessageStatus: persistmessage.Processed,
			Data:          "test data 2",
			DataType:      "string",
			CreatedAt:     time.Now(),
			DeliveryType:  persistmessage.Outbox,
		},
	}

	// seed data
	err := db.CreateInBatches(messages, len(messages)).Error
	if err != nil {
		return nil, errors.Wrap(err, "error in seed database")
	}

	return messages, nil
}
