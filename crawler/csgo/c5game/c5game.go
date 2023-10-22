package c5game

import (
	"compress/gzip"
	"fmt"
	"github.com/dop251/goja"
	"github.com/tidwall/gjson"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

var sellHistoryJsRege = regexp.MustCompile(`window.__NUXT__=(.+?)<`)
var existsIdList = &sync.Map{}

type C5Game struct {
	httpClient    *http.Client
	sellDataChan  chan string
	leaseDataChan chan string
	reqMain       bool
}

func NewC5Game(sellDataChan, leaseDataChan chan string) *C5Game {
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: nil})
	if err != nil {
		log.Fatal(err)
	}
	client := &http.Client{}
	client.Timeout = time.Second * 10
	client.Jar = jar
	return &C5Game{
		httpClient:    client,
		sellDataChan:  sellDataChan,
		leaseDataChan: leaseDataChan,
	}
}

func Start() {
	sellDataFileName := "c5game_sell.csv"
	leaseDataFileName := "c5game_lease.csv"
	sellDataChan := make(chan string, 10)
	leaseDataChan := make(chan string, 10)
	c1 := NewC5Game(sellDataChan, leaseDataChan)
	c2 := NewC5Game(sellDataChan, leaseDataChan)
	c3 := NewC5Game(sellDataChan, leaseDataChan)

	// 获取总页数
	pageRegex := regexp.MustCompile(`pages:\s*?(\d+)`)
	listPageContent := c1.GetMainPage()
	pageRegexResult := pageRegex.FindStringSubmatch(listPageContent)
	if len(pageRegexResult) != 2 {
		os.WriteFile("c3gamin_list_pagt.html", []byte(listPageContent), 0755)
		log.Fatal("pageRegexResult error", pageRegexResult)
	}
	pageTotal, _ := strconv.Atoi(pageRegexResult[1])
	log.Println("pageTotal", pageTotal)

	// 打开文件
	sellDataFile, err := os.OpenFile(sellDataFileName, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0755)
	if err != nil {
		log.Fatal(err)
	}
	leaseDataFile, err := os.OpenFile(leaseDataFileName, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0755)
	if err != nil {
		log.Fatal(err)
	}
	defer sellDataFile.Close()
	defer leaseDataFile.Close()

	// 写入文件
	fileover := make(chan bool)
	go func() {
		for val := range sellDataChan {
			//log.Println("write file sell")
			_, err = sellDataFile.WriteString(val)
			if err != nil {
				log.Println(err)
			}
		}
		for val := range leaseDataChan {
			//log.Println("write file lease")
			_, err = leaseDataFile.WriteString(val)
			if err != nil {
				log.Println(err)
			}
		}
		fileover <- true
	}()

	// 分页获取数据
	c5Games := []*C5Game{c1, c2, c3}
	getC5Game := func(num int) *C5Game {
		return c5Games[num]
	}
	count := 0
	num := make(chan int, 3)
	num <- 1
	num <- 1
	num <- 1
	wg := &sync.WaitGroup{}
	for page := 1; page <= pageTotal; page++ {
		<-num
		if count > 2 {
			count = 0
		}
		wg.Add(1)
		go func(p, c int) {
			defer wg.Done()
			getC5Game(c).FetchList(p)
			num <- 1
		}(page, count)
		count++
	}
	wg.Wait()
	close(sellDataChan)
	close(leaseDataChan)
	<-fileover
	log.Println("finish")
}

func (c *C5Game) FetchList(page int) {
	log.Println("###################fetch page list:", page)
	if !c.reqMain {
		_ = c.GetMainPage()
	}

	req := newReq(fmt.Sprintf("https://www.c5game.com/napi/trade/search/v2/items/730/search?limit=42&appId=730&page=%d&sort=0", page), true)
	req.Header.Set("Referer", fmt.Sprintf("https://www.c5game.com/csgo?appId=730&page=%d&limit=42&sort=0", page))
	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Println("req list page", page, resp.StatusCode)
		return
	}
	rawBody := readBody(resp)
	list := gjson.Get(rawBody, "data.list").Array()
	for _, item := range list {
		productId := item.Get("itemId").String()
		rowData := make([]string, 0)
		log.Println("productId", productId)

		// 略过已获取
		if _, ok := existsIdList.Load(productId); ok {
			log.Println("exists productId", productId)
			continue
		}
		existsIdList.Store(productId, 1)

		name := item.Get("itemName").String()
		rowData = append(rowData, name)
		rowData = append(rowData, item.Get("itemInfo.exteriorName").String())
		rowData = append(rowData, item.Get("itemInfo.rarityName").String())
		sellData := append(rowData, item.Get("quantity").String())
		sellData = append(sellData, c.getSellMinPrice(name, productId))
		sellData = append(sellData, c.getSellHistory(name, productId)...)

		c.sellDataChan <- strings.Join(sellData, ",") + "\n"
	}

}

