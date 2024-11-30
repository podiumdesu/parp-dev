// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package rpc

import (
	"context"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
)

const MetadataApi = "rpc"
const EngineApi = "engine"

// CodecOption specifies which type of messages a codec supports.
//
// Deprecated: this option is no longer honored by Server.
type CodecOption int

const (
	// OptionMethodInvocation is an indication that the codec supports RPC method calls
	OptionMethodInvocation CodecOption = 1 << iota

	// OptionSubscriptions is an indication that the codec supports RPC notifications
	OptionSubscriptions = 1 << iota // support pub sub
)

// Server is an RPC server.
type Server struct {
	services serviceRegistry
	idgen    func() ID

	mutex              sync.Mutex
	codecs             map[ServerCodec]struct{}
	run                atomic.Bool
	batchItemLimit     int
	batchResponseLimit int
	httpBodyLimit      int
}

// NewServer creates a new server instance with no registered handlers.
func NewServer() *Server {
	server := &Server{
		idgen:         randomIDGenerator(),
		codecs:        make(map[ServerCodec]struct{}),
		httpBodyLimit: defaultBodyLimit,
	}
	server.run.Store(true)
	// Register the default service providing meta information about the RPC service such
	// as the services and methods it offers.
	rpcService := &RPCService{server}
	server.RegisterName(MetadataApi, rpcService)
	return server
}

// SetBatchLimits sets limits applied to batch requests. There are two limits: 'itemLimit'
// is the maximum number of items in a batch. 'maxResponseSize' is the maximum number of
// response bytes across all requests in a batch.
//
// This method should be called before processing any requests via ServeCodec, ServeHTTP,
// ServeListener etc.
func (s *Server) SetBatchLimits(itemLimit, maxResponseSize int) {
	s.batchItemLimit = itemLimit
	s.batchResponseLimit = maxResponseSize
}

// SetHTTPBodyLimit sets the size limit for HTTP requests.
//
// This method should be called before processing any requests via ServeHTTP.
func (s *Server) SetHTTPBodyLimit(limit int) {
	s.httpBodyLimit = limit
}

// RegisterName creates a service for the given receiver type under the given name. When no
// methods on the given receiver match the criteria to be either a RPC method or a
// subscription an error is returned. Otherwise a new service is created and added to the
// service collection this server provides to clients.
func (s *Server) RegisterName(name string, receiver interface{}) error {
	return s.services.registerName(name, receiver)
}

// ServeCodec reads incoming requests from codec, calls the appropriate callback and writes
// the response back using the given codec. It will block until the codec is closed or the
// server is stopped. In either case the codec is closed.
//
// Note that codec options are no longer supported.
func (s *Server) ServeCodec(codec ServerCodec, options CodecOption) {
	defer codec.close()

	if !s.trackCodec(codec) {
		return
	}
	defer s.untrackCodec(codec)

	cfg := &clientConfig{
		idgen:              s.idgen,
		batchItemLimit:     s.batchItemLimit,
		batchResponseLimit: s.batchResponseLimit,
	}
	c := initClient(codec, &s.services, cfg)
	<-codec.closed()
	c.Close()
}

func (s *Server) trackCodec(codec ServerCodec) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if !s.run.Load() {
		return false // Don't serve if server is stopped.
	}
	s.codecs[codec] = struct{}{}
	return true
}

func (s *Server) untrackCodec(codec ServerCodec) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	delete(s.codecs, codec)
}

// serveSingleRequest reads and processes a single RPC request from the given codec. This
// is used to serve HTTP connections. Subscriptions and reverse calls are not allowed in
// this mode.
func (s *Server) serveSingleRequest(ctx context.Context, codec ServerCodec) {
	// Don't serve if server is stopped.
	if !s.run.Load() {
		return
	}

	h := newHandler(ctx, codec, s.idgen, &s.services, s.batchItemLimit, s.batchResponseLimit)
	h.allowSubscribe = false
	defer h.close(io.EOF, nil)

	reqs, batch, err := codec.readBatch()
	if err != nil {
		if err != io.EOF {
			resp := errorMessage(&invalidMessageError{"parse error"})
			codec.writeJSON(ctx, resp, true)
		}
		return
	}

	// Extract and log metadata
	// peerInfo := PeerInfoFromContext(ctx)
	// for _, req := range reqs {
	// 	// Convert params to a readable JSON string
	// 	paramsStr, err := json.Marshal(req.Params)
	// 	if err != nil {
	// 		log.Error("Failed to marshal RPC params", "err", err)
	// 		continue
	// 	}

	// 	log.Info("Received RPC request",
	// 		"method", req.Method,
	// 		"params", string(paramsStr),
	// 		"from", peerInfo.RemoteAddr)

	// 	// Save to file
	// 	logToFile(peerInfo, req.Method, string(paramsStr))
	// }

	if batch {
		h.handleBatch(reqs)
	} else {
		h.handleMsg(reqs[0])
	}
}

