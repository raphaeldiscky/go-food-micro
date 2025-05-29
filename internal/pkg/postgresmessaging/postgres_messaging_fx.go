package postgresmessaging

import (
	"go.uber.org/fx"
	"gorm.io/gorm"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/persistmessage"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/postgresmessaging/messagepersistence"
)

var Module = fx.Module(
	"postgresmessagingfx",
	fx.Provide(
		messagepersistence.NewPostgresMessagePersistenceDBContext,
		messagepersistence.NewPostgresMessageService,
	),
	fx.Invoke(migrateMessaging),
)

func migrateMessaging(db *gorm.DB) error {
	err := db.Migrator().AutoMigrate(&persistmessage.StoreMessage{})

	return err
}
