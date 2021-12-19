package storage

import (
	"encoding/json"
	"io"
	"log"
	"sync"

	"github.com/hashicorp/raft"
)

type fsm struct {
	mutex    sync.Mutex
	userInfo map[string]string
	follow   map[string]string
	post     map[string]string
}

type request struct {
	Operation string
	Key       string
	Value     string
}

//fsm implementation
func (fsm *fsm) Apply(logEntry *raft.Log) interface{} {
	var r request
	if err := json.Unmarshal(logEntry.Data, &r); err != nil {
		log.Println("Failed unmarshaling request")
		return nil
	}

	//apply the operation in fsm
	switch r.Operation {
	case "set":
		fsm.mutex.Lock()
		defer fsm.mutex.Unlock()
		fsm.userInfo[r.Key] = r.Value
		return nil
	case "delete":
		fsm.mutex.Lock()
		defer fsm.mutex.Unlock()
		delete(fsm.userInfo, r.Key)
		return nil
	default:
		log.Println("Unrecognized operation:", r.Operation)
	}
	return nil
}

func (fsm *fsm) Snapshot() (raft.FSMSnapshot, error) {
	fsm.mutex.Lock()
	defer fsm.mutex.Unlock()
	//copy the store
	userinfo := make(map[string]string)
	for k, v := range fsm.userInfo {
		userinfo[k] = v
	}
	return &fsmSnapshot{UserInfo: userinfo}, nil

}

func (fsm *fsm) Restore(serialize io.ReadCloser) error {
	userinfo := make(map[string]string)
	if err := json.NewDecoder(serialize).Decode(&userinfo); err != nil {
		return err
	}
	//restore
	fsm.userInfo = userinfo
	return nil
}

//fsmSnapshot implementation
type fsmSnapshot struct {
	UserInfo map[string]string `json:"user_info"`
}

func (f *fsmSnapshot) Persist(sink raft.SnapshotSink) error {

	snapshot, err := json.Marshal(f)
	if err != nil {
		sink.Cancel()
		return err
	}

	//write snapshot
	_, err = sink.Write(snapshot)
	if err != nil {
		sink.Cancel()
		return err
	}

	err = sink.Close()
	if err != nil {
		sink.Cancel()
		return err
	}

	return nil
}

func (f *fsmSnapshot) Release() {}
