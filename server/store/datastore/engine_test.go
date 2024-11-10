package datastore_test

import (
    "context"
    "os"

    "github.com/joho/godotenv"

    "goamin/server/store"
    "goamin/server/store/datastore"
)

var (
    ctx    = context.Background()
    engine store.Store
)

func init() {
    if err := godotenv.Load("../../../.env"); err != nil {
        panic(err)
    }

    opts := &store.Opts{
        Driver: os.Getenv("DATABASE_DRIVER"),
        Config: os.Getenv("DATABASE_DATASOURCE"),
        GORM: store.GORM{
            Log:     false,
            ShowSQL: false,
        },
    }
    e, err := datastore.NewEngine(ctx, opts)
    if err != nil {
        panic(err)
    }

    engine = e
}
