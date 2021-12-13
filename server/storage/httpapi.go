package main

import (
	"io"
	"log"
	"net/http"
	"strconv"
	"encoding/json"
)

// Handler for a http based key-value store backed by raft

func (kv *Userkvstore) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	username := r.RequestURI
	defer r.Body.Close()
	switch r.Method {
		case http.MethodPut:
			user, err := io.ReadAll(r.Body)
			if err != nil {
				log.Printf("Failed to read on PUT (%v)\n", err)
				http.Error(w, "Failed on PUT", http.StatusBadRequest)
				return
			}
			var userinfo Userinfo
			json.Unmarshal(user, &userinfo)
			res := kv.Propose(username, userinfo)
			w.Write([]byte(res))

		case http.MethodGet:
			if user, ok := kv.Lookup(username); ok {
				userinfo, _ := json.Marshal(user)
				w.Write(userinfo)
			} else {
				http.Error(w, "Failed to GET", http.StatusNotFound)
			}

		default:
			w.Header().Set("Allow", http.MethodPut)
			w.Header().Add("Allow", http.MethodGet)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}


func (kv *Postkvstore) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	username := r.RequestURI
	defer r.Body.Close()
	switch r.Method {
	case http.MethodPut:
		post, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Failed to read on PUT (%v)\n", err)
			http.Error(w, "Failed on PUT", http.StatusBadRequest)
			return
		}
		var newpost Post
		json.Unmarshal(post, &newpost)
		res := kv.Propose(username, newpost)
		w.Write([]byte(res))

	case http.MethodGet:
		if posts, ok := kv.Lookup(username); ok {
			userposts, _ := json.Marshal(posts)
			w.Write(userposts)
		} else {
			http.Error(w, "Failed to GET", http.StatusNotFound)
		}

	default:
		w.Header().Set("Allow", http.MethodPut)
		w.Header().Add("Allow", http.MethodGet)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
/*

func (kv *Followkvstore) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := r.RequestURI
	defer r.Body.Close()
	switch r.Method {
	case http.MethodPut:
		user, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Failed to read on PUT (%v)\n", err)
			http.Error(w, "Failed on PUT", http.StatusBadRequest)
			return
		}
		var userinfo Userinfo
		json.Unmarshal(user, &userinfo)
		kv.Propose(key, userinfo)
		w.WriteHeader(http.StatusNoContent)

	case http.MethodGet:
		if user, ok := kv.Lookup(key); ok {
			userinfo, _ := json.Marshal(user)
			w.Write(userinfo)
		} else {
			http.Error(w, "Failed to GET", http.StatusNotFound)
		}

	default:
		w.Header().Set("Allow", http.MethodPut)
		w.Header().Add("Allow", http.MethodGet)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
*/
func serveHttpKVAPIUser(kv *Userkvstore, port int, errorC <-chan error) {
	srv := http.Server{
		Addr: ":" + strconv.Itoa(port),
		Handler: kv,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	if err, ok := <-errorC; ok {
		log.Fatal(err)
	}
}

func serveHttpKVAPIPost(kv *Postkvstore, port int, errorC <-chan error) {
	srv := http.Server{
		Addr: ":" + strconv.Itoa(port),
		Handler: kv,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	if err, ok := <-errorC; ok {
		log.Fatal(err)
	}
}

/*
func serveHttpKVAPIFollow(kv *Followkvstore, port int, errorC <-chan error) {
	srv := http.Server{
		Addr: ":" + strconv.Itoa(port),
		Handler: kv,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	if err, ok := <-errorC; ok {
		log.Fatal(err)
	}
}
*/