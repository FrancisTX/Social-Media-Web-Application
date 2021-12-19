package storage

import (
	"encoding/json"
	"net"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb"
)

//A node(peer) in raft
type Store struct {
	RaftDir  string
	RaftBind string

	mu       sync.Mutex
	raft     *raft.Raft
	fsm      fsm
	userInfo map[string]string
	post     map[string]string
	follow   map[string]string
}

func InitStore(config *Config) (*Store, error) {
	if err := os.MkdirAll(config.StoreDir, 0700); err != nil {
		return nil, err
	}

	raftConfig := raft.DefaultConfig()

	raftConfig.LocalID = raft.ServerID(config.RaftAdd.String())

	address, err := net.ResolveTCPAddr("tcp", config.RaftAdd.String())
	if err != nil {
		return nil, err
	}
	transport, err := raft.NewTCPTransport(address.String(), address, 3, 10*time.Second, os.Stdout)
	if err != nil {
		return nil, err
	}

	logBolt, err := raftboltdb.NewBoltStore(filepath.Join(config.StoreDir, "log_store.bolt"))
	if err != nil {
		return nil, err
	}

	stableBolt, err := raftboltdb.NewBoltStore(filepath.Join(config.StoreDir, "stable_store.bolt"))
	if err != nil {
		return nil, err
	}

	newFsm := fsm{
		userInfo: make(map[string]string),
		post:     make(map[string]string),
		follow:   make(map[string]string),
	}

	snapshot, err := raft.NewFileSnapshotStore(config.StoreDir, 1, os.Stdout)
	if err != nil {
		return nil, err
	}

	raftNode, err := raft.NewRaft(raftConfig, &newFsm, logBolt, stableBolt, snapshot, transport)
	if err != nil {
		return nil, err
	}

	if config.Bootstrap {
		bootstrapConfig := raft.Configuration{
			Servers: []raft.Server{{ID: raftConfig.LocalID, Address: transport.LocalAddr()}},
		}
		raftNode.BootstrapCluster(bootstrapConfig)
	}

	return &Store{
		RaftDir:  config.StoreDir,
		RaftBind: config.RaftAdd.String(),
		raft:     raftNode,
		fsm:      newFsm,
	}, nil
}

func (s *Store) Get(key string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.fsm.userInfo[key], nil
}

func (s *Store) Set(key, value string) error {
	c := &request{
		Operation: "set",
		Key:       key,
		Value:     value,
	}
	b, err := json.Marshal(c)
	if err != nil {
		return err
	}

	f := s.raft.Apply(b, 5*time.Second)
	return f.Error()
}

func (s *Store) Delete(key string) error {
	c := &request{
		Operation: "delete",
		Key:       key,
	}
	b, err := json.Marshal(c)
	if err != nil {
		return err
	}

	f := s.raft.Apply(b, 5*time.Second)
	return f.Error()
}
