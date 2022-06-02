package handlers

import (
	"github.com/Dmk88/go_practice/currencymonitor/monitor"
	"github.com/gocql/gocql"
)

type Handler struct {
	config monitor.Config
	scylla *gocql.Session
	daemon monitor.Daemon
}

func NewHandler(
	config monitor.Config,
	scylla *gocql.Session,
	daemon monitor.Daemon,
) *Handler {
	return &Handler{
		config: config,
		scylla: scylla,
		daemon: daemon,
	}
}
