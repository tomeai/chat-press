## bak

```
//func main() {
//	// www.pushcode.cn
//	//startURL := "https://www.tsinghua.edu.cn/"
//
//	startURL := "https://blog.ashes.vip/"
//	// 创建主 Collector
//	c := colly.NewCollector(
//		//colly.AllowedDomains("*.tsinghua.edu.cn"), // 限定本域名
//
//		//
//		colly.MaxDepth(100),
//		colly.Async(true),
//	)
//
//	// redis://baiwang:Baiwang!@r-bp12ggg21n0g6khw6opd.redis.rds.aliyuncs.com:6379/1
//	// 去重
//	// 初始化 Redis 存储（默认使用 colly:visited 前缀）
//	storage := &redisstorage.Storage{
//		Address:  "r-bp12ggg21n0g6khw6opd.redis.rds.aliyuncs.com:6379",
//		Password: "416798Gao!",
//		DB:       5,
//		Prefix:   "colly:",
//	}
//
//	err := storage.Init()
//	if err != nil {
//		log.Fatal("Redis 初始化失败：", err)
//	}
//
//	// 挂载到 Collector 上
//	c.SetStorage(storage)
//
//	defer storage.Client.Close()
//
//	// http客户端
//	client := &http.Client{
//		Timeout: 5 * time.Second,
//		Transport: &http.Transport{
//			DialContext: (&net.Dialer{
//				Timeout:   5 * time.Second,
//				KeepAlive: 30 * time.Second,
//			}).DialContext,
//			TLSHandshakeTimeout: 5 * time.Second,
//		},
//	}
//
//	// 注入自定义 Transport 到 Colly（使用同一个 Transport 对象）
//	c.WithTransport(client.Transport)
//
//	// 配置代理
//	//err := c.SetProxy("")
//	//if err != nil {
//	//	return
//	//}
//
//	c.Limit(&colly.LimitRule{
//		DomainGlob:  "*", // 对所有域名生效（可指定 "example.com"）
//		Parallelism: 100, // 同时最多 2 个请求
//		//Delay:       2 * time.Second, // 每两个请求间隔 2 秒
//		//RandomDelay: 1 * time.Second, // 附加随机延迟，避免固定节奏
//	})
//
//	// 打印访问的 URL
//	c.OnRequest(onRequest)
//
//	// 遇到 a[href] 时，继续访问链接
//	c.OnHTML("a[href]", onHtml)
//
//	// 错误处理
//	c.OnError(onError)
//
//	// 启动爬虫
//	err = c.Visit(startURL)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// 并发
//	c.Wait()
//}
```