package main

import (
	"github.com/gin-gonic/gin"
	"main/web/auth"
	"net/http"
	"net/url"
)

func LoginAuth(c *gin.Context) {
	var username, _ = c.GetPostForm("username")
	var password, _ = c.GetPostForm("password")

	if err := auth.Auth(username, password); err == nil {
		location := url.URL{Path: "/home",}
    	c.Redirect(http.StatusFound, location.RequestURI())
		return
	} else {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{
			"error": err,
		})
		return
	}
}

func LogOut(c *gin.Context) {
	location := url.URL{Path: "/login",}
    c.Redirect(http.StatusFound, location.RequestURI())
}

func NavHome(c *gin.Context) {
	location := url.URL{Path: "/home",}
    c.Redirect(http.StatusFound, location.RequestURI())
}

func NavProfile(c *gin.Context) {
	location := url.URL{Path: "/profile",}
    c.Redirect(http.StatusFound, location.RequestURI())

}

func LoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

func MainPage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func ProfilePage(c *gin.Context) {
	c.HTML(http.StatusOK, "profile.html", nil)
}

func main() {
	server := gin.Default()
	server.LoadHTMLGlob("src/html/*")
	server.Static("/assets", "./src/assets")
	server.GET("/home", MainPage)
	server.GET("/profile", ProfilePage)
	server.GET("/login", LoginPage)
	server.POST("/login", LoginAuth)
	server.POST("/logout", LogOut)
	server.POST("/navprofile", NavProfile)
	server.POST("/navhome", NavHome)
	server.Run(":8888")
}