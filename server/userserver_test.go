package main

import (
	"testing"
	"context"
	pb "main/proto"
	"github.com/stretchr/testify/assert"
	"time"
    "strconv"
	"html/template"
	"errors"
)


var	c = UserServer{}

func TestLogin(t *testing.T) {
	r, _ := c.Login(context.Background(), &pb.LoginRequest{Username: "bot", Password: "bot"})
	assert.Equal(t, r.Status, "Success")

	r, _ = c.Login(context.Background(), &pb.LoginRequest{Username: "bot", Password: "gg"})
	assert.Equal(t, r.Status, "Fail")
	assert.Equal(t, r.Msg, "password is not correct")

	newuser := strconv.FormatInt(time.Now().UnixNano(), 10)
	r, _ = c.Login(context.Background(), &pb.LoginRequest{Username: newuser, Password: "bot"})
	assert.Equal(t, r.Status, "Fail")
	assert.Equal(t, r.Msg, "user does not exist")
}

func TestSignUp(t *testing.T) {
	newuser := strconv.FormatInt(time.Now().UnixNano(), 10)
	r, _ := c.SignUp(context.Background(), &pb.SignUpRequest{Username: newuser, Password: newuser, Profilename: newuser, Profileimg: template.URL("../assets/img/image.png")})
	assert.Equal(t, r.Status, "Success")
	time.Sleep(1 * time.Second)

	r, _ = c.SignUp(context.Background(), &pb.SignUpRequest{Username: newuser, Password: newuser, Profilename: newuser, Profileimg: template.URL("../assets/img/image.png")})
	assert.Equal(t, r.Status, "Fail")
	assert.Equal(t, r.Msg, "User already exists")

}

func TestCreatePost(t *testing.T) {
	newcontent := "Test: " + strconv.FormatInt(time.Now().UnixNano(), 10)
	r, _ := c.CreatePost(context.Background(), &pb.PostRequest{Username: "bot", Text: newcontent, Img: template.URL(""), Time: time.Now().Format("2006-01-02 15:04:05")})
	assert.Equal(t, r.Status, "Success")
}

func TestGetPosts(t *testing.T) {
	_, err := c.GetPosts(context.Background(), &pb.CommRequest{Username: "bot"})
	assert.Equal(t, err, nil)
}

func TestGetUserInfo(t *testing.T) {
	r, _ := c.GetUserInfo(context.Background(), &pb.CommRequest{Username: "bot"})
	assert.Equal(t, r.Username, "bot")
	assert.Equal(t, r.Profilename, "Andrew Hamilton")

	newuser := strconv.FormatInt(time.Now().UnixNano(), 10)
	_, err := c.GetUserInfo(context.Background(), &pb.CommRequest{Username: newuser})
	assert.Equal(t, err, errors.New("user not found"))
}

func TestFollow(t *testing.T) {
	r, _ := c.Follow(context.Background(), &pb.FollowRequest{Username1: "bot", Username2: "bot2"})
	assert.Equal(t, r.Status, "Success")
	time.Sleep(1 * time.Second)

	_, err := c.Follow(context.Background(), &pb.FollowRequest{Username1: "bot", Username2: "bot2"})
	assert.Equal(t, err, errors.New("/bot has already followed bot2"))
}

func TestUnfollow(t *testing.T) {
	r, _ := c.Unfollow(context.Background(), &pb.FollowRequest{Username1: "bot", Username2: "bot2"})
	assert.Equal(t, r.Status, "Success")
	time.Sleep(1 * time.Second)

	_, err := c.Unfollow(context.Background(), &pb.FollowRequest{Username1: "bot", Username2: "bot2"})
	assert.Equal(t, err, errors.New("/bot is not following bot2"))
}