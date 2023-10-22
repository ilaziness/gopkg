package main

import (
	"context"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"log"
	"regexp"
	"strings"
	"time"
)

var regexPrice = regexp.MustCompile(`unit_price":.*?"([\d\\.]+)?"`)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
		chromedp.ExecPath("/usr/bin/google-chrome-stable"),
		chromedp.WindowSize(1400, 1200),
	)

	allocctx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	ctx, cancel := chromedp.NewContext(
		allocctx,
		// chromedp.WithDebugf(log.Printf),
	)
	ctx, cancel = context.WithTimeout(ctx, time.Second*20)
	defer cancel()

	chromedp.ListenTarget(ctx, func(event interface{}) {
		switch ev := event.(type) {
		case *network.EventResponseReceived:
			go func() {
				if !strings.HasPrefix(ev.Response.URL, "https://www.igxe.cn/product/trade/730") || strings.Index(ev.Response.URL, "sort") == -1 {
					return
				}
				log.Println(ev.Response.URL, strings.Index(ev.Response.URL, "sort") == -1)
				log.Println(!strings.HasPrefix(ev.Response.URL, "https://www.igxe.cn/product/trade/730") || strings.Index(ev.Response.URL, "sort") == -1)
				c := chromedp.FromContext(ctx)
				body, err := network.GetResponseBody(ev.RequestID).Do(cdp.WithExecutor(ctx, c.Target))
				if err != nil {
					log.Println(err)
					return
				}
				res := regexPrice.FindStringSubmatch(string(body))
				if len(res) == 2 {
					log.Println("response:", res[1])
				}
			}()
		}
	})

	// 模拟鼠标移动到元素上
	var res bool
	err := chromedp.Run(ctx,
		chromedp.Navigate("https://www.igxe.cn/product/730/657111?cur_page=1"),
		//chromedp.WaitReady("body"),
		chromedp.WaitReady(`div.relatedInfo`),
		chromedp.Evaluate(`document.querySelector('div.relatedInfo div.drop-down-menu2').dispatchEvent(new Event('mouseover'))`, &res),
	)
	if err != nil {
		log.Fatal(err)
	}
	var listContent string
	err = chromedp.Run(ctx,
		chromedp.Click(`div.relatedInfo div.drop-down-menu2 li.com-menu:nth-child(2)`, chromedp.ByQueryAll),
		//chromedp.Text(`table.sell-data span.price`, &listContent),
	)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(listContent)

	err = chromedp.Run(ctx,
		chromedp.Sleep(time.Second),
		chromedp.Click(`div.relatedInfo div.drop-down-menu2 li.com-menu:nth-child(3)`, chromedp.ByQueryAll),
		//.TextContent(`table.sell-data span.price`, &listContent),
	)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(listContent)

	_ = chromedp.Run(ctx, chromedp.Sleep(time.Second*3))
}
