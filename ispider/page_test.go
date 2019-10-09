package ispider

import (
	"reader/library/tools"
	"testing"
)

func TestPage_PageInit(t *testing.T) {
	page := Page{
		Host: biQuGeTw,
	}
	tools.Log.Debug(page.PageInit())
	//time.Sleep(time.Second * 5)
	//page.PageInit()

}
