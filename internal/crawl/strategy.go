package crawl

import (
	"errors"
	"fmt"
	"github.com/gocolly/colly"
	"log"
	"time"
)

func onRequest(r *colly.Request) {
	// 添加请求信息
	//fmt.Printf("Depth %d: Visiting %s\n", r.Depth, r.URL)
	//fmt.Println("Visiting:", r.URL)
	log.Println("Visiting", r.URL.String())
}

func onHtml(e *colly.HTMLElement) {
	// 判断是不是详情页面  入库   hostname、url、content、datetime
	// 如何高效抽取   标题 时间 正文等信息
	//hostname := e.Response.Request.URL.Hostname()
	//fmt.Println("hostname: ", e.Request.URL.String())

	link := e.Request.AbsoluteURL(e.Attr("href"))
	if link != "" && isTsinghuaSubdomain(link) {
		err := e.Request.Visit(link) // 默认是 DFS
		if err != nil {
			// 是否配置
			if errors.Is(err, colly.ErrAlreadyVisited) {
				// 已访问过，正常，不用管
				return
			}
			fmt.Println("visit error:", err)
		}
	}
}

func onError(r *colly.Response, err error) {
	status := r.StatusCode
	url := r.Request.URL.String()

	switch {
	case status == 404:
		fmt.Println("🚫 404 页面不存在，跳过：", url)
	case status >= 500 && status < 600:
		// 服务器错误，可尝试重试
		// 记录在 redis
		//if retryCount[url] < maxRetries {
		//	retryCount[url]++
		//	fmt.Printf("🔁 第 %d 次重试（服务器错误 %d）：%s\n", retryCount[url], status, url)
		//	time.AfterFunc(1*time.Second, func() {
		//		_ = r.Request.Retry()
		//	})
		//} else {
		//	fmt.Println("❌ 重试次数用尽（服务器错误）：", url)
		//}
	case status == 429:
		// 请求太频繁，建议等待长一点再重试
		fmt.Println("⏳ 被限速（429），等待后重试：", url)
		time.AfterFunc(5*time.Second, func() {
			_ = r.Request.Retry()
		})
	default:
		fmt.Printf("❗ 其他错误 [%d]：%s %v\n", status, url, err)
	}
}
