package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	cli2 "github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"net/http"
)

type Server struct {
	HttpServer *http.Server
}

func NewServer(cli *cli2.Context, taskHandler *TaskHandler, logger *zap.Logger) *Server {

	handler := gin.Default()
	// 日志记录（暂时使用中间件记录）
	//handler.Use(middleware.LoggerMiddleware)

	// 业务路由
	taskHandler.RegisterRouter(handler)

	addr := fmt.Sprintf("%s:%s", cli.String("host"), cli.String("port"))
	logger.Info(fmt.Sprintf("listening on -> %s", addr))

	srv := &Server{
		HttpServer: &http.Server{
			Addr:    addr,
			Handler: handler,
		},
	}

	return srv
}
