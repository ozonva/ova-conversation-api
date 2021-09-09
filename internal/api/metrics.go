package api

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	promCreateCntr = promauto.NewCounter(prometheus.CounterOpts{Name: "conversation_created"})
	promUpdateCntr = promauto.NewCounter(prometheus.CounterOpts{Name: "conversation_updated"})
	promRemoveCntr = promauto.NewCounter(prometheus.CounterOpts{Name: "conversation_removed"})
)
