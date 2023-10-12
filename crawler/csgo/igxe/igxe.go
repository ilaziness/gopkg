package igxe

import (
	"compress/gzip"
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/tidwall/gjson"
)

type Igxe struct {
	httpClient       *http.Client
	serverTimestamp  int64
	pageCount        int
	csrfToken        string
	csrfTokenSetTime int64
	igxeSellFile     *os.File
	igxeLeaseFile    *os.File
	existsIdList     map[string]string
	igxeSellChan     chan string
	igxeLeaseChan    chan string
}

func NewIgxe() *Igxe {
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: nil})
	if err != nil {
		log.Fatal(err)
	}
	c := http.DefaultClient
	c.Jar = jar
	c.Timeout = time.Second * 15
	return &Igxe{
		httpClient:   c,
		existsIdList: make(map[string]string),
	}
}

func (i *Igxe) Fetch() {
	i.setParams()
	i.getIgexSell()
}

func (i *Igxe) setParams() {
	listMainUrl := "https://www.igxe.cn/market/csgo?sort=3"
	req := i.newReqest(listMainUrl, true)
	resp, err := i.httpClient.Do(req)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Fatal(resp.Status)
	}
	content, _ := io.ReadAll(resp.Body)
	regx := regexp.MustCompile(`timestamp:\s*?(\d+)`)
	page := regexp.MustCompile(`page_count:\s*?(\d+)`)

	// 服务器时间戳
	result := regx.FindStringSubmatch(string(content))
	if len(result) != 2 {
		log.Fatal("timestamp not found")
	}
	i.serverTimestamp, err = strconv.ParseInt(result[1], 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	// 总页数
	p := page.FindStringSubmatch(string(content))
	if len(p) != 2 {
		log.Fatal("page_count not found")
	}
	i.pageCount, err = strconv.Atoi(p[1])
	if err != nil {
		log.Fatal(err)
	}
	log.Println("timestamp:", i.serverTimestamp, "pageCount:", i.pageCount)
}

func (i *Igxe) getIgexSell() {
	var err error
	var fileName = "igxe_sell.csv"
	var leaseFile = "igxe_lease.csv"
	i.igxeSellFile, err = os.OpenFile(fileName, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0755)
	if err != nil {
		log.Fatal(err)
	}
	defer i.igxeSellFile.Close()

	i.igxeLeaseFile, err = os.OpenFile(leaseFile, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0755)
	if err != nil {
		log.Fatal(err)
	}
	defer i.igxeLeaseFile.Close()

	for page := 1; page <= i.pageCount; page++ {
		i.getList(page)
	}
}

func (i *Igxe) getList(page int) {
	if page%10 == 0 {
		log.Println("sleep 5s")
		time.Sleep(time.Second * 5)
	}
	reqUrl := fmt.Sprintf("https://www.igxe.cn/api/v2/product/search/730?app_id=730&sort=3&page_no=%d&page_size=20", page)
	log.Println(reqUrl)
	req := i.newReqest(reqUrl, false)
	if page <= 2 {
		req.Header.Set("Referer", "https://www.igxe.cn/market/csgo?sort=3")
	} else {
		req.Header.Set("Referer", fmt.Sprintf("https://www.igxe.cn/market/csgo?sort=3&page_no=%d&page_size=20", page-1))
	}
	resp, err := i.httpClient.Do(req)
	if err != nil {
		log.Println(err, reqUrl)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Println(resp.Status)
		return
	}
	rawBody := i.readBody(resp)
	if gjson.Get(rawBody, "code").Int() != 200 {
		log.Println("response data is exception", rawBody)
	}
	data := gjson.Get(rawBody, "data.data").Array()
	if len(data) == 0 {
		log.Println("data is empty")
		return
	}
	rowData := make(map[string]string)
	for _, row := range data {
		time.Sleep(time.Millisecond * 700)

		productId := row.Get("id").String()
		if _, ok := i.existsIdList[productId]; ok {
			log.Println("skip exists id:", productId)
			continue
		}
		i.existsIdList[productId] = productId

		log.Println("product id:", productId)
		rowData["title"] = row.Get("title").String()        // 名称
		rowData["exterior"] = ""                            // 外观
		rowData["rarity"] = row.Get("rarity_name").String() // 品质
		rowData["sale_count"] = ""                          // 在售数量
		rowData["min_price"] = ""                           // 最低价
		if exterior := row.Get("exterior_name"); exterior.Exists() {
			rowData["exterior"] = exterior.String() // 外观
		}
		if saleCount := row.Get("sale_count"); saleCount.Exists() {
			rowData["sale_count"] = saleCount.String() // 在售数量
		}
		if minPrice := row.Get("min_price"); minPrice.Exists() {
			rowData["min_price"] = minPrice.String() // 最低价
		}
		//rowData["max_price"] = getSellMaxPrice(productId) // 最高价
		//log.Printf("%#v\n", rowData)
		line := []string{
			productId,
			rowData["title"],
			rowData["exterior"],
			rowData["rarity"],
		}

		// 在售
		linesell := append(line, rowData["sale_count"])
		linesell = append(linesell, rowData["min_price"])
		linesell = append(linesell, i.getSellHistory(productId)...)
		_, err = i.igxeSellFile.WriteString(strings.Join(linesell, ",") + "\n")
		if err != nil {
			log.Println(err)
		}

		//租赁
		long, short := i.getLeasePriceMin(productId)
		linelease := append(line, row.Get("lease_sale_count").String())
		linelease = append(linelease, short)
		linelease = append(linelease, long)
		linelease = append(linelease, i.getLeaseHistory(productId)...)
		_, err = i.igxeLeaseFile.WriteString(strings.Join(linelease, ",") + "\n")
		if err != nil {
			log.Println(err)
		}
	}
}

