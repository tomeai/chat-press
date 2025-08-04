package crawl

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode"
)

// https://github.com/GeneralNewsExtractor/GeneralNewsExtractor?tab=readme-ov-file
// 基于文本及符号密度的网页正文提取方法
// https://claude.ai/public/artifacts/4421f300-f29d-490d-a23c-2f78355fe508

// ContentExtractor 内容提取器
type ContentExtractor struct {
	// 文本密度阈值
	TextDensityThreshold float64
	// 最小文本长度
	MinTextLength int
	// 标题相关选择器
	TitleSelectors []string
	// 时间相关选择器
	TimeSelectors []string
	// 作者相关选择器
	AuthorSelectors []string
}

// ExtractedContent 提取的内容结构
type ExtractedContent struct {
	Title     string     `json:"title"`
	Author    string     `json:"author"`
	PubTime   time.Time  `json:"pub_time"`
	Content   string     `json:"content"`
	TextNodes []TextNode `json:"text_nodes"`
}

// TextNode 文本节点信息
type TextNode struct {
	Text        string  `json:"text"`
	Density     float64 `json:"density"`
	TagName     string  `json:"tag_name"`
	WordCount   int     `json:"word_count"`
	LinkCount   int     `json:"link_count"`
	LinkDensity float64 `json:"link_density"`
}

// NewContentExtractor 创建新的内容提取器
func NewContentExtractor() *ContentExtractor {
	return &ContentExtractor{
		TextDensityThreshold: 1.0,
		MinTextLength:        50,
		TitleSelectors: []string{
			"title",
			"h1",
			".title",
			"#title",
			".article-title",
			".post-title",
			".entry-title",
			"[property='og:title']",
			"[name='title']",
		},
		TimeSelectors: []string{
			"time",
			".time",
			".date",
			".publish-time",
			".pub-time",
			".post-date",
			".article-date",
			"[datetime]",
			"[property='article:published_time']",
			"[name='pubdate']",
		},
		AuthorSelectors: []string{
			".author",
			".by-author",
			".post-author",
			".article-author",
			".writer",
			"[rel='author']",
			"[property='article:author']",
			"[name='author']",
		},
	}
}

// ExtractFromURL 从URL提取内容
func (ce *ContentExtractor) ExtractFromURL(url string) (*ExtractedContent, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	return ce.ExtractFromDocument(doc)
}

// ExtractFromHTML 从HTML字符串提取内容
func (ce *ContentExtractor) ExtractFromHTML(html string) (*ExtractedContent, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}

	return ce.ExtractFromDocument(doc)
}

// ExtractFromDocument 从goquery文档提取内容
func (ce *ContentExtractor) ExtractFromDocument(doc *goquery.Document) (*ExtractedContent, error) {
	content := &ExtractedContent{}

	// 提取标题
	content.Title = ce.extractTitle(doc)

	// 提取作者
	content.Author = ce.extractAuthor(doc)

	// 提取发布时间
	content.PubTime = ce.extractPubTime(doc)

	// 提取正文内容
	textNodes := ce.extractTextNodes(doc)
	content.TextNodes = textNodes
	content.Content = ce.extractMainContent(textNodes)

	return content, nil
}

// extractTitle 提取标题
func (ce *ContentExtractor) extractTitle(doc *goquery.Document) string {
	for _, selector := range ce.TitleSelectors {
		if title := strings.TrimSpace(doc.Find(selector).First().Text()); title != "" {
			// 对于meta标签，获取content属性
			if strings.Contains(selector, "[") {
				if content, exists := doc.Find(selector).First().Attr("content"); exists && content != "" {
					return content
				}
			}
			return title
		}
	}
	return ""
}

// extractAuthor 提取作者
func (ce *ContentExtractor) extractAuthor(doc *goquery.Document) string {
	for _, selector := range ce.AuthorSelectors {
		if author := strings.TrimSpace(doc.Find(selector).First().Text()); author != "" {
			// 对于meta标签，获取content属性
			if strings.Contains(selector, "[") {
				if content, exists := doc.Find(selector).First().Attr("content"); exists && content != "" {
					return content
				}
			}
			return author
		}
	}
	return ""
}

