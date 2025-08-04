package main

import (
	"context"
	"errors"
	"fmt"
	cli2 "github.com/urfave/cli/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"seed-detect/internal/api"
	"seed-detect/internal/crawl"
	"seed-detect/internal/utils"
	"syscall"
	"time"
)

type App struct {
	ctx      context.Context
	cancel   context.CancelFunc
	exitChan <-chan os.Signal
}

func NewApp(exitChan <-chan os.Signal) *App {
	ctx, cancel := context.WithCancel(context.Background())
	return &App{
		ctx:      ctx,
		cancel:   cancel,
		exitChan: exitChan,
	}
}

func (app *App) Run(args []string) {
	cli := cli2.NewApp()
	cli.Name = "pipeline"
	cli.Flags = []cli2.Flag{
		&cli2.StringFlag{
			Name:  "host",
			Value: "0.0.0.0",
		},
		&cli2.StringFlag{
			Name:  "port",
			Value: "6003",
		},
	}
	cli.Action = func(c *cli2.Context) error {
		options := []fx.Option{
			// go context
			fx.Provide(func() context.Context {
				return app.ctx
			}),
			// fx context
			fx.Provide(func() *cli2.Context {
				return c
			}),
			// log
			fx.Provide(func() *zap.Logger {
				return utils.ZlogInit()
			}),
		}
		options = append(options,
			fx.Provide(crawl.NewSpider),
			fx.Provide(api.NewTaskHandler),
			// 数据接收服务
			fx.Provide(api.NewServer),
			fx.Invoke(NewHttpServer),
		)
		depInj := fx.New(options...)
		if err := depInj.Start(app.ctx); err != nil {
			return err
		}

		<-app.exitChan
		stopCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := depInj.Stop(stopCtx); err != nil {
			fmt.Printf("[Fx] ERROR: Failed to stop cleanly: %v\n", err)
		}
		app.cancel()
		fmt.Printf("[Fx] Cleanly stopped\n")
		return nil
	}
	_ = cli.RunContext(app.ctx, args)
}

func NewHttpServer(lc fx.Lifecycle, server *api.Server, logger *zap.Logger) {

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := server.HttpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
					logger.Error("HTTP server error", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return server.HttpServer.Shutdown(ctx)
		},
	})
}

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	NewApp(c).Run(os.Args)
}
