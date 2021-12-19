package test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"main/client"
	"net/http"
	"testing"
)

func createUser(name string) map[string]string {
	return client.SignUp(map[string]string{"username": name, "password": "test_password", "profilename": "test_profilename"}, "")
}

func TestSearch(t *testing.T) {
	newUser := createUser("User1")
	response, err := client.GetUserInfo(newUser["username"])
	if err != nil || response == nil {
		t.Fatalf("Failed in TestSearch")
	}
	fmt.Printf("  ... Passed\n")
}

func TestSearchUnexistedUser(t *testing.T) {
	response, err := client.GetUserInfo("Unexisted")
	if err == nil || response != nil {
		t.Fatalf("Failed in TestSearchUnexistedUser")
	}
	fmt.Printf("  ... Passed\n")
}

func TestFollowOneUser(t *testing.T) {
	newUser2 := createUser("User2")
	newUser3 := createUser("User3")
	response, err := client.Follow(newUser2["username"], newUser3["username"])
	if err != nil || response == nil {
		t.Fatalf("Failed in TestFollowOneUser when follow")
	}
	resp, _ := http.Get("http://127.0.0.1:10379" + "/" + newUser3["username"])
	body, _ := ioutil.ReadAll(resp.Body)
	var followers []string
	json.Unmarshal(body, &followers)
	if len(followers) != 1 {
		t.Fatalf("Failed in TestFollowOneUser")
	}
	fmt.Printf("  ... Passed\n")
}

func TestFollowMultipleUser(t *testing.T) {
	newUser4 := createUser("User4")
	newUser5 := createUser("User5")
	newUser6 := createUser("User6")
	newUser7 := createUser("User7")
	response, err := client.Follow(newUser4["username"], newUser5["username"])
	if err != nil || response == nil {
		t.Fatalf("Failed in TestFollowMultipleUser when follow user 5")
	}

	response, err = client.Follow(newUser4["username"], newUser6["username"])
	if err != nil || response == nil {
		t.Fatalf("Failed in TestFollowMultipleUser when follow user 6")
	}

	response, err = client.Follow(newUser4["username"], newUser7["username"])
	if err != nil || response == nil {
		t.Fatalf("Failed in TestFollowMultipleUser when follow user 7")
	}

	resp, _ := http.Get("http://127.0.0.1:10379" + "/" + newUser4["username"])
	body, _ := ioutil.ReadAll(resp.Body)
	var followers []string
	json.Unmarshal(body, &followers)
	if len(followers) != 3 {
		t.Fatalf("Failed in TestFollowMultipleUser")
	}
	fmt.Printf("  ... Passed\n")
}

func TestFollowDuplicate(t *testing.T) {
	user4Info, err := client.GetUserInfo("User4")
	if err != nil || user4Info == nil {
		t.Fatalf("Failed in TestFollowDuplicate in get user4 info")
	}
	user5Info, err := client.GetUserInfo("User5")
	if err != nil || user5Info == nil {
		t.Fatalf("Failed in TestFollowDuplicate in get user5 info")
	}
	response, err := client.Follow(user4Info.Username, user5Info.Username)
	if err != nil || response != nil {
		t.Fatalf("Failed in TestFollowDuplicate")
	}
	fmt.Printf("  ... Passed\n")
}

func TestUnfollow(t *testing.T) {
	user2Info, err := client.GetUserInfo("User2")
	if err != nil || user2Info == nil {
		t.Fatalf("Failed in TestUnfollow in get user2 info")
	}
	user3Info, err := client.GetUserInfo("User3")
	if err != nil || user3Info == nil {
		t.Fatalf("Failed in TestUnfollow in get user3 info")
	}

	response, err := client.Unfollow(user2Info.Username, user3Info.Username)
	if err != nil || response == nil {
		t.Fatalf("Failed in TestUnfollow when unfollow")
	}

	resp, _ := http.Get("http://127.0.0.1:10379" + "/" + user3Info.Username)
	body, _ := ioutil.ReadAll(resp.Body)
	var followers []string
	json.Unmarshal(body, &followers)
	if len(followers) != 0 {
		t.Fatalf("Failed in TestFollowOneUser")
	}
	fmt.Printf("  ... Passed\n")
}

func TestUnfollowDuplicate(t *testing.T) {
	user2Info, err := client.GetUserInfo("User2")
	if err != nil || user2Info == nil {
		t.Fatalf("Failed in TestUnfollowDuplicate in get user2 info")
	}
	user3Info, err := client.GetUserInfo("User3")
	if err != nil || user3Info == nil {
		t.Fatalf("Failed in TestUnfollowDuplicate in get user3 info")
	}

	response, err := client.Unfollow(user2Info.Username, user3Info.Username)
	if err == nil || response != nil {
		t.Fatalf("Failed in TestUnfollowDuplicate when unfollow")
	}
	fmt.Printf("  ... Passed\n")
}
