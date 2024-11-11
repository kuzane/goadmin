package datastore

import (
    "context"
    "database/sql"
    "fmt"
    "sync"
    "time"

    "gorm.io/driver/mysql"
    "gorm.io/driver/postgres"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"

    "github.com/kuzane/goadmin/server/store"
    "github.com/kuzane/goadmin/server/store/datastore/migration"
)

var _ store.Store = (*storage)(nil)

type storage struct {
    engine *gorm.DB
    lock   sync.Mutex
}

const (
    DriverSqlite   = "sqlite"
    DriverMysql    = "mysql"
    DriverPostgres = "postgres"
)

func SupportedDriver(driver string) bool {
    switch driver {
    case DriverMysql, DriverPostgres, DriverSqlite:
        return true
    default:
        return false
    }
}

func NewEngine(ctx context.Context, opts *store.Opts) (store.Store, error) {
    var (
        s   = storage{}
        err error
    )

    s.lock.Lock()
    defer s.lock.Unlock()

    var log_level logger.LogLevel
    if opts.ShowSQL {
        log_level = logger.Info
    } else {
        log_level = logger.Silent
    }

    if s.engine == nil {
        pool, err := getConnPool(ctx, opts)
        if err != nil {
            return nil, err
        }
        switch opts.Driver {
        case DriverMysql:
            s.engine, err = gorm.Open(mysql.New(mysql.Config{Conn: pool}), &gorm.Config{
                PrepareStmt:            true,
                SkipDefaultTransaction: true,
                Logger:                 newGORMLogger(log_level),
            })
        case DriverSqlite:
            s.engine, err = gorm.Open(sqlite.New(sqlite.Config{Conn: pool}), &gorm.Config{
                PrepareStmt:            true,
                SkipDefaultTransaction: true,
                Logger:                 newGORMLogger(log_level),
            })
        case DriverPostgres:
            s.engine, err = gorm.Open(postgres.New(postgres.Config{Conn: pool}), &gorm.Config{
                PrepareStmt:            true,
                SkipDefaultTransaction: true,
                Logger:                 newGORMLogger(log_level),
            })
        }
    }

    return &s, err
}

// 获取连接池
func getConnPool(ctx context.Context, opts *store.Opts) (db *sql.DB, err error) {
    db, err = sql.Open(opts.Driver, opts.Config)
    if err != nil {
        return nil, err
    }
    // 对连接池进行设置
    db.SetMaxOpenConns(opts.MaxOpenConn)
    db.SetMaxIdleConns(opts.MaxIdleConn)
    if opts.MaxLifeTime != 0 {
        db.SetConnMaxLifetime(time.Second * time.Duration(opts.MaxLifeTime))
    }
    if opts.MaxIdleTime != 0 {
        db.SetConnMaxIdleTime(time.Second * time.Duration(opts.MaxIdleTime))
    }

    if err := db.PingContext(ctx); err != nil {
        return nil, fmt.Errorf("ping database <%s> error,%s", opts.Driver, err.Error())
    }

    return db, nil
}

func (s *storage) Close() error {
    if s.engine == nil {
        return nil
    }
    sqldb, _ := s.engine.DB()
    return sqldb.Close()
}

func (s *storage) Migrate(ctx context.Context) error {
    return migration.Migrate(ctx, s.engine)
}
