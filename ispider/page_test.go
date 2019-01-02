package ispider

import (
	"testing"
)

func TestPage_PageInit(t *testing.T) {
	page := Page{
		Host: biQuGeTw,
	}
	page.PageInit()
	//time.Sleep(time.Second * 5)
	//page.PageInit()

}
