package main

import (
	"errors"
	"log"
	"main/raft/storage"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"
)

func main() {
	//parse the config
	config, err := storage.ParseConfig()
	if err != nil {
		log.Println("Error resolve config:", err)
		os.Exit(1)
	}
	log.Println("Config: ", config)

	//initialize the node
	node, err := storage.InitStore(config)
	if err != nil {
		log.Println("Error configuring node:", err)
		os.Exit(1)
	}

	//join into cluster if join address valid
	if config.JoinAdd != "" {
		go func() {
			for {
				err := join(config)
				if err != nil {
					log.Println(err)
					time.Sleep(2 * time.Second)
				} else {
					break
				}
			}
		}()
	}

	httpServer := &storage.Server{
		Node: node,
		Addr: config.HTTPAdd,
	}

	//start the http server
	err = httpServer.Start()
	if err != nil {
		log.Println("Http server error")
	}
	log.Println("Start raft")

	//exit if keyboard interrupt
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt)
	<-exit
	log.Println("Exit raft")
}

func join(config *storage.Config) error {
	//make up a join post request
	url := url.URL{Scheme: "http", Host: config.JoinAdd, Path: "join"}
	req, err := http.NewRequest(http.MethodPost, url.String(), nil)
	if err != nil {
		log.Println("Http request err: ", err)
		return err
	}
	req.Header.Add("Peer", config.RaftAdd.String())

	//request to join
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Do request err: ", err)
		return err
	}

	if resp.StatusCode != http.StatusOK {
		log.Println("Error code", resp.StatusCode)
		return errors.New("Return Error")
	}
	return nil
}