// extractPubTime 提取发布时间
func (ce *ContentExtractor) extractPubTime(doc *goquery.Document) time.Time {
	for _, selector := range ce.TimeSelectors {
		element := doc.Find(selector).First()

		// 尝试从datetime属性获取
		if datetime, exists := element.Attr("datetime"); exists {
			if t, err := ce.parseTime(datetime); err == nil {
				return t
			}
		}

		// 尝试从content属性获取
		if content, exists := element.Attr("content"); exists {
			if t, err := ce.parseTime(content); err == nil {
				return t
			}
		}

		// 尝试从文本内容获取
		if text := strings.TrimSpace(element.Text()); text != "" {
			if t, err := ce.parseTime(text); err == nil {
				return t
			}
		}
	}

	return time.Time{}
}

// parseTime 解析时间字符串
func (ce *ContentExtractor) parseTime(timeStr string) (time.Time, error) {
	// 常见的时间格式
	formats := []string{
		time.RFC3339,
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02 15:04:05",
		"2006-01-02",
		"2006/01/02 15:04:05",
		"2006/01/02",
		"01/02/2006 15:04:05",
		"01/02/2006",
		"Jan 2, 2006 3:04 PM",
		"Jan 2, 2006",
		"January 2, 2006 3:04 PM",
		"January 2, 2006",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, timeStr); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("无法解析时间: %s", timeStr)
}

// extractTextNodes 提取所有文本节点并计算密度
func (ce *ContentExtractor) extractTextNodes(doc *goquery.Document) []TextNode {
	var textNodes []TextNode

	// 移除不需要的元素
	doc.Find("script, style, nav, header, footer, aside, .sidebar, .ad, .advertisement").Remove()

	// 遍历所有可能包含正文的元素
	contentSelectors := []string{"p", "div", "article", "section", "main", "content"}

	for _, selector := range contentSelectors {
		doc.Find(selector).Each(func(i int, s *goquery.Selection) {
			text := ce.cleanText(s.Text())
			if len(text) < ce.MinTextLength {
				return
			}

			// 计算文本密度
			density := ce.calculateTextDensity(s)

			// 计算链接密度
			linkDensity := ce.calculateLinkDensity(s)

			// 获取标签名
			tagName := goquery.NodeName(s)

			textNode := TextNode{
				Text:        text,
				Density:     density,
				TagName:     tagName,
				WordCount:   ce.countWords(text),
				LinkCount:   s.Find("a").Length(),
				LinkDensity: linkDensity,
			}

			textNodes = append(textNodes, textNode)
		})
	}

	// 按文本密度排序
	sort.Slice(textNodes, func(i, j int) bool {
		return textNodes[i].Density > textNodes[j].Density
	})

	return textNodes
}

// calculateTextDensity 计算文本密度
func (ce *ContentExtractor) calculateTextDensity(s *goquery.Selection) float64 {
	text := ce.cleanText(s.Text())
	if len(text) == 0 {
		return 0
	}

	// 计算HTML标签数量
	htmlContent, _ := s.Html()
	tagCount := strings.Count(htmlContent, "<")

	// 文本密度 = 文本长度 / (文本长度 + 标签数量)
	textLen := float64(len(text))
	return textLen / (textLen + float64(tagCount))
}

// calculateLinkDensity 计算链接密度
func (ce *ContentExtractor) calculateLinkDensity(s *goquery.Selection) float64 {
	text := ce.cleanText(s.Text())
	if len(text) == 0 {
		return 0
	}

	// 计算链接文本长度
	linkTextLen := 0
	s.Find("a").Each(func(i int, link *goquery.Selection) {
		linkTextLen += len(ce.cleanText(link.Text()))
	})

	return float64(linkTextLen) / float64(len(text))
}

// extractMainContent 提取主要正文内容
func (ce *ContentExtractor) extractMainContent(textNodes []TextNode) string {
	var contentParts []string

	for _, node := range textNodes {
		// 过滤条件：
		// 1. 文本密度大于阈值
		// 2. 链接密度不能太高（避免导航菜单等）
		// 3. 文本长度足够
		if node.Density >= ce.TextDensityThreshold &&
			node.LinkDensity < 0.5 &&
			len(node.Text) >= ce.MinTextLength {
			contentParts = append(contentParts, node.Text)
		}

		// 限制内容长度，避免过长
		if len(contentParts) >= 10 {
			break
		}
	}

	return strings.Join(contentParts, "\n\n")
}

