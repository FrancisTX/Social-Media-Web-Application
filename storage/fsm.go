package storage

import (
	"encoding/json"
	"io"
	"log"
	"sync"

	"github.com/hashicorp/raft"
)

type fsm struct {
	mutex      sync.Mutex
	stateValue string
}

type request struct {
	operation string
	value     string
}

func (fsm *fsm) Apply(logEntry *raft.Log) interface{} {
	var r request
	if err := json.Unmarshal(logEntry.Data, &r); err != nil {
		log.Panicln("[FSM Apply] Failed unmarshaling log entry")
	}

	switch r.operation {
	case "SET":
		fsm.mutex.Lock()
		defer fsm.mutex.Unlock()
		fsm.stateValue = r.value

		return nil
	default:
		log.Panicln("Unrecognized operation: %s.", r.operation)
	}
}

func (fsm *fsm) Snapshot() (raft.FSMSnapshot, error) {
	fsm.mutex.Lock()
	defer fsm.mutex.Unlock()

	return &fsmSnapshot{stateValue: fsm.stateValue}, nil
}

func (fsm *fsm) Restore(serialized io.ReadCloser) error {
	var snapshot fsmSnapshot
	if err := json.NewDecoder(serialized).Decode(&snapshot); err != nil {
		return err
	}

	fsm.stateValue = snapshot.stateValue
	return nil
}

type fsmSnapshot struct {
	stateValue string `json:"value"`
}

func (f *fsmSnapshot) Persist(sink raft.SnapshotSink) error {
	err := func() error {
		snapshotBytes, err := json.Marshal(f)
		if err != nil {
			return err
		}

		if _, err := sink.Write(snapshotBytes); err != nil {
			return err
		}

		if err := sink.Close(); err != nil {
			return err
		}

		return nil
	}()

	if err != nil {
		sink.Cancel()
		return err
	}

	return nil
}

func (f *fsmSnapshot) Release() {}
