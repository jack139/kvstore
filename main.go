package main

/*
	external app 启动额外 tendermint

	./kvstore
	
	for TCP:
	tendermint node

	for UNIX socket:
	tendermint node --proxy_app=unix://example.sock
*/

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	//"kvstore/kv1"
	//"kvstore/kv2"
	"kvstore/kv3"

	"github.com/dgraph-io/badger"

	abciserver "github.com/tendermint/tendermint/abci/server"
	"github.com/tendermint/tendermint/libs/log"
)

var socketAddr string

func init() {
	flag.StringVar(&socketAddr, "socket-addr", "unix://example.sock", "Unix domain socket address")
}

func main() {
	fmt.Println("starting")
	db, err := badger.Open(badger.DefaultOptions("/tmp/badger"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open badger db: %v", err)
		os.Exit(1)
	}
	defer db.Close()

	//app := kv1.NewKVStoreApplication(db)
	//app := kv2.NewApplication()
	app := kv3.NewPersistentKVStoreApplication("./data")

	flag.Parse()

	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))

	//server := abciserver.NewSocketServer(socketAddr, app)
	server := abciserver.NewSocketServer(":26658", app)
	server.SetLogger(logger)
	if err := server.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "error starting socket server: %v", err)
		os.Exit(1)
	}
	defer server.Stop()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	os.Exit(0)
}
