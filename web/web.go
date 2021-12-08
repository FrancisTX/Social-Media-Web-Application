package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"main/client"

	"github.com/gin-gonic/gin"
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
	client.CreatePost(map[string]string{"username": USERNAME, "text": text, "img": "", "time": time.Now().String()})
	posts, err := client.GetPosts(map[string]string{"username": USERNAME})
	if err == nil {
		c.HTML(http.StatusOK, "index.html", gin.H{"posts": posts, "curProfileimg": PROFILEIMG})
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
	posts, err := client.GetPosts(map[string]string{"username": USERNAME})
	if err != nil {
		log.Print("Error in mainpage: ", err)
	}
	log.Println("Post: ", posts)
	c.HTML(http.StatusOK, "index.html", gin.H{"posts": posts, "curProfileimg": PROFILEIMG})
}

func ProfilePage(c *gin.Context) {
	userInfo, err := client.GetUserInfo(USERNAME)
	if err == nil {
		c.HTML(http.StatusOK, "profile.html", gin.H{"Profilename": userInfo.Profilename, "Username": userInfo.Username, "Profileimg": userInfo.Profileimg})
	}
}

func UserSearch(c *gin.Context) {
	usrname := c.Query("usrname")
	userInfo, err := client.GetUserInfo(usrname)
	if err != nil {
		c.HTML(http.StatusOK, "search.html", gin.H{"Error": "User Not Found"})
	} else {
		c.HTML(http.StatusOK, "search.html", gin.H{"Profilename": userInfo.Profilename, "Username": userInfo.Username, "Profileimg": userInfo.Profileimg})
	}
}

func Follow(c *gin.Context) {
	usrname := c.PostForm("follow")
	_, err := client.Follow(USERNAME, usrname)
	if err != nil {
		userInfo, _ := client.GetUserInfo(usrname)
		log.Println("Unfollow failed: ", err.Error())
		c.HTML(http.StatusOK, "search.html", gin.H{"error": err.Error(), "Username": usrname, "Profilename": userInfo.Profilename, "Profileimg": userInfo.Profileimg})
	}
	location := url.URL{Path: "/home"}
	c.Redirect(http.StatusFound, location.RequestURI())
}

func Unfollow(c *gin.Context) {
	usrname := c.PostForm("unfollow")
	_, err := client.Unfollow(USERNAME, usrname)
	if err != nil {
		userInfo, _ := client.GetUserInfo(usrname)
		log.Println("Unfollow failed: ", err.Error())
		c.HTML(http.StatusOK, "search.html", gin.H{"error": err.Error(), "Username": usrname, "Profilename": userInfo.Profilename, "Profileimg": userInfo.Profileimg})
	}
	location := url.URL{Path: "/home"}
	c.Redirect(http.StatusFound, location.RequestURI())
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
	server.GET("/search", UserSearch)
	server.POST("/follow", Follow)
	server.POST("/unfollow", Unfollow)
	//server.GET("/follower", GetFollower)
	//server.GET("/following", GetFollowing)
	server.Run(":8888")
}
