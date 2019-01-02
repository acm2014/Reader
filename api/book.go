package api

import (
	"Reader/ispider"
	"Reader/tools"
	"net/http"
	"strconv"

	"github.com/didi/gendry/builder"
	"github.com/didi/gendry/scanner"
	"github.com/gin-gonic/gin"
)

const bookRouter = "/book"

type Book ispider.Book

type Response struct {
	ErrCode    int         `json:"err_code"`
	ErrMessage string      `json:"err_message"`
	Result     interface{} `json:"result"`
}

//获取书籍列表
func (b Book) GetBookList(c *gin.Context) {
	//生成mysql对象
	db, err := tools.NewMysqlConnection()
	if err != nil {
		tools.SystemOutput.Error("数据库连接失败", err)
		c.JSON(http.StatusOK, Response{ErrCode: 1, ErrMessage: "数据库连接失败"})
		return
	}
	//defer db.Close()
	//构建生成SQL语句
	table := "book"
	selectFields := []string{"name", "author", "image", "book_path"}
	cond, values, err := builder.BuildSelect(table, nil, selectFields)
	rows, err := db.Query(cond, values...)
	if err != nil {
		tools.SystemOutput.Error("数据库查询出错 ", err)
		c.JSON(http.StatusOK, Response{ErrCode: 1, ErrMessage: "数据库查询出错"})
		return
	}
	type book struct {
		Name     string `json:"name" ddb:"name"`
		ImageUrl string `json:"image_url" ddb:"image"`
		Author   string `json:"author" ddb:"author"`
		BookUrl  string `json:"book_url" ddb:"book_path"`
	}
	var bookList []book
	err = scanner.Scan(rows, &bookList)
	if err != nil {
		tools.SystemOutput.Error("数据库读取出错 ", err)
		c.JSON(http.StatusOK, Response{ErrCode: 1, ErrMessage: "数据库读取出错"})
		return
	}
	for i := range bookList {
		bookList[i].BookUrl = bookRouter + "/" + bookList[i].BookUrl
	}

	c.JSON(http.StatusOK, Response{
		ErrCode:    0,
		ErrMessage: "success",
		Result:     bookList,
	})
	return
}

//搜索书籍
func (b Book) SearchBook(c *gin.Context) {

}

//获取本书所有章节, 包括书籍信息
func (b Book) GetChapterList(c *gin.Context) {
	bookId := c.Param("bookId")
	db, err := tools.NewMysqlConnection()

	if err != nil {
		tools.SystemOutput.Error("数据库连接失败", err)
		c.JSON(http.StatusOK, Response{ErrCode: 1, ErrMessage: "数据库连接失败"})
	}
	//defer db.Close()
	sql := "select id,`name`,image,category,author,latest_chapters,abstract,source from book where book_path = ?"
	bookRows, err := db.Query(sql, bookId)
	if err != nil {
		tools.SystemOutput.Error("数据库查询出错 ,sql =", sql, err)
		c.JSON(http.StatusOK, Response{ErrCode: 1, ErrMessage: "数据库查询出错"})
		return
	}
	defer bookRows.Close()
	type (
		chapter struct {
			ChapterName string `json:"chapter_name"`
			ChapterId   int    `json:"chapter_id"`
		}
		bookDetail struct {
			BookName       string    `json:"book_name"`
			ImageUrl       string    `json:"image_url"`
			Category       string    `json:"category"`
			Author         string    `json:"author"`
			LatestChapters string    `json:"latest_chapters"`
			Abstract       string    `json:"abstract"`
			ChapterList    []chapter `json:"chapter_list"`
		}
	)
	aBook := bookDetail{}
	var id int
	for bookRows.Next() {
		var source string
		err := bookRows.Scan(&id, &aBook.BookName, &aBook.ImageUrl, &aBook.Category, &aBook.Author, &aBook.LatestChapters, &aBook.Abstract, &source)
		if err != nil {
			tools.SystemOutput.Error("数据库读取出错 ,sql =", sql, err)
			c.JSON(http.StatusOK, Response{ErrCode: 1, ErrMessage: "数据库读取出错"})
			return
		}
		break
	}
	if id == 0 {
		tools.SystemOutput.Info("找不到对应的书籍")
		c.JSON(http.StatusOK, Response{ErrCode: 1, ErrMessage: "找不到对应的书籍"})
		return
	}

	sql = "select `name`, chapter_id from chapter where book_id = ? order by chapter_id asc"
	chapterRows, err := db.Query(sql, id)
	if err != nil {
		tools.SystemOutput.Error("数据库查询出错 ,sql =", sql, err)
		c.JSON(http.StatusOK, Response{ErrCode: 1, ErrMessage: "数据库查询出错"})
		return
	}
	defer chapterRows.Close()

	chapterList := make([]chapter, 0)
	for chapterRows.Next() {
		var name string
		var id int
		err := chapterRows.Scan(&name, &id)
		if err != nil {
			tools.SystemOutput.Error("数据库读取出错 ,sql =", sql, err)
			c.JSON(http.StatusOK, Response{ErrCode: 1, ErrMessage: "数据库读取出错"})
			return
		}
		chapterList = append(chapterList, chapter{
			ChapterName: name,
			ChapterId:   id,
		})
	}
	aBook.ChapterList = chapterList
	c.JSON(http.StatusOK, Response{
		ErrCode:    0,
		ErrMessage: "success",
		Result:     aBook,
	})
}

