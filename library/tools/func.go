package tools

import (
	"fmt"

	"github.com/axgle/mahonia"
)

func Convert(src string, srcCode string, targetCode string) string {
	srcCoder := mahonia.NewDecoder(srcCode)
	srcResult := srcCoder.ConvertString(src)
	tagCoder := mahonia.NewDecoder(targetCode)
	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
	return string(cdata)
}

func ChineseGBKEncode(name string) string {
	enc := mahonia.NewEncoder("GBK")
	s := enc.ConvertString(name)
	res := ""
	for i := 0; i < len(s); i++ {
		res = res + fmt.Sprintf("%%%X", s[i])
	}
	return res
}
