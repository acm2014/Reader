package dao

type Book struct {
	ID             int    `gorm:"primary_key;column:id;type:int(10);not null" json:"-"`           // 书籍id,实际等于书名和source的hash
	Name           string `gorm:"column:name;type:varchar(50);not null" json:"name"`              // 书名
	Image          string `gorm:"column:image;type:varchar(100)" json:"image"`                    // 图片url
	Category       string `gorm:"column:category;type:varchar(50)" json:"category"`               // 书籍分类
	Author         string `gorm:"column:author;type:varchar(20)" json:"author"`                   // 作者
	LatestChapters string `gorm:"column:latest_chapters;type:varchar(50)" json:"latest_chapters"` // 最新章节
	Abstract       string `gorm:"column:abstract;type:text" json:"abstract"`                      // 书籍摘要
	BookPath       string `gorm:"column:book_path;type:varchar(50)" json:"book_path"`             // 数据的路径
	Source         string `gorm:"column:source;type:varchar(50)" json:"source"`                   // 爬虫的来源站点
}
