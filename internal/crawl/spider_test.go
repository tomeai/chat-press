package crawl

import (
	"fmt"
	"log"
	"testing"
)

func TestParseUrl(t *testing.T) {
	//target := "https://plan.music.tsinghua.edu.cn/"
	//target := "https://www.baidu.com"
	target := "https://blog.ashes.vip/"
	log.Println(extractRootDomain(target))
}

func TestIsDetail(t *testing.T) {
	//target := "https://www.tsinghua.edu.cn/info/1182/119870.htm"
	target := "https://blog.ashes.vip/2022/04/15/%E8%BD%AF%E4%BB%B6%E9%A1%B9%E7%9B%AE%E7%AE%A1%E7%90%86/day11/"
	fmt.Println(isDetailURL(target))
}

func TestExtract(t *testing.T) {
	// 创建提取器
	extractor := NewContentExtractor()

	// 从URL提取
	content, err := extractor.ExtractFromURL("https://www.tsinghua.edu.cn/info/1182/119870.htm")

	// 从HTML字符串提取
	//htmlContent := `
	//<html>
	//<head>
	//	<title>示例文章标题</title>
	//	<meta name="author" content="张三">
	//	<meta property="article:published_time" content="2024-01-01T10:00:00Z">
	//</head>
	//<body>
	//	<article>
	//		<h1>示例文章标题</h1>
	//		<div class="author">作者：张三</div>
	//		<div class="time">2024年1月1日</div>
	//		<div class="content">
	//			<p>这是文章的第一段内容，包含了主要的信息。</p>
	//			<p>这是文章的第二段内容，继续描述相关内容。</p>
	//			<p>这是文章的第三段内容，提供更多详细信息。</p>
	//		</div>
	//	</article>
	//</body>
	//</html>
	//`

	//content, err := extractor.ExtractFromHTML(htmlContent)
	if err != nil {
		log.Fatal(err)
	}

	// 输出结果
	fmt.Printf("标题: %s\n", content.Title)
	fmt.Printf("作者: %s\n", content.Author)
	fmt.Printf("发布时间: %s\n", content.PubTime.Format("2006-01-02 15:04:05"))
	fmt.Printf("正文内容:\n%s\n", content.Content)

	// 输出文本节点分析
	fmt.Println("\n文本节点分析:")
	for i, node := range content.TextNodes {
		fmt.Printf("节点 %d: 密度=%.2f, 链接密度=%.2f, 字数=%d, 标签=%s\n",
			i+1, node.Density, node.LinkDensity, node.WordCount, node.TagName)
	}
}

func TestPageList(t *testing.T) {
	//url := "https://www.sppm.tsinghua.edu.cn/syxx/xshd/3.htm"

	//url := "https://www.tsinghua.edu.cn/yxsz.htm"

	//url := "https://www.tsinghua.edu.cn"

	url := "https://www.sppm.tsinghua.edu.cn/info/1006/3088.htm"

	// 方法1：基于URL判断
	if isListPageByURL(url) {
		fmt.Println("根据URL判断，这是一个列表页面")
	}

	// 方法2：基于内容判断
	isList, err := isListPageByContent(url)
	if err != nil {
		log.Fatal(err)
	}
	if isList {
		fmt.Println("根据内容判断，这是一个列表页面")
	}

	isListByDOM, err := isListPageByDOM(url)
	if err != nil {
		log.Fatal(err)
	}
	if isListByDOM {
		fmt.Println("dom判断")
	}

	// 方法3：综合判断
	analyzer := NewPageAnalyzer()
	isList, confidence, err := analyzer.IsListPage(url)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("综合判断结果：%t，置信度：%.2f\n", isList, confidence)
}
