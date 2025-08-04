package crawl

import (
	"fmt"
	"github.com/gocolly/colly"
	"go.uber.org/zap"
	"net"
	"net/http"
	"strings"
	"time"
)

type Spider struct {
	collyCollector *colly.Collector
	logger         *zap.Logger
}

func NewSpider(logger *zap.Logger) *Spider {

	// Redis 去重配置
	//storage := &redisstorage.Storage{
	//	Address:  "r-bp12ggg21n0g6khw6opd.redis.rds.aliyuncs.com:6379",
	//	Password: "416798Gao!",
	//	DB:       5,
	//	Prefix:   "colly",
	//}
	//
	//if err := storage.Init(); err != nil {
	//	log.Fatal("Redis 初始化失败：", err)
	//}
	//defer storage.Client.Close()
	//
	//c.SetStorage(storage)

	return &Spider{
		logger: logger,
	}
}

func (spider *Spider) Start(target string, maxDepth int) error {
	logger := spider.logger.Named("Spider Start")
	//rePattern := fmt.Sprintf(`^https?://([a-zA-Z0-9-]+\.)*%s(/|$)`, regexp.QuoteMeta(target))
	//re := regexp.MustCompile(rePattern)

	c := colly.NewCollector(
		colly.Async(true),
		colly.MaxDepth(maxDepth),
		colly.IgnoreRobotsTxt(),
		//colly.URLFilters(re),
	)
	c.WithTransport(&http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout: 5 * time.Second,
	})

	// 限流配置
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 100,
		Delay:       500 * time.Millisecond,
		RandomDelay: 500 * time.Millisecond,
	})

	allowDomain, err := extractRootDomain(target)
	if err != nil {
		return err
	}

	// 请求日志
	c.OnRequest(func(r *colly.Request) {
		logger.Info(fmt.Sprintf("🔍 Visiting: %s", r.URL.String()))

		if strings.Contains(r.URL.String(), "mailto:") {
			r.Abort()
		}

		if !strings.HasSuffix(strings.TrimRight(r.URL.String(), "/"), allowDomain) {
			r.Abort()
		}

		// 下载器替换（替换为rod）
	})

	// 错误日志
	c.OnError(func(r *colly.Response, err error) {
		logger.Error(r.Request.URL.String())
	})

	c.OnResponse(func(r *colly.Response) {
		// 如何存储到 s3
		//url := r.Request.URL.String()

		//header := &http.Header{}
		//header.Add("x-cos-meta-spider-url", url)
		//header.Add("x-cos-meta-spider-datetime", time.Now().String())
		//
		//key := fmt.Sprintf("%s/%s.html", r.Request.URL.Hostname(), genMD5(r.Request.URL.String()))
		//err := putFile(bytes.NewReader(r.Body), key, header)
		//if err != nil {
		//	fmt.Println(err)
		//}
	})

	// 处理链接
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		rawLink := e.Request.AbsoluteURL(e.Attr("href"))
		//link := normalizeURL(e.Request.AbsoluteURL(rawLink))
		//if link == "" {
		//	return
		//}
		err := e.Request.Visit(rawLink)
		if err != nil {
			if err == colly.ErrAlreadyVisited {
				logger.Info(fmt.Sprintf("🟡 已访问，跳过: %s", rawLink))
			} else {
				logger.Error(fmt.Sprintf("⚠️ 访问失败: %s", rawLink))
			}
		}

	})

	if err := c.Visit(target); err != nil {
		logger.Error("首次访问失败")
		return err
	}
	c.Wait()
	// todo: 当采集完成进行记录
	logger.Info("✅ 所有采集任务已完成！")
	return nil
}
