package models

const (
	MonitoringStatusActive    = 1
	MonitoringStatusError     = 2
	MonitoringStatusCompleted = 3
)

type Monitoring struct {
	ID        string `json:"id"`
	Status    int    `json:"-"`
	Period    string `json:"-"`
	Frequency string `json:"-"`
	Data      string `json:"-"`
}
