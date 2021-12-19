
package main

import (
	"reflect"
	"testing"
	"github.com/stretchr/testify/assert"
	"time"
)

func TestUserKVStore(t *testing.T) {
	info := Userinfo{Password: "testpwd", Profilename: "testpn"}
	user := map[string]Userinfo{"test": info}
	s := &Userkvstore{UserkvStore: user}

	v, _ := s.Lookup("test")
	assert.Equal(t, v, user["test"])

	data, err := s.GetSnapshot()
	assert.Equal(t, err, nil)

	s.UserkvStore = nil

	err = s.recoverFromSnapshot(data)
	assert.Equal(t, err, nil)

	v, _ = s.Lookup("test")
	assert.Equal(t, v, user["test"])
	if !reflect.DeepEqual(s.UserkvStore, user) {
		t.Fatalf("store expected %+v, got %+v", user, s.UserkvStore)
	}
}

func TestPostKVStore(t *testing.T) {
	post := Post{Text: "testing", Time: time.Now().Format("2006-01-02 15:04:05")}
	var posts []Post
	posts = append(posts, post)
	testpost := map[string][]Post{"test": posts}
	s := &Postkvstore{PostkvStore: testpost}

	v, _ := s.Lookup("test")
	assert.Equal(t, v, testpost["test"])

	data, err := s.GetSnapshot()
	assert.Equal(t, err, nil)

	s.PostkvStore = nil

	err = s.recoverFromSnapshot(data)
	assert.Equal(t, err, nil)

	v, _ = s.Lookup("test")
	assert.Equal(t, v, testpost["test"])
	if !reflect.DeepEqual(s.PostkvStore, testpost) {
		t.Fatalf("store expected %+v, got %+v", testpost, s.PostkvStore)
	}
}

func TestFollowKVStore(t *testing.T) {
	follow := map[string][]string{"test": []string{"test2"}}
	s := &Followkvstore{FollowkvStore: follow}

	v, _ := s.Lookup("test")
	assert.Equal(t, v, follow["test"])

	data, err := s.GetSnapshot()
	assert.Equal(t, err, nil)

	s.FollowkvStore = nil

	err = s.recoverFromSnapshot(data)
	assert.Equal(t, err, nil)

	v, _ = s.Lookup("test")
	assert.Equal(t, v, follow["test"])
	if !reflect.DeepEqual(s.FollowkvStore, follow) {
		t.Fatalf("store expected %+v, got %+v", follow, s.FollowkvStore)
	}
}