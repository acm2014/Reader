package dao

type Chapter struct {
	ID        string `gorm:"primary_key;column:id;type:varchar(50);not null" json:"-"` //
	BookID    int    `gorm:"column:book_id;type:int(10)" json:"book_id"`               // 书籍id
	ChapterID int    `gorm:"column:chapter_id;type:int(10)" json:"chapter_id"`         //
	Name      string `gorm:"column:name;type:varchar(50)" json:"name"`                 // 章节名
	Content   string `gorm:"column:content;type:longtext" json:"content"`              // 章节内容
	Path      string `gorm:"column:path;type:varchar(50)" json:"path"`                 // 章节路径url
	Source    string `gorm:"column:source;type:varchar(50)" json:"source"`             // 爬虫来源站点
}
