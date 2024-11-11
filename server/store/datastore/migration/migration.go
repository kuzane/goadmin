package migration

import (
    "context"
    "fmt"
    "reflect"

    "gorm.io/gorm"

    "github.com/kuzane/goadmin/server/model"
)

var allBeans = []any{
    new(model.ServerConfig),
    new(model.User),
    new(model.Endpoint),
    new(model.Role),
    new(model.Userlog),
}

func Migrate(ctx context.Context, e *gorm.DB) error {
    if err := syncAll(ctx, e); err != nil {
        return fmt.Errorf("msg: %w", err)
    }
    return nil
}

func syncAll(ctx context.Context, e *gorm.DB) error {
    for _, bean := range allBeans {
        if err := e.WithContext(ctx).AutoMigrate(bean); err != nil {
            return fmt.Errorf("AutoMigrate error '%s': %w", reflect.TypeOf(bean), err)
        }
    }
    return nil
}
