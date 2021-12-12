package storage

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"log"
	"sync"

	"go.etcd.io/etcd/raft/v3/raftpb"
	"go.etcd.io/etcd/server/v3/etcdserver/api/snap"
)

// a key-value store backed by raft
type kvstore struct {
	proposeC    chan<- string // channel for proposing updates
	mu          sync.RWMutex
	kvStore     map[string]string // current committed key-value pairs
	snapshotter *snap.Snapshotter
}
 
type kv struct {
	Key string
	Val string
}

type Userinfo struct {
	Password    string
	Profilename string
	Profileimg  string
}

type Userkvstore struct {
	proposeC      chan<- string // channel for proposing updates
	mu            sync.RWMutex
	userkvStore   map[string]Userinfo // current committed key-value pairs
	snapshotter   *snap.Snapshotter
}

type userkv struct {
	Username    string
	Userinfo    Userinfo
}

type Post struct {
	Text string
	Time string
}

type Postkvstore struct {
	proposeC      chan<- string // channel for proposing updates
	mu            sync.RWMutex
	postkvStore   map[string][]Post // current committed key-value pairs
	snapshotter *snap.Snapshotter
}

type postkv struct {
	Username string
	Text string
	Time string
}

type Followkvstore struct {
	proposeC        chan<- string // channel for proposing updates
	mu              sync.RWMutex
	followkvStore   map[string][]string // current committed key-value pairs
	snapshotter *snap.Snapshotter
}

type followkv struct {
	Username1 string
	Username2 string
}

func NewUserKVStore(snapshotter *snap.Snapshotter, proposeC chan<- string, commitC <-chan *commit, errorC <-chan error) *Userkvstore {
	s := &Userkvstore{proposeC: proposeC, userkvStore: make(map[string]Userinfo), snapshotter: snapshotter}
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

func (s *Userkvstore) Lookup(key string) (Userinfo, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.userkvStore[key]
	return v, ok
}

func (s *Userkvstore) Propose(k string, v Userinfo) {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(userkv{k, v}); err != nil {
		log.Fatal(err)
	}
	s.proposeC <- buf.String()
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
			var dataKv userkv
			dec := gob.NewDecoder(bytes.NewBufferString(data))
			if err := dec.Decode(&dataKv); err != nil {
				log.Fatalf("raftexample: could not decode message (%v)", err)
			}
			s.mu.Lock()
			log.Printf("loading commit key %v and value %v", dataKv.Username, dataKv.Userinfo)
			s.userkvStore[dataKv.Username] = dataKv.Userinfo
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
	return json.Marshal(s.userkvStore)
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
	s.userkvStore = store
	return nil
}
