package ispider

import (
	"fmt"
	"net/http"
)

func main() {
	res, err := http.Get("http://www.biquge.com.tw/")
	if err != nil {
		fmt.Println("http get failed", err)
	}
	fmt.Println(res)
}
