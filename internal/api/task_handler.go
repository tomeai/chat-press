package api

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"seed-detect/internal/crawl"
)

type TaskHandler struct {
	spider *crawl.Spider
	logger *zap.Logger
}

type TaskInfo struct {
	Url      string `json:"url"`
	MaxDepth int    `json:"maxDepth"`
}

func NewTaskHandler(spider *crawl.Spider, logger *zap.Logger) *TaskHandler {
	return &TaskHandler{
		spider: spider,
		logger: logger,
	}
}

func (h *TaskHandler) RegisterRouter(server *gin.Engine) {
	// 小必姐消息查询
	xbj := server.Group("/task")
	xbj.POST("/submit", h.submitTask)
}

func (h *TaskHandler) submitTask(ctx *gin.Context) {
	logger := h.logger.Named("TaskHandler submitTask")
	var req TaskInfo

	logger.Info(req.Url)

	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, Result{
			Code: InvalidBody,
			Msg:  "参数不合法",
		})
		return
	}

	go func() {
		maxDepth := req.MaxDepth
		if maxDepth == 0 {
			maxDepth = 10
		}
		// todo: 如果没有返回错误 则说明采集成功  记录数据库
		err := h.spider.Start(req.Url, maxDepth)
		if err != nil {
			logger.Error(err.Error())
		}
	}()

	ctx.JSON(http.StatusOK, Result{
		Data: req,
	})
	return
}
