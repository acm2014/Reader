package ispider

import (
	"testing"
)

func TestBook_SearchBook(t *testing.T) {
	b := Book{
		Name: "完美世界",
	}
	b.BiQuGeTwInit()
	b.SearchBook()
}
