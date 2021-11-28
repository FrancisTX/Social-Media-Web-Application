package main

import (
	"github.com/gin-gonic/gin"
	"main/web/auth"
	"net/http"
	"net/url"
	"errors"
)

func LoginAuth(c *gin.Context) {
	var (
		username string
		password string
	)
	if in, isExist := c.GetPostForm("username"); isExist && in != "" {
		username = in
	} else {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{
			"error": errors.New("Must enter username"),
		})
		return
	}
	if in, isExist := c.GetPostForm("password"); isExist && in != "" {
		password = in
	} else {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{
			"error": errors.New("Must enter password"),
		})
		return
	}
	if err := auth.Auth(username, password); err == nil {
		//c.HTML(http.StatusOK, "login.html", gin.H{
		//	"success": "Login successful",
		//})
		location := url.URL{Path: "/",}
    	c.Redirect(http.StatusFound, location.RequestURI())
		return
	} else {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{
			"error": err,
		})
		return
	}
}

func LoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

func MainPage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func main() {
	server := gin.Default()
	server.LoadHTMLGlob("src/html/*")
	server.Static("/assets", "./src/assets")
	server.GET("/", MainPage)
	server.GET("/login", LoginPage)
	server.POST("/login", LoginAuth)
	server.Run(":8888")
}