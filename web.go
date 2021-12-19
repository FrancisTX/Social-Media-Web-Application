package main

import (
	"log"
	"net/http"
	"net/url"
	"time"
	"io/ioutil"
	"encoding/base64"
	"main/client"
	"html/template"
	"mime/multipart"

	"github.com/gin-gonic/gin"
)

var USERNAME, PROFILENAME string
var PROFILEIMG template.URL

func imgProcess(img *multipart.FileHeader) template.URL {
	imgfile, _ := img.Open()
	defer imgfile.Close()

	pimg, _ := ioutil.ReadAll(imgfile)
	var base64Encoding string
	imgType := http.DetectContentType(pimg)

	switch imgType {
	case "image/jpeg":
		base64Encoding += "data:image/jpeg;base64,"
	case "image/png":
		base64Encoding += "data:image/png;base64,"
	}

	base64Encoding += base64.StdEncoding.EncodeToString(pimg)
	return template.URL(base64Encoding)
}

func LoginAuth(c *gin.Context) {
	var username, _ = c.GetPostForm("username")
	var password, _ = c.GetPostForm("password")

	r, pimg := client.Login(map[string]string{"username": username, "password": password})

	if r["status"] == "Success" {
		USERNAME = r["username"]
		PROFILENAME = r["profilename"]
		PROFILEIMG = pimg
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
	var profileimg, _ = c.FormFile("profileimg")

	var encodedprofileimg template.URL

	if profileimg == nil {
		encodedprofileimg = template.URL("../assets/img/image.png")
	} else {
		encodedprofileimg = imgProcess(profileimg)
	}
	
	r := client.SignUp(map[string]string{"username": username, "password": password, "profilename": profilename,}, encodedprofileimg)

	if r["status"] == "Success" {
		c.HTML(http.StatusOK, "login.html", gin.H{"success": "User created! Please sign in.",})
	} else {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{"error": r["msg"],})
	}
}

func LogOut(c *gin.Context) {
	location := url.URL{Path: "/login"}
	c.Redirect(http.StatusFound, location.RequestURI())
}

func CreatePost(c *gin.Context) {
	var text = c.PostForm("content")
	var postimg, _ = c.FormFile("postimg")

	var encodedpostimg template.URL

	if postimg == nil {
		encodedpostimg = template.URL("")
	} else {
		encodedpostimg = imgProcess(postimg)
	}
	client.CreatePost(map[string]string{"username": USERNAME, "text": text, "time": time.Now().Format("2006-01-02 15:04:05")}, encodedpostimg)
	time.Sleep(1 * time.Second)
	posts, err := client.GetPosts(map[string]string{"username": USERNAME})
	if err == nil {
		c.HTML(http.StatusOK, "index.html", gin.H{"posts": posts, "curProfileimg": PROFILEIMG})
	} else {
		log.Print("Error in mainpage: ", err)
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
	c.HTML(http.StatusOK, "index.html", gin.H{"posts": posts, "curProfileimg": PROFILEIMG})
}

func ProfilePage(c *gin.Context) {
	userInfo, err := client.GetUserInfo(USERNAME)
	if err == nil {
		c.HTML(http.StatusOK, "profile.html", gin.H{"Profilename": userInfo.Profilename, "Username": userInfo.Username, "Profileimg": userInfo.Profileimg})
	}
}

func UserSearch(c *gin.Context) {
	username := c.PostForm("username")
	userInfo, err := client.GetUserInfo(username)
	if err != nil {
		c.HTML(http.StatusUnauthorized, "search.html", gin.H{"error": err})
	} else {
		c.HTML(http.StatusOK, "search.html", gin.H{"Profilename": userInfo.Profilename, "Username": userInfo.Username, "Profileimg": userInfo.Profileimg})
	}
}

func Follow(c *gin.Context) {
	username := c.PostForm("follow")
	_, err := client.Follow(USERNAME, username)
	if err != nil {
		userInfo, _ := client.GetUserInfo(username)
		log.Println("Unfollow failed: ", err.Error())
		c.HTML(http.StatusInternalServerError, "search.html", gin.H{"error": err, "Username": username, "Profilename": userInfo.Profilename, "Profileimg": userInfo.Profileimg})
	} else {
		time.Sleep(1 * time.Second)
		location := url.URL{Path: "/home"}
		c.Redirect(http.StatusFound, location.RequestURI())
	}
}

func Unfollow(c *gin.Context) {
	username := c.PostForm("unfollow")
	_, err := client.Unfollow(USERNAME, username)
	if err != nil {
		userInfo, _ := client.GetUserInfo(username)
		log.Println("Unfollow failed: ", err.Error())
		c.HTML(http.StatusInternalServerError, "search.html", gin.H{"error": err, "Username": username, "Profilename": userInfo.Profilename, "Profileimg": userInfo.Profileimg})
	} else {
		time.Sleep(1 * time.Second)
		location := url.URL{Path: "/home"}
		c.Redirect(http.StatusFound, location.RequestURI())
	}
	
}

func setUpRouter() *gin.Engine {
	server := gin.Default()
	server.LoadHTMLGlob("src/html/*")
	server.Static("/assets", "./src/assets")
	server.GET("/home", MainPage)
	server.GET("/profile", ProfilePage)
	server.GET("/login", LoginPage)
	server.POST("/login", LoginAuth)
	server.POST("/signup", SignUp)
	server.POST("/logout", LogOut)
	server.POST("/post", CreatePost)
	server.POST("/navprofile", NavProfile)
	server.POST("/navhome", NavHome)
	server.POST("/search", UserSearch)
	server.POST("/follow", Follow)
	server.POST("/unfollow", Unfollow)
	return server
}

func main() {
	server := setUpRouter()
	server.Run(":8888")
}
