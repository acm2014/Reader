package ispider

import (
	"errors"
	"fmt"
	"log"
	"reader/library/cache"
	"reader/library/dao"
	"reader/library/tools"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

type Book struct {
	Page
	Path           string
	Name           string
	Image          string
	Category       string
	Author         string
	LatestChapters string
	Abstract       string
}

type Chapter struct {
	Page
	Name    string
	Content string
	Path    string
}

func (b *Book) BiQuGeTwInit() {
	b.Host = biQuGeTw
	b.Page.Path = "modules/article/soshu.php"
	b.Param = map[string]string{"searchkey": "+" + tools.ChineseGBKEncode(b.Name)}
}

func (b *Book) BookTxtInit() {
	b.Host = dingDian
	b.Page.Path = "search.php"
	b.Param = map[string]string{"keyword": b.Name}
}

func (b *Book) Skip(doc *goquery.Document) *goquery.Document {
	if b.Host == biQuGeTw {
		return doc
	} else if b.Host == dingDian {
		var docSkip *goquery.Document
		doc.Find("div.result-item.result-game-item").Each(func(i int, s *goquery.Selection) {
			if docSkip != nil {
				return
			}
			title := s.Find("h3").Find("a")
			//fmt.Println(title.Attr("name"))
			if name, ok := title.Attr("title"); ok == true && name == b.Name {
				url, ok := title.Attr("href")
				fmt.Println(url)
				if ok == false {
					tools.Log.Error("没有找到书籍的url", b.Name)
					return
				} else {
					url = strings.TrimRight(url, "/")
					slices := strings.Split(url, "/")
					if len(slices) > 0 {
						b.Page.Path = slices[len(slices)-1]
						b.Param = nil
						log.Println(b)
					} else {
						tools.Log.Error("url 错误", url)
					}
					var err error

					docSkip, err = b.GetPage()

					if err != nil {
						tools.Log.Error("get page failed,", err)
						docSkip = nil
						return
					}
				}
			}
		})
		return docSkip
	} else {
		return nil
	}
}

func (b *Book) SearchBook() error {
	if b.Name == "" || b.Host == "" {
		tools.Log.Error("book name or page host can't be nil")
		return errors.New("book name or page host can't be nil")
	}
	doc, err := b.GetPage()
	if err != nil {
		tools.Log.Error("get page failed", err)
		return err
	}
	doc = b.Skip(doc)
	//fmt.Println(doc.Html())
	if doc == nil {
		tools.Log.Error("get page failed")
		return errors.New("get page failed")
	}
	b.getImage(doc)
	b.getInfo(doc)
	db, err := cache.NewMysql()
	if err != nil {
		tools.Log.Error("数据库连接失败", err)
		return err
	}
	var book dao.Book
	db = db.Where(&dao.Book{Name: b.Name, Source: b.Page.Path}).First(&book)
	if db.Error != nil && !db.RecordNotFound() {
		tools.Log.Error("get book info err", db.Error)
		return db.Error
	}
	if db.RecordNotFound() {
		book = dao.Book{
			Name:           b.Name,
			Image:          b.Image,
			Category:       b.Category,
			Author:         b.Author,
			LatestChapters: b.LatestChapters,
			Abstract:       b.Abstract,
			BookPath:       b.Path,
			Source:         b.Host,
		}
		db = db.Create(&book)
		if err = db.Error; err != nil {
			tools.Log.Error("数据库插入失败", err)
		}
	} else {
		tools.Log.Info("该书籍已经存在", b.Name, b.Image, b.Category, b.Author, b.LatestChapters, b.Abstract, b.Path, b.Host)
	}
	// TODO ,需要检验id是否更新
	b.getAllChapters(doc, int64(book.ID))
	return nil
}

//获取书籍封面图, 同时获取页面路径
func (b *Book) getImage(doc *goquery.Document) {
	url, ok := doc.Find("#fmimg").Find("img").Attr("src")
	if ok == false {
		b.Image = ""
	} else {
		b.Image = url
	}

	url, ok = doc.Find("#list").Find("a").First().Attr("href")
	if ok == false {
		b.Path = ""
	} else {
		url = strings.TrimLeft(url, "/")
		b.Path = strings.Split(url, "/")[0]
	}
	tools.Log.Info(b.Path)
	tools.Log.Info(b.Image)
}

//获取书籍信息: 类别, 作者, 最新章节, 摘要
func (b *Book) getInfo(doc *goquery.Document) {
	//tools.Log.Info(doc.Html())
	txt := doc.Find("div.con_top").Text()
	tools.Log.Info(txt)
	txt = strings.Replace(txt, " ", "", -1)
	b.Category = strings.Split(txt, ">")[1]
	info := doc.Find("div#maininfo").Find("#info").Find("p")
	b.Author = strings.TrimLeft(info.First().Text(), "作    者：")
	b.LatestChapters = info.Last().Find("a").Text()
	b.Abstract = doc.Find("div#maininfo").Find("#intro").Text()
	b.Abstract = strings.TrimLeft(b.Abstract, " ")
	b.Abstract = strings.TrimLeft(b.Abstract, "\n")
	b.Abstract = strings.TrimLeft(b.Abstract, "\t")
	b.Abstract = strings.TrimRight(b.Abstract, " ")
	b.Abstract = strings.TrimRight(b.Abstract, "\n")
	b.Abstract = strings.TrimRight(b.Abstract, "\t")
	tools.Log.Info(b.Category)
	tools.Log.Info(b.Author)
	tools.Log.Info(b.LatestChapters)
	tools.Log.Info(b.Abstract)
}

func (b *Book) getAllChapters(doc *goquery.Document, bookId int64) {
	var wg sync.WaitGroup
	doc.Find("#list").Find("a").Each(func(i int, s *goquery.Selection) {
		var c Chapter
		c.Name = s.Text()
		path, ok := s.Attr("href")
		if ok == true {
			c.Path = path
		} else {
			tools.Log.Error("can't get chapter's path chapter's name = ", c.Name)
			return
		}
		//页面跳转
		c.Page.Path = c.Path
		c.Page.Cookie = b.Page.Cookie
		c.Page.Host = b.Page.Host
		wg.Add(1)
		go c.GetChapter(bookId, &wg, int64(i))
	})
	wg.Wait()
}

func (c *Chapter) GetChapter(bookId int64, wg *sync.WaitGroup, chapterId int64) {
	defer wg.Done()
	doc, err := c.GetPage()
	if err != nil {
		tools.Log.Error("get page failed", err)
		return
	}
	c.Content, err = doc.Find("#content").Html()
	fmt.Println(c.Content, err)
	db, err := cache.NewMysql()
	if err != nil {
		tools.Log.Error("数据库连接失败", err)
		return
	}
	//defer db.Close()
	chapter := dao.Chapter{
		ID:        fmt.Sprintf("%d.%d", bookId, chapterId),
		BookID:    int(bookId),
		ChapterID: int(chapterId),
		Name:      c.Name,
		Content:   c.Content,
		Path:      c.Path,
		Source:    c.Host,
	}
	db.Create(&chapter)
	tools.Log.Debug("success", c.Name)
	return
}
