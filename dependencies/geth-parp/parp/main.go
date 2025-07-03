// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore
// +build ignore

package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/ethereum/go-ethereum/parp/handlers"
	"github.com/ethereum/go-ethereum/parp/manager"

	"github.com/ethereum/go-ethereum/parp/resmsg"
	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

type Client struct {
	ID   string
	Conn *websocket.Conn
}

// http.ResponseWriter is an interface, while Request is a concrete struct

func home(w http.ResponseWriter, r *http.Request) {

}

func main() {
	log.Println("Starting server...")

	// flag.Parse()
	// log.SetFlags(0)

	clientManager := manager.NewManager()

	log.Println("Address: ", clientManager.Address())

	http.HandleFunc("/ws/", handlers.HandleWebSocket(clientManager))
	// http.HandleFunc("/", handlers.HomeHandler())

	log.Fatal(http.ListenAndServe(":8888", nil))
	a := resmsg.HandshakeMsgBody{
		Type:             "gg",
		ServerPublicKeyB: []byte(""),
	}
	log.Println(a)

}
