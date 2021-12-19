package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
	"strconv"
	"github.com/stretchr/testify/assert"
)


var	newUser string
var	newUser2 string = createUser(strconv.FormatInt(time.Now().UnixNano(), 10)+"2")
var	newUser3 string = createUser(strconv.FormatInt(time.Now().UnixNano(), 10)+"3")
var	newUser4 string = createUser(strconv.FormatInt(time.Now().UnixNano(), 10)+"4")
var	newUser5 string = createUser(strconv.FormatInt(time.Now().UnixNano(), 10)+"5")
var	newUser6 string = createUser(strconv.FormatInt(time.Now().UnixNano(), 10)+"6")
var	newUser7 string = createUser(strconv.FormatInt(time.Now().UnixNano(), 10)+"7")


func createUser(name string) string {
	SignUp(map[string]string{"username": name, "password": "test_password", "profilename": "test_profilename"}, "")
	return name
}

func TestSignUp(t *testing.T) {
	newUser = strconv.FormatInt(time.Now().UnixNano(), 10)+"1"
	r := SignUp(map[string]string{"username": newUser, "password": "test_password", "profilename": "test_profilename"}, "")
	assert.Equal(t, r["status"], "Success")

	time.Sleep(1 * time.Second)
	r = SignUp(map[string]string{"username": newUser, "password": "test_password", "profilename": "test_profilename"}, "")
	assert.Equal(t, r["status"], "Fail")
	assert.Contains(t, r["msg"], "User already exists")

}

func TestLogin(t *testing.T) {
	r, _ := Login(map[string]string{"username": newUser, "password": "test_password"})
	assert.Equal(t, r["status"], "Success")

	r, _ = Login(map[string]string{"username": newUser, "password": "password"})
	assert.Equal(t, r["status"], "Fail")
	assert.Contains(t, r["msg"], "password is not correct")

	r, _ = Login(map[string]string{"username": strconv.FormatInt(time.Now().UnixNano(), 10), "password": "password"})
	assert.Equal(t, r["status"], "Fail")
	assert.Contains(t, r["msg"], "user does not exist")
}

func TestCreateGetPost(t *testing.T) {
	newcontent := "Test: " + strconv.FormatInt(time.Now().UnixNano(), 10)
	r := CreatePost(map[string]string{"username": newUser, "text": newcontent, "time": time.Now().Format("2006-01-02 15:04:05")}, "")
	assert.Equal(t, r["status"], "Success")

	time.Sleep(1 * time.Second)
	posts, err := GetPosts(map[string]string{"username": newUser})
	assert.Equal(t, err, nil)
	assert.Equal(t, posts[0].Text, newcontent)
}

func TestSearch(t *testing.T) {
	time.Sleep(1 * time.Second)
	
	response, err := GetUserInfo(newUser)
	if err != nil || response == nil {
		t.Fatalf("Failed in TestSearch, %v", err)
	}
	fmt.Printf("  ... Passed\n")
}

func TestSearchUnexistedUser(t *testing.T) {
	response, err := GetUserInfo("Unexisted")
	if err == nil || response != nil {
		t.Fatalf("Failed in TestSearchUnexistedUser")
	}
	fmt.Printf("  ... Passed\n")
}

func TestFollowOneUser(t *testing.T) {
	response, err := Follow(newUser2, newUser3)
	if err != nil || response == nil {
		t.Fatalf("Failed in TestFollowOneUser when follow")
	}
	time.Sleep(2 * time.Second)
	
	resp, _ := http.Get("http://127.0.0.1:12380" + "/" + newUser2)
	body, _ := ioutil.ReadAll(resp.Body)
	var followers []string
	json.Unmarshal(body, &followers)
	if len(followers) != 1 {
		t.Fatalf("Failed in TestFollowOneUser %v", followers)
	}
	fmt.Printf("  ... Passed\n")
}

func TestFollowMultipleUser(t *testing.T) {
	response, err := Follow(newUser4, newUser5)
	if err != nil || response == nil {
		t.Fatalf("Failed in TestFollowMultipleUser when follow user 5")
	}

	response, err = Follow(newUser4, newUser6)
	if err != nil || response == nil {
		t.Fatalf("Failed in TestFollowMultipleUser when follow user 6")
	}

	response, err = Follow(newUser4, newUser7)
	if err != nil || response == nil {
		t.Fatalf("Failed in TestFollowMultipleUser when follow user 7")
	}

	time.Sleep(2 * time.Second)

	resp, _ := http.Get("http://127.0.0.1:12380" + "/" + newUser4)
	body, _ := ioutil.ReadAll(resp.Body)
	var followers []string
	json.Unmarshal(body, &followers)
	if len(followers) != 3 {
		t.Fatalf("Failed in TestFollowMultipleUser %v", followers)
	}
	fmt.Printf("  ... Passed\n")
}

func TestFollowDuplicate(t *testing.T) {
	response, err := Follow(newUser4, newUser5)
	if err == nil || response != nil {
		t.Fatalf("Failed in TestFollowDuplicate")
	}
	fmt.Printf("  ... Passed\n")
}

func TestUnfollow(t *testing.T) {
	response, err := Unfollow(newUser2, newUser3)
	if err != nil || response == nil {
		t.Fatalf("Failed in TestUnfollow when unfollow %v", err)
	}

	time.Sleep(2 * time.Second)

	resp, _ := http.Get("http://127.0.0.1:12380" + "/" + newUser3)
	body, _ := ioutil.ReadAll(resp.Body)
	var followers []string
	json.Unmarshal(body, &followers)
	if len(followers) != 0 {
		t.Fatalf("Failed in TestFollowOneUser")
	}
	fmt.Printf("  ... Passed\n")
}

func TestUnfollowDuplicate(t *testing.T) {
	response, err := Unfollow(newUser2, newUser3)
	if err == nil || response != nil {
		t.Fatalf("Failed in TestUnfollowDuplicate when unfollow")
	}
	fmt.Printf("  ... Passed\n")
}