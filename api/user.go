package api

import "github.com/gin-gonic/gin"

type User struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

func (u User) Register(c gin.Context) {

}

func (u User) Login(c gin.Context) {

}