const (
	colorRed   = "\033[31m"
	colorReset = "\033[0m"
)

func decodeRawTransaction(rawTxHex string) (*types.Transaction, error) {
	// Remove the "0x" prefix if it exists
	log.Info("Decoding raw transaction-------")
	log.Info(rawTxHex)
	if len(rawTxHex) > 2 && rawTxHex[:2] == "0x" {
		rawTxHex = rawTxHex[2:]
	}

	// Decode the hex string into bytes
	rawTxBytes, err := hex.DecodeString(rawTxHex)
	if err != nil {
		return nil, err
	}

	// Create a new empty Transaction object
	var tx types.Transaction

	// Unmarshal the raw bytes into the Transaction object
	if err := tx.UnmarshalBinary(rawTxBytes); err != nil {
		return nil, err
	}

	return &tx, nil
}

func logToFile(peerInfo PeerInfo, method, params string) {
	f, err := os.OpenFile("rpc_requests.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Error("Failed to open log file", "err", err)
		return
	}
	defer f.Close()

	const (
		colorRed   = "\033[31m"
		colorReset = "\033[0m"
	)

	methodConsole := method
	methodFile := method
	if method == "eth_getBalance" {
		methodConsole = fmt.Sprintf("%s%s%s", colorRed, method, colorReset)
		methodFile = fmt.Sprintf("*** %s ***", method)
	}

	if method == "eth_sendRawTransaction" {
		log.Info("This is params")

		cleanedParams := strings.Trim(params, `[]"`)
		tx, err := decodeRawTransaction(cleanedParams)

		if err != nil {
			fmt.Printf("Failed to decode transaction: %v", err)
		}

		// Determine the signer
		signer := types.LatestSignerForChainID(tx.ChainId())

		// Recover the sender address
		from, err := types.Sender(signer, tx)
		if err != nil {
			log.Error("Failed to recover sender: %v", err)
		}

		// Print out the transaction details
		fmt.Println("Transaction details:")
		fmt.Printf("Nonce: %d\n", tx.Nonce())
		fmt.Printf("Gas Price: %s\n", tx.GasPrice().String())
		fmt.Printf("Gas Limit: %d\n", tx.Gas())
		fmt.Printf("To: %s\n", tx.To().Hex())
		fmt.Printf("From: %s\n", from.Hex())
		fmt.Printf("Value: %s\n", tx.Value().String())
		fmt.Printf("Data: %x\n", tx.Data())
	}

	// Log entry for the file (with brief formatting)
	logEntryFile := fmt.Sprintf(
		"Received RPC request from %-15s; Method: %-30s Params: %-100s\n",
		peerInfo.RemoteAddr, methodFile, params,
	)

	// Log entry for the console (with color codes)
	logEntryConsole := fmt.Sprintf(
		"Received RPC request from %-15s; Method: %-30s Params: %-100s\n",
		peerInfo.RemoteAddr, methodConsole, params,
	)

	// Write the log entry to the file (with brief formatting)
	if _, err := f.WriteString(logEntryFile); err != nil {
		log.Error("Failed to write to log file", "err", err)
	}

	// Print the log entry to the console (with color if applicable)
	fmt.Print(logEntryConsole)
}

// Stop stops reading new requests, waits for stopPendingRequestTimeout to allow pending
// requests to finish, then closes all codecs which will cancel pending requests and
// subscriptions.
func (s *Server) Stop() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.run.CompareAndSwap(true, false) {
		log.Debug("RPC server shutting down")
		for codec := range s.codecs {
			codec.close()
		}
	}
}

// RPCService gives meta information about the server.
// e.g. gives information about the loaded modules.
type RPCService struct {
	server *Server
}

// Modules returns the list of RPC services with their version number
func (s *RPCService) Modules() map[string]string {
	s.server.services.mu.Lock()
	defer s.server.services.mu.Unlock()

	modules := make(map[string]string)
	for name := range s.server.services.services {
		modules[name] = "1.0"
	}
	return modules
}

// PeerInfo contains information about the remote end of the network connection.
//
// This is available within RPC method handlers through the context. Call
// PeerInfoFromContext to get information about the client connection related to
// the current method call.
type PeerInfo struct {
	// Transport is name of the protocol used by the client.
	// This can be "http", "ws" or "ipc".
	Transport string

	// Address of client. This will usually contain the IP address and port.
	RemoteAddr string

	// Additional information for HTTP and WebSocket connections.
	HTTP struct {
		// Protocol version, i.e. "HTTP/1.1". This is not set for WebSocket.
		Version string
		// Header values sent by the client.
		UserAgent string
		Origin    string
		Host      string
	}
}

type peerInfoContextKey struct{}

// PeerInfoFromContext returns information about the client's network connection.
// Use this with the context passed to RPC method handler functions.
//
// The zero value is returned if no connection info is present in ctx.
func PeerInfoFromContext(ctx context.Context) PeerInfo {
	info, _ := ctx.Value(peerInfoContextKey{}).(PeerInfo)
	return info
}
