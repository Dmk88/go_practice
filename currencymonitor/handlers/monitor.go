package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Dmk88/go_practice/currencymonitor/models"
	"github.com/Dmk88/go_practice/currencymonitor/providers"
	uuid "github.com/satori/go.uuid"
)

func (h *Handler) Start(resp http.ResponseWriter, req *http.Request) {
	var err error
	period, err := time.ParseDuration(req.URL.Query().Get("period"))
	if err != nil {
		http.Error(resp, "bad period: "+err.Error(), http.StatusBadRequest)
		return
	}
	frequency, err := time.ParseDuration(req.URL.Query().Get("frequency"))
	if err != nil {
		http.Error(resp, "bad frequency: "+err.Error(), http.StatusBadRequest)
		return
	}

	var monitoring = models.Monitoring{
		ID:        uuid.NewV4().String(),
		Status:    models.MonitoringStatusActive,
		Period:    period.String(),
		Frequency: frequency.String(),
	}

	mp := providers.NewMonitoringProvider(h.scylla)
	err = mp.Save(&monitoring, h.config.Scylla.DefaultTTL)
	if err != nil {
		http.Error(resp, "add monitoring to scylladb: "+err.Error(), http.StatusBadRequest)
		return
	}
	go func() { h.daemon.ProcessMonitoring(monitoring) }()

	monitoringData, err := json.Marshal(monitoring)
	if err != nil {
		http.Error(resp, "json Marshal: "+err.Error(), http.StatusBadRequest)
		return
	}

	resp.Write(monitoringData)
}

func (h *Handler) MonitoringResults(resp http.ResponseWriter, req *http.Request) {
	var err error
	monitoringID := req.URL.Query().Get("monitoring_id")
	if monitoringID == "" {
		http.Error(resp, "bad monitoring_id", http.StatusBadRequest)
		return
	}
	mp := providers.NewMonitoringProvider(h.scylla)
	monitoring, err := mp.FindByID(monitoringID)
	if err != nil {
		http.Error(resp, "find monitoring in scylladb: "+err.Error(), http.StatusBadRequest)
		return
	}
	if monitoring.Status != models.MonitoringStatusCompleted {
		http.Error(resp, "monitoring is not completed", http.StatusBadRequest)
		return
	}

	resp.Write([]byte(monitoring.Data))
}