// cleanText 清理文本
func (ce *ContentExtractor) cleanText(text string) string {
	// 移除多余的空白字符
	re := regexp.MustCompile(`\s+`)
	text = re.ReplaceAllString(text, " ")

	// 移除前后空白
	text = strings.TrimSpace(text)

	return text
}

// countWords 计算单词数量
func (ce *ContentExtractor) countWords(text string) int {
	// 对于中文文本，按字符计算
	if ce.containsCJK(text) {
		return len([]rune(text))
	}

	// 对于英文文本，按单词计算
	words := strings.Fields(text)
	return len(words)
}

// containsCJK 检查是否包含中日韩文字
func (ce *ContentExtractor) containsCJK(text string) bool {
	for _, r := range text {
		if unicode.Is(unicode.Han, r) ||
			unicode.Is(unicode.Hiragana, r) ||
			unicode.Is(unicode.Katakana, r) {
			return true
		}
	}
	return false
}

// AdvancedExtractor 高级提取器，使用更复杂的算法
type AdvancedExtractor struct {
	*ContentExtractor
}

// NewAdvancedExtractor 创建高级提取器
func NewAdvancedExtractor() *AdvancedExtractor {
	return &AdvancedExtractor{
		ContentExtractor: NewContentExtractor(),
	}
}

// ExtractWithBoilerpipe 使用类似Boilerpipe的算法提取内容
func (ae *AdvancedExtractor) ExtractWithBoilerpipe(doc *goquery.Document) (*ExtractedContent, error) {
	content := &ExtractedContent{}

	// 提取标题、作者、时间
	content.Title = ae.extractTitle(doc)
	content.Author = ae.extractAuthor(doc)
	content.PubTime = ae.extractPubTime(doc)

	// 使用Boilerpipe类似算法
	blocks := ae.extractTextBlocks(doc)
	blocks = ae.classifyBlocks(blocks)
	content.Content = ae.extractContentFromBlocks(blocks)

	return content, nil
}

// TextBlock 文本块
type TextBlock struct {
	Text        string
	WordCount   int
	LinkCount   int
	TextDensity float64
	LinkDensity float64
	IsContent   bool
	TagName     string
}

// extractTextBlocks 提取文本块
func (ae *AdvancedExtractor) extractTextBlocks(doc *goquery.Document) []TextBlock {
	var blocks []TextBlock

	// 移除不需要的元素
	doc.Find("script, style, nav, header, footer, aside").Remove()

	// 提取文本块
	doc.Find("*").Each(func(i int, s *goquery.Selection) {
		text := ae.cleanText(s.Clone().Children().Remove().End().Text())
		if len(text) < 30 {
			return
		}

		wordCount := ae.countWords(text)
		linkCount := s.Find("a").Length()

		block := TextBlock{
			Text:        text,
			WordCount:   wordCount,
			LinkCount:   linkCount,
			TextDensity: ae.calculateTextDensity(s),
			LinkDensity: ae.calculateLinkDensity(s),
			TagName:     goquery.NodeName(s),
		}

		blocks = append(blocks, block)
	})

	return blocks
}

// classifyBlocks 分类文本块
func (ae *AdvancedExtractor) classifyBlocks(blocks []TextBlock) []TextBlock {
	if len(blocks) == 0 {
		return blocks
	}

	// 计算平均值
	var totalWords, totalLinks float64
	for _, block := range blocks {
		totalWords += float64(block.WordCount)
		totalLinks += float64(block.LinkCount)
	}

	avgWords := totalWords / float64(len(blocks))
	avgLinks := totalLinks / float64(len(blocks))

	// 分类规则
	for i := range blocks {
		block := &blocks[i]

		// 判断是否为正文内容
		block.IsContent = block.WordCount > int(avgWords*0.5) &&
			block.LinkDensity < 0.3 &&
			block.TextDensity > 0.3 &&
			float64(block.LinkCount) < avgLinks*2
	}

	return blocks
}

// extractContentFromBlocks 从文本块中提取内容
func (ae *AdvancedExtractor) extractContentFromBlocks(blocks []TextBlock) string {
	var contentParts []string

	for _, block := range blocks {
		if block.IsContent {
			contentParts = append(contentParts, block.Text)
		}
	}

	return strings.Join(contentParts, "\n\n")
}
