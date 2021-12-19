package storage

import (
	"errors"
	"net"
	"path/filepath"

	template "github.com/hashicorp/go-sockaddr/template"
	flag "github.com/ogier/pflag"
)

type Config struct {
	RaftAdd   net.Addr
	HTTPAdd   net.Addr
	JoinAdd   string
	StoreDir  string
	Bootstrap bool
}

func ParseConfig() (*Config, error) {
	//Get config from command line
	bindAddress := flag.StringP("address", "a", "127.0.0.1", "Bind IP address")
	raftport := flag.IntP("raft-port", "r", 7000, "Bind raft port")
	httpPort := flag.IntP("http-port", "h", 8000, "Bind HTTP port")
	dir := flag.StringP("store-dir", "d", "", "Raft data store directory")
	joinAddress := flag.String("join", "", "Address that other node can join")
	bootStrap := flag.Bool("bootstrap", false, "if the first node in the cluster")
	flag.Parse()

	//Address
	var bindIP net.IP
	parsedAddress, err := template.Parse(*bindAddress)
	if err != nil {
		return nil, err
	}

	bindIP = net.ParseIP(parsedAddress)
	if bindIP == nil {
		return nil, err
	}

	// Raft and HTTP port
	if *raftport < 1 || *raftport > 65536 || *httpPort < 1 || *httpPort > 65536 {
		return nil, errors.New("port invalid")
	}
	raftAddr := &net.TCPAddr{IP: bindIP, Port: *raftport}
	httpAddr := &net.TCPAddr{IP: bindIP, Port: *httpPort}

	// Data directory
	dataDir, err := filepath.Abs(*dir)
	if err != nil {
		return nil, errors.New("store directory error")
	}

	return &Config{
		RaftAdd:   raftAddr,
		HTTPAdd:   httpAddr,
		JoinAdd:   *joinAddress,
		StoreDir:  dataDir,
		Bootstrap: *bootStrap,
	}, nil
}