func (i *Igxe) newReqest(url string, html bool, method ...string) *http.Request {
	m := http.MethodGet
	if len(method) > 0 {
		m = method[0]
	}
	req, err := http.NewRequest(m, url, nil)
	if err != nil {
		log.Println(err, url)
		return nil
	}
	if html {
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
		req.Header.Set("Upgrade-Insecure-Requests", "1")
	} else {
		req.Header.Set("Accept", "application/json, text/plain, */*")
		req.Header.Set("Accept-Encoding", "gzip, deflate, br")
		req.Header.Set("Client-Type", "web")
		req.Header.Set("X-Requested-with", "XMLHttpRequest")
		random, timestamp, err := i.getRandomString(url)
		if err != nil {
			log.Println(err)
		} else {
			req.Header.Set("Random-String", random)
			req.Header.Set("Timestamp", strconv.FormatInt(timestamp, 10))
		}
	}
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.8,zh-TW;q=0.7,zh-HK;q=0.5,en-US;q=0.3,en;q=0.2")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/115.0")

	return req
}

func (i *Igxe) getRandomString(u string) (string, int64, error) {
	urlObj, err := url.Parse(u)
	if err != nil {
		return "", 0, err
	}
	now := time.Now().Add(30).Unix()
	str := fmt.Sprintf("%s:(%d)", urlObj.Path, now)
	return fmt.Sprintf("%x", md5.Sum([]byte(str))), now, nil
}

