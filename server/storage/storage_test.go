package main

import (
	"testing"
	"io"
	"net/http"
	"net/http/httptest"
	"time"
	"bytes"
	"fmt"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"os"
)

func TestPutAndGetUserKeyValue(t *testing.T) {
	clusters := []string{"http://127.0.0.1:9021"}

	proposeC := make(chan string)
	defer close(proposeC)

	var userkvs *Userkvstore
	getSnapshot := func() ([]byte, error) { return userkvs.GetSnapshot() }
	commitC, errorC, snapshotterReady := newRaftNode(1, "test", clusters, getSnapshot, proposeC)
	userkvs = NewUserKVStore(<-snapshotterReady, proposeC, commitC, errorC)

	srv := httptest.NewServer(userkvs)
	defer srv.Close()

	// wait server started
	<-time.After(time.Second * 3)
	user := Userinfo {
		Password: "test-pwd",
		Profilename: "test-pn",
		Profileimg: "",
	}
	userinfo, _ := json.Marshal(user)
	wantKey, wantValue := "test-username", userinfo
	url := fmt.Sprintf("%s/%s", srv.URL, wantKey)
	cli := srv.Client()

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(userinfo))
	if err != nil {
		t.Fatal(err)
	}
	_, err = cli.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	// wait for a moment for processing message, otherwise get would be failed.
	<-time.After(time.Second)

	resp, err := cli.Get(url)
	if err != nil {
		t.Fatal(err)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	assert.Equal(t, data, wantValue)
	os.RemoveAll("data/testserver-1")
	os.RemoveAll("data/testserver-1-snap")
}

func TestPutAndGetPostKeyValue(t *testing.T) {
	clusters := []string{"http://127.0.0.1:9021"}

	proposeC := make(chan string)
	defer close(proposeC)

	var postkvs *Postkvstore
	getSnapshot := func() ([]byte, error) { return postkvs.GetSnapshot() }
	commitC, errorC, snapshotterReady := newRaftNode(1, "test", clusters, getSnapshot, proposeC)
	postkvs = NewPostKVStore(<-snapshotterReady, proposeC, commitC, errorC)

	srv := httptest.NewServer(postkvs)
	defer srv.Close()

	// wait server started
	<-time.After(time.Second * 3)
	post := Post {
		Text: "test-content",
		Time: time.Now().Format("2006-01-02 15:04:05"),
		Img:  "",
	}
	newpost, _ := json.Marshal(post)
	wantKey := "test-username"
	url := fmt.Sprintf("%s/%s", srv.URL, wantKey)
	cli := srv.Client()

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(newpost))
	if err != nil {
		t.Fatal(err)
	}
	_, err = cli.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	// wait for a moment for processing message, otherwise get would be failed.
	<-time.After(time.Second)

	resp, err := cli.Get(url)
	if err != nil {
		t.Fatal(err)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()		

	var posts []Post
	json.Unmarshal(data, &posts)

	assert.Equal(t, posts[0], post)
	os.RemoveAll("data/testserver-1")
	os.RemoveAll("data/testserver-1-snap")
}

func TestPutAndGetFollowKeyValue(t *testing.T) {
	clusters := []string{"http://127.0.0.1:9021"}

	proposeC := make(chan string)
	defer close(proposeC)

	var followkvs *Followkvstore		
	getSnapshot := func() ([]byte, error) { return followkvs.GetSnapshot() }
	commitC, errorC, snapshotterReady := newRaftNode(1, "test", clusters, getSnapshot, proposeC)
	followkvs = NewFollowKVStore(<-snapshotterReady, proposeC, commitC, errorC)

	srv := httptest.NewServer(followkvs)
	defer srv.Close()

	// wait server started
	<-time.After(time.Second * 3)
	
	wantKey, wantValue := "test-username", "test-username2"
	url := fmt.Sprintf("%s/%s", srv.URL, wantKey)
	cli := srv.Client()

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer([]byte(wantValue)))
	if err != nil {
		t.Fatal(err)
	}
	_, err = cli.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	// wait for a moment for processing message, otherwise get would be failed.
	<-time.After(time.Second)

	resp, err := cli.Get(url)
	if err != nil {
		t.Fatal(err)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	var users []string
	json.Unmarshal(data, &users)

	assert.Equal(t, users[0], wantValue)
	os.RemoveAll("data/testserver-1")
	os.RemoveAll("data/testserver-1-snap")
}