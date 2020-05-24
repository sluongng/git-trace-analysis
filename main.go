package main

import (
	"bufio"
	"context"
	"encoding/json"
	"net"
	"os"

	"go.uber.org/zap"
)

const (
	SockAddr = "/tmp/git_trace.sock"
)

type RunContext struct {
	context.Context
	log *zap.Logger
}

func NewRunContext(log *zap.Logger) *RunContext {
	return &RunContext{context.Background(), log}
}

func (ctx *RunContext) handleEvent(b []byte) {
	event := &GitTrace2Event{}
	if err := json.Unmarshal(b, &event); err != nil {
		ctx.log.Error("unable to unmarshal event", zap.Error(err))
	}

	ctx.log.Info("event", zap.String("type", string(event.Event)))
}

func main() {
	log, _ := zap.NewDevelopment()
	defer log.Sync()

	if err := os.RemoveAll(SockAddr); err != nil {
		log.Fatal("failed to remove socket", zap.Error(err))
	}

	uAddr, err := net.ResolveUnixAddr("unixgram", SockAddr)
	if err != nil {
		log.Fatal(
			"could not resolve unix socket address",
			zap.String("socketAddress", SockAddr),
			zap.Error(err),
		)
	}

	conn, err := net.ListenUnixgram("unixgram", uAddr)
	if err != nil {
		log.Fatal(
			"could not listen on unix socket datagram",
			zap.String("socketAddress", SockAddr),
			zap.Error(err),
		)
	}
	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	for {
		if ok := scanner.Scan(); !ok {
			log.Error("scan not ok")
			continue
		}

		ctx := NewRunContext(log)

		b := scanner.Bytes()
		go ctx.handleEvent(b)
	}
}
