package main

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"log"
	"sync"
	"html/template"

	"go.etcd.io/etcd/raft/v3/raftpb"
	"go.etcd.io/etcd/server/v3/etcdserver/api/snap"
)

type Userinfo struct {
	Password    string
	Profilename string
	Profileimg  template.URL
}

type Userkvstore struct {
	proposeC      chan<- string
	mu            sync.RWMutex
	UserkvStore   map[string]Userinfo
	snapshotter   *snap.Snapshotter
}

type Userkv struct {
	Username    string
	Userinfo    Userinfo
}


/////////////////////////////////////////
/////////////////////////////////////////
// User
/////////////////////////////////////////
/////////////////////////////////////////


func NewUserKVStore(snapshotter *snap.Snapshotter, proposeC chan<- string, commitC <-chan *commit, errorC <-chan error) *Userkvstore {
	s := &Userkvstore{proposeC: proposeC, UserkvStore: make(map[string]Userinfo), snapshotter: snapshotter}
	snapshot, err := s.loadSnapshot()
	if err != nil {
		log.Panic(err)
	}
	if snapshot != nil {
		log.Printf("loading snapshot at term %d and index %d", snapshot.Metadata.Term, snapshot.Metadata.Index)
		if err := s.recoverFromSnapshot(snapshot.Data); err != nil {
			log.Panic(err)
		}
	}
	// read commits from raft into kvStore map until error
	go s.readCommits(commitC, errorC)
	return s
}

func (s *Userkvstore) Lookup(username string) (Userinfo, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	user, ok := s.UserkvStore[username]
	return user, ok
}

func (s *Userkvstore) Propose(username string, userinfo Userinfo) string {
	if _, ok := s.Lookup(username); ok {
		return "User already exists"
	} else {
		var buf bytes.Buffer
		if err := gob.NewEncoder(&buf).Encode(Userkv{username, userinfo}); err != nil {
			log.Fatal(err)
		}
		s.proposeC <- buf.String()
		return ""
	}
}

func (s *Userkvstore) readCommits(commitC <-chan *commit, errorC <-chan error) {
	for commit := range commitC {
		if commit == nil {
			// signaled to load snapshot
			snapshot, err := s.loadSnapshot()
			if err != nil {
				log.Panic(err)
			}
			if snapshot != nil {
				log.Printf("loading snapshot at term %d and index %d", snapshot.Metadata.Term, snapshot.Metadata.Index)
				if err := s.recoverFromSnapshot(snapshot.Data); err != nil {
					log.Panic(err)
				}
			}
			continue
		}

		for _, data := range commit.data {
			var dataKv Userkv
			dec := gob.NewDecoder(bytes.NewBufferString(data))
			if err := dec.Decode(&dataKv); err != nil {
				log.Fatalf("raftexample: could not decode message (%v)", err)
			}
			s.mu.Lock()
			s.UserkvStore[dataKv.Username] = dataKv.Userinfo
			s.mu.Unlock()
		}
		close(commit.applyDoneC)
	}
	if err, ok := <-errorC; ok {
		log.Fatal(err)
	}
}

func (s *Userkvstore) GetSnapshot() ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return json.Marshal(s.UserkvStore)
}


func (s *Userkvstore) loadSnapshot() (*raftpb.Snapshot, error) {
	snapshot, err := s.snapshotter.Load()
	if err == snap.ErrNoSnapshot {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return snapshot, nil
}

func (s *Userkvstore) recoverFromSnapshot(snapshot []byte) error {
	var store map[string]Userinfo
	if err := json.Unmarshal(snapshot, &store); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.UserkvStore = store
	return nil
}

/////////////////////////////////////////
/////////////////////////////////////////
// Post 
/////////////////////////////////////////
/////////////////////////////////////////


type Post struct {
	Text string
	Time string
	Img  template.URL
}

type Postkvstore struct {
	proposeC      chan<- string 
	mu            sync.RWMutex
	PostkvStore   map[string][]Post
	snapshotter *snap.Snapshotter
}

type postkv struct {
	Username string
	Post     Post
}


func NewPostKVStore(snapshotter *snap.Snapshotter, proposeC chan<- string, commitC <-chan *commit, errorC <-chan error) *Postkvstore {
	s := &Postkvstore{proposeC: proposeC, PostkvStore: make(map[string][]Post), snapshotter: snapshotter}
	snapshot, err := s.loadSnapshot()
	if err != nil {
		log.Panic(err)
	}
	if snapshot != nil {
		log.Printf("loading snapshot at term %d and index %d", snapshot.Metadata.Term, snapshot.Metadata.Index)
		if err := s.recoverFromSnapshot(snapshot.Data); err != nil {
			log.Panic(err)
		}
	}
	// read commits from raft into kvStore map until error
	go s.readCommits(commitC, errorC)
	return s
}

func (s *Postkvstore) Lookup(username string) ([]Post, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	post, ok := s.PostkvStore[username]
	return post, ok
}

func (s *Postkvstore) Propose(username string, post Post) string {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(postkv{username, post}); err != nil {
		log.Fatal(err)
	}
	s.proposeC <- buf.String()
	return ""
}

func (s *Postkvstore) readCommits(commitC <-chan *commit, errorC <-chan error) {
	for commit := range commitC {
		if commit == nil {
			snapshot, err := s.loadSnapshot()
			if err != nil {
				log.Panic(err)
			}
			if snapshot != nil {
				log.Printf("loading snapshot at term %d and index %d", snapshot.Metadata.Term, snapshot.Metadata.Index)
				if err := s.recoverFromSnapshot(snapshot.Data); err != nil {
					log.Panic(err)
				}
			}
			continue
		}

		for _, data := range commit.data {
			var dataKv postkv
			dec := gob.NewDecoder(bytes.NewBufferString(data))
			if err := dec.Decode(&dataKv); err != nil {
				log.Fatalf("raftexample: could not decode message (%v)", err)
			}
			s.mu.Lock()
			s.PostkvStore[dataKv.Username] = append(s.PostkvStore[dataKv.Username], dataKv.Post)
			s.mu.Unlock()
		}
		close(commit.applyDoneC)
	}
	if err, ok := <-errorC; ok {
		log.Fatal(err)
	}
}

func (s *Postkvstore) GetSnapshot() ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return json.Marshal(s.PostkvStore)
}


