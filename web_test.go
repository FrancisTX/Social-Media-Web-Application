package main

import (
	"bytes"
	"net/http"
    "net/url"
	"net/http/httptest"
	"testing"
	"github.com/stretchr/testify/assert"
    "time"
    "strconv"
)

var server = setUpRouter()

func generateReq(method string, loc string, data url.Values) (*httptest.ResponseRecorder, *http.Request) {
	b := bytes.NewBufferString(data.Encode())
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(method, loc, b)
	if (b!=nil) {req.Header.Add("Content-Type", "application/x-www-form-urlencoded")}
	return w, req
}

func TestHomeGet(t *testing.T) {
	w, req := generateReq("GET", "/home", nil)
	server.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestLoginGet(t *testing.T) {
	w, req := generateReq("GET", "/login", nil)
	server.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestProfileGet(t *testing.T) {
	w, req := generateReq("GET", "/home", nil)
	server.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestNavProfile(t *testing.T) {
	w, req := generateReq("POST", "/navprofile", nil)
	server.ServeHTTP(w, req)

	assert.Equal(t, 302, w.Code)
}

func TestNavHome(t *testing.T) {
	w, req := generateReq("POST", "/navhome", nil)
	server.ServeHTTP(w, req)

	assert.Equal(t, 302, w.Code)
}

func TestLogOut(t *testing.T) {
	w, req := generateReq("POST", "/logout", nil)
	server.ServeHTTP(w, req)

	assert.Equal(t, 302, w.Code)
}

func TestLoginPost(t *testing.T) {
	// Not existed user
	data := url.Values{}
	data.Set("username", time.Now().String())
	data.Set("password", "abcd")
	w, req := generateReq("POST", "/login", data)
	server.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
	assert.Contains(t, w.Body.String(), "user does not exist")

	// Wrong password
	data = url.Values{}
	data.Set("username", "bot")
	data.Set("password", "abcd")
	w, req = generateReq("POST", "/login", data)
	server.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
	assert.Contains(t, w.Body.String(), "password is not correct")

	// Correct user
	data = url.Values{}
	data.Set("username", "bot")
	data.Set("password", "bot")
	w, req = generateReq("POST", "/login", data)
	server.ServeHTTP(w, req)

	assert.Equal(t, 302, w.Code)
}

func TestSignUpPost(t *testing.T) {

	newuser := strconv.FormatInt(time.Now().UnixNano(), 10)
	data := url.Values{}
	data.Set("username", newuser)
	data.Set("profilename", newuser)
	data.Set("password", newuser)
	w, req := generateReq("POST", "/signup", data)
	server.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	// sign up again using same user profile
	time.Sleep(1 * time.Second)
	w, req = generateReq("POST", "/signup", data)
	server.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
	assert.Contains(t, w.Body.String(), "User already exists")
}

func TestCreatePost(t *testing.T) {
	// Login testing account
	data := url.Values{}
	data.Set("username", "bot")
	data.Set("password", "bot")
	w, req := generateReq("POST", "/login", data)
	server.ServeHTTP(w, req)

	assert.Equal(t, 302, w.Code)

	// create post
	newcontent := "Test: " + strconv.FormatInt(time.Now().UnixNano(), 10)
	data = url.Values{}
	data.Set("content", newcontent)
	w, req = generateReq("POST", "/post", data)
	server.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), newcontent)
}

func TestUserSearch(t *testing.T) {
	data := url.Values{}
	data.Set("username", "bot")
	w, req := generateReq("POST", "/search", data)
	server.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	// Search non-existed user
	newuser := strconv.FormatInt(time.Now().UnixNano(), 10)
	data = url.Values{}
	data.Set("username", newuser)
	w, req = generateReq("POST", "/search", data)
	server.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
	assert.Contains(t, w.Body.String(), "user not found")
}

func TestFollow(t *testing.T) {
	// Login testing account
	data := url.Values{}
	data.Set("username", "bot")
	data.Set("password", "bot")
	w, req := generateReq("POST", "/login", data)
	server.ServeHTTP(w, req)

	assert.Equal(t, 302, w.Code)

	// follow
	data = url.Values{}
	data.Set("follow", "bot2")
	w, req = generateReq("POST", "/follow", data)
	server.ServeHTTP(w, req)

	assert.Equal(t, 302, w.Code)

	// follow again
	w, req = generateReq("POST", "/follow", data)
	server.ServeHTTP(w, req)

	assert.Equal(t, 500, w.Code)
	assert.Contains(t, w.Body.String(), "bot has already followed bot2")
}

func TestUnFollow(t *testing.T) {
	// Login testing account
	data := url.Values{}
	data.Set("username", "bot")
	data.Set("password", "bot")
	w, req := generateReq("POST", "/login", data)
	server.ServeHTTP(w, req)

	assert.Equal(t, 302, w.Code)

	// unfollow
	data = url.Values{}
	data.Set("unfollow", "bot2")
	w, req = generateReq("POST", "/unfollow", data)
	server.ServeHTTP(w, req)

	assert.Equal(t, 302, w.Code)

	// unfollow again
	w, req = generateReq("POST", "/unfollow", data)
	server.ServeHTTP(w, req)

	assert.Equal(t, 500, w.Code)
	assert.Contains(t, w.Body.String(), "bot is not following bot2")
}