func (c *C5Game) GetMainPage() string {
	//请求一下页面，得到cookie
	resp, err := c.httpClient.Do(newReq("https://www.c5game.com/csgo?appId=730&page=1&limit=42&sort=0", false))
	if err != nil {
		log.Println(err)
		return ""
	}
	defer resp.Body.Close()
	log.Println("req main page", resp.StatusCode)
	c.reqMain = true
	return readBody(resp)
}

// getSellMinPrice 获取在售最低价
func (c *C5Game) getSellMinPrice(name, productId string) string {
	retryCount := 0
retry:
	req := newReq(fmt.Sprintf("https://www.c5game.com/napi/trade/steamtrade/sga/sell/v3/list?itemId=%s&orderBy=2&page=1&limit=10", productId), true)
	req.Header.Set("Referer", fmt.Sprintf("https://www.c5game.com/csgo/%s/%s/sell?orderBy=2", productId, url.QueryEscape(name)))
	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Println(err, req.URL)
		return ""
	}
	if resp.StatusCode != 200 {
		log.Println("req min price", resp.StatusCode, req.URL)
		if retryCount < 2 {
			time.Sleep(time.Second * 2)
			retryCount++
			goto retry
		}
		return ""
	}
	defer resp.Body.Close()
	return gjson.Get(readBody(resp), "data.list.0.price").String()
}

func (c *C5Game) getSellHistory(name, productId string) (history []string) {
	req := newReq(
		fmt.Sprintf("https://www.c5game.com/csgo/%s/%s/record", productId, url.QueryEscape(name)),
		false,
	)
	req.Header.Set("Referer", fmt.Sprintf("https://www.c5game.com/csgo/%s/%s/sell", productId, url.QueryEscape(name)))
	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Println(err, req.URL)
		return
	}
	if resp.StatusCode != 200 {
		log.Println("req min price", resp.StatusCode, req.URL)
		return
	}
	defer resp.Body.Close()
	rawBody := readBody(resp)

	// 页面找到js数据，执行js获取数据量列表json字符串
	result := sellHistoryJsRege.FindStringSubmatch(rawBody)
	if len(result) != 2 {
		log.Println("result error", result)
		return
	}
	jsVm := goja.New()
	historyJs := result[1]
	js := fmt.Sprintf("const result = %s\nconst list = JSON.stringify(result.data[0].list);", historyJs)
	_, err = jsVm.RunString(js)
	if err != nil {
		log.Println("run js error", err)
		return
	}
	resJson := ""
	err = jsVm.ExportTo(jsVm.Get("list"), &resJson)
	if err != nil {
		log.Println("export to json error", err)
		return
	}
	//log.Println(resJson)
	list := gjson.Get(resJson, "data").Array()
	count := 0
	for _, item := range list {
		if count >= 10 {
			break
		}
		timestamp := item.Get("updateTime").Int()
		if timestamp == 0 {
			history = append(history, "")
			continue
		} else {
			history = append(history, time.Unix(timestamp, 0).Format("2006年1月2日"))
		}
		history = append(history, item.Get("price").String())
		count++
	}
	return
}

func newReq(url string, json bool) *http.Request {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.8,zh-TW;q=0.7,zh-HK;q=0.5,en-US;q=0.3,en;q=0.2")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/115.0")
	if json {
		req.Header.Set("Accept", "application/json, text/plain, */*")
	} else {
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
	}

	return req
}

func readBody(resp *http.Response) string {
	body := resp.Body
	var err error
	if strings.ToLower(resp.Header.Get("Content-Encoding")) == "gzip" {
		body, err = gzip.NewReader(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
	}
	if body == nil {
		return ""
	}
	rawBody, err := io.ReadAll(body)
	return string(rawBody)
}
