package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v3"

	"goamin/pkg/logger"
	"goamin/pkg/utils"
	"goamin/server/router"
	"goamin/server/router/middleware"
)

const (
	shutdownTimeout = time.Second * 5
)

var (
	stopServerFunc     context.CancelCauseFunc = func(error) {}
	shutdownCancelFunc context.CancelFunc      = func() {}
	shutdownCtx                                = context.Background()
)

func Run(ctx context.Context, c *cli.Command) error {
	// 初始化日志
	if err := logger.SetupGlobalLogger(ctx, c, true); err != nil {
		return err
	}
	// 参数校验
	if err := setupServerHost(c); err != nil {
		return err
	}
	// 初始化翻译器
	if err := utils.InitTrans(c.String("server-trans")); err != nil {
		return err
	}
	// 基于日志级别设置gin日志
	if zerolog.GlobalLevel() > zerolog.DebugLevel {
		gin.SetMode(gin.ReleaseMode)
	}
	// 设置上下文ctx
	ctx, ctxCancel := context.WithCancelCause(ctx)
	stopServerFunc = func(err error) {
		if err != nil {
			log.Error().Err(err).Msg("shutdown of whole server")
		}
		stopServerFunc = func(error) {}
		shutdownCtx, shutdownCancelFunc = context.WithTimeout(shutdownCtx, shutdownTimeout)
		ctxCancel(err)
	}
	defer stopServerFunc(nil)
	defer shutdownCancelFunc()
	// 初始化后端存储
	_store, err := setupStore(ctx, c)
	if err != nil {
		return fmt.Errorf("can't setup store: %w", err)
	}
	defer func() {
		if err := _store.Close(); err != nil {
			log.Error().Err(err).Msg("could not close store")
		}
	}()
	// 设置severconfig
	if err := setupEvilGlobals(ctx, c, _store); err != nil {
		return err
	}
	// 启动服务
	handler := router.Load(_store, middleware.Logger(time.RFC3339, true), middleware.Store(_store))
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", c.String("server-port")),
		Handler: handler,
	}

	go func() {
		<-ctx.Done()
		log.Info().Msg("shutdown http server ...")
		if err = httpServer.Shutdown(ctx); err != nil {
			log.Error().Err(err).Msg("shutdown http server failed")
		} else {
			log.Info().Msg("http server stopped")
		}
	}()

	log.Info().Msg("starting http server ...")
	if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Error().Err(err).Msg("http server failed")
		stopServerFunc(fmt.Errorf("http server failed: %w", err))
	}

	return nil
}