func (s *Postkvstore) loadSnapshot() (*raftpb.Snapshot, error) {
	snapshot, err := s.snapshotter.Load()
	if err == snap.ErrNoSnapshot {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return snapshot, nil
}

func (s *Postkvstore) recoverFromSnapshot(snapshot []byte) error {
	var store map[string][]Post
	if err := json.Unmarshal(snapshot, &store); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.PostkvStore = store
	return nil
}


/////////////////////////////////////////
/////////////////////////////////////////
// Follow
/////////////////////////////////////////
/////////////////////////////////////////


type Followkvstore struct {
	proposeC        chan<- string // channel for proposing updates
	mu              sync.RWMutex
	FollowkvStore   map[string][]string // current committed key-value pairs
	snapshotter *snap.Snapshotter
}

type followkv struct {
	Username string
	Username2 string
	Event     string
}


func NewFollowKVStore(snapshotter *snap.Snapshotter, proposeC chan<- string, commitC <-chan *commit, errorC <-chan error) *Followkvstore {
	s := &Followkvstore{proposeC: proposeC, FollowkvStore: make(map[string][]string), snapshotter: snapshotter}
	snapshot, err := s.loadSnapshot()
	if err != nil {
		log.Panic(err)
	}
	if snapshot != nil {
		log.Printf("loading snapshot at term %d and index %d", snapshot.Metadata.Term, snapshot.Metadata.Index)
		if err := s.recoverFromSnapshot(snapshot.Data); err != nil {
			log.Panic(err)
		}
	}
	// read commits from raft into kvStore map until error
	go s.readCommits(commitC, errorC)
	return s
}

func (s *Followkvstore) Lookup(username string) ([]string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	followedusers, ok := s.FollowkvStore[username]
	return followedusers, ok
}

func (s *Followkvstore) LookupFollowing(username string, username2 string) int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if followedusers, ok := s.FollowkvStore[username]; ok {
		for idx, user := range followedusers {
	        if (user == username2) {
	            return idx
	        }
	    }   
	}
	return -1
}

func (s *Followkvstore) Propose(username string, username2 string, event string) string {
	if event == "follow" {
		if s.LookupFollowing(username, username2) != -1 {
			return username + " has already followed " + username2
		}
	} else {
		if s.LookupFollowing(username, username2) == -1 {
			return username + " is not following " + username2
		}
	}
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(followkv{username, username2, event}); err != nil {
		log.Fatal(err)
	}
	s.proposeC <- buf.String()
	return ""
}

func (s *Followkvstore) readCommits(commitC <-chan *commit, errorC <-chan error) {
	for commit := range commitC {
		if commit == nil {
			snapshot, err := s.loadSnapshot()
			if err != nil {
				log.Panic(err)
			}
			if snapshot != nil {
				log.Printf("loading snapshot at term %d and index %d", snapshot.Metadata.Term, snapshot.Metadata.Index)
				if err := s.recoverFromSnapshot(snapshot.Data); err != nil {
					log.Panic(err)
				}
			}
			continue
		}

		for _, data := range commit.data {
			var dataKv followkv
			dec := gob.NewDecoder(bytes.NewBufferString(data))
			if err := dec.Decode(&dataKv); err != nil {
				log.Fatalf("raftexample: could not decode message (%v)", err)
			}
			s.mu.Lock()
			if dataKv.Event == "follow" {
				s.FollowkvStore[dataKv.Username] = append(s.FollowkvStore[dataKv.Username], dataKv.Username2)
			} else {
				for idx, user := range s.FollowkvStore[dataKv.Username] {
	        		if (user == dataKv.Username2) {
	        			s.FollowkvStore[dataKv.Username][idx] = s.FollowkvStore[dataKv.Username][len(s.FollowkvStore[dataKv.Username])-1]
	        			s.FollowkvStore[dataKv.Username][len(s.FollowkvStore[dataKv.Username])-1] = ""
	        			s.FollowkvStore[dataKv.Username] = s.FollowkvStore[dataKv.Username][:len(s.FollowkvStore[dataKv.Username])-1]
			            break
			        }
			    }   
			}
			s.mu.Unlock()
		}
		close(commit.applyDoneC)
	}
	if err, ok := <-errorC; ok {
		log.Fatal(err)
	}
}

func (s *Followkvstore) GetSnapshot() ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return json.Marshal(s.FollowkvStore)
}


func (s *Followkvstore) loadSnapshot() (*raftpb.Snapshot, error) {
	snapshot, err := s.snapshotter.Load()
	if err == snap.ErrNoSnapshot {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return snapshot, nil
}

func (s *Followkvstore) recoverFromSnapshot(snapshot []byte) error {
	var store map[string][]string
	if err := json.Unmarshal(snapshot, &store); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.FollowkvStore = store
	return nil
}