//获取章节具体信息
func (b Book) GetChapterById(c *gin.Context) {
	bookId := c.Param("bookId")
	chapterId, err := strconv.Atoi(c.Param("chapterId"))
	if err != nil || chapterId < 0 {
		tools.SystemOutput.Info("找不到对应的章节")
		c.JSON(http.StatusOK, Response{ErrCode: 1, ErrMessage: "找不到对应的章节"})
		return
	}
	db, err := tools.NewMysqlConnection()
	if err != nil {
		tools.SystemOutput.Error("数据库连接失败", err)
		c.JSON(http.StatusOK, Response{ErrCode: 1, ErrMessage: "数据库连接失败"})
	}
	//defer db.Close()

	sql := "select id from book where book_path = ?"
	bookRows, err := db.Query(sql, bookId)
	if err != nil {
		tools.SystemOutput.Error("数据库查询出错 ,sql =", sql, err)
		c.JSON(http.StatusOK, Response{ErrCode: 1, ErrMessage: "数据库查询出错"})
		return
	}
	defer bookRows.Close()
	var id int
	if bookRows.Next() {
		err := bookRows.Scan(&id)
		if err != nil {
			tools.SystemOutput.Error("数据库读取出错 ,sql =", sql, err)
			c.JSON(http.StatusOK, Response{ErrCode: 1, ErrMessage: "数据库读取出错"})
			return
		}
	} else {
		tools.SystemOutput.Info("找不到对应的书籍")
		c.JSON(http.StatusOK, Response{ErrCode: 1, ErrMessage: "找不到对应的书籍"})
		return
	}

	sql = "select `name`, content from chapter where book_id = ? and chapter_id = ?"
	chapterRows, err := db.Query(sql, id, chapterId)
	if err != nil {
		tools.SystemOutput.Error("数据库查询出错 ,sql =", sql, err)
		c.JSON(http.StatusOK, Response{ErrCode: 1, ErrMessage: "数据库查询出错"})
		return
	}

	type chapter struct {
		Content       string `json:"content"`
		Name          string `json:"name"`
		PreChapterId  int    `json:"pre_chapter_id"`
		NextChapterId int    `json:"next_chapter_id"`
	}
	var aChapter chapter
	if chapterRows.Next() {
		err = chapterRows.Scan(&aChapter.Name, &aChapter.Content)
		if err != nil {
			tools.SystemOutput.Error("数据库读取出错 ,sql =", sql, err)
			c.JSON(http.StatusOK, Response{ErrCode: 1, ErrMessage: "数据库读取出错"})
			return
		}
	} else {
		tools.SystemOutput.Info("找不到对应的章节")
		c.JSON(http.StatusOK, Response{ErrCode: 1, ErrMessage: "找不到对应的章节"})
		return
	}
	aChapter.NextChapterId = chapterId + 1
	aChapter.PreChapterId = chapterId - 1
	c.JSON(http.StatusOK, Response{
		ErrCode:    0,
		ErrMessage: "success",
		Result:     aChapter,
	})
}
