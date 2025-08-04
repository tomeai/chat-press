package crawl

type PageAnalyzer struct {
	urlWeight     float64
	contentWeight float64
	domWeight     float64
}

func NewPageAnalyzer() *PageAnalyzer {
	return &PageAnalyzer{
		urlWeight:     0.3,
		contentWeight: 0.3,
		domWeight:     0.4,
	}
}

func (pa *PageAnalyzer) IsListPage(pageURL string) (bool, float64, error) {
	var totalScore float64

	// URL判断
	if isListPageByURL(pageURL) {
		totalScore += pa.urlWeight
	}

	// 内容判断
	isListByContent, err := isListPageByContent(pageURL)
	if err != nil {
		return false, 0, err
	}
	if isListByContent {
		totalScore += pa.contentWeight
	}

	// DOM判断
	isListByDOM, err := isListPageByDOM(pageURL)
	if err != nil {
		return false, 0, err
	}
	if isListByDOM {
		totalScore += pa.domWeight
	}

	return totalScore >= 0.5, totalScore, nil
}
