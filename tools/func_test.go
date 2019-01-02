package tools

import (
	"fmt"
	"testing"
)

func TestGetDbConnection(t *testing.T) {
	for i := 0; i < 1000; i++ {
		db, err := NewMysqlConnection()
		fmt.Printf("%p %v\n", db, err)
	}
}

func TestChineseGBKEncode(t *testing.T) {
	fmt.Println(ChineseGBKEncode("永夜君王"))
}
