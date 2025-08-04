## CharPress

## todo

```
1. 构建采集多叉树？
2. 失败的请求 流式  golang接受
3. withOptions学习
4. 是不是只需要body里面的url就可以了
5. golang发现种子  然后交给 scrapy-redis
    1. 如果是详情 就不需要下探了
    2. 基于 规则
    3. 判断是不是列表 比较好判断
    4. 有可能跳转到 站外 比如公众号
6. 先判断有没有 sitmap 或者 rss 没有在递归
7. 后期构建模型判断
8. 政务、医院、高校
9. 统一使用rod??

## 策略
1. 先采集 最后清洗
2. 基础数据入库 使用dataflow  
    1. 进入dataflow 在操作cos吗
    2. dataflow 构建 数据入库 pipeline
```

## 种子

```
# 采集
https://gaokao.chsi.com.cn/sch/schoolInfo--schId-1,categoryId-26177,mindex-2.dhtml

高校：院系  政务：下面导航  医院：好大夫 https://www.haodf.com/hospital/list-11.html
    1. 高校新闻 -> 院系新闻

https://www.govdir.cn/favorites/shandong
https://gaokao.chsi.com.cn/sch/search--ss-on,option-qg,searchType-1,start-0.dhtml
```

## 学习

```
https://github.com/BruceDone/awesome-crawler?utm_source=chatgpt.com
https://github.com/Nandakumartc/scraper-crawler?utm_source=chatgpt.com
https://github.com/lizongying/go-crawler
https://lizongying.github.io/go-crawler/docs/docs/usage/example/

# python
https://github.com/pcatattacks/domain-web-crawler/tree/master
https://github.com/Dineshs91/crawler?utm_source=chatgpt.com
https://github.com/khyatig0206/BFS-DFS-WebCrawler/?utm_source=chatgpt.com
https://github.com/danhilse/web-scraper/blob/main/contxt/scraper.py
```