func (i *Igxe) readBody(resp *http.Response) string {
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

// 获取在售列表， 价格降序排序
//func getSellMaxPrice(productId string) string {
//	reqUrl := fmt.Sprintf("https://www.igxe.cn/product/trade/730/%s?sort=2&sort_rule=2", productId)
//	req := newReqest(reqUrl, false)
//	req.Header.Set("Referer", fmt.Sprintf("https://www.igxe.cn/product/730/%s?sticker_product_ids&sticker_slot&cur_page=1", productId))
//	resp, err := httpClient.Do(req)
//	if err != nil {
//		log.Println(reqUrl, err)
//		return ""
//	}
//	defer resp.Body.Close()
//	rawBody := readBody(resp)
//	price := gjson.Get(rawBody, "d_list.0.unit_price")
//	if price.Exists() {
//		return price.String()
//	}
//	return ""
//}

func (i *Igxe) getSellHistory(productId string) []string {
	retryCount := 0
	if i.csrfToken == "" {
		i.setCsrfToken(productId)
	}
	if time.Now().Unix()-i.csrfTokenSetTime > 3600*3.5 {
		i.setParams()
		i.setCsrfToken(productId)
	}

retry:
	req := i.newReqest(fmt.Sprintf("https://www.igxe.cn/product/get_product_sales_history/730/%s", productId), false, http.MethodPost)
	if i.csrfToken != "" {
		req.Header.Set("X-CSRFToken", i.csrfToken)
	}
	req.Header.Set("Referer", fmt.Sprintf("https://www.igxe.cn/product/730/%s?cur_page=3", productId))
	resp, err := i.httpClient.Do(req)
	if err != nil {
		log.Println(req.URL, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Println(resp.Status, req.Method, req.URL)
		time.Sleep(time.Second * 3)
		if retryCount <= 1 {
			retryCount++
			goto retry
		}
		return []string{}
	}
	rawBody := i.readBody(resp)
	list := gjson.Get(rawBody, "data").Array()
	if len(list) == 0 {
		log.Println("data is empty", req.URL)
		return []string{}
	}
	count := 0
	sells := make([]string, 0)
	for _, r := range list {
		if count >= 10 {
			break
		}
		sells = append(sells, r.Get("last_updated").String())
		sells = append(sells, r.Get("unit_price").String())
		count++
	}
	return sells
}

func (i *Igxe) setCsrfToken(productId string) {
	reqUrl := fmt.Sprintf("https://www.igxe.cn/product/730/%s?cur_page=3", productId)
	resp, err := i.httpClient.Do(i.newReqest(reqUrl, true))
	if err != nil {
		log.Println(reqUrl, err)
		return
	}
	defer resp.Body.Close()
	rawBody := i.readBody(resp)
	csrfRegex := regexp.MustCompile(`csrfmiddlewaretoken:\s*'(\w+)'`)
	result := csrfRegex.FindStringSubmatch(rawBody)
	if len(result) != 2 {
		log.Println("csrf token not found")
	}
	i.csrfToken = result[1]
	i.csrfTokenSetTime = time.Now().Unix()
}

func (i *Igxe) getLeasePriceMin(productId string) (long, short string) {
	retryCount := 0

retry:
	req := i.newReqest(fmt.Sprintf("https://www.igxe.cn/api/v2/lease/trade-list/730/%s?sort=3&sort_rule=3", productId), false)
	req.Header.Set("Referer", fmt.Sprintf("https://www.igxe.cn/product/730/%s?cur_page=6&sort_rule=1", productId))
	resp, err := i.httpClient.Do(req)
	if err != nil {
		log.Println(err, req.URL)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Println(resp.Status, req.URL)
		time.Sleep(time.Second * 3)
		if retryCount <= 1 {
			retryCount++
			goto retry
		}
		return
	}
	rawBody := i.readBody(resp)
	short = gjson.Get(rawBody, "data.rows.0.unit_price").String()

	// 长租
	req = i.newReqest(fmt.Sprintf("https://www.igxe.cn/api/v2/lease/trade-list/730/%s?sort=11&sort_rule=11", productId), false)
	req.Header.Set("Referer", fmt.Sprintf("https://www.igxe.cn/product/730/%s?cur_page=6&sort_rule=1", productId))
	resp, err = i.httpClient.Do(req)
	if err != nil {
		log.Println(err, req.URL)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Println(resp.Status, req.URL)
		return
	}
	rawBody = i.readBody(resp)
	long = gjson.Get(rawBody, "data.rows.0.long_term_price").String()
	return
}

func (i *Igxe) getLeaseHistory(productId string) []string {
	req := i.newReqest(fmt.Sprintf("https://www.igxe.cn/api/v2/lease/product-history/%s", productId), false)
	req.Header.Set("Referer", fmt.Sprintf("https://www.igxe.cn/product/730/%s?cur_page=3", productId))
	resp, err := i.httpClient.Do(req)
	if err != nil {
		log.Println(req.URL, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Println(resp.Status, req.URL)
		return []string{}
	}
	rawBody := i.readBody(resp)
	list := gjson.Get(rawBody, "data.rows").Array()
	if len(list) == 0 {
		log.Println("data is empty", req.URL)
		return []string{}
	}
	count := 0
	lease := make([]string, 0)
	for _, r := range list {
		if count >= 10 {
			break
		}
		lease = append(lease, r.Get("last_updated").String())
		lease = append(lease, r.Get("unit_price").String())
		lease = append(lease, r.Get("cash_pledge").String())
		count++
	}
	return lease
}
