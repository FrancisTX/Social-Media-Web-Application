package main

import (
	"fmt"
	"net/http"
	"net/url"
<<<<<<< HEAD

	"main/client"

	"github.com/gin-gonic/gin"
=======
	"time"
	"fmt"
>>>>>>> 95eadfe4f3543d74eafd524d5fb5b0005c7c9e94
)

var USERNAME = ""
var PROFILENAME = ""
var PROFILEIMG = ""

func LoginAuth(c *gin.Context) {
	var username, _ = c.GetPostForm("username")
	var password, _ = c.GetPostForm("password")

	r := client.Login(map[string]string{"username": username, "password": password})

	if r["status"] == "Success" {
		USERNAME = r["username"]
		PROFILENAME = r["profilename"]
		PROFILEIMG = r["profileimg"]
		location := url.URL{Path: "/home"}
		c.Redirect(http.StatusFound, location.RequestURI())
		return
	} else {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{
			"error": r["msg"],
		})
		return
	}
}

func SignUp(c *gin.Context) {
	var username, _ = c.GetPostForm("username")
	var password, _ = c.GetPostForm("password")
	var profilename, _ = c.GetPostForm("profilename")
	var profileimg, _ = c.GetPostForm("profileimg")

	r := client.SignUp(map[string]string{"username": username, "password": password, "profilename": profilename, "profileimg": profileimg})

	if r["status"] == "Success" {
		USERNAME = r["username"]
		PROFILENAME = r["profilename"]
		PROFILEIMG = r["profileimg"]
		c.HTML(http.StatusOK, "login.html", gin.H{
			"success": "User created! Please sign in.",
		})
		return
	} else {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{
			"error": r["msg"],
		})
		return
	}
}

func LogOut(c *gin.Context) {
	fmt.Println(USERNAME, PROFILENAME, PROFILEIMG)
	location := url.URL{Path: "/login"}
	c.Redirect(http.StatusFound, location.RequestURI())
}

func CreatePost(c *gin.Context) {
	var text = c.PostForm("content")
	client.CreatePost(map[string]string{"username":USERNAME, "profilename":PROFILENAME, "profileimg": PROFILEIMG, "text": text, "img": "", "time": time.Now().String()})
	posts, err := client.GetPosts(map[string]string{"username":USERNAME})
	if err == nil {
		c.HTML(http.StatusOK, "index.html", gin.H{"posts": posts, "curProfileimg": PROFILEIMG,})
	}
}

func NavHome(c *gin.Context) {
	location := url.URL{Path: "/home"}
	c.Redirect(http.StatusFound, location.RequestURI())
}

func NavProfile(c *gin.Context) {
	location := url.URL{Path: "/profile"}
	c.Redirect(http.StatusFound, location.RequestURI())

}

func LoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

func MainPage(c *gin.Context) {
	posts, err := client.GetPosts(map[string]string{"username":USERNAME})
	if err == nil {
		c.HTML(http.StatusOK, "index.html", gin.H{"posts": posts, "curProfileimg": PROFILEIMG,})
	}
}

func ProfilePage(c *gin.Context) {
	c.HTML(http.StatusOK, "profile.html", nil)
}

func main() {
	server := gin.Default()
	server.LoadHTMLGlob("web/src/html/*")
	server.Static("/assets", "./web/src/assets")
	server.GET("/home", MainPage)
	server.GET("/profile", ProfilePage)
	server.GET("/login", LoginPage)
	server.POST("/login", LoginAuth)
	server.POST("/signup", SignUp)
	server.POST("/logout", LogOut)
	server.POST("/post", CreatePost)
	server.POST("/navprofile", NavProfile)
	server.POST("/navhome", NavHome)
	server.Run(":8888")
}
