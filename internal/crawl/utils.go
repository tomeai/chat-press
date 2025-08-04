package crawl

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/publicsuffix"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

func normalizeURL(raw string) string {
	u, err := url.Parse(raw)
	if err != nil {
		return ""
	}
	u.Fragment = "" // 去除锚点

	// 清理追踪参数
	q := u.Query()
	for key := range q {
		if strings.HasPrefix(key, "utm_") || key == "ref" {
			q.Del(key)
		}
	}
	u.RawQuery = q.Encode()

	// 去除尾部斜杠（非根路径）
	if strings.HasSuffix(u.Path, "/") && u.Path != "/" {
		u.Path = strings.TrimRight(u.Path, "/")
	}

	return u.String()
}

func genMD5(s string) string {
	hash := md5.Sum([]byte(s))
	return hex.EncodeToString(hash[:])
}

func isTsinghuaSubdomain(link string) bool {
	u, err := url.Parse(link)
	if err != nil {
		return false
	}
	return strings.HasSuffix(u.Hostname(), ".ashes.vip")
	//return strings.HasSuffix(u.Hostname(), ".tsinghua.edu.cn")
}

func extractRootDomain(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	host := u.Hostname()
	eTLDPlusOne, err := publicsuffix.EffectiveTLDPlusOne(host)
	if err != nil {
		return "", err
	}

	return eTLDPlusOne, nil
}

// isDetailURL 基于 URL 的启发式判断
func isDetailURL(u string) bool {
	lower := strings.ToLower(u)
	detailKeywords := []string{"article", "detail", "view", "post", "read", "content"}
	for _, kw := range detailKeywords {
		if strings.Contains(lower, kw) {
			return true
		}
	}

	// 匹配 /2024/06/27/12345.html 或类似数字型 URL
	matched, _ := regexp.MatchString(`/\d{4}/\d{2}/\d{2}/\d+`, u)
	if matched {
		return true
	}

	// 末尾是 .html 并包含数字
	matched, _ = regexp.MatchString(`\d+\.html?$`, u)
	return matched
}

func isListPageByURL(pageURL string) bool {
	u, err := url.Parse(pageURL)
	if err != nil {
		return false
	}

	path := strings.ToLower(u.Path)

	// 常见列表页面URL特征
	listPatterns := []string{
		"/list",
		"/index",
		"/category",
		"/products",
		"/articles",
		"/news",
		"/blog",
	}

	for _, pattern := range listPatterns {
		if strings.Contains(path, pattern) {
			return true
		}
	}

	// 检查分页参数
	query := u.Query()
	if query.Get("page") != "" || query.Get("p") != "" {
		return true
	}

	return false
}

func isListPageByContent(pageURL string) (bool, error) {
	resp, err := http.Get(pageURL)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	content := string(body)
	return analyzeHTMLContent(content), nil
}

func analyzeHTMLContent(content string) bool {
	content = strings.ToLower(content)

	// 检查列表相关的HTML结构
	listIndicators := []string{
		"<ul", "<ol", // 无序和有序列表
		"class=\"list", "class='list",
		"id=\"list", "id='list",
		"<table", // 表格形式的列表
	}

	indicatorCount := 0
	for _, indicator := range listIndicators {
		if strings.Contains(content, indicator) {
			indicatorCount++
		}
	}

	// 检查分页相关元素
	paginationPatterns := []string{
		"pagination", "pager", "page-nav",
		"next", "prev", "上一页", "下一页",
	}

	hasPagination := false
	for _, pattern := range paginationPatterns {
		if strings.Contains(content, pattern) {
			hasPagination = true
			break
		}
	}

	// 检查重复的链接结构
	linkPattern := regexp.MustCompile(`<a[^>]*href=[^>]*>`)
	links := linkPattern.FindAllString(content, -1)

	// 如果有多个指标符合条件，认为是列表页面
	return indicatorCount >= 2 || hasPagination || len(links) > 10
}

func isListPageByDOM(pageURL string) (bool, error) {
	resp, err := http.Get(pageURL)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return false, err
	}

	return analyzeDOM(doc), nil
}

func analyzeDOM(doc *goquery.Document) bool {
	score := 0

	// 检查列表元素数量
	listItems := doc.Find("ul li, ol li").Length()
	if listItems > 5 {
		score += 2
	}

	// 检查表格行数
	tableRows := doc.Find("table tr").Length()
	if tableRows > 3 {
		score += 2
	}

	// 检查重复的结构化内容
	articles := doc.Find("article, .item, .product, .news-item").Length()
	if articles > 3 {
		score += 3
	}

	// 检查分页元素
	pagination := doc.Find(".pagination, .pager, .page-nav").Length()
	if pagination > 0 {
		score += 2
	}

	// 检查标题中的列表关键词
	title := strings.ToLower(doc.Find("title").Text())
	listKeywords := []string{"列表", "list", "index", "目录", "分类"}
	for _, keyword := range listKeywords {
		if strings.Contains(title, keyword) {
			score += 1
			break
		}
	}

	// 检查链接密度
	totalLinks := doc.Find("a[href]").Length()
	contentLength := len(doc.Find("body").Text())
	if contentLength > 0 && totalLinks > contentLength/100 {
		score += 1
	}

	return score >= 4
}
