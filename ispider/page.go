package ispider

import (
	"Reader/tools"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

const biQuGeTw = "http://www.biquyun.com"
const dingDian = "http://www.booktxt.com"

var control chan struct{}

func init() {
	control = make(chan struct{}, 30)
}

type Page struct {
	Host   string
	Path   string
	Param  map[string]string
	Cookie string
}

func (p *Page) PageInit() (res string, err error) {
	if p.Host == "" {
		tools.SystemOutput.Error("page host 不能为空")
		return "", errors.New("page host 不能为空")
	}
	// create context
	ctxt, cancel := context.WithCancel(context.Background())
	defer cancel()

	// create chrome instance
	c, err := chromedp.New(ctxt, chromedp.WithLog(log.Printf))
	if err != nil {
		log.Fatal(err)
	}
	// run task list
	defer func() {
		// 关闭chrome实例
		err = c.Shutdown(ctxt)
		if err != nil {
			log.Fatal(err)
		}

		// 等待chrome实例关闭
		err = c.Wait()
		if err != nil {
			log.Fatal(err)
		}
	}()

	err = c.Run(ctxt, chromedp.Tasks{
		chromedp.Navigate(p.getUrl()),
		chromedp.WaitVisible(`div.footer_cont`, chromedp.ByQuery),
		chromedp.OuterHTML("body", &res),
		chromedp.ActionFunc(func(ctx context.Context, h cdp.Executor) error {
			cookies, err := network.GetAllCookies().Do(ctx, h)
			var c string
			for _, v := range cookies {
				if v.Name == "BAIDUID" {
					continue
				}
				c = c + v.Name + "=" + v.Value + ";"
			}
			p.Cookie = c
			log.Println(c)
			if err != nil {
				return err
			}
			return nil
		}),
	})
	if err != nil {
		log.Fatal(err)
	}
	return res, nil
}

func (p *Page) HttpGet() (res *http.Response, err error) {
	control <- struct{}{}
	defer func() {
		<-control
	}()
	//fmt.Printf("Page_Pointer %p\n", p)
	req, _ := http.NewRequest("GET", p.getUrl(), nil)
	req.Header.Add("Accept-Encoding", "gzip, deflate")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.110 Safari/537.36")
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")
	req.Header.Add("Accept-Language", "en-US,en;q=0.9")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Cookie", p.Cookie)
	req.Header.Add("Upgrade-Insecure-Requests", "1")

	res, err = http.DefaultClient.Do(req)
	if err != nil {
		tools.SystemOutput.Error("http get failed", err)
		return nil, err
	}
	return res, nil
}

func (p *Page) GetPage() (doc *goquery.Document, err error) {
	// 使用传统方式获取页面, 如果获取不到, 使用chromedp模拟浏览器行为
	res, err := p.HttpGet()
	if err != nil || res.StatusCode == 521 {
		if err != nil {
			fmt.Println("???????", err)
		} else {
			fmt.Println("??????", res.StatusCode)
		}
		if len(p.Param) == 0 {
			return nil, err
		}
		tools.SystemOutput.Info("method a")
		res, err := p.PageInit()
		if err != nil {
			return nil, err
		} else {
			return goquery.NewDocumentFromReader(strings.NewReader(res))
		}
	} else {
		for err != nil || res.StatusCode != 200 {
			if res != nil {
				tools.SystemOutput.Info(res.StatusCode)
			}
			time.Sleep(time.Second)
			res, err = p.HttpGet()
		}
		defer res.Body.Close()
		tools.SystemOutput.Info("method b")
		reader, err := gzip.NewReader(res.Body)
		if err != nil {
			tools.SystemOutput.Error("gzip decode failed", err)
			return nil, err
		}
		body, err := ioutil.ReadAll(reader)
		if err != nil {
			tools.SystemOutput.Error("read html failed", err)
			return nil, err
		}
		if p.Host == dingDian && len(p.Param) != 0 {
			return goquery.NewDocumentFromReader(strings.NewReader(string(body)))
		} else {
			return goquery.NewDocumentFromReader(strings.NewReader(tools.Convert(string(body), "GBK", "UTF-8")))
		}
	}
}

func (p *Page) getUrl() string {
	u := p.Host + "/" + p.Path + "?"
	for k, v := range p.Param {
		u = u + k + "=" + v
		u += "&"
	}
	u = strings.TrimRight(u, "&")
	u = strings.TrimRight(u, "?")
	u = strings.TrimRight(u, "/")
	if len(p.Param) == 0 && strings.Contains(u, ".html") == false {
		u += "/"
	}
	tools.SystemOutput.Info("u-sss", u)
	us, _ := url.Parse(u)
	q := us.Query()
	us.RawQuery = q.Encode() //urlEncode
	tools.SystemOutput.Info("request", us.String())
	return us.String()
}
