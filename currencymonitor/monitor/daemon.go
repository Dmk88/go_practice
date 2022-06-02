package monitor

import (
	"encoding/json"
	"log"
	"time"

	"github.com/Dmk88/go_practice/currencymonitor/helpers"

	"github.com/Dmk88/go_practice/currencymonitor/models"
	"github.com/Dmk88/go_practice/currencymonitor/providers"
	"github.com/gocql/gocql"
)

type Daemon struct {
	scylla *gocql.Session
}

func NewDaemon(scylla *gocql.Session) Daemon {
	return Daemon{scylla: scylla}
}

func (d *Daemon) CheckMonitoring() {
	mp := providers.NewMonitoringProvider(d.scylla)
	monitoringRecords := mp.FindByStatus(models.MonitoringStatusActive)
	for _, monitoring := range monitoringRecords {
		d.ProcessMonitoring(monitoring)
	}
}

func (d *Daemon) ProcessMonitoring(monitoring models.Monitoring) {
	period, errPeriod := time.ParseDuration(monitoring.Period)
	if errPeriod != nil {
		log.Printf("monitoring ID: %s error parsing period Duration\n", monitoring.ID)
	}
	frequency, errFrequency := time.ParseDuration(monitoring.Frequency)
	if errFrequency != nil {
		log.Printf("monitoring ID: %s error parsing frequency Duration\n", monitoring.ID)
	}
	if errPeriod != nil || errFrequency != nil {
		d.updateMonitoringStatus(monitoring.ID, models.MonitoringStatusError)
		return
	}

	go d.process(monitoring.ID, period, frequency)
}

func (d *Daemon) process(monitoringID string, period, frequency time.Duration) {
	var respData string
	var resp models.APIResponse
	var err error

	result := make(map[string]float64)
	abort := make(chan struct{})
	go func() {
		time.AfterFunc(period, func() {
			abort <- struct{}{}
		})
	}()
	ticker := time.NewTicker(frequency)
	for alive := true; alive; {
		select {
		case <-abort:
			alive = false
			ticker.Stop()
			resultData, err := json.Marshal(result)
			if err != nil {
				d.updateMonitoringStatus(monitoringID, models.MonitoringStatusError)
				return
			}
			d.updateMonitoringData(monitoringID, string(resultData))
		case <-ticker.C:
			respData = genRsp()
			if respData == "" {
				d.updateMonitoringStatus(monitoringID, models.MonitoringStatusError)
				alive = false
				ticker.Stop()
			}
			resp, err = helpers.ParseAPIResponse(genRsp())
			if err != nil {
				d.updateMonitoringStatus(monitoringID, models.MonitoringStatusError)
				alive = false
				ticker.Stop()
			} else {
				result[time.Now().Format(time.RFC3339)] = resp.Amount
			}
		}
	}
}

func (d *Daemon) updateMonitoringStatus(monitoringID string, monitoringStatus int) {
	mp := providers.NewMonitoringProvider(d.scylla)
	err := mp.UpdateStatus(monitoringID, monitoringStatus)
	if err != nil {
		log.Printf("monitoring ID: %s error update in scylla: %s\n", monitoringID, err.Error())
	}
}

func (d *Daemon) updateMonitoringData(monitoringID string, monitoringData string) {
	mp := providers.NewMonitoringProvider(d.scylla)
	err := mp.UpdateData(monitoringID, monitoringData)
	if err != nil {
		d.updateMonitoringStatus(monitoringID, models.MonitoringStatusError)
		log.Printf("monitoring ID: %s error update in scylla: %s\n", monitoringID, err.Error())
	} else {
		d.updateMonitoringStatus(monitoringID, models.MonitoringStatusCompleted)
	}
}
