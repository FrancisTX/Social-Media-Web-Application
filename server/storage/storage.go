package main

import (
	"flag"
	"strings"
)

func main() {
	storage := flag.String("storage", "user", "type of storage")
	cluster := flag.String("cluster", "http://127.0.0.1:12379", "comma separated cluster peers")
	id := flag.Int("id", 1, "node ID")
	port := flag.Int("port", 12380, "key-value server port")
	flag.Parse()

	proposeC := make(chan string)
	defer close(proposeC)

	// raft provides a commit stream for the proposals from the http api
	if *storage == "user"{
		var userkvs *Userkvstore
		getSnapshot := func() ([]byte, error) { return userkvs.GetSnapshot() }
		commitC, errorC, snapshotterReady := newRaftNode(*id, *storage, strings.Split(*cluster, ","), getSnapshot, proposeC)
		userkvs = NewUserKVStore(<-snapshotterReady, proposeC, commitC, errorC)
		serveHttpKVAPIUser(userkvs, *port, errorC)
	} else if *storage == "post" {
		var postkvs *Postkvstore
		getSnapshot := func() ([]byte, error) { return postkvs.GetSnapshot() }
		commitC, errorC, snapshotterReady := newRaftNode(*id, *storage, strings.Split(*cluster, ","), getSnapshot, proposeC)
		postkvs = NewPostKVStore(<-snapshotterReady, proposeC, commitC, errorC)
		serveHttpKVAPIPost(postkvs, *port, errorC)
	} 
//	else if *storage == "follow" {
//		var followkvs *Followkvstore		
//		getSnapshot := func() ([]byte, error) { return followkvs.GetSnapshot() }
//		commitC, errorC, snapshotterReady := newRaftNode(*id, strings.Split(*cluster, ","), getSnapshot, proposeC)
//		followkvs = newKVStore(<-snapshotterReady, proposeC, commitC, errorC)
//		serveHttpKVAPIFollow(followkvs, *port, errorC)
//	}
}
