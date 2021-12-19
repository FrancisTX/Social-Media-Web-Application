package storage

import (
	"encoding/json"
	"io"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/hashicorp/raft"
)

type Server struct {
	Addr     net.Addr
	Node     *Store
	listener net.Listener
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.Addr.String())
	if err != nil {
		return err
	}
	server := http.Server{Handler: s}
	s.listener = listener
	http.Handle("/", s)

	go func() {
		err := server.Serve(s.listener)
		if err != nil {
			log.Fatalf("HTTP start error: %s", err)
		}
	}()
	return nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/key") {
		s.keyRequestHandler(w, r)
	} else if r.URL.Path == "/join" {
		s.joinHandler(w, r)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func (server *Server) keyRequestHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		request := map[string]string{}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		for key, val := range request {
			err := server.Node.Set(key, val)
			if err != nil {
				log.Println("Node set Error: ", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		w.WriteHeader(http.StatusOK)
	case http.MethodGet:
		//get key
		parameters := strings.Split(r.URL.Path, "/")
		if len(parameters) != 3 {
			log.Println("Error Parameter")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		key := parameters[2]

		//get
		val, err := server.Node.Get(key)
		if err != nil {
			log.Println("Node get Error:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		response, err := json.Marshal(map[string]string{key: val})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		io.WriteString(w, string(response))
	case http.MethodDelete:
		//get key
		parameters := strings.Split(r.URL.Path, "/")
		if len(parameters) != 3 {
			log.Println("Error Parameter")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		key := parameters[2]
		if key == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		//delete
		err := server.Node.Delete(key)
		if err != nil {
			log.Println("Node delete Error:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
	return
}

func (server *Server) joinHandler(w http.ResponseWriter, r *http.Request) {
	//get the peer address in the request
	clusterAddress := r.Header.Get("Peer")
	if clusterAddress == "" {
		log.Println("Get nil peer address")
		w.WriteHeader(http.StatusBadRequest)
	}

	//add as a voter in the cluster
	addVoter := server.Node.raft.AddVoter(raft.ServerID(clusterAddress), raft.ServerAddress(clusterAddress), 0, 0)
	if addVoter.Error() != nil {
		log.Println("Error joining into cluster")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Println("Join into cluster")
	w.WriteHeader(http.StatusOK)
}
