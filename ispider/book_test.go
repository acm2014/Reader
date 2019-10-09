package ispider

import (
	"reader/library/tools"
	"testing"
)

func TestBook_SearchBook(t *testing.T) {
	b := Book{
		Name: "完美世界",
	}
	b.BiQuGeTwInit()
	tools.Log.Error(b.SearchBook())
}
