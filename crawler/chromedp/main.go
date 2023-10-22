package main

import (
	"context"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"log"
	"strconv"
	"strings"
	"time"
)

var Host = "https://www.igxe.cn"

var cookies []*network.Cookie

func main() {
	opts := append(
		chromedp.DefaultExecAllocatorOptions[:],
		chromedp.ExecPath("C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe"),
		chromedp.Flag("headless", false),
		chromedp.Flag("hide-scrollbars", false),
		//chromedp.WindowSize(1920, 1080),
	)
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	// create chrome instance
	ctx, cancel := chromedp.NewContext(
		allocCtx,
		//chromedp.WithDebugf(log.Printf),
	)
	defer cancel()

	chromedp.ListenTarget(ctx, func(event interface{}) {
		switch ev := event.(type) {
		case *network.EventResponseReceived:
			if !strings.HasPrefix(ev.Response.URL, "https://www.igxe.cn/api/v2/product/search/730?app_id=730") ||
				!strings.HasPrefix(ev.Response.URL, "https://www.igxe.cn/product/trade/730") {
				return
			}
			log.Println(ev.Response.URL, ev.RequestID)
			go func() {
				c := chromedp.FromContext(ctx)
				body, err := network.GetResponseBody(ev.RequestID).Do(cdp.WithExecutor(ctx, c.Target))
				if err != nil {
					log.Println("get body error:", err)
					return
				}
				log.Println(string(body))
			}()
		}
	})

	// create a timeout
	ctx, cancel = context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	// 首页
	list := ""
	pageCount := ""
	err := chromedp.Run(ctx,
		chromedp.Navigate(`https://www.igxe.cn/market/csgo?sort=3`),
		// wait for footer element is visible (ie, page is loaded)
		chromedp.WaitReady(`body`),
		// click the button to open the dialogchromedp.Click("ul.el-pager > li:nth-child(2)", chromedp.NodeVisible),
		// retrieve the text of the textarea
		chromedp.InnerHTML(`div.list`, &list),
		chromedp.Text("ul.el-pager > li:last-child", &pageCount),
	)

	c := chromedp.FromContext(ctx)
	cookies, err = network.GetCookies().Do(cdp.WithExecutor(ctx, c.Target))
	if err != nil {
		log.Println("get cookies error:", err)
	}
	if len(cookies) > 0 {
		log.Printf("cookies: %#v", cookies[0])
	}

	if err != nil {
		log.Fatalln(err)
	}
	if pageCount == "" {
		log.Fatalln("获取总页数错误")
	}
	tPage, err := strconv.Atoi(strings.TrimSpace(pageCount))
	if err != nil {
		log.Fatalln("转换数据类型错误", pageCount, err)
	}

	mainList := parseList(list)
	if len(mainList) == 0 {
		log.Fatalln("获取首页数据为空")
	}

	getDetail(ctx, mainList[0])
	log.Printf("首页数据：%d", len(parseList(list)))

	for p := 2; p <= tPage; p++ {
		getList(ctx, p)
	}

	//quit := make(chan os.Signal)
	//signal.Notify(quit, os.Interrupt)
	//<-quit
	//cancel()
	//log.Println("Shutdown ...")
}

func getList(ctx context.Context, no int) {

}

func getDetail(ctx context.Context, link string) {
	//q := strings.Split(link, "?")
	//p := strings.Split(q[0], "/")
	//productId := p[len(p)-1]

	infobok := ""
	partsbok := "" //皮肤外观
	sealNum := ""
	err := chromedp.Run(ctx,
		chromedp.Navigate(Host+link),
		chromedp.WaitReady("body"),
		chromedp.InnerHTML("div.productInfo div.txt", &infobok),
		chromedp.InnerHTML("div.relatedInfo > div:first-child", &sealNum),
	)
	if err != nil {
		log.Println(err)
		return

	}
	err = chromedp.Run(ctx,
		chromedp.WaitVisible(`div.parts-bok`),
		chromedp.Text("div.parts-bok > a.select", &partsbok),
	)
	if err != nil {
		log.Println("获取皮肤外观错误", err)
	}

	//最低价
	var res bool
	var selldata string
	err = chromedp.Run(ctx,
		//div.productDetailsCn span.steam-database--select
		chromedp.Evaluate(`document.querySelector('div.productDetailsCn span.steam-database--select > div.drop-down-menu2').dispatchEvent(new Event('mouseover'))`, &res),
		//chromedp.WaitVisible(`div.productDetailsCn > div.mod-dropMenu > div.leftItems > div.relatedInfo div.csgoInfo2 > div.filter com-ul`),
		chromedp.Click(`div.filter ul.com-ul > li.com-menu:nth-child(2)`),
		chromedp.InnerHTML(`div.sell-data`, &selldata),
	)
	if err != nil {
		log.Println("获取皮肤外观错误", err)
	}

	log.Println(infobok)
	log.Println(sealNum)
	log.Println(partsbok)
	log.Println(res, selldata)
	name := nameRegexp.FindStringSubmatch(infobok)
	log.Printf("%#v", name[1])
	rarity := rarityRegexp.FindStringSubmatch(infobok)
	log.Printf("%#v", rarity[1])
	onSeal := onsealRegexp.FindStringSubmatch(sealNum)
	log.Printf("%#v", onSeal[1])
}
